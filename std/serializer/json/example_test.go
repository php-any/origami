package json

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/php-any/origami/data"
)

func TestJsonSerializer(t *testing.T) {
	// 测试 IntValue
	intVal := data.NewIntValue(42)
	serializer := NewJsonSerializer()

	// 序列化
	data, err := intVal.Marshal(serializer)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 验证 JSON 格式
	var result int
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON 解析失败: %v", err)
	}
	if result != 42 {
		t.Fatalf("期望 42，得到 %d", result)
	}

	// 反序列化
	newIntVal := &data.IntValue{}
	if err := newIntVal.Unmarshal(data, serializer); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}
	if newIntVal.Value != 42 {
		t.Fatalf("期望 42，得到 %d", newIntVal.Value)
	}

	fmt.Printf("IntValue 测试通过: %d -> %s -> %d\n", intVal.Value, string(data), newIntVal.Value)
}

func TestJsonSerializerArray(t *testing.T) {
	// 测试 ArrayValue
	arrVal := data.NewArrayValue([]data.Value{
		data.NewIntValue(1),
		data.NewStringValue("hello"),
		data.NewBoolValue(true),
	})
	serializer := NewJsonSerializer()

	// 序列化
	data, err := arrVal.Marshal(serializer)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 验证 JSON 格式
	var result []interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON 解析失败: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("期望长度 3，得到 %d", len(result))
	}

	// 反序列化
	newArrVal := &data.ArrayValue{}
	if err := newArrVal.Unmarshal(data, serializer); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}
	if len(newArrVal.Value) != 3 {
		t.Fatalf("期望长度 3，得到 %d", len(newArrVal.Value))
	}

	fmt.Printf("ArrayValue 测试通过: %s -> %s\n", arrVal.AsString(), string(data))
}
