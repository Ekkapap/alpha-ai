---
name: graphify_build
description: "The full pipeline to turn a folder of files into a knowledge graph. Includes build, update, watch, and ingestion."
trigger: /graphify
---

# graphify_build

The full pipeline to turn a folder of files into a knowledge graph.

## Usage
- `/graphify <path>`: Build full graph
- `/graphify <path> --update`: Incremental update
- `/graphify <path> --watch`: Watch for changes
- `/graphify add <url>`: Add new source

## Step 0 - Clone (if GitHub URL)
```bash
LOCAL_PATH=$(graphify clone <github-url> [--branch <branch>])
```

## Step 1 - Ensure Installation
```bash
# (Interpreter resolution logic...)
mkdir -p graphify-out
python3 -c "import sys; open('graphify-out/.graphify_python', 'w').write(sys.executable)"
```

## Step 2 - Detect
```bash
$(cat graphify-out/.graphify_python) -c "import json; from graphify.detect import detect; from pathlib import Path; print(json.dumps(detect(Path('INPUT_PATH'))))" > graphify-out/.graphify_detect.json
```

## Step 3 - Extract (AST & Semantic)
- **Part A (AST)**: For code files.
- **Part B (Semantic)**: Parallel subagents for docs/images/papers.
  - Use `subagent_type="general-purpose"`.
  - Use Agent tool for chunks of 20-25 files.

## Step 4 - Build & Cluster
```bash
$(cat graphify-out/.graphify_python) -c "from graphify.build import build_from_json; from graphify.cluster import cluster; ... to_json(G, communities, 'graphify-out/graph.json')"
```

## Step 6 - Exports
- `graph.html`: Interactive visualization.
- `graph.svg`: Embedding.
- `cypher.txt`: Neo4j.

## Incremental Update (--update)
- Detect changed files.
- Re-extract and merge with `graph.json`.

## Background Watcher (--watch)
```bash
python3 -m graphify.watch INPUT_PATH --debounce 3
```

## Git Hook
```bash
graphify hook install
```

## Claude.md Integration
```bash
graphify claude install
```

## Honesty Rules
- Never invent an edge.
- Show token cost.
- Show raw cohesion scores.
- Warn if >5,000 nodes for HTML.
