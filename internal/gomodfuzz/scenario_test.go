// Copyright (C) 2019 The gomodfuzz Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gomodfuzz_test

import (
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	tp_algo "github.com/codeactual/gomodfuzz/internal/third_party/stackexchange/algo"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	cage_exec_mocks "github.com/codeactual/gomodfuzz/internal/cage/os/exec/mocks"
	testkit_file "github.com/codeactual/gomodfuzz/internal/cage/testkit/os/file"
	testkit_require "github.com/codeactual/gomodfuzz/internal/cage/testkit/testify/require"
	"github.com/codeactual/gomodfuzz/internal/gomodfuzz"
)

type ScenarioSuite struct {
	suite.Suite
	executor *cage_exec_mocks.Executor
}

func (s *ScenarioSuite) SetupTest() {
	t := s.T()
	s.executor = new(cage_exec_mocks.Executor)
	testkit_file.ResetTestdata(t)
}

func (s *ScenarioSuite) TestPermute() {
	t := s.T()

	expectScenarios := []gomodfuzz.Scenario{
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "auto",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "off",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "-mod=vendor",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.EmptyGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UsableGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   true,
			WD:          gomodfuzz.WdOutsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdInsideGopath,
		},
		{
			GO111MODULE: "on",
			GOFLAGS:     "",
			GOPATH:      gomodfuzz.UnusedGopath,
			IN_MODULE:   false,
			WD:          gomodfuzz.WdOutsideGopath,
		},
	}

	expectRootDir := filepath.Join(testkit_file.DynamicDataDir(), "scenario_applied")
	baseScenario := gomodfuzz.NewScenario(s.executor, expectRootDir)
	permutations := tp_algo.Permute(&baseScenario)

	require.Exactly(t, len(expectScenarios), len(permutations))

	for n, expect := range expectScenarios {
		actual := permutations[n].(gomodfuzz.Scenario)

		require.Exactly(t, expectRootDir, actual.GetRootDir(), expect.String())

		// assert monotonic permutation IDs were assigned
		require.Exactly(t, n, actual.Id())

		// assign the expected ID manually (normally done by Permute)
		expect = expect.PermuteId(expect, n).(gomodfuzz.Scenario)

		// assert permutation values
		require.Exactly(t, expect.GO111MODULE, actual.GO111MODULE, expect.String())
		require.Exactly(t, expect.GOFLAGS, actual.GOFLAGS, expect.String())
		require.Exactly(t, expect.GOPATH, actual.GOPATH, expect.String())
		require.Exactly(t, expect.IN_MODULE, actual.IN_MODULE, expect.String())
		require.Exactly(t, expect.WD, actual.WD, expect.String())

		permuteId := strconv.Itoa(expect.Id())
		expectUsableGopath := filepath.Join(expectRootDir, permuteId, "usable_gopath")

		// assert GOPATH string computed based on the mode

		var expectedGopath string
		switch expect.GOPATH {
		case gomodfuzz.EmptyGopath:
		case gomodfuzz.UsableGopath:
			expectedGopath = expectUsableGopath
		case gomodfuzz.UnusedGopath:
			expectedGopath = filepath.Join(expectRootDir, permuteId, "unused_gopath")
		default:
			t.Fatalf("unexpected GOPATH mode [%d]\n", expect.GOPATH)
		}

		require.Exactly(t, expectedGopath, actual.Gopath())

		// assert working directory computed based on the mode

		var expectedWd string
		switch expect.WD {
		case gomodfuzz.WdInsideGopath:
			expectedWd = filepath.Join(expectUsableGopath, "wd")
		case gomodfuzz.WdOutsideGopath:
			expectedWd = filepath.Join(expectRootDir, permuteId, "wd")
		default:
			t.Fatalf("unexpected WD mode [%d]\n", expect.WD)
		}

		require.Exactly(t, expectedWd, actual.Wd())

		// assert the String() value has the expected permutation value strings

		expectStrPatterns := []string{
			"GO111MODULE=" + expect.GO111MODULE,
			"GOFLAGS=" + expect.GOFLAGS,
		}

		if expect.GOPATH == gomodfuzz.EmptyGopath {
			expectStrPatterns = append(expectStrPatterns, "GOPATH= ")
		} else {
			expectStrPatterns = append(expectStrPatterns, "GOPATH="+filepath.Join(expectRootDir, expect.Gopath()))
		}
		expectStrPatterns = append(
			expectStrPatterns,
			fmt.Sprintf("IN_MODULE=%t", expect.IN_MODULE),
			"Wd="+filepath.Join(expectRootDir, expect.Wd()),
		)
		testkit_require.StringContains(t, actual.String(), expectStrPatterns...)
	}
}

func TestScenarioSuite(t *testing.T) {
	suite.Run(t, new(ScenarioSuite))
}
