---
name: skills-index
description: Master skill index. Load on-demand only — never auto-load. Match task to skill path, then fetch that specific SKILL.md.
---

[no-skill] 
**DO NOT USE SKILL AT SESSION START**

# 🗺️ SKILLS_MANIFEST (v3.4.11)
[Mode: ON-DEMAND | Status: INDEX-ONLY]

## 🚨 CRITICAL & MANDATORY RULE
DO NOT scan or auto-read files or include any skill content into your active context.
DO NOT scan or load any files in the `.agents/skills/` directory automatically. 
Refer to this index only. Fetch full content ONLY when a task explicitly is assigned or requires it.

## 🗂️ SKILL INDEX MAP
reference-path = ./

- **auth**: [better-auth, authentication] → `skills/auth/SKILL.md`
- **data**: [prisma-orm, electric-sql, tanstack-db] → `skills/data/SKILL.md`
- **framework**: [next-js, golang-expert] → `skills/framework/SKILL.md`
- **runtime**: [bun-runtime, bun] → `skills/runtime/SKILL.md`
- **state**: [tanstack-state, state-management] → `skills/state/SKILL.md`
- **ui**: [designer, mantine-ui, screenshot-to-code, tailwind-css] → `skills/ui/SKILL.md`
- **ux**: [advanced-ux, user-experience] → `skills/ux/SKILL.md`
- **graphify**: [graphify-core, knowledge-graph] → `skills/graphify/SKILL.md`

## 🎯 EXECUTION_LOGIC
1. Identify User Intent.
2. Match with **Skill Name** or **TAG**.
3. Call `rtk read` on the specific Path.
4. Wipe content from memory after task completion.