package node

import (
	"os"
	"strings"

	"github.com/php-any/origami/data"
)

// $_SERVER

type ServerVariable struct {
	*Node `pp:"-"`
}

var serverValue *data.ObjectValue

func NewServerVariable(from data.From) data.Variable {
	return &ServerVariable{Node: NewNode(from)}
}

func (v *ServerVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if serverValue == nil {
		serverValue = data.NewObjectValue()

		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			serverValue.SetProperty("REQUEST_METHOD", data.NewStringValue(httpReq.Method))
			serverValue.SetProperty("REQUEST_URI", data.NewStringValue(httpReq.RequestURI))
			serverValue.SetProperty("QUERY_STRING", data.NewStringValue(httpReq.URL.RawQuery))
			serverValue.SetProperty("HTTP_HOST", data.NewStringValue(httpReq.Host))
			serverValue.SetProperty("SERVER_NAME", data.NewStringValue(httpReq.Host))
			serverValue.SetProperty("SERVER_PORT", data.NewStringValue(httpReq.URL.Port()))
			serverValue.SetProperty("REMOTE_ADDR", data.NewStringValue(httpReq.RemoteAddr))
			// SCRIPT_NAME 勿用 URL.Path（如 "/"），否则 Symfony/Laravel 的 baseUrl 计算错误导致路由不匹配
			serverValue.SetProperty("SCRIPT_NAME", data.NewStringValue("/index.php"))
			serverValue.SetProperty("PHP_SELF", data.NewStringValue("/index.php"))

			for key, values := range httpReq.Header {
				if len(values) > 0 {
					headerKey := "HTTP_" + strings.ReplaceAll(strings.ToUpper(key), "-", "_")
					serverValue.SetProperty(headerKey, data.NewStringValue(values[0]))
				}
			}
		} else {
			serverValue.SetProperty("SERVER_SOFTWARE", data.NewStringValue("Origami"))
			registerArgcArgv := os.Getenv("ORIGAMI_PHPT_REGISTER_ARGC_ARGV")
			if registerArgcArgv != "0" && len(os.Args) > 1 {
				arr := make([]data.Value, 0, len(os.Args)-1)
				for _, s := range os.Args[1:] {
					arr = append(arr, data.NewStringValue(s))
				}
				serverValue.SetProperty("argv", data.NewArrayValue(arr))
				serverValue.SetProperty("argc", data.NewIntValue(len(arr)))
			}
		}
	}

	return serverValue, nil
}

func (v *ServerVariable) GetIndex() int       { return 0 }
func (v *ServerVariable) GetName() string     { return "$_SERVER" }
func (v *ServerVariable) GetType() data.Types { return nil }
func (v *ServerVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	// 允许设置 $_SERVER 值，使 Symfony/Laravel 的 Request::capture() 等能正常工作
	if objectValue, ok := value.(*data.ObjectValue); ok {
		serverValue = objectValue
		return nil
	}
	return data.NewErrorThrow(v.from, nil)
}
