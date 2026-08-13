package http

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	jsonSerializer "github.com/php-any/origami/std/serializer/json"
	"github.com/php-any/origami/utils"
)

// RequestBindMethod 绑定请求数据到指定 DTO 类。
type RequestBindMethod struct {
	source *httpsrc.Request
}

func (h *RequestBindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	vm := ctx.GetVM()
	if vm == nil {
		return data.NewObjectValue(), nil
	}

	classStmt, acl := vm.GetOrLoadClass(param0)
	if acl != nil {
		return nil, acl
	}
	if classStmt == nil {
		return nil, utils.NewThrowf("无法加载 DTO 类: %s", param0)
	}

	classInstance, acl := classStmt.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	classValue, ok := classInstance.(*data.ClassValue)
	if !ok {
		return data.NewObjectValue(), nil
	}

	serializer := jsonSerializer.NewJsonSerializer()
	contentType := h.source.Header.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		if err := ensureMultipartParsed(h.source); err != nil {
			return nil, utils.NewThrow(err)
		}
		if err := bindMultipartToClass(ctx, classValue, h.source); err != nil {
			return nil, utils.NewThrow(err)
		}
		return classValue, nil
	}

	if strings.Contains(contentType, "application/json") && h.source.Body != nil {
		body, err := io.ReadAll(h.source.Body)
		if err == nil && len(body) > 0 {
			if err := serializer.UnmarshalClass(body, classValue); err != nil {
				return nil, utils.NewThrow(err)
			}
			return classValue, nil
		}
	}

	if len(h.source.Form) > 0 {
		if err := bindFlatMapToClass(classValue, stringMapFromValues(h.source.Form)); err != nil {
			return nil, utils.NewThrow(err)
		}
		return classValue, nil
	}

	if query := h.source.URL.Query(); len(query) > 0 {
		if err := bindFlatMapToClass(classValue, stringMapFromValues(query)); err != nil {
			return nil, utils.NewThrow(err)
		}
		return classValue, nil
	}

	return classValue, nil
}

func stringMapFromValues(src map[string][]string) map[string]string {
	out := make(map[string]string, len(src))
	for key, values := range src {
		if len(values) > 0 {
			out[key] = values[0]
		}
	}
	return out
}

// bindFlatMapToClass 将 query/form 等扁平字符串映射到 DTO 属性（按属性类型转换）。
func bindFlatMapToClass(classValue *data.ClassValue, flat map[string]string) error {
	for key, raw := range flat {
		propStmt, ok := classValue.GetPropertyStmt(key)
		if !ok {
			continue
		}
		cp, ok := propStmt.(*node.ClassProperty)
		if !ok {
			continue
		}
		val, err := coerceFlatInput(raw, cp.Type)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
		classValue.SetProperty(key, val)
	}
	return nil
}

func bindMultipartToClass(ctx data.Context, classValue *data.ClassValue, r *httpsrc.Request) error {
	if r.MultipartForm != nil && len(r.MultipartForm.Value) > 0 {
		if err := bindFlatMapToClass(classValue, stringMapFromValues(r.MultipartForm.Value)); err != nil {
			return err
		}
	}
	for _, prop := range classValue.Class.GetPropertyList() {
		cp, ok := prop.(*node.ClassProperty)
		if !ok {
			continue
		}
		typeFQN, nullable := typeFQNFromTypes(cp.Type)
		if !isUploadedFileTypeFQN(typeFQN) {
			continue
		}
		header, err := uploadedFileFromForm(r, cp.Name)
		if err != nil || header == nil {
			if nullable {
				classValue.SetProperty(cp.Name, data.NewNullValue())
			} else {
				classValue.SetProperty(cp.Name, newUploadedFileValue(ctx, nil))
			}
			continue
		}
		classValue.SetProperty(cp.Name, newUploadedFileValue(ctx, header))
	}
	return nil
}

func coerceFlatInput(raw string, ty data.Types) (data.Value, error) {
	if ty == nil {
		return data.NewStringValue(raw), nil
	}
	if nt, ok := ty.(data.NullableType); ok {
		if strings.TrimSpace(raw) == "" {
			return data.NewNullValue(), nil
		}
		return coerceFlatInput(raw, nt.BaseType)
	}
	switch ty.(type) {
	case data.Int:
		if strings.TrimSpace(raw) == "" {
			return data.NewIntValue(0), nil
		}
		n, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil {
			return nil, fmt.Errorf("无法转换为 int: %q", raw)
		}
		return data.NewIntValue(n), nil
	case data.Float:
		if strings.TrimSpace(raw) == "" {
			return data.NewFloatValue(0), nil
		}
		f, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, fmt.Errorf("无法转换为 float: %q", raw)
		}
		return data.NewFloatValue(f), nil
	case data.Bool:
		if strings.TrimSpace(raw) == "" {
			return data.NewBoolValue(false), nil
		}
		b, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			return nil, fmt.Errorf("无法转换为 bool: %q", raw)
		}
		return data.NewBoolValue(b), nil
	case data.String:
		return data.NewStringValue(raw), nil
	default:
		return data.NewStringValue(raw), nil
	}
}

func (h *RequestBindMethod) GetName() string            { return "bind" }
func (h *RequestBindMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestBindMethod) GetIsStatic() bool          { return false }
func (h *RequestBindMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
	}
}
func (h *RequestBindMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
	}
}
func (h *RequestBindMethod) GetReturnType() data.Types { return data.NewBaseType("object") }
