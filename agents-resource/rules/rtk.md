---
trigger: model_decision
description: Apply whenever running shell commands.
---

# RTK — TOKEN OPTIMIZATION

RTK filters and compresses shell command output before it reaches the AI context, saving 60–90% tokens on common operations.

## Rule

Always prefix shell commands with `rtk` to minimize token consumption:

```bash
rtk ls src/
rtk read "filepath"                # Read whole file
rtk read -n "filepath"             # Read with line numbers
rtk read -m <max-lines> "filepath" # Read up to N lines
rtk grep "pattern" src/
rtk find "*.ts" .
rtk cargo test
rtk docker ps
rtk git status
rtk git diff
rtk gh pr list
```

## Reading a Specific Section

To read starting from a known term/keyword (e.g. deep in a large file), use `alpha focus` or `/alpha-focus`:

```bash
alpha focus <path> <term> [max_lines]   # CLI
# or: /alpha-focus <path> <term>        # MCP slash command
```

`term` is searched in the file; returns content from that point up to `max_lines`.

## Meta Commands (run directly — no `rtk` prefix)

```bash
rtk gain              # Show cumulative token savings
rtk gain --history    # Command history with per-command savings
rtk discover          # Scan history for missed RTK opportunities
rtk proxy <cmd>       # Run raw command without filtering (debug only)
```

## Applies To

All agents running shell commands: Claude Code, Gemini CLI, Cursor, Antigravity, Cline, Codex, and any terminal-integrated AI tool.

> If `rtk` is not installed, run `α/install.sh` — RTK is installed in Step 1.
