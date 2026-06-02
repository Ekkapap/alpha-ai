ทดสอบ G10 แล้ว flow ดูไม่ค่อยถูกต้อง นี่คือ flow ที่ควรจะเป็น

1. docker ps ดูครับ มีส่วนเกิน ลบ container สว่นเกินออก รวมทั้ง  script ที่สร้างมันขึ้นมาจาก install.sh
2. logic ไม่ถูกต้อง ใน project ไม่ควรมี graphify-out -> α/knowledge-graph/graphify-out, α/knowledge-graph/*, ยกเว้น memories ให้คงไว้ใน project นั้นๆของใครของมัน , 
  เอาง่ายๆ สุด ก็ ไม่ต้องเอาอะไรมาลงในนี้เลยครับ ให้ลงไว้ที่ global  
  
# rtk tree -L 2 ~/.alpha.ai/
.
├── CONTRIBUTING.md
├── LICENSE
├── README.md
├── agents-resource
│   ├── PRODUCTION_CLAUDE.md
│   ├── PRODUCTION_GEMINI.md
│   ├── commands
│   ├── config.json
│   ├── hooks
│   ├── rules
│   ├── skills
│   ├── tools
│   └── workflows
├── alpha.json
├── config.json
├── docker
│   ├── Dockerfile.alpha
│   ├── Dockerfile.graphify
│   ├── Dockerfile.understand
│   ├── dashboard.html
│   ├── dashboard.nginx.conf
│   └── understand-start.sh
├── docker-compose.global.yml
├── docker-compose.yml
├── docs
│   ├── SESSION_CONTEXT.md
│   ├── TASK.md
│   ├── USER_REQUIREMENT.md
│   ├── archived
│   ├── ctx_arch.md
│   ├── ctx_docker.md
│   └── ctx_tasks.md
├── graphify-out -> knowledge-graph/graphify-out
├── knowledge-graph
│   ├── graphify-out
│   ├── memories
│   ├── raw-knowledge
│   └── understand-anything
├── logs
│   └── nginx
└── scripts
    ├── dashboard.sh
    ├── graphify.sh
    ├── install.sh
    ├── setup-hooks.cmd
    └── setup-hooks.sh

Create symlink
3. α/ -> ~/.alpha.ai/
4. .claude/* -> ~/.alpha.ai/[skills|commands|hooks|rules]
5. .agents/* -> ~/.alpha.ai/[skills|tools|workflows]
6. graphify-out/ -> ~/.alpha.ai/knowledge-graph/graphify-out/[project-id]
7. .understand-anything/ -> ~/.alpha.ai/knowledge-graph/understand-anything/[project-id]
8. memories/*.md (Self Memories)
9. CLAUDE.md -> ~/.alpha.ai/agents-resource/PRODUCTION_CLAUDE.md
10. GEMINI.md -> ~/.alpha.ai/agents-resource/PRODUCTION_GEMINI.md
11. .mcp.json (ไม่แน่ใจว่าหากใช้จาก global ที่เดียวเป็น mcp global เลยได้ไหมไม่ต้องติดตั้งแยก project)
12. .understandignore (cp from ~/.alpha.ai/knowledge-graph/understand-anything/.understandignore)
13. .graphifyignore (cp from ~/.alpha.ai/knowledge-graph/graphify-out/.graphifyignore)
14. .gitignore (cp from ~/.alpha.ai/.gitignore)