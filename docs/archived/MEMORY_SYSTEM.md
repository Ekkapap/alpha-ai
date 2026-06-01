# Agent Memory Operating System (AMOS)

## Vision

สร้างระบบความทรงจำสำหรับ AI Agent ที่ทำงานกับ Codebase ขนาดใหญ่ โดยไม่พึ่ง LLM ในการตีความซ้ำทุกครั้ง และไม่ต้อง Scan/Rebuild ทั้งโปรเจกต์ตลอดเวลา

แนวคิดหลักคือแยก

* Structure
* Meaning
* Experience
* Knowledge
* History

ออกจากกันอย่างชัดเจน และใช้ MCP เป็นตัวกลางในการจัดการทั้งหมด

---

# เป้าหมาย

Agent ควรสามารถเข้าใจโปรเจกต์ได้จากข้อมูลที่ถูกสกัดไว้แล้ว โดยไม่จำเป็นต้องเปิด Source Code ทุกครั้ง

ลด

* Token Usage
* Context Window Consumption
* จำนวน Query Hop
* เวลาในการ Reasoning

เพิ่ม

* Context Quality
* Knowledge Reuse
* Historical Learning
* Cross-Agent Collaboration

---

# Architecture

Source Code

↓

Graphify

↓

Understand

↓

Meta Graph

↓

Muscle Memory

↓

Knowledge Vault

---

# Layer 1 : Graphify

## หน้าที่

สกัดโครงสร้างของระบบ

ตัวอย่าง

* Folder
* File
* Class
* Function
* Import
* Export
* Call Graph
* Dependency Graph

## ตัวอย่างข้อมูล

Function

* createUser()

รู้ว่า

* อยู่ไฟล์ไหน
* เรียกใครบ้าง
* ถูกเรียกโดยใครบ้าง

## จุดเด่น

ไม่ใช้ LLM

ข้อมูลมาจาก AST และ Static Analysis

---

# Layer 2 : Understand

## หน้าที่

สกัดความหมายของ Code

ตัวอย่าง

Function

createUser()

Meaning

* Create User Account
* Initialize Default Permissions

## ข้อมูลที่เก็บ

* Summary
* Domain
* Capability
* Responsibility
* Keywords

## จุดเด่น

ตอบคำถามเชิง Business Meaning ได้

---

# Layer 3 : Meta Graph

## หน้าที่

เป็น Graph กลางที่ Wrap Graphify และ Understand

Agent จะไม่อ่าน Graphify หรือ Understand โดยตรง

Agent อ่าน Meta Graph เท่านั้น

---

## หมวดหมู่ของ Node

* codebase
* folder
* file
* class
* function
* route
* api
* component
* concept
* knowledge
* bug
* solution
* decision

---

## ข้อมูลใน Node

เก็บเฉพาะข้อมูลระดับสรุป

ไม่เก็บรายละเอียดลึก

ตัวอย่าง

* name
* type
* domain
* capability
* tags
* importance
* refs

---

## Refs

เชื่อมไปยังระบบอื่น

ตัวอย่าง

* graphifyRef
* understandRef
* muscleRef
* knowledgeRef

---

## หลักการ

Node ต้องเข้าใจได้ในตัวเองระดับหนึ่ง

Agent ควรตอบคำถามได้ส่วนใหญ่โดยไม่ต้อง Follow Ref

แต่ยังสามารถ Follow Ref ได้เมื่อจำเป็น

---

# Layer 4 : Muscle Memory

## แนวคิด

เก็บประสบการณ์การทำงานจริง

ไม่ใช่แค่ความหมายของ Code

แต่รวมถึง

* สิ่งที่เคยทำ
* สิ่งที่เคยแก้
* ปัญหาที่เคยพบ
* วิธีแก้
* ข้อสังเกต
* คำอธิบายเพิ่มเติม
* Context จาก Human
* Context จาก Agent

---

## โครงสร้าง

เก็บเป็น JSON

แยกตาม Symbol

เช่น

* Function
* Class
* Route
* Component

---

## ตัวอย่างข้อมูล

createUser()

Experiences

* เคยแก้เรื่อง Email Verification
* เคยมี Bug เรื่อง Permission
* เคย Refactor Authentication Flow

---

## Learning History

ไม่ลบความรู้เดิม

เพิ่มเป็น Version ใหม่เสมอ

แต่ละ Entry มี

* createdAt
* confidence
* source
* tags
* content

---

## Query Strategy

เลือกข้อมูลที่มี

* confidence สูงสุด
* active = true

เป็นค่าเริ่มต้น

---

## Tags

ใช้ค้นหา

ตัวอย่าง

* auth
* user
* registration
* email
* jwt
* permission

---

# Layer 5 : Knowledge Vault

## แนวคิด

เก็บ Knowledge ที่นิ่งแล้ว

ไม่เปลี่ยนบ่อย

---

## Storage

SQLite

---

## การย้ายข้อมูล

ใช้หลายปัจจัยร่วมกัน

* Stability
* Confidence
* Access Frequency
* Human Approval

ไม่ใช้เวลาอย่างเดียว

---

## ตัวอย่าง

ความรู้ที่

* ไม่เปลี่ยนมานาน
* Confidence สูง
* ถูกยืนยันโดย Human

สามารถย้ายเข้า Vault ได้

---

## Meta Graph

เก็บเพียง Reference

ไม่เก็บข้อมูลทั้งหมด

---

# MCP Driven Architecture

## Agent ไม่มีสิทธิ์เขียนไฟล์โดยตรง

Agent ต้องเรียก MCP เท่านั้น

ตัวอย่าง

* search_symbol
* get_context
* update_memory
* create_node
* update_node
* archive_knowledge

---

## MCP เป็นตัวจัดการ

MCP รับผิดชอบ

* Validation
* Normalization
* Deduplication
* Scoring
* Versioning
* Conflict Resolution

---

# Event Driven Memory

## แนวคิด

ไม่มีการ Sync ทั้งโปรเจกต์

ทุกอย่างเกิดจาก Event

---

## ตัวอย่าง Event

Read

Agent เปิดดู Function

Modify

Agent แก้ไข Function

Knowledge Update

Agent เพิ่มคำอธิบาย

Decision

Human บันทึกเหตุผลการตัดสินใจ

Bug

Agent หรือ Human บันทึกปัญหา

Solution

Agent หรือ Human บันทึกวิธีแก้

---

## ประโยชน์

Knowledge โตตามการใช้งานจริง

ไม่ต้อง Rebuild ทั้งระบบบ่อย

---

# Agent Workflow

## แบบเดิม

Question

↓

Search

↓

Open Node

↓

Search Related Node

↓

Open Related Node

↓

Search Semantic

↓

Open Semantic

↓

Reason

---

## แบบใหม่

Question

↓

Open Meta Graph Node

↓

Reason

---

## กรณีต้องการรายละเอียด

Question

↓

Meta Graph

↓

Follow Ref

↓

Muscle Memory

หรือ

↓

Knowledge Vault

↓

Reason

---

# Separation of Concerns

## Graphify

Structure

---

## Understand

Meaning

---

## Muscle Memory

Experience

---

## Knowledge Vault

Stable Knowledge

---

## MCP

Governance

---

# Long Term Goal

สร้างระบบที่ Agent ทุกตัวสามารถเรียนรู้ร่วมกันได้

ไม่ขึ้นกับ

* Claude
* GPT
* Gemini
* Qwen
* DeepSeek

Agent ทุกตัวเป็นเพียงผู้สร้าง Event

ส่วนความรู้จริงถูกจัดการโดย MCP และ Memory System

---

# Core Principle

Agent ไม่ใช่เจ้าของความรู้

Agent เป็นเพียงผู้สร้าง Context และ Event

MCP Memory System คือ Source of Truth

---

# Final Concept

Agent Memory Operating System (AMOS)

ประกอบด้วย

* Graphify สำหรับ Structure
* Understand สำหรับ Meaning
* Meta Graph สำหรับ Unified Context
* Muscle Memory สำหรับ Experience
* Knowledge Vault สำหรับ Stable Knowledge
* MCP สำหรับ Governance และ Synchronization

เป้าหมายคือทำให้ Agent เข้าใจโปรเจกต์ได้ใกล้เคียงมนุษย์มากขึ้น โดยใช้ข้อมูลที่ถูกสะสมจากการทำงานจริง ไม่ใช่เพียงการวิเคราะห์ Source Code เพียงอย่างเดียว



---

# Incremental Update Architecture

## แนวคิด

AMOS ไม่ควร Scan และ Rebuild ทั้ง Project ทุกครั้งที่มีการเปลี่ยนแปลง

เป้าหมายคือ

- ลดเวลา Update
- ลด Token
- ลด CPU
- ลด Disk IO
- ลดการ Generate Knowledge ซ้ำ

โดยใช้ Incremental Update เป็นหลัก

---

# Initialization

## First Run

คำสั่ง

```bash
amos start
```

หน้าที่

- Scan Project ครั้งแรก
- Build Graphify Structure
- Build Understand Semantic
- Build Meta Graph
- Build Muscle Memory
- Create Project Metadata

---

## Default Ignore

AMOS จะ Skip อัตโนมัติ

```text
node_modules/
.git/
.next/
dist/
build/
coverage/
.cache/

.env
.env.*
.gitignore
.DS_Store
```

รวมถึง File และ Folder ที่ถูก Ignore โดย Git

---

## Project Metadata

สร้าง

```text
.amos/project.json
```

ตัวอย่าง

```json
{
  "projectId": "...",
  "initialized": true,
  "createdAt": "...",
  "version": 1
}
```

---

## Prevent Re-Initialization

หากพบว่า Project ถูก Initialize แล้ว

```bash
amos start
```

จะ Return

```text
AMOS already initialized.

Use:
amos update
```

เพื่อป้องกันการ Build ซ้ำโดยไม่จำเป็น

---

# Incremental Update

## Philosophy

Update เฉพาะสิ่งที่เปลี่ยน

ไม่ Update ทั้ง Project

---

## Command

```bash
amos update
```

---

## Detection Priority

ลำดับการตรวจสอบความเปลี่ยนแปลง

### Priority 1

Git Diff

```bash
git diff
```

---

### Priority 2

RTK Diff

หาก Project ไม่อยู่ภายใต้ Git

---

### Priority 3

File Modified Time

Fallback กรณีไม่สามารถใช้ Diff ได้

---

# Root Path Protection

กรณีมีการเรียก

```bash
amos update .
```

หรือ

```bash
amos update /
```

AMOS จะไม่ Rebuild ทั้ง Project ทันที

แต่จะ

1. ตรวจสอบ Diff
2. วิเคราะห์ไฟล์ที่เปลี่ยน
3. วิเคราะห์ Symbol ที่ได้รับผลกระทบ
4. Update เฉพาะส่วนที่เกี่ยวข้อง

---

# Symbol-Level Update

## Principle

AMOS ทำงานระดับ Symbol

ไม่ใช่ระดับ File

---

ตัวอย่าง

File

```text
src/services/user.service.ts
```

มี

```ts
createUser()
updateUser()
deleteUser()
```

---

หากแก้ไข

```ts
createUser()
```

เพียงตัวเดียว

AMOS ควร Detect

```json
{
  "changedSymbols": [
    "createUser"
  ]
}
```

---

และ Update เฉพาะ

- Graph Node
- Semantic Node
- Muscle Memory
- Knowledge Links

ของ

```text
createUser
```

เท่านั้น

---

# Symbol Hash

## Objective

ลดการ Rebuild ที่ไม่จำเป็น

---

เก็บ Hash ระดับ Symbol

ตัวอย่าง

```json
{
  "symbol": "createUser",
  "hash": "..."
}
```

---

เมื่อ Update

AMOS จะ

1. Parse AST
2. Generate Symbol Hash ใหม่
3. เปรียบเทียบกับ Hash เดิม

หาก Hash ไม่เปลี่ยน

```text
SKIP
```

---

แม้ว่า File จะถูกแก้ไขก็ตาม

---

# Patch-Based Memory Update

## Principle

Memory ต้องถูก Patch

ไม่ใช่ Replace

---

ข้อมูลเดิม

```json
{
  "notes": [
    "jwt migration",
    "email verification"
  ]
}
```

---

Agent เพิ่มข้อมูลใหม่

```text
oauth support
```

---

ผลลัพธ์

```json
{
  "notes": [
    "jwt migration",
    "email verification",
    "oauth support"
  ]
}
```

---

ไม่ควรลบข้อมูลเดิมโดยอัตโนมัติ

ยกเว้นได้รับคำสั่ง Explicit จาก Human หรือ MCP Policy

---

# MCP Controlled Update

## Agent Responsibility

Agent มีหน้าที่เพียง

- Read
- Analyze
- Generate Context
- Send Event

---

## MCP Responsibility

MCP เป็นผู้จัดการ

- Merge
- Normalize
- Validate
- Score
- Version
- Archive

---

Agent ไม่มีสิทธิ์แก้ไข Memory Store โดยตรง

---

# Event Driven Synchronization

## Traditional Approach

```text
Code Change
    ↓
Rescan Entire Project
    ↓
Rebuild Everything
```

---

## AMOS Approach

```text
Code Change
    ↓
Diff
    ↓
Changed Symbol
    ↓
Patch
    ↓
Update Memory
```

---

# Core Rule

ทุกการ Update ต้องเริ่มจาก

```text
Diff
```

ก่อนเสมอ

ไม่ว่าจะถูกเรียกโดย

- Human
- Agent
- MCP
- Automation

เพื่อให้ AMOS สามารถ Update ได้อย่างแม่นยำในระดับ Symbol และหลีกเลี่ยงการ Rebuild ข้อมูลที่ไม่เกี่ยวข้อง

# Project Root Coordinates
- โฟลเดอร์ชื่อ `α` คือ source of truth ของโปรเจกต์
- โฟลเดอร์ `α` จำเป็นต้องอยู่ที่ root ของโปรเจกต์ เพื่อใช้เป็นตัวระบุตำแหน่งของ project root
