import pdfplumber, os, sys, re

pdf_path = r"Books\Go_Recipes_Aaron_Torres_2019_v2\source\original.pdf"
output_dir = r"Books\Go_Recipes_Aaron_Torres_2019_v2\raw"

with pdfplumber.open(pdf_path) as pdf:
    total = len(pdf.pages)
    print(f"Total pages: {total}")
    for i, page in enumerate(pdf.pages, start=1):
        text = page.extract_text() or ""
        lines = text.split("\n")
        # noise filter: drop 'Tlgm: @it_boooks' watermark line
        lines = [ln for ln in lines if ln.strip() != "Tlgm: @it_boooks"]
        text = "\n".join(lines)
        with open(os.path.join(output_dir, f"{i:03d}.txt"), "w", encoding="utf-8") as f:
            f.write(text)
        if i % 50 == 0:
            print(f"Processed {i}/{total}")
print("DONE")
