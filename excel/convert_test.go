package excel

import (
	"strconv"
	"testing"
)

func TestToBool(t *testing.T) {
	tests := []struct {
		args    interface{}
		wantRes bool
	}{
		{
			"true",
			true,
		},
		{
			"false",
			false,
		},
		{
			"1",
			true,
		},
		{
			"0",
			false,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gotRes, err := ToBool(tt.args)
			if err != nil {
				t.Errorf("ToBool() error = %v ", err)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("ToBool() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		args    interface{}
		wantRes int
	}{
		{
			"1",
			1,
		},
		{
			"-1",
			-1,
		},
		{
			"true",
			1,
		},
		{
			"false",
			0,
		},
		{
			6.666,
			6,
		},
		{
			int64(99999999),
			99999999,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gotRes, err := ToInt(tt.args)
			if err != nil {
				t.Errorf("ToInt() error = %v ", err)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("ToInt() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestToInt8(t *testing.T) {
	tests := []struct {
		args    interface{}
		wantRes int8
	}{
		{
			"1",
			1,
		},
		{
			"-1",
			-1,
		},
		{
			"true",
			1,
		},
		{
			"false",
			0,
		},
		{
			6.666,
			6,
		},
		{
			int64(99999999),
			-1,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gotRes, err := ToInt8(tt.args)
			if err != nil {
				t.Errorf("ToInt8() error = %v ", err)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("ToInt8() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestToUint8(t *testing.T) {
	tests := []struct {
		args    interface{}
		wantRes uint8
	}{
		{
			"1",
			1,
		},
		{
			"-1",
			255,
		},
		{
			"true",
			1,
		},
		{
			"false",
			0,
		},
		{
			6.666,
			6,
		},
		{
			int64(99999999),
			255,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gotRes, err := ToUint8(tt.args)
			if err != nil {
				t.Errorf("ToUint8() error = %v ", err)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("ToUint8() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		args    interface{}
		wantRes float64
	}{
		{
			"1.0",
			1,
		},
		{
			"-1.90",
			-1.90,
		},
		{
			"true",
			1,
		},
		{
			"false",
			0,
		},
		{
			6.6606,
			6.6606,
		},
		{
			int64(99999999),
			99999999,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gotRes, err := ToFloat64(tt.args)
			if err != nil {
				t.Errorf("ToFloat64() error = %v ", err)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("ToFloat64() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
