package to

import . "github.com/zhao520a1a/go-utils/types"

// TSet converts []T to map[T]interface{}, with all values equal
// to struct{}{}, to simulate a Set.
func TSet(ts []T) (s map[T]interface{}) {
	s = make(map[T]interface{}, len(ts))
	for _, t := range ts {
		s[t] = struct{}{}
	}
	return
}
