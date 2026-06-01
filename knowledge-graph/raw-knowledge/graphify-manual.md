# Graphify & Understand Tools — Usage Manual

This guide describes the custom, token-efficient knowledge graph and analysis toolchain developed for this project.

---

## 🚀 3-Phase Query Strategy (AI Agent & CLI)

To keep context window clean and avoid heavy token consumption, **DO NOT** read `GRAPH_REPORT.md` directly. Instead, follow the optimized **3-Phase Query Flow**:

```
[ Phase 0: overview ] ───> [ Phase 1: sketch ] ───> [ Phase 2: detail ]
(Architecture Pillars)      (Targeted BFS Subgraph)   (Detailed Callers/Callees)
     <200 tokens                ~500 tokens             Full contextual info
```

### Phase 0: Overview (General Structure)
Get project statistics, top communities, and primary architecture pillars (God Nodes).
* **CLI Command**: `/overview`
* **MCP Tool**: `overview()`

### Phase 1: Sketch (Targeted Search)
Perform BFS traversal starting from top matches for a specific query (depth-limited, default: 3).
* **CLI Command**: `/sketch --query "<natural language query>"`
* **MCP Tool**: `sketch(query: string, depth: int)`

### Phase 2: Detail (Context Isolation)
Fetch callers, callees, file type, and community info for specific nodes returned by Phase 1.
* **CLI Command**: `/detail --ids "<id1,id2,...>"`
* **MCP Tool**: `detail(ids: string)`

---

## 🔍 Core CLI Command Guide

All tools include **Interactive Terminal Fallbacks** to prevent process freezing when running manual commands.

### 📁 Setup & Lifecycle
* `/awake`                          : Initialize current session memory, latest status, and Phase 0 overview.
* `/sync [-s "Summary text"]`       : Sync graph database, save memory checkpoint, and open graph visualizer.
* `/forget [pattern] [-y]`          : Delete specific past session summaries and matching graph database memories.
* `/map`                            : Open interactive web visualizer in browser.

### ⚡ Compilation & Clustering
* `/graphify-update`                : Incremental AST update of recently changed code files.
* `/graphify-cluster-only`          : Re-cluster community structure without scanning files.
* `/understand --start`            : Full understand-anything AST parse and Louvain clustering run.
* `/understand --update`           : Incremental understand-anything update.

### 🎯 Workspace Focus & Ripple Effects
* `/focus <file_path> <search_term>` : Locate a term in a file and view context lines around it.
* `/understand --diff`             : Check uncommitted git changes and estimate blast radius/impact.
* `/understand --explain <file>`    : Deep-dive explain a specific file's purpose.
