// Copyright (C) 2019 The gomodfuzz Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gomodfuzz

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	cage_algo "github.com/codeactual/gomodfuzz/internal/cage/algo"
	cage_exec "github.com/codeactual/gomodfuzz/internal/cage/os/exec"
	cage_file_stage "github.com/codeactual/gomodfuzz/internal/cage/os/file/stage"
)

const (
	newDirPerm  = 0755
	newFilePerm = 0666 // Match os.Create for use with ioutil.WriteFile

	// Scenario.GOPATH selection modes
	EmptyGopath = iota
	UsableGopath
	UnusedGopath

	// Scenario.WD selection modes
	WdInsideGopath = iota
	WdOutsideGopath
)

// Scenario defines how a command should be executed in a scenario.
type Scenario struct {
	// GO111MODULE is the environment variable value applied to the scenario.
	//
	// It is assigned a value by a permutation generator. The generator assigns one of
	// three values: "auto", "off", "on".
	GO111MODULE string

	// GOFLAGS is the environment variable value applied to the scenario.
	//
	// It is assigned a value by a permutation generator. The generator assigns one of
	// two values: empty string or "-mod=vendor".
	GOFLAGS string

	// GOPATH is a mode of selecting environment variable value applied to the scenario.
	//
	// It is assigned a value by a permutation generator. The generator assigns one of three modes
	// which select these path types: a path which may contain the working directory as a descendant,
	// a path which never contains the working directory, and an empty string.
	GOPATH int

	// IN_MODULE is true if the command should in a working directory with a go.mod.
	//
	// It is assigned a value by a permutation generator.
	IN_MODULE bool

	// Wd is a mode of selecting the working directory in which the command runs.
	//
	// It is assigned a value by a permutation generator. The generator assigns one of
	// two modes which select these path types: "<Scenario.rootDir>/<scenario dir>/wd" or
	// "<Scenario.rootDir>/<scenario dir>/gopath/wd".
	WD int

	// executor implementations run os/exec commands, allowing tests to mock their execution.
	executor cage_exec.Executor

	// rootDir is the top of the file tree dedicated to the testing of this particular Scenario/permutation.
	// The command's working directory, and GOPATH (if enabled), will be under this root directory.
	rootDir string

	// permuteId is assigned the current value of nextPermuteId when the Wd value is generated in PermuteValues.
	//
	// It is used to generate unique Wd values.
	permuteId int
}

// NewScenario returns an initialized value.
func NewScenario(executor cage_exec.Executor, rootDir string) Scenario {
	return Scenario{executor: executor, rootDir: rootDir}
}

// BeforeRun sets up the environment in preparation for Run.
func (s Scenario) BeforeRun(stage *cage_file_stage.Stage) error {
	// Create the go.mod file to simulate running the input command from a module's directory.
	if s.IN_MODULE {
		modFilePath := filepath.Join(s.Wd(), "go.mod")
		relPath, pathErr := filepath.Rel(s.rootDir, modFilePath)
		if pathErr != nil {
			return errors.Wrapf(pathErr,
				"failed to get relative path from [%s] to [%s]", s.rootDir, modFilePath,
			)
		}

		gomodFile, createErr := stage.CreateFileAll(relPath, newFilePerm, newDirPerm)
		if createErr != nil {
			return errors.Wrapf(createErr,
				"failed to create go.mod in scenario [%s] working directory [%s]", s.String(), s.Wd(),
			)
		}

		if _, writeErr := gomodFile.WriteString("module wd\n"); writeErr != nil {
			return errors.Wrapf(writeErr,
				"failed to update go.mod in scenario [%s] working directory [%s]", s.String(), s.Wd(),
			)
		}
	} else {
		relPath, pathErr := filepath.Rel(stage.Path(), s.Wd())
		if pathErr != nil {
			return errors.Wrapf(pathErr,
				"failed to get relative path from [%s] to [%s]", stage.Path(), s.Wd(),
			)
		}

		if mkdirErr := stage.MkdirAll(relPath, newDirPerm); mkdirErr != nil {
			return errors.Wrapf(mkdirErr, "failed to create go.mod in scenario [%s] working directory [%s]", s.String(), s.Wd())
		}
	}

	return nil
}

// Run applies the permutation-defined fields, runs the command, and returns the result.
func (s Scenario) Run(ctx context.Context, args []string) (res Result, err error) {
	// collectCmdRes runs the command with permutation-defined config applied.
	collectCmdRes := func(cmd *exec.Cmd) (stdout, stderr string, pipeRes cage_exec.PipelineResult, err error) {
		cmd.Env = append(os.Environ(), []string{
			"GO111MODULE=" + s.GO111MODULE,
			"GOFLAGS=" + s.GOFLAGS,
			"GOPATH=" + s.Gopath(),
		}...)
		cmd.Dir = s.Wd()

		stdoutBuf, stderrBuf, pipeRes, cmdErr := s.executor.Buffered(ctx, cmd)

		ctxErr := ctx.Err()
		if ctxErr != nil {
			cmdErr = ctxErr
		}

		return stdoutBuf.String(), stderrBuf.String(), pipeRes, cmdErr
	}

	name := s.String()
	res = NewResult(s)

	// Collect `go env` output to display if the scenario fails.

	goEnvCmd := s.executor.Command("go", "env")
	if goEnvStdout, _, _, err := collectCmdRes(goEnvCmd); err == nil {
		res.GoEnv = goEnvStdout
	} else {
		return Result{}, errors.Wrapf(err, "failed to run 'go env' for scenario [%s]", name)
	}

	// Run the input command.

	var subjectCmd *exec.Cmd
	if len(args) == 1 {
		subjectCmd = s.executor.Command(args[0]) // #nosec
	} else {
		subjectCmd = s.executor.Command(args[0], args[1:]...) // #nosec
	}
	subjectStdout, subjectStderr, pipeRes, err := collectCmdRes(subjectCmd)

	res.Code = pipeRes.Cmd[subjectCmd].Code
	res.Err = err
	res.Stderr = strings.TrimSpace(subjectStderr)
	res.Stdout = strings.TrimSpace(subjectStdout)
	res.Scenario = s

	return res, nil
}

// GetRootDir returns the top of the scenario's file tree.
func (s Scenario) GetRootDir() string {
	return s.rootDir
}

// String returns a scenario identifier for display.
func (s Scenario) String() string {
	return fmt.Sprintf(
		"GO111MODULE=%s "+
			"GOFLAGS=%s "+
			"GOPATH=%s "+
			"IN_MODULE=%t "+
			"Wd=%s",
		s.GO111MODULE,
		s.GOFLAGS,
		s.Gopath(),
		s.IN_MODULE,
		s.Wd(),
	)
}

// PermuteAxes enumerates all the fields whose possible values should yield permutations, e.g. "size" and "color".
//
// It implements Permutator.
func (s *Scenario) PermuteAxes() (axes []interface{}) {
	axes = append(axes, "GO111MODULE", "GOFLAGS", "GOPATH", "IN_MODULE", "WD")
	return axes
}

// PermuteSubject returns a zero-valued subject from which permutations are created.
//
// It implements Permutator.
func (s *Scenario) PermuteSubject() interface{} {
	scenario := Scenario{
		executor: s.executor,
		rootDir:  s.rootDir,
	}
	return scenario
}

// PermuteNew returns a new permutation with the input axis assigned the input value.
//
// It implements Permutator.
func (s *Scenario) PermuteNew(subject, axis, value interface{}) interface{} {
	n := subject.(Scenario) //nolint:errcheck
	switch axis.(string) {
	case "GO111MODULE":
		n.GO111MODULE = value.(string) //nolint:errcheck
	case "GOFLAGS":
		n.GOFLAGS = value.(string) //nolint:errcheck
	case "GOPATH":
		n.GOPATH = value.(int) //nolint:errcheck
	case "IN_MODULE":
		n.IN_MODULE = value.(bool) //nolint:errcheck
	case "WD":
		n.WD = value.(int) //nolint:errcheck
	}
	return n
}

func (s Scenario) Id() int {
	return s.permuteId
}

func (s Scenario) UsableGopath() string {
	return filepath.Join(s.rootDir, strconv.Itoa(s.permuteId), "usable_gopath")
}

func (s Scenario) Gopath() string {
	switch s.GOPATH {
	case UsableGopath:
		// UsableGopath is an path to a file tree which may contain the working directory value (Wd)
		// as a descendant. This enables permutations where the command "runs from in the GOPATH."
		return s.UsableGopath()
	case UnusedGopath:
		// UnusedGopath complements UsableGopath by enabling permutations where the working directory value (Wd)
		// is never a descendant. This enables permutations where the environment variable is non-empty/valid but the
		// command "runs from outside the GOPATH".
		return filepath.Join(s.rootDir, strconv.Itoa(s.permuteId), "unused_gopath")
	case EmptyGopath:
		return ""
	default:
		panic(errors.Errorf("scenario generator used an invalid GOPATH mode [%d]", s.GOPATH))
	}
}

func (s Scenario) Wd() string {
	switch s.WD {
	case WdOutsideGopath:
		// WdOutsideGopath aligns with UsableGopath/UnusedGopath, by being a descendent of neither, to enable
		// permutations where the command "runs from outside the GOPATH."
		return filepath.Join(s.rootDir, strconv.Itoa(s.permuteId), "wd")
	case WdInsideGopath:
		// WdInsideGopath aligns with UsableGopath to enable permutations where the command "runs from in the GOPATH."
		return filepath.Join(s.UsableGopath(), "wd")
	default:
		panic(errors.Errorf("scenario generator used an invalid Wd mode [%d]", s.WD))
	}
}

// PermuteValues returns all possible values of the input axis, e.g. "red" and "green" for axis "colors".
//
// It implements Permutator.
func (s *Scenario) PermuteValues(axis interface{}) (values []interface{}) {
	switch axis.(string) {
	case "GO111MODULE":
		values = append(values, "auto", "off", "on")
	case "GOFLAGS":
		values = append(values, "-mod=vendor", "")
	case "GOPATH":
		values = append(values, EmptyGopath, UsableGopath, UnusedGopath)
	case "IN_MODULE":
		values = append(values, true, false)
	case "WD":
		values = append(values, WdInsideGopath, WdOutsideGopath)
	}
	return values
}

// PermuteId stores the permutation ID.
//
// It implements Permutator.
func (s Scenario) PermuteId(subject interface{}, id int) interface{} {
	scenario := subject.(Scenario) //nolint:errcheck
	scenario.permuteId = id
	return scenario
}

var _ cage_algo.Permutator = (*Scenario)(nil)
