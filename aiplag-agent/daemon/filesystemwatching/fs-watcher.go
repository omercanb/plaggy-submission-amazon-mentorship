// Containts the functionality for watching the users filesystem and calling functions on filesystem events
package filesystemwatching

import (
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/fsnotify/fsnotify"
)

// FSWatcher is a wrapper on top of fsnotify.Watcher that ignores umimportant events and relays important events to the
// event handler
type FSWatcher struct {
	watcher                  *fsnotify.Watcher
	watchedDirectories       []string
	ignoredGlobsForDirectory map[string][]string
	ignoredFiles             []string
	eventHandler             FSEventHandler
}

func NewFSWatcher(eventHandler FSEventHandler) *FSWatcher {
	fsw := &FSWatcher{}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	fsw.watcher = watcher
	fsw.eventHandler = eventHandler
	return fsw
}

// Starts watching for the directory non recursively
func (fsw *FSWatcher) AddDirectoryNonRecursive(dir string) error {
	err := fsw.watcher.Add(dir)
	return err
}

// Starts watching the directory recursively
// Currently no ignoring for ignored globs
func (fsw *FSWatcher) AddDirectory(dir string) error {
	err := fsw.watcher.Add(dir)
	if err != nil {
		return err
	}
	entriesInDirectory, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entriesInDirectory {
		if entry.IsDir() {
			err = fsw.AddDirectory(filepath.Join(dir, entry.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Stops watching of the directory recursively
func (fsw *FSWatcher) StopWatchingDirectory(dir string) error {
	err := fsw.watcher.Remove(dir)
	if err != nil {
		return err
	}
	entriesInDirectory, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entriesInDirectory {
		if entry.IsDir() {
			err = fsw.StopWatchingDirectory(filepath.Join(dir, entry.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Adds the file to be ignored for filesystem events
func (fsw *FSWatcher) IgnoreFile(path string) {
	fsw.ignoredFiles = append(fsw.ignoredFiles, path)
}

func (fsw *FSWatcher) Close() {
	fsw.watcher.Close()
}

// Starts running the watcher in a blocking fashion
func (fsw *FSWatcher) Run() {
	for {
		select {
		case event, ok := <-fsw.watcher.Events:
			if !ok {
				return
			}

			switch {
			case event.Has(fsnotify.Create):
				if fsw.shouldNotifyForPath(event.Name) {
					fsw.eventHandler.FileAdded(event.Name)
				}

			case event.Has(fsnotify.Remove):
				if fsw.shouldNotifyForPath(event.Name) {
					fsw.eventHandler.FileDeleted(event.Name)
				}
			case event.Has(fsnotify.Rename):
				if fsw.shouldNotifyForPath(event.Name) {
					fsw.eventHandler.FileRenamed(event.Name)
				}
			case event.Has(fsnotify.Write):
				if fsw.shouldNotifyForPath(event.Name) {
					fsw.eventHandler.FileModified(event.Name)
				}
			}

		case err, ok := <-fsw.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

// Checking if an event for a path should be relayed to the event handler
func (fsw *FSWatcher) shouldNotifyForPath(path string) bool {
	if slices.Contains(fsw.ignoredFiles, path) {
		return false
	}
	return true
}
