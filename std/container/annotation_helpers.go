package container

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func annotationStringArg(ctx data.Context, index int) (string, data.Control) {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return "", nil
	}
	if s, ok := v.(data.AsString); ok {
		return s.AsString(), nil
	}
	return "", utils.NewThrow(errors.New("注解参数必须是字符串"))
}

func annotationTargetClass(ctx data.Context) (*node.ClassStatement, data.Control) {
	idx := 1
	if _, ok := ctx.GetIndexValue(1); !ok {
		idx = 0
	}
	tv, ok := ctx.GetIndexValue(idx)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 target 参数"))
	}
	anyT, ok := tv.(*data.AnyValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("target 类型错误"))
	}
	if cls := classStatementFromAny(anyT.Value); cls != nil {
		return cls, nil
	}
	return nil, utils.NewThrow(fmt.Errorf("注解只能用于类 (target=%T)", anyT.Value))
}

func classStatementFromAny(v any) *node.ClassStatement {
	switch t := v.(type) {
	case *node.ClassStatement:
		return t
	case *node.ClassRegisterStmt:
		return t.Class
	case *node.AbstractClassStatement:
		return t.ClassStatement
	case *node.ClassGeneric:
		return t.ClassStatement
	case data.GetValue:
		// 解开二次装箱的 GetValue，避免无限递归
		switch u := t.(type) {
		case *node.ClassStatement:
			return u
		case *node.ClassRegisterStmt:
			return u.Class
		case *node.AbstractClassStatement:
			return u.ClassStatement
		case *node.ClassGeneric:
			return u.ClassStatement
		default:
			return nil
		}
	default:
		return nil
	}
}

func annotationTargetParameter(ctx data.Context) (*node.Parameter, string, data.Control) {
	idx := 1
	if _, ok := ctx.GetIndexValue(1); !ok {
		idx = 0
	}
	tv, ok := ctx.GetIndexValue(idx)
	if !ok {
		return nil, "", utils.NewThrow(errors.New("缺少 target 参数"))
	}
	anyT, ok := tv.(*data.AnyValue)
	if !ok {
		return nil, "", utils.NewThrow(errors.New("target 类型错误"))
	}
	switch t := anyT.Value.(type) {
	case *node.Parameter:
		return t, "", nil
	case *node.PromotedParameter:
		return t.Parameter, "", nil
	case *node.ClassProperty:
		return nil, "", utils.NewThrow(errors.New("Container\\Annotation\\Inject 只能用于构造器参数"))
	default:
		return nil, "", utils.NewThrow(errors.New("注解目标类型不支持"))
	}
}

func scanDirectory(vm data.VM, dir string) data.Control {
	dir = filepath.Clean(dir)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".php" && ext != ".zy" {
			return nil
		}
		if _, acl := vm.LoadAndRun(path); acl != nil {
			return errors.New("scan load failed")
		}
		return nil
	})
	if err != nil {
		return data.NewErrorThrow(nil, err)
	}
	return nil
}

func RegisterClassAnnotation(ctx data.Context, lifetime Lifetime) data.Control {
	alias, _ := annotationStringArg(ctx, 0)
	cls, acl := annotationTargetClass(ctx)
	if acl != nil {
		return acl
	}
	metadataSetLifetime(cls.Name, lifetime, alias)
	if e := activeEngine(ctx); e != nil {
		e.RegisterClass(cls.Name, lifetime, alias)
	}
	return nil
}
