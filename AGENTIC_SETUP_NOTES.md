# Agentic Engineering Stack — Setup Notes

**Scope:** These tools are for the **Alice Suite** codebase (this repo) only. Use them only when working on alice-suite-go unless you explicitly use them elsewhere.

## What’s installed

- **Homebrew tap:** `steipete/tap`
- **CLI tools:** Watchman, Poltergeist, Oracle, Summarize, yt-dlp, ffmpeg, MCPorter, Peekaboo
- **Menu bar apps:** Trimmy, CodexBar, RepoBar (in `/Applications`)
- **Agent pointer:** This repo’s `AGENTS.MD` points to `~/Projects/agent-scripts/AGENTS.MD`

**Not installed (by design):** OpenClaw  
**VibeTunnel:** `npm install -g vibetunnel` failed (native build; needs Xcode Command Line Tools / build env). You can retry later with build tools installed.

---

## macOS permissions to grant

Grant these in **System Settings → Privacy & Security** (and Automation where listed). They can’t be set from the command line.

1. **Full Disk Access**  
   - Needed for: **imsg** (read `chat.db`), **CodexBar** (read Safari cookies).  
   - Add: Terminal (and/or Cursor), CodexBar.

2. **Accessibility**  
   - Needed for: **Peekaboo** (GUI automation).

3. **Screen Recording**  
   - Needed for: **Peekaboo** (screenshots / “see” commands).  
   - Path: **Screen & System Audio Recording** — enable for Terminal/Cursor as needed.

4. **Automation**  
   - Needed for: **imsg** to control the Messages app.  
   - In **Automation**, allow the app you use (e.g. Terminal/Cursor) to control **Messages**.

---

## Environment variables (API keys)

In your shell profile (`~/.zshrc` or `~/.bash_profile`), set:

- `OPENAI_API_KEY`
- `ANTHROPIC_API_KEY`
- `GEMINI_API_KEY`
- `X_AI_API_KEY`

**Check if they’re set:**

```bash
echo $OPENAI_API_KEY
echo $ANTHROPIC_API_KEY
echo $GEMINI_API_KEY
echo $X_AI_API_KEY
```

If any are empty, add lines like:

```bash
export OPENAI_API_KEY="your-key-here"
export ANTHROPIC_API_KEY="your-key-here"
# ... etc
```

Then run `source ~/.zshrc` (or your profile) or open a new terminal.

---

## gogcli (Google Workspace)

When you have the OAuth JSON file, run:

```bash
gog auth credentials <path-to-json>
```

(Install with `brew install steipete/tap/gogcli` if needed.)

---

## MCPorter

Run `mcporter list` to confirm it sees your MCP configs (e.g. under `~/.cursor/mcp.json`). Servers may show “offline” until their processes are running; that’s expected.
