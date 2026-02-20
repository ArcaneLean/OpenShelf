# OpenShelf Core Specification

## Version: 0.2.0-draft

---

## 1. Overview

OpenShelf is a filesystem-based specification for interoperable reading state storage.

It defines:

* A standard library directory structure
* A content-addressed book identity model
* A location for reading state files
* Conflict resolution rules
* A capability declaration mechanism

OpenShelf does not define reader behavior, synchronization protocols, or server infrastructure.

The OpenShelf Core Specification defines only the minimal interoperable requirements. Additional behavior is defined in separate normative extensions.

---

## 2. Terminology

---

## 3. Directory Structure

A compliant library follows this structure:

```
/OpenShelf/
  library.json
  /books/
  /.state/
```

* Adapters **MAY** read book files from any location on the device. The `/books/` directory is optional, provided for implementations that want a central location for all books.
* Reading state files **MUST** be stored under `/.state/`. All adapters **MUST** read and write state files from this location to ensure interoperability.
* Directories whose names start with `.` are **implementation-owned**. Implementations **MUST NOT** require users to manually edit dot-prefixed directories.
* Implementations **MAY** create additional directories as needed, provided they do not conflict with the mandated structure.

---

# 4. Book Identifiers

## 4.1 Overview

Each book in an OpenShelf library is uniquely identified by a `bookId`.

The `bookId` serves as the canonical identifier for:

* Reading state files
* Cross-device synchronization
* Adapter operations
* Library consistency checks

The `bookId` namespace is global within a library.

---

## 4.2 Normative Definition

The derivation and validation of `bookId` values are defined by the:

> **OpenShelf Book Identity Specification**

Implementations **MUST** derive `bookId` values strictly according to that specification.

Implementations **MUST NOT**:

* Generate random identifiers
* Use file paths as identifiers
* Use adapter-defined identifiers
* Derive identifiers using non-standard hashing methods

Only identifiers produced according to the Book Identity Specification are valid.

---

## 4.3 Identity Stability

A book’s `bookId` **MUST** remain stable for the lifetime of that book within the library.

Implementations:

* **MUST NOT** regenerate a `bookId` unless the underlying identity truly changes.
* **MUST** treat identical `bookId` values as referring to the same book.
* **MUST NOT** assign multiple `bookId` values to the same canonical identity.

If the canonical identity of a publication changes (e.g., different edition, modified content, or distinct logical source), a new `bookId` **MUST** be generated according to the Book Identity Specification.

---

## 4.4 Reading State Association

Reading state files located in `/.state/` **MUST** be associated with books using their `bookId`.

The filename of a reading state file **MUST** correspond exactly to the associated `bookId`.

Example:

```
.state/
  9f2c1a...e4.json
```

Where `9f2c1a...e4` is a valid `bookId` derived according to the Book Identity Specification.

---

## 4.5 Conformance

An implementation is not compliant with the OpenShelf Core Specification if it derives, assigns, or interprets `bookId` values in a manner inconsistent with the OpenShelf Book Identity Specification.

---

## 5. Reading State File Format

* A reading state file **MUST** conform to the **OpenShelf Reading State Format Specification**.
* The version of the format used **MUST** be compatible with the `spec.version` declared in `library.json`.

---

## 6. Conflict Resolution

When multiple reading state updates conflict:

* Implementations **MUST** select the state with the latest `updatedAt` timestamp.
* Timestamps **MUST** be compared as UTC.
* If timestamps are equal, behavior is implementation-defined.

---

## 7. Capability Declaration

Declared capabilities describe the features supported by the library format itself, not by any specific reader or adapter.

* Implementations **MUST** ignore unknown capabilities
* The library **MUST** declare supported capabilities in `library.json`:

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

## 8. Forward Compatibility

* Implementations **MUST** ignore unknown fields in all OpenShelf-defined files.
* Future versions of this specification **MAY** introduce additional fields or capabilities.

---

## 9. Non-Goals

This version does **not** define:

* Annotations
* Highlights
* Bookmarks
* Metadata indexing
* Sync protocols
* Server APIs
* Multi-user libraries

---

## 10. Related Specifications

The following specifications extend OpenShelf:

* **OpenShelf Reading State Format Specification** — defines the structure and naming of reading state files.
* **OpenShelf Neutral Reading State Model Specification** — defines the conceptual in-memory model and merge semantics.
* **OpenShelf Adapter Responsibility Levels Specification** — defines behavioral requirements for adapters interacting with an OpenShelf library.
* **OpenShelf Book Identity Specification** — defines the derivation, validation, and stability requirements of `bookId` values.

These documents define normative extensions to the OpenShelf Core Specification.

Implementations claiming conformance to OpenShelf **MUST** comply with all referenced specifications applicable to the features they implement.