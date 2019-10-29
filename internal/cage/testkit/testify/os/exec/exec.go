// Copyright (C) 2019 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package exec

import (
	"bytes"
	"context"
	"errors"
	std_exec "os/exec"

	"github.com/stretchr/testify/mock"

	cage_exec "github.com/codeactual/gomodfuzz/internal/cage/os/exec"
	cage_exec_mocks "github.com/codeactual/gomodfuzz/internal/cage/os/exec/mocks"
)

// NewExecutorMethodAnyCall returns a mock Executor method call that can receive any arguments.
//
// Example context type strings: "*context.timerCtx", "*context.emptyCtx=context.Background"
func NewExecutorMethodAnyCall(executor *cage_exec_mocks.Executor, method string, ctxType interface{}) *mock.Call {
	return executor.On(
		method,
		ctxType,
		mock.AnythingOfType("*exec.Cmd"),
	)
}

// NewSingleCmdPassRes returns a command a pipeline result initialized to describe an execution success.
//
// It may simplify test case boilerplate when configuring a mock to return a passing result
// but the specifics (e.g. standard error content) can just be sane defaults.
//
// It returns the command, in addition to the result, for easier access to the result because
// the former value is the only key of the result map.
func NewSingleCmdPassRes() (*std_exec.Cmd, cage_exec.PipelineResult) {
	expectedCmd := std_exec.Command("(pass cmd, mocked out, not executed)") //nolint:staticcheck
	return expectedCmd, cage_exec.PipelineResult{
		Cmd: map[*std_exec.Cmd]cage_exec.Result{
			expectedCmd: {
				Code:   0,
				Err:    nil,
				Stderr: bytes.NewBufferString("pass stderr"),
				Stdout: bytes.NewBufferString("pass stdout"),
			},
		},
	}
}

// NewSingleCmdFailRes returns a command a pipeline result initialized to describe an execution error.
//
// It may simplify test case boilerplate when configuring a mock to return a failing result
// but the specifics (e.g. standard error content) can just be sane defaults.
//
// It returns the command, in addition to the result, for easier access to the result because
// the former value is the only key of the result map.
func NewSingleCmdFailRes() (*std_exec.Cmd, cage_exec.PipelineResult) {
	expectedCmd := std_exec.Command("(fail cmd, mocked out, not executed)") //nolint:staticcheck
	return expectedCmd, cage_exec.PipelineResult{
		Cmd: map[*std_exec.Cmd]cage_exec.Result{
			expectedCmd: {
				Code:   1,
				Err:    errors.New("mock exec error"),
				Stderr: bytes.NewBufferString("fail stderr"),
				Stdout: bytes.NewBufferString("fail stdout"),
			},
		},
	}
}

// NewExecutorMethodAnyPass returns a mock call to an Executor method which is configured
// to return a success response.
//
// It also returns the command, expected return values, Command/CommandContext's mock call, and
// the input method's mock call.
//
// If the context is non-nil, Executor.CommandContext is used instead of Executor.Command.
func NewExecutorMethodAnyPass(executor *cage_exec_mocks.Executor, method string, ctx context.Context, ctxType interface{}, argsLen int) (_ *std_exec.Cmd, _ cage_exec.PipelineResult, newCmdCall, methodCall *mock.Call) {
	cmd, pipeRes := NewSingleCmdPassRes()
	cmdRes := pipeRes.Cmd[cmd]

	var expectedArgs []interface{}
	for n := 0; n < argsLen; n++ {
		expectedArgs = append(expectedArgs, mock.AnythingOfType("string"))
	}

	// On arguments must align with those used in NewSingleCmdFailRes
	if ctx == nil {
		newCmdCall = executor.On("Command", expectedArgs...)
	} else {
		newCmdCall = executor.On("CommandContext", expectedArgs...)
	}

	methodCall = NewExecutorMethodAnyCall(executor, method, ctxType).
		Return(
			cmdRes.Stdout,
			cmdRes.Stderr,
			pipeRes,
			cmdRes.Err,
		)

	return cmd, pipeRes, newCmdCall, methodCall
}

// NewExecutorMethodAnyFail returns a mock call to an Executor method which is configured
// to return an error response.
//
// It also returns the command, expected return values, Command/CommandContext's mock call, and
// the input method's mock call.
//
// If the context is non-nil, Executor.CommandContext is used instead of Executor.Command.
func NewExecutorMethodAnyFail(executor *cage_exec_mocks.Executor, method string, ctx context.Context, ctxType interface{}, argsLen int) (_ *std_exec.Cmd, _ cage_exec.PipelineResult, newCmdCall, methodCall *mock.Call) {
	cmd, pipeRes := NewSingleCmdFailRes()
	cmdRes := pipeRes.Cmd[cmd]

	var expectedArgs []interface{}
	for n := 0; n < argsLen; n++ {
		expectedArgs = append(expectedArgs, mock.AnythingOfType("string"))
	}

	// On arguments must align with those used in NewSingleCmdFailRes
	if ctx == nil {
		newCmdCall = executor.On("Command", expectedArgs...)
	} else {
		newCmdCall = executor.On("CommandContext", expectedArgs...)
	}

	methodCall = NewExecutorMethodAnyCall(executor, method, ctxType).
		Return(
			cmdRes.Stdout,
			cmdRes.Stderr,
			pipeRes,
			cmdRes.Err,
		)

	return cmd, pipeRes, newCmdCall, methodCall
}
