import pdfplumber
import os
import sys

pdf_path = r"D:\Tasks\ai_code\books_ai\Books\_media_books_go_razrabotka_prilozhenij_v_mikroservisnoj_arhitekture.pdf"
output_dir = r"D:\Tasks\ai_code\books_ai\Books\_media_books_go_razrabotka_prilozhenij_v_mikroservisnoj_arhitekture\raw"

os.makedirs(output_dir, exist_ok=True)

try:
    with pdfplumber.open(pdf_path) as pdf:
        total_pages = len(pdf.pages)
        print(f"Total pages: {total_pages}")
        
        for i, page in enumerate(pdf.pages, start=1):
            text = page.extract_text()
            if text is None:
                text = ""
            
            filename = f"{i:03d}.txt"
            filepath = os.path.join(output_dir, filename)
            
            with open(filepath, "w", encoding="utf-8") as f:
                f.write(text)
            
            if i % 50 == 0:
                print(f"Processed {i}/{total_pages} pages...")
        
        print(f"Done. Extracted {total_pages} pages to {output_dir}")
except Exception as e:
    print(f"Error: {e}", file=sys.stderr)
    sys.exit(1)
