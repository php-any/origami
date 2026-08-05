package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	netdata "github.com/php-any/origami/std/net/data"
	"github.com/php-any/origami/std/validation"
	"github.com/php-any/origami/utils"
)

// resolveHandlerArgs 按 HandlerSpec 从请求上下文解析控制器实参。
// 返回 nil args 表示校验失败且响应已写入（如未来 validate 失败）。
func resolveHandlerArgs(
	vm data.VM,
	ctx data.Context,
	spec netdata.HandlerSpec,
	reqProxy data.Value,
	resProxy data.Value,
) ([]data.Value, data.Control) {
	if len(spec.Params) == 0 {
		return []data.Value{}, nil
	}

	args := make([]data.Value, len(spec.Params))
	for _, binding := range spec.Params {
		if binding.Index < 0 || binding.Index >= len(args) {
			continue
		}
		switch binding.Source {
		case netdata.SourceRequest:
			args[binding.Index] = reqProxy
		case netdata.SourceResponse:
			args[binding.Index] = resProxy
		case netdata.SourceDTO:
			dto, acl := bindDTO(vm, ctx, reqProxy, resProxy, binding)
			if acl != nil {
				return nil, acl
			}
			if dto == nil {
				return nil, nil
			}
			args[binding.Index] = dto
		case netdata.SourcePath:
			val, acl := resolvePathParam(reqProxy, binding)
			if acl != nil {
				return nil, acl
			}
			if done, acl := validateScalarBinding(resProxy, binding, val); done {
				return nil, acl
			}
			args[binding.Index] = val
		case netdata.SourceQuery:
			val, acl := resolveQueryParam(reqProxy, binding)
			if acl != nil {
				return nil, acl
			}
			if done, acl := validateScalarBinding(resProxy, binding, val); done {
				return nil, acl
			}
			args[binding.Index] = val
		case netdata.SourceFormFile:
			val, acl := resolveFormFileParam(vm, ctx, reqProxy, resProxy, binding)
			if acl != nil {
				return nil, acl
			}
			if val == nil {
				return nil, nil
			}
			args[binding.Index] = val
		default:
			return nil, utils.NewThrowf("无法解析控制器参数 %q", binding.Name)
		}
	}
	return args, nil
}

func bindDTO(vm data.VM, ctx data.Context, reqProxy, resProxy data.Value, binding netdata.ParamBinding) (data.Value, data.Control) {
	reqCV, ok := reqProxy.(*data.ProxyValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("DTO 绑定需要 Request 代理对象"))
	}
	bindMethod, has := reqCV.GetMethod("bind")
	if !has {
		return nil, utils.NewThrow(errors.New("Request 缺少 bind 方法"))
	}
	vars := bindMethod.GetVariables()
	if len(vars) == 0 {
		return nil, utils.NewThrow(errors.New("Request::bind 参数无效"))
	}

	bindCtx := reqCV.CreateContext(vars)
	className := strings.TrimPrefix(binding.TypeFQN, "\\")
	bindCtx.SetVariableValue(vars[0], data.NewStringValue(className))
	dto, acl := bindMethod.Call(bindCtx)
	if acl != nil {
		return nil, acl
	}
	if dto == nil {
		return data.NewObjectValue(), nil
	}
	val, ok := dto.(data.Value)
	if !ok {
		val = data.NewAnyValue(dto)
	}

	if binding.Validate {
		if cv, ok := val.(*data.ClassValue); ok {
			if violations := validation.ValidateObject(cv); len(violations) > 0 {
				if acl := writeValidationError(resProxy, violations); acl != nil {
					return nil, acl
				}
				return nil, nil
			}
		}
	}
	return val, nil
}

func writeValidationError(resProxy data.Value, violations []validation.Violation) data.Control {
	resCV, ok := resProxy.(*data.ProxyValue)
	if !ok {
		return utils.NewThrow(errors.New("校验失败响应需要 Response 代理对象"))
	}
	errMethod, has := resCV.GetMethod("error")
	if !has {
		return utils.NewThrow(errors.New("Response 缺少 error 方法"))
	}
	vars := errMethod.GetVariables()
	errCtx := resCV.CreateContext(vars)
	errCtx.SetVariableValue(vars[0], data.NewStringValue("validation failed"))
	if len(vars) > 1 {
		errCtx.SetVariableValue(vars[1], data.NewIntValue(http.StatusUnprocessableEntity))
	}
	if len(vars) > 2 {
		items := make([]data.Value, 0, len(violations))
		for _, v := range violations {
			item := data.NewObjectValue()
			item.SetProperty("field", data.NewStringValue(v.Field))
			item.SetProperty("message", data.NewStringValue(v.Message))
			items = append(items, item)
		}
		errCtx.SetVariableValue(vars[2], data.NewArrayValue(items))
	}
	_, acl := errMethod.Call(errCtx)
	return acl
}

func resolvePathParam(reqProxy data.Value, binding netdata.ParamBinding) (data.Value, data.Control) {
	reqCV, ok := reqProxy.(*data.ProxyValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("路径参数解析需要 Request 代理对象"))
	}
	pathMethod, has := reqCV.GetMethod("pathValue")
	if !has {
		return nil, utils.NewThrow(errors.New("Request 缺少 pathValue 方法"))
	}
	vars := pathMethod.GetVariables()
	if len(vars) == 0 {
		return nil, utils.NewThrow(errors.New("Request::pathValue 参数无效"))
	}

	pathCtx := reqCV.CreateContext(vars)
	pathCtx.SetVariableValue(vars[0], data.NewStringValue(binding.PathKey))
	raw, acl := pathMethod.Call(pathCtx)
	if acl != nil {
		return nil, acl
	}

	rawStr := ""
	if raw != nil {
		if sv, ok := raw.(data.AsString); ok {
			rawStr = sv.AsString()
		}
	}
	if rawStr == "" && binding.Nullable {
		return data.NewNullValue(), nil
	}
	return resolveScalarValue(rawStr, binding, "路径")
}

func resolveQueryParam(reqProxy data.Value, binding netdata.ParamBinding) (data.Value, data.Control) {
	reqCV, ok := reqProxy.(*data.ProxyValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("查询参数解析需要 Request 代理对象"))
	}
	inputMethod, has := reqCV.GetMethod("input")
	if !has {
		return nil, utils.NewThrow(errors.New("Request 缺少 input 方法"))
	}
	vars := inputMethod.GetVariables()
	if len(vars) == 0 {
		return nil, utils.NewThrow(errors.New("Request::input 参数无效"))
	}

	inputCtx := reqCV.CreateContext(vars)
	inputCtx.SetVariableValue(vars[0], data.NewStringValue(binding.QueryKey))
	raw, acl := inputMethod.Call(inputCtx)
	if acl != nil {
		return nil, acl
	}

	rawStr := inputRawString(raw)
	if rawStr == "" && binding.Nullable {
		return data.NewNullValue(), nil
	}
	return resolveScalarValue(rawStr, binding, "查询")
}

func resolveFormFileParam(
	vm data.VM,
	ctx data.Context,
	reqProxy data.Value,
	resProxy data.Value,
	binding netdata.ParamBinding,
) (data.Value, data.Control) {
	_ = vm
	r := requestFromProxy(reqProxy)
	header, err := uploadedFileFromForm(r, binding.FileKey)
	if err != nil || header == nil {
		if binding.Nullable {
			return data.NewNullValue(), nil
		}
		if binding.Validate && len(binding.Constraints) > 0 {
			if done, acl := validateScalarBinding(resProxy, binding, data.NewNullValue()); done {
				return nil, acl
			}
		}
		return newUploadedFileValue(ctx, nil), nil
	}
	val := newUploadedFileValue(ctx, header)
	if done, acl := validateScalarBinding(resProxy, binding, val); done {
		return nil, acl
	}
	return val, nil
}

// resolveScalarValue 将路径/查询标量绑定为 Value；有约束时空值留给校验层处理。
func resolveScalarValue(rawStr string, binding netdata.ParamBinding, kind string) (data.Value, data.Control) {
	if rawStr == "" && binding.TypeFQN == "string" {
		return data.NewStringValue(""), nil
	}
	if rawStr == "" && binding.TypeFQN != "string" && binding.Validate && len(binding.Constraints) > 0 {
		return data.NewNullValue(), nil
	}
	return coerceScalar(rawStr, binding.TypeFQN, binding.Name, kind)
}

func inputRawString(raw data.GetValue) string {
	if raw == nil {
		return ""
	}
	if _, ok := raw.(*data.NullValue); ok {
		return ""
	}
	if av, ok := raw.(*data.AnyValue); ok && av.Value == nil {
		return ""
	}
	if sv, ok := raw.(data.AsString); ok {
		return sv.AsString()
	}
	return ""
}

func validateScalarBinding(resProxy data.Value, binding netdata.ParamBinding, val data.Value) (bool, data.Control) {
	if !binding.Validate || len(binding.Constraints) == 0 {
		return false, nil
	}
	label := binding.Label
	if label == "" {
		label = binding.Name
	}
	violations := validation.ValidateConstraints(label, binding.Name, binding.Constraints, val)
	if len(violations) == 0 {
		return false, nil
	}
	if acl := writeValidationError(resProxy, violations); acl != nil {
		return true, acl
	}
	return true, nil
}

func coerceScalar(raw, typeFQN, paramName, kind string) (data.Value, data.Control) {
	switch typeFQN {
	case "string":
		return data.NewStringValue(raw), nil
	case "int":
		if raw == "" {
			return nil, utils.NewThrow(fmt.Errorf("%s参数 $%s 无法转换为 int：值为空", kind, paramName))
		}
		n, err := strconv.Atoi(raw)
		if err != nil {
			return nil, utils.NewThrow(fmt.Errorf("%s参数 $%s 无法转换为 int：%q", kind, paramName, raw))
		}
		return data.NewIntValue(n), nil
	case "float":
		if raw == "" {
			return nil, utils.NewThrow(fmt.Errorf("%s参数 $%s 无法转换为 float：值为空", kind, paramName))
		}
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, utils.NewThrow(fmt.Errorf("%s参数 $%s 无法转换为 float：%q", kind, paramName, raw))
		}
		return data.NewFloatValue(f), nil
	case "bool":
		if raw == "" {
			return nil, utils.NewThrow(fmt.Errorf("%s参数 $%s 无法转换为 bool：值为空", kind, paramName))
		}
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, utils.NewThrow(fmt.Errorf("%s参数 $%s 无法转换为 bool：%q", kind, paramName, raw))
		}
		return data.NewBoolValue(b), nil
	default:
		return nil, utils.NewThrow(fmt.Errorf("不支持的%s参数类型 %s", kind, typeFQN))
	}
}
