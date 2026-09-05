import json
import os
import shutil
import sys
import re

# Build toc.json + copy chapter source pages from raw/ into source_pages/chapter-NN/

workspace = sys.argv[1]
meta = json.load(open(os.path.join(workspace, "metadata.json"), encoding="utf-8"))
raw_dir = os.path.join(workspace, "raw")

chapters = []
for ch in meta["📑 chapters"]:
    start, end = [int(x) for x in ch["pdf_pages"].split("-")]
    chapters.append({
        "chapter": ch["chapter"],
        "title": ch["title"],
        "book_pages": ch["pages"],
        "start_pdf_page": start,
        "end_pdf_page": end,
    })

os.makedirs(os.path.join(workspace, "toc"), exist_ok=True)
with open(os.path.join(workspace, "toc", "toc.json"), "w", encoding="utf-8") as f:
    json.dump({"book": meta["📘 title_original"], "page_offset": meta["page_offset"], "chapters": chapters}, f, ensure_ascii=False, indent=2)

sp_root = os.path.join(workspace, "source_pages")
for ch in chapters:
    cdir = os.path.join(sp_root, f"chapter-{ch['chapter']:02d}")
    os.makedirs(cdir, exist_ok=True)
    for p in range(ch["start_pdf_page"], ch["end_pdf_page"] + 1):
        src = os.path.join(raw_dir, f"{p:03d}.txt")
        if os.path.exists(src):
            shutil.copy2(src, os.path.join(cdir, f"page-{p:03d}.txt"))

# chapter folders
def slug(t):
    s = re.sub(r"[^a-z0-9]+", "-", t.lower()).strip("-")
    return s[:40]

ch_root = os.path.join(workspace, "chapters")
for ch in chapters:
    os.makedirs(os.path.join(ch_root, f"{ch['chapter']:02d}-{slug(ch['title'])}"), exist_ok=True)

print("toc + source_pages + chapter folders done")
