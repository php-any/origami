package cookie

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const cookieJarClassName = "Illuminate\\Cookie\\CookieJar"
const symfonyCookie = "Symfony\\Component\\HttpFoundation\\Cookie"

type jarState struct {
	path     string
	domain   data.Value
	secure   data.Value
	sameSite string
	queued   map[string]map[string]*data.ClassValue // name -> path -> cookie
}

var (
	jarStates sync.Map
	jarNextID atomic.Int64
)

func jarOf(cv *data.ClassValue) *jarState {
	if cv == nil {
		return &jarState{path: "/", sameSite: "lax", queued: map[string]map[string]*data.ClassValue{}}
	}
	if v, ctl := cv.GetProperty("__origami_jar_id"); ctl == nil && v != nil {
		if iv, ok := kit.Unwrap(v).(*data.IntValue); ok && iv.Value > 0 {
			if s, ok := jarStates.Load(int64(iv.Value)); ok {
				return s.(*jarState)
			}
		}
	}
	id := jarNextID.Add(1)
	s := &jarState{path: "/", sameSite: "lax", queued: map[string]map[string]*data.ClassValue{}}
	jarStates.Store(id, s)
	_ = cv.SetProperty("__origami_jar_id", data.NewIntValue(int(id)))
	return s
}

type CookieJarClass struct {
	node.Node
	methods map[string]data.Method
}

func NewCookieJarClass() data.ClassStmt {
	c := &CookieJarClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *CookieJarClass) GetName() string    { return cookieJarClassName }
func (c *CookieJarClass) GetExtend() *string { return nil }
func (c *CookieJarClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Cookie\\QueueingFactory"}
}
func (c *CookieJarClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "path", "domain", "secure", "sameSite", "queued", "__origami_jar_id":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *CookieJarClass) GetPropertyList() []data.Property {
	out := make([]data.Property, 0, 6)
	for _, n := range []string{"path", "domain", "secure", "sameSite", "queued", "__origami_jar_id"} {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *CookieJarClass) GetConstruct() data.Method                { return nil }
func (c *CookieJarClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	st := jarOf(cv)
	_ = cv.SetProperty("path", data.NewStringValue(st.path))
	_ = cv.SetProperty("sameSite", data.NewStringValue(st.sameSite))
	return cv, nil
}
func (c *CookieJarClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *CookieJarClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *CookieJarClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *CookieJarClass) register() {
	c.methods["make"] = kit.InstanceMethodOpt("make", []string{"name", "value", "minutes", "path", "domain", "secure", "httpOnly", "raw", "sameSite"}, 2, jarMake)
	c.methods["queue"] = kit.InstanceMethod("queue", []string{"parameters"}, jarQueue)
	c.methods["getqueuedcookies"] = kit.InstanceMethod("getQueuedCookies", nil, jarGetQueued)
	c.methods["forget"] = kit.InstanceMethodOpt("forget", []string{"name", "path", "domain"}, 1, jarForget)
	c.methods["forever"] = kit.InstanceMethodOpt("forever", []string{"name", "value", "path", "domain", "secure", "httpOnly", "raw", "sameSite"}, 2, jarForever)
	c.methods["hasqueued"] = kit.InstanceMethodOpt("hasQueued", []string{"key", "path"}, 1, jarHasQueued)
	c.methods["queued"] = kit.InstanceMethodOpt("queued", []string{"key", "default", "path"}, 1, jarQueued)
	c.methods["expire"] = kit.InstanceMethodOpt("expire", []string{"name", "path", "domain"}, 1, jarExpire)
	c.methods["unqueue"] = kit.InstanceMethodOpt("unqueue", []string{"name", "path"}, 1, jarUnqueue)
	c.methods["setdefaultpathanddomain"] = kit.InstanceMethodOpt("setDefaultPathAndDomain", []string{"path", "domain", "secure", "sameSite"}, 2, jarSetDefaultPathAndDomain)
	c.methods["flushqueuedcookies"] = kit.InstanceMethod("flushQueuedCookies", nil, jarFlushQueued)
	kit.RegisterMacroable(c.methods, cookieJarClassName)
}

func jarRecv(ctx data.Context) (*data.ClassValue, *jarState, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return nil, nil, data.NewErrorThrow(nil, fmt.Errorf("CookieJar missing $this"))
	}
	return cv, jarOf(cv), nil
}

func jarMake(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return jarMakeWith(cv, st, ctx,
		kit.Arg(ctx, 0), kit.Arg(ctx, 1), kit.Arg(ctx, 2),
		kit.Arg(ctx, 3), kit.Arg(ctx, 4), kit.Arg(ctx, 5),
		kit.Arg(ctx, 6), kit.Arg(ctx, 7), kit.Arg(ctx, 8),
	)
}

func jarMakeWith(cv *data.ClassValue, st *jarState, ctx data.Context, nameV, value, minutesV, pathV, domainV, secureV, httpOnlyV, rawV, sameSiteV data.Value) (data.GetValue, data.Control) {
	name := ""
	if nameV != nil {
		name = nameV.AsString()
	}
	minutes := 0
	if minutesV != nil {
		if iv, ok := minutesV.(data.AsInt); ok {
			minutes, _ = iv.AsInt()
		}
	}
	path := st.path
	if pathV != nil && !kit.IsNull(pathV) {
		path = pathV.AsString()
	}
	domain := domainV
	if domain == nil {
		domain = st.domain
	}
	secure := secureV
	if secure == nil {
		secure = st.secure
	}
	httpOnly := data.NewBoolValue(true)
	if httpOnlyV != nil {
		httpOnly = httpOnlyV
	}
	raw := data.NewBoolValue(false)
	if rawV != nil {
		raw = rawV
	}
	sameSite := st.sameSite
	if sameSiteV != nil && !kit.IsNull(sameSiteV) {
		sameSite = sameSiteV.AsString()
	}
	var expire data.Value = data.NewIntValue(0)
	if minutes != 0 {
		expire = data.NewIntValue(int(time.Now().Unix()) + minutes*60)
	}
	return newSymfonyCookie(ctx, cv, name, value, expire, data.NewStringValue(path), domain, secure, httpOnly, raw, data.NewStringValue(sameSite))
}

func jarForever(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return jarMakeWith(cv, st, ctx,
		kit.Arg(ctx, 0), kit.Arg(ctx, 1), data.NewIntValue(576000),
		kit.Arg(ctx, 2), kit.Arg(ctx, 3), kit.Arg(ctx, 4),
		kit.Arg(ctx, 5), kit.Arg(ctx, 6), kit.Arg(ctx, 7),
	)
}

func jarForget(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return jarMakeWith(cv, st, ctx,
		kit.Arg(ctx, 0), data.NewNullValue(), data.NewIntValue(-2628000),
		kit.Arg(ctx, 1), kit.Arg(ctx, 2), nil, nil, nil, nil,
	)
}

func jarQueue(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	args := ctx.GetFlatCallArgs()
	var cookie *data.ClassValue
	if len(args) > 0 {
		if c, ok := kit.Unwrap(args[0]).(*data.ClassValue); ok && c != nil && c.Class.GetName() == symfonyCookie {
			cookie = c
		}
	}
	if cookie == nil {
		ret, ctl := jarMake(ctx)
		if ctl != nil {
			return nil, ctl
		}
		var ok bool
		cookie, ok = ret.(*data.ClassValue)
		if !ok {
			return data.NewNullValue(), nil
		}
	}
	name := cookieGetName(cookie)
	path := cookieGetPath(cookie)
	if st.queued[name] == nil {
		st.queued[name] = map[string]*data.ClassValue{}
	}
	st.queued[name][path] = cookie
	syncQueuedProperty(cv, st)
	return data.NewNullValue(), nil
}

func jarGetQueued(ctx data.Context) (data.GetValue, data.Control) {
	_, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	arr := data.NewArrayValue(nil).(*data.ArrayValue)
	i := 0
	for _, byPath := range st.queued {
		for _, c := range byPath {
			arr.SetIntKey(i, c)
			i++
		}
	}
	return arr, nil
}

func jarHasQueued(ctx data.Context) (data.GetValue, data.Control) {
	_, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := ""
	if v := kit.Arg(ctx, 0); v != nil {
		key = v.AsString()
	}
	pathArg := kit.Arg(ctx, 1)
	byPath, ok := st.queued[key]
	if !ok || len(byPath) == 0 {
		return data.NewBoolValue(false), nil
	}
	if pathArg == nil || kit.IsNull(pathArg) {
		return data.NewBoolValue(true), nil
	}
	_, ok = byPath[pathArg.AsString()]
	return data.NewBoolValue(ok), nil
}

func jarQueued(ctx data.Context) (data.GetValue, data.Control) {
	_, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := ""
	if v := kit.Arg(ctx, 0); v != nil {
		key = v.AsString()
	}
	def := kit.Arg(ctx, 1)
	pathArg := kit.Arg(ctx, 2)
	byPath, ok := st.queued[key]
	if !ok || len(byPath) == 0 {
		if def == nil {
			return data.NewNullValue(), nil
		}
		return def, nil
	}
	if pathArg == nil || kit.IsNull(pathArg) {
		var last *data.ClassValue
		for _, c := range byPath {
			last = c
		}
		if last == nil {
			if def == nil {
				return data.NewNullValue(), nil
			}
			return def, nil
		}
		return last, nil
	}
	ps := pathArg.AsString()
	if c, ok := byPath[ps]; ok {
		return c, nil
	}
	if def == nil {
		return data.NewNullValue(), nil
	}
	return def, nil
}

func jarExpire(ctx data.Context) (data.GetValue, data.Control) {
	cookie, ctl := jarForget(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cv, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	c, ok := cookie.(*data.ClassValue)
	if !ok || c == nil {
		return data.NewNullValue(), nil
	}
	name := cookieGetName(c)
	path := cookieGetPath(c)
	if st.queued[name] == nil {
		st.queued[name] = map[string]*data.ClassValue{}
	}
	st.queued[name][path] = c
	syncQueuedProperty(cv, st)
	return data.NewNullValue(), nil
}

func jarUnqueue(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	name := ""
	if v := kit.Arg(ctx, 0); v != nil {
		name = v.AsString()
	}
	pathArg := kit.Arg(ctx, 1)
	if pathArg == nil || kit.IsNull(pathArg) {
		delete(st.queued, name)
	} else {
		ps := pathArg.AsString()
		if byPath, ok := st.queued[name]; ok {
			delete(byPath, ps)
			if len(byPath) == 0 {
				delete(st.queued, name)
			}
		}
	}
	syncQueuedProperty(cv, st)
	return data.NewNullValue(), nil
}

func jarSetDefaultPathAndDomain(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if v := kit.Arg(ctx, 0); v != nil {
		st.path = v.AsString()
		_ = cv.SetProperty("path", data.NewStringValue(st.path))
	}
	st.domain = kit.Arg(ctx, 1)
	_ = cv.SetProperty("domain", st.domain)
	st.secure = kit.Arg(ctx, 2)
	_ = cv.SetProperty("secure", st.secure)
	if v := kit.Arg(ctx, 3); v != nil && !kit.IsNull(v) {
		st.sameSite = v.AsString()
		_ = cv.SetProperty("sameSite", data.NewStringValue(st.sameSite))
	}
	return cv, nil
}

func jarFlushQueued(ctx data.Context) (data.GetValue, data.Control) {
	cv, st, ctl := jarRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	st.queued = map[string]map[string]*data.ClassValue{}
	syncQueuedProperty(cv, st)
	return cv, nil
}

func syncQueuedProperty(cv *data.ClassValue, st *jarState) {
	_ = cv.SetProperty("queued", data.NewArrayValue(nil))
}

func newSymfonyCookie(ctx data.Context, _ *data.ClassValue, name string, value, expire, path, domain, secure, httpOnly, raw, sameSite data.Value) (*data.ClassValue, data.Control) {
	stmt, ctl := ctx.GetVM().GetOrLoadClass(symfonyCookie)
	if ctl != nil {
		return nil, ctl
	}
	cv := data.NewClassValue(stmt, ctx.CreateBaseContext())
	cons := stmt.GetConstruct()
	if cons == nil {
		return cv, nil
	}
	cctx := cv.CreateContext(cons.GetVariables())
	vars := cons.GetVariables()
	args := []data.Value{data.NewStringValue(name), value, expire, path, domain, secure, httpOnly, raw, sameSite}
	for i, v := range vars {
		if i < len(args) {
			_ = cctx.SetVariableValue(v, args[i])
		}
	}
	cctx.SetFlatCallArgs(args)
	if _, ctl := cons.Call(cctx); ctl != nil {
		return nil, ctl
	}
	return cv, nil
}

func cookieGetName(cv *data.ClassValue) string {
	v, _ := cv.GetProperty("name")
	if v != nil {
		return v.AsString()
	}
	return ""
}

func cookieGetPath(cv *data.ClassValue) string {
	v, _ := cv.GetProperty("path")
	if v != nil {
		return v.AsString()
	}
	return "/"
}

func argString(ctx data.Context, i int, def string) string {
	if v := kit.Arg(ctx, i); v != nil && !kit.IsNull(v) {
		return v.AsString()
	}
	return def
}
