# OpenShelf Roadmap

## Phase 1: Engine — EPUB-only working product

- [ ] **EPUB metadata extractor** — new `internal/epub` package using `archive/zip` + `encoding/xml`. Extracts title, authors, language, publisher, and identifiers (ISBN etc.) into a `Canonical`.
- [ ] **Complete `GetBookID()`** — implement the cache-miss path: extract EPUB metadata → `Canonical.ComputeBookID()` → write `.metadata/{hash}.json` and `.canonical/{bookId}.json`. For non-EPUB files or extraction failure: return an error directing the user to `openshelf register`.
- [ ] **`register` command** — lets the user manually supply canonical fields for a book file (or non-file publication), writes the metadata and canonical files. This is also the primitive for non-file publications.
- [ ] **Fix `cli/fetch.go`** — replace deleted `state` package with `library` + `model`, fix ignored error.
- [ ] **Implement `cmd/merge`** — read two state file paths, merge with `merge.MergeReadingStates()`, write result.
- [ ] **End-to-end smoke test** of all commands with a real EPUB.

## Phase 2: KOReader plugin

- [ ] Replace all Lua hash/library/state logic with calls to the `openshelf` binary.
- [ ] `onReaderReady` → `openshelf fetch <file>`, parse output, restore position.
- [ ] `onCloseDocument` → `openshelf update <file> <type> <value>`.
- [ ] Plugin becomes thin glue — no direct file I/O or hashing in Lua.

## Design notes

- **BookId derivation**: All publications (file-based and non-file-based) use the Canonical Identity Descriptor approach. `fileSha256` is a cache key only, not the identity. This ensures stability even if EPUB metadata is edited by the user.
- **Metadata cache flow**: `sha256(file bytes)` → check `.metadata/{hash}.json` → hit: return `bookId` / miss: extract metadata → build `Canonical` → `ComputeBookID()` → write `.metadata` and `.canonical` files.
- **Spec**: Defer updating `spec/BookID.md` and related docs until after a working implementation exists.
