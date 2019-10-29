// Copyright (C) 2019 The CodeActual Go Environment Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package algo

// Permutator implementations enumerate the permutation axes and all possible values of each.
//
// It assumes a permutation only has one value per axis.
//
// The method naming convention it assumes subject type T may implement this interface itself,
// rather than through a separate Permutator type where the names could stutter less
// (https://blog.golang.org/package-names).
type Permutator interface {
	// PermuteAxes enumerates all the fields whose possible values should yield permutations, e.g. "size" and "color".
	PermuteAxes() []interface{}

	// PermuteSubject returns a zero-valued subject from which permutations are created.
	PermuteSubject() interface{}

	// PermuteNew returns a new permutation with the input axis assigned the input value.
	//
	// The subject is a pointer type, this function should should initialize the new
	// permutation from the subject without modifying the latter.
	//
	// The subject may or may not have other axis values defined. Implementations take the
	// current subject value and only assign one axis/value pair, allowing PermuteSet
	// to be used to compose single subject from multiple PermuteSet calls to set different axes.
	PermuteNew(subject, axis, value interface{}) interface{}

	// PermuteValues returns all possible values of the input axis, e.g. "red" and "green" for axis "colors".
	PermuteValues(axis interface{}) []interface{}

	// PermuteId stores the permutation ID.
	PermuteId(subject interface{}, id int) interface{}
}
