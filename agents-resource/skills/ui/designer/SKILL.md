---
name: ui-designer
description: Extract design systems from reference UI images and generate implementation-ready UI design prompts. Use when users provide UI screenshots/mockups and want to create consistent designs, generate design systems, or build MVP UIs matching reference aesthetics.
---

This skill implements a 6-step workflow: gather reference images and project idea → extract design system (colors, typography, spacing, components) via subagent → generate MVP PRD interactively → compose final implementation prompt combining design system and PRD → verify React environment → implement UI with multiple variations (3 mobile, 2 web).

## References

| File | Purpose |
|------|---------|
| references/workflow.md | Full step-by-step workflow, template asset descriptions, best practices, example |
| assets/design-system.md | Template for extracting visual design patterns from images |
| assets/app-overview-generator.md | Template for collaborative PRD generation |
| assets/vibe-design-template.md | Final implementation prompt template |
