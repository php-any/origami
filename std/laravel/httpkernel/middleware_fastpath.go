package httpkernel

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
)

// 全局中间件里「本次请求逻辑上等价于 return $next($request)」的条目，直接在 Go 侧判定并跳过。
//
// 为什么值得做：这些中间件每次请求都要走一遍容器解析 + 若干次 PHP 调用，换来的结论却是
// 「什么都不做」。实测（examples/laravel13 /hello，c=32，单请求 14.03ms CPU）：
// HandleCors +1.61ms、PreventRequestsDuringMaintenance +1.34ms、TrimStrings +0.28ms、
// ConvertEmptyStringsToNull +0.47ms、ValidatePostSize +0.20ms、InvokeDeferredCallbacks +0.19ms。
// 对无 query、无请求体的 GET，9 个全局中间件里有 6 个被跳过，单请求降到 10.36ms（−26%），
// QPS 1024 → 1371（+34%）。
//
// 判定纪律：**只有能证明「必然透传」才跳过**，证明不了的一律保留原路径。
// 请求无关的部分（配置形态、pattern 列表、down 文件路径）只在进程内算一次；
// 依赖请求的部分一律用与中间件**同一个** PHP 调用（$request->is(...)）来判，不自己重写
// pattern 匹配，避免任何语义漂移。
//
// 跳过只影响 Pipeline 里跑哪些中间件；kernelTerminate 仍按完整的 s.middleware 遍历 terminate，
// 所以像 InvokeDeferredCallbacks 这种 handle 透明但 terminate 有副作用的条目依然会正常收尾。
const (
	mwInvokeDeferredCallbacks          = "Illuminate\\Foundation\\Http\\Middleware\\InvokeDeferredCallbacks"
	mwHandleCors                       = "Illuminate\\Http\\Middleware\\HandleCors"
	mwPreventRequestsDuringMaintenance = "Illuminate\\Foundation\\Http\\Middleware\\PreventRequestsDuringMaintenance"
	mwValidatePostSize                 = "Illuminate\\Http\\Middleware\\ValidatePostSize"
	mwTrimStrings                      = "Illuminate\\Foundation\\Http\\Middleware\\TrimStrings"
	mwConvertEmptyStringsToNull        = "Illuminate\\Foundation\\Http\\Middleware\\ConvertEmptyStringsToNull"
)

// middlewareBypass 缓存「哪些全局中间件在本进程的配置下会透明」这一请求无关的判定，
// 由常驻 kernelState 持有、请求沙箱共享同一个指针（见 request_policy.go 的 NativeClone）。
// 字段在 once.Do 完成后只读，并发请求无锁可读。
type middlewareBypass struct {
	once sync.Once

	// corsKnown 为真表示 config('cors.paths') 是一组可判定的普通字符串 pattern
	// （见 prover 注释）。HandleCors::handle 透传 ⟺ 该列表一个都不匹配。
	corsKnown bool
	corsPaths []string

	// maintenanceKnown / maintenanceExcluded 对应 PreventRequestsDuringMaintenance
	// 的排除表；maintenanceDownFile 非空表示维护模式用的是文件驱动，
	// FileBasedMaintenanceMode::active() 就是 file_exists(该路径)。
	maintenanceKnown    bool
	maintenanceExcluded []string
	maintenanceDownFile string
}

// apply 去掉本次请求上必然透传的全局中间件。没有可去掉的条目时原样返回入参切片（零分配）。
func (b *middlewareBypass) apply(ctx data.Context, s *kernelState, request data.Value, list []string) []string {
	if b == nil || len(list) == 0 {
		return list
	}
	b.once.Do(func() { b.resolve(ctx, s) })
	req := goHTTPRequest(ctx)
	for i, name := range list {
		if !b.transparent(name, req, ctx, request) {
			continue
		}
		out := make([]string, 0, len(list))
		out = append(out, list[:i]...)
		for _, rest := range list[i+1:] {
			if !b.transparent(rest, req, ctx, request) {
				out = append(out, rest)
			}
		}
		return out
	}
	return list
}

// transparent 判断单个中间件在本次请求上是否必然等价于 return $next($request)。
func (b *middlewareBypass) transparent(name string, req *http.Request, ctx data.Context, request data.Value) bool {
	switch name {
	case mwInvokeDeferredCallbacks:
		// 方法体就是 `return $next($request);`，无条件透传。
		// terminate() 另有副作用，但 terminate 遍历的是完整的 s.middleware，与此无关。
		return true
	case mwHandleCors:
		// hasMatchingPath() 之外的两条出口（skipCallbacks 循环、paths 不匹配）都是透传，
		// 所以「一个 pattern 都不匹配」就等价于整个 handle 透传。
		return b.corsKnown && !patternMatches(ctx, request, b.corsPaths)
	case mwPreventRequestsDuringMaintenance:
		// 排除表命中或 down 文件不存在时，handle() 都是 `return $next($request)`。
		// stat 出错（权限等）一律保留中间件走原路径。
		if !b.maintenanceKnown || b.maintenanceDownFile == "" {
			return false
		}
		if patternMatches(ctx, request, b.maintenanceExcluded) {
			return false
		}
		_, err := os.Stat(b.maintenanceDownFile)
		return err != nil && os.IsNotExist(err)
	case mwValidatePostSize:
		// `if ($max > 0 && $request->server('CONTENT_LENGTH') > $max)`：
		// 无请求体时 CONTENT_LENGTH 要么缺失（null）、要么为 0，两个分支都取不到真值，
		// 无论 post_max_size 是多少都到不了 throw。
		return req != nil && req.ContentLength <= 0
	case mwTrimStrings, mwConvertEmptyStringsToNull:
		// 二者都只改写「请求输入」，作用域是 query bag + request bag（以及 JSON body）：
		//   TransformsRequest::clean → cleanParameterBag($bag) → $bag->replace(cleanArray($bag->all()))
		// bag 为空时 all() 是 []，cleanArray([]) 还是 []，replace([]) 不改任何东西。
		// 所以「没有 query string、也没有请求体」时二者必然等价于 return $next($request)。
		//
		// ContentLength 严格等于 0 才算无体：chunked 请求里它是 -1，此时有体但长度未知，
		// 一律保留中间件走原路径。
		return req != nil && req.ContentLength == 0 && req.URL.RawQuery == ""
	}
	return false
}

// patternMatches 用中间件同款调用判「本次请求路径是否命中这些 pattern」。
// 传入空列表直接短路，省掉一次 PHP 调用。
func patternMatches(ctx data.Context, request data.Value, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}
	if request == nil {
		// 取不到请求就按「命中」处理，让调用方保留中间件。
		return true
	}
	return requestMatches(ctx, request, patterns)
}

// goHTTPRequest 取当前请求的原生 *http.Request（PHP 侧那个 Request 就是从它构造的）。
func goHTTPRequest(ctx data.Context) *http.Request {
	if ctx == nil {
		return nil
	}
	host, ok := ctx.GetVM().(interface{ HTTPRequest() *http.Request })
	if !ok || host == nil {
		return nil
	}
	return host.HTTPRequest()
}

// resolve 只在进程内执行一次。任何一步取不到或形态不可判定就保持零值，
// 判定退回「全部保留中间件」的原路径，只损失快路径而不改变语义。
func (b *middlewareBypass) resolve(ctx data.Context, s *kernelState) {
	if ctx == nil || s == nil || s.app == nil {
		return
	}
	configObj := asValue(requestInstance(ctx, s.app, "config"))
	if configObj == nil {
		return
	}

	if ret, ctl := callObjectMethodInContext(ctx, configObj, "get",
		data.NewStringValue("cors.paths"), data.NewArrayValue(nil)); ctl == nil {
		b.corsPaths, b.corsKnown = provablePaths(asValue(ret))
	}

	// 维护模式：驱动必须是 file（MaintenanceModeManager::getDefaultDriver 读的是同一个 key）。
	driver, ctl := callObjectMethodInContext(ctx, configObj, "get",
		data.NewStringValue("app.maintenance.driver"), data.NewStringValue("file"))
	if ctl != nil || valueAsString(asValue(driver)) != "file" {
		return
	}
	mw, ctl := callObjectMethodInContext(ctx, s.app, "make", data.NewStringValue(mwPreventRequestsDuringMaintenance))
	if ctl != nil || mw == nil {
		return
	}
	excluded, ctl := callObjectMethodInContext(ctx, asValue(mw), "getExcludedPaths")
	if ctl != nil {
		return
	}
	b.maintenanceExcluded, b.maintenanceKnown = provablePaths(asValue(excluded))
	if !b.maintenanceKnown {
		return
	}
	down, ctl := callObjectMethodInContext(ctx, s.app, "storagePath", data.NewStringValue("framework/down"))
	if ctl != nil {
		return
	}
	b.maintenanceDownFile = valueAsString(asValue(down))
}

// provablePaths 判断一个配置数组是不是「可判定的普通字符串 pattern 列表」。
//
// 需要排除的两类形态：
//   - 字符串键（host 键控，如 ['api.example.com' => ['api/*']]）：Laravel 走
//     `$paths[$host] ?? array_filter($paths, is_string)`，结果依赖请求 Host，Go 侧不猜。
//   - 可能命中 fullUrl 的 pattern：中间件的判定是 `fullUrlIs($path) || is($path)`，
//     而 Go 侧只复核前者为假的那些（见 cantMatchFullURL），否则会漏判。
func provablePaths(v data.Value) ([]string, bool) {
	arr, ok := v.(*data.ArrayValue)
	if !ok || arr == nil {
		return nil, false
	}
	out := make([]string, 0, arr.Len())
	for i := 0; i < arr.Len(); i++ {
		z := arr.At(i)
		if z == nil || z.EmptyStrKey || !listKey(z.Name) {
			return nil, false
		}
		s, ok := z.ReadValue().(*data.StringValue)
		if !ok {
			return nil, false
		}
		// 与 ExcludesPaths::inExceptArray / HandleCors::hasMatchingPath 一致地 trim。
		path := s.Value
		if path != "/" {
			path = strings.Trim(path, "/")
		}
		if !cantMatchFullURL(path) {
			return nil, false
		}
		out = append(out, path)
	}
	return out, true
}

// listKey 判断槽位键是不是「列表位置键」：空键（数组字面量）或十进制整数键
// （array_merge / array_values 等会补成 "0"、"1"…）。字符串键一律视为 host 键控。
func listKey(name string) bool {
	if name == "" {
		return true
	}
	for i := 0; i < len(name); i++ {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}
	return true
}

// cantMatchFullURL 证明 Str::is(pattern, $request->fullUrl()) 必然为假。
//
// fullUrl 一定以 "http://" 或 "https://" 开头。pattern 含 "://" 时可能在字面上等于
// fullUrl（放弃）；不含 * 时正则整体锚定，只有完全相等才命中，而含 "://" 的 fullUrl
// 不可能等于不含 "://" 的 pattern（必然为假）；含 * 时把第一个 * 之前的字面前缀记作
// head，命中要求 fullUrl 以 head 开头，故 head 是 "http://"/"https://" 的前缀时无法
// 证明（放弃），head 为空（pattern 以 * 开头，含 "*" 本身）同样放弃。
func cantMatchFullURL(pattern string) bool {
	if strings.Contains(pattern, "://") {
		return false
	}
	head, _, hasStar := strings.Cut(pattern, "*")
	if !hasStar {
		return true
	}
	if head == "" {
		return false
	}
	return !strings.HasPrefix("http://", head) && !strings.HasPrefix("https://", head)
}
