#!/usr/bin/env python3
"""
Extract text from Alice in Wonderland PDF using page markers to correctly split content.
Page markers like "2 DOWN THE RABBIT-HOLE. 3" indicate where physical pages begin/end.
This fixes the column layout issue by using the actual page boundaries in the text.
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

def extract_with_page_markers(pdf_path, output_path):
    """Extract text and split using page markers found in the text"""
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
        
        # Find page markers: "N [CHAPTER TITLE]. M" where N is current page, M is next page
        # Example: "2 DOWN THE RABBIT-HOLE. 3"
        page_marker_pattern = r'(\d+)\s+([A-Z][A-Z\s]+?)\.\s+(\d+)'
        
        print("\nFinding page markers in text...")
        markers = []
        for match in re.finditer(page_marker_pattern, full_text):
            page_num = int(match.group(1))
            next_page_num = int(match.group(3))
            markers.append({
                'page': page_num,
                'next_page': next_page_num,
                'position': match.start(),
                'end_position': match.end(),
                'text': match.group(0)
            })
        
        print(f"✅ Found {len(markers)} page markers")
        if markers:
            print(f"   First marker: {markers[0]['text']} at position {markers[0]['position']}")
            print(f"   Last marker: {markers[-1]['text']} at position {markers[-1]['position']}")
        
        # Split text by markers
        pages = {}
        
        # Handle content before first marker (title pages, etc.)
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
            
            # Find where this page's content ends
            if i + 1 < len(markers):
                next_marker = markers[i + 1]
                # Content is from end of current marker to start of next marker
                content = full_text[marker['end_position']:next_marker['position']].strip()
            else:
                # Last marker - content goes to end of text
                content = full_text[marker['end_position']:].strip()
            
            if content:
                reflowed = reflow_text(content)
                if reflowed.strip():
                    pages[page_num] = reflowed
                    print(f"  Extracted page {page_num}", end='\r')
        
        print(f"\n✅ Extracted {len(pages)} physical pages")
        
        # Check for missing pages (gaps in page numbers)
        found_pages = sorted(pages.keys())
        missing = []
        for i in range(1, len(found_pages)):
            if found_pages[i] - found_pages[i-1] > 1:
                gap_start = found_pages[i-1] + 1
                gap_end = found_pages[i] - 1
                missing.extend(range(gap_start, gap_end + 1))
        
        if missing:
            print(f"⚠️  Missing pages detected: {missing[:10]}... (showing first 10)")
            print("   These pages may need manual extraction or are blank pages")
        
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

if __name__ == "__main__":
    pdf_path = "alice_wonderland.pdf"
    output_path = "alice_wonderland_by_pages.txt"
    
    print(f"Extracting text from {pdf_path}...")
    print(f"Output will be saved to {output_path}\n")
    
    if not extract_with_page_markers(pdf_path, output_path):
        sys.exit(1)

