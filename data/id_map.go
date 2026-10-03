package data

import "slices"

// IDMap is a persistent two-level paged table. Publishing copies one value
// page, one 64-page branch and the small root directory. Sparse symbol IDs
// must not make every request declaration copy the entire page directory.
type IDMap[T any] struct{ branches []*idBranch[T] }
type idBranch[T any] [64]*idPage[T]
type idPage[T any] struct {
	values  [64]T
	present uint64
}

func (m IDMap[T]) Get(id uint32) (T, bool) {
	page, index := id/64, id%64
	branch, offset := page/64, page%64
	if int(branch) < len(m.branches) {
		if b := m.branches[branch]; b != nil {
			if p := b[offset]; p != nil && p.present&(uint64(1)<<index) != 0 {
				return p.values[index], true
			}
		}
	}
	var zero T
	return zero, false
}
func (m IDMap[T]) With(id uint32, value T) IDMap[T] {
	page, index := id/64, id%64
	branch, offset := page/64, page%64
	length := max(len(m.branches), int(branch)+1)
	directory := make([]*idBranch[T], length)
	copy(directory, m.branches)
	b := new(idBranch[T])
	if directory[branch] != nil {
		*b = *directory[branch]
	}
	p := new(idPage[T])
	if b[offset] != nil {
		*p = *b[offset]
	}
	p.values[index] = value
	p.present |= uint64(1) << index
	b[offset] = p
	directory[branch] = b
	return IDMap[T]{branches: directory}
}

// Set is an immutable sorted nominal-ID set; membership allocates nothing.
type Set[T ~uint32] struct{ values []T }

func NewSet[T ~uint32](values ...T) Set[T] {
	copy := slices.Clone(values)
	slices.Sort(copy)
	return Set[T]{slices.Compact(copy)}
}
func (s Set[T]) Contains(value T) bool {
	_, found := slices.BinarySearch(s.values, value)
	return found
}
func (s Set[T]) View() Span[T] { return NewSpan(s.values) }
