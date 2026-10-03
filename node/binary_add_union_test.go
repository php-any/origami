package node

import (
	"testing"

	"github.com/php-any/origami/data"
)

func TestMergeArrayUnionKeepsStringKeys(t *testing.T) {
	left := data.NewArrayValueFromSlots([]*data.ZVal{data.NewNamedZVal("a", data.NewStringValue("login"))})
	right := data.NewArrayValueFromSlots([]*data.ZVal{})
	out := mergeArrayUnion(left, right)
	if out.Len() != 1 {
		t.Fatalf("len=%d", out.Len())
	}
	if out.At(0).Name != "a" {
		t.Fatalf("name=%q want a", out.At(0).Name)
	}
	if out.At(0).ReadValue().AsString() != "login" {
		t.Fatalf("val=%q", out.At(0).ReadValue().AsString())
	}
}

func TestMergeArrayUnionRightAddsMissing(t *testing.T) {
	left := data.NewArrayValueFromSlots([]*data.ZVal{
		data.NewZVal(data.NewStringValue("zero")),
		data.NewNamedZVal("a", data.NewStringValue("keep")),
	})
	right := data.NewArrayValueFromSlots([]*data.ZVal{
		data.NewZVal(data.NewStringValue("drop")),
		data.NewNamedZVal("b", data.NewStringValue("add")),
		data.NewNamedZVal("1", data.NewStringValue("one")),
	})
	out := mergeArrayUnion(left, right)
	got := map[string]string{}
	for arraySlots15, i := out.View(), 0; i < arraySlots15.Len(); i++ {
		z := arraySlots15.At(i)
		got[arraySlotKey(z, i)] = z.ReadValue().AsString()
	}
	if got["a"] != "keep" || got["0"] != "zero" || got["b"] != "add" || got["1"] != "one" {
		t.Fatalf("got=%v", got)
	}
}

func TestObjectPlusEmptyArrayKeepsKeys(t *testing.T) {
	obj := data.NewObjectValue()
	_ = obj.SetProperty("a", data.NewStringValue("login"))
	out := mergeArrayUnion(objectToNamedArray(obj), &data.ArrayValue{})
	if out.Len() != 1 || out.At(0).Name != "a" || out.At(0).ReadValue().AsString() != "login" {
		t.Fatalf("out=%v", out.Snapshot())
	}
}
