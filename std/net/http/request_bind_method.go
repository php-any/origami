package http

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	jsonSerializer "github.com/php-any/origami/std/serializer/json"
	"github.com/php-any/origami/utils"
)

// RequestBindMethod 绑定 JSON 数据到指定类
type RequestBindMethod struct {
	source *httpsrc.Request
}

func (h *RequestBindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	// 获取要绑定的类
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	// 获取 VM 来查找类定义
	vm := ctx.GetVM()
	if vm != nil {
		// 根据类名获取类定义
		if classStmt, ok := vm.GetClass(param0); ok {
			// 创建类实例
			classInstance, _ := classStmt.GetValue(ctx)
			if classValue, ok := classInstance.(*data.ClassValue); ok {
				// 检查内容类型
				contentType := h.source.Header.Get("Content-Type")

				if strings.Contains(contentType, "application/json") {
					// 处理 JSON 请求体，直接使用 UnmarshalClass
					if h.source.Body != nil {
						body, err := io.ReadAll(h.source.Body)
						if err == nil && len(body) > 0 {
							// 直接使用 JSON 序列化器处理原始 JSON 字节
							serializer := jsonSerializer.NewJsonSerializer()
							err = serializer.UnmarshalClass(body, classValue)
							if err != nil {
								return data.NewObjectValue(), nil
							}
							return classValue, nil
						}
					}
				} else {
					// 处理表单数据
					if h.source.Form != nil {
						// 将表单数据转换为 JSON 格式
						formData := make(map[string]interface{})
						for key, values := range h.source.Form {
							if len(values) > 0 {
								formData[key] = values[0]
							}
						}

						// 将表单数据转换为 JSON 字节
						jsonBytes, err := json.Marshal(formData)
						if err != nil {
							return data.NewObjectValue(), nil
						}

						// 使用 JSON 序列化器处理
						serializer := jsonSerializer.NewJsonSerializer()
						err = serializer.UnmarshalClass(jsonBytes, classValue)
						if err != nil {
							return data.NewObjectValue(), nil
						}
						return classValue, nil
					}
				}
			}
		}
	}

	// 如果无法创建类实例，返回空对象
	return data.NewObjectValue(), nil
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
