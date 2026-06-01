Call `mcp__ALPHA__sync` with `summary: "$ARGUMENTS"` to save session context.

Go handles automatically (no agent action needed):
1. Run graphify + understand incremental update
2. Write `α/knowledge-graph/memories/session-[timestamp].md` (archive)
3. Append + dedup + trim `α/knowledge-graph/memories/session-summary.md`
4. Append to `session-summary.md` + dedup same-day entries + trim to 50 entries max

Response is minimal: `Session saved: session-YYYYMMDD-HHMM.md | session-summary.md: N entries`

Optional — only when semantic cleanup is needed (outdated topics across different wording):
Call `mcp__ALPHA__update_session_summary` with the rewritten full content.
To read current content first: read `α/knowledge-graph/memories/session-summary.md`
