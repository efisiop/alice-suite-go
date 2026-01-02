-- Export sections data from localhost database
-- Run this on localhost: sqlite3 data/alice-suite.db < scripts/export-sections.sql > scripts/sections-data.sql

-- Export sections in INSERT format
.mode insert sections
SELECT * FROM sections ORDER BY page_number, section_number;

