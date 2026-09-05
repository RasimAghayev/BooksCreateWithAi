import sys
import os
import re

# Usage: python clean_raw.py <raw_dir> <pattern_file_regex>
# Removes license/watermark lines from every raw page file.

raw_dir = sys.argv[1]
pattern = re.compile(r"Licensed to .*$", re.IGNORECASE)

count = 0
for name in sorted(os.listdir(raw_dir)):
    if not name.endswith(".txt"):
        continue
    path = os.path.join(raw_dir, name)
    with open(path, encoding="utf-8") as f:
        lines = f.read().splitlines()
    kept = [ln for ln in lines if not pattern.search(ln)]
    if kept != lines:
        count += 1
        with open(path, "w", encoding="utf-8") as f:
            f.write("\n".join(kept) + ("\n" if kept else ""))

print(f"cleaned {count} files")
