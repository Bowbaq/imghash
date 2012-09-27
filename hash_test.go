// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package imghash

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"
)

// Maximum Hamming-distance at which we consider images to be un-equal.
const MaxDistance = 3

func TestAverage(t *testing.T) {
	var avg Average
	a := getHash(t, avg, "testdata/gopher_large.png")
	b := getHash(t, avg, "testdata/gopher_small.png")

	dist := Distance(a, b)
	if dist >= MaxDistance {
		t.Fatalf("Hash mismatch: %d", dist)
	}
}

func getHash(t *testing.T, h Hasher, file string) uint64 {
	fd, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}

	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		t.Fatal(err)
	}

	return h.Compute(img)
}
