# Graph Report - workspace  (2026-06-01)

## Corpus Check
- 32 files · ~15,873 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 171 nodes · 328 edges · 10 communities
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 56 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]

## God Nodes (most connected - your core abstractions)
1. `main()` - 22 edges
2. `assembleGraph()` - 16 edges
3. `finalizeGraph()` - 11 edges
4. `main()` - 11 edges
5. `loadFullGraph()` - 8 edges
6. `updatePipeline()` - 8 edges
7. `main()` - 7 edges
8. `cliAwake()` - 7 edges
9. `sketchGraph()` - 7 edges
10. `graphOverview()` - 7 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `cliAlpha()`  [INFERRED]
  α/agents-resource/tools/graphify/main.go → α/agents-resource/tools/graphify/cmd_alpha.go
- `main()` --calls--> `registerMCPAlpha()`  [INFERRED]
  α/agents-resource/tools/graphify/main.go → α/agents-resource/tools/graphify/cmd_alpha.go
- `cliAwake()` --calls--> `getGodNodes()`  [INFERRED]
  α/agents-resource/tools/graphify/cmd_awake.go → α/agents-resource/tools/graphify/graph.go
- `cliAwake()` --calls--> `getGraphStats()`  [INFERRED]
  α/agents-resource/tools/graphify/cmd_awake.go → α/agents-resource/tools/graphify/graph.go
- `main()` --calls--> `cliBuild()`  [INFERRED]
  α/agents-resource/tools/graphify/main.go → α/agents-resource/tools/graphify/cmd_build.go

## Import Cycles
- None detected.

## Communities (10 total, 0 thin omitted)

### Community 1 - "Community 1"
Cohesion: 0.27
Nodes (11): handleKnowledgeGraph(), runGraphifyUpdate(), runGraphifyUpdateTarget(), binPath(), findRoots(), inDocker(), main(), runTool() (+3 more)

### Community 8 - "Community 8"
Cohesion: 0.13
Nodes (26): GodNode, analysisJSON, BuildResult, assembleGraph(), buildGraph(), cliUpdate(), dirPrefix(), filterEmpty() (+18 more)

### Community 10 - "Community 10"
Cohesion: 0.40
Nodes (9): appendAndMaintainSummary(), cliSync(), countSummaryEntries(), extractEntryDate(), registerMCPSync(), registerMCPUpdateSessionSummary(), splitSummaryEntries(), generateProjectSummary() (+1 more)

### Community 12 - "Community 12"
Cohesion: 0.21
Nodes (8): alphaConfig, cliAlpha(), registerMCPAlpha(), alphaDisplay(), alphaReadyChecks(), loadAlphaConfig(), readyCheck, MCPServer

### Community 14 - "Community 14"
Cohesion: 0.31
Nodes (7): extractGeneric(), isKeyword(), lastImportSegment(), lastPathSegment(), langRules, Regexp, ExtractedFile

### Community 15 - "Community 15"
Cohesion: 0.33
Nodes (8): CallExpr, Expr, calleeLabel(), extractGo(), fileNodeLabel(), lastSegment(), receiverType(), ExtractedFile

### Community 16 - "Community 16"
Cohesion: 0.11
Nodes (29): cliAwake(), registerMCPAwake(), cliBuild(), cliFocus(), registerMCPDebugInfo(), registerMCPFocus(), cliForget(), cliDetail() (+21 more)

### Community 18 - "Community 18"
Cohesion: 0.29
Nodes (10): ExtractedFile, RawEdge, RawNode, extractFile(), globMatch(), loadIgnorePatterns(), scanProject(), shouldIgnore() (+2 more)

### Community 19 - "Community 19"
Cohesion: 0.53
Nodes (8): APIKey, callAnthropic(), callGemini(), callOpenAI(), GenerateCommunityLabels(), LoadAPIKey(), parseJSONLabels(), postJSON()

### Community 20 - "Community 20"
Cohesion: 0.15
Nodes (30): GraphEdge, GraphNode, KnowledgeGraph, Layer, cleanComment(), countFiles(), detectFrameworks(), detectLanguages() (+22 more)

## Knowledge Gaps
- **16 isolated node(s):** `CallToolResult`, `CallToolRequest`, `MCPServer`, `MCPServer`, `GodNode` (+11 more)
  These have ≤1 connection - possible missing edges or undocumented components.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `Community 16` to `Community 8`, `Community 10`, `Community 12`?**
  _High betweenness centrality (0.342) - this node is a cross-community bridge._
- **Why does `assembleGraph()` connect `Community 8` to `Community 19`, `Community 14`?**
  _High betweenness centrality (0.177) - this node is a cross-community bridge._
- **Why does `runFullBuild()` connect `Community 8` to `Community 18`?**
  _High betweenness centrality (0.129) - this node is a cross-community bridge._
- **Are the 21 inferred relationships involving `main()` (e.g. with `cliAlpha()` and `registerMCPAlpha()`) actually correct?**
  _`main()` has 21 INFERRED edges - model-reasoned connections that need verification._
- **Are the 6 inferred relationships involving `assembleGraph()` (e.g. with `AssignCommunities()` and `FindGodNodes()`) actually correct?**
  _`assembleGraph()` has 6 INFERRED edges - model-reasoned connections that need verification._
- **Are the 6 inferred relationships involving `loadFullGraph()` (e.g. with `cliAwake()` and `registerMCPAwake()`) actually correct?**
  _`loadFullGraph()` has 6 INFERRED edges - model-reasoned connections that need verification._
- **What connects `CallToolResult`, `CallToolRequest`, `MCPServer` to the rest of the system?**
  _16 weakly-connected nodes found - possible documentation gaps or missing edges._