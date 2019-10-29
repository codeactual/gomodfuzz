// Copyright (C) 2019 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package algo_test

import (
	"testing"

	tp_algo "github.com/codeactual/gomodfuzz/internal/third_party/stackexchange/algo"

	"github.com/stretchr/testify/require"

	cage_algo "github.com/codeactual/gomodfuzz/internal/cage/algo"
)

type ThreeAxis struct {
	A, B, C string
}

func (p *ThreeAxis) PermuteNew(subject, axis, value interface{}) interface{} {
	n := subject.(ThreeAxis)
	switch axis.(string) {
	case "A":
		n.A = value.(string)
	case "B":
		n.B = value.(string)
	case "C":
		n.C = value.(string)
	}
	return n
}

func (p *ThreeAxis) PermuteAxes() (axes []interface{}) {
	return append(axes, "A", "B", "C")
}

func (p *ThreeAxis) PermuteSubject() interface{} {
	return ThreeAxis{}
}

func (p *ThreeAxis) PermuteValues(axis interface{}) (values []interface{}) {
	switch axis.(string) {
	case "A":
		values = append(values, "a1", "a2")
	case "B":
		values = append(values, "b1", "b2")
	case "C":
		values = append(values, "c1", "c2")
	}
	return values
}

func (p *ThreeAxis) PermuteId(subject interface{}, id int) interface{} {
	return subject
}

var _ cage_algo.Permutator = (*ThreeAxis)(nil)

type FourAxis struct {
	A, B, C, D string
}

func (p *FourAxis) PermuteNew(subject, axis, value interface{}) interface{} {
	n := subject.(FourAxis)
	switch axis.(string) {
	case "A":
		n.A = value.(string)
	case "B":
		n.B = value.(string)
	case "C":
		n.C = value.(string)
	case "D":
		n.D = value.(string)
	}
	return n
}

func (p *FourAxis) PermuteAxes() (axes []interface{}) {
	return append(axes, "A", "B", "C", "D")
}

func (p *FourAxis) PermuteSubject() interface{} {
	return FourAxis{}
}

func (p *FourAxis) PermuteValues(axis interface{}) (values []interface{}) {
	switch axis.(string) {
	case "A":
		values = append(values, "a1", "a2", "a3")
	case "B":
		values = append(values, "b1", "b2")
	case "C":
		values = append(values, "c1")
	case "D":
		values = append(values, "d1", "d2")
	}
	return values
}

func (p *FourAxis) PermuteId(subject interface{}, id int) interface{} {
	return subject
}

var _ cage_algo.Permutator = (*FourAxis)(nil)

func TestPermute(t *testing.T) {
	var expected []interface{}

	expected = append(expected,
		ThreeAxis{
			A: "a1",
			B: "b1",
			C: "c1",
		},
		ThreeAxis{
			A: "a1",
			B: "b1",
			C: "c2",
		},
		ThreeAxis{
			A: "a1",
			B: "b2",
			C: "c1",
		},
		ThreeAxis{
			A: "a1",
			B: "b2",
			C: "c2",
		},
		ThreeAxis{
			A: "a2",
			B: "b1",
			C: "c1",
		},
		ThreeAxis{
			A: "a2",
			B: "b1",
			C: "c2",
		},
		ThreeAxis{
			A: "a2",
			B: "b2",
			C: "c1",
		},
		ThreeAxis{
			A: "a2",
			B: "b2",
			C: "c2",
		},
	)

	require.Exactly(
		t,
		expected,
		tp_algo.Permute(&ThreeAxis{}),
	)

	expected = nil
	expected = append(expected,
		FourAxis{
			A: "a1",
			B: "b1",
			C: "c1",
			D: "d1",
		},
		FourAxis{
			A: "a1",
			B: "b1",
			C: "c1",
			D: "d2",
		},
		FourAxis{
			A: "a1",
			B: "b2",
			C: "c1",
			D: "d1",
		},
		FourAxis{
			A: "a1",
			B: "b2",
			C: "c1",
			D: "d2",
		},
		FourAxis{
			A: "a2",
			B: "b1",
			C: "c1",
			D: "d1",
		},
		FourAxis{
			A: "a2",
			B: "b1",
			C: "c1",
			D: "d2",
		},
		FourAxis{
			A: "a2",
			B: "b2",
			C: "c1",
			D: "d1",
		},
		FourAxis{
			A: "a2",
			B: "b2",
			C: "c1",
			D: "d2",
		},
		FourAxis{
			A: "a3",
			B: "b1",
			C: "c1",
			D: "d1",
		},
		FourAxis{
			A: "a3",
			B: "b1",
			C: "c1",
			D: "d2",
		},
		FourAxis{
			A: "a3",
			B: "b2",
			C: "c1",
			D: "d1",
		},
		FourAxis{
			A: "a3",
			B: "b2",
			C: "c1",
			D: "d2",
		},
	)

	require.Exactly(
		t,
		expected,
		tp_algo.Permute(&FourAxis{}),
	)
}
