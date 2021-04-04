package errors

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithCaller sets the position at which the error was formed. If this is
// false, errors will be wrapped with no location information.
var WithCaller = true

// WithCallerVerbose just adds file:line components alongside the caller information.
var WithCallerVerbose = false

// WithStack provides a stack trace at the point the error was yielded, either
// via Error() method or direct function call. This may be set at any time to
// get a stack trace when calling Error() on a WithLocation error.
var WithStack = false

// StackBufferSize can impact performance and should be adjusted if you have
// small or larger stack traces than the default size, this is the bytes
// buffer that will be created for calculating stack traces.
var StackBufferSize = 8192

// WithLocation is a kind of error that stores the location at which it was
// created at. Dereference `Err` to get at the actual error.
type WithLocation struct {
	File string
	Line int
	Loc  uintptr

	Err error
}

func (wl WithLocation) Error() string {
	s := wl.Err.Error()

	if wl.Loc != 0 {
		verbose := ""
		if wl.File != "" && wl.Line != 0 && WithCallerVerbose {
			verbose = fmt.Sprintf("[%v:%v]", wl.File, wl.Line)
		}

		s = fmt.Sprintf("[%v]%v: %v", runtime.FuncForPC(wl.Loc).Name(), verbose, s)
	}

	sbuf := bytes.NewBufferString(s)
	if WithStack {
		buf := make([]byte, StackBufferSize)
		sbuf.WriteString("\nSTACK TRACE:\n-----------\n")
		sbuf.Write(buf[:runtime.Stack(buf, false)])
	}

	return sbuf.String()
}

// New creates a new error with caller identification.
func New(s string) error {
	return mkError(errors.New(s))
}

// Errorf returns a formatted error with caller identification.
func Errorf(f string, args ...interface{}) error {
	return mkError(fmt.Errorf(f, args...))
}

// WithError combines two errors to be compatible with errors.As() and
// similar functions in the stdlib errors package.
func WithError(err error, f string, args ...interface{}) error {
	return mkError(fmt.Errorf("%v: %w", fmt.Sprintf(f, args...), err))
}

// GRPC formats a GRPC error in our caller ID pattern.
func GRPC(code codes.Code, f string, args ...interface{}) error {
	return status.Errorf(code, f, args...)
}

func transformWithCaller(depth int, e error) error {
	if WithCaller {
		pc, file, line, ok := runtime.Caller(depth)
		if ok {
			return WithLocation{
				File: file,
				Line: line,
				Loc:  pc,
				Err:  e,
			}
		}
	}

	return WithLocation{Err: e}
}

func mkError(e error) error {
	return transformWithCaller(3, e)
}
