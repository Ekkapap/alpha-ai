# AMOS (Agent Memory Operating System)

## Core Vision

AMOS ไม่ใช่ Agent

AMOS คือ

* Memory System
* Knowledge System
* Graph System
* Task Orchestrator
* Context Orchestrator
* Event Store

ส่วน AI ต่างๆ เป็นเพียง Worker

* Claude Code
* Gemini CLI
* Antigravity
* Codex CLI
* Local Models

---

# Architecture

Human
↓
WebUI
↓
AMOS
↓
PostgreSQL (Source of Truth)
↓
Workers

Workers

* Claude Code
* Gemini CLI
* Antigravity Agent
* Codex CLI
* Future Workers

---

# Source of Truth

ทุกอย่างเก็บใน PostgreSQL

ตัวอย่าง

* events
* tasks
* memories
* knowledge
* workers
* channels
* reports

AMOS ไม่ใช้ Memory ของ Agent เป็น Source of Truth

AMOS เป็นเจ้าของข้อมูลทั้งหมด

---

# Worker Concept

Worker เป็น Stateless ให้มากที่สุด

Worker ไม่ต้องรู้ว่า

* Memory อยู่ไหน
* Graph อยู่ไหน
* Knowledge อยู่ไหน

Worker รู้เพียง

* รับงาน
* ทำงาน
* ส่งผลกลับ

---

# Context Package

AMOS ไม่ส่งแค่ Task

AMOS ส่ง Context Package

ประกอบด้วย

* Role
* Scope
* Rules
* Skills
* Memory
* Knowledge
* Graph Context
* Messages
* Task

ตัวอย่าง

Role

Senior TypeScript Refactor Engineer

Scope

src/auth

Rules

* strict typescript
* no any
* backward compatible

Memory

* auth refactor history
* jwt related experience

Knowledge

* auth architecture
* permission architecture

Graph Context

* related files
* dependencies
* references

---

# Task Distribution

Worker ไม่เรียก task_list()

Worker เรียก

next_task()

AMOS จะ

* เลือกงาน
* Lock งาน
* ส่ง Context Package

ทันที

Worker ตัวอื่นจะไม่ได้งานซ้ำ

---

# Polling Strategy

ไม่ใช้ Push

ใช้ Poll

ตัวอย่าง

Claude Code
ทุก 30 วินาที

Gemini CLI
ทุก 40 วินาที

Antigravity
ทุก 50 วินาที

เหลื่อมกันเพื่อให้ระบบมี Worker ตื่นอยู่ตลอด

---

# Human Communication

Human ไม่ส่ง Prompt ตรงเข้า Agent

Human ส่ง Message เข้า AMOS

WebUI
↓
PostgreSQL
↓
AMOS

Agent จะอ่านจาก AMOS ในรอบ Poll ถัดไป

---

# Event Driven Memory

AMOS ไม่เชื่อ Summary

AMOS เชื่อ Event

หลักการ

Trust the logs.
Don't trust the summary.

---

# Agent Report Format

Agent ควรส่ง

{
"taskId": "...",
"reason": [],
"event": [],
"action": [],
"result": {}
}

---

# Meaning of Fields

reason

เหตุผลที่ตัดสินใจทำ

ตัวอย่าง

* jwt validation duplicated
* reduce maintenance cost

event

สิ่งที่พบ

ตัวอย่าง

* duplicate function found
* circular dependency detected

action

สิ่งที่ทำจริง

ตัวอย่าง

* created jwt.service.ts
* moved verifyToken()

result

ผลลัพธ์

ตัวอย่าง

* tests passed
* duplicate removed

---

# Event Log is More Important Than Summary

AMOS ควรเก็บ

* edit file
* run command
* run test
* git diff
* graph update

ทั้งหมดเป็น Event

ตัวอย่าง

edit auth.ts
run test
edit auth.ts
run test
edit auth.ts
run test

Qwen สามารถสรุปได้เองว่า

Attempts = 3

Success = true

โดยไม่ต้องให้ Agent รายงาน

---

# Local Intelligence Layer

Mac Mini

MLX

รัน

* Qwen 3.5 4B
* Qwen 3.5 8B
* Qwen 2.5 7B
* Qwen 2.5 14B

เพื่อจัดการ

* Memory Extraction
* Knowledge Extraction
* Tag Generation
* Semantic Labeling
* Report Generation

ฟรี

---

# Why Local Model

งานเหล่านี้ไม่ต้องการ

* Claude
* GPT
* Gemini Pro

งานเหล่านี้ต้องการ

* อ่านเก่ง
* เชื่อมโยงเก่ง
* จัดหมวดเก่ง

Qwen เหมาะมาก

---

# Graph Layer

ใช้

Tree-sitter
↓
Graphify
↓
Structural Graph

Tree-sitter ทำหน้าที่

* parse source code
* function
* class
* import
* export
* call
* reference

---

# Understand Layer

Understand ทำหน้าที่

* semantic meaning
* business meaning
* domain meaning

ตัวอย่าง

auth.middleware.ts

ไม่ใช่แค่

imports jwt.ts

แต่คือ

Authentication Middleware

---

# Unified Knowledge Pipeline

Raw Context
↓
Qwen
↓
Memory Draft
↓
Graphify
↓
Structural Graph
↓
Understand
↓
Semantic Graph
↓
Qwen Merge
↓
Unified Report
↓
PostgreSQL

---

# Memory Hierarchy

Level 0

Raw Experience

เก็บ

* context
* events
* actions
* reports
* diffs

ห้ามทิ้ง

---

Level 1

Memory

สกัดจาก Raw Experience

---

Level 2

Knowledge

สกัดจาก Memory

---

Level 3

Graph

สกัดจาก Knowledge

---

# Why Keep Raw Experience

อนาคต

Qwen 6
Qwen 7
Gemma
DeepSeek
Llama

อาจเก่งกว่า

สามารถ Rebuild

* Memory
* Knowledge
* Graph

ทั้งหมดได้ใหม่

เพราะยังมี Raw Experience

---

# Worker Specialization

Claude Code

* coding
* refactor
* debugging

Gemini CLI

* architecture
* analysis
* review

Antigravity Flash

* monitoring
* scheduling
* dispatching

Local Qwen

* memory
* knowledge
* graph enrichment

---

# Context Reduction Strategy

Agent ไม่ควรถือทั้งโปรเจกต์

Claude

ถือแค่

src/auth

Gemini

ถือแค่

Architecture Context

Planner

ถือแค่

Task Context

ผลลัพธ์

Context เล็กลง

Quota ใช้น้อยลง

ประสิทธิภาพสูงขึ้น

---

# Multi Provider Strategy

ใช้หลาย Provider พร้อมกัน

Claude Code
+
Gemini CLI
+
Antigravity

ผลลัพธ์

* Quota รวมกัน
* Context แยกกัน
* Failure Isolation

Provider ใดล่ม
ระบบยังทำงานต่อได้

---

# AMOS Philosophy

Agent ไม่ใช่ Memory

Agent ไม่ใช่ Knowledge

Agent ไม่ใช่ Source of Truth

Agent คือ Worker

AMOS คือ สมองส่วนกลาง

ทุกสิ่งต้องสามารถสร้างใหม่ได้จาก

Raw Experience
+
Events
+
Graph
+
Knowledge

โดยไม่ผูกติดกับ AI ค่ายใดค่ายหนึ่ง
