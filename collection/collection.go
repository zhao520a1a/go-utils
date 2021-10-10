// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package collection

import "github.com/zhao520a1a/go-utils.git/to"

// StringReverseInPlace reverses the []String in place.
func StringReverseInPlace(ts []string) {
	l, m := len(ts), len(ts)/2
	for i := 0; i < m; i++ {
		j := l - i - 1
		ts[i], ts[j] = ts[j], ts[i]
	}
	return
}

// StringIntersect computes the intersection of two []String.
func StringIntersect(tsa, tsb []string) (ret []string) {
	sa, sb := to.StringSet(tsa), to.StringSet(tsb)

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

// StringForEach applies the function do on each element of []String.
func StringForEach(ts []string, do func(i int, t string, raw []string)) {
	for i, t := range ts {
		do(i, t, ts)
	}
}

// StringFilter keeps all elements of []String that get true from function filter.
func StringFilter(ts []string, filter func(t string) bool) (ret []string) {
	for _, t := range ts {
		if filter(t) {
			ret = append(ret, t)
		}
	}
	return ret
}

// StringIn checks whether the element is in []String.
func StringIn(ts []string, t string) bool {
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

// StringRemove removes the element from []String
func StringRemove(ts []string, predicate func(t string) bool) (ret []string, removed bool) {
	for _, tt := range ts {
		if predicate(tt) {
			removed = true
			continue
		}
		ret = append(ret, tt)
	}
	return
}

// IntReverseInPlace reverses the []Int in place.
func IntReverseInPlace(ts []int) {
	l, m := len(ts), len(ts)/2
	for i := 0; i < m; i++ {
		j := l - i - 1
		ts[i], ts[j] = ts[j], ts[i]
	}
	return
}

// IntIntersect computes the intersection of two []Int.
func IntIntersect(tsa, tsb []int) (ret []int) {
	sa, sb := to.IntSet(tsa), to.IntSet(tsb)

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

// IntForEach applies the function do on each element of []Int.
func IntForEach(ts []int, do func(i int, t int, raw []int)) {
	for i, t := range ts {
		do(i, t, ts)
	}
}

// IntFilter keeps all elements of []Int that get true from function filter.
func IntFilter(ts []int, filter func(t int) bool) (ret []int) {
	for _, t := range ts {
		if filter(t) {
			ret = append(ret, t)
		}
	}
	return ret
}

// IntIn checks whether the element is in []Int.
func IntIn(ts []int, t int) bool {
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

// IntRemove removes the element from []Int
func IntRemove(ts []int, predicate func(t int) bool) (ret []int, removed bool) {
	for _, tt := range ts {
		if predicate(tt) {
			removed = true
			continue
		}
		ret = append(ret, tt)
	}
	return
}

// Int64ReverseInPlace reverses the []Int64 in place.
func Int64ReverseInPlace(ts []int64) {
	l, m := len(ts), len(ts)/2
	for i := 0; i < m; i++ {
		j := l - i - 1
		ts[i], ts[j] = ts[j], ts[i]
	}
	return
}

// Int64Intersect computes the intersection of two []Int64.
func Int64Intersect(tsa, tsb []int64) (ret []int64) {
	sa, sb := to.Int64Set(tsa), to.Int64Set(tsb)

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

// Int64ForEach applies the function do on each element of []Int64.
func Int64ForEach(ts []int64, do func(i int, t int64, raw []int64)) {
	for i, t := range ts {
		do(i, t, ts)
	}
}

// Int64Filter keeps all elements of []Int64 that get true from function filter.
func Int64Filter(ts []int64, filter func(t int64) bool) (ret []int64) {
	for _, t := range ts {
		if filter(t) {
			ret = append(ret, t)
		}
	}
	return ret
}

// Int64In checks whether the element is in []Int64.
func Int64In(ts []int64, t int64) bool {
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

// Int64Remove removes the element from []Int64
func Int64Remove(ts []int64, predicate func(t int64) bool) (ret []int64, removed bool) {
	for _, tt := range ts {
		if predicate(tt) {
			removed = true
			continue
		}
		ret = append(ret, tt)
	}
	return
}

// Float64ReverseInPlace reverses the []Float64 in place.
func Float64ReverseInPlace(ts []float64) {
	l, m := len(ts), len(ts)/2
	for i := 0; i < m; i++ {
		j := l - i - 1
		ts[i], ts[j] = ts[j], ts[i]
	}
	return
}

// Float64Intersect computes the intersection of two []Float64.
func Float64Intersect(tsa, tsb []float64) (ret []float64) {
	sa, sb := to.Float64Set(tsa), to.Float64Set(tsb)

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

// Float64ForEach applies the function do on each element of []Float64.
func Float64ForEach(ts []float64, do func(i int, t float64, raw []float64)) {
	for i, t := range ts {
		do(i, t, ts)
	}
}

// Float64Filter keeps all elements of []Float64 that get true from function filter.
func Float64Filter(ts []float64, filter func(t float64) bool) (ret []float64) {
	for _, t := range ts {
		if filter(t) {
			ret = append(ret, t)
		}
	}
	return ret
}

// Float64In checks whether the element is in []Float64.
func Float64In(ts []float64, t float64) bool {
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

// Float64Remove removes the element from []Float64
func Float64Remove(ts []float64, predicate func(t float64) bool) (ret []float64, removed bool) {
	for _, tt := range ts {
		if predicate(tt) {
			removed = true
			continue
		}
		ret = append(ret, tt)
	}
	return
}

// ByteReverseInPlace reverses the []Byte in place.
func ByteReverseInPlace(ts []byte) {
	l, m := len(ts), len(ts)/2
	for i := 0; i < m; i++ {
		j := l - i - 1
		ts[i], ts[j] = ts[j], ts[i]
	}
	return
}

// ByteIntersect computes the intersection of two []Byte.
func ByteIntersect(tsa, tsb []byte) (ret []byte) {
	sa, sb := to.ByteSet(tsa), to.ByteSet(tsb)

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

// ByteForEach applies the function do on each element of []Byte.
func ByteForEach(ts []byte, do func(i int, t byte, raw []byte)) {
	for i, t := range ts {
		do(i, t, ts)
	}
}

// ByteFilter keeps all elements of []Byte that get true from function filter.
func ByteFilter(ts []byte, filter func(t byte) bool) (ret []byte) {
	for _, t := range ts {
		if filter(t) {
			ret = append(ret, t)
		}
	}
	return ret
}

// ByteIn checks whether the element is in []Byte.
func ByteIn(ts []byte, t byte) bool {
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

// ByteRemove removes the element from []Byte
func ByteRemove(ts []byte, predicate func(t byte) bool) (ret []byte, removed bool) {
	for _, tt := range ts {
		if predicate(tt) {
			removed = true
			continue
		}
		ret = append(ret, tt)
	}
	return
}

// InterfaceReverseInPlace reverses the []Interface in place.
func InterfaceReverseInPlace(ts []interface{}) {
	l, m := len(ts), len(ts)/2
	for i := 0; i < m; i++ {
		j := l - i - 1
		ts[i], ts[j] = ts[j], ts[i]
	}
	return
}

// InterfaceIntersect computes the intersection of two []Interface.
func InterfaceIntersect(tsa, tsb []interface{}) (ret []interface{}) {
	sa, sb := to.InterfaceSet(tsa), to.InterfaceSet(tsb)

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

// InterfaceForEach applies the function do on each element of []Interface.
func InterfaceForEach(ts []interface{}, do func(i int, t interface{}, raw []interface{})) {
	for i, t := range ts {
		do(i, t, ts)
	}
}

// InterfaceFilter keeps all elements of []Interface that get true from function filter.
func InterfaceFilter(ts []interface{}, filter func(t interface{}) bool) (ret []interface{}) {
	for _, t := range ts {
		if filter(t) {
			ret = append(ret, t)
		}
	}
	return ret
}

// InterfaceIn checks whether the element is in []Interface.
func InterfaceIn(ts []interface{}, t interface{}) bool {
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

// InterfaceRemove removes the element from []Interface
func InterfaceRemove(ts []interface{}, predicate func(t interface{}) bool) (ret []interface{}, removed bool) {
	for _, tt := range ts {
		if predicate(tt) {
			removed = true
			continue
		}
		ret = append(ret, tt)
	}
	return
}
