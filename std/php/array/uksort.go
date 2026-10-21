package array

import (
	"github.com/php-any/origami/data"
	"sort"
)

type UksortFunction struct{ UasortFunction }

func NewUksortFunction() data.FuncStmt    { return &UksortFunction{} }
func (f *UksortFunction) GetName() string { return "uksort" }
func (f *UksortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	callback, _ := ctx.GetIndexValue(1)
	array, ok := value.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	copy := data.CloneArrayValue(array)
	var failure data.Control
	copy.EditPreservingKeys(func(slots []*data.ZVal) {
		keys := make(map[*data.ZVal]data.Value, len(slots))
		for i, slot := range slots {
			keys[slot] = slot.PHPArrayKey(i)
		}
		sort.SliceStable(slots, func(i, j int) bool {
			if failure != nil {
				return false
			}
			result, ctl := invokeCallback(ctx, callback, []data.Value{keys[slots[i]], keys[slots[j]]})
			if returned, ok := ctl.(data.ReturnControl); ok {
				result, _ = returned.ReturnValue().(data.Value)
				ctl = nil
			}
			if ctl != nil {
				failure = ctl
				return false
			}
			return compareSign(result) < 0
		})
	})
	data.CowAssign(ctx.GetIndexZVal(0), copy)
	if failure != nil {
		return nil, failure
	}
	return data.NewBoolValue(true), nil
}
