# Contributing to α ALPHA

ขอบคุณที่สนใจ contribute! / Thanks for your interest in contributing!

---

## ภาษาไทย

### วิธี Contribute

1. **Fork** repo นี้
2. สร้าง branch ใหม่: `git checkout -b feature/your-feature`
3. ทำการเปลี่ยนแปลง แล้ว commit: `git commit -m "add: your feature"`
4. Push ขึ้น fork ของคุณ: `git push origin feature/your-feature`
5. เปิด **Pull Request** มาที่ `main`

### สิ่งที่ต้องการ Contribute

- **skills/** — เพิ่ม skill ใหม่สำหรับ framework / library ต่างๆ
- **rules/** — เพิ่ม rules สำหรับ AI tool อื่นๆ
- **install.sh** — แก้ bug หรือเพิ่ม AI tool ใหม่
- **docs** — แก้ไข README หรือ documentation

### Skill Structure

```
skills/
└── your-skill/
    ├── SKILL.md        # metadata + description
    └── references/
        └── details.md  # เนื้อหาหลัก
```

### Commit Convention

```
add: เพิ่มฟีเจอร์ใหม่
fix: แก้ bug
update: อัปเดตของเดิม
docs: แก้ documentation
```

### ทดสอบก่อน PR

```bash
rm -rf /tmp/test-alpha && mkdir -p /tmp/test-alpha
bash install.sh /tmp/test-alpha
```

---

## English

### How to Contribute

1. **Fork** this repository
2. Create a new branch: `git checkout -b feature/your-feature`
3. Make your changes and commit: `git commit -m "add: your feature"`
4. Push to your fork: `git push origin feature/your-feature`
5. Open a **Pull Request** targeting `main`

### What to Contribute

- **skills/** — Add skills for new frameworks / libraries
- **rules/** — Add rules for other AI tools
- **install.sh** — Bug fixes or support for new AI tools
- **docs** — README or documentation improvements

### Skill Structure

```
skills/
└── your-skill/
    ├── SKILL.md        # metadata + description
    └── references/
        └── details.md  # main content
```

### Commit Convention

```
add: new feature
fix: bug fix
update: update existing
docs: documentation changes
```

### Test Before PR

```bash
rm -rf /tmp/test-alpha && mkdir -p /tmp/test-alpha
bash install.sh /tmp/test-alpha
```

---

## Issues

- **Bug report** — เปิด issue พร้อม OS, output ที่เห็น, และ steps to reproduce
- **Feature request** — อธิบาย use case ว่าทำไมถึงมีประโยชน์

## Questions

เปิด Discussion หรือ issue label `question` ได้เลยครับ
