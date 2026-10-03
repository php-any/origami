package data

import (
	"sync"
	"testing"
)

type policyCloneFixture struct {
	count int
	self  *policyCloneFixture
	value Value
}
type policyNewFixture struct{ count int }
type policyReadFixture struct{ name string }
type policyMissingFixture struct{}

func init() {
	RegisterNativeRequestPolicy[*policyCloneFixture](NativeClone, func(scope *RequestObjectScope, source *policyCloneFixture) *policyCloneFixture {
		clone := &policyCloneFixture{count: source.count}
		scope.RememberNativeState(source, clone)
		clone.self = scope.BindNativeState(source.self).(*policyCloneFixture)
		clone.value = scope.Bind(source.value)
		return clone
	})
	RegisterNativeRequestPolicy[*policyNewFixture](NativeNew, func(*RequestObjectScope, *policyNewFixture) *policyNewFixture { return &policyNewFixture{} })
	RegisterNativeRequestPolicy[*policyReadFixture](NativeReadOnly, nil)
}

func TestNativePoliciesAliasesCyclesAndConcurrentIsolation(t *testing.T) {
	array := NewArrayValue([]Value{NewIntValue(1)}).(*ArrayValue)
	original := &policyCloneFixture{count: 4, value: array}
	original.self = original
	fresh := &policyNewFixture{count: 9}
	readonly := &policyReadFixture{name: "startup"}
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			scope := NewRequestObjectScope(nil)
			clone := scope.BindNativeState(original).(*policyCloneFixture)
			if clone == original || clone.self != clone || scope.BindNativeState(original) != clone {
				t.Error("native graph identity lost")
			}
			clone.count++
			clone.value.(*ArrayValue).SetIntKey(0, NewIntValue(8))
			if scope.BindNativeState(fresh).(*policyNewFixture).count != 0 {
				t.Error("new policy copied retained state")
			}
			if scope.BindNativeState(readonly) != readonly {
				t.Error("read-only state copied")
			}
		}()
	}
	group.Wait()
	if original.count != 4 || array.At(0).ReadValue().AsString() != "1" {
		t.Fatal("retained native state mutated")
	}
}

func TestNativePolicyFailureCannotReturnPartialObjects(t *testing.T) {
	scope := NewRequestObjectScope(nil)
	var missing *policyMissingFixture
	if scope.BindNativeState(missing) != missing {
		t.Fatal("typed nil must remain absent")
	}
	for i := 0; i < 2; i++ {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("unknown native state accepted")
				}
			}()
			scope.BindNativeState(&policyMissingFixture{})
		}()
	}
}
