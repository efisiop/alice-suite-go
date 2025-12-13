#!/usr/bin/env python3
"""
Extract text from Alice in Wonderland PDF by physical book page numbers.
Splits content by physical page numbers (1-192) and reflows text within each page.
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
    """Extract text and split by physical book page numbers"""
    try:
        import pdfplumber
        
        print("Extracting all text from PDF...")
        all_text = []
        
        with pdfplumber.open(pdf_path) as pdf:
            total_pdf_pages = len(pdf.pages)
            print(f"Processing {total_pdf_pages} PDF pages...")
            
            for i, page in enumerate(pdf.pages, start=1):
                print(f"Extracting PDF page {i}/{total_pdf_pages}...", end='\r')
                text = page.extract_text()
                if text:
                    all_text.append(text)
        
        print(f"\n✅ Extracted text from {total_pdf_pages} PDF pages")
        
        # Combine all text
        full_text = '\n'.join(all_text)
        
        # Find physical page number markers
        # Pattern: number followed by chapter title, then another number
        # Example: "2 DOWN THE RABBIT-HOLE. 3"
        page_pattern = r'^(\d+)\s+([A-Z][^.]*?)\.\s+(\d+)'
        
        # Also look for standalone page numbers at start of lines
        # Pattern: number at start of line, possibly followed by chapter title
        standalone_pattern = r'^(\d+)(?:\s+[A-Z][^.]*?\.\s+(\d+))?'
        
        lines = full_text.split('\n')
        pages = {}  # page_num -> content
        current_page = 1
        current_content = []
        
        print("\nParsing physical page numbers...")
        
        i = 0
        while i < len(lines):
            line = lines[i]
            
            # Check for page number pattern
            match = re.search(page_pattern, line, re.MULTILINE)
            if not match:
                match = re.search(standalone_pattern, line)
            
            if match:
                page_num = int(match.group(1))
                
                # Save previous page content
                if current_content:
                    content_text = '\n'.join(current_content)
                    reflowed = reflow_text(content_text)
                    if reflowed.strip():
                        pages[current_page] = reflowed
                
                # Start new page
                current_page = page_num
                current_content = []
                
                # Add the rest of this line (after page number) to new page
                remaining = line[match.end():].strip()
                if remaining:
                    current_content.append(remaining)
            else:
                # Add line to current page
                if line.strip():
                    current_content.append(line)
            
            i += 1
        
        # Save last page
        if current_content:
            content_text = '\n'.join(current_content)
            reflowed = reflow_text(content_text)
            if reflowed.strip():
                pages[current_page] = reflowed
        
        # Write output file
        print(f"\nWriting {len(pages)} pages to {output_path}...")
        with open(output_path, 'w', encoding='utf-8') as output_file:
            # Write pages in order
            for page_num in sorted(pages.keys()):
                output_file.write(f"{'='*80}\n")
                output_file.write(f"PAGE {page_num}\n")
                output_file.write(f"{'='*80}\n\n")
                output_file.write(pages[page_num])
                output_file.write("\n\n")
        
        print(f"✅ Successfully extracted {len(pages)} physical pages")
        print(f"   Page range: {min(pages.keys())} - {max(pages.keys())}")
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



