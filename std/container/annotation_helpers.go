package container

import (
	"errors"
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
	cls, ok := anyT.Value.(*node.ClassStatement)
	if !ok {
		return nil, utils.NewThrow(errors.New("注解只能用于类"))
	}
	return cls, nil
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
		return nil, "", utils.NewThrow(errors.New("Container\\Inject 只能用于构造器参数"))
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

func registerClassAnnotation(ctx data.Context, lifetime Lifetime) data.Control {
	alias, _ := annotationStringArg(ctx, 0)
	cls, acl := annotationTargetClass(ctx)
	if acl != nil {
		return acl
	}
	metadataSetLifetime(cls.Name, lifetime, alias)
	activeEngine(ctx).RegisterClass(cls.Name, lifetime, alias)
	return nil
}
