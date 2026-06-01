---
trigger: model_decision
description: Apply whenever running shell commands.
---

# RTK — TOKEN OPTIMIZATION

RTK filters and compresses shell command output before it reaches the AI context, saving 60–90% tokens on common operations.

## Rule

Always prefix shell commands with `rtk` to minimize token consumption:

```bash
rtk git status
rtk git diff
rtk ls src/
rtk grep "pattern" src/
rtk find "*.ts" .
rtk cargo test
rtk docker ps
rtk gh pr list
```

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
