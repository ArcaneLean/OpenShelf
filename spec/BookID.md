# 4. Book Identity and Canonical Identity Descriptor

## 4.1 Book Identifier

Each book in an OpenShelf library is identified by a `bookId`.

The `bookId` **MUST** be the SHA-256 hash (hex-encoded, lowercase) of the book’s **Canonical Identity Descriptor**.

This ensures:

* Deterministic identity
* Collision resistance
* Uniform identifier format
* Media-agnostic compatibility

---

## 4.2 Canonical Identity Descriptor

A Canonical Identity Descriptor (CID) is the deterministic byte sequence from which the `bookId` is derived.

The CID **MUST** be defined according to the book’s identity type.

---

## 4.3 Identity Types

### 4.3.1 File-Based Publications

For file-based publications (e.g., EPUB, PDF, CBZ):

The Canonical Identity Descriptor **MUST** be the raw, unmodified byte contents of the file.

No transformation, compression normalization, or metadata stripping is permitted.

```
CID = raw file bytes
bookId = sha256(CID)
```

This guarantees content-derived identity.

---

### 4.3.2 Non-File-Based Publications

For publications that are not represented by a local file (e.g., online books, webtoons, audiobooks streamed via URL, manually tracked content), the Canonical Identity Descriptor **MUST** be a canonicalized JSON object encoded as UTF-8.

This JSON object is referred to as the **Identity Descriptor Object (IDO)**.

---

## 4.4 Identity Descriptor Object (IDO)

The IDO **MUST**:

* Be a valid JSON object
* Contain only deterministic fields
* Exclude transient or user-specific metadata
* Be serialized using canonical JSON rules defined below

### Required Fields

The IDO **MUST** include:

```
{
  "type": "<identityType>",
  "source": "<primary source identifier>"
}
```

Where:

* `type` identifies the identity category (e.g., `"web"`, `"stream"`, `"manual"`).
* `source` uniquely identifies the publication within that category.

### Optional Fields

Optional fields MAY be included if they contribute to stable identity, such as:

* `provider`
* `edition`
* `language`

Optional fields **MUST NOT** include:

* Access tokens
* Session identifiers
* Timestamps
* User-specific data

---

## 4.5 Canonical JSON Serialization

When deriving the Canonical Identity Descriptor from an IDO:

1. The JSON object **MUST** be serialized using:

   * UTF-8 encoding
   * No insignificant whitespace
   * Lexicographically sorted object keys
   * No trailing zeros in numbers
2. The serialized byte sequence becomes the CID.

```
CID = UTF8(canonicalJSON(IDO))
bookId = sha256(CID)
```

This ensures deterministic identity across implementations.

---

## 4.6 URL Normalization Requirements

If the IDO includes a URL in `source`, implementations **SHOULD** apply canonicalization rules before inclusion:

* Remove URL fragments (`#...`)
* Remove known tracking query parameters (e.g., `utm_*`)
* Normalize scheme and host to lowercase
* Remove default ports
* Resolve redundant path segments

If canonicalization cannot be performed reliably, the URL **MAY** be used as-is, but identity stability is not guaranteed.

---

## 4.7 Identity Stability Rules

Implementations:

* **MUST NOT** change a book’s identity descriptor once created.
* **MUST NOT** regenerate a bookId unless the underlying identity truly changes.
* **MUST** treat two identical CIDs as representing the same book.

Different identity types (e.g., a live web novel vs. a downloaded EPUB) are considered distinct publications and MUST produce distinct `bookId` values.