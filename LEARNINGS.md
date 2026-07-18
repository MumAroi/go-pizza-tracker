# สิ่งที่เรียนรู้จาก Pizza Tracker Project

## 1. Project Structure (Clean Architecture)

```
cmd/                    → Entry point (main.go)
internal/               → Business logic (ไม่ export ออกภายนอก)
├── admin/              → Admin module
├── order/              → Order module
├── user/               → User module
├── middleware/          → Auth middleware
├── shared/
│   ├── notification/   → Shared services
│   └── util/           → Utility functions
├── app/                → App initialization
├── config/             → Configuration
├── database/           → Database setup
├── route/              → Route registration
└── session/            → Session store
templates/              → HTML templates
```

**หลักการ:** แยก module ตาม domain แต่ละ module มี handlers, models, repository ของตัวเอง

---

## 2. Interface

```go
type OrderRepository interface {
    GetOrder(id string) (*Order, error)
    CreateOrder(order *Order) error
}
```

**ใช้เมื่อไหร่:**
- ต้อง mock ใน unit test
- มีหลาย implementation (MySQL, PostgreSQL, Mock)
- ลด coupling ระหว่าง layer

**ไม่จำเป็นต้องใช้:**
- Function เล็กๆ ใช้ครั้งเดียว
- มี implementation เดียว (YAGNI)

**หลักการ:** เขียน struct ก่อน → ค่อยแยก interface เมื่อมีเหตุผล

---

## 3. Dependency Injection

```go
// Deps struct สำหรับ constructor
type OrderDeps struct {
    OrderRepo       OrderRepository
    NotificationMgr *notification.NotificationManager
}

// Constructor รับ Deps
func NewHandler(deps OrderDeps) Handler {
    return &handler{
        OrderDeps: deps,
    }
}

// เรียกใช้ - ส่ง dependencies จากภายนอก
orderH := order.NewHandler(order.OrderDeps{
    OrderRepo:       app.OrderRepo,
    NotificationMgr: app.NotificationMgr,
})
```

**ข้อดี:**
- Test ง่าย (ส่ง mock เข้าไปได้)
- Swap implementation ได้ (MySQL → PostgreSQL)
- แยกส่วนชัดเจน

---

## 4. Repository Pattern

```go
type UserRepository interface {
    Authenticate(username, password string) (*User, error)
    GetByID(id string) (*User, error)
}

type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) GetByID(id string) (*User, error) {
    var user User
    if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("Not found user")
        }
        return nil, err
    }
    return &user, nil
}
```

**หลักการ:** Handler ไม่ควรรู้เรื่อง database → เรียกผ่าน Repository interface

---

## 5. Struct Embedding

```go
type OrderDeps struct {
    OrderRepo       OrderRepository
    NotificationMgr *notification.NotificationManager
}

type handler struct {
    OrderDeps  // embed (ไม่มี field name)
}

// เรียกใช้ได้ 2 แบบ (เหมือนกัน)
h.OrderRepo           // promoted field (สั้น)
h.OrderDeps.OrderRepo // full path (ยาว)
```

**ข้อดี:** ลด code ซ้ำซ้อน เมื่อ struct ใหญ่ขึ้น เพิ่ม dependency ได้ง่าย

---

## 6. Go Concurrency

### Channel (ท่อส่งข้อมูล)

```go
ch := make(chan string)  // สร้าง channel
ch <- "hello"            // ส่งข้อมูล
msg := <-ch              // รับข้อมูล
```

### map[chan string]bool (set ของ channels)

```go
// ใช้ map แทน slice เพราะลบง่าย O(1)
clients map[string]map[chan string]bool
//              ↑ outer map (key = order ID)
//                 ↑ inner map (set of channels)

delete(clients, client)  // ลบ channel ออกง่าย
```

### sync.RWMutex (ป้องกัน race condition)

```go
func (n *NotificationManager) AddClient(...) {
    n.mu.Lock()        // เขียน → Lock
    defer n.mu.Unlock()
}

func (n *NotificationManager) Notify(...) {
    n.mu.RLock()       // อ่าน → RLock (หลายคนอ่านพร้อมกันได้)
    defer n.mu.RUnlock()
}
```

| Type | อ่านพร้อมกัน | เขียนพร้อมกัน | เหมาะกับ |
|---|---|---|---|
| `Mutex` | ไม่ได้ | ไม่ได้ | Write-heavy |
| `RWMutex` | ได้หลายคน | ไม่ได้ | Read-heavy |

---

## 7. SSE (Server-Sent Events)

```go
func (n *NotificationManager) StreamSSE(c *gin.Context, client chan string) {
    // ตั้ง headers สำหรับ SSE
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")

    // Stream ข้อมูลแบบ real-time
    c.Stream(func(w io.Writer) bool {
        if msg, ok := <-client; ok {
            c.SSEvent("message", msg)
            return true   // ทำซ้ำต่อ
        }
        return false      // หยุด
    })
}
```

**Non-blocking send (select + default):**

```go
select {
case client <- message:  // ส่งได้ → ส่ง
default:                  // ส่งไม่ได้ (channel เต็ม) → ข้าม
}
```

**SSE vs WebSocket:**

| เทคโนโลยี | Direction | ใช้กับ |
|---|---|---|
| SSE | Server → Client เท่านั้น | Notifications, dashboards |
| WebSocket | Server ↔ Client (2 ทาง) | Chat, gaming |

---

## 8. Session Management

```go
// สร้าง session store
func NewSessionStore(db *gorm.DB, secretKey []byte) sessions.Store {
    store := gormsessions.NewStore(db, true, secretKey)
    store.Options(sessions.Options{
        Path:     "/",
        MaxAge:   86400,
        HttpOnly: true,
        Secure:   false,  // true ใน production (HTTPS)
        SameSite: http.SameSiteLaxMode,
    })
    return store
}

// ตั้งค่า session
func SetSessionValue(c *gin.Context, key string, value any) error {
    session := sessions.Default(c)
    session.Set(key, value)
    return session.Save()
}
```

**ข้อควรระวัง:**
- `Secure: true` → Cookie ต้อง HTTPS เท่านั้น (dev ใช้ `false`)
- `SESSION_SECRET` ต้องไม่ว่าง (ไม่งั้น `securecookie: hash key is not set`)

---

## 9. Auth Middleware

```go
func AuthMiddleware(userRepo user.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := util.GetSessionString(c, "userID")
        if userID == "" {
            c.Redirect(http.StatusSeeOther, "/login")
            c.Abort()
            return
        }
        _, err := userRepo.GetByID(userID)
        if err != nil {
            util.ClearSession(c)
            c.Redirect(http.StatusSeeOther, "/login")
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**หลักการ:**
- Middleware ต้องอยู่ **ก่อน** routes ที่ต้องการ protect
- Login routes ไม่ต้อง auth → แยกเป็น public group

---

## 10. GORM

### Model & Relationships

```go
type Order struct {
    ID    string      `gorm:"primaryKey;size:14"`
    Items []OrderItem `gorm:"foreignKey:OrderID"`
}
```

### Preload (Eager Loading)

```go
// โหลด Items มาพร้อมกับ Order
db.Preload("Items").First(&order, "id = ?", id)
```

### Delete with Association

```go
// ลบ Order + Items ทั้งหมด
db.Select("Items").Delete(&Order{ID: id})
```

### Update

```go
// อัปเดต field เดียว
db.Model(&Order{}).Where("id = ?", id).Update("status", "delivered")
```

### Ordering

```go
db.Order("created_at DESC").Find(&orders)
```

---

## 11. Type Assertion

```go
val, _ := someAnyValue.(string)
//        ↑ แปลงจาก any → string
//        ถ้าแปลงไม่ได้ → val = "" (zero value)
```

---

## 12. make() vs new()

| | `make()` | `new()` |
|---|---|---|
| ใช้กับ | map, slice, channel | struct, primitive types |
| คืนค่า | value ที่ initialize แล้ว | pointer (`*T`) |

```go
m := make(map[string]int)   // พร้อมใช้
ch := make(chan string)      // พร้อมใช้
p := new(Order)              // ได้ *Order (pointer)
```

---

## 13. bcrypt (Password Hashing)

```go
// Hash password
hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Compare password
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
```

---

## 14. Template Functions

```go
functions := template.FuncMap{
    "add": func(a, b int) int {
        return a + b
    },
    "json": func(v interface{}) template.JS {
        b, _ := json.Marshal(v)
        return template.JS(b)
    },
}
```

---

## 15. Constructor Pattern ใน Go

```go
// ไม่มี constructor ใน Go → ใช้ factory function
func NewHandler(deps OrderDeps) Handler {
    return &handler{OrderDeps: deps}
}
```

**ข้อดี:**
- ซ่อน implementation
- เพิ่ม validation ได้ก่อนสร้าง
- Return interface ได้

---

## 16. สิ่งที่ได้เรียนรู้เรื่อง Error Handling

```go
// ❌ Ignore error
util.SetSessionValue(c, "userID", "123")

// ✅ Handle error
if err := util.SetSessionValue(c, "userID", "123"); err != nil {
    slog.Error("Failed to set session", "error", err)
}
```

---

## 17. ข้อควรระวังที่พบในโปรเจค

| ปัญหา | สาเหตุ | วิธีแก้ |
|---|---|---|
| Redirect กลับ login ตลอด | `Secure: true` บน HTTP | ตั้ง `Secure: false` สำหรับ dev |
| `securecookie: hash key is not set` | `SESSION_SECRET` ว่าง | เพิ่มใน `.env` |
| Nil pointer dereference | ลืม inject dependency | ส่ง repo ครบใน `AdminDeps` |
| Admin ไม่ได้รับ notification | Key ไม่ตรงกัน | เช็ค key ให้ตรง (new_order vs new_orders) |
| Interface mismatch | ชื่อ method ไม่ตรง | `DashboardData` vs `Dashboard` |

---

## สรุป Pattern ที่ใช้ในโปรเจค

```
Request → Router → Middleware → Handler → Repository → Database
                                                        ↓
Response ← Router ← Middleware ← Handler ← Repository ←─┘
                                              ↓
                                    NotificationManager → SSE → Browser
```
