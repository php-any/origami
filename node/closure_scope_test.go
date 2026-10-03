package node

import (
	"github.com/php-any/origami/data"
	"testing"
)

func TestStaticMethodClosureHasNoRequestObjectIdentity(t *testing.T) {
	declaration := &ClassStatement{Name: "MagicActions"}
	unrelated := &ClassStatement{Name: "ReleaseTokens"}
	owner := &data.ClassMethodContext{ClassValue: &data.ClassValue{Class: declaration}, SelfClass: declaration, StaticClass: declaration}
	closure := &LambdaExpression{ctx: owner}
	if len(closure.RequestScopeObjects()) != 0 {
		t.Fatal("static method scope exposed as a request object")
	}
	// A nil identity must never rebind the lexical class to another feature.
	replacement := &data.ClassValue{Class: unrelated}
	bound := closure.BindRequestScope(nil, map[*data.ObjectValue]*data.ClassValue{nil: replacement}).(*LambdaExpression)
	if bound.ctx != owner {
		t.Fatal("static method closure changed declaration scope")
	}
}

func TestRequestCaptureGraphPreservesNestedAliasesAndReferences(t *testing.T) {
	source := &data.ClassValue{Class: &ClassStatement{Name: "Asset"}, ObjectValue: data.NewObjectValue()}
	target := &data.ClassValue{Class: source.Class, ObjectValue: data.NewObjectValue()}
	assets := data.NewArrayValue([]data.Value{source, data.NewThisValue(source)})
	shared := data.NewZVal(assets)
	inner := &LambdaExpression{captured: map[int]data.Value{0: assets}, capturedRefs: map[int]*data.ZVal{1: shared}}
	outer := &LambdaExpression{captured: map[int]data.Value{0: data.NewFuncValue(inner)}, capturedRefs: map[int]*data.ZVal{1: shared}}
	if objects := outer.RequestScopeObjects(); len(objects) != 1 || objects[0] != source {
		t.Fatalf("nested captures: %v", objects)
	}
	scope := data.NewRequestCaptureScope(nil, map[*data.ObjectValue]*data.ClassValue{source.ObjectValue: target})
	bound := outer.BindRequestCapture(scope).(*LambdaExpression)
	nested := bound.captured[0].(*data.FuncValue).Value.(*LambdaExpression)
	if bound.capturedRefs[1] != nested.capturedRefs[1] || bound.capturedRefs[1] == shared {
		t.Fatal("request capture reference aliases were lost or reused the original slot")
	}
	array := nested.captured[0].(*data.ArrayValue)
	if array == assets || bound.capturedRefs[1].ReadValue() != array || array.At(0).ReadValue() != target || array.At(1).ReadValue().(*data.ThisValue).ClassValue != target {
		t.Fatal("nested array object aliases were not mapped to the request")
	}
	// A separately registered callback belongs to the same request graph.
	other := (&LambdaExpression{capturedRefs: map[int]*data.ZVal{0: shared}}).BindRequestCapture(scope).(*LambdaExpression)
	if other.capturedRefs[0] != bound.capturedRefs[1] {
		t.Fatal("separate callbacks lost their shared PHP reference")
	}
	bound.capturedRefs[1].StoreRaw(data.NewIntValue(7))
	if shared.ReadValue() != assets || other.capturedRefs[0].ReadValue().AsString() != "7" {
		t.Fatal("request reference state escaped or aliases diverged")
	}
	second := outer.BindRequestCapture(data.NewRequestCaptureScope(nil, map[*data.ObjectValue]*data.ClassValue{source.ObjectValue: target})).(*LambdaExpression)
	if second.capturedRefs[1] == bound.capturedRefs[1] {
		t.Fatal("requests shared a capture slot")
	}
}

func TestRequestCaptureGraphPreservesCyclesAndArrayKeyHistory(t *testing.T) {
	array := data.NewArrayValue(nil).(*data.ArrayValue)
	array.SetKey(data.NewIntValue(100), data.NewIntValue(1))
	array.UnsetKey(data.NewIntValue(100))
	closure := &LambdaExpression{captured: map[int]data.Value{0: array}}
	callback := data.NewFuncValue(closure)
	array.SetStringKey("callback", callback)
	scope := data.NewRequestCaptureScope(nil, nil)
	bound := scope.Bind(callback).(*data.FuncValue)
	copy := bound.Value.(*LambdaExpression).captured[0].(*data.ArrayValue)
	if slot, ok := copy.LookupZValByStringKey("callback"); !ok || slot.ReadValue() != bound {
		t.Fatal("cyclic closure capture lost its identity")
	}
	copy.AppendValue(data.NewIntValue(2))
	if slot, _ := copy.FindSlotByIntKey(101); slot == nil || array.Len() != 1 {
		t.Fatal("request copy lost array key history or changed the source")
	}
}
