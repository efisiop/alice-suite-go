-- Reader reading-experience reaction to a consultant AI prompt.
ALTER TABLE consultant_prompts ADD COLUMN IF NOT EXISTS reading_reaction TEXT;
ALTER TABLE consultant_prompts ADD COLUMN IF NOT EXISTS reading_reaction_at TEXT;
