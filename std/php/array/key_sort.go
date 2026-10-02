package array

import (
	"github.com/php-any/origami/data"
	"sort"
)

func sortArrayKeys(a *data.ArrayValue, flags int, descending bool) {
	a.EditPreservingKeys(func(slots []*data.ZVal) {
		sort.SliceStable(slots, func(i, j int) bool {
			left, right := slots[i].PHPArrayKey(i), slots[j].PHPArrayKey(j)
			if descending {
				left, right = right, left
			}
			if flags == 0 {
				return data.Compare(left, right) < 0
			}
			return compareValues(left, right, flags)
		})
	})
}
