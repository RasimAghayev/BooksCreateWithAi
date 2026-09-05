import sys
import os
import json

import pymupdf

# Usage: python convert_raw.py <pdf_path> <workspace_dir>
# Creates workspace/raw/NNN.txt for every page. Page text extracted from
# embedded text layer (no OCR fallback in this script).

pdf_path = sys.argv[1]
workspace = sys.argv[2]

raw_dir = os.path.join(workspace, "raw")
os.makedirs(raw_dir, exist_ok=True)

doc = pymupdf.open(pdf_path)
total = len(doc)
stats = {"pages": total, "text_pages": 0, "empty_pages": 0}

for i in range(total):
    page = doc[i]
    text = page.get_text("text")
    if text and text.strip():
        stats["text_pages"] += 1
    else:
        stats["empty_pages"] += 1
    with open(os.path.join(raw_dir, f"{i+1:03d}.txt"), "w", encoding="utf-8") as f:
        f.write(text or "")

with open(os.path.join(raw_dir, "_stats.json"), "w", encoding="utf-8") as f:
    json.dump(stats, f, ensure_ascii=False, indent=2)

print(json.dumps(stats))
