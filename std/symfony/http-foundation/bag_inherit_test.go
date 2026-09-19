package httpfoundation

import "testing"

func TestServerBagInheritsParameterBagGet(t *testing.T) {
	c := NewServerBagClass()
	if _, ok := c.GetMethod("get"); !ok {
		t.Fatal("ServerBag missing inherited ParameterBag::get")
	}
	if _, ok := c.GetMethod("has"); !ok {
		t.Fatal("ServerBag missing inherited ParameterBag::has")
	}
	if _, ok := c.GetMethod("all"); !ok {
		t.Fatal("ServerBag missing inherited ParameterBag::all")
	}
}

func TestInputBagInheritsParameterBagHas(t *testing.T) {
	c := NewInputBagClass()
	if _, ok := c.GetMethod("has"); !ok {
		t.Fatal("InputBag missing inherited ParameterBag::has")
	}
}

func TestFileBagInheritsParameterBagGet(t *testing.T) {
	c := NewFileBagClass()
	if _, ok := c.GetMethod("get"); !ok {
		t.Fatal("FileBag missing inherited ParameterBag::get")
	}
}

func TestResponseHeaderBagInheritsHeaderBagGet(t *testing.T) {
	c := NewResponseHeaderBagClass()
	if _, ok := c.GetMethod("get"); !ok {
		t.Fatal("ResponseHeaderBag missing inherited HeaderBag::get")
	}
}
