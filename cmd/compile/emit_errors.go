package compile

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EmitError 表示 AST 节点无法转译的错误。
type EmitError struct {
	File    string
	Type    string
	Line    int
	Column  int
	Message string
}

func (e *EmitError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("compile error: %s\n  unsupported AST node %s at line %d, col %d\n  %s",
			e.File, e.Type, e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("compile error: %s\n  unsupported AST node %s\n  %s",
		e.File, e.Type, e.Message)
}

func newEmitError(file string, v data.GetValue, message string) *EmitError {
	err := &EmitError{
		File:    file,
		Type:    fmt.Sprintf("%T", v),
		Message: message,
	}
	if from, ok := v.(node.GetFrom); ok && from.GetFrom() != nil {
		if line, col := from.GetFrom().GetStartPosition(); line >= 0 {
			err.Line = line + 1
			err.Column = col + 1
		}
	}
	return err
}

func wrapEmitError(file string, v data.GetValue, err error) error {
	if err == nil {
		return nil
	}
	if ee, ok := err.(*EmitError); ok {
		return ee
	}
	return fmt.Errorf("compile error: %s\n  node %T: %w", file, v, err)
}
