// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package imghash

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

// SearchResult is returned by Database.Find.
type SearchResult struct {
	Path     string // Image path, relative to Database.Root
	Hash     uint64 // Perceptual Image hash.
	Distance uint   // Hamming Distance to search term.
}

// ResultSet holds search results, sorted by Hamming Distance.
type ResultSet []*SearchResult

func (r ResultSet) Len() int           { return len(r) }
func (r ResultSet) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ResultSet) Less(i, j int) bool { return r[i].Distance < r[j].Distance }

// Entry represents a single database entry.
type Entry struct {
	Path    string // Image path, relative to Database.Root
	Hash    uint64 // Perceptual Image hash.
	ModTime int64  // Last-Modified timestamp for this file.
}

// A Database holds a listing of Perceptual hashes, mapped
// to image file paths.
//
// Note: This is a very naive implementation that can benefit
// a great deal from optimization.
type Database struct {
	Root    string   // Database root path.
	Entries []*Entry // List of entries.
}

// NewDatabase creates a new, empty database.
func NewDatabase() *Database { return new(Database) }

// Find finds all entries which have a Hamming Diance <= to the
// specified distance with the given hash.
// The list is sorted by relevance.
func (d *Database) Find(hash uint64, distance uint) ResultSet {
	var rs ResultSet
	var dist uint

	for _, e := range d.Entries {
		dist = Distance(e.Hash, hash)

		if dist <= distance {
			rs = append(rs, &SearchResult{
				Path:     e.Path,
				Hash:     e.Hash,
				Distance: dist,
			})
		}
	}

	sort.Sort(rs)
	return rs
}

// Load loads a database from the given file.
// Leave the filename empty to use the default file.
func (d *Database) Load(file string) (err error) {
	if len(file) == 0 {
		file = os.Getenv("IMGHASH_DB")
	}

	fd, err := os.Open(file)
	if err != nil {
		return
	}

	defer fd.Close()

	r := bufio.NewReader(fd)

	var line []byte
	var entry *Entry

	line, err = r.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return
	}

	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return errors.New("Invalid database file.")
	}

	d.Root = string(line)

	for {
		line, err = r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		line = bytes.TrimSpace(line)
		if len(line) < 34 {
			continue
		}

		entry = new(Entry)
		entry.Path = string(line[33:])

		entry.Hash, err = strconv.ParseUint(string(line[:16]), 16, 64)
		if err != nil {
			return
		}

		entry.ModTime, err = strconv.ParseInt(string(line[18:32]), 16, 64)
		if err != nil {
			return
		}

		d.Entries = append(d.Entries, entry)
	}

	return
}

// Save saves the database to the given file.
// Leave the filename empty to use the default file.
func (d *Database) Save(file string) (err error) {
	if len(file) == 0 {
		file = os.Getenv("IMGHASH_DB")
	}

	fd, err := os.Create(file)
	if err != nil {
		return
	}

	defer fd.Close()

	fmt.Fprintf(fd, "%s\n", d.Root)

	for _, e := range d.Entries {
		fmt.Fprintf(fd, "%016x %015x %s\n", e.Hash, e.ModTime, e.Path)
	}

	return
}

// Set adds the given file if it doesn't already exist.
// Otherwise it overwrites the existing one.
func (d *Database) Set(file string, modtime int64, hash uint64) {
	index := d.IndexFile(file)

	if index == -1 {
		d.Entries = append(d.Entries, &Entry{file, hash, modtime})
		return
	}

	f := d.Entries[index]
	f.ModTime = modtime
	f.Hash = hash
}

// IsNew returns true if the given file has been updated
// since it was last stored in the database.
func (d *Database) IsNew(file string, modtime int64) bool {
	index := d.IndexFile(file)

	if index == -1 {
		// Non-existant entry is always new.
		return true
	}

	return d.Entries[index].ModTime != modtime
}

// IndexFile returns the index for the given file.
func (d *Database) IndexFile(file string) int {
	for i, e := range d.Entries {
		if e.Path == file {
			return i
		}
	}

	return -1
}

// IndexHash returns the indices for files with the given hash.
// There can be more than one of them.
func (d *Database) IndexHash(hash uint64) []int {
	var list []int

	for i, e := range d.Entries {
		if e.Hash == hash {
			list = append(list, i)
		}
	}

	return list
}
