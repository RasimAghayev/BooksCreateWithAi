# Continue Book Processing — Plan

## Current State
- **Book 1** (`_media_books_go_razrabotka_prilozhenij_v_mikroservisnoj_arhitekture`): workspace deleted from working tree; raw PDF still in `Books/`; `telegram_upload/` has the ready package.
- **Book 2** (`Explore_Go_Cryptography - John Arundel - 2026 - v1`): PDF present in `Books/`, no workspace created yet.
- `reader/current_book.json` and `reader/books.json` are empty.
- Existing plan file: `.kilo/plans/1788554893864-book-reader-agent-plan.md`.

## Goal
Continue the `next` pipeline for Book 2 (`Explore_Go_Cryptography`) and prepare it for Telegram upload.

## Steps

### 1. State Recovery / Reset
- Restore `reader/current_book.json` to point to Book 2 workspace.
- Initialize `reader/books.json` entry for Book 2 with `status: reading`.
- Verify `.gitignore` is correct (already excludes `source/`, `raw/`, `external/source-code/`).

### 2. Book 2 Workspace Creation
- Create workspace: `Books/Explore_Go_Cryptography - John Arundel - 2026 - v1/`
- Move/copy PDF into `source/original.pdf`.
- Run `scripts/extract_pdf.py` to create `raw/001.txt`, `raw/002.txt`, ...
- Run `scripts/build_toc.py` (or `build_toc3.py`) to extract TOC.
- Create `metadata.json` with:
  - `title_original` (preserved)
  - `filename_safe` (sanitized)
  - `title`, `author`, `year`, `version`, `page_count`
  - `classification` (primary_language: Go, level, domains, technologies, tags)
  - `chapters` array

### 3. Chapter Processing
- Create `chapters/` folders per TOC (01-introduction, 02-... through 18).
- For each chapter:
  - Create `index.md` with chapter title, page range, sections.
  - Extract Azerbaijani knowledge summary with code block explanations.
  - Create `cheatsheet.md` for CLI commands, code snippets.
- Commit after each chapter: `git add` → `git commit -m "Book 2: Chapter N processed"` → `git push`.

### 4. Info Files
- `info/summary.az.md` — overall book summary in Azerbaijani.
- `info/terminology.az.md` — extracted terms with `English (Azərbaycanca qarşılıq)` format.
- `info/teacher-notes.az.md` — teacher notes with book quotes vs agent recommendations.
- `info/cheatsheet.az.md` — global cheat sheet from all chapters.

### 5. Telegram Package
- Run the packaging step to create `telegram_upload/Explore_Go_Cryptography - John Arundel - 2026 - v1.pdf` + `.json`.
- Ensure metadata JSON uses icons in **values** (not keys), per SYSTEM_PROMPT.md invariant.
- Verify `telegram/channels/RA_ALL/.env`, `RA_GO/.env`, `RA_JS/.env` have real credentials before upload.

### 6. Upload & Cleanup
- Run `uploader.exe` from `telegram/uploader/`.
- Verify per-channel upload status in store.
- Delete local files only after ALL matched channels succeed.

## Validation
- Check `raw/` files exist and are non-empty.
- Check `metadata.json` is valid JSON with required fields.
- Check at least 2 chapter `index.md` files exist with Azerbaijani content.
- Check `telegram_upload/` contains `.pdf` + `.json` pair.
- Verify git status is clean (only expected files staged/committed).

## Blockers / Dependencies
- Telegram credentials must be filled in `telegram/channels/RA_*/.env` for upload step.
- Book 1 workspace recovery is out of scope (already packaged for Telegram).
