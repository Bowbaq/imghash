// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package imghash

import "image"

// A Hasher computes a Perceptual Hash for a given image.
type Hasher interface {
	Compute(image.Image) uint64
}

// Distance calculates the Hamming Distance between the two input hashes.
func Distance(a, b uint64) int {
	var dist int

	for val := a ^ b; val != 0; val &= val - 1 {
		dist++
	}

	return dist
}
