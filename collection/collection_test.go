package collection

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringFilter(t *testing.T) {
	tests := map[string]struct {
		givenSlice  []string
		givenFilter func(string) bool
		wantSlice   []string
	}{
		"simple": {
			[]string{"a", "b", "c"},
			func(ele string) bool { return ele != "a" },
			[]string{"b", "c"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotSlice := StringFilter(tt.givenSlice, tt.givenFilter)
			assert.Equal(t, tt.wantSlice, gotSlice)
		})
	}
}

type item struct {
	I   int
	Ele string
	Raw []string
}

type recorder struct {
	Items []item
}

func (r *recorder) do(i int, ele string, raw []string) {
	r.Items = append(r.Items, item{i, ele, raw})
}

func TestStringForEach(t *testing.T) {
	tests := map[string]struct {
		givenSlice    []string
		givenRecorder *recorder
		wantItems     []item
	}{
		"simple": {
			[]string{"a", "b", "c"},
			&recorder{},
			[]item{
				{0, "a", []string{"a", "b", "c"}},
				{1, "b", []string{"a", "b", "c"}},
				{2, "c", []string{"a", "b", "c"}},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			StringForEach(tt.givenSlice, tt.givenRecorder.do)
			assert.Equal(t, tt.wantItems, tt.givenRecorder.Items)
		})
	}
}

func TestReverseStringSliceInPlace(t *testing.T) {
	tests := map[string]struct {
		givenSlice []string
		wantSlice  []string
	}{
		"simple": {
			[]string{"a", "b", "c"},
			[]string{"c", "b", "a"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			StringReverseInPlace(tt.givenSlice)
			assert.Equal(t, tt.wantSlice, tt.givenSlice)
		})
	}
}

func TestIntersectStringSlice(t *testing.T) {
	tests := map[string]struct {
		givenSliceA []string
		givenSliceB []string
		wantSlice   []string
	}{
		"same size": {
			[]string{"a", "b", "e", "h", "k"},
			[]string{"b", "c", "d", "h"},
			[]string{"b", "h"},
		},
		"small and large": {
			[]string{"a"},
			[]string{"a", "b", "c", "d", "e", "f", "g"},
			[]string{"a"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotSlice := StringIntersect(tt.givenSliceA, tt.givenSliceB)
			sort.Strings(gotSlice)
			assert.Equal(t, tt.wantSlice, gotSlice)
		})
	}
}

func TestStringIn(t *testing.T) {
	tests := map[string]struct {
		givenTS []string
		givenT  string
		wantRet bool
	}{
		"empty ts":                      {[]string{}, "a", false},
		"empty t":                       {[]string{"a", "b"}, "", false},
		"basic in":                      {[]string{"a", "b"}, "a", true},
		"basic not in":                  {[]string{"a", "b"}, "c", false},
		"in with repeated elements":     {[]string{"a", "a"}, "a", true},
		"not in with repeated elements": {[]string{"a", "a"}, "b", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRet := StringIn(tt.givenTS, tt.givenT)
			assert.Equal(t, tt.wantRet, gotRet)
		})
	}
}

func TestStringRemove(t *testing.T) {
	tests := map[string]struct {
		givenTS        []string
		givenPredicate func(string) bool
		wantRet        []string
		wantRemoved    bool
	}{
		"empty ts": {[]string{}, nil, nil, false},
		"remove one": {
			[]string{"a", "b", "c"},
			func(s string) bool {
				return s == "a"
			},
			[]string{"b", "c"},
			true,
		},
		"remove multiple": {
			[]string{"a", "b", "c", "d"},
			func(s string) bool {
				return s > "a"
			},
			[]string{"a"},
			true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ret, removed := StringRemove(tc.givenTS, tc.givenPredicate)
			assert.Equal(t, tc.wantRet, ret)
			assert.Equal(t, tc.wantRemoved, removed)
		})
	}
}
