package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func AppDataDir() string {
	if runtime.GOOS == "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// fallback to current directory
			homeDir = "."
		}
		return filepath.Join(homeDir, ".plagai")
	}
	return "/var/lib/plaggy"
}

func DBPath() string {
	path := filepath.Join(AppDataDir(), "app.db")
	return path
}

func ConfigPath() string {
	path := filepath.Join(AppDataDir(), "config.yaml")
	return path
}

func DaemonLogPath() string {
	path := filepath.Join(AppDataDir(), "daemon.log")
	return path
}

func DaemonExecutablePath() string {
	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(AppBinDir(), "plaggydaemon.exe")
	} else {
		path = filepath.Join(AppBinDir(), "plaggydaemon")
	}
	return path
}

func AppBinDir() string {
	path := filepath.Join(AppDataDir(), "bin")
	return path
}

func TCPPortFilePath() string {
	path := filepath.Join(AppDataDir(), "daemon.port")
	return path
}

func UsedTCPAddress() (string, error) {
	data, err := os.ReadFile(TCPPortFilePath())
	if err != nil {
		return "", err
	}
	port := strings.TrimSpace(string(data))
	return "127.0.0.1:" + port, nil
}

func UserBinDir() string {
	if runtime.GOOS == "windows" {
		return AppBinDir()
	}

	// Assume Unix
	return "/usr/local/bin"
}
