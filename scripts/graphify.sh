#!/bin/bash
# .agents/scripts/graphify.sh
# Core logic for Graphify Pipeline
# Relies on $PROJECT_ROOT being provided by the caller (hooks/bin or proxy)

TARGET_DIR="${1:-.}"
MODE="$2"

if [[ -n "$MODE" && "$MODE" != "--update" && "$MODE" != "--cluster-only" ]]; then
    if [[ "$MODE" == "--help" || "$MODE" == "-h" || "$MODE" == "--h" ]]; then
        echo "Graphify Pipeline Wrapper"
        echo ""
        echo "Usage:"
        echo "  /graphify [flag]"
        echo ""
        echo "Flags:"
        echo "  --update        Update the knowledge graph"
        echo "  --cluster-only  Update community clusters only"
        echo "  --help, -h      Show this help message"
        exit 0
    else
        echo "Error: Unrecognized flag '$MODE'"
        echo "Run '/graphify --help' to see usage."
        exit 1
    fi
fi

# 1. Resolve Python Interpreter
if [ -f "$PROJECT_ROOT/α/knowledge-graph/graphify-out/.graphify_python" ]; then
    PY_PATH=$(cat "$PROJECT_ROOT/α/knowledge-graph/graphify-out/.graphify_python")
else
    PY_PATH=$(which python3)
fi

echo "🚀 Running Graphify ($MODE) on $TARGET_DIR..."

# 2. Run Base Graphify Command
case "$MODE" in
    "--update")
        graphify update "$TARGET_DIR"
        ;;
    "--cluster-only")
        graphify cluster-only "$TARGET_DIR"
        ;;
    *)
        graphify update "$TARGET_DIR" --force
        ;;
esac

# 3. AI Post-Processing (Force Knowledge + Prune Docs)
$PY_PATH -c "
import json, os
from pathlib import Path
from networkx.readwrite import json_graph
from graphify.cluster import cluster
from graphify.export import to_json

graph_path = Path(os.environ['PROJECT_ROOT']) / 'α/knowledge-graph/graphify-out/graph.json'
if not graph_path.exists():
    exit(0)

data = json.loads(graph_path.read_text())
G = json_graph.node_link_graph(data, edges='links')

# A. Force Ingest Knowledge
knowledge_path = Path(os.environ['PROJECT_ROOT']) / 'α/knowledge-graph/raw-knowledge'
if knowledge_path.exists():
    if not any(n == 'wisdom_central_kb' for n in G.nodes()):
        G.add_node('wisdom_central_kb', label='Project Wisdom Base', file_type='concept')
    for f in knowledge_path.glob('*.md'):
        f_id = f'wisdom_{str(f).replace(\"/\", \"_\").replace(\".\", \"_\")}'
        if f_id not in G.nodes():
            G.add_node(f_id, label=f.name, file_type='document', source_file=str(f))
            G.add_edge(f_id, 'wisdom_central_kb', relation='belongs_to')

# B. Prune Ignored Docs (README, GEMINI, etc.)
targets = ['README.md', 'GEMINI.md', 'AGENTS.md', 'CLAUDE.md']
nodes_to_remove = [
    n for n, d in G.nodes(data=True) 
    if any(t in str(d.get('source_file', '')) for t in targets) or 
       any(t.lower().replace('.', '_') in str(n).lower() for t in targets)
]
G.remove_nodes_from(nodes_to_remove)

# C. Save & Export
communities = cluster(G)
to_json(G, communities, str(Path(os.environ['PROJECT_ROOT']) / 'α/knowledge-graph/graphify-out/graph.json'), force=True)
"

# 4. Final Export
graphify export html
echo "✅ Graph updated and exported to graphify-out/graph.html"
