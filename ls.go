package main

import (
	"fmt"
	"io/fs"
	"iter"
	"log"
	"os"
	"path/filepath"
)

func envValByKey(key string) string { return os.Getenv(key) }

func dirname() string { return envValByKey("ENV_DIR_NAME") }

type FsLike struct{ fs.FS }

type Dirent struct {
	Dirname  string
	Basename string
}

func (d Dirent) ToPath() string {
	return filepath.Join(d.Dirname, d.Basename)
}

type DirentLike struct{ fs.DirEntry }

func (l DirentLike) ToDirent(dirname string) Dirent {
	return Dirent{
		Dirname:  dirname,
		Basename: l.DirEntry.Name(),
	}
}

func (l FsLike) ReadDir(dirname string) ([]fs.DirEntry, error) {
	return fs.ReadDir(l.FS, dirname)
}

func (l FsLike) ToDirents(dirname string) iter.Seq2[Dirent, error] {
	return func(yield func(Dirent, error) bool) {
		var empty Dirent
		dirents, e := l.ReadDir(dirname)
		if nil != e {
			yield(empty, e)
			return
		}

		for _, dirent := range dirents {
			var entry Dirent = DirentLike{DirEntry: dirent}.ToDirent(dirname)
			if !yield(entry, nil) {
				return
			}
		}
	}
}

func printDirents(dirents iter.Seq2[Dirent, error]) error {
	for dirent, e := range dirents {
		if nil != e {
			return e
		}

		_, e := fmt.Println(dirent.ToPath())
		if nil != e {
			return e
		}
	}

	return nil
}

type RootLike struct {
	Dirname string
}

func (l RootLike) ToFs() fs.FS { return os.DirFS(l.Dirname) }

type Dirname string

func (d Dirname) ToRoot() RootLike {
	return RootLike{Dirname: string(d)}
}

func (d Dirname) PrintDirents() error {
	var fsys fs.FS = d.ToRoot().ToFs()

	var dirents iter.Seq2[Dirent, error] = FsLike{FS: fsys}.ToDirents(".")
	return printDirents(dirents)
}

func env2dirname2dirents2names2stdout() error {
	var dname string = dirname()
	return Dirname(dname).PrintDirents()
}

func main() {
	var e error = env2dirname2dirents2names2stdout()
	if nil != e {
		log.Printf("%v\n", e)
	}
}
