-- ============================================================
-- Migration 014: Admin settings (e.g. login email notifications)
-- ============================================================

CREATE TABLE IF NOT EXISTS admin_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Default: login email notifications off, recipient for when enabled
INSERT OR IGNORE INTO admin_settings (key, value) VALUES ('login_email_notifications', '0');
INSERT OR IGNORE INTO admin_settings (key, value) VALUES ('login_email_recipient', 'efisio@mylivemail.net');
