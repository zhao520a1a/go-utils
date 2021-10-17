package collection

import (
	"github.com/zhao520a1a/go-utils/to"
	. "github.com/zhao520a1a/go-utils/types"
)

// TReverseInPlace reverses the []T in place.
func TReverseInPlace(ts []T) {
	l, m := len(ts), len(ts)/2
	for i := 0; i < m; i++ {
		j := l - i - 1
		ts[i], ts[j] = ts[j], ts[i]
	}
	return
}

// TIntersect computes the intersection of two []T.
func TIntersect(tsa, tsb []T) (ret []T) {
	sa, sb := to.TSet(tsa), to.TSet(tsb)

	if len(sa) > len(sb) {
		sa, sb = sb, sa
	}

	for k := range sa {
		if _, ok := sb[k]; ok {
			ret = append(ret, k)
		}
	}
	return
}

// TForEach applies the function do on each element of []T.
func TForEach(ts []T, do func(i int, t T, raw []T)) {
	for i, t := range ts {
		do(i, t, ts)
	}
}

// TFilter keeps all elements of []T that get true from function filter.
func TFilter(ts []T, filter func(t T) bool) (ret []T) {
	for _, t := range ts {
		if filter(t) {
			ret = append(ret, t)
		}
	}
	return ret
}

// TIn checks whether the element is in []T.
func TIn(ts []T, t T) bool {
	if len(ts) == 0 {
		return false
	}

	for _, tt := range ts {
		if tt == t {
			return true
		}
	}
	return false
}

// TRemove removes the element from []T
func TRemove(ts []T, predicate func(t T) bool) (ret []T, removed bool) {
	for _, tt := range ts {
		if predicate(tt) {
			removed = true
			continue
		}
		ret = append(ret, tt)
	}
	return
}
