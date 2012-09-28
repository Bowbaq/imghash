## img-index

img-index accepts the path to a given directory.
It constructs Perceptual hashes for all PNG, GIF and JPEG images it
finds and saves them in a database. This database can be used to quickly
find similar images using the `img-find` tool.

This command should be run at least once before using `index-find`.
After that, one should run it whenever the images in the given path
change or new ones are added. This is kept separate from `img-find`,
because it takes a while to hash large libraries. Specially with
large image files. The tool is smart enough not to re-calculate hashes
for files which have not changed. So repeatedly running the tool over
the same set of files, will only update those files which have been
altered since the last time it was run.


## Database

The generated database is a flat file. It is stored in the location
specified either in the `-db` commandline parameter, or in the `IMGHASH_DB`
environment variable.


### Usage

    go get github.com/jteeuwen/imghash/img-index


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

