package node

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
)

// $_SERVER

type ServerVariable struct {
	*Node `pp:"-"`
}

func NewServerVariable(from data.From) data.Variable {
	return &ServerVariable{Node: NewNode(from)}
}

func (v *ServerVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return superglobalArray(ctx, "_SERVER", func() *data.ArrayValue {
		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			obj := data.NewArrayValueFromSlots(nil)
			obj.SetStringKey("REQUEST_METHOD", data.NewStringValue(httpReq.Method))
			uri := httpReq.RequestURI
			if uri == "" {
				uri = httpReq.URL.RequestURI()
			}
			obj.SetStringKey("REQUEST_URI", data.NewStringValue(uri))
			obj.SetStringKey("QUERY_STRING", data.NewStringValue(httpReq.URL.RawQuery))
			obj.SetStringKey("HTTP_HOST", data.NewStringValue(httpReq.Host))
			host, port, err := net.SplitHostPort(httpReq.Host)
			if err != nil {
				host = strings.Trim(httpReq.Host, "[]")
				port = "80"
				if httpReq.TLS != nil || httpReq.URL.Scheme == "https" {
					port = "443"
				}
			}
			obj.SetStringKey("SERVER_NAME", data.NewStringValue(host))
			obj.SetStringKey("SERVER_PORT", data.NewStringValue(port))
			remote, remotePort, err := net.SplitHostPort(httpReq.RemoteAddr)
			if err != nil {
				remote = httpReq.RemoteAddr
			}
			obj.SetStringKey("REMOTE_ADDR", data.NewStringValue(remote))
			obj.SetStringKey("REMOTE_PORT", data.NewStringValue(remotePort))
			obj.SetStringKey("SERVER_PROTOCOL", data.NewStringValue(httpReq.Proto))
			if httpReq.ContentLength >= 0 && httpReq.Body != nil && httpReq.Body != http.NoBody {
				obj.SetStringKey("CONTENT_LENGTH", data.NewStringValue(strconv.FormatInt(httpReq.ContentLength, 10)))
			}
			if httpReq.TLS != nil || httpReq.URL.Scheme == "https" {
				obj.SetStringKey("HTTPS", data.NewStringValue("on"))
			}
			obj.SetStringKey("SCRIPT_NAME", data.NewStringValue("/index.php"))
			obj.SetStringKey("PHP_SELF", data.NewStringValue("/index.php"))
			for key, values := range httpReq.Header {
				if len(values) > 0 {
					headerKey := "HTTP_" + strings.ReplaceAll(strings.ToUpper(key), "-", "_")
					if strings.EqualFold(key, "Content-Type") || strings.EqualFold(key, "Content-Length") {
						headerKey = strings.TrimPrefix(headerKey, "HTTP_")
					}
					obj.SetStringKey(headerKey, data.NewStringValue(strings.Join(values, ", ")))
				}
			}
			return obj
		}

		server := data.NewArrayValueFromSlots(nil)
		for _, entry := range EnvironmentEntries(ctx) {
			if name, value, ok := strings.Cut(entry, "="); ok {
				server.SetStringKey(name, data.NewStringValue(value))
			}
		}
		server.SetStringKey("SERVER_SOFTWARE", data.NewStringValue("Origami"))
		if len(os.Args) > 1 {
			server.SetStringKey("PHP_SELF", data.NewStringValue(os.Args[1]))
			server.SetStringKey("SCRIPT_NAME", data.NewStringValue(os.Args[1]))
			server.SetStringKey("SCRIPT_FILENAME", data.NewStringValue(os.Args[1]))
		}
		ensureCLIArgv(server)
		return server
	}), nil
}

func ensureCLIArgv(server *data.ArrayValue) {
	if server == nil || os.Getenv("ORIGAMI_PHPT_REGISTER_ARGC_ARGV") == "0" {
		return
	}
	if slot, _ := server.LookupZValByStringKey("argv"); slot != nil {
		return
	}
	arr := make([]data.Value, 0)
	if len(os.Args) > 1 {
		arr = make([]data.Value, 0, len(os.Args)-1)
		for _, s := range os.Args[1:] {
			arr = append(arr, data.NewStringValue(s))
		}
	}
	server.SetStringKey("argv", data.NewArrayValue(arr))
	server.SetStringKey("argc", data.NewIntValue(len(arr)))
}

func (v *ServerVariable) GetIndex() int       { return -1 }
func (v *ServerVariable) GetName() string     { return "$_SERVER" }
func (v *ServerVariable) GetType() data.Types { return nil }
func (v *ServerVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return setSuperglobalArray(ctx, "_SERVER", value)
}

func (v *ServerVariable) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	if _, ctl := v.GetValue(ctx); ctl != nil {
		return nil, ctl
	}
	return ctx.GetVM().EnsureGlobalZVal("_SERVER"), nil
}
func (v *ServerVariable) SuperglobalName() string { return "_SERVER" }
