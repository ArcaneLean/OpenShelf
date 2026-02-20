# OpenShelf Reading State Format Specification

## Version: 0.2.0-draft

This normative specification defines the **structure, naming, and schema** of reading state files.

A reading state file represents the last known reading position for a single book.

---

## 1. File Location and Naming

Reading state files **MUST** be stored in:

```
/OpenShelf/.state/
```

The filename **MUST**:

* Be the book’s unique Book ID (See Core Specification section 4).
* Use the `.json` extension.

Each book **MAY have at most one** reading state file.

Each reading state file **MUST correspond to exactly one book**.

### Example

If a book’s Book ID (in lowercase hexadecimal) is:

```
abc123
```

The reading state file **MUST** be:

```
/OpenShelf/.state/abc123.json
```

---

## 2. File Format

A reading state file **MUST** be a UTF-8 encoded JSON object.

The object **MUST** contain the required top-level fields defined below.

---

## 3. Required Top-Level Fields

| Field         | Type   | Description                                                  |
| ------------- | ------ | ------------------------------------------------------------ |
| `specVersion` | string | The OpenShelf specification version this file conforms to    |
| `bookId`      | string | The SHA-256 hash of the book’s contents                      |
| `updatedAt`   | string | ISO 8601 UTC timestamp of the last modification to this file |
| `location`    | object | Container for zero or more location representations          |

### Field Requirements

* `specVersion` **MUST** match the specification version understood by the writing implementation.
* `bookId` **MUST** equal the Book ID derived from the book’s contents.
* `updatedAt` **MUST** be an ISO 8601 UTC timestamp (e.g., `"2026-01-25T10:15:00Z"`).
* `location` **MUST** be present, but **MAY** be empty.

The top-level `updatedAt` **MUST** be updated whenever any location representation changes.

---

## 4. Location Representations

The `location` object contains zero or more **location representations**.

Each entry inside `location`:

* Uses the representation type as its key (e.g., `percentage`, `epubcfi`, `page`)
* **MUST** be an object with the following fields:

| Field       | Type             | Description                                             |
| ----------- | ---------------- | ------------------------------------------------------- |
| `value`     | number or string | The location value                                      |
| `updatedAt` | string           | ISO 8601 UTC timestamp for this specific representation |

### Example

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

---

## 5. Optional and Unknown Fields

* Location representations **MAY** be omitted.
* Additional top-level or nested fields **MAY** be present.
* Implementations **MUST ignore unknown fields**.
* Implementations **MUST preserve unknown fields** during write-back.

---

## 6. Timestamp Requirements

All timestamps:

* **MUST** be ISO 8601 formatted
* **MUST** be in UTC
* **MUST** include the `Z` suffix

---

# 7. Interoperable Location Representations

This section defines location representation types that are intended to promote interoperability across adapters and reading systems.

Location representations not defined in this section remain valid under the open vocabulary rules of this specification. Implementations **MUST** preserve unknown representation types.

The representations defined here are considered **interoperable representations**.

Future versions of this specification **MAY** define additional interoperable representations.

---

## 7.1 General Requirements

Interoperable representations:

* **MUST** have deterministic interpretation.
* **MUST** be platform-independent.
* **SHOULD** represent logical reading progression rather than layout-dependent position.
* **MUST** be safe to ignore if unsupported.
* **MUST NOT** invalidate a reading state file if malformed or unusable.

If multiple interoperable representations are present, they are assumed to represent the same logical position unless otherwise indicated.

---

## 7.2 Universal Interoperable Representation

### `percentage`

#### Type

`number`

#### Range

0.0 to 1.0 inclusive.

Values outside this range **MUST** be treated as invalid.

#### Meaning

Represents logical progression through the publication’s primary reading order.

* `0.0` indicates the start of the reading order.
* `1.0` indicates the end of the reading order.
* Intermediate values indicate proportional progression.

This representation:

* **MUST NOT** depend on viewport size, font size, pagination mode, or screen dimensions.
* **SHOULD** reflect position relative to the total logical content length.

#### Precision

Writers **SHOULD** limit precision to a reasonable number of decimal places (e.g., 4–6).
Readers **SHOULD** tolerate reasonable floating-point variance.

#### Interpretation

A reading system **SHOULD** seek to the closest internally representable position corresponding to the given percentage.
If exact mapping is not possible, the system **SHOULD** approximate to the nearest valid position.

**Notes**:

* `percentage` is the **universal fallback representation** for all content types.
* All adapters **SHOULD** preserve this representation when writing reading state.

---

## 7.3 Media-Specific Interoperable Representations

Some interoperable representations are only meaningful for certain media types. Implementations **SHOULD** use representations appropriate to the publication’s media type.

---

### 7.3.1 `epubcfi`

#### Type

`string`

#### Format

MUST conform to the EPUB Canonical Fragment Identifier (CFI) specification (EPUB 3.3).

Invalid CFIs **MUST** be ignored.

#### Meaning

Represents a canonical structural position within an EPUB publication’s spine.

#### Interpretation

* If the publication is not an EPUB or does not support CFI, this representation **MUST** be ignored.
* If resolution fails, readers **SHOULD** fall back to another representation (e.g., `percentage`).

---

### 7.3.2 `pageNumber`

#### Type

`integer`

#### Range

1 to N (inclusive), where N is the total number of logical pages.

#### Meaning

Represents the fixed-layout page index of a publication.

#### Applicability

* Interoperable only for fixed-layout publications (e.g., PDF, fixed-layout EPUB, CBZ/CBR comics).
* **MUST NOT** be treated as interoperable for reflowable documents.

---

### 7.3.3 `timeSeconds`

#### Type

`number`

#### Range

0.0 to total duration (inclusive)

#### Meaning

Represents playback position in seconds from start.

* Independent of playback speed.
* Percentage can exist as fallback.

#### Applicability

* Interoperable only for time-based media (audiobooks, podcasts).

---

## 7.4 Coexistence With Other Representations

Reading state files **MAY** include additional proprietary or reader-specific location representations.

Implementations:

* **MUST** preserve all representations when updating state.
* **MUST NOT** delete representations they do not understand.
* **MAY** prefer interoperable representations when writing new state.

The absence of interoperable representations does not invalidate a reading state file, but may reduce cross-adapter compatibility.

---

## 7.5 Extensibility

Additional interoperable representations **MAY** be defined in future versions.

Such additions **MUST**:

* Define type and range precisely.
* Provide deterministic semantic interpretation.
* Be layout-independent unless explicitly stated otherwise.
* Preserve forward compatibility with existing readers.