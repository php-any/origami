package core

import (
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ComposerAutoloadFunction 实现 composer_autoload 函数
// 替代 require __DIR__.'/vendor/autoload.php'
// 读取 composer 的 autoload 配置，注册命名空间映射到解释器的 ClassPathManager 中
type ComposerAutoloadFunction struct{}

func NewComposerAutoloadFunction() data.FuncStmt {
	return &ComposerAutoloadFunction{}
}

func (f *ComposerAutoloadFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	dirValue, _ := ctx.GetIndexValue(0)

	var dir string
	switch d := dirValue.(type) {
	case data.AsString:
		dir = d.AsString()
	default:
		if dirValue == nil {
			return data.NewBoolValue(false), nil
		}
		dir = dirValue.AsString()
	}

	if dir == "" {
		return data.NewBoolValue(false), nil
	}

	// 转为绝对路径
	if !filepath.IsAbs(dir) {
		abs, err := filepath.Abs(dir)
		if err == nil {
			dir = abs
		}
	}

	vm := ctx.GetVM()
	result := data.NewObjectValue()

	// 加载 PSR-4 映射
	psr4File := filepath.Join(dir, "vendor", "composer", "autoload_psr-4.php")
	if _, err := os.Stat(psr4File); err == nil {
		v, acl := vm.LoadAndRun(psr4File)
		if acl != nil {
			return nil, acl
		}
		if arr, ok := v.(*data.ArrayValue); ok {
			count := 0
			for _, z := range arr.List {
				if z.Name == "" {
					continue
				}
				namespace := z.Name
				// 值是路径数组
				if paths, ok := z.Value.(*data.ArrayValue); ok {
					for _, pz := range paths.List {
						if pathStr, ok := pz.Value.(data.AsString); ok {
							p := pathStr.AsString()
							if p != "" {
								vm.AddNamespace(namespace, p)
								count++
							}
						}
					}
				}
			}
			result.SetProperty("psr4_count", data.NewIntValue(count))
		}
		if vv, ok := v.(data.Value); ok {
			result.SetProperty("psr4", vv)
		}
	}

	// 加载 classmap 映射
	classmapFile := filepath.Join(dir, "vendor", "composer", "autoload_classmap.php")
	if _, err := os.Stat(classmapFile); err == nil {
		v, acl := vm.LoadAndRun(classmapFile)
		if acl != nil {
			return nil, acl
		}
		if arr, ok := v.(*data.ArrayValue); ok {
			count := 0
			for _, z := range arr.List {
				if z.Name == "" {
					continue
				}
				className := z.Name
				if pathVal, ok := z.Value.(data.AsString); ok {
					p := pathVal.AsString()
					if p != "" {
						// 注册到类路径缓存中，避免后续动态查找
						vm.SetClassPathCache(className, p)
						count++
					}
				}
			}
			result.SetProperty("classmap_count", data.NewIntValue(count))
		}
		if vv, ok := v.(data.Value); ok {
			result.SetProperty("classmap", vv)
		}
	}

	return result, nil
}

func (f *ComposerAutoloadFunction) GetName() string {
	return "composer_autoload"
}

func (f *ComposerAutoloadFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "dir", 0, nil, nil),
	}
}

func (f *ComposerAutoloadFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "dir", 0, data.NewBaseType("string")),
	}
}
