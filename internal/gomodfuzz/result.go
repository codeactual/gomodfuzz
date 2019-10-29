// Copyright (C) 2019 The gomodfuzz Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package gomodfuzz

// Result is the outcome of one execution of the input command in one Scenario.
type Result struct {
	// Scenario is a copy of the executed scenario spec.
	Scenario Scenario

	// Err is non-nil if the scenario's command fails to start, exits with a non-zero code,
	// or its context times out.
	//
	// It is also non-nil if `go env` fails before running the scenario's command.
	Err error

	// Code is from the scenario's command.
	//
	// It is -1 if the command does not get an opporunity to run, e.g. if `go env` fails
	// for some reason.
	Code int

	// GoEnv is the output of `go env` prior to running the scenario.
	GoEnv string

	// Stderr is from the scenario's command.
	Stderr string

	// Stdout is from the scenario's command.
	Stdout string
}

// NewResult returns an initialized Result.
func NewResult(s Scenario) Result {
	return Result{
		Code:     -1,
		Scenario: s,
	}
}
