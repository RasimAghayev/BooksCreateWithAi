# Chapter 19 — Special Features (Xüsusi Xüsusiyyətlər)

## Bu fəsil nədən bəhs edir?

Build constraints (build tags + filename suffix-lər, GOOS/GOARCH cross-compile),
reflection (reflect.TypeOf/ValueOf, FieldByName, runtime tip yoxlaması),
reflect.DeepEqual (müqayisəolunmayan tiplər), wildcard pattern (`./...`), unsafe
paketi (Float32bits, unsafe.Pointer) və cgo (C kod çağırışı).

## Əsas fikirlər

### 1. Build Constraints — 2 Üsul
**Məqsəd:** faylın compile-a DAXİL olub-olmaması şərtləri; eyni funksiyanın
fərqli OS/arch üçün fərqli implementasiyası (standart kitabxanada geniş istifadə:
os/syscall/runtime paketləri).

**Üsul 1 — Build Tags (sətir şərhi):**
```go
// +build linux                    // YALNIZ linux

// +build amd64,darwin 386,!gccgo   // (amd64 AND darwin) OR (386 AND NOT gccgo)
                                     // vergül = AND, boşluq = OR, ! = NOT

package main                        // BOŞ SƏTİR ilə ayrılmalı!
```
Xüsusi tag: `// +build ignore` — fayl compile-a HEÇ VAXT daxil olmur.

**Üsul 2 — Filename Suffix:**
```
syscall_linux.go              → YALNIZ linux
syscall_windows.go            → windows
signal_darwin_amd64.go       → darwin + amd64
stat_aix.go                  → aix
```
Formalar: `*_GOOS`, `*_GOARCH`, `*_GOOS_GOARCH`.

### 2. GOOS/GOARCH — Cross-Compile
```bash
go env                    # bütün dəyişənlər
go env GOOS GOARCH        # cari OS/arch (məs. darwin, amd64)

GOOS=linux go build -o app    # Linux üçün (macOS-dan!) — cross-compile
GOARCH=386 go build           // 386 üçün
```
Tag/suffix GOOS-uyğun DEYİLSƏ → "cannot find module for path ." build xətası.
Cross-compile Go-nun nadir güclərindəndir — başqa dillərdə demək olar yoxdur.

### 3. Reflection — Runtime Tip İnspekti
**Nədir:** runtime-da obyektin TİPİNi yoxlamaq və dəyərlərini manipulyasiya;
giriş tipi QABAQCADAN BİLMƏDİKDƏ. encoding/json və fmt daxilində istifadə olunur.

**2 sütun:**
```go
reflect.TypeOf(input)   // Type — tip haqqında
reflect.ValueOf(input)  // Value — dəyər haqqında
```

**Kitabdan MyPrint nümunəsi:**
```go
func MyPrint(input interface{}) {
    t := reflect.TypeOf(input)
    v := reflect.ValueOf(input)
    switch {
    case t.Name() == "Animal":
        fmt.Println("I am a ", v.FieldByName("Name"))    // sahə adı ilə çıxış
    case t.Name() == "Object":
        fmt.Println("I am a ", v.FieldByName("Type"))
    default:
        fmt.Println("I got an unknown entity")
    }
}
MyPrint(Object{Type: "Chair"})   // I am a Chair
MyPrint(Animal{Name: "Tiger"})   // I am a Tiger
MyPrint(Person{Name: "Gobin"})   // I got an unknown entity
```

**area() — reflection ilə çoxşekilli funksiya (kitabdan):**
```go
func area(input interface{}) float64 {
    inputType := reflect.TypeOf(input)
    if inputType.Name() == "circle" {
        val := reflect.ValueOf(input)
        radius := val.FieldByName("radius")     // sahə dəyəri
        return math.Pi * math.Pow(radius.Float(), 2)   // .Float() konvertoru
    }
    if inputType.Name() == "rectangle" {
        val := reflect.ValueOf(input)
        return val.FieldByName("length").Float() *
            val.FieldByName("breadth").Float()
    }
    return 0
}
area(circle{radius: 3})        // 28.274334
area(rectangle{length: 3, breadth: 7})  // 21.000000
```
**Xəbərdarlıq:** yanlış çevirmə və ya dəstəklənməyən metod → PANİK; ehtiyatla.

### 4. reflect.DeepEqual
**Problem:** slice/map `==` ilə MÜQAYİSƏ OLUNA BİLMİR (compile xətası).

```go
reflect.DeepEqual(make([]int, 10), make([]int, 10))     // true — ölçü VACİB
reflect.DeepEqual([3]int{1,2,3}, [3]int{1,2,3})          // true — sıra VACİB
reflect.DeepEqual(map[int]string{1:"one"}, map[int]string{1:"one"})  // true
                                                         // map-də sıra ƏHƏMİYYƏTSİZ
reflect.DeepEqual(nil, nil)                               // true
```

### 5. Wildcard Pattern — `./...`
```bash
go list ./...              # cari + ALT qovluqlardakı bütün .go
go test ./...              # bütün testlər recursive
go list -f {{.GoFiles}}{{.Dir}} ./...
```
`./vendor` avtomatik İGNORLANIR; CI pipeline test avtomatizasiyasının standartı.

### 6. unsafe Paketi — Yaddaşa Birbaşa Çıxış
**Nədir:** runtime yaddaş idarəsini BYPASS edən alətlər. Adı kimi — TƏHLÜKƏSİZ;
Go 1 compat ZƏMANƏTSİZ (gələcək versiyalarda qırıla bilər).

**Standart kitabxana nümunələri:**
```go
func Float32bits(f float32) uint32 {
    return *(*uint32)(unsafe.Pointer(&f))    // float → uint32 BİT səviyyəsində
}
func Float32frombits(b uint32) float32 {
    return *(*float32)(unsafe.Pointer(&b))
}
```

### 7. cgo — C Kod Çağırışı
**Kitabdan kod nümunəsi:**
```go
package main

// #include <stdio.h>
// #include <stdlib.h>
// static void myprint(char* s) {
//   printf("%s\n", s);
// }
import "C"                     // PSEUDO-paket — şərh C başlığı kimi oxunur
import "unsafe"

func main() {
    cs := C.CString("Hello World!")    // Go string → C string
    C.myprint(cs)                       // C funksiyası çağır
    C.free(unsafe.Pointer(cs))          // C yaddaşını ƏL İLƏ azad et!
}
```

**Konvertorlar:**
- `C.CString(string) *C.char` — Go→C string
- `C.GoBytes(unsafe.Pointer, C.int) []byte` — C→Go baytları
- `C.free(unsafe.Pointer)` — C yaddaş təmizliyi (GC C-ni BİLMİR!)

**Kitabdan GoBytes nümunəsi:**
```go
cString = C.CString("Hello World!\n")
defer C.free(unsafe.Pointer(cString))          // cleanup MÜTLƏQ
b = C.GoBytes(unsafe.Pointer(cString), C.int(14))  // → []byte
fmt.Print(string(b))
```
Windows-da cgo üçün GCC (MinGW) lazımdır.

## Əsas terminlələr
- Build Constraint — faylın compile-a daxilolma şərti
- Build Tag — `// +build` sətri; AND/OR/NOT məntiqi
- Filename Suffix — `_GOOS`, `_GOARCH`, `_GOOS_GOARCH` ad konvensiyası
- GOOS/GOARCH — OS/arch build dəyişənləri; cross-compile
- `+build ignore` — faylı tam istisna
- Reflection — runtime tip/dəyər inspekti
- reflect.Type / reflect.Value — tip / dəyər deskriptorları
- TypeOf/ValueOf — reflection-a giriş nöqtələri
- FieldByName — struct sahəsinə ad ilə çıxış
- .Float() — Value-dan konkret tip dəyəri
- reflect.DeepEqual — dərin müqayisə; slice/map üçün
- Wildcard `./...` — recursive paket patterni; vendor istisna
- unsafe.Pointer — tip təhlükəsizliyini ləğv edən körpü pointer
- cgo — Go↔C FFI; `import "C"` pseudo-paketi
- C.CString/GoBytes/C.free — tip konvertorları + manual yaddaş idarəsi
- Go 1 Compatibility — unsafe ZƏMANƏTSİZ

## Praktik nətidə

(1) Fərqli OS kodu: tag VƏ YA suffix — standart kitabxananın öz üsulu; suffix daha
gizli (fayl adından görünür). (2) `// +build` və `package` arası BOŞ SƏTİR şərtidir —
yoxsa tag işləməz. (3) GOOS/GOARCH dəyişmək = cross-compile — dağıtım üçün binary
hazırlamaq bir əmrlə. (4) Reflection: tip BİLİNMİRSƏ yalnız — JSON/fmt kimi;
yoxsa interfeys/type switch daha təmizdir; səhv istifadə PANİK. (5) FieldByName +
.Float() kombinasiyası — struct sahəsini runtime-da oxu. (6) Slice/map == ile
MÜQAYİSƏ OLMAZ — DeepEqual (testlərdə əvəzsiz). (7) `go test ./...` CI standartı —
bütün paketləri bir əmrlə yoxla. (8) unsafe: YALNız son çarə — bit səviyyəli çevirmə
(Math paketi nümunəsi) və cgo; gələcək qırılma riski daşıyır. (9) cgo-da C yaddaşı
C.free ilə SƏN təmizləyirsən — GC C-nin malloc-unu görmür. (10) `import "C"`-dən
ƏVVƏLKI şərhlər C başlığı kimi compile olunur — funksiyaları orada define et.

## Mənbə
Pages: 661-680 (PDF 692-715)
