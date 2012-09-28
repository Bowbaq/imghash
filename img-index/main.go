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
	"path"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"
)

var (
	db     = imghash.NewDatabase()
	hasher = imghash.Average
	dbfile = flag.String("db", "", "")
	cpu    = flag.String("cpu", "", "")
)

func main() {
	parseArgs()

	// Enable CPU profiling of requested.
	if len(*cpu) > 0 {
		fd, err := os.Create(*cpu)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		defer fd.Close()

		pprof.StartCPUProfile(fd)
		defer pprof.StopCPUProfile()
	}

	// Load database file.
	if err := db.Load(*dbfile); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		// Do not quit here. The file may yet have to be created.
	}

	defer db.Save(*dbfile)

	// Begin hashing of files
	var fileCount, fileSize uint64
	start := time.Now()

	fmt.Printf("* Indexing %s...\n", db.Root)
	filepath.Walk(db.Root, func(file string, stat os.FileInfo, e error) (err error) {
		if e != nil {
			return e
		}

		if stat.IsDir() {
			return
		}

		switch strings.ToLower(path.Ext(file)) {
		case ".png", ".jpg", ".jpeg", ".gif":
			if writeDB(file, stat.ModTime().Unix()) {
				fileCount++
				fileSize += uint64(stat.Size())
			}
		}

		return
	})

	// Print some statistics.
	fmt.Println("* Done.")
	fmt.Printf("* %d image(s) (%s) hashed in %s\n",
		fileCount, prettySize(fileSize), time.Since(start))
}

// writeDB creates a hash for the given file and inserts/updates
// the database with its info if necessary.
func writeDB(file string, modtime int64) bool {
	fmt.Printf("* %s\n", file)

	sfile, err := filepath.Abs(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", file, err)
		return false
	}

	sfile = strings.Replace(sfile, db.Root, "", 1)

	if !db.IsNew(sfile, modtime) {
		return false
	}

	hash, err := getHash(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", file, err)
		return false
	}

	db.Set(sfile, modtime, hash)
	return true
}

// getHash creates a perceptual hash for the given file.
func getHash(file string) (uint64, error) {
	fd, err := os.Open(file)
	if err != nil {
		return 0, err
	}

	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return 0, err
	}

	return hasher(img), nil
}

// prettySize returns a human-friendly version of the given
// file size in bytes.
func prettySize(size uint64) string {
	var unit int
	fsize := float64(size)

	for fsize >= 1024.0 {
		fsize /= 1024
		unit++
	}

	units := [...]string{"byte(s)", "kb", "mb", "gb", "tb", "pb", "yb"}
	return fmt.Sprintf("%.2f %s", fsize, units[unit])
}

// parseArgs processes and validates commandline arguments.
func parseArgs() {
	var err error

	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <directory path>\n\n", os.Args[0])
		fmt.Printf("  -db: Location for the hash database. Alternatively, this can be set\n" +
			"       in the IMGHASH_DB environment variable.\n")
		fmt.Printf(" -cpu: File to write CPU profile to.\n")
		fmt.Printf("   -v: Display version information.\n")
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

	db.Root = filepath.Clean(flag.Arg(0))
	db.Root, err = filepath.Abs(db.Root)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
