// A simple wrapper over the users filesystem to make files fit to the File interface used by the filesystem store
package filesystemwatching

import "os"

type LocalFile string

func (f LocalFile) Read() (string, error) {
	b, err := os.ReadFile(string(f))
	return string(b), err
}

func (f LocalFile) Path() string {
	return string(f)
}
