#!/usr/bin/env python3
"""
Extract glossary terms, definitions, and references from Vocabulary.com HTML file.
Outputs JSON and SQL for importing into the database.
"""

import re
import json
from html.parser import HTMLParser
from html import unescape

class GlossaryExtractor(HTMLParser):
    def __init__(self):
        super().__init__()
        self.entries = []
        self.current_entry = None
        self.in_entry = False
        self.in_word_link = False
        self.in_definition = False
        self.in_example = False
        self.current_text = ""
        
    def handle_starttag(self, tag, attrs):
        attrs_dict = dict(attrs)
        
        # Check if this is a vocabulary entry
        if tag == 'li' and 'class' in attrs_dict and 'entry' in attrs_dict['class']:
            self.in_entry = True
            self.current_entry = {
                'word': attrs_dict.get('word', ''),
                'frequency': float(attrs_dict.get('freq', 0)) if attrs_dict.get('freq') else None,
                'learnable': 'learnable' in attrs_dict.get('class', ''),
                'definition': '',
                'examples': []
            }
        
        # Word link
        elif tag == 'a' and 'class' in attrs_dict and attrs_dict['class'] == 'word':
            self.in_word_link = True
            # Get definition from title attribute
            if 'title' in attrs_dict and not self.current_entry['definition']:
                self.current_entry['definition'] = unescape(attrs_dict['title']).strip()
        
        # Definition div
        elif tag == 'div' and 'class' in attrs_dict and attrs_dict['class'] == 'definition':
            self.in_definition = True
            self.current_text = ""
        
        # Example div
        elif tag == 'div' and 'class' in attrs_dict and attrs_dict['class'] == 'example':
            self.in_example = True
            self.current_text = ""
        
        # Strong tag in example (highlighted word)
        elif tag == 'strong' and self.in_example:
            pass  # We'll handle the text content
    
    def handle_endtag(self, tag):
        if tag == 'li' and self.in_entry:
            # Save entry if it has a word and definition
            if self.current_entry and self.current_entry['word'] and self.current_entry['definition']:
                self.entries.append(self.current_entry)
            self.in_entry = False
            self.current_entry = None
        
        elif tag == 'a' and self.in_word_link:
            self.in_word_link = False
        
        elif tag == 'div' and self.in_definition:
            # Save definition text
            if self.current_entry and self.current_text.strip():
                if not self.current_entry['definition']:
                    self.current_entry['definition'] = unescape(self.current_text.strip())
            self.in_definition = False
            self.current_text = ""
        
        elif tag == 'div' and self.in_example:
            # Save example text
            if self.current_entry and self.current_text.strip():
                example_text = unescape(self.current_text.strip())
                # Clean up the example text
                example_text = re.sub(r'\s+', ' ', example_text)
                if example_text:
                    self.current_entry['examples'].append(example_text)
            self.in_example = False
            self.current_text = ""
    
    def handle_data(self, data):
        if self.in_entry:
            if self.in_definition:
                self.current_text += data
            elif self.in_example:
                self.current_text += data

def extract_glossary(html_path, json_output, sql_output):
    """Extract glossary from HTML file"""
    print(f"Reading HTML file: {html_path}")
    
    with open(html_path, 'r', encoding='utf-8') as f:
        html_content = f.read()
    
    print("Parsing HTML...")
    parser = GlossaryExtractor()
    parser.feed(html_content)
    
    entries = parser.entries
    print(f"✅ Extracted {len(entries)} glossary entries")
    
    # Write JSON output
    print(f"\nWriting JSON to {json_output}...")
    with open(json_output, 'w', encoding='utf-8') as f:
        json.dump(entries, f, indent=2, ensure_ascii=False)
    print(f"✅ JSON written: {len(entries)} entries")
    
    # Write SQL output
    print(f"\nWriting SQL to {sql_output}...")
    with open(sql_output, 'w', encoding='utf-8') as f:
        f.write("-- Alice's Adventures in Wonderland Glossary\n")
        f.write("-- Extracted from Vocabulary.com\n")
        f.write("-- Book ID: alice-in-wonderland\n\n")
        f.write("BEGIN TRANSACTION;\n\n")
        
        book_id = 'alice-in-wonderland'
        
        for i, entry in enumerate(entries):
            # Generate ID (simple format: glossary-{index})
            entry_id = f"glossary-{i+1}"
            
            # Escape single quotes for SQL
            word = entry['word'].replace("'", "''")
            definition = entry['definition'].replace("'", "''")
            
            # Get first example (or empty string)
            example_text = entry['examples'][0] if entry['examples'] else ''
            example_text = example_text.replace("'", "''")
            
            # Try to extract chapter reference from example if it contains "CHAPTER"
            chapter_ref = None
            for ex in entry['examples']:
                chapter_match = re.search(r'CHAPTER\s+([IVX]+)', ex, re.IGNORECASE)
                if chapter_match:
                    chapter_ref = chapter_match.group(1)
                    break
            
            # Build SQL INSERT statement
            chapter_ref_sql = f"'{chapter_ref}'" if chapter_ref else "NULL"
            example_sql = f"'{example_text}'" if example_text else "NULL"
            
            sql = f"""INSERT OR IGNORE INTO alice_glossary (id, book_id, term, definition, example, chapter_reference)
VALUES ('{entry_id}', '{book_id}', '{word}', '{definition}', {example_sql}, {chapter_ref_sql});
"""
            f.write(sql)
        
        f.write("\nCOMMIT;\n")
    
    print(f"✅ SQL written: {len(entries)} entries")
    
    # Print summary
    print(f"\n📊 Summary:")
    print(f"   Total entries: {len(entries)}")
    print(f"   Learnable entries: {sum(1 for e in entries if e['learnable'])}")
    print(f"   Entries with examples: {sum(1 for e in entries if e['examples'])}")
    print(f"   Total examples: {sum(len(e['examples']) for e in entries)}")
    
    # Show first few entries
    print(f"\n📝 First 5 entries:")
    for i, entry in enumerate(entries[:5], 1):
        print(f"\n{i}. {entry['word']}")
        print(f"   Definition: {entry['definition'][:80]}...")
        if entry['examples']:
            print(f"   Example: {entry['examples'][0][:80]}...")
    
    return entries

if __name__ == "__main__":
    import sys
    
    html_path = "../alice-suite/ALICE'S ADVENTURES IN WONDERLAND - Vocabulary List | Vocabulary.com.html"
    json_output = "alice_glossary.json"
    sql_output = "alice_glossary.sql"
    
    if len(sys.argv) > 1:
        html_path = sys.argv[1]
    if len(sys.argv) > 2:
        json_output = sys.argv[2]
    if len(sys.argv) > 3:
        sql_output = sys.argv[3]
    
    extract_glossary(html_path, json_output, sql_output)

