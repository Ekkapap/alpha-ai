# Alpha Knowledge Graph Management

`/alpha-knowledge-graph` หรือ `alpha --knowledge-graph <subcommand>` ใช้จัดการ Docker services และ graph data

**เมื่อไหรควรใช้**: ต้องการรู้ว่า path/file นั้นคืออะไร เกี่ยวข้องกับอะไร มีความหมายยังไงในโปรเจ็ค — ใช้ `update` แล้วตาม 3-phase query flow

## Subcommands

**`start`** — เริ่ม dashboard + understand services (nginx + understand-server)
```bash
alpha --knowledge-graph start
```

**`stop`** — หยุด dashboard services
```bash
alpha --knowledge-graph stop
```

**`restart`** — stop แล้ว start ใหม่

**`status`** — ดูสถานะ containers ทั้งหมด (ชื่อ, state, port)
```bash
alpha --knowledge-graph status
```

**`logs [-f] [--grep <pattern>]`** — ดู logs ของ services
```bash
alpha --knowledge-graph logs                  # logs ล่าสุด
alpha --knowledge-graph logs -f               # follow real-time
alpha --knowledge-graph logs --grep "error"   # filter pattern
```

**`update [path]`** — rebuild knowledge graph (graphify scan + understand update)
```bash
alpha --knowledge-graph update          # scan จาก project root
alpha --knowledge-graph update ./src    # scan เฉพาะ path นั้น
```
ใช้เมื่อ: เพิ่ม/แก้ไข code แล้วต้องการให้ graph อัปเดต

**`init [path] [--force]`** — initialize graph ครั้งแรก
```bash
alpha --knowledge-graph init            # init จาก project root
alpha --knowledge-graph init --force    # ล้างแล้ว rebuild ใหม่
```
ถ้า graph.json มีอยู่แล้วจะแจ้งให้ใช้ `update` แทน

## Workflow: ต้องการเข้าใจ path หรือ file

เมื่อเจอ path/file ที่ไม่รู้ว่าคืออะไร เกี่ยวข้องกับอะไร ทำตามนี้:

```
1. /alpha-knowledge-graph update
   → rebuild graph ให้เป็นปัจจุบัน

2. mcp__ALPHA__sketch  query: "ชื่อไฟล์หรือ function"
   → BFS หา nodes ที่เกี่ยวข้อง ดูว่าถูกเรียกจากไหน เรียกอะไร

3. mcp__ALPHA__detail  ids: "<node-id-จาก-sketch>"
   → ดู callers, callees, file location แบบละเอียด

4. /alpha-focus <path> <term>
   → อ่านไฟล์จากจุดที่ term ปรากฏ (efficient กว่าอ่านทั้งไฟล์)
```

## MCP Tool ที่ใช้ร่วม

- `mcp__ALPHA__sketch` — Phase 1: BFS จาก query term, ได้ nodes + neighbors
- `mcp__ALPHA__detail` — Phase 2: callers/callees/location ของ node IDs ที่เจาะจง
- `mcp__ALPHA__focus` — อ่านไฟล์จาก keyword ที่ระบุ (args: path, term)
- `mcp__ALPHA__overview` — ดู god nodes + communities ทั้งหมด (<200 tokens)
