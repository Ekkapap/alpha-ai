# 🗂️ ALPHA SYSTEM: โครงสร้างโฟลเดอร์และองค์ประกอบ (System Structure)

เอกสารนี้ระบุรายละเอียดโครงสร้างของโฟลเดอร์และไฟล์ทั้งหมดภายในระบบ **Alpha (α)** ซึ่งตั้งอยู่ใน `[PROJECT_ROOT]/α/` เพื่อให้ AI และผู้พัฒนาเข้าใจขอบเขตการทำงานของแต่ละส่วนได้อย่างรวดเร็ว

---

## 📂 แผนผังโครงสร้างของระบบ (System Directory Tree)

```text
/cockpit-new/α/
├── alpha.json                  # ไฟล์คอนฟิกหลักของระบบ Alpha (Version, Tools, MCP)
├── README.md                   # คู่มือแนะนำเป้าหมายและการทำงานเชิงสถาปัตยกรรม (Architecture Design)
├── STRUCTURE.md                # เอกสารระบุโครงสร้างโฟลเดอร์และองค์ประกอบระบบ (ไฟล์นี้)
│
├── bin/                        # ไบนารีหลักของ Alpha แยกตามแพลตฟอร์ม (darwin/linux/windows)
│
├── commands/                   # ข้อกำหนดและเอกสารคำสั่ง (Command Spec)
│   └── awake.md                # เอกสารการเรียกใช้งานและคืนสภาพ Context ของคำสั่ง awake
│
├── hooks/                      # สคริปต์สั่งการเชิงรุก (Executable Hooks & CLI Command Wrappers)
│   └── bin/                    # ชุดสคริปต์ห่อหุ้มคำสั่ง (Mac/Linux Shell & Windows .cmd)
│       ├── awake / awake.cmd
│       ├── detail / detail.cmd
│       ├── focus / focus.cmd
│       ├── forget / forget.cmd
│       ├── graphify / graphify.cmd
│       ├── graphify-cluster-only / graphify-cluster-only.cmd
│       ├── graphify-update / graphify-update.cmd
│       ├── map / map.cmd
│       ├── alpha               # ตัวจัดการหลัก (Command Dispatcher ไปยัง my-graphify Go Binary)
│       ├── overview / overview.cmd
│       ├── sketch / sketch.cmd
│       ├── sync / sync.cmd
│       └── understand / understand.cmd
│
├── knowledge/                  # คลังความรู้ส่วนบุคคล แอดดัปเตอร์ภายนอก และเอกสารวิเคราะห์
│
├── memories/                   # ประวัติความจำและการบันทึกสภาวะสถานะล่าสุด
│   └── latest_state.md         # บันทึกความจำและข้อสรุปล่าสุดจากการซิงค์ระบบ
│
├── rules/                      # กฎสากลควบคุมการทำงานของ AI Agent ในระบบ Alpha
│   ├── execution.md            # กฎควบคุมความรวดเร็วในการเขียนโค้ดและดำเนินการ
│   ├── graphify.md             # กฎสามระยะ (3-Phase Query) เพื่อประหยัด Token สูงสุด
│   └── skill.md                # กฎควบคุมการโหลด Skill เมื่อต้องการเท่านั้น (On-demand)
│
├── scripts/                    # สคริปต์สำหรับผู้พัฒนา เพื่อติดตั้งและจัดระเบียบระบบ
│   ├── graphify.sh             # สคริปต์ตั้งค่าระบบความสัมพันธ์ Graphify เบื้องต้น
│   ├── setup-hooks.cmd         # สคริปต์ติดตั้ง Command Wrappers สำหรับ Windows
│   └── setup-hooks.sh          # สคริปต์ติดตั้ง Command Wrappers สำหรับ macOS & Linux
│
├── skills/                     # โมดูลทักษะเฉพาะด้านที่ได้รับการปรับปรุงให้มีขนาดเล็ก (Minimal SKILL Index)
│   ├── SKILL.md                # ดัชนีทักษะหลัก (Master Skill Index) ขนาดเบาพิเศษ
│   ├── auth/                   # ทักษะเกี่ยวกับระบบรักษาความปลอดภัยและการยืนยันตัวตน
│   ├── data/                   # ทักษะการจัดการฐานข้อมูล (Prisma, ElectricSQL)
│   ├── framework/              # ทักษะการควบคุมเฟรมเวิร์กหลัก (Next.js, Golang)
│   ├── graphify/               # ทักษะการจัดการและสืบค้นกราฟความสัมพันธ์
│   ├── runtime/                # ทักษะเกี่ยวกับรันไทม์ (Bun JavaScript/TypeScript)
│   ├── state/                  # ทักษะการจัดการสถานภาพตัวแปร (TanStack State)
│   ├── ui/                     # ทักษะการสร้างและตกแต่ง UI (Tailwind CSS, Mantine)
│   └── ux/                     # ทักษะการออกแบบและพัฒนาประสบการณ์ผู้ใช้ขั้นสูง
│
└── tools/                      # โค้ดต้นฉบับ (Golang) และไบนารีระบบย่อย
    ├── bin/                    # พื้นที่เก็บไบนารีที่คอมไพล์แล้วแยกตามแพลตฟอร์ม
    │   ├── darwin/             # macOS Binaries (my-graphify, my-understand)
    │   ├── linux/              # Linux Binaries (my-graphify, my-understand)
    │   └── windows/            # Windows Binaries (my-graphify.exe, my-understand.exe)
    ├── my-graphify/            # ระบบย่อย Golang สำหรับจัดการ Graph และควบคุมหน่วยความจำ
    │   ├── go.mod / go.sum
    │   └── main.go             # โค้ดส่วนหลักที่จัดการ CLI Argument และพฤติกรรมของ `alpha`
    └── my-understand/          # ระบบย่อย Golang สำหรับวิเคราะห์โครงสร้างโค้ดเชิงลึก (AST Depth Analysis)
        ├── go.mod / go.sum
        └── main.go             # โค้ดส่วนหลักของตัวตรวจสอบ AST และ Call Graph
```

---

## 🧩 รายละเอียดโมดูลและการทำงาน (Detailed File & Module Specifications)

### 1. ไฟล์การกำหนดค่าหลัก: `alpha.json`
- **หน้าที่**: ควบคุมความปลอดภัยและบอกตำแหน่งเครื่องมือ (Tools Mapping) ให้กับ AI และ Go Orchestrator
- **โครงสร้างข้อมูลหลัก**:
  - `version`: เวอร์ชั่นของระบบ
  - `tools`: แผนผังบอกตำแหน่งซอร์สโค้ดและโฟลเดอร์เป้าหมายไบนารีตามสถาปัตยกรรม CPU/OS
  - `mcp`: รายชื่อระบบ MCP ที่ได้รับการปรับปรุงให้ใช้งานร่วมกันได้ (`MY_GRAPHIFY`, `MY_UNDERSTAND`)
  - `memories`: ตำแหน่งเก็บข้อมูลแคชสภาวะของ AI

### 2. ชุดโมดูลความจำ: `/memories/`
- **latest_state.md**: ทำหน้าที่เหมือน "RAM ความทรงจำระยะยาว" ของ AI ในระดับโปรเจกต์ ซึ่งตัว Go Binary จะคัดกรองข้อมูลสำคัญที่ AI สรุปไว้มาอัปเดตแบบย้อนกลับ (Backward Update) ทุกครั้งที่รัน `alpha --memo-sync`
- **เป้าหมายเชิง Token**: หลีกเลี่ยงการจัดเก็บประวัติการสนทนาทั้งหมด และจัดเก็บเฉพาะข้อสรุปการตัดสินใจในอดีตเท่านั้น ทำให้เวลา AI ตื่นขึ้นมาด้วย `alpha --awake` จะประหยัด Token ได้มากกว่า 80% เมื่อเทียบกับการอ่านแชทประวัติศาสตร์ทั้งหมด

### 3. ชั้นประมวลผลคำสั่งเชิงรุก: `/hooks/bin/`
- **บทบาท**: เป็น Command Interface ที่คอยแปลงรูปและดักจับการพิมพ์คำสั่งใน Terminal
- **กลไกการส่งทอด (Delegation Flow)**:
  1. เมื่อผู้พัฒนาหรือ AI รันสคริปต์ `α/hooks/bin/alpha --awake` (หรือผ่าน CLI Command/Slash command)
  2. ตัวสคริปต์ตรวจสอบระบบปฏิบัติการ (Darwin / Linux / Windows)
  3. สกริปต์ส่งคำสั่งและพารามิเตอร์ทั้งหมดไปยังไบนารีที่ถูกคอมไพล์สำเร็จแล้ว เช่น `α/tools/bin/{platform}/my-graphify alpha awake`
  4. ไบนารี Go ประมวลผลดึงข้อมูลความจำและกราฟความสัมพันธ์แบบ Local และคืนค่า (Return) ข้อมูลที่ถูกแปลงให้กระชับและไร้ขยะ Token สู่ AI ทันที

### 4. กฎสากลสำหรับ AI: `/rules/`
- ทำงานสอดประสานกับระบบความปลอดภัยและหลักความประหยัด Token โดยที่:
  - `execution.md`: บังคับให้ AI ลงมือเขียนโค้ดและทดสอบทันทีเพื่อลดปริมาณการบ่นหรือการวางแผนที่ซับซ้อนเกินจำเป็นในไฟล์แชท
  - `graphify.md`: บังคับให้ใช้ **3-Phase Query Flow** เพื่อให้ค้นหาจุดเชื่อมโยง (Seed) ก่อนดึงรายละเอียด
  - `skill.md`: บังคับให้ AI ไม่โหลดสคริปต์/คู่มือที่ไม่มีความจำเป็นเข้ามาในบริบท (Context) ช่วยตัดเสียงรบกวน (Noise) ลงอย่างมหาศาล

### 5. ซอร์สโค้ดของเครื่องมือ: `/tools/`
- **my-graphify**: ตัวขับเคลื่อนระบบกราฟหลัก แปลงข้อมูล `graph.json` ออกมาเป็นภาพรวม, BFS เครือข่ายใกล้เคียง และข้อมูลความคุ้นเคยของ Module
- **my-understand**: วิเคราะห์เชิงลึกวิถีการเรียกใช้งาน (Call Graph) จากผลลัพธ์ของระบบ AST
