# Chapter 12 — Go for other architectures (səh. 581-625)

## Bu chapter nədən bəhs edir?

Go kodunun brauzerdə (WebAssembly) və mikrokontrollerlərdə (TinyGo) işləməsi:
Wasm compile, DOM manipulyasiyası (`syscall/js`), funksiya register etmə,
koşullu main dayandırma; TinyGo ilə Arduino traffic light layihəsi — pin
konfiqurasiyası, build tags, flashing, button listener, debug.

## Əsas fikirlər

### 1. WebAssembly (Wasm) nədir?
**Necə işləyir:** Go kodu → bytecode compile → **Wasm interpreter** (brauzerin
JS kitabxanasından) işlədir. JS-in yavaş interpretasiyasından qaçmaq üçün
wasm-binary yüngül və sürətlidir.

**Compile (böyük hərflə):**
```bash
GOOS=js GOARCH=wasm go build -o main.wasm
# Windows-da:
set GOOS=js & set GOARCH=wasm & go build -o main.wasm
```

**wasm_exec.js:** Go quraşdırmasından kopyalanan standart Wasm runner — hər
layihədə eynidir (Go 1.23-də `/misc/` və ya `/lib/` müzakirəsi).

**HTML yükləməsi:**
```html
<script src="wasm_exec.js"></script>
<script>
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject)
        .then((result) => go.run(result.instance));
</script>
```

**Təhlükəsizlik xəbərdarlığı:** .wasm faylı brauzerə yüklənəndə hər kəs onu
əldə edib **retro-engineer** edə bilər (parollar, tokendlər çıxarıla bilər).
Obfuskasiya çətinləşdirir, amma QORUMUR. Həssas məlumat heç vaxt client-side
kodda olmamalıdır.

### 2. HTTP fayl serveri (embed ilə)
`file://` yerinə HTTP daha yaxşıdır (production-a bənzəyir):
```go
//go:embed index.xhtml
//go:embed main.wasm
//go:embed wasm_exec.js
var assets embed.FS

func main() {
    fs := http.FileServer(http.FS(assets))
    http.Handle("/", fs)
    http.ListenAndServe("127.0.0.1:30001", nil)
}
```

### 3. DOM manipulyasiyası — syscall/js
```go
import "syscall/js"

type multiplication struct {
    opLeft, opRight int
}

func main() {
    m := &multiplication{opLeft: rand.IntN(11), opRight: rand.IntN(11)}
    document := js.Global().Get("document")
    document.Call("getElementById", "operand1").Set("innerHTML", m.opLeft)
    document.Call("getElementById", "operand2").Set("innerHTML", m.opRight)
}
```
**Sub-kod izahı:**
- `js.Global()` — qlobal JS scope (window)
- `.Get("document")` — property oxu; `.Call("metod", args...)` — metod çağır
- `.Set("innerHTML", dəyər)` — property yaz

### 4. Go funksiyasını JS-ə aç — js.FuncOf
```go
js.Global().Set("generate", js.FuncOf(m.generate))
```
- Callback imzası: `func(this Value, args []Value) any`
- Xəta idarəsi YOXDUR — funksiya error qaytarsa `panic` baş verir
- `js.Func` çox çağırılınca **leak** olur

**Vacib proqram fiqrəsi — əsas kanal:** Wasm-da `main` çıxanda proqam dayanır:
```
Error: Go program has already exited
```
**Həll:** main-i sonsuz gözlət:
```go
func main() {
    js.Global().Set("generate", js.FuncOf(m.generate))
    <-make(chan struct{})  // sonsuz gözlə — proqram canlı qalsın
    // və ya: select {}
}
```

**Namespace (s Scope izolyasiyası):**
```go
multiply := js.Global().Get("Object").New()   // yeni JS obyekt
multiply.Set("generate", js.FuncOf(m.generate))
js.Global().Set("multiplyApp", multiply)       // qlobal scope-a qoy
// istifadə: onclick="multiplyApp.generate()"
```

### 5. Form input oxuma
```html
<button onclick="validate(providedAnswer.value)">Validate</button>
```
JS tərəfdən dəyəri ötür (təmiz üsul), yaxud Go-da:
`dom.Call("getElementById", "providedAnswer").Get("value")`.
Cavabı `strconv.Atoi` ilə parse et — `text` inputda `-`, `+`, `e` də girilə
bilər; parse xətasında istifadəçiyə bildir (innerHTML ilə), konsola YOX.

Input təmizləmək: `defer` ilə `Set("value", "")`.

### 6. TinyGo — mikrokontrollerlər üçün Go
**Nədir:** Kiçik cihazlar üçün Go compileri — standart Go-dan fərqli, daha
məhdud runtime (garbage collector optimallaşdırılıb, bəzi paketlər yoxdur).

**Virtual mikrocontroller:** play.tinygo.org — kodu paste et, Flash düyməsinə
bas, LED-lərin real simulyasiyasını gör (direnclər əlavə etmək lazım deyil!).

### 7. Build tags (qurğu spesifik kod)
```go
//go:build arduino_nano33
const (
    D13 Pin = PA17
    LED = D13       // alias
)
```
- `//go:build {tag}` — faylın yalnız həmin mühitdə compile olunması
- İnkar: `//go:build !windows`
- Bir neçə şərt: `//go:build linux && amd64`
- TinyGo `--target=arduino-nano33` özü düzgün tag ağacını reallaşdırır

### 8. LED blink — ən sadə TinyGo proqramı
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
**Sub-kod izahı:**
- `machine` paketi — pin idarəsi (hər board üçün fərqli mapləmə faylı)
- `Configure(PinOutput)` — pini çıxış kimi hazırla
- `High()/Low()` — yandır/söndür

### 9. Traffic light layihəsi (TinyCity)
**Konstruksiyalar:**
```go
func newCarLight(redPin, amberPin, greenPin machine.Pin) *carLight {
    c := &carLight{red: redPin, amber: amberPin, green: greenPin}
    c.red.Configure(machine.PinConfig{Mode: machine.PinOutput})
    c.amber.Configure(machine.PinConfig{Mode: machine.PinOutput})
    c.green.Configure(machine.PinConfig{Mode: machine.PinOutput})
    c.red.High()   // başlanğıcda qırmızı
    c.amber.Low()
    c.green.Low()
    return c
}
```
**Əsas döngü:**
```go
for {
    walk.Stop()
    car.Go()
    time.Sleep(time.Second * 5)
    car.Stop()
    walk.Go()
    time.Sleep(time.Second * 5)
}
```

**Piəstrian düyməsi (D7):**
```go
type crossing struct {
    cars, walks   *carLight / *walkLight
    button        *machine.Pin
    pedestriansGo bool
    buttonPressed chan struct{}
}

go c.listenButton()      // ayrıca goroutine — düyməni izlə

for {
    c.Switch()
    select {
    case <-c.buttonPressed:      // piəstrian basdı → tez keç
    case <-time.After(time.Second * 5):  // vaxt bitdi → normal keç
    }
}
```

### 10. Flashing (yeridilmə)
```bash
tinygo build --target=arduino-nano33 main.go      # test build
tinygo flash -target=arduino-nano33 -port=/dev/ttyUSB0  # yeridilmə
```
- `--target` özü build tag-ləri həll edir
- Real cihazda: datasheet-dən pin sxemi oxu (analog/digital/GND ayrılır)
- LED polarizəlidir: anod → resistor → cərəyan; katod → GND

### 11. Fərqlər və məhdudiyyətlər
- **Goroutine overhead:** bir neçə KB — müasir PC-də əhəmiyyətsiz, amma
  Arduino Nano 33 IoT-in 256 kiB RAM-ında çoxlu goroutine = yaddaş bitiri
- **Debug:** serial port (UART) çapı və ya TinyGo debugger (GDB əsaslı) —
  panic stack trace çap edə bilir
- Binary kiçiltmək: `go build -ldflags "-w -s"` (debug məlumatını at)

## Əsas terminlər

- WebAssembly (Wasm)
- wasm_exec.js
- syscall/js
- DOM (Document Object Model)
- js.FuncOf
- TinyGo
- Microcontroller (mikrokontroller)
- Build Tag (build constraint)
- Flashing (proqram yeridilməsi)
- Pin ( GPIO pin)
- Datasheet (data cədvəli)

## Praktik nəticə

- Wasm: `GOOS=js GOARCH=wasm` + `wasm_exec.js` + `js.Global()` ilə DOM-a çıxış
- Wasm main-i `select {}` / `<-make(chan struct{})` ilə sağ saxla
- Həssas kod/məlumat client-side .wasm-da saxlanmır — çıxarıla bilər
- TinyGo: `machine` paketi + `//go:build` tag-ləri + `tinygo flash --target=...`
- Embed ilə bütün veb layihəni tək binary-yə yığ (`//go:embed` + `http.FS`)

## Mənbə

Pages: 581-625 (Chapter 12, Learn Go with Pocket-Sized Projects)
