package data

import "iter"

// Span is a borrowed, structurally read-only view. It does not copy elements or
// expose a mutable slice header. The owner must remain unchanged during use.
type Span[T any] struct{ values []T }

func NewSpan[T any](values []T) Span[T]        { return Span[T]{values: values} }
func (s Span[T]) Len() int                     { return len(s.values) }
func (s Span[T]) At(i int) T                   { return s.values[i] }
func (s Span[T]) Slice(start, end int) Span[T] { return NewSpan(s.values[start:end]) }
func (s Span[T]) AppendTo(dst []T) []T         { return append(dst, s.values...) }
func (s Span[T]) Range() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, value := range s.values {
			if !yield(i, value) {
				return
			}
		}
	}
}
