---
name: graph_stats
description: "Return summary statistics: node count, edge count, communities, confidence breakdown, and cost tracking."
---

# graph_stats

Return summary statistics: node count, edge count, communities, confidence breakdown.

## General Stats
```bash
$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
from networkx.readwrite import json_graph

data = json.loads(Path('graphify-out/graph.json').read_text())
analysis = json.loads(Path('graphify-out/.graphify_analysis.json').read_text())

print(f'Nodes: {len(data["nodes"])}')
print(f'Edges: {len(data["links"])}')
print(f'Communities: {len(analysis.get("communities", {}))}')
"
```

## Step 8 - Token reduction benchmark
If `total_words` from `graphify-out/.graphify_detect.json` is greater than 5,000, run:

```bash
$(cat graphify-out/.graphify_python) -c "
import json
from graphify.benchmark import run_benchmark, print_benchmark
from pathlib import Path

detection = json.loads(Path('graphify-out/.graphify_detect.json').read_text())
result = run_benchmark('graphify-out/graph.json', corpus_words=detection['total_words'])
print_benchmark(result)
"
```

## Cost Tracking (Step 9)
```bash
$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
cost_path = Path('graphify-out/cost.json')
if cost_path.exists():
    cost = json.loads(cost_path.read_text())
    print(f'Total Input Tokens: {cost["total_input_tokens"]:,}')
    print(f'Total Output Tokens: {cost["total_output_tokens"]:,}')
    print(f'Total Runs: {len(cost["runs"])}')
else:
    print('No cost data found.')
"
```
