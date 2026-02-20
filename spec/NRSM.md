
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
