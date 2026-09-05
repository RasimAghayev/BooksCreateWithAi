# Book Reader Agent — Redesign Plan

## 1. Current Architecture Summary

The existing system has:
- `SYSTEM_PROMPT.md` as the agent prompt with full pipeline rules
- `telegram/uploader/` Go binary for Telegram uploads with per-channel routing
- `reader/current_book.json` and `reader/books.json` for tracking
- `progress.json` per book workspace for resumability
- Telegram uploader with `telegram_upload/` staging area

### Critical Problems Identified
1. **Filename-as-identity**: Some code (uploader `store.go`) uses `book_file` (filename) as primary key instead of stable `book_id`
2. **Metadata icon placement**: `SYSTEM_PROMPT.md` puts emojis in JSON keys (`"📘 title_original"`) instead of values
3. **Missing state machine**: No explicit states like DISCOVERED, IDENTIFIED, UPLOADING, PARTIALLY_UPLOADED
4. **Rename breaks lookups**: After `RENAME` step, old filename-based lookups fail
5. **Partial upload tracking**: No per-channel upload status in `reader/` state — only Telegram uploader internal store
6. **No recovery from corrupted state**: If `progress.json` is corrupted, system has no safe fallback
7. **NEXT command ambiguity**: Currently described as "one step" but should mean "full pipeline for one book"

---

## 2. Target State Machine

```text
DISCOVERED          → Book found in Books/
IDENTIFIED          → Title/author/year extracted, book_id assigned
READING             → Raw pages extracted, chapters being processed
METADATA_READY      → metadata.json complete
CONTENT_READY       → All chapters, cheatsheets, info/ complete
READY_FOR_UPLOAD    → Git committed, telegram_upload/ package created
UPLOADING           → Upload in progress
PARTIALLY_UPLOADED  → Some channels succeeded, some failed
UPLOADED            → All matched channels succeeded
ARCHIVED            → Moved to books_read/
COMPLETED           → Archive + cleanup done, _book_complete.json exists
FAILED              → Unrecoverable error
```

---

## 3. Book Identity (Stable)

Every book gets a deterministic `book_id` at IDENTIFY stage:

```json
{
  "book_id": "sha256(title + author + year + version + page_count + sample_pages)",
  "title_original": "Go in Action",
  "title_normalized": "go-in-action",
  "author": "William Kennedy",
  "year": 2016,
  "version": 1,
  "page_count": 266,
  "original_filename": "original.pdf",
  "current_filename": "Go in Action - William Kennedy - 2016 - v1.pdf",
  "current_path": "Books/05 go-in-action/",
  "filename_safe": "05 go-in-action",
  "identity_fingerprint": "sha256:abc123..."
}
```

Rules:
- `book_id` NEVER changes, even if file is renamed or moved
- `original_filename` preserved forever in metadata
- `current_filename` and `current_path` updated atomically on rename
- All state files reference `book_id`, never raw filename

---

## 4. Reader Recovery Strategy

On `NEXT` or `ALL_START`:

```
1. Scan reader/current_book.json → is there an active book?
2. If yes → read book_id
3. Check workspace: Books/{filename_safe}/progress.json
4. If progress.json exists and valid → resume from last_step
5. If progress.json missing/corrupt → check _book_complete.json
6. If _book_complete.json exists → mark COMPLETED, skip to next
7. If no valid state → SAFE RECOVERY: restart from CREATE_WORKSPACE
   BUT preserve existing: metadata.json, chapters/, toc/, telegram_upload/
```

---

## 5. Rename Strategy

Rename happens ONCE after IDENTIFY:

```text
Before: Books/{random-original-name}/
After:  Books/{filename_safe}/
```

Atomic updates:
- `metadata.json`: `current_filename`, `current_path`, `filename_safe`
- `reader/current_book.json`: `current_path`
- `reader/books.json`: entry path updated

After rename, NO code uses the original folder name. All lookups use `book_id` or `filename_safe`.

---

## 6. Metadata Structure (Fixed)

Icons MUST be in values, NEVER in keys:

```json
{
  "book_id": "sha256:...",
  "title_original": "📘 Go in Action",
  "filename_safe": "05 go-in-action",
  "book": {
    "title": "📖 Go in Action",
    "author": "🧑‍💻 William Kennedy",
    "year": "📅 2016",
    "version": "📖 v1",
    "pages": "📄 266",
    "isbn": "🔢 978-...",
    "publisher": "🏢 Manning",
    "series": "📚 ..."
  },
  "classification": {
    "primary_language": "💻 Go",
    "level": "🎯 2",
    "level_name": "Elementary",
    "domains": ["Backend Development", "Software Development"]
  },
  "technologies": ["Go", "Golang", "Web Services"],
  "tags": ["#go", "#golang", "#backend"],
  "chapters": [
    {"chapter": 1, "title": "Introduction", "start_page": 1, "end_page": 20, "slug": "01-introduction"}
  ],
  "summary": "📝 ...",
  "telegram_uploads": [
    {
      "channel": "RA_ALL",
      "status": "success",
      "message_id": 1234,
      "file_id": "AgAC...",
      "uploaded_at": "2026-09-05T10:00:02Z"
    }
  ],
  "processing_state": "UPLOADED",
  "last_updated": "2026-09-05T10:00:02Z"
}
```

Rule: `KEY = semantic field name`, `VALUE = presentation-ready string (with icons if needed)`

---

## 7. Per-Channel Upload State

Each book tracks uploads per channel:

```json
{
  "uploads": {
    "RA_ALL": {
      "status": "success",
      "message_id": 1234,
      "file_id": "AgAC...",
      "uploaded_at": "..."
    },
    "RA_GO": {
      "status": "failed",
      "retry_count": 2,
      "last_error": "timeout",
      "last_attempt": "..."
    }
  }
}
```

Rules:
- UPLOAD step checks existing status BEFORE attempting
- If `status == success` → SKIP that channel
- If `status == failed` and `retry_count < 2` → RETRY
- Local file deleted ONLY when ALL matched channels = success

---

## 8. Observer & Logging

Structured, append-only logs:

**logs/process.log**:
```
timestamp | event | book_id | book_title | source_path | current_path | step | status | error | retry_count
```

**logs/upload.log**:
```
timestamp | book_id | book_title | channel | message_id | file_id | status | duration_ms
```

Minimum events:
```
BOOK_DISCOVERED, BOOK_IDENTIFIED, WORKSPACE_CREATED, READ_STARTED, READ_COMPLETED,
METADATA_CREATED, RENAME_STARTED, RENAME_COMPLETED, CONTENT_READY, UPLOAD_STARTED,
UPLOAD_SUCCESS, UPLOAD_FAILED, ARCHIVE_STARTED, ARCHIVE_COMPLETED, CLEANUP_STARTED,
CLEANUP_COMPLETED, BOOK_COMPLETED, BOOK_FAILED
```

---

## 9. Archive vs Cleanup

**ARCHIVE**: `Books/{filename_safe}/` → `books_read/{filename_safe}/`
- Copies entire workspace
- Happens after ALL channels succeed
- `books_read/` is append-only history

**CLEANUP**: Deletes `source/` binary from `Books/` AFTER archive validation
- `source/` deleted only
- `raw/`, `chapters/`, `metadata.json` etc. moved to `books_read/`

---

## 10. NEXT Command (Full Pipeline)

`NEXT` = full pipeline for ONE book, then stop:

```
NEXT
  ↓
detect active/incomplete book
  ↓
if incomplete → recover state
  ↓
select next unprocessed book from Books/
  ↓
CHECK_IDENTITY (deduplication)
  ↓
CREATE_WORKSPACE
  ↓
CONVERT_TO_RAW
  ↓
IDENTIFY_BOOK → assign book_id
  ↓
RENAME → atomic path update
  ↓
BUILD_METADATA + CLASSIFICATION
  ↓
FIND_TOC → BUILD_INDEX
  ↓
BUILD_CHAPTER_FOLDERS
  ↓
EXTRACT_KNOWLEDGE (chapter by chapter, AZ)
  ↓
BUILD_CHEATSHEETS
  ↓
BUILD_TEACHER_NOTES
  ↓
GIT_COMMIT_PUSH
  ↓
TELEGRAM_PACKAGE (book.json + pdf)
  ↓
TELEGRAM_ROUTE → matched channels
  ↓
UPLOAD (per channel, skip successes, retry failures)
  ↓
VERIFY all channels
  ↓
ARCHIVE → books_read/
  ↓
CLEANUP → delete source/
  ↓
mark COMPLETED (_book_complete.json)
  ↓
STOP
```

---

## 11. ALL_START Command

Loop `NEXT` until:
- `Books/` has no unprocessed books
- OR user interrupts

On interruption:
- Current book state preserved
- Next `ALL_START` resumes from current book

---

## 12. Idempotency Guarantees

| Operation | Idempotency Mechanism |
|-----------|----------------------|
| Read | Skip if `_book_complete.json` exists |
| Metadata | Skip if `metadata.json` exists and valid |
| Chapters | Skip if `chapters/` complete |
| Git | Skip if already committed (check git log) |
| Telegram upload | Skip if `telegram_uploads[channel].status == success` |
| Archive | Skip if `books_read/{filename_safe}/` exists and validated |
| Cleanup | Skip if `source/` already deleted |

---

## 13. Failure Matrix

| Failure | State After | Recovery Action |
|---------|-------------|-----------------|
| READ_FAILED | READING | Retry page extraction from `progress.json` page |
| METADATA_FAILED | IDENTIFIED | Retry IDENTIFY → BUILD_METADATA |
| RENAME_FAILED | IDENTIFIED | Retry RENAME (atomic) |
| GIT_FAILED | CONTENT_READY | Retry git commit (content preserved) |
| TELEGRAM_FILE_FAILED | UPLOADING | Retry failed channels only, keep local file |
| TELEGRAM_METADATA_FAILED | UPLOADING | Retry metadata upload |
| ONE_CHANNEL_FAILED | PARTIALLY_UPLOADED | Next run retries only failed channels |
| MULTIPLE_CHANNELS_FAILED | PARTIALLY_UPLOADED | Retry all failed channels |
| AGENT_CRASH | Varies | Read `reader/current_book.json` + workspace `progress.json` |
| PC_RESTART | Varies | Same as AGENT_CRASH |
| DUPLICATE_BOOK | IDENTIFIED | SKIP, log to `reader/books.json` |
| RENAMED_FILE | IDENTIFIED | Lookup by `book_id`, not filename |
| MISSING_STATE | Varies | SAFE RECOVERY: restart from CREATE_WORKSPACE |
| CORRUPTED_STATE | Varies | SAFE RECOVERY: preserve metadata/chapters if valid |
| ARCHIVE_FAILED | UPLOADED | Retry archive (upload already succeeded) |
| CLEANUP_FAILED | ARCHIVED | Retry cleanup |

---

## 14. Git Strategy

- Chapter-by-chapter commits with message: `feat: chapter N — {title}`
- `.gitignore` excludes:
  - `Books/**/source/`
  - `Books/**/raw/`
  - `Books/**/external/source-code/`
  - `books_read/`
  - `telegram_upload/`
  - `logs/`
- Git tracks: `metadata.json`, `toc/`, `chapters/**/*.md`, `info/*.md`, `progress.json`

---

## 15. Telegram Package

Created at `telegram_upload/{filename_safe}.pdf` + `telegram_upload/{filename_safe}.json`

Routing:
```
classification.primary_language == "Go" → RA_GO
classification.domains contains "Backend" → RA_BACKEND
match: "*" → RA_ALL
```

Caption built deterministically from metadata JSON (icons in values).

---

## 16. Implementation Tasks

### Task 1: Update SYSTEM_PROMPT.md Metadata Section
- Move all emojis from JSON keys to values
- Add `book_id`, `current_path`, `processing_state`, `telegram_uploads` to metadata schema

### Task 2: Add State Machine to SYSTEM_PROMPT.md
- Define all states in §8
- Add recovery logic to §8

### Task 3: Update NEXT Command Definition
- Change from "one step" to "full pipeline for one book"
- Add recovery as first step

### Task 4: Update ALL_START Definition
- Add recovery and idempotency rules

### Task 5: Add Archive vs Cleanup Distinction
- Add §15 with clear separation

### Task 6: Add Failure Matrix to SYSTEM_PROMPT.md
- Add §16 with table

### Task 7: Add Idempotency Section to SYSTEM_PROMPT.md
- Add §17 with guarantees table

### Task 8: Update Telegram Uploader store.go
- Change primary key from `book_file` to `book_id`
- Add per-channel status tracking
- Add retry_count field

### Task 9: Update Scanner to Handle Nested Metadata
- Already done (verified in scanner.go)

### Task 10: Add Observer Events to Uploader
- Emit structured log events for each state transition

---

## 17. Validation Steps

1. **Identity test**: Rename book folder → verify system still finds it by `book_id`
2. **Recovery test**: Kill process mid-chapter → verify resume from `progress.json`
3. **Partial upload test**: Simulate RA_GO failure → verify next run retries only RA_GO
4. **Duplicate test**: Drop same book with different filename → verify SKIP
5. **Icon test**: Verify JSON keys have no emojis, values have emojis
6. **Archive test**: Verify cleanup only after all channels succeed
7. **State test**: Verify `_book_complete.json` only created after all steps

---

## 18. Files to Modify

| File | Changes |
|------|---------|
| `SYSTEM_PROMPT.md` | Metadata icons, state machine, recovery, NEXT/ALL_START, archive/cleanup, failure matrix, idempotency |
| `.kilo/plans/1788554893864-book-reader-agent-plan.md` | This file |
| `telegram/uploader/internal/uploader/store.go` | book_id primary key, per-channel status, retry_count |
| `telegram/uploader/internal/uploader/service.go` | Use book_id for lookups, structured observer events |
| `telegram/uploader/internal/scanner/scanner.go` | Already updated for nested metadata |

---

## 19. Out of Scope (Defer)

- OCR engine implementation details
- EPUB/DJVU parsing libraries
- Telegram rate limit handling
- Conflict resolution for `books.json` entries
