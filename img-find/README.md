## img-find

img-find accepts the path to a single image.
It constructs a Perceptual Hash for it and checks a hash database
for any matches. It returns the filenames associated with these
matches, sorted by relevance.


## Matches

A match is determined by the [Hamming Distance][hd] between two hashes.
The threashold for this value can be specified through a commandline
parameter. The smaller the value, the more restrictive the matches.

[hd]: http://en.wikipedia.org/wiki/Hamming_distance


## Database

The hash database it uses is a flat file. It is stored in the location
specified either in the `-db` commandline parameter, or in the `IMGHASH_DB`
environment variable.


### Usage

    go get github.com/jteeuwen/imghash/img-find


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

