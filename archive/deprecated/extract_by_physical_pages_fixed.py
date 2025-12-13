#!/usr/bin/env python3
"""
Extract text from Alice in Wonderland PDF by physical book page numbers.
Maps PDF pages to physical pages correctly:
- Physical page 2 = PDF pages 2+3
- Physical page 3 = PDF pages 4+5 (missing odd pages)
- Physical page 4 = PDF pages 6+7
- Physical page 5 = PDF pages 8+9
- etc.
"""

import sys
import re

def reflow_text(text):
    """
    Reflow text by removing line breaks within paragraphs.
    - Preserves paragraph breaks (double newlines)
    - Removes single line breaks (replaces with spaces)
    - Cleans up extra whitespace
    """
    if not text:
        return ""
    
    # Filter out common PDF navigation and metadata text
    filtered_lines = []
    skip_patterns = [
        r'^Fit Page',
        r'^Full Screen',
        r'^Close Book',
        r'^Navigate Control',
        r'^Internet',
        r'^Digital Interface',
        r'^BookVirtual',
        r'^U\.S\. Patent',
        r'^All Rights Reserved',
        r'^© \d{4}',
        r'DDiiggiittaall',
        r'InItnetrefrafcaec',
        r'oBookoVkiVritrutaul',
        r'SP\.a tePnatt enPte ndPeinndgi',
        r'ARllig hRtsi ghRtess erRveesde',
    ]
    
    lines = text.split('\n')
    for line in lines:
        should_skip = False
        line_stripped = line.strip()
        
        # Skip very short lines
        if len(line_stripped) < 3:
            should_skip = True
        
        # Check against skip patterns
        for pattern in skip_patterns:
            if re.search(pattern, line, re.IGNORECASE):
                should_skip = True
                break
        
        # Skip lines that are mostly non-alphabetic
        if line_stripped and not should_skip:
            alpha_count = sum(1 for c in line_stripped if c.isalpha())
            if len(line_stripped) > 10 and alpha_count / len(line_stripped) < 0.3:
                should_skip = True
        
        if not should_skip:
            filtered_lines.append(line)
    
    text = '\n'.join(filtered_lines)
    
    # Split by paragraph breaks
    paragraphs = re.split(r'\n\s*\n+', text)
    
    reflowed = []
    for para in paragraphs:
        if not para.strip():
            continue
        
        # Remove single line breaks within paragraph
        para_text = para.replace('\n', ' ')
        para_text = re.sub(r' +', ' ', para_text)
        para_text = para_text.strip()
        
        if para_text:
            reflowed.append(para_text)
    
    return '\n\n'.join(reflowed)

def extract_by_physical_pages(pdf_path, output_path):
    """Extract text and map PDF pages to physical book pages"""
    try:
        import pdfplumber
        
        print("Extracting text from PDF pages...")
        pdf_pages_text = []
        
        with pdfplumber.open(pdf_path) as pdf:
            total_pdf_pages = len(pdf.pages)
            print(f"Processing {total_pdf_pages} PDF pages...")
            
            for i, page in enumerate(pdf.pages, start=1):
                print(f"Extracting PDF page {i}/{total_pdf_pages}...", end='\r')
                text = page.extract_text()
                pdf_pages_text.append(text if text else "")
        
        print(f"\n✅ Extracted text from {total_pdf_pages} PDF pages")
        
        # Map PDF pages to physical pages
        # Based on user's findings:
        # - Physical page 2 = PDF pages 2+3
        # - Physical page 4 = PDF pages 4+5
        # - Physical page 6 = PDF pages 6+7
        # - Physical page 8 = PDF pages 8+9
        # So pattern: Physical page N (even) = PDF pages N and N+1
        # And Physical page N+1 (odd) = PDF pages N+2 and N+3
        
        physical_pages = {}
        
        # Physical page 1 = PDF page 1 (title/intro)
        if len(pdf_pages_text) > 0:
            text = reflow_text(pdf_pages_text[0])
            if text.strip():
                physical_pages[1] = text
        
        # Map remaining PDF pages to physical pages
        # Starting from PDF page 2, pair them up
        pdf_idx = 1  # Start from PDF page 2 (index 1)
        physical_page = 2  # Start from physical page 2
        
        print("\nMapping PDF pages to physical pages...")
        
        while pdf_idx < len(pdf_pages_text):
            # Combine two consecutive PDF pages into one physical page
            if pdf_idx + 1 < len(pdf_pages_text):
                combined_text = pdf_pages_text[pdf_idx] + '\n' + pdf_pages_text[pdf_idx + 1]
            else:
                combined_text = pdf_pages_text[pdf_idx]
            
            reflowed = reflow_text(combined_text)
            if reflowed.strip():
                physical_pages[physical_page] = reflowed
                print(f"  Physical page {physical_page} = PDF pages {pdf_idx + 1} + {pdf_idx + 2 if pdf_idx + 1 < len(pdf_pages_text) else pdf_idx + 1}", end='\r')
            
            pdf_idx += 2  # Move to next pair
            physical_page += 1  # Next physical page
        
        print(f"\n✅ Mapped to {len(physical_pages)} physical pages")
        
        # Write output file
        print(f"\nWriting {len(physical_pages)} pages to {output_path}...")
        with open(output_path, 'w', encoding='utf-8') as output_file:
            # Write pages in order
            for page_num in sorted(physical_pages.keys()):
                output_file.write(f"{'='*80}\n")
                output_file.write(f"PAGE {page_num}\n")
                output_file.write(f"{'='*80}\n\n")
                output_file.write(physical_pages[page_num])
                output_file.write("\n\n")
        
        print(f"✅ Successfully extracted {len(physical_pages)} physical pages")
        print(f"   Page range: {min(physical_pages.keys())} - {max(physical_pages.keys())}")
        return True
        
    except ImportError:
        print("pdfplumber not available. Please install:")
        print("  pip3 install pdfplumber")
        return False
    except Exception as e:
        print(f"Error: {e}")
        import traceback
        traceback.print_exc()
        return False

def main():
    pdf_path = "alice_wonderland.pdf"
    output_path = "alice_wonderland_by_pages.txt"
    
    print(f"Extracting text from {pdf_path}...")
    print(f"Output will be saved to {output_path}\n")
    
    if not extract_by_physical_pages(pdf_path, output_path):
        sys.exit(1)

if __name__ == "__main__":
    main()


