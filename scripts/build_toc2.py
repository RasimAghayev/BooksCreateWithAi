import os
import re

raw_dir = r"D:\Tasks\ai_code\books_ai\Books\Explore_Go_Cryptography - John Arundel - 2026 - v1\raw"

files = sorted(os.listdir(raw_dir))

toc_titles = [
    "Ciphers",
    "Enciphering",
    "Deciphering",
    "Cracking",
    "Keys",
    "Cribs",
    "Passwords",
    "Blocks",
    "Modes",
    "Padding",
    "Enumeration",
    "Entropy",
    "Randomness",
    "Chains",
    "Hashing",
    "Coins",
    "Authentication",
    "Cryptography"
]

chapters = []

for i, title in enumerate(toc_titles, 1):
    start_page = None
    end_page = None
    
    for filename in files:
        filepath = os.path.join(raw_dir, filename)
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()
        
        if title in content and start_page is None:
            start_page = int(filename.replace(".txt", ""))
    
    chapters.append({
        "chapter": i,
        "title": title,
        "start_page": start_page
    })

for i, ch in enumerate(chapters):
    if i < len(chapters) - 1:
        ch["end_page"] = chapters[i + 1]["start_page"] - 1
    else:
        ch["end_page"] = int(files[-1].replace(".txt", ""))

for ch in chapters:
    print(f"Chapter {ch['chapter']}: {ch['title']} (pages {ch['start_page']}-{ch['end_page']})")
