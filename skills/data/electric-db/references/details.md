## Concept

```
Postgres → Electric (HTTP Stream) → TanStack DB Collections → UI (reactive)
```

ElectricSQL sync ข้อมูลจาก Postgres มาที่ client ผ่าน HTTP Streaming
TanStack DB รับมาเก็บเป็น reactive collection — join เกิดบน **client**

---

## แหล่งความรู้

| หัวข้อ | URL |
|--------|-----|
| ElectricSQL Docs | https://electric-sql.com/docs |
| Shapes API | https://electric-sql.com/docs/api/http |
| TanStack DB | https://tanstack.com/db/latest |
| TanStack DB + Electric | https://tanstack.com/db/latest/docs/framework/react/guides/electric-sql |
| Examples | https://github.com/electric-sql/electric/tree/main/examples |

---

## Core Concepts

### Shape — หน่วย sync

```ts
const queryShape = createElectricShape({
  url: `${ELECTRIC_URL}/v1/shape`,
  table: 'query',
  where: `status != 'Completed'`, // filter ที่ต้นทาง
})
```

### Collection — client-side table

```ts
const queryCollection = createCollection({
  id: 'queries',
  shape: queryShape,
  getKey: (row) => row.id, // ต้องเป็น unique key เสมอ
})
```

### useQuery — reactive join บน client

```ts
const results = useQuery(queryCollection, (q) =>
  q
    .from(queryCollection)
    .select((item) => ({
      ...item,
      approvers: approverCollection
        .filter((a) => a.queryId === item.id)
        .toArray(),
      assignments: assignmentCollection
        .filter((a) => a.queryId === item.id)
        .toArray(),
    }))
    .orderBy((q) => q.createdAt, 'desc')
)
```

---

## Best Practices

### 1. Filter Shape ให้แคบที่สุด
```ts
// ❌
createElectricShape({ table: 'query_assignment' })

// ✅
createElectricShape({
  table: 'query_assignment',
  where: `hoscode = '${currentHoscode}'`,
})
```

### 2. One Shape Per Table — join บน client
```ts
const queryShape      = createElectricShape({ table: 'query', ... })
const approverShape   = createElectricShape({ table: 'query_approver', ... })
const assignmentShape = createElectricShape({ table: 'query_assignment', ... })
```

### 3. Composite key
```ts
getKey: (row) => `${row.queryId}-${row.hoscode}`
```

### 4. แทน refetchInterval
```ts
// ❌ เดิม
useQuery({ queryFn: fetchQueries, refetchInterval: 30_000 })

// ✅ ใหม่
const queries = useQuery(queryCollection, (q) => q.from(queryCollection))
```

### 5. Optimistic Update
```ts
const transact = useTransact()

await transact(({ mutate }) => {
  mutate(queryCollection).update({ key: id, value: { status } })
})
await api.patch(`/queries/${id}`, { status })
```

---

## ตัวอย่างเต็ม — Query Dashboard

```ts
// shapes.ts
export const queryShape = createElectricShape({
  url: `${import.meta.env.VITE_ELECTRIC_URL}/v1/shape`,
  table: 'query',
  where: `status IN ('Pending', 'Processing')`,
})
export const approverShape = createElectricShape({
  url: `${import.meta.env.VITE_ELECTRIC_URL}/v1/shape`,
  table: 'query_approver',
})
export const assignmentShape = createElectricShape({
  url: `${import.meta.env.VITE_ELECTRIC_URL}/v1/shape`,
  table: 'query_assignment',
  where: `hoscode = '${userHoscode}'`,
})
```

```ts
// collections.ts
export const queryCollection = createCollection({
  id: 'queries',
  shape: queryShape,
  getKey: (row) => row.id,
})
export const approverCollection = createCollection({
  id: 'query-approvers',
  shape: approverShape,
  getKey: (row) => row.id,
})
export const assignmentCollection = createCollection({
  id: 'query-assignments',
  shape: assignmentShape,
  getKey: (row) => `${row.queryId}-${row.hoscode}`,
})
```

```tsx
// QueryDashboard.tsx
function QueryDashboard() {
  const queries = useQuery(queryCollection, (q) =>
    q
      .from(queryCollection)
      .select((query) => ({
        ...query,
        approvers: approverCollection
          .filter((a) => a.queryId === query.id)
          .toArray(),
        assignments: assignmentCollection
          .filter((a) => a.queryId === query.id)
          .toArray(),
        assignmentCount: assignmentCollection
          .filter((a) => a.queryId === query.id)
          .count(),
      }))
      .orderBy((q) => q.createdAt, 'desc')
  )

  return (
    <div>
      {queries.map((q) => (
        <QueryCard key={q.id} data={q} />
      ))}
    </div>
  )
}
```

---

## เปรียบเทียบ

| | REST + refetchInterval | ElectricSQL + TanStack DB |
|--|--|--|
| Update latency | ทุก 30s | < 1s |
| Network | request ทุกรอบ | streaming ตลอด |
| Join | server-side | client-side |
| Offline | ❌ | ✅ |
| Optimistic UI | ยาก | built-in |

---

## ข้อจำกัด

- Postgres ต้องเปิด logical replication
- Shape `where` รองรับแค่ simple conditions (ไม่มี JOIN บน server)
- TanStack DB ยังอยู่ใน beta — API อาจเปลี่ยน
