# OpenShelf Adapter Responsibility Levels

## Version: 0.1.0-draft
## Status: Normative Extension

Adapters bridge a OpenShelf Library filesystem with a specific reader or device. Different adapters may support different levels of interaction depending on platform constraints, reader APIs, and implementation maturity.

To ensure interoperability and predictable behavior, adapters are classified into **Responsibility Levels**.

Adapters MUST declare their supported level(s).
An adapter MAY support multiple responsibility levels simultaneously and select
the appropriate behavior based on context or configuration.

---

## Level 0 — Filesystem Access Only

### Description

A Level 0 adapter provides **access to book files only**. It does not read or write reading state.

### Responsibilities

* Expose `/books/` contents to the reader
* MUST NOT read from `/.state/`
* MUST NOT write to `/.state/`
* MUST NOT modify `library.json`

### Use cases

* Simple file browsers
* Legacy readers
* Read-only environments
* “Bring your own sync” workflows

### Guarantees

* Zero risk of state corruption
* Full compatibility with all libraries

---

## Level 1 — Read-Only Reading State

### Description

A Level 1 adapter can **consume** reading state to resume reading, but does not write updates back.

### Responsibilities

* Read `/.state/<bookId>.json`
* Normalize reading state according to Section 10
* Select an appropriate location representation
* Resume reading when possible

### Constraints

* MUST NOT write to `/.state/`
* MUST ignore unsupported location representations
* MUST gracefully fall back if state is missing or invalid

### Use cases

* Readers without write access
* Shared or read-only devices
* Testing and validation tools

---

## Level 2 — Write-Back Reading Progress

### Description

A Level 2 adapter can both **read and write** reading state, but only for basic progress.

### Responsibilities

* All Level 1 responsibilities
* Write reading state updates to `/.state/`
* Update `updatedAt` on meaningful progress change
* Preserve unknown fields when writing
* Follow merge semantics defined in Section 10

A meaningful progress change is implementation-defined, but SHOULD avoid
frequent updates caused by pagination, layout changes, or minor position jitter.

### Constraints

* MUST NOT delete existing location representations
* MUST NOT assume ownership of the state file
* MUST tolerate concurrent updates from other devices

### Use cases

* Most desktop and mobile readers
* KOReader-style integrations
* E-ink devices with writable storage

---

## Level 3 — Extended State (Future)

⚠️ *Not defined in v0.1*

Reserved for future specification versions that introduce:

* bookmarks
* annotations
* highlights

Adapters at this level will be required to:

* preserve unknown extended fields
* respect feature-specific merge rules

---

## Level 4 — Full Fidelity Adapter (Future)

⚠️ *Not defined in v0.1*

A Level 4 adapter provides:

* full bidirectional sync
* lossless annotation handling
* reader-native feature parity

This level may require:

* reader-specific extensions
* optional server assistance

---

## Adapter Capability Declaration

Adapters MUST declare their responsibility level(s) and supported location representations.

Example (conceptual):

```json
{
  "adapter": "koreader",
  "levels": [1, 2],
  "supportedLocations": [
    "percentage",
    "epubcfi"
  ]
}
```

This declaration may be:

* documented
* logged
* exposed via UI
* embedded in adapter metadata

The exact mechanism is implementation-defined.

---

## Adapter Safety Rules (Applies to All Levels)

Regardless of level, all adapters:

* MUST treat the filesystem as the source of truth
* MUST NOT rewrite files unnecessarily
* MUST preserve unknown fields
* MUST fail gracefully
* MUST NOT block reading due to missing or invalid state
* MUST NOT require exclusive access to the library

---

## Why Responsibility Levels Matter

This system allows:

* Incremental adoption
* Partial implementations
* Safe experimentation
* Clear user expectations
* Reader diversity

A minimal adapter can exist in a few hundred lines of code, while advanced adapters can evolve over time without breaking compatibility.
