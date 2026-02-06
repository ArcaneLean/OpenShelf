YES. This is the perfect first adapter. 🎯
KOReader is flexible, open, file-based, and already lives close to your philosophy.

Let’s design this **before** writing code so you don’t fight KOReader later.

I’ll walk you through this in **layers**, and we’ll stop at clear, implementable decisions.

---

# OpenShelf × KOReader Adapter — Design

## 0️⃣ Adapter Goals (explicit)

The KOReader adapter should:

* Be **Level 2** (read + write progress)
* Require **no KOReader patches** if possible
* Work **offline-first**
* Tolerate external sync (Dropbox, Syncthing, etc.)
* Never break KOReader’s native behavior
* Never assume exclusive ownership of state

This is a *good citizen* adapter.

---

## 1️⃣ What KOReader gives us (important constraints)

KOReader already has:

* A local **library directory**
* Per-book **document settings** (Lua tables)
* Support for:

  * EPUB CFI
  * Page-based progress
  * Percentage-like progress (via location / total)
* Optional KOReader sync (we will **not** use it)

KOReader state lives roughly in:

```
koreader/
  settings/
    document/
      <hashed-filename>.lua
```

These files already contain:

* last position
* current page
* zoom
* rendering prefs

We must **not replace** these — only map to/from them.

---

## 2️⃣ Adapter Shape (high level)

The adapter has **three responsibilities**:

```
OpenShelf filesystem
        ⇅
Neutral Reading State (in-memory)
        ⇅
KOReader document state
```

This is where your Section 10 model shines.

---

## 3️⃣ Where does the adapter run?

You have two viable architectures. Pick one deliberately.

### Option A — In-process KOReader plugin (recommended)

* Written in Lua
* Runs inside KOReader
* Reads OpenShelf files directly
* Updates KOReader state on open / close

**Pros**

* Best UX
* No external process
* No sync race guessing

**Cons**

* Lua
* KOReader API familiarity needed

### Option B — External sync helper

* Runs outside KOReader
* Translates states before/after reading

**Pros**

* Any language
* Easier debugging

**Cons**

* Worse UX
* More race conditions
* Harder to know “book opened” events

👉 **Strong recommendation: Option A**

KOReader plugins are powerful and this fits naturally.

---

## 4️⃣ Book Identity Mapping (critical design)

OpenShelf book ID = `sha256(file_bytes)`

KOReader identifies books by:

* file path
* internal hash of filename/path

### Adapter strategy

When a book is opened in KOReader:

1. Compute `sha256(file_bytes)`
   (cache this aggressively)
2. Look for:

   ```
   /.state/sha256_<hash>.json
   ```
3. If found → import reading state
4. If not → do nothing

📌 The adapter **never invents IDs**.

---

## 5️⃣ Mapping location models (the heart)

### KOReader → OpenShelf (write-back)

KOReader can provide:

| KOReader data    | OpenShelf field       |
| ---------------- | --------------------- |
| CFI              | `location.epubcfi`    |
| current page     | `location.page`       |
| position / total | `location.percentage` |

Write rules:

* Update `updatedAt` only on:

  * page turn
  * chapter change
  * reader close
* Preserve unknown fields in JSON
* Never delete other location types

---

### OpenShelf → KOReader (resume)

When opening a book:

Priority order (recommended):

1. `epubcfi` (if EPUB + supported)
2. `percentage`
3. `page`

Rules:

* Use **first supported** representation
* Never assume consistency
* If resume fails → fall back silently

---

## 6️⃣ Conflict handling (simple and safe)

You do **not** implement sync.

Instead:

* On book open:

  * Read OpenShelf state
  * If `updatedAt` > KOReader local timestamp → apply
* On progress change:

  * Write OpenShelf state
  * Let filesystem sync resolve conflicts

If two devices write:

* Spec’s LWW applies
* Adapter simply re-reads next time

No locks. No coordination.

---

## 7️⃣ Adapter Responsibility Declaration

Document this clearly (README or metadata):

```json
{
  "adapter": "openshelf-koreader",
  "levels": [1, 2],
  "supportedLocations": [
    "epubcfi",
    "percentage",
    "page"
  ]
}
```

This matches your spec perfectly.

---

## 8️⃣ Failure modes (explicitly allowed)

Your adapter **must** survive:

* Missing `/.state/`
* Invalid JSON
* Partial writes
* Sync conflicts
* Read-only filesystem
* Hash computation failures

Behavior:

> Log, ignore, continue reading.

Never block KOReader.

---

## 9️⃣ Minimal feature set (v0.1 adapter)

Do **only this** first:

✔ Read OpenShelf state on book open
✔ Resume reading if possible
✔ Write progress on close
✔ Preserve unknown fields

Do **not**:

* sync annotations
* touch KOReader settings
* optimize prematurely

---

## 10️⃣ Suggested file layout (repo)

```
openshelf-koreader/
  README.md
  plugin/
    openshelf.lua
    state.lua
    hash.lua
    json.lua
```

Small. Focused. Replaceable.

---

## Next concrete step

Pick **one** of these and we’ll go deeper:

1️⃣ KOReader plugin lifecycle (when hooks fire)
2️⃣ Exact Lua data structures to read/write
3️⃣ JSON read/write safety patterns
4️⃣ How to cache hashes efficiently
5️⃣ UX: should this be invisible or toggleable?

Say the number — we’ll design it properly.

