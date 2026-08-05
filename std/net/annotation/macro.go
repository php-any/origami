package annotation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	netdata "github.com/php-any/origami/std/net/data"
	"github.com/php-any/origami/std/validation"
	"github.com/php-any/origami/utils"
)

var pathVarPattern = regexp.MustCompile(`\{([^}/:]+)\}`)

// MacroExpander 宏注解在扫描阶段对目标 AST 进行分析或变换。
type MacroExpander interface {
	Expand(target data.Method, ctx data.Context, routePath string) (effective data.Method, spec netdata.HandlerSpec, acl data.Control)
}

// ExpandHTTPHandlerMethod 分析 HTTP 映射目标方法的形参，产出 HandlerSpec。
func ExpandHTTPHandlerMethod(target data.Method, ctx data.Context, routePath string) (data.Method, netdata.HandlerSpec, data.Control) {
	if target == nil {
		return nil, netdata.HandlerSpec{}, nil
	}
	var vm data.VM
	if ctx != nil {
		vm = ctx.GetVM()
	}
	spec, acl := AnalyzeHandlerParams(target, vm, routePath)
	if acl != nil {
		return nil, netdata.HandlerSpec{}, acl
	}
	return target, spec, nil
}

// AnalyzeHandlerParams 根据方法形参类型与路由路径生成参数绑定计划。
func AnalyzeHandlerParams(method data.Method, vm data.VM, routePath string) (netdata.HandlerSpec, data.Control) {
	params := method.GetParams()
	pathVars := extractPathVars(routePath)
	spec := netdata.HandlerSpec{Params: make([]netdata.ParamBinding, 0, len(params))}
	for i, param := range params {
		p, ok := param.(*node.Parameter)
		if !ok {
			continue
		}
		binding, acl := analyzeParameter(p, i, vm, pathVars, method.GetName(), routePath)
		if acl != nil {
			return netdata.HandlerSpec{}, acl
		}
		spec.Params = append(spec.Params, binding)
	}
	return spec, nil
}

func analyzeParameter(
	p *node.Parameter,
	index int,
	vm data.VM,
	pathVars map[string]bool,
	methodName, routePath string,
) (netdata.ParamBinding, data.Control) {
	typeFQN, nullable := unwrapTypeFQN(p.Type)
	binding := netdata.ParamBinding{
		Name:     p.Name,
		Label:    validation.AnnotationDisplayName(p.Annotations, p.Name),
		TypeFQN:  typeFQN,
		Index:    index,
		Nullable: nullable,
	}

	if typeFQN == "" {
		return binding, utils.NewThrow(fmt.Errorf(
			"控制器方法 %s 的参数 $%s 缺少类型声明",
			methodName, p.Name,
		))
	}

	switch typeFQN {
	case "Net\\Http\\Request":
		binding.Source = netdata.SourceRequest
	case "Net\\Http\\Response":
		binding.Source = netdata.SourceResponse
	case "Net\\Http\\UploadedFile":
		binding.Source = netdata.SourceFormFile
		binding.FileKey = p.Name
	default:
		if isScalarTypeFQN(typeFQN) {
			constraints := validation.ConstraintAnnotations(p.Annotations)
			if len(constraints) > 0 {
				binding.Constraints = constraints
				binding.Validate = true
			}
			if pathVars[p.Name] {
				binding.Source = netdata.SourcePath
				binding.PathKey = p.Name
			} else {
				binding.Source = netdata.SourceQuery
				binding.QueryKey = p.Name
			}
		} else {
			binding.Source = netdata.SourceDTO
			binding.Validate = hasValidationConstraints(vm, typeFQN)
		}
	}
	return binding, nil
}

func extractPathVars(routePath string) map[string]bool {
	vars := make(map[string]bool)
	for _, m := range pathVarPattern.FindAllStringSubmatch(routePath, -1) {
		if len(m) > 1 {
			vars[m[1]] = true
		}
	}
	return vars
}

func unwrapTypeFQN(ty data.Types) (string, bool) {
	if ty == nil {
		return "", false
	}
	nullable := false
	if nt, ok := ty.(data.NullableType); ok {
		nullable = true
		ty = nt.BaseType
	}
	return strings.TrimPrefix(ty.String(), "\\"), nullable
}

func isScalarTypeFQN(typeFQN string) bool {
	switch typeFQN {
	case "string", "int", "float", "bool":
		return true
	default:
		return false
	}
}

func hasValidationConstraints(vm data.VM, typeFQN string) bool {
	if vm == nil || typeFQN == "" {
		return false
	}
	var cls data.ClassStmt
	if loaded, ok := vm.GetClass(typeFQN); ok && loaded != nil {
		cls = loaded
	} else {
		pkg, acl := vm.LoadPkg(typeFQN)
		if acl != nil || pkg == nil {
			return false
		}
		cast, ok := pkg.(data.ClassStmt)
		if !ok {
			return false
		}
		cls = cast
	}
	for _, prop := range cls.GetPropertyList() {
		cp, ok := prop.(*node.ClassProperty)
		if !ok {
			continue
		}
		for _, ann := range cp.Annotations {
			if ann == nil || ann.Class == nil {
				continue
			}
			if strings.HasPrefix(ann.Class.GetName(), "Validation\\Annotation\\") {
				return true
			}
		}
	}
	return false
}

// AsMacroExpander 若注解类实现 MacroExpander 则返回其实例。
func AsMacroExpander(class data.ClassStmt) MacroExpander {
	if exp, ok := class.(MacroExpander); ok {
		return exp
	}
	return nil
}
