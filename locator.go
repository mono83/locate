package locate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// Locator is helper structure, used to locate files
type Locator struct {
	Paths      []string // Lookup paths
	Extensions []string // Allowed file extensions
}

// Find searches for file in all configured locations
func (l Locator) Find(name string) (string, error) {
	var files []string
	ext := filepath.Ext(name)
	if len(ext) > 0 {
		files = []string{name}
	} else {
		for _, e := range l.Extensions {
			files = append(files, concatExt(name, e))
		}
	}

	// Empty paths list
	if len(l.Paths) == 0 {
		l.Paths = []string{"./"}
	}

	// Searching in all paths
	misses := []string{}
	for _, path := range l.Paths {
		for _, file := range files {
			n := concatPath(path, file)
			_, err := os.Stat(n)
			if err == nil {
				return n, nil
			} else if !os.IsNotExist(err) {
				return "", err
			}

			misses = append(misses, n)
		}
	}

	return "", locatingError{name: name, misses: misses}
}

// ReadFile locates file and reads it
func (l Locator) ReadFile(name string) ([]byte, error) {
	name, err := l.Find(name)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(name)
}

func concatExt(name, ext string) string {
	if strings.HasPrefix(ext, ".") {
		return name + ext
	}

	return name + "." + ext
}

func concatPath(path, file string) string {
	if e, err := homedir.Expand(path); err == nil {
		path = e
	}
	if strings.HasSuffix(path, "/") {
		return path + file
	}

	return path + "/" + file
}

type locatingError struct {
	name   string
	misses []string
}

func (l locatingError) Error() string {
	return "unable to locate " + l.name + "; lookup locations were " + strings.Join(l.misses, ",")
}
