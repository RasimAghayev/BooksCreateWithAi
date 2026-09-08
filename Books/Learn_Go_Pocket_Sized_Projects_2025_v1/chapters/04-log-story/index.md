# Chapter 4 — A log story: Creating a library (səh. 122-153)

## Bu fəsil nədən bəhs edir?

İlk KİTABXANA (paket) layihəsi — pocketlog: 3-səviyyəli logger. Logging anlayışı
(gəmi logbook/chip log etimologiyası; rubric — qırmızı yazı), library = paket
(exported API; backward compatibiliy — public DƏYİŞMƏZ, private AZAD), paket
qaydaları (qovluq adı = paket adı; bir sözlük, lowercase; fayl bölgüsü —
"kiçik fayllardan qorxma" — level.go ayrı), Level enum-u (type Level byte —
int32-dən 4x kiçik!; iota — hər sətirdə artır; SIRALAMA vacibdir; sənədlər
"// Level represents..." adla BAŞLAYIR), 3 səviyyə (Debug — prod-da YOX;
Info — milestone-lar; Error — xəta+rekoveri), receiver metodları (func (l
*Logger) Debugf — pointer receiver; l.Log(Info,...) YOX, l.Info(...) —
oxunaqlılıq), variadic (args ...any; args[2] slice kimi; f hərfi = Printf
imzasına uyğun), New() konstruktoru (Go-da ctor YOXDUR — konvensiya;
pocketlog.New — stuttering QADAĞASI: NewPocketLog YOX; threshold UNEXPORTED —
"şübhələnəndə export ETMƏ"; pointer qaytarır — paylaşma asanlığı), doc.go
paket başlığı ("Package pocketlog..." — böyük P linter üçün), go doc (lokal;
go mod download; paket/simvol səviyyəsində), **closed-box testing** (paket_x
_test — İSTİFADƏÇİ nöqteyi-nəzəri; foo_test + foo EYNİ qovluqda İSTİSNADİR;
ExampleLogger_Debugf adı), io.Writer/Reader interfeysləri (Write(p []byte)
(n int, err error); implicit implementasiya; istənilən destination — stdout/
fayl/şəbəkə/printer), happy-path alignment (guard-clause: if threshold >
LevelDebug { return } — biznes məntiqi SOLDA), nil-qoruma (output == nil →
os.Stdout), _ = açar sözü (bilərəkdən rədd — "əlimdə var amma lazım deyil"),
refaktorinq logf (unexported; üç səviyyə eyni yazma kodu ÇAĞIRIR — "eyni
kodu 2 dəfədən çox saxlama"), **functional options pattern** (Option
func(*Logger); WithOutput(io.Writer) Option; New(threshold, opts ...Option)
— default + istəyə bağlı; yeni option = yeni funksiya, API DƏYİŞMİR),
testWriter mock (io.Writer implementasiya — contents yığır; strings.Builder/
bytes.Buffer praktikada; "interfeys standart deyilsə mock belə yazılır"),
TDT ilə 3 metodu bir testdə (səviyyə threshold-lar; %q ilə space yoxlaması),
exported Logf(lvl Level, ...) (runtime-da səviyyə seçimi — GDPR ssenarisi),
logging best practices (aydın mesaj — "step 1" YOX; QISA — map-in özü YOX
uzunluğu; milestone-larda — loop-da hər 10 000-də bir; log ≠ debug aləti —
kodua aydınlat/sənədləşdir/test; structured JSON loglar —{"time","level",
"message"}; mashine-readable dashboards).

## Əsas fikirlər

### 1. Paket = Kitabxana; API = Export Edilən Hissə
- **Paket qaydaları:** eyni qovluq + eyni package adı; qovluq adı = paket
  adı (bir söz, lowercase); faylları BÖL (level.go ayrı — "açılanda level
  axtarırsansa level.go-da")
- **Stability:** exported (böyük hərf) = MÜQAVİLƏ — versiyalar arası DƏYİŞMƏZ;
  unexported = azad dəyişə bilər
- **doc.go:** "Package pocketlog..." — paket README-si; böyük P linter
  tələbi

### 2. Level Enum — iota
```go
// Level represents an available logging level.
type Level byte          // int32 YOX — 4x yaddaş qənaəti
const (
    // LevelDebug represents the lowest level of log...
    LevelDebug Level = iota     // 0
    LevelInfo                    // 1
    LevelError                   // 2
)
```
- **iota:** hər sətirdə avtomatik artır; İLK dəyərdə tip göstərilir;
  **SIRA = önəm sırası** (müqayisə < bunun üzərində qurulur!)
- **Sənədlər:** hər exported simvol — şərh ADLA başlayır + nöqtə ilə bitir
- Yeni səviyyə = bir sətir (renumbering YOX)

### 3. Metod Dizaynı — l.Infof(...) > l.Log(Info, ...)
```go
func (l *Logger) Debugf(format string, args ...any) {
    if l.threshold > LevelDebug {    // guard: happy-path SOLDA
        return
    }
    l.logf(format, args...)
}
func (l *Logger) logf(format string, args ...any) {   // unexported core
    _, _ = fmt.Fprintf(l.output, format+"\n", args...)
}
```
- **Hər səviyyə = ayrıca METOD** (oxunaqlı çağırış; f = Printf imzası);
  hamısı logf-ı ÇAĞIRIR — dekorasiya dəyişərsə 1 yer
- **Guard clause üslubu:** şərt pozulursa RETURN — məntiqə indent YOX
- **_, _ =** — n/err bilərəkdən RƏDD ("bilirəm qayıdır, lazım deyil")
- **nil qoruması:** New istifadə olunmasa (zero Logger) → output nil →
  os.Stdout əvəzlə

### 4. New() Konstruktor Konvensiyası
```go
func New(threshold Level, opts ...Option) *Logger {
    lgr := &Logger{threshold: threshold, output: os.Stdout}
    for _, configFunc := range opts { configFunc(lgr) }
    return lgr
}
```
- **pocketlog.New** (NewPocketLog YOX — stuttering); məcburi parametrləri
  TƏTBİQ etmək üçün yeganə yol
- **Pointer qaytarır:** paylaşma (dependency injection) asanlaşır
- **threshold UNEXPORTED:** "şübhələndikdə export ETMƏ" — istifadəçi daxili
  strukturu bilə BİLMƏMƏLİ

### 5. io.Writer — İstənilən Çıxış
```go
type Writer interface { Write(p []byte) (n int, err error) }
```
- **Implicit implementasiya:** Write metodu VAR = Writer-ə UYĞUN (elan YOX)
- Fərdi çıxışlar: stdout, fayl, DB driver, printer, şəbəkə — hamısı eyni
  interfeys

### 6. Functional Options Pattern
```go
// Option defines a functional option to our logger.
type Option func(*Logger)
// WithOutput returns a configuration function...
func WithOutput(output io.Writer) Option {
    return func(lgr *Logger) { lgr.output = output }
}
lgr := pocketlog.New(pocketlog.LevelInfo, pocketlog.WithOutput(os.Stdout))
```
- **Default + özelleştirme:** New THRESHOLD məcburi (stdout default);
  opts variadic — sıfır və ya daha çox
- **Genişlənmə API-ni POZMUR:** yeni option = yeni WithXxx funksiyası;
  çağırış köhnə kodda DƏYİŞMİR
- **Kitabın prinsipi:** default-lar "ehtiyatla istifadə olunsa YAXŞI"

### 7. testWriter — Mock Əl ilə
```go
type testWriter struct{ contents string }
func (tw *testWriter) Write(p []byte) (n int, err error) {
    tw.contents += string(p)
    return len(p), nil
}
// istifadə: tw := &testWriter{}; New(level, WithOutput(tw)); ...; tw.contents YOXL
```
- **Mock dərsi:** standart interfeys YOXdursa BELƏ — kiçik struct + Write;
  praktikada strings.Builder/bytes.Buffer
- **Fayda:** stdout testdən ASILI deyil; ardıcıl çağırışların CƏMİ yoxlanılır

### 8. Closed-Box Test — paketfoo_test
```go
// logger_test.go:  package pocketlog_test   ← fərqli paket!
import "learngo-pockets/logger/pocketlog"
func ExampleLogger_Debugf() { ... }
func TestLogger_DebugfInfofErrorf(t *testing.T) {
    tt := map[string]testCase{ "debug": {LevelDebug, 3 mesaj}, "info": {...}, "error": {...} }
    for name, tc := range tt {
        t.Run(name, ...Debugf+Infof+Errorf; tw.contents == tc.expected...)
    }
}
```
- **foo + foo_test qovluqda birgə — YALNIZ bu istisna**; gizli simvollara
  giriş YOXDUR = istifadəçi nöqteyi-nəzəri
- **Ad formatı:** Example{Tip}_{Metod}[_ssena]

### 9. Exported Logf — Runtime Səviyyə
```go
func (l *Logger) Logf(lvl Level, format string, args ...any) {
    if l.threshold > lvl { return }
    l.logf(format, args...)
}
```
- **Ssenari:** GDPR — bir platformada email loglanır, digərində YOX;
  konfiqurasiyadan gəlir; "hamısı bunu çağırsın" təklifi — amma API
  qarışır; 2 variantın olması = clutter

### 10. Logging Best Practices
| Qayda | Nümunə |
|---|---|
| AYDIN mesaj | "step 1" YOX — "sənəd DB-yə yazıldı (took 3ms)" |
| QISA | minlərlə açarlı map-i YOX — len(map)/contains yaz |
| Milestone-larda | loop-da hər element YOX — hər 10 000-də bir |
| Log ≠ debug aləti | kod aydın DEYİLsə: funksiya böl / sənəd yaz / test yaz |
| Structured (JSON) | {"time","level","message"} — maşın oxunur; dashboards |
| Qiymət | log saxlama PULdur — həcm/aydınlıq kompromisi |

## Əsas terminlər
- API — paketin exported səthi; backward compatibility müqaviləsi
- doc.go — paket başlığı; "Package X..." böyük P
- Enum/iota — avtomatik artan konstant sırası; sıra = önəm
- Receiver method — (l *Logger) ilə bağlı funksiya
- Variadic (args ...any) — 0..N arqument slice kimi
- New() — konvensiya konstruktoru; stuttering-dən qaçınma
- Functional options — Option func(*Logger) + WithXxx + variadic opts
- io.Writer/Reader — yaz/oxu interfeysləri (implicit)
- Guard clause — əvvəl qayıt; happy path solda
- _, _ = — bilərəkdən dəyər rəddi
- Closed-box (external) test — {pkg}_test paketi
- testWriter/mock — interfeysin yaddaş yığan implementasiyası
- Structured logging — JSON formatlı maşın-oxunarlı loglar
- Stuttering — pocketlog.NewPocketLog təkrarı (antipattern)

## Praktik nəticə

1. **Kitabxana şablonu:** qovluq/paket + doc.go + logger.go + level.go +
   (gizli sahələr) + New(məcburi, opts...) + Example + closed-box test.
2. **Enum:** type X byte/iota; sənədlər adla; sıra müqayisə semantikasıdır.
3. **Metod strukturu:** səviyyə metodları = guard + logf; hamısı TİK
   logf-ı çağırır — dekorasiya 1 yerdə.
4. **Genişlənən konfiqurasiya:** functional options — default stdout +
   WithOutput + gələcək WithXxx-lər API-ni POZMUR.
5. **Test mock əl ilə:** testWriter (strings.Builder/bytes.Buffer kimi) —
   stdout Examplesini tam TestXxx ilə əvəz et.
6. **Export qərarı:** şübhə varsa UNEXPORTED; sənədləşdirmə və closed-box
   test istifadəçi gözləntilərini KİLİDLƏYİR.
7. **Log yazarkən:** aydın + qısa + milestone + strukturlu; log debugging
   VASİTƏSİ deyil, kod keyfiyyəti göstəricisidir.

## Mənbə
Pages: 122-153 (PDF 123-154)
