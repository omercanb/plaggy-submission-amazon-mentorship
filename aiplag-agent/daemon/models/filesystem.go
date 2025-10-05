package models

type File interface {
	Read() (string, error)
	Path() string
}

type Filesystem interface {
	Open(path string) (File, error)
}
