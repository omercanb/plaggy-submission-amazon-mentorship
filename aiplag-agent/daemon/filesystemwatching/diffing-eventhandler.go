package filesystemwatching

import (
	"aiplag-agent/common/db"
	"aiplag-agent/daemon/models"
	"log"
	"os"
)

// DiffingEventHandler handles filesystem events by updating the stored
// filesystem state and recording them in the edit history.
//
// It delegates to EditHistoryEventHandler for diffing and history storage,
// and to FilesystemStore for maintaining the canonical copy of file contents.
type DiffingEventHandler struct {
	editHistoryHandler EditHistoryEventHandler
	fsStore            *db.FilesystemStore
}

// EditHistoryEventHandler provides the underlying logic for diffing file states
// and writing events into the EditHistoryStore.
type EditHistoryEventHandler struct {
	editHistoryStore *db.EditHistoryStore
	// assignmentID     int
	storedFS   *db.FilesystemStore
	fileDiffer *FileDiffer
}

// NewDiffingEventHandler constructs a DiffingEventHandler wired to the provided
// stores and filesystem. A new FileDiffer is created internally.
func NewDiffingEventHandler(eh *db.EditHistoryStore, storedFS *db.FilesystemStore) *DiffingEventHandler {
	return &DiffingEventHandler{
		editHistoryHandler: EditHistoryEventHandler{
			editHistoryStore: eh,
			storedFS:         storedFS,
			fileDiffer:       NewFileDiffer(),
		},
		fsStore: storedFS,
	}
}

// FileAdded stores the new file in the FilesystemStore and records an "added"
// event in the EditHistoryStore. If reading or storing fails, errors are logged
// and the event may not be recorded.
func (h *DiffingEventHandler) FileAdded(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("FileAdded: failed to read file %s: %v", path, err)
		return
	}
	file := &db.StoredFile{
		Content:  string(content),
		Filepath: path,
	}
	err = h.fsStore.AddOrUpdateFile(file)
	if err != nil {
		log.Printf("FileAdded: failed to add file to db %s: %v", path, err)
	}
	err = h.editHistoryHandler.editHistoryStore.AddEvent(path, models.EventAdded, "")
	if err != nil {
		log.Printf("FileAdded: failed to log add event for %s: %v", path, err)
	}
}

// FileDeleted records a "deleted" event in the EditHistoryStore. No file
// contents are touched.
func (h *DiffingEventHandler) FileDeleted(path string) {
	err := h.editHistoryHandler.editHistoryStore.AddEvent(path, models.EventDeleted, "")
	if err != nil {
		log.Printf("FileDeleted: failed to log delete event for %s: %v", path, err)
	}
}

// FileRenamed records a "renamed" event in the EditHistoryStore for the old path.
func (h *DiffingEventHandler) FileRenamed(oldPath string) {
	err := h.editHistoryHandler.editHistoryStore.AddEvent(oldPath, models.EventRenamed, "")
	if err != nil {
		log.Printf("FileRenamed: failed to log rename event for %s: %v", oldPath, err)
	}
}

// FileModified computes a diff between the stored file state and the local file,
// then records a "modified" event with the patch in the EditHistoryStore.
// If diffing fails, the event may be missing or incomplete.
func (h *DiffingEventHandler) FileModified(path string) {
	log.Printf("EditHistoryEventHandler FileModified for %v", path)

	oldFileState, err := h.editHistoryHandler.storedFS.Open(path)
	if err != nil {
		log.Printf("FileModified: failed to open old file state for %s: %v", path, err)
		return
	}
	file := LocalFile(path)
	filePatch, err := h.editHistoryHandler.fileDiffer.Diff(oldFileState, file)
	if err != nil {
		log.Printf("FileModified: failed diffing for %v", path)
	}
	err = h.editHistoryHandler.editHistoryStore.AddEvent(path, models.EventModified, filePatch)
	if err != nil {
		log.Printf("FileModified: failed to log modify event for %s: %v", path, err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("FileAdded: failed to read file %s: %v", path, err)
		return
	}
	fileToStore := &db.StoredFile{
		Content:  string(content),
		Filepath: path,
	}
	err = h.fsStore.AddOrUpdateFile(fileToStore)
	if err != nil {
		log.Printf("FileAdded: failed to add file to db %s: %v", path, err)
	}
}
