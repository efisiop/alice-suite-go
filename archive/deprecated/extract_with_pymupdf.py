#!/usr/bin/env python3
"""
Extract text from Alice in Wonderland PDF using PyMuPDF (fitz).
PyMuPDF handles multi-column layouts better than pdfplumber in some cases.
"""

import sys
import re

def reflow_text(text):
    """Reflow text by removing line breaks within paragraphs"""
    if not text:
        return ""
    
    # Filter metadata
    filtered_lines = []
    skip_patterns = [
        r'^Fit Page', r'^Full Screen', r'^Close Book', r'^Navigate Control',
        r'^Internet', r'^Digital Interface', r'^BookVirtual', r'^U\.S\. Patent',
        r'^All Rights Reserved', r'^© \d{4}',
        r'DDiiggiittaall', r'InItnetrefrafcaec', r'oBookoVkiVritrutaul',
        r'SP\.a tePnatt enPte ndPeinndgi', r'ARllig hRtsi ghRtess erRveesde',
    ]
    
    lines = text.split('\n')
    for line in lines:
        should_skip = False
        line_stripped = line.strip()
        
        if len(line_stripped) < 3:
            should_skip = True
        
        for pattern in skip_patterns:
            if re.search(pattern, line, re.IGNORECASE):
                should_skip = True
                break
        
        if line_stripped and not should_skip:
            alpha_count = sum(1 for c in line_stripped if c.isalpha())
            if len(line_stripped) > 10 and alpha_count / len(line_stripped) < 0.3:
                should_skip = True
        
        if not should_skip:
            filtered_lines.append(line)
    
    text = '\n'.join(filtered_lines)
    paragraphs = re.split(r'\n\s*\n+', text)
    
    reflowed = []
    for para in paragraphs:
        if not para.strip():
            continue
        para_text = para.replace('\n', ' ')
        para_text = re.sub(r' +', ' ', para_text).strip()
        if para_text:
            reflowed.append(para_text)
    
    return '\n\n'.join(reflowed)

def extract_with_pymupdf(pdf_path, output_path):
    """Extract text using PyMuPDF"""
    try:
        import fitz  # PyMuPDF
        
        print("Extracting text from PDF using PyMuPDF...")
        all_text = []
        
        doc = fitz.open(pdf_path)
        total_pdf_pages = len(doc)
        print(f"Processing {total_pdf_pages} PDF pages...")
        
        for page_num in range(total_pdf_pages):
            print(f"Extracting PDF page {page_num + 1}/{total_pdf_pages}...", end='\r')
            page = doc[page_num]
            
            # Try different extraction methods
            # Method 1: Extract with blocks (preserves layout)
            try:
                blocks = page.get_text("blocks")
                # Sort blocks by position (top to bottom, left to right)
                blocks.sort(key=lambda b: (b[1], b[0]))  # Sort by y, then x
                text = '\n'.join([block[4] for block in blocks if block[4].strip()])
            except:
                # Fall back to regular text extraction
                text = page.get_text()
            
            if text:
                all_text.append(text)
        
        doc.close()
        print(f"\n✅ Extracted text from {total_pdf_pages} PDF pages")
        
        # Combine all text
        full_text = '\n'.join(all_text)
        
        # Find ALL page markers: "N [TITLE]. M" pattern
        page_marker_pattern = r'(\d+)\s+([A-Z][A-Z\s\-]+?)\.\s+(\d+)'
        
        print("\nFinding page markers in text...")
        markers = []
        for match in re.finditer(page_marker_pattern, full_text):
            page_num = int(match.group(1))
            next_page_num = int(match.group(3))
            # Only add if it looks like a valid page marker
            if next_page_num == page_num + 1 or next_page_num == page_num + 2:
                markers.append({
                    'page': page_num,
                    'next_page': next_page_num,
                    'position': match.start(),
                    'end_position': match.end(),
                    'text': match.group(0)
                })
        
        # Sort markers by position
        markers.sort(key=lambda x: x['position'])
        
        print(f"✅ Found {len(markers)} page markers")
        if markers:
            print(f"   Page range: {markers[0]['page']} - {markers[-1]['next_page']}")
        
        # Split text by markers
        pages = {}
        
        # Handle content before first marker (title pages, etc.) as page 1
        if markers:
            first_marker = markers[0]
            content_before = full_text[:first_marker['position']].strip()
            if content_before:
                reflowed = reflow_text(content_before)
                if reflowed.strip():
                    pages[1] = reflowed
        
        # Process each page based on markers
        for i, marker in enumerate(markers):
            page_num = marker['page']
            next_page_num = marker['next_page']
            
            # Content for page_num is from previous marker (or start) to this marker
            if i == 0:
                content = full_text[:marker['position']].strip()
            else:
                prev_marker = markers[i - 1]
                content = full_text[prev_marker['end_position']:marker['position']].strip()
            
            if content:
                reflowed = reflow_text(content)
                if reflowed.strip():
                    pages[page_num] = reflowed
            
            # Content for next_page_num is from this marker to next marker (or end)
            if i + 1 < len(markers):
                next_marker = markers[i + 1]
                content = full_text[marker['end_position']:next_marker['position']].strip()
            else:
                content = full_text[marker['end_position']:].strip()
            
            if content:
                reflowed = reflow_text(content)
                if reflowed.strip():
                    pages[next_page_num] = reflowed
        
        print(f"\n✅ Extracted {len(pages)} physical pages")
        
        # Write output
        print(f"\nWriting {len(pages)} pages to {output_path}...")
        with open(output_path, 'w', encoding='utf-8') as output_file:
            for page_num in sorted(pages.keys()):
                output_file.write(f"{'='*80}\n")
                output_file.write(f"PAGE {page_num}\n")
                output_file.write(f"{'='*80}\n\n")
                output_file.write(pages[page_num])
                output_file.write("\n\n")
        
        print(f"✅ Successfully extracted {len(pages)} physical pages")
        if pages:
            print(f"   Page range: {min(pages.keys())} - {max(pages.keys())}")
        return True
        
    except ImportError:
        print("PyMuPDF not available. Please install:")
        print("  pip3 install PyMuPDF")
        return False
    except Exception as e:
        print(f"Error: {e}")
        import traceback
        traceback.print_exc()
        return False

if __name__ == "__main__":
    pdf_path = "alice_wonderland.pdf"
    output_path = "alice_wonderland_by_pages.txt"
    
    print(f"Extracting text from {pdf_path}...")
    print(f"Output will be saved to {output_path}\n")
    
    if not extract_with_pymupdf(pdf_path, output_path):
        sys.exit(1)


