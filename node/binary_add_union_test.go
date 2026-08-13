package node

import (
	"testing"

	"github.com/php-any/origami/data"
)

func TestMergeArrayUnionKeepsStringKeys(t *testing.T) {
	left := &data.ArrayValue{List: []*data.ZVal{data.NewNamedZVal("a", data.NewStringValue("login"))}}
	right := &data.ArrayValue{List: []*data.ZVal{}}
	out := mergeArrayUnion(left, right)
	if len(out.List) != 1 {
		t.Fatalf("len=%d", len(out.List))
	}
	if out.List[0].Name != "a" {
		t.Fatalf("name=%q want a", out.List[0].Name)
	}
	if out.List[0].Value.AsString() != "login" {
		t.Fatalf("val=%q", out.List[0].Value.AsString())
	}
}

func TestMergeArrayUnionRightAddsMissing(t *testing.T) {
	left := &data.ArrayValue{List: []*data.ZVal{
		data.NewZVal(data.NewStringValue("zero")),
		data.NewNamedZVal("a", data.NewStringValue("keep")),
	}}
	right := &data.ArrayValue{List: []*data.ZVal{
		data.NewZVal(data.NewStringValue("drop")),
		data.NewNamedZVal("b", data.NewStringValue("add")),
		data.NewNamedZVal("1", data.NewStringValue("one")),
	}}
	out := mergeArrayUnion(left, right)
	got := map[string]string{}
	for i, z := range out.List {
		got[arraySlotKey(z, i)] = z.Value.AsString()
	}
	if got["a"] != "keep" || got["0"] != "zero" || got["b"] != "add" || got["1"] != "one" {
		t.Fatalf("got=%v", got)
	}
}

func TestObjectPlusEmptyArrayKeepsKeys(t *testing.T) {
	obj := data.NewObjectValue()
	_ = obj.SetProperty("a", data.NewStringValue("login"))
	out := mergeArrayUnion(objectToNamedArray(obj), &data.ArrayValue{})
	if len(out.List) != 1 || out.List[0].Name != "a" || out.List[0].Value.AsString() != "login" {
		t.Fatalf("out=%v", out.List)
	}
}
