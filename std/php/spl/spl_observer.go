package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewSplObserverInterface 返回 PHP SplObserver 接口
func NewSplObserverInterface() data.InterfaceStmt {
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "update", "public", []data.GetValue{
			node.NewParameter(nil, "subject", 0, nil, data.NewBaseType("SplSubject")),
		}, nil),
	}
	return node.NewInterfaceStatement(nil, "SplObserver", nil, methods)
}

// NewSplSubjectInterface 返回 PHP SplSubject 接口
func NewSplSubjectInterface() data.InterfaceStmt {
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "attach", "public", []data.GetValue{
			node.NewParameter(nil, "observer", 0, nil, data.NewBaseType("SplObserver")),
		}, nil),
		node.NewInterfaceMethod(nil, "detach", "public", []data.GetValue{
			node.NewParameter(nil, "observer", 0, nil, data.NewBaseType("SplObserver")),
		}, nil),
		node.NewInterfaceMethod(nil, "notify", "public", []data.GetValue{}, nil),
	}
	return node.NewInterfaceStatement(nil, "SplSubject", nil, methods)
}
