package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSliceToSet(t *testing.T) {
	tests := map[string]struct {
		givenSlice []string
		wantSet    map[string]interface{}
	}{
		"simple": {
			[]string{"a", "b", "c"},
			map[string]interface{}{
				"a": struct{}{},
				"b": struct{}{},
				"c": struct{}{},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotSet := StringSet(tt.givenSlice)
			assert.Equal(t, tt.wantSet, gotSet)
		})
	}
}
