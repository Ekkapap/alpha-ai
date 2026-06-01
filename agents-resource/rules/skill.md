---
trigger: model_decision
description: Apply at the start of a session.
---

# CONTEXT LOADING PROTOCOL

> Skills are in `<agent-home>/skills/`  e.g. `.claude/` · `.agents/`

- DO NOT load supplementary context by default
- Load only when the user specifies it or the task clearly requires it
- If unsure → DO NOT load
- Load specific files only — no directory scanning
- Never mix parent and child context at the same time
