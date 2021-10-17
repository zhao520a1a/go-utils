// mainly inspired by Ben Johnson, read the following blog for details:
// https://middlemost.com/failure-is-your-domain/
package errors

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/zhao520a1a/go-utils/collection"
)

// TODO:
// auto generate op by runtime info

// Error defines a standard application error.
type Error struct {
	// For application/machine
	Class Class
	// For users & operators, see methods ErrMsg (users) and Error (operators)
	Msg string
	// For operators
	Op    Op    // logical operation
	Code  int   // error code, which identifies an user-defined error
	Cause error // error from lower level
}

type Op string
type Class string

// E constructs an *Error with variable number of args, which corresponds to
// the Class, Msg, Op, Code and Cause fields, all args are optional.
func E(args ...interface{}) error {
	args = collection.InterfaceFilter(args, func(t interface{}) bool {
		return t != nil
	})

	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}

	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Class:
			e.Class = arg
		case string:
			e.Msg = arg
		case Op:
			e.Op = arg
		case int:
			e.Code = arg
		case int32:
			e.Code = int(arg)
		case *Error:
			cp := *arg
			e.Cause = &cp
		case error:
			e.Cause = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return fmt.Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	// deduplication
	if cause, ok := e.Cause.(*Error); ok && cause.Op == e.Op {
		return cause
	}

	// TODO: callstack
	return e
}

// F formats according to a format specifier and returns the string
// as a value that satisfies error.
func F(format string, a ...interface{}) error {
	return E(fmt.Sprintf(format, a...))
}

// Is determines whether the given error has the code along the chain, only the
// first non-empty Class encountered is considered.
func Is(err error, class Class) bool {
	if err == nil {
		return class == ""
	} else if e, ok := err.(*Error); ok {
		if e.Class == "" {
			return Is(e.Cause, class)
		}
		return e.Class == class
	}
	return false
}

func IsContextCanceled(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Cause == context.Canceled || IsContextCanceled(e.Cause)
	}
	return err == context.Canceled
}

func (e *Error) Error() string {
	b := bytes.NewBuffer(nil)
	if e.Op != "" {
		_, _ = fmt.Fprintf(b, "%s: ", e.Op)
	}
	// print operation info of the tail error
	if e.Cause == nil {
		e.writeOpInfo(b)
		return b.String()
	}

	// if the inner error is of type *Error, only print the Op,
	// otherwise, print operation info and the inner error
	if _, isError := e.Cause.(*Error); isError {
		b.WriteString(e.Cause.Error())
	} else {
		e.writeOpInfo(b)
		b.WriteByte(' ')
		b.WriteString(e.Cause.Error())
	}

	return b.String()
}

func (e *Error) writeOpInfo(b *bytes.Buffer) {
	if e.Code != 0 && len(e.Msg) > 0 {
		_, _ = fmt.Fprintf(b, "[%d] %s", e.Code, e.Msg)
	} else if e.Code != 0 {
		_, _ = fmt.Fprintf(b, "[%d]", e.Code)
	} else if len(e.Msg) > 0 {
		_, _ = fmt.Fprintf(b, "%s", e.Msg)
	}
}

// ErrCode returns the ErrCode of the first error along the chain,
// otherwise returns -1.
func ErrCode(err error) int {
	if err == nil {
		return -1
	} else if e, ok := err.(*Error); ok && e.Code != 0 {
		return e.Code
	} else if ok && e.Cause != nil {
		return ErrCode(e.Cause)
	}
	return -1
}

const defaultMsg = "An internal error has occurred. Please contact technical support."

// ErrMsg returns the first human-readable message along the chain,
// otherwise returns a default generic message.
func ErrMsg(err error) string {
	code := firstCode(err)
	msg := firstMsg(err)

	if msg != "" && code != 0 {
		return fmt.Sprintf("[%d] %s", code, msg)
	}

	return msg
}

func firstCode(err error) int {
	if err == nil {
		return 0
	}

	e, ok := err.(*Error)
	if !ok {
		return 0
	}

	if e.Code != 0 {
		return e.Code
	}
	return firstCode(e.Cause)
}

func firstMsg(err error) string {
	// only return empty string when err == nil
	if err == nil {
		return ""
	}

	e, ok := err.(*Error)
	if !ok {
		return defaultMsg
	}

	if e.Msg != "" {
		return e.Msg
	}

	if e.Cause == nil {
		return defaultMsg
	}

	return firstMsg(e.Cause)
}

func New(text string) error {
	return E(text)
}

func Combine(errs []error) error {
	numErrors := len(errs)

	if numErrors == 0 {
		return nil
	}

	if numErrors == 1 {
		return errs[0]
	}

	errStrs := make([]string, 0, numErrors)
	for _, err := range errs {
		errStrs = append(errStrs, err.Error())
	}

	combinedStr := fmt.Sprintf("[%s]", strings.Join(errStrs, "; "))
	return New(combinedStr)
}
