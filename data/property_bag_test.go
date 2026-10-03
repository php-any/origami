package data

import "testing"

func TestPropertyBagCannotEscapeAsPHPValue(t *testing.T) {
	bag := NewPropertyBag()
	bag.SetProperty("empty", NewNullValue())
	if _, ok := any(bag).(Value); ok {
		t.Fatal("internal property storage is a PHP value")
	}
	if _, ok := any(bag).(GetValue); ok {
		t.Fatal("internal property storage is executable")
	}
	if !bag.HasProperty("empty") {
		t.Fatal("null property lost its declaration")
	}
	object := NewStdClassValue(nil)
	if ValueKindOf(object) != ValueObject || !TypeObject.Matches(object, nil) || TypeArray.Matches(object, nil) {
		t.Fatal("stdClass identity is not an object")
	}
}
