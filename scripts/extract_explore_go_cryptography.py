import pdfplumber
import os
import sys

pdf_name = "Explore_Go_Cryptography - John Arundel - 2026 - v1.pdf"
pdf_path = rf"D:\Tasks\ai_code\books_ai\Books\{pdf_name}"
output_dir = rf"D:\Tasks\ai_code\books_ai\Books\Explore_Go_Cryptography - John Arundel - 2026 - v1"
raw_dir = os.path.join(output_dir, "raw")

os.makedirs(raw_dir, exist_ok=True)

try:
    with pdfplumber.open(pdf_path) as pdf:
        total_pages = len(pdf.pages)
        print(f"Total pages: {total_pages}")
        
        for i, page in enumerate(pdf.pages, start=1):
            text = page.extract_text()
            if text is None:
                text = ""
            
            filename = f"{i:03d}.txt"
            filepath = os.path.join(raw_dir, filename)
            
            with open(filepath, "w", encoding="utf-8") as f:
                f.write(text)
            
            if i % 50 == 0:
                print(f"Processed {i}/{total_pages} pages...")
        
        print(f"Done. Extracted {total_pages} pages to {raw_dir}")
except Exception as e:
    print(f"Error: {e}", file=sys.stderr)
    sys.exit(1)
