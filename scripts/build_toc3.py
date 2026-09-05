import os
import re

raw_dir = r"D:\Tasks\ai_code\books_ai\Books\Explore_Go_Cryptography - John Arundel - 2026 - v1\raw"

files = sorted(os.listdir(raw_dir))

chapter_pattern = re.compile(r'^(\d+)\.\s+(Ciphers|Enciphering|Deciphering|Cracking|Keys|Cribs|Passwords|Blocks|Modes|Padding|Enumeration|Entropy|Randomness|Chains|Hashing|Coins|Authentication|Cryptography)$')

chapters = []

for filename in files:
    page_num = int(filename.replace(".txt", ""))
    if page_num < 20:
        continue
    
    filepath = os.path.join(raw_dir, filename)
    with open(filepath, "r", encoding="utf-8") as f:
        lines = f.readlines()
    
    for line in lines:
        line = line.strip()
        match = chapter_pattern.match(line)
        if match:
            chapter_num = int(match.group(1))
            chapter_title = match.group(2)
            
            chapters.append({
                "chapter": chapter_num,
                "title": chapter_title,
                "start_page": page_num
            })
            break

for i, ch in enumerate(chapters):
    if i < len(chapters) - 1:
        ch["end_page"] = chapters[i + 1]["start_page"] - 1
    else:
        ch["end_page"] = int(files[-1].replace(".txt", ""))

for ch in chapters:
    print(f"Chapter {ch['chapter']}: {ch['title']} (pages {ch['start_page']}-{ch['end_page']})")
