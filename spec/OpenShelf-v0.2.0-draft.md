# OpenShelf Specification

## Version: 0.2.0-draft (WIP)

---

## 1. Directory Structure

A compliant library **MUST** follow this structure:

```
/OpenShelf/
  library.json
  /books/
  /.state/
```

### Rules

* **Book files**:
  Adapters **MAY** read book files from any location on the device.
  The `/books/` directory is optional, provided for implementations that want a central location for all books.

* **Reading state files**:
  Reading state files **MUST** be stored under `/.state/`.
  All adapters **MUST** read and write state files from this location to ensure interoperability.

* **Dot-prefixed directories**:
  Directories whose names start with `.` are **implementation-owned**.
  Implementations **MUST NOT** require users to manually edit dot-prefixed directories.

* **Optional directories**:
  Implementations **MAY** create additional directories as needed, provided they do not conflict with the mandated structure.

---

## 2. Book Placement

* The directory structure under `/books/` is **user-defined**.
* Implementations **MUST NOT** derive meaning from folder names.
* Book identity is independent of file path or filename.
* Not all books need to be stored in `/books/`; adapters may access book files from any location on the device.

This allows users to organize their library freely.

---

## 3. Book Identity

Each book is identified by the hash of its file contents.

### Rules

* Book ID **MUST** be computed as:
  ```
  sha256(file_bytes)
  ```

* The Book ID **MUST NOT** depend on:
  * filename
  * file path
  * metadata

* Book IDs **MUST** be treated as immutable

---

## 4. Reading State Files

* Each book may have **at most one** reading state file.
* Each reading state file **MUST** correspond to a single book.
* The filename **MUST** be the book's unique Book ID with a `.json` extension.
* The Book ID **MUST** be encoded in a filesystem-safe form.

**Example:**

If a book has Book ID:
```
abc123
```

then its reading state file **MUST** be:
```
/OpenShelf/.state/abc123.json
```

---

## 5. Reading State Schema

A reading state file **MUST** be a JSON object containing at least the required fields listed below. The following is a representative example of a fully populated state file:

```json
{
  "specVersion": "0.2.0-draft",
  "bookId": "abc123",
  "updatedAt": "2026-01-25T10:15:00Z",
  "location": {
    "percentage": {
      "value": 42.3,
      "updatedAt": "2026-01-25T10:15:00Z"
    },
    "epubcfi": {
      "value": "/6/2[chapter1]!/4/2/14",
      "updatedAt": "2026-01-25T10:15:00Z"
    },
    "page": {
      "value": 123,
      "updatedAt": "2026-01-25T10:15:00Z"
    }
  }
}
```

### Required top-level fields

* `specVersion` — The version of the OpenShelf spec this state file conforms to.
* `bookId` — The SHA-256 hash of the book’s contents (Book ID).
* `updatedAt` — ISO 8601 UTC timestamp representing the last update to the reading state file.
* `location` — Object containing all known location representations for the book (may be empty if no locations are recorded).

### Location fields

* Each location (`percentage`, `epubcfi`, `page`, etc.) **MUST** be an object containing:

  * `value` — The location value (number for percentage/page, string for epubcfi).
  * `updatedAt` — ISO 8601 UTC timestamp representing the last update to this specific location.

### Optional fields

* Any location representation may be omitted if the adapter does not support it.
* Implementations **MUST ignore unknown fields**.

### Notes

* Timestamps **MUST** be ISO 8601 UTC strings (e.g., `"2026-01-25T10:15:00Z"`).
* The top-level `updatedAt` **MUST** be updated whenever any location in the file changes.

---

## 6. Location Model

The `location` object represents the reader’s last known position in the book, using zero or more location representations.

Each field within the `location` object represents a distinct location representation. Implementations MAY provide multiple representations simultaneously.

* A **location representation** is a specific way of expressing a position in the book, such as:
  * `percentage` – the percentage of the book completed (0–100)
  * `page` – a linear page number within the book
  * `epubcfi` – a precise EPUB CFI (Canonical Fragment Identifier) path
* Zero or more location representations MAY be present.
* Implementations MUST handle missing fields gracefully.
* Percentage is RECOMMENDED as a universal fallback if no other location is available.

The behavioral interpretation, merging, and prioritization of location representations is defined in Section 10.

---

## 7. Capability Declaration

Declared capabilities describe the features supported by the library format itself, not by any specific reader or adapter.

The library **MUST** declare supported capabilities in `library.json`:
```json
{
  "spec": {
    "name": "OpenShelf",
    "version": "0.2.0-draft"
  },
  "capabilities": [
    "location"
  ]
}
```

### Required fields

* `spec`
* `spec.name`
* `spec.version`
* `capabilities`

---

## 8. Conflict Resolution

When multiple reading state updates conflict:

* Implementations **MUST** select the state with the latest `updatedAt` timestamp

---

## 9. Explicit Non-Goals

This version does **not** define:

* Annotations
* Highlights
* Bookmarks
* Metadata indexing
* Sync protocols
* Server APIs
* Multi-user libraries

---

## 10. Neutral Reading State Model

This section defines the **conceptual, in-memory model** that implementations SHOULD use when reading, merging, and writing reading state files. It exists to ensure consistent behavior across adapters and readers, even when different location representations are used.

This model is **not a wire format** and **not persisted directly**. It describes how implementations reason about reading state internally.

---

### 10.1 Conceptual Model

A reading state represents a **set of location representations** describing the reader’s last known position in a book.

Conceptually, a reading state consists of:

* A book identity
* A timestamp indicating when the state was last updated
* One or more location representations

In memory, implementations SHOULD treat locations as **independent entries**, even if they are serialized together in a single `location` object on disk.

Example conceptual structure:

```
ReadingState
├─ bookId
├─ updatedAt
├─ locations[]
│   ├─ type (percentage | epubcfi | page | unknown)
│   ├─ value
│   └─ metadata (optional, opaque)
```

---

### 10.2 Location Representations

Each location representation:

* Has a **type** (e.g., `percentage`, `epubcfi`, `page`)
* Has a **value** whose interpretation depends on the type
* Is considered a **hint**, not an absolute truth

No location representation is authoritative.

Implementations MUST assume that:

* Representations may be missing
* Representations may be inconsistent with each other
* Precision varies across readers and formats

---

### 10.3 Normalization (Disk → Memory)

When reading a reading state file:

* Each present field inside the `location` object MUST be mapped to a distinct in-memory location representation
* Missing fields MUST simply result in absent representations
* Unknown fields MUST be preserved as opaque data and carried forward during write-back
* `updatedAt` MUST be treated as the timestamp of the overall state, not of individual representations

No inference or conversion between representations is required during normalization.

---

### 10.4 Merge Semantics

When two reading states for the same book conflict (e.g., due to concurrent updates on different devices), implementations MUST apply **last-write-wins** semantics based on `updatedAt`.

After selecting the winning state:

* Location representations present in the newer state SHOULD replace representations of the same type from the older state
* Location representations present only in the older state MAY be preserved if not explicitly overwritten
* Unknown fields MUST NOT be discarded unless explicitly replaced

This allows different readers to contribute different location representations over time without data loss.

---

### 10.5 Write-Back Rules (Memory → Disk)

When writing a reading state file:

* Implementations MUST preserve unknown fields they do not understand
* Implementations SHOULD write all location representations they can compute
* Implementations MUST NOT remove location representations solely because they cannot interpret them
* `updatedAt` MUST be updated only when a meaningful reading position change occurs

Implementations SHOULD avoid unnecessary rewrites to minimize sync churn.

---

### 10.6 Precision and Fallback Behavior

Due to differences in rendering engines, screen sizes, and pagination strategies:

* EPUB CFI provides high precision but low portability
* Page numbers are reader-specific and least portable
* Percentage is the most portable and is RECOMMENDED as a universal fallback

Implementations SHOULD prefer the most precise supported representation when resuming reading, but MUST gracefully fall back to other available representations.

---

### 10.7 Lossy Operations

Any conversion between location representations MAY be lossy.

Implementations:

* MUST NOT assume reversible conversions
* SHOULD prefer preserving original representations over recomputation
* MUST treat all representations as advisory

---

### 10.8 Forward Compatibility

Future specification versions may introduce additional location representations or metadata.

To ensure compatibility:

* Unknown location types MUST be preserved
* Implementations MUST ignore fields they do not recognize
* No assumptions MUST be made about completeness or consistency of location data

---

## 11. Adapters

An **OpenShelf adapter** is a component that integrates a reader, device, or application
with an OpenShelf library.

Adapters are responsible for translating between:

* The OpenShelf filesystem and data formats
* Reader- or device-specific storage formats and capabilities

Adapters **MUST**:

* Treat the OpenShelf library as the source of truth
* Preserve all data they do not explicitly understand
* Respect the immutability of book identity
* Comply with the conflict resolution rules defined in this specification
* Tolerate the presence of other adapters interacting with the same library

Adapters **MUST NOT**:

* Modify or delete OpenShelf files outside their declared responsibility level
* Require exclusive access to the library
* Assume the presence of a server or network connectivity
* Fail or block reading due to unsupported or missing OpenShelf features

Adapters that integrate readers or devices with an OpenShelf library MUST comply
with the **OpenShelf Adapter Responsibility Levels** specification.

That specification defines adapter behavior, supported responsibility levels, and
safety constraints, and is a **normative extension** of this document.

