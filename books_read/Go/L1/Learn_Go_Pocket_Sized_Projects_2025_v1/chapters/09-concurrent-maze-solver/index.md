# Chapter 9 — Concurrent maze solver (səh. 393-451)

## Bu chapter nədən bəhs edir?

PNG labirint şəklini oxuyub çıxışa aparan yolu tapan concurreny CLI aləti:
goroutine-lərlə "Daedalus" kəşfiyyatı, channel-lərlə kommunikasiya, linked list
path, WaitGroup, select/quit siqnalı, xəzinə tapılanda dayandırma, explored
pixel-lərin qeydi və GIF animasiyası.

## Layihə məqsədi

- Loop-suz labirinti (hər pikselə yalnız 1 yol) PNG RGBA şəkli kimi oxu
- Girişdən xəzinəyə aparan yolu tap, nəticəni highlight edilmiş PNG kimi yaz
- Bonus: kəşif prosesinin GIF animasiyası

## Əsas fikirlər

### 1. Labirint = raster image (PNG RGBA)
- Divar = qara, yol = ağ, giriş = tünd göy, xəzinə = çəhrayı, həll = narıncı
- `image.RGBA` strukturu — piksel üzərində `RGBAAt(x,y)` / `Set(x,y,c)`
- JPEG YOX (lossy pozuntular), PNG — lossless
- `color.RGBA` strukturdur → const ola bilməz → `defaultPalette()` FUNKSİYA ilə
  qaytarılır (qlobal dəyişən təhlükəsindən uzaq)

### 2. PNG oxuma/latma
```go
f, _ := os.Open(imagePath)
defer f.Close()
img, err := png.Decode(f)        // image.Image interfeysi
rgbaImg, ok := img.(*image.RGBA) // type assertion: RGBA olduğuna EMİN ol
```
`image.Image.At(x,y)` yalnız `color.Color` qaytarır — tip idarəsi üçün
`*image.RGBA`-ya assert etmək lazımdır (`RGBAAt` sürətli yoldur).

### 3. Solver strukturu və API
```go
type Solver struct {
    maze           *image.RGBA
    palette        palette
    pathsToExplore chan *path
    solution       *path
    quit           chan struct{}
    // + mutex, exploredPixels kanalı
}

func (s *Solver) Solve() error          // giriş tap → kəşfi başlat
func (s *Solver) SaveSolution(path string) error  // PNG + GIF yaz
```
Struktur: `internal/solver` paketi (imagefile.go, solver.go, explore.go,
neighbours.go, path.go, animation.go, palette.go).

### 4. Linked list path (linked siyahı yolu)
```go
type path struct {
    previousStep *path      // əvvəlki addım (geriyə izləmək üçün)
    at           image.Point // bu addımın pikseli
}
```
Hər kəşfiyyat goroutine-u öz yolunu **yalnız öz növbəti addımı + pointer**
daşıyır — tam slice kopyasından ucuz.

### 5. Concurrency modeli — Daedalus/Theseus
**Məntiq:** `Solve` girişi tapır → ilk `path`-i `pathsToExplore` kanalına
göndərir → `listenToBranches` hər mesaj üçün yeni goroutine (`explore`) başladır.

**Kəsişmədə:** `explore` əvvəlki pikseli keçir, 3 qonşuya baxır:
- Yalnız 1 namizəd → özü davam edir
- 2+ namizəd → BİRİNİ özü araşdırır, qalanlarını `pathsToExplore`-a **publish** edir

**Qeyd:** Goroutine-lər başlanğıc vaxtı məlum, son vaxtı məlum deyil —
izləmək üçün WaitGroup və quit kanalı lazımdır.

### 6. Channel əsasları
- **Unbuffered channel yazmaq blocking-dir** — oxuyan yaranana qədər gözləyir
- Buffered kanal: ölçüsü = oxunuş başlamazdan əvvəl yazılacaq maksimum element sayı
- `for p := range s.pathsToExplore` — kanal bağlanana qədər oxu

### 7. Xəzinə tapıldı — dayandırma problemi
**Sadə həll (native yox):** `s.solution != nil` yoxlaması + mutex → amma
listenToBranches dayananda kanal dolu qalır → **DEADLOCK**.

**select ilə həll:**
```go
func (s *Solver) listenToBranches() {
    wg := sync.WaitGroup{}
    defer wg.Wait()
    for {
        select {
        case <-s.quit:            // xəzinə tapıldı — dayan
            return
        case p := <-s.pathsToExplore:  // yeni budaq araşdır
            wg.Add(1)
            go func(p *path) {
                defer wg.Done()
                s.explore(p)
            }(p)
        }
    }
}
```
- `select` — bir neçə channel əməliyyatından İLK hazır olanı icra edir;
  bir neçə hazır olsa RANDOM seçir
- `quit chan struct{}` — siqnal üçün boş strukturlu kanal

### 8. Loop dəyişəni tələsi (Go ≤1.22)
```go
for p := range s.pathsToExplore {
    go func() { s.explore(p) }()   // BUG: hamısı SON p görəcək (≤1.21)
}
```
**Həll:** `p := p` shadow copy və ya parametr ötürmə: `go func(p *path){...}(p)`.
(Go 1.22+ loop-var semantikası dəyişdi, amma köhnə kod üçün hələ vacibdir.)

### 9. Explored piksellər + GIF animasiyası
- Hər araşdırılan piksel `exploredPixels` kanalına göndərilir → `registerExploredPixels`
  goroutine-u onları rəngləyir və 30 frame-dən birini GIF-ə yazır
- 1 frame/piksel YOX (1000×1000 labirint = 40000 frame!) — threshold ilə
- GIF: `image.Paletted` + `palette.Plan9`; ad conflict → import alias:
  `plt "image/color/palette"`

### 10. Testlər
- Maze test case-ləri PNG kimi testdata/-da: cross (2 budaq), dead-end, iki budaq
- `explore` testi: Solver-də buffer-li kanal (3) qur → publish olunan budaq
  sayını say (`wantSize`)
- `neighbours` testi: random sıra → `cmpopts.EquateEmpty` / sort ilə determinizm

## Əsas terminlər

- Goroutine
- Channel (buffered/unbuffered)
- select
- sync.WaitGroup
- Linked List (zəncir siyahı)
- Raster Image (rastr şəkil)
- PNG RGBA / image.Paletted
- Deadlock (ölü kilidlənmə)
- Loop Variable Shadowing (loop dəyişən kölgələmə)
- Publish/Subscribe (kanal ilə)

## Praktik nəticə

- Concurrent kəşif: hər budaq = goroutine; kanal = dispatch növbəsi
- Dayandırma: quit kanalı + select + WaitGroup — 3-lük birlikdə işləyir
- Loop daxilində goroutine başladanda dəyişəni KOPYALA (`p := p` / parametr)
- Struktur rəngləri const-etmək olmaz → defaultPalette() funksiyası
- image manipulyasiyası: `png.Decode` → `*image.RGBA` assert; GIF frame-lər
  threshold ilə

## Mənbə

Pages: 393-451 (Chapter 9, Learn Go with Pocket-Sized Projects)
