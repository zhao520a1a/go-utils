package xsqlutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOmitAllEmpty(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		givenMap map[string]interface{}
		wantMap  map[string]interface{}
	}{
		"primitive": {
			map[string]interface{}{
				"v1": 0,
				"v2": "hello world",
				"v3": 0.3,
				"v4": "",
			},
			map[string]interface{}{
				"v2": "hello world",
				"v3": 0.3,
			},
		},
		"time": {
			map[string]interface{}{
				"t1": time.Time{},
				"t2": time.Time{},
				"t3": now,
			},
			map[string]interface{}{
				"t3": now,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantMap, OmitAllEmpty(tc.givenMap))
		})
	}
}
