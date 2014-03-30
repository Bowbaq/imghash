// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imghash

import (
	"github.com/nfnt/resize"
	"image"
)

// grayscale turns the image into a grayscale image.
func grayscale(img image.Image) image.Image {
	rect := img.Bounds()
	gray := image.NewGray(rect)

	var x, y int
	for y = rect.Min.Y; y < rect.Max.Y; y++ {
		for x = rect.Min.X; x < rect.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}

	return gray
}

// average converts the sums to averages and returns the result.
func average(sum []uint64, w, h int, n uint64) image.Image {
	ret := image.NewRGBA(image.Rect(0, 0, w, h))
	pix := ret.Pix

	var x, y, idx int
	for y = 0; y < h; y++ {
		for x = 0; x < w; x++ {
			idx = 4 * (y*w + x)
			pix[idx] = uint8(sum[idx] / n)
			pix[idx+1] = uint8(sum[idx+1] / n)
			pix[idx+2] = uint8(sum[idx+2] / n)
			pix[idx+3] = uint8(sum[idx+3] / n)
		}
	}

	return ret
}

func downscale(m image.Image, w, h uint) image.Image {
	return resize.Resize(w, h, m, resize.NearestNeighbor)
}
