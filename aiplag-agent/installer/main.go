package main

import (
	"aiplag-agent/common/config"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
)

func main() {
	isWindows := runtime.GOOS == "windows"
	// Create directory
	if err := os.MkdirAll(config.AppDataDir(), 0755); err != nil {
		fmt.Println("Error creating directory:", config.AppDataDir(), err)
		return
	}

	// Make the directory readable and writable by all
	if err := os.Chmod(config.AppDataDir(), 0777); err != nil {
		fmt.Println("Failed to set permissions on directory:", config.AppDataDir(), err)
		return
	}

	// Check if config file exists
	if _, err := os.Stat(config.ConfigPath()); os.IsNotExist(err) {
		// Create empty file with 0666 permissions (read/write for all)
		file, err := os.OpenFile(config.ConfigPath(), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Failed to create config file:", err)
			return
		}
		file.Close()
		// Make sure permissions are correct
		if err := os.Chmod(config.ConfigPath(), 0666); err != nil {
			fmt.Println("Failed to set permissions on config file:", err)
			return
		}
	}

	dbPath := config.DBPath()
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.OpenFile(dbPath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Failed to create DB file:", err)
			return
		}
		file.Close()
	}

	// Ensure permissions are 0666
	if err := os.Chmod(dbPath, 0666); err != nil {
		fmt.Println("Failed to set permissions on DB file:", err)
		return
	}

	// Build CLI
	fmt.Println("Building CLI...")
	cliOut := filepath.Join(config.UserBinDir(), "plaggy")
	if isWindows {
		cliOut += ".exe"
	}
	err := build("./cli", cliOut)

	// Build daemon
	fmt.Println("Building daemon...")
	daemonOut := config.DaemonExecutablePath()
	err = build("./daemon", daemonOut)
	if err != nil {
		if isWindows {
			fmt.Println("You need to run this installer with a root terminal")
		} else {
			fmt.Println("You need to run this installer with sudo:")
			fmt.Println("    sudo go run ./installer")
		}
	}

	portFile, err := os.Create(config.TCPPortFilePath())
	freePort, err := getFreePort()
	portFile.WriteString(strconv.Itoa(freePort))

	// Run daemon commands using exec
	actions := []string{"stop", "uninstall", "install", "start"}
	for _, action := range actions {
		if err := runDaemonCommand(daemonOut, action); err != nil {
			fmt.Printf("Daemon command %q failed: %v\n", action, err)
		}
	}
	fmt.Println("Started Plaggy Daemon")

	fmt.Println("Installation complete!")
}

func build(pkg, output string) error {
	cmd := exec.Command("go", "build", "-o", output, pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Build failed")
		return err
	}
	return nil
}

// runDaemonCommand runs the daemon binary with the given argument, suppressing stdout/stderr
func runDaemonCommand(daemonPath, arg string) error {
	cmd := exec.Command(daemonPath, arg)
	cmd.Stdout = nil // Supress output
	cmd.Stderr = nil
	return cmd.Run()
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
