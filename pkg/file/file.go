package file

import (
	"os"
)

// Exists returns a bool indicating if the specified file exists or not
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
