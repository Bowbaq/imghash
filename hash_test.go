// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package imghash

import (
	"image"
	_ "image/png"
	"os"
	"testing"
)

// Maximum Hamming-distance at which we consider images to be un-equal.
const MaxDistance = 3

func TestAverage(t *testing.T) {
	a := getHash(t, Average, "testdata/gopher_large.png")
	b := getHash(t, Average, "testdata/gopher_small.png")

	dist := Distance(a, b)
	if dist >= MaxDistance {
		t.Fatalf("Hash mismatch: %d", dist)
	}
}

func getHash(t *testing.T, hf HashFunc, file string) uint64 {
	fd, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}

	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		t.Fatal(err)
	}

	return hf(img)
}
