#!/usr/bin/env python3
"""
Extract text from Alice in Wonderland PDF with page numbers and reflowed text.
Preserves paragraph breaks but removes unnecessary line breaks within paragraphs.
"""

import sys
import re

def reflow_text(text):
    """
    Reflow text by removing line breaks within paragraphs.
    - Preserves paragraph breaks (double newlines)
    - Removes single line breaks (replaces with spaces)
    - Cleans up extra whitespace
    - Filters out common PDF navigation/metadata text
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
        # Skip lines that match skip patterns
        should_skip = False
        line_stripped = line.strip()
        
        # Skip very short lines that are likely metadata
        if len(line_stripped) < 3:
            should_skip = True
        
        # Check against skip patterns
        for pattern in skip_patterns:
            if re.search(pattern, line, re.IGNORECASE):
                should_skip = True
                break
        
        # Skip lines that are mostly non-alphabetic (likely corrupted metadata)
        if line_stripped and not should_skip:
            alpha_count = sum(1 for c in line_stripped if c.isalpha())
            if len(line_stripped) > 10 and alpha_count / len(line_stripped) < 0.3:
                should_skip = True
        
        if not should_skip:
            filtered_lines.append(line)
    
    text = '\n'.join(filtered_lines)
    
    # Split by paragraph breaks (double newlines or more)
    paragraphs = re.split(r'\n\s*\n+', text)
    
    reflowed = []
    for para in paragraphs:
        if not para.strip():
            continue
        
        # Remove single line breaks within paragraph
        # Replace single newlines with spaces
        para_text = para.replace('\n', ' ')
        
        # Clean up multiple spaces
        para_text = re.sub(r' +', ' ', para_text)
        
        # Trim whitespace
        para_text = para_text.strip()
        
        if para_text:
            reflowed.append(para_text)
    
    return '\n\n'.join(reflowed)

def extract_with_pdfplumber(pdf_path, output_path):
    """Extract using pdfplumber (better quality)"""
    try:
        import pdfplumber
        
        with pdfplumber.open(pdf_path) as pdf:
            with open(output_path, 'w', encoding='utf-8') as output_file:
                total_pages = len(pdf.pages)
                print(f"Processing {total_pages} pages...")
                
                for i, page in enumerate(pdf.pages, start=1):
                    print(f"Extracting page {i}/{total_pages}...", end='\r')
                    
                    # Extract text
                    text = page.extract_text()
                    
                    if text:
                        # Reflow the text
                        reflowed_text = reflow_text(text)
                        
                        # Write page header
                        output_file.write(f"{'='*80}\n")
                        output_file.write(f"PAGE {i}\n")
                        output_file.write(f"{'='*80}\n\n")
                        
                        # Write reflowed content
                        output_file.write(reflowed_text)
                        output_file.write("\n\n")
                
                print(f"\n✅ Successfully extracted {total_pages} pages to {output_path}")
                return True
                
    except ImportError:
        print("pdfplumber not available, trying alternative method...")
        return False
    except Exception as e:
        print(f"Error with pdfplumber: {e}")
        return False

def extract_with_pdftotext(pdf_path, output_path):
    """Extract using pdftotext command-line tool"""
    import subprocess
    
    try:
        # First, extract raw text
        result = subprocess.run(
            ['pdftotext', '-layout', pdf_path, '-'],
            capture_output=True,
            text=True,
            check=True
        )
        
        raw_text = result.stdout
        
        # Split by page breaks (pdftotext doesn't add page markers, so we need to estimate)
        # For now, let's use a simpler approach - extract page by page
        with open(output_path, 'w', encoding='utf-8') as output_file:
            # Try to get page count first
            page_count_result = subprocess.run(
                ['pdfinfo', pdf_path],
                capture_output=True,
                text=True
            )
            
            # Extract page count from pdfinfo output
            page_count = 1
            for line in page_count_result.stdout.split('\n'):
                if 'Pages:' in line:
                    try:
                        page_count = int(line.split(':')[1].strip())
                    except:
                        pass
                    break
            
            # Extract each page individually
            for page_num in range(1, page_count + 1):
                print(f"Extracting page {page_num}/{page_count}...", end='\r')
                
                page_result = subprocess.run(
                    ['pdftotext', '-f', str(page_num), '-l', str(page_num), '-layout', pdf_path, '-'],
                    capture_output=True,
                    text=True,
                    check=True
                )
                
                text = page_result.stdout
                if text.strip():
                    reflowed_text = reflow_text(text)
                    
                    output_file.write(f"{'='*80}\n")
                    output_file.write(f"PAGE {page_num}\n")
                    output_file.write(f"{'='*80}\n\n")
                    output_file.write(reflowed_text)
                    output_file.write("\n\n")
            
            print(f"\n✅ Successfully extracted {page_count} pages to {output_path}")
            return True
            
    except FileNotFoundError:
        print("pdftotext not found. Please install poppler-utils:")
        print("  brew install poppler")
        return False
    except subprocess.CalledProcessError as e:
        print(f"Error running pdftotext: {e}")
        return False
    except Exception as e:
        print(f"Error: {e}")
        return False

def main():
    pdf_path = "alice_wonderland.pdf"
    output_path = "alice_wonderland_extracted.txt"
    
    print(f"Extracting text from {pdf_path}...")
    print(f"Output will be saved to {output_path}\n")
    
    # Try pdfplumber first (better quality)
    if not extract_with_pdfplumber(pdf_path, output_path):
        # Fall back to pdftotext
        if not extract_with_pdftotext(pdf_path, output_path):
            print("\n❌ Failed to extract text. Please install pdfplumber:")
            print("  pip3 install pdfplumber")
            sys.exit(1)

if __name__ == "__main__":
    main()


