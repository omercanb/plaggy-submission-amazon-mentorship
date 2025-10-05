package filesystemwatching_test

// import (
// 	"os"
// 	"path/filepath"
// 	"testing"
// 	"time"

// 	"aiplag-agent/common/db"
// 	"aiplag-agent/daemon/filesystemwatching"
// 	"aiplag-agent/daemon/models"
// )

// // An end to end test for when a user starts up the filesystem watching and does some file editing
// func TestFSWatcherEvents(t *testing.T) {
// 	// Temporary directories and DBs
// 	testDir := "fs_test_dir"
// 	storeDB := "store_test.db"
// 	historyDB := "history_test.db"

// 	// Cleanup old test artifacts
// 	os.RemoveAll(testDir)
// 	if err := os.Mkdir(testDir, 0755); err != nil {
// 		t.Fatalf("failed to create test directory: %v", err)
// 	}
// 	os.Remove(storeDB)
// 	os.Remove(historyDB)

// 	// Initialize stored filesystem
// 	storedFS, err := db.NewFilesystemStore(storeDB)
// 	if err != nil {
// 		t.Fatalf("failed to init stored filesystem: %v", err)
// 	}
// 	defer storedFS.Close()

// 	// Initialize edit history store
// 	editHistory, err := db.NewEditHistoryStore(historyDB)
// 	if err != nil {
// 		t.Fatalf("failed to init edit history: %v", err)
// 	}
// 	defer editHistory.Close()

// 	handler := filesystemwatching.NewDiffingEventHandler(editHistory, storedFS)
// 	watcher := filesystemwatching.NewFSWatcher(handler)
// 	defer watcher.Close()

// 	watcher.IgnoreFile(storeDB)
// 	watcher.IgnoreFile(historyDB)

// 	if err := watcher.AddDirectory(testDir); err != nil {
// 		t.Fatalf("failed to add directory to watcher: %v", err)
// 	}

// 	// Run watcher in background
// 	go watcher.Run()

// 	// Small delay helper to allow fsnotify to catch events
// 	wait := func() { time.Sleep(200 * time.Millisecond) }

// 	// --- Test sequence ---
// 	file1 := filepath.Join(testDir, "file1.txt")

// 	// 1. Add file
// 	if err := os.WriteFile(file1, []byte("Hello World"), 0644); err != nil {
// 		t.Fatalf("failed to write file: %v", err)
// 	}
// 	wait()

// 	// 2. Modify file
// 	if err := os.WriteFile(file1, []byte("Hello Go"), 0644); err != nil {
// 		t.Fatalf("failed to modify file: %v", err)
// 	}
// 	wait()

// 	// 3. Rename file
// 	file1Renamed := filepath.Join(testDir, "file1_renamed.txt")
// 	if err := os.Rename(file1, file1Renamed); err != nil {
// 		t.Fatalf("failed to rename file: %v", err)
// 	}
// 	wait()

// 	// 4. Delete file
// 	if err := os.Remove(file1Renamed); err != nil {
// 		t.Fatalf("failed to delete file: %v", err)
// 	}
// 	wait()

// 	// --- Validate results ---
// 	events, err := editHistory.GetEventsByAssignment(assignmentID)
// 	if err != nil {
// 		t.Fatalf("failed to get events: %v", err)
// 	}

// 	if len(events) == 0 {
// 		t.Fatal("no events were captured")
// 	}

// 	editHistory.DebugPrint()

// 	t.Logf("Captured %d events:", len(events))
// 	for _, e := range events {
// 		t.Logf("Event: %s, File: %s, Patch: %.200s", e.EventType, e.FilePath, e.Patch)
// 		if e.EventType == models.EventModified && len(e.Patch) == 0 {
// 			t.Fatalf("Patch of size 0 for Modified Event")
// 		}
// 	}

// 	// Cleanup
// 	os.RemoveAll(testDir)
// 	os.Remove(storeDB)
// 	os.Remove(historyDB)
// }
