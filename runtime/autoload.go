package runtime

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func AddAutoLoad(fun *data.FuncValue) {
	parser.AddAutoLoad(fun)
}

func RemoveAutoLoad(fun *data.FuncValue) {
	parser.RemoveAutoLoad(fun)
}

// ClearAutoLoad 清空全部 autoload 回调，供开发模式热重载使用。
func ClearAutoLoad() {
	parser.ClearAutoLoad()
}

func GetAutoLoad() []*data.FuncValue {
	return parser.GetAutoLoad()
}

func CallAutoLoad(name string, ctx data.Context) (bool, data.Control) {
	return parser.CallAutoLoad(name, ctx)
}
