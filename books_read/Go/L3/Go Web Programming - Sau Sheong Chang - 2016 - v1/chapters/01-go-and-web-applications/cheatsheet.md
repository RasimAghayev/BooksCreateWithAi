### `http.HandleFunc("/", handler)`

**Nə edir:** Root URL path (`/`) üçün `handler` funksiyasını qeyd edir. Sorğu gəldikdə avtomatik olaraq bu funksiya çağırılır.

**Sub-komanda/parametr izahı:**
- `/` → Path pattern (route)
- `handler` → `func(http.ResponseWriter, *http.Request)` imzalı funksiya

**Mənbə:** Chapter 1, page 40

---

### `http.ListenAndServe(":8080", nil)`

**Nə edir:** 8080 portunda TCP dinləyicisi başladır və gələn sorğuları handler-lərə yönləndirir.

**Sub-komanda/parametr izahı:**
- `:8080` → Port nömrəsi
- `nil` → Default multiplexer (ServeMux) istifadə et

**Mənbə:** Chapter 1, page 40

---

### `fmt.Fprintln(w, "Hello, Go!")`

**Nə edir:** `http.ResponseWriter`-ə mətn yazır və avtomatik olaraq yeni sətir əlavə edir.

**Sub-komanda/parametr izahı:**
- `w` → `http.ResponseWriter` interfezi
- `"Hello, Go!"` → cavab body-si

**Mənbə:** Chapter 1, page 40

---

### `package main`

**Nə edir:** Faylın `main` paketinə aid olduğunu bildirir. İcra olunan proqramlar üçün məcburi paket adıdır.

**Mənbə:** Chapter 1, page 40

---

### `import "net/http"`

**Nə edir:** HTTP server və client funksiyalarını təmin edən standart kitabxananı daxil edir.

**Sub-komanda/parametr izahı:**
- `net/http` → HTTP protokolunun Go implementasiyası

**Mənbə:** Chapter 1, page 40

---

### `func handler(w http.ResponseWriter, r *http.Request)`

**Nə edir:** HTTP handler funksiyasını təyin edir. İlk parametr cavab yazmaq üçün, ikinci parametr gələn sorğunun məlumatlarıdır.

**Sub-komanda/parametr izahı:**
- `w` → Cavab yazmaq üçün `http.ResponseWriter`
- `r` → Sorğu məlumatları (`Method`, `URL`, `Header`, `Body`)

**Mənbə:** Chapter 1, page 40

---

### `$ go run main.go`

**Nə edir:** `main.go` faylını kompilyasiya edib birbaşa icra edir. `go build` ilə yaradılmış binary-lərə ehtiyac qalmadan test etmək üçün istifadə olunur.

**Sub-komanda/parametr izahı:**
- `run` → Kompilyasiya + icra
- `main.go` → Giriş faylı

**Mənbə:** Chapter 1, page 40
