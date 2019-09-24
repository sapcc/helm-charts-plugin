package charts

import (
	"fmt"
	"os"
	"path"
)

// EnsureFileExists ensures all directories and the file itself exist.
func EnsureFileExists(absPath, filename string) (*os.File, error) {
	if err := os.MkdirAll(absPath, os.ModeDir); err != nil {
		return nil, err
	}

	filepath := path.Join(absPath, filename)
	fmt.Println("Using file: ", filepath)

	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return f, f.Truncate(0)
}
