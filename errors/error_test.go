package errors

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestE(t *testing.T) {
	tests := map[string]struct {
		givenCode  Class
		givenMsg   string
		givenOp    Op
		givenCause error
		wantError  error
	}{
		"simple": {
			Internal, "data inconsistent", Op("SetUser"), nil,
			&Error{
				Class: Internal,
				Msg:   "data inconsistent",
				Op:    Op("SetUser"),
				Cause: nil,
			},
		},
		"nested": {
			NotFound, "user joe not found", Op("GetUser"), sql.ErrNoRows,
			&Error{
				Class: NotFound,
				Msg:   "user joe not found",
				Op:    Op("GetUser"),
				Cause: sql.ErrNoRows,
			},
		},
		"msg only": {
			Class(""), "user joe not found", Op("GetUser"), nil,
			&Error{
				Msg: "user joe not found",
				Op:  Op("GetUser"),
			},
		},
		"dedup": {
			Internal, "data inconsistent", Op("SetUser"), &Error{Op: Op("SetUser")},
			&Error{
				Op: Op("SetUser"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantError, E(tc.givenOp, tc.givenCode, tc.givenMsg, tc.givenCause))
		})
	}

	t.Run("panic with no valid arguments", func(t *testing.T) {
		assert.Panics(t, func() { E() })
		assert.Panics(t, func() { E(nil) })
		assert.Panics(t, func() { E(nil, nil) })
	})

	t.Run("invalid type", func(t *testing.T) {
		assert.True(t, strings.Contains(E(bytes.NewBuffer(nil)).Error(), "unknown type"))
	})

	t.Run("int32", func(t *testing.T) {
		assert.Equal(t, &Error{Msg: "foo", Code: -1}, E(int32(-1), "foo"))
	})
}

func TestIs(t *testing.T) {
	tests := map[string]struct {
		givenErr  error
		givenCode Class
		wantIs    bool
	}{
		"nil error": {
			nil, Internal, false,
		},
		"simple is": {
			E(Op("DelUser"), Internal, "user joe not found"),
			Internal, true,
		},
		"simple is not": {
			E(Op("DelUser"), Internal, "user joe not found"),
			Invalid, false,
		},
		"types other than *Error": {
			errors.New("invalid username"), Invalid, false,
		},
		"wrapped is": {
			E(Op("HandleDelUser"), "userService error",
				E(Op("DelUser"), NotFound, "user joe not found")),
			NotFound, true,
		},
		"wrapped is not": {
			E(Op("HandleDelUser"), "userService error",
				E(Op("DelUser"), NotFound, "user joe not found")),
			Invalid, false,
		},
		"first non-empty is not": {
			E(Op("HandleDelUser"), Internal, "userService error",
				E(Op("DelUser"), NotFound, "user joe not found")),
			NotFound, false,
		},
		"is empty": {
			E(Op("HandleDelUser"), "userService error"),
			"", true,
		},
		"wrapped is empty": {
			E(Op("HandleDelUser"), "userService error",
				E(Op("DelUser"), "user joe not found")),
			"", true,
		},
		"is not empty": {
			E(Op("HandleDelUser"), "userService error",
				E(Op("DelUser"), Internal, "user joe not found")),
			"", false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantIs, Is(tc.givenErr, tc.givenCode))
		})
	}
}

func TestErrCode(t *testing.T) {
	tests := map[string]struct {
		givenError error
		wantErrNo  int
	}{
		"nil": {
			nil,
			-1,
		},
		"not Error": {
			sql.ErrNoRows,
			-1,
		},
		"simple": {
			E(Op("DelUser"), NotFound, 10001, "user joe not found"),
			10001,
		},
		"nested": {
			E(Op("HandleDelUser"), Internal, "userService error",
				E(Op("DelUser"), NotFound, 10001, "user joe not found")),
			10001,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantErrNo, ErrCode(tc.givenError))
		})
	}
}

func TestErrMsg(t *testing.T) {
	tests := map[string]struct {
		givenError error
		wantErrMsg string
	}{
		"nil": {
			nil,
			"",
		},
		"not of type *Error": {
			sql.ErrNoRows,
			defaultMsg,
		},
		"simple": {
			E(Op("DelUser"), NotFound, "user joe not found"),
			"user joe not found",
		},
		"simple code only": {
			E(Op("DelUser"), NotFound, 10001),
			"[10001] An internal error has occurred. Please contact technical support.",
		},
		"simple msg only": {
			E(Op("DelUser"), NotFound, "user joe not found"),
			"user joe not found",
		},
		"simple msg and code": {
			E(Op("DelUser"), NotFound, 10001, "user joe not found"),
			"[10001] user joe not found",
		},
		"wrapped": {
			E(Op("HandleDelUser"), Internal, "userService error",
				E(Op("DelUser"), NotFound, "user joe not found")),
			"userService error",
		},
		"wrapped and first without msg": {
			E(Op("HandleDelUser"), Internal,
				E(Op("DelUser"), NotFound, "user joe not found")),
			"user joe not found",
		},
		"without msg, code and cause": {
			E(Op("HandleDelUser"), Internal,
				E(Op("DelUser"), NotFound)),
			defaultMsg,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantErrMsg, ErrMsg(tc.givenError))
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := map[string]struct {
		givenError error
		wantStr    string
	}{
		"simple": {
			E(Op("DelUser"), NotFound, "user joe not found"),
			"DelUser: user joe not found",
		},
		"simple with Code": {
			E(Op("DelUser"), NotFound, 10001, "user joe not found"),
			"DelUser: [10001] user joe not found",
		},
		"simple with Code, without Msg": {
			E(Op("DelUser"), NotFound, 10001),
			"DelUser: [10001]",
		},
		"wrap external error": {
			E(Op("DelUser"), NotFound, "user joe not found", sql.ErrNoRows),
			"DelUser: user joe not found sql: no rows in result set",
		},
		"wrap *Error": {
			E(
				Op("HandleDelUser"), "", "userService error",
				E(Op("DelUser"), NotFound, "user joe not found")),
			"HandleDelUser: DelUser: user joe not found",
		},
		"wrap *Error and external error": {
			E(
				Op("HandleDelUser"), "", "userService error",
				E(Op("DelUser"), NotFound, "user joe not found", sql.ErrNoRows)),
			"HandleDelUser: DelUser: user joe not found sql: no rows in result set",
		},
		"wrap errors.New": {
			E(Op("DelUser"), NotFound, New("user joe not found")),
			"DelUser: user joe not found",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantStr, tc.givenError.Error())
		})
	}
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		givenText string
		wantError error
	}{
		"name": {"message", &Error{
			Class: "",
			Msg:   "message",
			Op:    "",
			Code:  0,
			Cause: nil,
		}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantError, New(tc.givenText))
		})
	}
}

func TestCombine(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, Combine(nil))
	})

	tests := map[string]struct {
		givenErrs  []error
		wantErrStr string
	}{
		"single": {
			[]error{
				E(Op("DelUser"), NotFound, New("user joe not found")),
			},
			"DelUser: user joe not found",
		},
		"double": {
			[]error{
				E(Op("DelUser"), NotFound, New("user joe not found")),
				E(Op("SendMsg"), Internal, New("internal server error")),
			},
			"[DelUser: user joe not found; SendMsg: internal server error]",
		},
		"multiple": {
			[]error{
				E(Op("A"), NotFound, New("message a")),
				E(Op("B"), Internal, New("message b")),
				E(Op("C"), Conflict, New("message c")),
				E(Op("D"), Other, New("message d")),
			},
			"[A: message a; B: message b; C: message c; D: message d]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := Combine(tc.givenErrs)
			assert.Equal(t, tc.wantErrStr, err.Error())
		})
	}
}

func TestF(t *testing.T) {
	type args struct {
		format string
		a      []interface{}
	}
	tests := []struct {
		name     string
		args     args
		checkErr func(error) bool
	}{
		{
			name: "plain format",
			args: args{
				format: "error for plain format",
				a:      nil,
			},
			checkErr: func(err error) bool {
				return err.Error() == "error for plain format"
			},
		},
		{
			name: "vsd verbs",
			args: args{
				format: "failed to %v, expect %d, got %s",
				a: []interface{}{
					time.Date(2021, 3, 9, 10, 56, 0, 0, time.UTC),
					1056,
					"foo bar baz",
				},
			},
			checkErr: func(err error) bool {
				return err.Error() == "failed to 2021-03-09 10:56:00 +0000 UTC, expect 1056, got foo bar baz"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := F(tt.args.format, tt.args.a...); !tt.checkErr(err) {
				t.Errorf("F() error = %v", err)
			}
		})
	}
}

func TestIsCausedByContextCanceled(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "nil", args: args{nil}, want: false},
		{name: "context canceled", args: args{context.Canceled}, want: true},
		{name: "context deadline exceeded", args: args{context.DeadlineExceeded}, want: false},
		{name: "range error", args: args{strconv.ErrRange}, want: false},
		{name: "wrap", args: args{E(context.Canceled)}, want: true},
		{name: "wrap with class", args: args{E(Invalid, context.Canceled)}, want: true},
		{name: "wrap", args: args{E(Invalid, E("foo bar", E(context.Canceled, "baz")))}, want: true},
		{name: "wrap strings", args: args{E(Invalid, E("foo bar", E("baz")))}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsContextCanceled(tt.args.err); got != tt.want {
				t.Errorf("IsContextCanceled() = %v, want %v", got, tt.want)
			}
		})
	}
}
