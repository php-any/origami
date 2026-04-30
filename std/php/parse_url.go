package php

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PHP_URL_* 常量值
const (
	PHP_URL_SCHEME   = 0
	PHP_URL_HOST     = 1
	PHP_URL_PORT     = 2
	PHP_URL_USER     = 3
	PHP_URL_PASS     = 4
	PHP_URL_PATH     = 5
	PHP_URL_QUERY    = 6
	PHP_URL_FRAGMENT = 7
)

func NewParseUrlFunction() data.FuncStmt {
	return &ParseUrlFunction{}
}

type ParseUrlFunction struct{}

func (f *ParseUrlFunction) GetName() string {
	return "parse_url"
}

func (f *ParseUrlFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "url", 0, nil, nil),
		node.NewParameter(nil, "component", 1, data.NewIntValue(-1), nil),
	}
}

func (f *ParseUrlFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "url", 0, nil),
		node.NewVariable(nil, "component", 1, nil),
	}
}

func (f *ParseUrlFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	urlValue, exists := ctx.GetIndexValue(0)
	if !exists || urlValue == nil {
		return data.NewBoolValue(false), nil
	}

	rawURL := urlValue.AsString()

	// 获取 component 参数
	component := -1
	componentValue, exists := ctx.GetIndexValue(1)
	if exists && componentValue != nil {
		if asInt, ok := componentValue.(data.AsInt); ok {
			component, _ = asInt.AsInt()
		}
	}

	// 解析 URL
	// PHP 的 parse_url 对于没有 scheme 的 URL 也可以解析
	// 但 Go 的 url.Parse 对此行为不同，需要特殊处理
	parsed, scheme, err := phpParseURL(rawURL)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	// 提取各组件
	host := parsed.Hostname()
	portStr := parsed.Port()
	user := ""
	pass := ""
	hasPass := false
	if parsed.User != nil {
		user = parsed.User.Username()
		pass, hasPass = parsed.User.Password()
	}
	path := parsed.Path
	query := parsed.RawQuery
	fragment := parsed.Fragment

	// 如果指定了 component，返回单个组件
	if component != -1 {
		switch component {
		case PHP_URL_SCHEME:
			if scheme == "" {
				return data.NewNullValue(), nil
			}
			return data.NewStringValue(scheme), nil
		case PHP_URL_HOST:
			if host == "" {
				return data.NewNullValue(), nil
			}
			return data.NewStringValue(host), nil
		case PHP_URL_PORT:
			if portStr == "" {
				return data.NewNullValue(), nil
			}
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return data.NewNullValue(), nil
			}
			return data.NewIntValue(port), nil
		case PHP_URL_USER:
			if user == "" && parsed.User == nil {
				return data.NewNullValue(), nil
			}
			return data.NewStringValue(user), nil
		case PHP_URL_PASS:
			if !hasPass {
				return data.NewNullValue(), nil
			}
			return data.NewStringValue(pass), nil
		case PHP_URL_PATH:
			if path == "" {
				return data.NewNullValue(), nil
			}
			return data.NewStringValue(path), nil
		case PHP_URL_QUERY:
			if query == "" {
				return data.NewNullValue(), nil
			}
			return data.NewStringValue(query), nil
		case PHP_URL_FRAGMENT:
			if fragment == "" {
				return data.NewNullValue(), nil
			}
			return data.NewStringValue(fragment), nil
		default:
			return data.NewNullValue(), nil
		}
	}

	// 返回关联数组
	result := data.NewObjectValue()
	if scheme != "" {
		result.SetProperty("scheme", data.NewStringValue(scheme))
	}
	if host != "" {
		result.SetProperty("host", data.NewStringValue(host))
	}
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			result.SetProperty("port", data.NewIntValue(port))
		}
	}
	if parsed.User != nil {
		result.SetProperty("user", data.NewStringValue(user))
		if hasPass {
			result.SetProperty("pass", data.NewStringValue(pass))
		}
	}
	if path != "" {
		result.SetProperty("path", data.NewStringValue(path))
	}
	if query != "" {
		result.SetProperty("query", data.NewStringValue(query))
	}
	if fragment != "" {
		result.SetProperty("fragment", data.NewStringValue(fragment))
	}

	return result, nil
}

// phpParseURL 以 PHP parse_url 兼容方式解析 URL
// 返回解析结果、scheme 和错误
func phpParseURL(rawURL string) (*url.URL, string, error) {
	// 提取 scheme
	scheme := ""
	urlToParse := rawURL

	// 检查是否有 scheme (如 http://, https://, ftp://, etc.)
	if idx := strings.Index(rawURL, "://"); idx > 0 {
		// 验证 scheme 部分只包含合法字符
		potentialScheme := rawURL[:idx]
		validScheme := true
		for i, c := range potentialScheme {
			if i == 0 {
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
					validScheme = false
					break
				}
			} else {
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '+' || c == '-' || c == '.') {
					validScheme = false
					break
				}
			}
		}
		if validScheme {
			scheme = potentialScheme
		}
	}

	// 如果没有 scheme，添加一个临时的以便 Go 的 url.Parse 正确解析
	if scheme == "" {
		// 对于 //host/path 格式
		if strings.HasPrefix(rawURL, "//") {
			urlToParse = "placeholder:" + rawURL
		} else {
			// 纯路径，直接解析
			parsed, err := url.Parse(rawURL)
			if err != nil {
				return nil, "", err
			}
			return parsed, "", nil
		}
	}

	parsed, err := url.Parse(urlToParse)
	if err != nil {
		return nil, "", err
	}

	return parsed, scheme, nil
}

// InitParseUrlConstants 初始化 parse_url 相关常量
func InitParseUrlConstants(vm data.VM) {
	vm.SetConstant("PHP_URL_SCHEME", data.NewIntValue(PHP_URL_SCHEME))
	vm.SetConstant("PHP_URL_HOST", data.NewIntValue(PHP_URL_HOST))
	vm.SetConstant("PHP_URL_PORT", data.NewIntValue(PHP_URL_PORT))
	vm.SetConstant("PHP_URL_USER", data.NewIntValue(PHP_URL_USER))
	vm.SetConstant("PHP_URL_PASS", data.NewIntValue(PHP_URL_PASS))
	vm.SetConstant("PHP_URL_PATH", data.NewIntValue(PHP_URL_PATH))
	vm.SetConstant("PHP_URL_QUERY", data.NewIntValue(PHP_URL_QUERY))
	vm.SetConstant("PHP_URL_FRAGMENT", data.NewIntValue(PHP_URL_FRAGMENT))
}
