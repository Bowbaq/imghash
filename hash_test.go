// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package imghash

import (
	"image"
	"image/png"
	"os"
	"testing"
)

// Maximum Hamming-distance at which we consider images to be equal.
const MaxDistance = 3

func TestResize(t *testing.T) {
	img, err := loadImg("testdata/gopher_large.png")
	if err != nil {
		t.Fatal(err)
	}

	img = resize(img, 32, 32)

	err = saveImg(img, "testdata/gopher_32x32.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAverage(t *testing.T) {
	a := getHash(t, Average, "testdata/gopher_large.png")
	b := getHash(t, Average, "testdata/gopher_small.png")

	dist := Distance(a, b)
	if dist > MaxDistance {
		t.Fatalf("Hash mismatch: 0x%x 0x%x %d\n", a, b, dist)
	}
}

func getHash(t *testing.T, hf HashFunc, file string) uint64 {
	img, err := loadImg(file)

	if err != nil {
		t.Fatal(err)
	}

	return hf(img)
}

func loadImg(file string) (image.Image, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func saveImg(img image.Image, file string) error {
	fd, err := os.Create(file)
	if err != nil {
		return err
	}

	defer fd.Close()

	err = png.Encode(fd, img)
	if err != nil {
		return err
	}

	return nil
}
