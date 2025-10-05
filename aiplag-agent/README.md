# AIPlag Agent - Anti-Plagiarism Tool for Coding Assignments: CLI & Daemon

AIPlag Agent is the command-line client and daemon component 
of the AI Plagiarism Prevention for Programming Homeworks project.
It is designed to help educators ensure academic integrity by
tracking assignment edits and detecting suspicious usage patterns
that may indicate plagiarized code.

This repository contains:
- **Command-line Interface:** Used by students to log in, fetch assignments, start code tracking, and submit completed work.
- **Daemon:** Runs in the background to monitor file edits and securely record assignment histories.
- **Common folder:** Contains the common packages for CLI and daemon such as the backend API.

The backend and instructor portal are maintained in seperate repositories.

## Features

**CLI Client:**
- Secure login (via magic link).
- Viewing available assignments and deadlines.
- Initializing assignments: notifies the daemon to begin tracking files.
- Submitting assignments with their edit history.

**Daemon:**
- Monitors file system changes inside tracked assignment directories.
- On each edit, generates and appends a diff with timestamp and integrity hash.
- Stores encrypted copies of diffs and protects them from tampering.

## Prerequisites

This project requires:
- [Go](https://golang.org/doc/install) (version 1.24.5)
- [kardianos/service](https://github.com/kardianos/service):
  Used to install and manage the daemon on Windows, MacOS, and Linux.
- [sergi/go-diff/diffmatchpatch](https://github.com/sergi/go-diff):
  Provides algorithms for computing file diffs and patches.
- [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify):
  Used by the daemon to detect file saves and trigger diff creation.
- [spf13/cobra](https://github.com/spf13/cobra):
  Go framework for building CLI.
- [manifoldco/promptui](https://github.com/manifoldco/promptui):
  Selection prompts in the CLI
- [spf13/viper](https://github.com/spf13/viper):
  Config Management
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3):
  Sqlite Diver 

The dependencies can be installed automatically with ``` go mod tidy ``` or ``` go build ```.

## Build
   
**1. Build the Daemon:**
Navigate to the daemon directory and build the daemon executable.
      
   ```bash
   cd daemon
   go build -o daemon.exe
   ```
      
**2. Build the CLI:**
Now, head over to the cli folder and build your CLI.
      
   ```bash
   cd ../cli
   go build -o cli.exe .
   ```

## Running the Daemon

**Start the Daemon:**
Navigate to the daemon folder and run the following commands to launch the daemon.

```bash
./daemon.exe install
./daemon.exe start
```

**Warning:** Installing the daemon requires admin permission. On Windows, run the terminal as administrator. On MacOS and Linux, use the ```sudo``` keyword.

**Stop the Daemon:**
Run these commands (on daemon directory) to kill the daemon process.

```bash
./daemon.exe stop
./daemon.exe uninstall
```

## Usage
INSERT & EXPLAIN CLI COMMANDS HERE! MELİS, ÖMER CAN


