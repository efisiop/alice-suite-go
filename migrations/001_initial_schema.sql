-- SQLite Schema for Alice Suite Go
-- Physical Book Companion App
-- Initial schema for first 3 chapters test ground

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Users/Profiles
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  first_name TEXT,
  last_name TEXT,
  role TEXT CHECK (role IN ('reader', 'consultant')) DEFAULT 'reader',
  is_verified INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now'))
);

-- Books
CREATE TABLE IF NOT EXISTS books (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  description TEXT,
  total_pages INTEGER NOT NULL,
  created_at TEXT DEFAULT (datetime('now'))
);

-- Chapters (First 3 chapters only for test ground)
CREATE TABLE IF NOT EXISTS chapters (
  id TEXT PRIMARY KEY,
  book_id TEXT NOT NULL,
  title TEXT NOT NULL,
  number INTEGER NOT NULL,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  UNIQUE(book_id, number)
);

-- Sections (Book sections for reference/word clarification)
CREATE TABLE IF NOT EXISTS sections (
  id TEXT PRIMARY KEY,
  chapter_id TEXT NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  start_page INTEGER NOT NULL,
  end_page INTEGER NOT NULL,
  number INTEGER NOT NULL,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
  UNIQUE(chapter_id, number)
);

-- Alice Glossary (Alice-specific definitions)
CREATE TABLE IF NOT EXISTS alice_glossary (
  id TEXT PRIMARY KEY,
  book_id TEXT NOT NULL,
  term TEXT NOT NULL,
  definition TEXT NOT NULL,
  source_sentence TEXT,
  example TEXT,
  chapter_reference TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  UNIQUE(book_id, term)
);

-- Verification Codes (for book access)
CREATE TABLE IF NOT EXISTS verification_codes (
  code TEXT PRIMARY KEY,
  book_id TEXT NOT NULL,
  is_used INTEGER DEFAULT 0,
  used_by TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (book_id) REFERENCES books(id),
  FOREIGN KEY (used_by) REFERENCES users(id)
);

-- Reading Progress (tracks progress in physical book)
CREATE TABLE IF NOT EXISTS reading_progress (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  chapter_id TEXT,
  section_id TEXT,
  last_page INTEGER,
  last_read_at TEXT DEFAULT (datetime('now')),
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
  UNIQUE(user_id, book_id)
);

-- Vocabulary Lookups (words looked up from physical book)
CREATE TABLE IF NOT EXISTS vocabulary_lookups (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  word TEXT NOT NULL,
  definition TEXT NOT NULL,
  chapter_id TEXT,
  section_id TEXT,
  context TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);

-- AI Interactions (AI assistance requests)
CREATE TABLE IF NOT EXISTS ai_interactions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  section_id TEXT,
  interaction_type TEXT CHECK (interaction_type IN ('explain', 'quiz', 'simplify', 'definition', 'chat')) DEFAULT 'chat',
  question TEXT,
  prompt TEXT,
  response TEXT NOT NULL,
  context TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);

-- Help Requests (Tier 3: Human consultant support)
CREATE TABLE IF NOT EXISTS help_requests (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  section_id TEXT,
  status TEXT CHECK (status IN ('pending', 'assigned', 'resolved')) DEFAULT 'pending',
  content TEXT NOT NULL,
  context TEXT,
  assigned_to TEXT,
  response TEXT,
  resolved_at TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE,
  FOREIGN KEY (assigned_to) REFERENCES users(id)
);

-- Consultant Assignments
CREATE TABLE IF NOT EXISTS consultant_assignments (
  id TEXT PRIMARY KEY,
  consultant_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  active INTEGER DEFAULT 1,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (consultant_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  UNIQUE(consultant_id, user_id, book_id)
);

-- Consultant Triggers (Prompts sent by consultants)
CREATE TABLE IF NOT EXISTS consultant_triggers (
  id TEXT PRIMARY KEY,
  consultant_id TEXT,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  trigger_type TEXT NOT NULL,
  message TEXT,
  is_processed INTEGER DEFAULT 0,
  processed_at TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (consultant_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

-- User Feedback
CREATE TABLE IF NOT EXISTS user_feedback (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  section_id TEXT,
  feedback_type TEXT NOT NULL,
  content TEXT NOT NULL,
  is_public INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);

-- Reading Statistics
CREATE TABLE IF NOT EXISTS reading_stats (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  book_id TEXT NOT NULL,
  total_reading_time INTEGER DEFAULT 0,
  pages_read INTEGER DEFAULT 0,
  vocabulary_words INTEGER DEFAULT 0,
  last_session_date TEXT DEFAULT (datetime('now')),
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  UNIQUE(user_id, book_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_chapters_book_id ON chapters(book_id);
CREATE INDEX IF NOT EXISTS idx_sections_chapter_id ON sections(chapter_id);
CREATE INDEX IF NOT EXISTS idx_alice_glossary_book_term ON alice_glossary(book_id, term);
CREATE INDEX IF NOT EXISTS idx_vocabulary_lookups_user_book ON vocabulary_lookups(user_id, book_id);
CREATE INDEX IF NOT EXISTS idx_ai_interactions_user_book ON ai_interactions(user_id, book_id);
CREATE INDEX IF NOT EXISTS idx_help_requests_status ON help_requests(status);
CREATE INDEX IF NOT EXISTS idx_consultant_assignments_consultant ON consultant_assignments(consultant_id);



