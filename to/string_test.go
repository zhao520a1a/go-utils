package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := map[string]struct {
		givenObj interface{}
		wantStr  string
	}{
		"int":    {1, "1"},
		"string": {"abc", "abc"},
		"slice":  {[]int{1, 2, 3}, "[1 2 3]"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotStr := String(tt.givenObj)
			assert.Equal(t, tt.wantStr, gotStr)
		})
	}
}

func TestJSON(t *testing.T) {
	tests := map[string]struct {
		givenObj interface{}
		wantStr  string
	}{
		"int slice": {[]int{1, 2, 3}, "[1,2,3]"},
		"map[string]string": {
			map[string]string{
				"name":  "jack",
				"email": "jack@ipalfish.com",
			},
			`{"email":"jack@ipalfish.com","name":"jack"}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotStr, err := JSON(tt.givenObj)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStr, gotStr)
		})
	}
}

func TestLogString(t *testing.T) {
	tests := map[string]struct {
		givenObj interface{}
		wantStr  string
	}{
		"[]int": {[]int{1, 2, 3}, "[1,2,3]"},
		"map[string]string": {
			map[string]string{
				"name":  "jack",
				"email": "jack@ipalfish.com",
			},
			`{"email":"jack@ipalfish.com","name":"jack"}`,
		},
		"string": {"not found", `"not found"`},
		"nil":    {nil, `null`},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotStr := LogString(tt.givenObj)
			assert.Equal(t, tt.wantStr, gotStr)
		})
	}
}
