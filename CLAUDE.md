# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

OpenShelf is an open, filesystem-based reading state standard and CLI tool. It defines how reading progress (locations, timestamps) is stored and synced across devices and reader apps without a central server. The repository contains:

- **`/spec/`** — Normative markdown specifications (Core, BookID, RSF, NRSM, AdapterLevels)
- **`/engine/`** — Go CLI implementation (`github.com/ArcaneLean/openshelf`, Go 1.25)
- **`/adapters/`** — Planned reader-specific integrations (KOReader stub)

## Commands

All commands run from `/engine/`:

```bash
# Build
go build -o openshelf .

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/model

# Vet
go vet ./...
```

## Architecture

### Core Data Flow

```
Book File → SHA-256 hash → metadata lookup → bookId → .state/{bookId}.json → ReadingState
```

1. `pkg/utils.HashFile()` — hashes a file to derive its content identity
2. `internal/metadata` — maps file hash → bookId (stored in `.metadata/{hash}.json`)
3. `internal/library` — resolves library root from `OPENSHELF_LIBRARY` env var; provides `.state/`, `.metadata/`, `.canonical/` paths
4. `internal/model` — `ReadingState` and `Location` structs; JSON I/O; location filtering
5. `internal/merge` — last-write-wins merge of two `ReadingState` values (validates bookIds match)
6. `internal/canonical` — for non-file publications: normalizes metadata → deterministic JSON → SHA-256 bookId
7. `internal/identity` — top-level bookId resolution from a file path (ties above together)

### CLI Layer

- **`cmd/`** — Cobra command definitions (root, fetch, show, update, merge)
- **`internal/cli/`** — Implementation logic called by commands

Commands: `fetch [book-file|book-id]`, `show <state-file>`, `update <state-file> <type> <value>`, `merge` (stub).

### Spec ↔ Code Alignment

The spec defines three location types: `percentage`, `epubcfi`, `pageNumber`, `timeSeconds`. `model.IsInteroperable()` validates against these. The `Canonical` struct in `internal/canonical` maps directly to the IDO (Identity Descriptor Object) defined in `spec/BookID.md`.

## Key Conventions

- BookIds are lowercase hex SHA-256 strings — either of file content (for file-based books) or of deterministic canonical JSON (for non-file publications).
- State files live at `{library}/.state/{bookId}.json`; metadata at `{library}/.metadata/{hash}.json`.
- Conflict resolution is always last-write-wins on `updatedAt` timestamps (RFC3339).
- The `OPENSHELF_LIBRARY` environment variable sets the library root path.
