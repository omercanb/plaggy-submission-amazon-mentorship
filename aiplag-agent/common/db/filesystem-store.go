// A SQL databse that stores a copy of required files in the watched directories of the user
package db

import (
	"aiplag-agent/daemon/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// FilesystemStore represents a persistent storage for files backed by a SQLite database
type FilesystemStore struct {
	db                  *sql.DB
	openFileStmt        *sql.Stmt
	addOrUpdateFileStmt *sql.Stmt
	getAllFilepathsStmt *sql.Stmt
	deleteFileStmt      *sql.Stmt
}

// StoredFile represents a file stored in the FilesystemStore.
type StoredFile struct {
	Content  string
	Filepath string
}

func (f *StoredFile) Read() (string, error) {
	return f.Content, nil
}

func (f *StoredFile) Path() string {
	return f.Filepath
}

// NewFilesystemStore initializes a new FilesystemStore with a SQLite database
// at the given dbPath. It creates the required schema if it doesn't exist.
func NewFilesystemStore(dbPath string) (*FilesystemStore, error) {
	db, err := InitDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("NewFilesystemStore: failed to create db: %v", err)
	}

	repo := &FilesystemStore{db: db}
	repo.initSchema()
	repo.prepareStatements()

	return repo, nil
}

// AddOrUpdateFile inserts a new file or updates the content if the file already exists.
func (fsstore *FilesystemStore) AddOrUpdateFile(file models.File) error {
	content, err := file.Read()
	if err != nil {
		return err
	}
	path := file.Path()
	_, err = fsstore.addOrUpdateFileStmt.Exec(path, content)
	if err != nil {
		return err
	}
	return nil
}

// GetAllFilepaths returns a list of all file paths stored in the database.
func (fsstore *FilesystemStore) GetAllFilepaths() []string {
	filepaths := []string{}
	rows, err := fsstore.getAllFilepathsStmt.Query()
	if err != nil {
		return filepaths
	}
	defer rows.Close()
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			continue
		}
		filepaths = append(filepaths, path)
	}

	return filepaths
}

// DeleteFile removes a file from the database by path.
func (fsstore *FilesystemStore) DeleteFile(path string) error {
	_, err := fsstore.deleteFileStmt.Exec(path)
	if err != nil {
		log.Printf("DeleteFile: failed to delete file from db %s: %v", path, err)
		return err
	}
	return nil
}

// AddDirectory walks through the given directory, reads all files, and adds or updates
// them in the FilesystemStore.
// This goes againts separation of concerns and should be refactored later
func (fsstore *FilesystemStore) AddDirectory(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open the file and wrap it in a StoredFile
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			log.Printf("AddDirectory: failed to read file %s: %v", path, err)
			return nil // skip this file, continue walking
		}

		file := &StoredFile{
			Filepath: path,
			Content:  string(contentBytes),
		}

		if err := fsstore.AddOrUpdateFile(file); err != nil {
			log.Printf("AddDirectory: failed to add/update file %s: %v", path, err)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// prepareStatements prepares all SQL statements used by the FilesystemStore.
func (fsstore *FilesystemStore) prepareStatements() {
	var err error
	fsstore.openFileStmt, err = fsstore.db.Prepare(`SELECT content FROM files WHERE path == ?`)
	if err != nil {
		log.Fatalf("openFileStmt: %v", err)
	}
	fsstore.addOrUpdateFileStmt, err = fsstore.db.Prepare(`INSERT INTO files (path, content) VALUES (?, ?) ON CONFLICT(path) DO UPDATE SET content=excluded.content`)
	if err != nil {
		log.Fatalf("updateFileStmt: %v", err)
	}
	fsstore.getAllFilepathsStmt, err = fsstore.db.Prepare(`SELECT path FROM files`)
	if err != nil {
		log.Fatalf("getAllFilepathsStmt: %v", err)
	}
	fsstore.deleteFileStmt, err = fsstore.db.Prepare(`DELETE FROM files WHERE path == ?`)
	if err != nil {
		log.Fatalf("deleteFileStmt: %v", err)
	}
}

// initSchema creates the necessary tables in the SQLite database if they don't exist.
func (fsstore *FilesystemStore) initSchema() {
	initSchemaStmt, err := fsstore.db.Prepare(`
		CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL UNIQUE,
		content TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatalf("initSchemaStmt: %v", err)
	}
	initSchemaStmt.Exec()
	initSchemaStmt.Close()
}

// Close closes all prepared statements and the underlying database connection.
func (fsstore *FilesystemStore) Close() error {
	usedStatements := []*sql.Stmt{fsstore.addOrUpdateFileStmt, fsstore.openFileStmt, fsstore.getAllFilepathsStmt, fsstore.deleteFileStmt}
	for _, stmt := range usedStatements {
		if stmt != nil {
			stmt.Close()
		}
	}
	if fsstore.db != nil {
		return fsstore.db.Close()
	}
	return nil
}

// Open retrieves the stored file with the given path. Returns an error if the file does not exist.
func (fsstore *FilesystemStore) Open(path string) (*StoredFile, error) {
	var content string
	row := fsstore.openFileStmt.QueryRow(path)
	if err := row.Scan(&content); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No such file: %s, %w", path, err)
		}
		return nil, err
	}
	return &StoredFile{Content: content}, nil
}

// DebugPrint prints all stored files with their IDs, paths, and truncated content
// (first 50 characters) for debugging purposes.
func (fs *FilesystemStore) DebugPrint() {
	rows, err := fs.db.Query(`SELECT id, path, content FROM files`)
	if err != nil {
		log.Printf("DebugPrint: failed to query for everything: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var path, content string
		err := rows.Scan(&id, &path, &content)
		if err != nil {
			log.Printf("DebugPrint: failed to scan row: %v", err)
			continue
		}
		truncatedContent := content
		if len(truncatedContent) > 50 {
			truncatedContent = truncatedContent[:50]
		}
		fmt.Printf("ID: %d, Path: %s, Content: %s\n", id, path, truncatedContent)
	}

	if err = rows.Err(); err != nil {
		log.Printf("DebugPrint: rows iteration error: %v", err)
	}
}

// func main() {
// 	repo, err := NewDBFileSystemRepo("file_history.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("DB connection:", repo.db)
// }
