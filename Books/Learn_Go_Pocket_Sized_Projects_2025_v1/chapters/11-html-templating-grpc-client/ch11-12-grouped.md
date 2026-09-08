# Chapters 11-12 — HTML Templating + gRPC Client, Wasm & TinyGo (səh. 519-660)

## Bu fəsillər nədən bəhs edir?

(11) Habits backend-inə HTML UI qoşmaq: go:embed, text/template (dəyişən,
range, if, custom funksiyalar, define), form emalı, CSS, gRPC client + minimock.
(12) Wasm və TinyGo: brauzerdə Go (multiplication quiz), js.FuncOf, DOM
manipulyasiyası, Arduino Nano 33-də TinyGo (machine paketi, build tags, LED
blink, traffic lights).

## Əsas fikirlər

### 1. go:embed — şablonu binary-ya daxil et
```go
import "embed"

//go:embed index.xhtml
var indexPage string
```
- Global var + build zamanı fayl məzmunu string-ə EMBED olunur; `embed` paketi
  boş import ( `_ "embed"` tələb etmir — dəyişən üstündə direktiv kifayətdir)

### 2. Template əsasları
```go
tpl, err := template.New("index").Funcs(template.FuncMap{
    "statusCSSClass": statusCSSClass,  // custom funksiya
}).Parse(indexPage)
err = tpl.Execute(w, data)   // data = "." (dot)
```
- **Action:** `{{ . }}` — dot = verilən data; boşluq istənilən: `{{.}}` ≡ `{{ . }}`
- **Sahələr:** `{{ .Name }}` — yalnız EXPOSED sahələr; metod: `{{ .Method }}`
  (parametrsiz, qaytarmaqlı — amma ZİNCİR ÇƏKMƏ: `{{ .A.B.C }}` qadağan tövsiyə)
- **Range:** `{{ range .Habits }}...{{ end }}` — daxildə dot = element;
  index: `{{ range $i, $h := .Habits }}`
- **Şərtlər (truthiness cədvəli):**
  | Go | Template if |
  |---|---|
  | bool | `{{ if . }}` |
  | 0 olmayan ədəd | `{{ if . }}` |
  | boş olmayan string/slice/map | `{{ if . }}` |
  | eq: `{{ if eq .WeeklyFrequency 1 }}` ... `{{ else }}` |
- **Whitespace:** `{{- ... }}` / `{{ ... -}}` — sətir boşluqlarını yemək

### 3. Custom funksiyalar + pipline
```go
template.FuncMap{"statusCSSClass": statusCSSClass}
// istifadə: {{ statusCSSClass . }}
// pipeline: {{ . | len }}  ≡  {{ len . }}
```

### 4. define — alt-şablon
```html
{{- define "habitItem" }}
<li class="habit {{ statusCSSClass . }}">...</li>
{{- end }}
```
- `tpl.ExecuteTemplate(w, "index")` — konkret şablon adı ilə icra
- html/template (escape edir — XSS qoruması) vs text/template (xam mətn)

### 5. Form + POST endpoint
```html
<form action="/create" method="post">
  <input type="text" name="habitName">
  <input type="number" name="weeklyFrequency">
  <button>Create</button>
</form>
```
```go
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
    habitName := r.FormValue("habitName")
    weeklyFreq, err := strconv.Atoi(r.FormValue("weeklyFrequency"))
    // ...
    http.Redirect(w, r, "/", http.StatusSeeOther)  // PRG pattern
}
```

### 6. gRPC client + minimock
```go
// local kiçik interfeys (generated API-ni domain-ə BURAXMA):
type HabitsClient interface {
    ListHabits(ctx context.Context) ([]habit.Habit, error)
}
//go:generate minimock -s "_mock.go" -o "mocks"
```
- Mock testi: `mocks.NewHabitsClientMock(t).ListHabitsMock.Expect(...).Return(...)`
- `replace learngo-pockets/habits => ../10-habits/...` — lokal modul istifadəsi

### 7. Golden files (test)
- Gözlənilən çıxışı faylda saxla, `//go:generate` ilə yenilə — tam müqayisə
  testi struktur dəyişikliyində siqnal verir

### 8. Wasm — brauzerdə Go (Ch12)
```bash
GOOS=js GOARCH=wasm go build -o main.wasm
```
```go
type multiplication struct{ opLeft, opRight int }

func (m *multiplication) generate(this js.Value, args []js.Value) any {
    m.opLeft = rand.IntN(11)
    m.opRight = rand.IntN(11)
    document := js.Global().Get("document")
    document.Call("getElementById", "operandLeft").Set("textContent", m.opLeft)
    return nil
}

func main() {
    m := &multiplication{}
    multiply := js.Global().Get("Object").New()   // namespace obyekti
    multiply.Set("generate", js.FuncOf(m.generate))
    js.Global().Set("multiplyApp", multiply)
    <-make(chan struct{})   // ƏSAS: main bitməz — Wasm callback-lər gözləyir
}
```
**Sub-kod izahı:**
- `js.FuncOf(func(this js.Value, args []js.Value) any)` — JS-ə çağırıla bilən
  funksiya qeydiyyatı
- DOM: `js.Global().Get("document").Call(...)` / `.Set(...)`
- **`<-make(chan struct{})`** — main exit Wasm-i ÖLDÜRÜR; sonsuz blokla!
- wasm_exec.js + `WebAssembly.instantiateStreaming(fetch(...))` — yükləmə
- **Təhlükəsizlik:** .wasm faylı oxuna biləndir — parollar/secret-lər compile
  olunmuş fayldan çıxarıla bilər!

### 9. TinyGo — mikrokontrolerdə Go
```go
package main

import (
    "machine"
    "time"
)

func main() {
    led := machine.LED
    led.Configure(machine.PinConfig{Mode: machine.PinOutput})
    for {
        led.High()
        time.Sleep(time.Second / 2)
        led.Low()
        time.Sleep(time.Second / 2)
    }
}
```
```bash
tinygo build --target=arduino-nano33 main.go   # target = build tag-ləri avtomatik
tinygo flash -target=arduino-nano33 -port=/dev/ttyUSB0
```
- `//go:build arduino_nano33` — build tag: platforma xüsusi pin xəritələri
  (machine_XXX.go faylları hər board üçün ayrıca)
- `machine.D2`, `machine.PinConfig{Mode: machine.PinOutput}`, `led.High()/Low()`
- **Playground:** play.tinygo.org — virtual mikrokontroler + LED/sxem (wasm
  kompilyasiya, real hardware lazım deyil)
- **Traffic lights:** crossing struct (carLight + walkLight + button);
  `go c.listenButton()` — mikrokontrolerdə belə goroutine mümkün (amma ehtiyatlı)
- **Debug:** serial port (println), GDB-əsaslı debugger, ldflags "-w -s" —
  binary kiçiltmə
- **Fərqlər:** runtime GC fərqli, bəzi paketlər yoxdur (reflect məhdud),
  binary ölçüsü kritikdir

## Əsas terminlər

- go:embed (fayl daxiletmə)
- text/template / html/template
- Action / Dot ({{ . }})
- FuncMap (custom funksiya)
- Pipeline ({{ . | fn }})
- Golden File (referens çıxış testi)
- WebAssembly / Wasm
- js.FuncOf (JS körpüsü)
- TinyGo
- Build Tag (//go:build)
- Flashing (mikrokontrolerə yazma)
- machine paketi (pin/LED API)

## Praktik nəticə

- Şablonlar binary-da embed — deploy = tək fayl
- Template-də məntiqi minimal saxla: custom funksiya + CSS class adı qaytar
- Domain paketlər generated API tiplərini import ETMƏZ — lokal mini-interfeys
- Wasm: main sonsuz blokla; js.FuncOf ilə qeydiyyat; .wasm-da secret YOX
- TinyGo: board-özlü kod build tag-lə; playground-da sına

## Mənbə

Pages: 519-660 (Chapters 11-12 + appendices, Learn Go with Pocket-Sized Projects)
