// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/imghash"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

var (
	file   string
	dist   = flag.Uint64("dist", 5, "")
	dbfile = flag.String("db", "", "")
)

func main() {
	var db imghash.Database

	parseArgs()

	// Compute averahe hash for the input image.
	hash, err := getHash(imghash.Average, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Load the database file.
	err = db.Load(*dbfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Perform query for our hash and distance.
	results := db.Find(hash, *dist)

	// Do we have results?
	if len(results) == 0 {
		fmt.Printf("No matches where found.\n")
		return
	}

	// Display results.
	for _, res := range results {
		fmt.Printf("%d %s\n", res.Distance, filepath.Join(db.Root, res.Path))
	}
}

func getHash(hf imghash.HashFunc, file string) (uint64, error) {
	fd, err := os.Open(file)
	if err != nil {
		return 0, err
	}

	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return 0, err
	}

	return hf(img), nil
}

func parseArgs() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <filename>\n\n", os.Args[0])
		fmt.Printf("   -db: Location for the hash database. Alternatively, this can be set\n" +
			"        in the IMGHASH_DB environment variable.\n")
		fmt.Printf(" -dist: Hamming Distance to use when matching hashes. Defaults to 5.\n" +
			"        Smaller distance provides more restrictive matches.\n" +
			"        Distance can be in the range: 0-64.\n" +
			"        Where 0 means the hash is identical and 64 means we match every hash.\n")
		fmt.Printf("    -v: Display version information.\n")
	}

	version := flag.Bool("v", false, "Display version information.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	file = filepath.Clean(flag.Arg(0))

	if *dist > 64 {
		*dist = 64
	}
}
