package http

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ServerStaticMethod: $server->static(prefix: "/assets/", dir: "./pages")
type ServerStaticMethod struct {
	server *ServerClass
}

func (h *ServerStaticMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 参数：prefix, dir
	prefixV, _ := ctx.GetIndexValue(0)
	dirV, _ := ctx.GetIndexValue(1)

	var prefix, dir string
	if s, ok := prefixV.(data.AsString); ok {
		prefix = s.AsString()
	}
	if s, ok := dirV.(data.AsString); ok {
		dir = s.AsString()
	}
	if prefix == "" {
		prefix = "/assets/"
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	if dir == "" {
		dir = "."
	}

	// 将目录转换为绝对路径，便于提前校验并在报错时给出明确信息
	if !filepath.IsAbs(dir) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, utils.NewThrow(err)
		}
		dir = filepath.Join(cwd, dir)
	}
	// 规范化路径
	dir = filepath.Clean(dir)
	// 提前校验目录是否存在且为目录
	info, err := os.Stat(dir)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	if !info.IsDir() {
		return nil, utils.NewThrow(errors.New("静态资源目录不是目录: " + dir))
	}

	// 组合路由前缀
	route := h.server.Prefix + prefix
	// 规范化：必须以 / 结尾用于子树匹配
	if !strings.HasSuffix(route, "/") {
		route += "/"
	}

	// StripPrefix 必须与注册的前缀完全一致（包含末尾斜杠）
	fs := http.StripPrefix(route, http.FileServer(http.Dir(dir)))
	var final http.Handler = fs
	if len(h.server.Middlewares) > 0 {
		final = applyMiddlewares(final, h.server.Middlewares)
	}

	// 仅注册 GET/HEAD
	h.server.source.Handle("GET "+route, final)
	h.server.source.Handle("HEAD "+route, final)
	return nil, nil
}

func (h *ServerStaticMethod) GetName() string            { return "static" }
func (h *ServerStaticMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerStaticMethod) GetIsStatic() bool          { return false }
func (h *ServerStaticMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "prefix", 0, data.NewStringValue("/assets/"), data.NewBaseType("string")),
		node.NewParameter(nil, "dir", 1, data.NewStringValue("."), data.NewBaseType("string")),
	}
}
func (h *ServerStaticMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "prefix", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "dir", 1, data.NewBaseType("string")),
	}
}
func (h *ServerStaticMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
