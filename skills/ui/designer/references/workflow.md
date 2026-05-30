# UI Designer Workflow

## Step 1: Gather Inputs

Request from user:
- **Reference images directory**: Path to folder containing UI screenshots/mockups
- **Project idea file**: Document describing the product concept and goals
- **Existing PRD** (optional): If PRD already exists, skip Step 3

## Step 2: Extract Design System from Images

**Use Task tool with general-purpose subagent**, providing the template from `assets/design-system.md`:
- Analyze color palettes (primary, secondary, accent, functional colors)
- Extract typography (font families, sizes, weights, line heights)
- Identify component styles (buttons, cards, inputs, icons)
- Document spacing system
- Note animations/transitions patterns
- Include dark mode variants if present

**Output**: Complete design system markdown  
**Save to**: `documents/designs/{image_dir_name}_design_system.md`

## Step 3: Generate MVP PRD (if not provided)

**Use Task tool with general-purpose subagent**, using template from `assets/app-overview-generator.md`:
- Replace `{项目背景}` with content from project idea file
- Guides through: elevator pitch, problem statement, target audience, USP, features list, UX/UI considerations

**Interact with user** to refine and clarify product requirements  
**Save as variable** for Step 4 (optionally save to `documents/prd/`)

## Step 4: Compose Final UI Implementation Prompt

Combine design system and PRD using `assets/vibe-design-template.md`:

**Substitutions:**
- `{项目设计指南}` → Design system from Step 2
- `{项目MVP PRD}` → PRD from Step 3 or provided PRD file

**Save to**: `documents/ux-design/{idea_file_name}_design_prompt_{timestamp}.md`

## Step 5: Verify React Environment

```bash
find . -name "package.json" -exec grep -l "react" {} \;
```

If none found:
```bash
npx create-react-app my-app
cd my-app
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
npm install lucide-react
```

## Step 6: Implement UI

Use the final composed prompt from Step 4. The prompt instructs to:
- Create multiple design variations (3 for mobile, 2 for web)
- Organize as separate components: `[solution-name]/pages/[page-name].jsx`
- Aggregate all variations in showcase page

## Template Assets

### assets/design-system.md
Template for extracting visual design patterns. Covers: color palette, typography, component styles, spacing system, animations, dark mode variants.

### assets/app-overview-generator.md
Template for collaborative PRD generation. Guides through: elevator pitch, problem statement, target audience, USP, platform targets, feature list, UX/UI considerations per screen.

### assets/vibe-design-template.md
Final implementation prompt template combining design system and PRD. Includes: aesthetic principles, Tailwind CSS + Lucide icons requirements, task specifications.

## Best Practices

- Read all images before starting analysis; look for patterns across multiple screens
- Use specific values (hex codes, px sizes) not generic descriptions
- Engage user interactively during PRD generation to clarify ambiguities
- Save all outputs to `documents/` directory for easy reference
- Save final prompt with timestamp for version tracking

## Example

**User provides:** `reference-images/saas-dashboard/` + `ideas/project-management-app.md`

1. Read 5 images → use design-system.md template → save `documents/designs/saas-dashboard_design_system.md`
2. Use app-overview-generator.md → refine PRD with user
3. Combine with vibe-design-template.md → save `documents/ux-design/project-management-app_design_prompt_20251025_153000.md`
4. Check React environment → implement UI using final prompt
