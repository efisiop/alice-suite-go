-- Migration 004: Link Glossary Terms to Sections
-- Creates junction table to link glossary terms with sections where they appear

PRAGMA foreign_keys = ON;

-- Junction table linking glossary terms to sections
CREATE TABLE IF NOT EXISTS glossary_section_links (
  id TEXT PRIMARY KEY,
  glossary_id TEXT NOT NULL,
  section_id TEXT NOT NULL,
  page_number INTEGER NOT NULL,
  section_number INTEGER NOT NULL,
  term TEXT NOT NULL, -- Denormalized for quick lookup
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (glossary_id) REFERENCES alice_glossary(id) ON DELETE CASCADE,
  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
  UNIQUE(glossary_id, section_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_glossary_links_term ON glossary_section_links(term);
CREATE INDEX IF NOT EXISTS idx_glossary_links_section ON glossary_section_links(section_id);
CREATE INDEX IF NOT EXISTS idx_glossary_links_page ON glossary_section_links(page_number, section_number);
CREATE INDEX IF NOT EXISTS idx_glossary_term ON alice_glossary(term);



