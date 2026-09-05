import os
import re

raw_dir = r"D:\Tasks\ai_code\books_ai\Books\Explore_Go_Cryptography - John Arundel - 2026 - v1\raw"

files = sorted(os.listdir(raw_dir))

chapter_pattern = re.compile(r'^(\d+)\.\s+(.+)$')

chapters = []
current_chapter = None

for filename in files:
    filepath = os.path.join(raw_dir, filename)
    with open(filepath, "r", encoding="utf-8") as f:
        lines = f.readlines()
    
    page_num = int(filename.replace(".txt", ""))
    
    for line in lines:
        line = line.strip()
        match = chapter_pattern.match(line)
        if match:
            chapter_num = int(match.group(1))
            chapter_title = match.group(2)
            
            if current_chapter:
                current_chapter["end_page"] = page_num - 1
                chapters.append(current_chapter)
            
            current_chapter = {
                "chapter": chapter_num,
                "title": chapter_title,
                "start_page": page_num
            }

if current_chapter:
    current_chapter["end_page"] = int(files[-1].replace(".txt", ""))
    chapters.append(current_chapter)

for ch in chapters:
    print(f"Chapter {ch['chapter']}: {ch['title']} (pages {ch['start_page']}-{ch['end_page']})")
