CREATE TABLE IF NOT EXISTS consultant_request_shares (
  id TEXT PRIMARY KEY,
  help_request_id TEXT NOT NULL,
  consultant_id TEXT NOT NULL,
  shared_by TEXT NOT NULL,
  note TEXT,
  created_at TEXT DEFAULT (datetime('now')),
  FOREIGN KEY (help_request_id) REFERENCES help_requests(id) ON DELETE CASCADE,
  FOREIGN KEY (consultant_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (shared_by) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(help_request_id, consultant_id)
);

CREATE INDEX IF NOT EXISTS idx_consultant_request_shares_request ON consultant_request_shares(help_request_id);
