## imghash

imghash computes the Perceptual Hash for a given input image.
It supports PNG, GIF and JPEG. The Perceptual Hash is returned
as a 64 bit integer.

This is an implementation of an article on [hackerfactor.com][hf].

Comparing two images can be done by constructing the hash from each image
and counting the number of bit positions that are different. This is a
[Hamming distance][hd]. A distance of zero indicates that it is likely a very
similar picture (or a variation of the same picture). A distance of 5 means a
few things may be different, but they are probably still close enough to be
similar. But a distance of 10 or more? That's probably a very different picture.

[hf]: http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
[hd]: http://en.wikipedia.org/wiki/Hamming_distance

The package supports the hashing modes:

* **Average**: Average Hash is a great algorithm if you are looking for
  something specific. For example, if we have a small thumbnail of an image
  and we know if the big one exists somewhere in our collection.
  Average Hash will find it very quickly. However, if there are modifications
  -- like text was added or a head was spliced into place, then Average Hash
  probably won't do the job. The Average Hash is quick and easy, but it can
  generate false-misses if gamma correction or color histogram is applied to
  the image. This is because the colors move along a non-linear scale --
  changing where the "average" is located and therefore changing which bits
  are above/below the average.


### Usage

    go get github.com/jteeuwen/imghash


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

