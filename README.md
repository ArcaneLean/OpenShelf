# OpenShelf

**An open, filesystem-based reading state standard.
Readers plug in. Servers are optional.**

---

## What is OpenShelf?

OpenShelf is an **open specification and ecosystem** for syncing reading progress across devices, readers, and file formats — **without requiring a central server**.

At its core, an OpenShelf library is just a **folder**:

```
/OpenShelf/
  library.json
  /books/
  /.state/
```

Put that folder in **any synced filesystem** (cloud storage, Syncthing, rclone, USB, etc.), and compatible readers can share reading progress automatically.

---

## The Problem OpenShelf Solves

Today’s reading ecosystems are fragmented:

* Progress is locked to a specific reader or service
* Sync usually requires a proprietary server
* Cross-reader support is rare or impossible
* Power users resort to hacks and scripts

OpenShelf flips this model.

---

## The OpenShelf Model

> **A reading library is just a folder.
> Readers adapt to it.
> Servers are optional.**

Key ideas:

* 📁 **Filesystem-first** — no mandatory backend
* 🔓 **Open & inspectable** — JSON, not opaque databases
* 🔁 **Reader-agnostic** — KOReader, Kobo, desktop readers, future apps
* ☁️ **Sync-agnostic** — use Dropbox, Drive, WebDAV, Syncthing, or none
* 🧱 **Composable** — adapters are replaceable, not privileged

---

## What OpenShelf Is (and Is Not)

### OpenShelf **is**

* A **specification** for storing books and reading state
* A **shared contract** between readers and devices
* A foundation for cross-device reading sync
* Friendly to offline-first workflows

### OpenShelf **is not**

* ❌ A reader application
* ❌ A cloud service
* ❌ A sync engine
* ❌ A walled garden

---

## How It Works (High Level)

1. **Books** live in `/books/`
   Users organize them however they like.

2. **Reading state** lives in `/.state/`
   One JSON file per book, keyed by a content hash.

3. **Adapters** translate between:

   * OpenShelf’s neutral state
   * A specific reader or device’s internal format

4. **Sync** is handled externally
   OpenShelf reacts to file changes — it does not manage sync itself.

---

## Reading State (Example)

```json
{
  "specVersion": "0.1.0",
  "bookId": "sha256_abc123",
  "updatedAt": "2026-02-03T10:15:00Z",
  "location": {
    "percentage": 42.3,
    "epubcfi": "/6/2[chapter1]!/4/2/14",
    "page": 123
  }
}
```

Multiple location formats may coexist.
No single representation is considered authoritative.

---

## Adapters

An **OpenShelf adapter** bridges a reader or device to an OpenShelf library.

Examples:

* KOReader adapter
* Kobo device adapter
* Desktop EPUB reader adapter
* “Filesystem only” adapter (read-only)

Adapters declare **responsibility levels** that define:

* What data they can read
* What data they can write
* What features they support

Adapters are:

* Optional
* Replaceable
* Non-exclusive

---

## Why No Mandatory Server?

Because files already sync.

OpenShelf intentionally:

* **Does not reinvent sync**
* **Does not centralize trust**
* **Does not lock users in**

A server can exist *later* — but only as a consumer of the same filesystem format.

---

## Current Status

* 🧪 **Specification:** v0.1.0-draft
* 📄 **Core focus:** book identity + reading progress
* 🚫 **Explicitly out of scope (for now):**

  * Annotations
  * Highlights
  * Bookmarks
  * Metadata indexing
  * Multi-user libraries

---

## Who Is This For?

* Readers who use **multiple devices**
* People who switch between **different reading apps**
* Developers tired of re-implementing sync
* Anyone who wants **ownership of their reading data**

---

## Project Structure (planned)

```
openshelf-spec/        # Core specification
openshelf-core/        # Neutral state engine
openshelf-cli/         # Setup & validation tool
openshelf-adapters/
  ├─ koreader/
  ├─ kobo/
  └─ desktop/
```

---

## Philosophy

* Files over APIs
* Standards over platforms
* Ownership over convenience
* Interoperability over control

---

## License

OpenShelf is intended to be **free, open, and implementation-neutral**.
(Exact license TBD.)

---

## One-Line Summary

> **OpenShelf lets your reading progress follow you — not your reader.**
