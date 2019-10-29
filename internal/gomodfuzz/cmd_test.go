// Copyright (C) 2019 The gomodfuzz Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gomodfuzz_test

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	tp_algo "github.com/codeactual/gomodfuzz/internal/third_party/stackexchange/algo"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	cage_exec "github.com/codeactual/gomodfuzz/internal/cage/os/exec"
	cage_exec_mocks "github.com/codeactual/gomodfuzz/internal/cage/os/exec/mocks"
	cage_file "github.com/codeactual/gomodfuzz/internal/cage/os/file"
	cage_file_stage "github.com/codeactual/gomodfuzz/internal/cage/os/file/stage"
	testkit_file "github.com/codeactual/gomodfuzz/internal/cage/testkit/os/file"
	testkit_filepath "github.com/codeactual/gomodfuzz/internal/cage/testkit/path/filepath"
	cage_testify_exec "github.com/codeactual/gomodfuzz/internal/cage/testkit/testify/os/exec"
	"github.com/codeactual/gomodfuzz/internal/gomodfuzz"
)

type CmdSuite struct {
	suite.Suite
	executor *cage_exec_mocks.Executor
}

func (s *CmdSuite) SetupTest() {
	t := s.T()
	s.executor = new(cage_exec_mocks.Executor)
	testkit_file.ResetTestdata(t)
}

func (s *CmdSuite) TestScenarioApplied() {
	t := s.T()

	ctx := context.Background()
	ctxType := mock.AnythingOfType("*context.emptyCtx") // context.Background()
	expectRootDir := filepath.Join(testkit_file.DynamicDataDir(), "scenario_applied")
	stage := cage_file_stage.NewStage(expectRootDir)
	stagePath := testkit_filepath.Abs(t, stage.Path())

	baseScenario := gomodfuzz.NewScenario(s.executor, expectRootDir)
	permutations := tp_algo.Permute(&baseScenario)

	subjectCmdArgs := []string{"arg0", "arg1", "arg2"}
	expectedGoEnvStdout := "fake 'go env' output"

	// requireGoEnvCmdExactly asserts the executed os/exec.Cmd had the expected environment.
	requireGoEnvCmdExactly := func(cmd *exec.Cmd, scenario gomodfuzz.Scenario) func(mock.Arguments) {
		return func(mockArgs mock.Arguments) {
			mockCmd := mockArgs.Get(1)
			require.Contains(t, mockCmd.(*exec.Cmd).Env, "GO111MODULE="+scenario.GO111MODULE)
			require.Contains(t, mockCmd.(*exec.Cmd).Env, "GOFLAGS="+scenario.GOFLAGS)
			require.Contains(t, mockCmd.(*exec.Cmd).Env, "GOPATH="+scenario.Gopath())
			require.Contains(t, mockCmd.(*exec.Cmd).Dir, scenario.Wd())
			require.Exactly(t, cmd.Args, mockCmd.(*exec.Cmd).Args)
		}
	}

	// requireSubjectCmdExactly asserts the executed os/exec.Cmd had the expected environment.
	requireSubjectCmdExactly := func(cmd *exec.Cmd, scenario gomodfuzz.Scenario) func(mock.Arguments) {
		return func(mockArgs mock.Arguments) {
			mockCmd := mockArgs.Get(1)
			require.Contains(t, mockCmd.(*exec.Cmd).Env, "GO111MODULE="+scenario.GO111MODULE)
			require.Contains(t, mockCmd.(*exec.Cmd).Env, "GOFLAGS="+scenario.GOFLAGS)
			require.Contains(t, mockCmd.(*exec.Cmd).Env, "GOPATH="+scenario.Gopath())
			require.Contains(t, mockCmd.(*exec.Cmd).Dir, scenario.Wd())
			require.Exactly(t, cmd.Args, mockCmd.(*exec.Cmd).Args)
		}
	}

	mockCallIdx := 0 // track so we can customize mocks created by NewExecutorMethodAny*

	for _, p := range permutations {
		// setup for both sub-scenarios

		scenario := p.(gomodfuzz.Scenario)
		sid := scenario.String()
		goModPath := testkit_filepath.Abs(t, filepath.Join(scenario.Wd(), "go.mod"))

		// expect go.mod created if IN_MODULE is true
		if err := scenario.BeforeRun(stage); err != nil {
			require.NoError(t, err, sid)
		}

		exists, _, err := cage_file.Exists(goModPath)
		require.NoError(t, err, sid)
		require.Exactly(t, scenario.IN_MODULE, exists, sid)

		for subScenario := 0; subScenario < 2; subScenario++ { // run the scenario twice (pass + fail)
			subScenarioExpectPass := subScenario == 0

			// expect two Executor calls for running 'go env': create an os/exec.Cmd and execute it

			goEnvCmd, goEnvRes, newCmdCall, goEnvCall := cage_testify_exec.NewExecutorMethodAnyPass(s.executor, "Buffered", nil, ctxType, 2)

			newCmdCall.Return(goEnvCmd).Once()
			mockCallIdx++

			goEnvCallReturnStdout := s.executor.ExpectedCalls[mockCallIdx].ReturnArguments[0].(*bytes.Buffer)
			goEnvCallReturnStdout.Reset()
			goEnvCallReturnStdout.WriteString(expectedGoEnvStdout) // replace the default from NewExecutorMethodAnyPass
			goEnvCall.Run(requireGoEnvCmdExactly(goEnvCmd, scenario)).Once()
			mockCallIdx++

			// expect two Executor calls for running the subject command: create an os/exec.Cmd and execute it

			var subjectCmd *exec.Cmd
			var subjectRes cage_exec.PipelineResult
			var subjectExecCall *mock.Call

			if subScenarioExpectPass {
				subjectCmd, subjectRes, newCmdCall, subjectExecCall = cage_testify_exec.NewExecutorMethodAnyPass(s.executor, "Buffered", nil, ctxType, len(subjectCmdArgs))
			} else {
				subjectCmd, subjectRes, newCmdCall, subjectExecCall = cage_testify_exec.NewExecutorMethodAnyFail(s.executor, "Buffered", nil, ctxType, len(subjectCmdArgs))
			}

			newCmdCall.Return(subjectCmd).Once()
			mockCallIdx++

			subjectExecCall.Run(requireSubjectCmdExactly(subjectCmd, scenario)).Once()
			mockCallIdx++

			runRes, err := scenario.Run(ctx, subjectCmdArgs) // runs both above commands
			require.NoError(t, err, sid)

			require.Exactly(t, expectedGoEnvStdout, goEnvRes.Cmd[goEnvCmd].Stdout.String(), sid)
			require.Exactly(t, expectedGoEnvStdout, runRes.GoEnv, sid)

			require.Exactly(t, subjectRes.Cmd[subjectCmd].Code, runRes.Code, sid)
			require.Exactly(t, subjectRes.Cmd[subjectCmd].Err, runRes.Err, sid)
			require.Exactly(t, subjectRes.Cmd[subjectCmd].Stderr.String(), runRes.Stderr, sid)
			require.Exactly(t, subjectRes.Cmd[subjectCmd].Stdout.String(), runRes.Stdout, sid)

			s.executor.AssertExpectations(t)

			require.NoError(t, cage_file.RemoveAllSafer(testkit_filepath.Abs(t, scenario.Wd())))
		}
	}

	require.NoError(t, cage_file.RemoveAllSafer(stagePath))
}

func TestCmdSuite(t *testing.T) {
	suite.Run(t, new(CmdSuite))
}
