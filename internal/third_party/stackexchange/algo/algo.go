package algo

import cage_algo "github.com/codeactual/gomodfuzz/internal/cage/algo"

// Permute returns all permutations defined by the Permutator.
//
// Based on this JS implementation:
//   https://stackoverflow.com/questions/32838388/how-can-i-create-all-combinations-of-this-objects-keys-values-in-javascript/32839413#32839413
//   https://stackoverflow.com/users/236660/dmytro-shevchenko
func Permute(p cage_algo.Permutator) (perms []interface{}) {
	axes := p.PermuteAxes()

	// use a long declaration to make permute identifier available in the function's scope for recursion
	var permute func(axesIdx int, current interface{})
	var nextPermuteId int

	permute = func(axesIdx int, current interface{}) {
		vals := p.PermuteValues(axes[axesIdx])

		for n := 0; n < len(vals); n++ {
			// Create a new permutation for each possible value of the current axis.
			current = p.PermuteNew(current, axes[axesIdx], vals[n])

			if axesIdx+1 < len(axes) {
				permute(axesIdx+1, current) // collect permutations, based on the latest one, from values of the next axis
			} else {
				current = p.PermuteId(current, nextPermuteId)
				nextPermuteId++
				perms = append(perms, current)
			}
		}
	}

	permute(0, p.PermuteSubject()) // start with the first axis and a zero-valued subject

	return perms
}
