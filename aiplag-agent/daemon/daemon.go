package main

import (
	"aiplag-agent/common/config"
	"aiplag-agent/common/db"
	"aiplag-agent/daemon/commandListener"
	"aiplag-agent/daemon/filesystemwatching"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

type Daemon struct {
	watcher     *filesystemwatching.FSWatcher
	logFile     *os.File
	editHistory *db.EditHistoryStore
}

// Start is called when the service starts
func (d *Daemon) Start(s service.Service) error {
	var err error

	logPath := config.DaemonLogPath()
	err = os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		fmt.Println("Failed to create log directory:", err)
	}
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file, fallback to stdout:", err)
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(logFile)
		d.logFile = logFile
	}

	freePort, err := getFreePort()
	if err != nil {
		log.Println("Failed to get free TCP port:", err)
		return err
	}

	portFilePath := config.TCPPortFilePath()
	err = os.WriteFile(portFilePath, []byte(fmt.Sprintf("%d", freePort)), 0644)
	if err != nil {
		log.Println("Failed to write TCP port file:", err)
		return err
	}
	log.Printf("Daemon TCP port written to %s (port %d)\n", portFilePath, freePort)

	log.Println("Daemon starting...")
	dbPath := config.DBPath()

	// Initialize stores
	storedFS, err := db.NewFilesystemStore(dbPath)
	if err != nil {
		log.Println("Failed to initialize stored filesystem:", err)
		return err
	}

	d.editHistory, err = db.NewEditHistoryStore(dbPath)
	if err != nil {
		log.Println("Failed to initialize edit history:", err)
		return err
	}

	// Event handler + watcher
	diffingHandler := filesystemwatching.NewDiffingEventHandler(d.editHistory, storedFS)
	d.watcher = filesystemwatching.NewFSWatcher(diffingHandler)

	// TCP command listener (new signature includes editHistory)
	tcpAdress, err := config.UsedTCPAddress()
	if err != nil {
		log.Println("No tcp address file")
		return nil
	}
	socket := commandListener.NewTCPWatcher(tcpAdress, d.watcher, storedFS, d.editHistory)

	// Start components
	go d.watcher.Run()
	go socket.Run()

	assignmentPaths, err := d.editHistory.GetAssignmentFullPaths()
	for _, path := range assignmentPaths {
		d.watcher.AddDirectory(path)
	}

	log.Println("Daemon started successfully.")
	return nil
}

// Stop is called when the service stops
func (d *Daemon) Stop(s service.Service) error {
	log.Println("Daemon stopping...")

	if d.logFile != nil {
		_ = d.logFile.Close()
		d.logFile = nil
	}
	return nil
}

func GetServiceConfig() *service.Config {
	svcConfig := &service.Config{
		Name:        "PlaggyDaemon",
		DisplayName: "Plaggy Daemon",
		Description: "Daemon for monitoring filesystem changes and assignments.",
		Option: map[string]any{
			"StartType": "automatic",
			// Windows is the only OS where StartType: "automatic" actually auto-starts without extra commands, after the service is installed.
			// linux and macos require some shell commands of their own for launch on boot, could work on this later
			// but it seems like that will be done when we actually have installer stuff instead of dev commands
		},
	}
	return svcConfig
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
