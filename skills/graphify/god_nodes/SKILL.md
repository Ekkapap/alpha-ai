---
name: god_nodes
description: "Return the most connected nodes - the core abstractions of the knowledge graph."
---

# god_nodes

Return the most connected nodes - the core abstractions of the knowledge graph.

```bash
$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
analysis = json.loads(Path('graphify-out/.graphify_analysis.json').read_text())
gods = analysis.get('gods', [])
print('God Nodes (most connected):')
for i, god in enumerate(gods[:10]):
    print(f'{i+1}. {god["label"]} ({god["degree"]} edges)')
"
```
