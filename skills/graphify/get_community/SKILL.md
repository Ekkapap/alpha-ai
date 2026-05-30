---
name: get_community
description: "Explore communities, their cohesion, and members. Includes labeling and wiki generation."
---

# get_community

Explore communities, their cohesion, and members.

## List all communities
```bash
$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
analysis = json.loads(Path('graphify-out/.graphify_analysis.json').read_text())
labels = json.loads(Path('graphify-out/.graphify_labels.json').read_text()) if Path('graphify-out/.graphify_labels.json').exists() else {}
communities = analysis.get('communities', {})
cohesion = analysis.get('cohesion', {})

for cid, members in communities.items():
    label = labels.get(cid, f'Community {cid}')
    score = cohesion.get(cid, 0)
    print(f'ID {cid}: {label} (Cohesion: {score:.2f}, {len(members)} nodes)')
"
```

## Get members of a community
```bash
$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
analysis = json.loads(Path('graphify-out/.graphify_analysis.json').read_text())
communities = analysis.get('communities', {})
cid = 'COMMUNITY_ID'
if cid in communities:
    print(f'Members of Community {cid}:')
    for nid in communities[cid]:
        print(f'  - {nid}')
else:
    print(f'Community {cid} not found.')
"
```

## Step 5 - Label communities
Read `graphify-out/.graphify_analysis.json`. For each community key, look at its node labels and write a 2-5 word plain-language name (e.g. "Attention Mechanism", "Training Pipeline", "Data Loading").

Then regenerate the report and save the labels for the visualizer:

```bash
$(cat graphify-out/.graphify_python) -c "
import sys, json
from graphify.build import build_from_json
from graphify.cluster import score_all
from graphify.analyze import god_nodes, surprising_connections, suggest_questions
from graphify.report import generate
from pathlib import Path

extraction = json.loads(Path('graphify-out/.graphify_extract.json').read_text())
detection  = json.loads(Path('graphify-out/.graphify_detect.json').read_text())
analysis   = json.loads(Path('graphify-out/.graphify_analysis.json').read_text())

G = build_from_json(extraction)
communities = {int(k): v for k, v in analysis['communities'].items()}
cohesion = {int(k): v for k, v in analysis['cohesion'].items()}
tokens = {'input': extraction.get('input_tokens', 0), 'output': extraction.get('output_tokens', 0)}

# LABELS - replace these with the names you chose above
labels = LABELS_DICT

# Regenerate questions with real community labels (labels affect question phrasing)
questions = suggest_questions(G, communities, labels)

report = generate(G, communities, cohesion, labels, analysis['gods'], analysis['surprises'], detection, tokens, 'INPUT_PATH', suggested_questions=questions)
Path('graphify-out/GRAPH_REPORT.md').write_text(report)
Path('graphify-out/.graphify_labels.json').write_text(json.dumps({str(k): v for k, v in labels.items()}))
print('Report updated with community labels')
"
```

## Step 6b - Wiki
Only run this step if `--wiki` was explicitly given.

```bash
$(cat graphify-out/.graphify_python) -c "
import json
from graphify.build import build_from_json
from graphify.wiki import to_wiki
from graphify.analyze import god_nodes
from pathlib import Path

extraction = json.loads(Path('graphify-out/.graphify_extract.json').read_text())
analysis   = json.loads(Path('graphify-out/.graphify_analysis.json').read_text())
labels_raw = json.loads(Path('graphify-out/.graphify_labels.json').read_text()) if Path('graphify-out/.graphify_labels.json').exists() else {}

G = build_from_json(extraction)
communities = {int(k): v for k, v in analysis['communities'].items()}
cohesion = {int(k): v for k, v in analysis['cohesion'].items()}
labels = {int(k): v for k, v in labels_raw.items()}
gods = god_nodes(G)

n = to_wiki(G, communities, 'graphify-out/wiki', community_labels=labels or None, cohesion=cohesion, god_nodes_data=gods)
print(f'Wiki: {n} articles written to graphify-out/wiki/')
print('  graphify-out/wiki/index.md  ->  agent entry point')
"
```
