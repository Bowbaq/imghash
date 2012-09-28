// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
img-index accepts the path to a given directory.
It constructs Perceptual hashes for all PNG, GIF and JPEG images it
finds and saves them in a database. This database case be used to quickly
find similar images using the `img-find` tool.
*/
package main
