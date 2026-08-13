package array

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	sortAsc  = 4
	sortDesc = 3
)

type multisortSpec struct {
	arr   *data.ArrayValue
	order int
	flags int
}

// ArrayMultisortFunction 实现 array_multisort
// array_multisort(array &$array, mixed $sort_order = SORT_ASC, mixed $sort_flags = SORT_REGULAR, ...): bool
type ArrayMultisortFunction struct{}

func NewArrayMultisortFunction() data.FuncStmt {
	return &ArrayMultisortFunction{}
}

func (f *ArrayMultisortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	callArgs := ctx.GetCallArgs()
	if len(callArgs) == 0 {
		return data.NewBoolValue(false), nil
	}

	specs := make([]multisortSpec, 0)
	i := 0
	for i < len(callArgs) {
		arr, ok, ctl := resolveArrayRef(ctx, callArgs[i])
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			i++
			continue
		}
		spec := multisortSpec{arr: arr, order: sortAsc, flags: 0}
		i++
		for i < len(callArgs) {
			val, acl := callArgs[i].GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			v, ok := val.(data.Value)
			if !ok {
				break
			}
			if order, ok := readSortOrder(v); ok {
				spec.order = order
				i++
				continue
			}
			if flags, ok := readSortFlags(v); ok {
				spec.flags = flags
				i++
				continue
			}
			break
		}
		specs = append(specs, spec)
	}

	if len(specs) == 0 {
		return data.NewBoolValue(false), nil
	}

	n := len(specs[0].arr.List)
	if n <= 1 {
		return data.NewBoolValue(true), nil
	}

	indices := make([]int, n)
	for j := range indices {
		indices[j] = j
	}

	sort.SliceStable(indices, func(a, b int) bool {
		ia, ib := indices[a], indices[b]
		for _, spec := range specs {
			if ia >= len(spec.arr.List) || ib >= len(spec.arr.List) {
				continue
			}
			va := spec.arr.List[ia].Value
			vb := spec.arr.List[ib].Value
			if valuesEqual(va, vb) {
				continue
			}
			less := compareValues(va, vb, spec.flags)
			if spec.order == sortDesc {
				return !less
			}
			return less
		}
		return ia < ib
	})

	for _, spec := range specs {
		reorderArrayInPlace(spec.arr, indices)
	}
	return data.NewBoolValue(true), nil
}

func resolveArrayRef(ctx data.Context, arg data.GetValue) (*data.ArrayValue, bool, data.Control) {
	if arg == nil {
		return nil, false, nil
	}
	if v, ok := arg.(data.Variable); ok {
		zv := ctx.GetIndexZVal(v.GetIndex())
		if zv == nil {
			return nil, false, nil
		}
		if arr, ok := zv.Value.(*data.ArrayValue); ok {
			return arr, true, nil
		}
	}
	type zvalGetter interface {
		GetZVal(data.Context) (*data.ZVal, data.Control)
	}
	if g, ok := arg.(zvalGetter); ok {
		zv, ctl := g.GetZVal(ctx)
		if ctl != nil {
			return nil, false, ctl
		}
		if zv != nil {
			if arr, ok := zv.Value.(*data.ArrayValue); ok {
				return arr, true, nil
			}
		}
	}
	val, ctl := arg.GetValue(ctx)
	if ctl != nil {
		return nil, false, ctl
	}
	if arr, ok := val.(*data.ArrayValue); ok {
		return arr, true, nil
	}
	return nil, false, nil
}

func readSortOrder(v data.Value) (int, bool) {
	iv, ok := v.(data.AsInt)
	if !ok {
		return 0, false
	}
	n, err := iv.AsInt()
	if err != nil {
		return 0, false
	}
	if n == sortAsc || n == sortDesc {
		return n, true
	}
	return 0, false
}

func readSortFlags(v data.Value) (int, bool) {
	iv, ok := v.(data.AsInt)
	if !ok {
		return 0, false
	}
	n, err := iv.AsInt()
	if err != nil {
		return 0, false
	}
	switch n {
	case 0, 1, 2, 3, 5, 6:
		return n, true
	default:
		return 0, false
	}
}

func reorderArrayInPlace(arr *data.ArrayValue, order []int) {
	old := arr.List
	newList := make([]*data.ZVal, len(order))
	for i, idx := range order {
		if idx < len(old) {
			newList[i] = old[idx]
		} else {
			newList[i] = data.NewZVal(data.NewNullValue())
		}
	}
	arr.List = newList
}

func (f *ArrayMultisortFunction) GetName() string { return "array_multisort" }

func (f *ArrayMultisortFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewCallerContextParameter(nil)}
}

func (f *ArrayMultisortFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}
