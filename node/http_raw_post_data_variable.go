package node

import (
	"fmt"
	"os"

	"github.com/php-any/origami/data"
)

// HttpRawPostDataVariable 表示已移除的 $HTTP_RAW_POST_DATA（PHP 7+ 使用 php://input）。
type HttpRawPostDataVariable struct {
	*Node `pp:"-"`
}

func NewHttpRawPostDataVariable(from data.From) *HttpRawPostDataVariable {
	return &HttpRawPostDataVariable{Node: NewNode(from)}
}

func (v *HttpRawPostDataVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	file := "Unknown"
	line := 0
	if v.GetFrom() != nil {
		if src := v.GetFrom().GetSource(); src != "" {
			file = src
		}
		if sl, _ := v.GetFrom().GetStartPosition(); sl >= 0 {
			line = sl + 1
		}
	}
	if os.Getenv("ORIGAMI_HTTP_RAW_WARN_NL") == "1" {
		fmt.Fprint(os.Stderr, "\n")
	}
	fmt.Fprintf(os.Stderr, "Warning: Undefined variable $HTTP_RAW_POST_DATA in %s on line %d\n", file, line)
	return data.NewNullValue(), nil
}

func (v *HttpRawPostDataVariable) GetIndex() int       { return 0 }
func (v *HttpRawPostDataVariable) GetName() string     { return "HTTP_RAW_POST_DATA" }
func (v *HttpRawPostDataVariable) GetType() data.Types { return nil }

func (v *HttpRawPostDataVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return nil
}
