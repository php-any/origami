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

func GetAutoLoad() []*data.FuncValue {
	return parser.GetAutoLoad()
}

func CallAutoLoad(name string, ctx data.Context) (bool, data.Control) {
	return parser.CallAutoLoad(name, ctx)
}
