import os
import json

base_dir = r"D:\Tasks\ai_code\books_ai\Books\Explore_Go_Cryptography - John Arundel - 2026 - v1"
metadata_path = os.path.join(base_dir, "metadata.json")

with open(metadata_path, "r", encoding="utf-8") as f:
    metadata = json.load(f)

for ch in metadata["chapters"]:
    chapter_num = ch["chapter"]
    slug = ch["slug"]
    title = ch["title"]
    start_page = ch["start_page"]
    end_page = ch["end_page"]
    
    chapter_dir = os.path.join(base_dir, "chapters", f"{chapter_num:02d}-{slug}")
    os.makedirs(chapter_dir, exist_ok=True)
    
    index_path = os.path.join(chapter_dir, "index.md")
    content = f"""# {chapter_num}. {title}

**Səhifələr:** {start_page}-{end_page}

## Bu fəsil nədən bəhs edir?

(İçerik buraya əlavə ediləcək)

## Əsas fikirlər

(İçerik buraya əlavə ediləcək)

## Əsas terminlər

(İçerik buraya əlavə ediləcək)

## Praktik nəticə

(İçerik buraya əlavə ediləcək)

## Mənbə

Pages: {start_page}-{end_page}
"""
    with open(index_path, "w", encoding="utf-8") as f:
        f.write(content)
    
    cheatsheet_path = os.path.join(chapter_dir, "cheatsheet.md")
    cheatsheet_content = f"""# Cheat Sheet — {title}

(İçerik buraya əlavə ediləcək)
"""
    with open(cheatsheet_path, "w", encoding="utf-8") as f:
        f.write(cheatsheet_content)

print(f"Created {len(metadata['chapters'])} chapter folders")
