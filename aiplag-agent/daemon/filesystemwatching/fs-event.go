// The filesystem events that represent what happens on the watched parts of the user's filesystem
// Currently acting as a shallow wrapper on top of fsnotify.Op
package filesystemwatching

type FSEventType int

const (
	FileAdded FSEventType = iota
	FileDeleted
	FileRenamed
	FileModified
)

func (t FSEventType) String() string {
	switch t {
	case FileAdded:
		return "Added"
	case FileDeleted:
		return "Deleted"
	case FileRenamed:
		return "Renamed"
	case FileModified:
		return "Modified"
	default:
		return "Unknown"
	}
}

type FSEvent struct {
	Type    FSEventType
	Path    string // current or original path
	OldPath string // only set if renamed
}

type FSEventHandler interface {
	FileAdded(path string)
	FileDeleted(path string)
	// A rename is always sent with the old path as Event.Name, and a Create event will be sent with the new name.
	FileRenamed(oldPath string)
	FileModified(path string)
}
