# Plaggy Backend

The **Plaggy Backend** is the API server for the Plagiarism Prevention for Programming Homework system.
It provides REST endpoints for the instructor portal and integrates with the student-side **CLI & daemon agent** to enable assignment distribution, submission, and plagiarism detection.

---

## Features

* **Authentication & Authorization**

  * Magic link login for students and instructors.
  * Session handling and middleware-based authentication.

* **Assignments**

  * Create and manage assignments.
  * Send assignments to students.
  * Handle submissions, including file uploads and diffs.

* **Detection & Flagging**

  * Rule-based plagiarism detection engine (e.g., "no deletions," "typing too fast").
  * Centralized flag repository for suspicious edits.

* **Infrastructure**

  * Built-in mail server for sending magic login links.
  * Postgres persistence layer with repository pattern.

---

## Directory Structure

```text
.
├── api/              # API layer – route handlers and auth helpers
│   └── routeHandles/ # Individual HTTP handlers (assignments, login, submissions, etc.)
├── core/             # Core utilities for file processing & token handling
├── flagging/         # Plagiarism detection rules and rule engine
├── mail-server/      # Outbound email (magic links, notifications)
├── middleware/       # HTTP middleware (auth, etc.)
├── models/           # Data models, domain objects, mappers, and responses
├── repository/       # Database repositories for persistence
├── scripts/          # Database scripts and mock data generators
├── security/         # Password hashing & security utilities
├── server/           # HTTP server setup and routing
├── service/          # Application services (e.g., file building)
├── tmp/              # Temporary build files
└── main.go           # Entry point of the backend server
```

---

## Prerequisites

* **Go 1.24.5+**

  * Windows: `choco install golang`
  * macOS: `brew install golang`
  * Linux: use your distro’s package manager or download from [golang.org](https://golang.org/doc/install).

* Recommended for development:
  [air](https://github.com/air-verse/air) for hot-reloading.

  ```bash
  go install github.com/air-verse/air@latest
  ```

---

## Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/plagai/plagai-backend.git
cd plaggy-backend
go get .
```

---

## Running

**With air (recommended for dev):**

```bash
air
```

**Without air:**

```bash
go run main.go
```

> Note: without hot reload, you need to restart the server manually after code changes.

---

## Development Notes

* Database models live in `models/database/`.
* Domain models (business logic representations) are in `models/domain/`.
* Repository layer abstracts persistence.
* Detection rules can be added in `flagging/rules/`.

---

## Related Repositories

* [Plaggy Agent (CLI + Daemon)](https://github.com/plagai/aiplag-agent.git) – Student-side tracking and submission tool.
* Instructor portal (web frontend) (https://github.com/plagai/plagai-frontend.git).

