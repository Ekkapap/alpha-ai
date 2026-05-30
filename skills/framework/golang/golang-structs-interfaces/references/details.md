# Go Structs & Interfaces

## Interface Design Principles

### Keep Interfaces Small

> "The bigger the interface, the weaker the abstraction." — Go Proverbs

Interfaces SHOULD have 1-3 methods. Compose larger interfaces from smaller ones:

```go
type ReadWriteCloser interface {
    io.Reader
    io.Writer
    io.Closer
}
```

### Define Interfaces Where They're Consumed

Interfaces MUST be defined where consumed, not where implemented. This keeps the consumer in control of the contract.

```go
// package notification — defines only what it needs
type Sender interface {
    Send(to, body string) error
}
```

### Accept Interfaces, Return Structs

```go
// Good — accepts interface, returns concrete
func NewService(store UserStore) *Service { ... }

// BAD — NEVER return interfaces from constructors
func NewService(store UserStore) ServiceInterface { ... }
```

### Don't Create Interfaces Prematurely

NEVER create interfaces prematurely — wait for 2+ implementations or a testability requirement. Start with concrete types; extract an interface when a second consumer or a test mock demands it.

## Make the Zero Value Useful

```go
// Good — zero value is ready to use
var buf bytes.Buffer
buf.WriteString("hello")

// Good — lazy initialization guards the zero value
func (r *Registry) Register(name string, item Item) {
    if r.items == nil {
        r.items = make(map[string]Item)
    }
    r.items[name] = item
}
```

## Avoid `any` / `interface{}` When a Specific Type Will Do

```go
// Bad — loses type safety
func Contains(slice []any, target any) bool { ... }

// Good — generic, type-safe
func Contains[T comparable](slice []T, target T) bool { ... }
```

## Key Standard Library Interfaces

| Interface | Package | Method |
| --- | --- | --- |
| `Reader` | `io` | `Read(p []byte) (n int, err error)` |
| `Writer` | `io` | `Write(p []byte) (n int, err error)` |
| `Closer` | `io` | `Close() error` |
| `Stringer` | `fmt` | `String() string` |
| `error` | builtin | `Error() string` |
| `Handler` | `net/http` | `ServeHTTP(ResponseWriter, *Request)` |
| `Marshaler` | `encoding/json` | `MarshalJSON() ([]byte, error)` |
| `Unmarshaler` | `encoding/json` | `UnmarshalJSON([]byte) error` |

## Compile-Time Interface Check

```go
var _ io.ReadWriter = (*MyBuffer)(nil)
```

## Type Assertions & Type Switches

```go
// Safe type assertion — use comma-ok form
s, ok := val.(string)
if !ok { /* handle */ }

// Type switch
switch v := val.(type) {
case string:
    fmt.Println(v)
case int:
    fmt.Println(v * 2)
case io.Reader:
    io.Copy(os.Stdout, v)
default:
    fmt.Printf("unexpected type %T\n", v)
}

// Optional behavior pattern
func writeData(w io.Writer, data []byte) error {
    if _, err := w.Write(data); err != nil { return err }
    if f, ok := w.(Flusher); ok { return f.Flush() }
    return nil
}
```

## Struct Embedding

```go
type Server struct {
    Logger           // promoted methods accessible on Server
    addr string
}

// Embed when outer type IS an enhanced version
type Server struct { http.Handler }

// Named field when outer type only HAS the dependency
type Server struct { store *DataStore }
```

## Struct Field Tags

```go
type Order struct {
    ID        string    `json:"id"         db:"id"`
    UserID    string    `json:"user_id"    db:"user_id"`
    Total     float64   `json:"total"      db:"total"`
    Items     []Item    `json:"items"      db:"-"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    Internal  string    `json:"-"          db:"-"`
}
```

| Directive | Meaning |
| --- | --- |
| `json:"name"` | Field name in JSON output |
| `json:"name,omitempty"` | Omit if zero value |
| `json:"-"` | Always exclude |
| `json:",string"` | Encode number/bool as JSON string |
| `db:"column"` | Database column mapping |
| `validate:"required"` | Struct validation |

## Pointer vs Value Receivers

Use pointer `(s *Server)` when: method modifies receiver, receiver contains sync.Mutex, receiver is large struct.
Use value `(s Server)` when: receiver is small and immutable, method is read-only accessor.

Receiver type MUST be consistent across all methods of a type.

## Preventing Struct Copies with `noCopy`

```go
type noCopy struct{}
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type ConnPool struct {
    noCopy noCopy
    mu     sync.Mutex
    conns  []*Conn
}
// go vet will flag any ConnPool passed by value
```

## Common Mistakes

| Mistake | Fix |
| --- | --- |
| Large interfaces (5+ methods) | Split into focused 1-3 method interfaces |
| Defining interfaces in implementor package | Define where consumed |
| Returning interfaces from constructors | Return concrete types |
| Bare type assertions without comma-ok | Always use `v, ok := x.(T)` |
| Missing field tags on serialized structs | Tag all exported fields in marshaled types |
| Mixing pointer and value receivers | Pick one and be consistent |
| Missing compile-time interface check | Add `var _ Interface = (*Type)(nil)` |
| Using `any` for type-safe operations | Use generics `[T comparable]` instead |
| Premature interface with single implementation | Start concrete, extract when needed |
| Nil map/slice in zero value struct | Use lazy initialization in methods |
