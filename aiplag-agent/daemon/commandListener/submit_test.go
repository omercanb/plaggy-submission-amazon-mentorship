package commandListener

import (
	"aiplag-agent/common/db"
	"aiplag-agent/daemon/filesystemwatching"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Instructions on how to run this test
// 1) Login with magic link
// 2) Copy paste the token you recieved in the config.yaml into the token variable below
func TestSubmitEdits(t *testing.T) {
	// Temporary directories and DBs
	testDir := "fs_test_dir"
	storeDB := "store_test.db"
	historyDB := "history_test.db"

	// Cleanup old test artifacts
	os.RemoveAll(testDir)
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	os.Remove(storeDB)
	os.Remove(historyDB)

	// Initialize stored filesystem
	storedFS, err := db.NewFilesystemStore(storeDB)
	if err != nil {
		t.Fatalf("failed to init stored filesystem: %v", err)
	}
	defer storedFS.Close()

	// Initialize edit history store
	editHistory, err := db.NewEditHistoryStore(historyDB)
	if err != nil {
		t.Fatalf("failed to init edit history: %v", err)
	}
	defer editHistory.Close()

	assignmentID := 1
	handler := filesystemwatching.NewDiffingEventHandler(editHistory, storedFS)
	watcher := filesystemwatching.NewFSWatcher(handler)
	defer watcher.Close()

	watcher.IgnoreFile(storeDB)
	watcher.IgnoreFile(historyDB)

	if err := watcher.AddDirectory(testDir); err != nil {
		t.Fatalf("failed to add directory to watcher: %v", err)
	}
	testDirFullPath, err := filepath.Abs(testDir)
	t.Log(testDirFullPath)
	if err = editHistory.MapFullPathToAssignmentID(testDirFullPath, assignmentID); err != nil {
		t.Fatal(err)
	}

	// Run watcher in background
	go watcher.Run()

	tcpWatcher := NewTCPWatcher("127.0.0.1:9090", watcher, storedFS, editHistory)
	go tcpWatcher.Run()

	// Small delay helper to allow fsnotify to catch events
	wait := func() { time.Sleep(200 * time.Millisecond) }

	// --- Test sequence ---
	file1 := filepath.Join(testDir, "file1.txt")

	// 1. Add file
	if err := os.WriteFile(file1, []byte("Hello World"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	wait()

	// 2. Modify file
	if err := os.WriteFile(file1, []byte("Hello Go"), 0644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}
	wait()

	// 3. Rename file
	file1Renamed := filepath.Join(testDir, "file1_renamed.txt")
	if err := os.Rename(file1, file1Renamed); err != nil {
		t.Fatalf("failed to rename file: %v", err)
	}
	wait()

	// 4. Delete file
	if err := os.Remove(file1Renamed); err != nil {
		t.Fatalf("failed to delete file: %v", err)
	}
	wait()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNydEBicnQuY29tIiwiZXhwIjoxNzU2NzI2ODY5LCJpYXQiOjE3NTY3MjU5Njl9.qUDvfvIvoSwbn4ZDdaKTSexJFNkXQtKj-wQdJhXIXbg"
	err = tcpWatcher.submit(testDirFullPath, token)
	if err != nil {
		t.Fatalf("%v", err)
	}
	// t.Log("Submission successful")
}
