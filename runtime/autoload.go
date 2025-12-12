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
