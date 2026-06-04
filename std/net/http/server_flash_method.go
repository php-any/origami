package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

// ServerFlashMethod 启动时扫描注解路由并直接注册到 Server
// 支持多个独立应用：扫描目录下所有 main.php，每个 main.php 的 #[Application] 独立扫描各自的控制器目录
// 用法: $server->flash("./src") 或 $server->flash("./apps")
type ServerFlashMethod struct {
	server *ServerClass
}

func (h *ServerFlashMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 解析 dir 参数
	dir, acl := resolveAppDirPath(ctx, 0, "./src")
	if acl != nil {
		return nil, acl
	}

	vm := ctx.GetVM()

	// 扫描目录下所有 main.php 文件（支持多应用）
	mainFiles, err := findMainFiles(dir)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	if len(mainFiles) == 0 {
		return nil, utils.NewThrowf("flash: 目录 %s 中未找到 main.php", dir)
	}

	// 加载每个 main.php，触发各自的 #[Application] 注解扫描
	for _, mainFile := range mainFiles {
		if _, acl := vm.LoadAndRun(mainFile); acl != nil {
			return nil, acl
		}
	}

	// 从 VM 获取所有已注册的注解路由
	routes := runtime.HTTPRoutes(vm)
	if len(routes) == 0 {
		return nil, utils.NewThrowf("flash: 目录 %s 未发现注解路由，请确认控制器带有 @Controller/@*Mapping 注解", dir)
	}

	// 将每条路由直接注册到 Server 的 ServeMux
	for _, rt := range routes {
		rt := rt
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request := NewRequestClassFrom(r)
			response := NewResponseWriterClassFrom(w)

			reqProxy := data.NewProxyValue(request, ctx)
			resProxy := data.NewProxyValue(response, ctx)

			// 执行控制器级拦截器 preHandle
			if !executePreHandle(vm, ctx, rt.Middlewares, reqProxy, resProxy) {
				return
			}

			// 执行控制器方法
			var lastACL data.Control
			executeControllerMethod(rt, reqProxy, resProxy, ctx, &lastACL)

			// 执行控制器级拦截器 postHandle
			executePostHandle(vm, ctx, rt.Middlewares, reqProxy, resProxy)

			// 执行控制器级拦截器 afterCompletion
			executeAfterCompletion(vm, ctx, rt.Middlewares, reqProxy, resProxy)
		})

		// 包装服务器级中间件
		var final http.Handler = handler
		if len(h.server.Middlewares) > 0 {
			final = applyMiddlewares(final, h.server.Middlewares)
		}

		// 使用 Go 1.22+ 方法路由语法注册
		methodPath := strings.ToUpper(rt.Method) + " " + rt.Path
		h.server.source.Handle(methodPath, final)
	}

	// 构建返回的路由信息数组，供 PHP 侧使用（如打印启动日志）
	routeList := make([]data.Value, 0, len(routes))
	for _, rt := range routes {
		obj := data.NewObjectValue()
		obj.SetProperty("method", data.NewStringValue(rt.Method))
		obj.SetProperty("path", data.NewStringValue(rt.Path))
		routeList = append(routeList, obj)
	}

	return data.NewArrayValue(routeList), nil
}

// findMainFiles 递归扫描目录下所有 main.php 文件
func findMainFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "main.php" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (h *ServerFlashMethod) GetName() string            { return "flash" }
func (h *ServerFlashMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerFlashMethod) GetIsStatic() bool          { return false }
func (h *ServerFlashMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "dir", 0, data.NewStringValue("./src"), data.String{}),
	}
}
func (h *ServerFlashMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "dir", 0, nil),
	}
}
func (h *ServerFlashMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
