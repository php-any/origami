package data

import (
	"fmt"
	"reflect"
	"sync"
)

type NativeRequestPolicy uint8

const (
	NativeReadOnly NativeRequestPolicy = iota + 1
	NativeClone
	NativeNew
)

type nativePolicy struct {
	mode NativeRequestPolicy
	bind func(*RequestObjectScope, any) any
}

// Registration is a startup operation. Unknown retained native state is
// rejected instead of becoming an implicit mutable worker singleton.
var nativeRequestPolicies sync.Map // reflect.Type -> nativePolicy

func RegisterNativeRequestPolicy[T any](mode NativeRequestPolicy, bind func(*RequestObjectScope, T) T) {
	typeOf := reflect.TypeFor[T]()
	if mode < NativeReadOnly || mode > NativeNew {
		panic("invalid native request policy")
	}
	if mode != NativeReadOnly && bind == nil {
		panic("native request clone/new policy requires a binder")
	}
	policy := nativePolicy{mode: mode}
	if bind != nil {
		policy.bind = func(scope *RequestObjectScope, value any) any { return bind(scope, value.(T)) }
	}
	if _, loaded := nativeRequestPolicies.LoadOrStore(typeOf, policy); loaded {
		panic("duplicate native request policy for " + typeOf.String())
	}
}

func (s *RequestObjectScope) BindNativeState(source any) any {
	if s.failure != nil {
		panic(s.failure)
	}
	defer func() {
		if caught := recover(); caught != nil {
			s.failure = caught
			panic(caught)
		}
	}()
	if source == nil {
		return nil
	}
	switch source.(type) {
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		return source
	}
	typeOf := reflect.TypeOf(source)
	value := reflect.ValueOf(source)
	switch value.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface:
		if value.IsNil() {
			return source
		}
	}
	policy, ok := nativeRequestPolicies.Load(typeOf)
	if !ok {
		panic(NewErrorThrow(nil, fmt.Errorf("Retained native state %s has no request isolation policy", typeOf)))
	}
	entry := policy.(nativePolicy)
	if entry.mode == NativeReadOnly {
		return source
	}
	if typeOf.Comparable() {
		if previous, ok := s.native[source]; ok {
			return previous
		}
	}
	target := entry.bind(s, source)
	if typeOf.Comparable() {
		s.native[source] = target
	}
	return target
}

func (s *RequestObjectScope) RememberNativeState(source, target any) {
	if source != nil && reflect.TypeOf(source).Comparable() {
		s.native[source] = target
	}
}
