## 1. Overview

* Filesystem-based specification defining interoperable reading state storage.
* Introduces a **content-addressed book identity model**.
* Core responsibilities: minimal interoperable requirements, conflict resolution, capability declaration.
* Does **not** define reader behavior, sync protocols, or server infrastructure.

---

## 2. Terminology

* Define key terms such as `bookId`, reading state, canonical metadata, mutable metadata, adapter, etc.

---

## 3. Book Identity

* Short overview of the **book identity model**.
* References the **Book Identity Specification** for derivation rules.
* Explains that a `bookId` uniquely identifies a book and ties reading state and metadata.
* Emphasizes stability: identical content → same `bookId`, changes to canonical identity → new `bookId`.

---

## 4. Directory Structure

* Minimal filesystem layout:

```
/OpenShelf/
  library.json
  /books/
  /.state/
  /.canonical/
  /.metadata/
```

* **/books/** — optional, central location for book files.

* **/.state/** — required, reading state files.

* **/.canonical/** — immutable canonical metadata defining identity.

* **/.metadata/** — mutable, non-identity-defining metadata (user tags, custom titles).

* Dot-prefixed directories are **implementation-owned**; users should not manually edit.

* Purpose of folders is briefly explained; full details and normative behavior are in separate specifications.

---

## 5. Reading State File Format

* References **OpenShelf Reading State Format Specification**.
* Core spec only gives a brief overview: each reading state file is tied to a `bookId` and stored under `/.state/`.

---

## 6. Conflict Resolution

* If multiple updates conflict: use **latest `updatedAt` timestamp**.
* Timestamps must be compared in UTC.
* Equal timestamps: behavior is implementation-defined.

---

## 7. Capability Declaration

* Libraries declare their supported capabilities in `library.json`.
* Unknown capabilities must be ignored.
* Example snippet:

```json
{
  "spec": {
    "name": "OpenShelf",
    "version": "0.2.0-draft"
  },
  "capabilities": ["location"]
}
```

---

## 8. Forward Compatibility

* Unknown fields **MUST** be ignored.
* Future versions may introduce additional fields or capabilities.

---

## 9. Non-Goals

* Core spec does **not** define:

  * Annotations, highlights, bookmarks
  * Metadata indexing beyond minimal structure
  * Sync protocols, server APIs
  * Multi-user libraries

---

## 10. Related Specifications

* **OpenShelf Reading State Format Specification** — reading state file structure.
* **OpenShelf Neutral Reading State Model Specification** — conceptual in-memory model, merge semantics.
* **OpenShelf Adapter Responsibility Levels Specification** — adapter behavior requirements.
* **OpenShelf Book Identity Specification** — derivation, validation, and stability of `bookId` values.
