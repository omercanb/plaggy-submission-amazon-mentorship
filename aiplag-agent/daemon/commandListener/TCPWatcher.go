package commandListener

import (
	"aiplag-agent/common/db"
	"aiplag-agent/daemon/filesystemwatching"
	"encoding/binary"
	"io"
	"log"
	"net"
)

// TCPWatcher handles TCP commands from CLI using length-prefixed protocol
type TCPWatcher struct {
	addr             string
	watcher          *filesystemwatching.FSWatcher
	storedFS         *db.FilesystemStore
	edithistoryStore *db.EditHistoryStore
}

// NewTCPWatcher creates a new TCPWatcher
func NewTCPWatcher(addr string, watcher *filesystemwatching.FSWatcher, storedFS *db.FilesystemStore, editHistoryStore *db.EditHistoryStore) *TCPWatcher {
	return &TCPWatcher{
		addr:             addr,
		watcher:          watcher,
		storedFS:         storedFS,
		edithistoryStore: editHistoryStore,
	}
}

// Run starts the TCP server (blocks until CLI connects)
func (tcp *TCPWatcher) Run() {
	listener, err := net.Listen("tcp", tcp.addr)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("TCPWatcher listening on %s", tcp.addr)

	for {
		conn, err := listener.Accept() // Blocking IO
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go tcp.handleConnection(conn)
	}
}

// handleConnection reads length-prefixed messages and executes them
func (tcp *TCPWatcher) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		cmd, payload, err := readMessage(conn)
		if err != nil {
			if err == io.EOF {
				// log.Println("End of message") // <-- treat normal EOF as normal
				return
			}
			log.Println("Connection closed or error:", err)
			return
		}

		var resp byte = 'A' // success by default

		switch cmd {
		case 'W': // watch
			path := payload
			log.Printf("Started watching path: %s", path)
			if err := tcp.storedFS.AddDirectory(path); err != nil {
				log.Printf("failed to add directory to stored filesystem: %s, err: %v", path, err)
				resp = 'R' // reject
			}
			if err := tcp.watcher.AddDirectory(payload); err != nil {
				log.Printf("failed to add directory to watcher: %s, err: %v", path, err)
				resp = 'R' // reject
			}
			tcp.edithistoryStore.AddAssignment(path)
			// TODO
		case 'X': // stop watching
			log.Printf("Stop watching path: %s", payload)
			if err := tcp.watcher.StopWatchingDirectory(payload); err != nil {
				log.Println(err)
				resp = 'R' // reject
			}
		default:
			resp = 'R' // unknown command
		}

		// send response back to CLI
		length := uint16(1) // 1 byte for the response
		lengthBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(lengthBytes, length)

		_, err = conn.Write(append(lengthBytes, resp))
		if err != nil {
			log.Println("Failed to send response:", err)
			return
		}
	}
}

// readMessage reads length-prefixed message from the TCP connection
func readMessage(conn net.Conn) (byte, string, error) {
	header := make([]byte, 2)
	if _, err := conn.Read(header); err != nil {
		return 0, "", err
	}
	length := binary.BigEndian.Uint16(header) // length = command + payload

	data := make([]byte, length)
	total := 0
	for total < int(length) {
		n, err := conn.Read(data[total:])
		if err != nil {
			return 0, "", err
		}
		total += n
	}

	cmd := data[0]              // first byte = command ('W', 'S', 'X') /watch /send
	payload := string(data[1:]) // payload
	return cmd, payload, nil
}
