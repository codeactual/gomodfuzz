// Copyright (C) 2019 The gomodfuzz Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Command gomodfuzz assists testing of Go program compatability with 1.11+ module support. It runs
// the input program with permutations of GO111MODULE, GOFLAGS, GOPATH, execution from a module directory
// (affecting GOMOD), and working directory (e.g. under GOPATH or not).
//
// Usage:
//
//   gomodfuzz --help
//
// Basic test:
//
//   gomodfuzz -- /path/to/subject --subject_flag0  --subject_flag1 subject_arg0 subject_arg1
//
// Run subject command with a timeout of 10 seconds:
//
//   gomodfuzz --timeout 10 -- /path/to/subject
//
// Display verbose results (passes, full errors, etc.)
//
//   gomodfuzz -v -- /path/to/subject
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tp_algo "github.com/codeactual/gomodfuzz/internal/third_party/stackexchange/algo"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/codeactual/gomodfuzz/internal/cage/cli/handler"
	handler_cobra "github.com/codeactual/gomodfuzz/internal/cage/cli/handler/cobra"
	log_zap "github.com/codeactual/gomodfuzz/internal/cage/cli/handler/mixin/log/zap"
	cage_exec "github.com/codeactual/gomodfuzz/internal/cage/os/exec"
	cage_file "github.com/codeactual/gomodfuzz/internal/cage/os/file"
	cage_file_stage "github.com/codeactual/gomodfuzz/internal/cage/os/file/stage"
	cage_reflect "github.com/codeactual/gomodfuzz/internal/cage/reflect"
	"github.com/codeactual/gomodfuzz/internal/gomodfuzz"
)

const (
	progName = "gomodfuzz"
)

func main() {
	err := handler_cobra.NewHandler(&Handler{
		Session: &handler.DefaultSession{},
	}).Execute()
	if err != nil {
		panic(errors.WithStack(err))
	}
}

// Handler defines the sub-command flags and logic.
type Handler struct {
	handler.Session

	Timeout uint `usage:"Number of seconds to allow the command to run in each scenario"`
	Stdout  bool `usage:"Display standard output from scenarios that fail"`
	Verbose bool `usage:"Display additional status/result information"`

	// example holds command usage examples.
	example []string

	log *log_zap.Mixin

	// stage creates the scenario file trees.
	stage *cage_file_stage.Stage
}

// Init defines the command, its environment variable prefix, etc.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) Init() handler_cobra.Init {
	var err error

	if h.stage, err = cage_file_stage.NewTempDirStage(progName); err != nil {
		h.ExitOnErrShort(err, "init failed", 1)
	}

	h.example = []string{
		progName + " -- /path/to/cmd -flag1 arg1 arg2",
	}

	h.log = &log_zap.Mixin{}

	return handler_cobra.Init{
		Cmd: &cobra.Command{
			Use:     progName,
			Short:   "Asserts Go code compatibility with v1.11+ module scenarios",
			Example: strings.Join(h.example, "\n"),
		},
		EnvPrefix: handler.EnvPrefix(progName),
		Mixins: []handler.Mixin{
			h.log,
		},
	}
}

// BindFlags binds the flags to Handler fields.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) BindFlags(cmd *cobra.Command) []string {
	cmd.Flags().UintVarP(&h.Timeout, "timeout", "t", 30, cage_reflect.GetFieldTag(*h, "Timeout", "usage"))
	cmd.Flags().BoolVarP(&h.Verbose, "verbose", "v", false, cage_reflect.GetFieldTag(*h, "Verbose", "usage"))
	cmd.Flags().BoolVarP(&h.Stdout, "stdout", "o", false, cage_reflect.GetFieldTag(*h, "Stdout", "usage"))
	return []string{}
}

// Run performs the sub-command logic.
//
// It implements cli/handler/cobra.Handler.
func (h *Handler) Run(ctx context.Context, input handler.Input) {
	if len(input.Args) == 0 {
		h.log.Exitf(1, "command not specified (example: %s)", h.example[0])
	}

	// Generate all scenario permutations and run them serially.

	var results []gomodfuzz.Result

	baseScenario := gomodfuzz.NewScenario(cage_exec.CommonExecutor{}, h.stage.Path())

	for _, permutation := range tp_algo.Permute(&baseScenario) {
		s := permutation.(gomodfuzz.Scenario) //nolint:errcheck

		if err := s.BeforeRun(h.stage); err != nil {
			h.log.ExitOnErr(1, errors.Wrapf(err, "failed to run prepare environment for scenario [%s]", s))
		}

		cmdCtx, cmdCancel := context.WithTimeout(ctx, time.Duration(h.Timeout)*time.Second)
		defer cmdCancel()

		r, err := s.Run(cmdCtx, input.Args)
		if err != nil {
			h.log.ExitOnErr(1, errors.Wrapf(err, "failed to run scenario [%s]", s))
		}

		results = append(results, r)
	}

	// Display scenario results.

	hr := func(n int) {
		if n > 0 {
			fmt.Print("\n----\n")
		}
	}

	// passCauses indexes occurrence counts first by variable name (e.g. "GO111MODULE") then by variable value.
	// It supports the pass-cause summary.
	passCauses := map[string]map[string]int{
		"GO111MODULE": {},
		"GOFLAGS":     {},
		"GOPATH":      {},
		"IN_MODULE":   {},
		"WD":          {},
	}

	// failCauses indexes occurrence counts first by variable name (e.g. "GO111MODULE") then by variable value.
	// It supports the failure-cause summary.
	failCauses := map[string]map[string]int{
		"GO111MODULE": {},
		"GOFLAGS":     {},
		"GOPATH":      {},
		"IN_MODULE":   {},
		"WD":          {},
	}

	updateCauses := func(current map[string]map[string]int, s gomodfuzz.Scenario) {
		current["GO111MODULE"][s.GO111MODULE]++
		if s.GOFLAGS == "" {
			current["GOFLAGS"]["<empty>"]++
		} else {
			current["GOFLAGS"][s.GOFLAGS]++
		}
		switch s.GOPATH {
		case gomodfuzz.EmptyGopath:
			current["GOPATH"]["<empty>"]++
		case gomodfuzz.UsableGopath:
			current["GOPATH"]["a file tree that may contain WD"]++
		case gomodfuzz.UnusedGopath:
			current["GOPATH"]["a file that never contains WD"]++
		}
		if s.IN_MODULE {
			current["IN_MODULE"]["inside a module"]++
		} else {
			current["IN_MODULE"]["outside a module"]++
		}
		switch s.WD {
		case gomodfuzz.WdInsideGopath:
			current["WD"]["inside the GOPATH"]++
		case gomodfuzz.WdOutsideGopath:
			current["WD"]["outside the GOPATH"]++
		}
	}

	printCauses := func(title string, causes map[string]map[string]int, samples int) {
		fmt.Println(title)
		for varName, valueCounts := range causes {
			fmt.Println("\t" + varName)
			for val, count := range valueCounts {
				fmt.Printf("\t\t%s: %.2f%%\n", val, (float64(count)/float64(samples))*float64(100))
			}
		}
	}

	var passes int
	for n, r := range results {
		if r.Code == 0 && r.Err == nil {
			if h.Verbose {
				hr(n)
				fmt.Fprintf(h.Out(), "PASS: %s\n", r.Scenario.String())
			}

			updateCauses(passCauses, r.Scenario)

			passes++
		} else {
			hr(n)

			updateCauses(failCauses, r.Scenario)

			fmt.Fprintf(h.Out(), "FAIL (exit code %d): %s\n", r.Code, r.Scenario.String())
			if r.Err != nil && h.Verbose {
				fmt.Fprintf(h.Out(), "\tErr: %+v\n", r.Err)
			}
			fmt.Fprintf(h.Out(), "\tStderr (len=%d): %+v\n", len(r.Stderr), r.Stderr)
			if h.Stdout {
				fmt.Fprintf(h.Out(), "\tStdout (len=%d): %+v\n", len(r.Stdout), r.Stdout)
			}
			if h.Verbose {
				fmt.Fprintf(h.Out(), "\tgo env: %+v\n", strings.TrimSpace(r.GoEnv))
			}
		}
	}

	fmt.Fprintf(h.Out(), "\n- %d/%d scenarios passed\n", passes, len(results))

	if h.Verbose && passes > 0 {
		printCauses("- Occurrences in passes:", passCauses, passes)
	}
	if len(results) != passes {
		printCauses("- Occurrences in failures:", failCauses, len(results)-passes)
	}

	if passes == len(results) {
		h.log.ExitOnErr(1, cage_file.RemoveAllSafer(h.stage.Path()))
	} else {
		fmt.Printf("- Scenario stage will not be deleted so it can be inspected or used for manual tests. Location: %s\n", h.stage.Path())
		os.Exit(2)
	}
}

var _ handler_cobra.Handler = (*Handler)(nil)
