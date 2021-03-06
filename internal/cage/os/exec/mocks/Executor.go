// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import bytes "bytes"
import context "context"
import exec "os/exec"
import io "io"
import mock "github.com/stretchr/testify/mock"
import osexec "github.com/codeactual/gomodfuzz/internal/cage/os/exec"

// Executor is an autogenerated mock type for the Executor type
type Executor struct {
	mock.Mock
}

// Buffered provides a mock function with given fields: ctx, cmds
func (_m *Executor) Buffered(ctx context.Context, cmds ...*exec.Cmd) (*bytes.Buffer, *bytes.Buffer, osexec.PipelineResult, error) {
	_va := make([]interface{}, len(cmds))
	for _i := range cmds {
		_va[_i] = cmds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *bytes.Buffer
	if rf, ok := ret.Get(0).(func(context.Context, ...*exec.Cmd) *bytes.Buffer); ok {
		r0 = rf(ctx, cmds...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bytes.Buffer)
		}
	}

	var r1 *bytes.Buffer
	if rf, ok := ret.Get(1).(func(context.Context, ...*exec.Cmd) *bytes.Buffer); ok {
		r1 = rf(ctx, cmds...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*bytes.Buffer)
		}
	}

	var r2 osexec.PipelineResult
	if rf, ok := ret.Get(2).(func(context.Context, ...*exec.Cmd) osexec.PipelineResult); ok {
		r2 = rf(ctx, cmds...)
	} else {
		r2 = ret.Get(2).(osexec.PipelineResult)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, ...*exec.Cmd) error); ok {
		r3 = rf(ctx, cmds...)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Command provides a mock function with given fields: name, arg
func (_m *Executor) Command(name string, arg ...string) *exec.Cmd {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *exec.Cmd
	if rf, ok := ret.Get(0).(func(string, ...string) *exec.Cmd); ok {
		r0 = rf(name, arg...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*exec.Cmd)
		}
	}

	return r0
}

// CommandContext provides a mock function with given fields: ctx, name, arg
func (_m *Executor) CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	_va := make([]interface{}, len(arg))
	for _i := range arg {
		_va[_i] = arg[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *exec.Cmd
	if rf, ok := ret.Get(0).(func(context.Context, string, ...string) *exec.Cmd); ok {
		r0 = rf(ctx, name, arg...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*exec.Cmd)
		}
	}

	return r0
}

// Pty provides a mock function with given fields: cmd
func (_m *Executor) Pty(cmd *exec.Cmd) error {
	ret := _m.Called(cmd)

	var r0 error
	if rf, ok := ret.Get(0).(func(*exec.Cmd) error); ok {
		r0 = rf(cmd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Standard provides a mock function with given fields: ctx, stdout, stderr, stdin, cmds
func (_m *Executor) Standard(ctx context.Context, stdout io.Writer, stderr io.Writer, stdin io.Reader, cmds ...*exec.Cmd) (osexec.PipelineResult, error) {
	_va := make([]interface{}, len(cmds))
	for _i := range cmds {
		_va[_i] = cmds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, stdout, stderr, stdin)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 osexec.PipelineResult
	if rf, ok := ret.Get(0).(func(context.Context, io.Writer, io.Writer, io.Reader, ...*exec.Cmd) osexec.PipelineResult); ok {
		r0 = rf(ctx, stdout, stderr, stdin, cmds...)
	} else {
		r0 = ret.Get(0).(osexec.PipelineResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, io.Writer, io.Writer, io.Reader, ...*exec.Cmd) error); ok {
		r1 = rf(ctx, stdout, stderr, stdin, cmds...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
