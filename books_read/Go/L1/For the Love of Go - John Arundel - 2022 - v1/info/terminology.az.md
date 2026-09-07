# For the Love of Go — Terminologiya (Azərbaycanca)

Format: **English Term (Azərbaycanca qarşılıq)** — qısa izah

## Layihə və Alətlər (Ch1-3)
- **Module (modul)** — kod vahidi; go.mod ilə müəyyən; modul içində modul YOX
- **go mod init** — yeni modul başlat
- **go test / go run / go build** — testləri/ proqramı / binary-nı işə sal
- **gofmt (-d / -w)** — format diff / yerində düzəlt
- **gofumpt** — daha sıxı gofmt variantı
- **cmd/ Pattern** — executable üçün konvensional qovluq
- **GOOS / GOARCH ("goose"/"gorch")** — cross-compilation hədəf OS/CPU
- **go tool dist list** — mövcud GOOS/GOARCH kombinasiyaları
- **os.Exit(0/1)** — statuslu çıxış (0 = OK konvensiyası)

## Test Əsasları (Ch1-3, 5)
- **Want and Got Pattern** — gözlənti/nəticə müqayisə test strukturu
- **Null Implementation** — `return 0` / `Book{}` — testin özünü yoxlama
- **Table Test** — testCase struct + testCases slice + range
- **"One Behaviour, One Test"** — hər test tək davranış qanunu
- **Valid / Invalid Input Testləri** — 2 davranış = 2 test
- **t.Parallel / t.Error / t.Errorf / t.Fatal / t.Fatalf** — test idarəçiliyi
- **Red, Green, Refactor** — TDD tsikli
- **Compile-Only Test** — `_ = T{...}` — struct yaratma testi
- **Covered vs Tested** — icra ≠ yoxlama (TestBuyTrivial)
- **Coverage Ratchet** — coverage-i SALAN dəyişiklik qadağası

## Tiplər və Data (Ch4, 6-8)
- **String / int / float64 / bool** — əsas tiplər
- **Value / Variable / Type** — data / adlı yer / "kind of thing"
- **Zero Value** — default: 0 / "" / false / boş struct
- **var vs `:=`** — zero dəyər / başlanğıc dəyərlə
- **Literal** — string/int/struct/slice/map/function literal-ları
- **Composite Type** — çoxparçalı vahid
- **Struct / Field** — strukturlaşmış qeyd / sahələr
- **Type Definition** — `type Ad struct {...}`
- **Exported / Unexported** — Böyük hərf public / kiçik private
- **Slice ([]T)** — eyni tipli sıralı toplusu; index, len, append
- **Map (map[K]V)** — key→element; lookup, comma-ok; sırası RANDOM
- **Comma-Ok (b, ok := m[k])** — mövcudluq yoxlaması idiomu
- **go-cmp / cmp.Equal / cmp.Diff** — struktur müqayisəsi / element diff
- **cmpopts.IgnoreUnexported** — unexported sahəli müqayisə üçün

## Metodlar və Pointerlər (Ch8-10)
- **Method** — receiver-lü funksiya
- **Receiver** — `func (b Book)` — xüsusi ilk parametr
- **Pointer Method / Receiver** — dəyişən metod mütləq pointer alır
- **Value Receiver Tələsi** — kopya dəyişir; SA4005 staticcheck xəbərdarlığı
- **Automatic Dereferencing** — b.Field → (*b).Field öz-özünə
- **Non-Local Type** — paketin olmayan tipə metod QADAĞA
- **Custom Type (type MyInt int)** — underlying əsaslı LOCAL tip
- **Underlying Type** — yeni tipin əsası; MyInt ≠ int (DISTINCT)
- **Type Conversion** — `mytypes.MyInt(9)` — TipName(dəyər)
- **Wrapping** — struct daxilində underlying tip sahəsi (Contents)
- **Metod MİRASI YOXDUR** — type X Y underlying metodları GƏTİRMİZ

## Validasiya (Ch10)
- **Validating Accessor** — yoxlayan setter (SetPriceCents)
- **Always Valid Field** — unexported + getter/setter
- **Always Valid Struct** — unexported tip + constructor
- **Constructor (New)** — yalnız VALID dəyərlər yaradan funksiya
- **Map-as-Set** — map[T]bool; missing key → false
- **Constants** — compiler yoxlamalı adlar (http.StatusOK)
- **iota** — 0,1,2... avtomatik konstant enumerator
- **type Category int** — mənalı int tipi

## Control Flow (Ch11-12)
- **Statement / Declaration / Assignment** — əmr / yer ayırma / dəyər qoyma
- **Tuple Assignment** — a, b, c := 1, 2, 3 — əlaqəli qrup
- **Blank Identifier (_)** — "at" yer tutucusu
- **Happy Path** — hər şey normal olandakı axın
- **Left-Align Happy Path** — minimum indent prinsipi
- **Flip Şərtlər** — mümkünsüzləri əvvələ al (x <= 0 → return)
- **Early Return** — if + return = else-sizlik (Go tərzi)
- **&& / || / !** — and / inklüziv or / not
- **switch / case / default** — N yollu seçim; ilk-match qaydası
- **fallthrough** — növbəti case-i də icra (nadir)
- **Switch Expression** — switch x { case 1, 2, 3: }
- **for / range** — conditional / forever / collection loop-ları
- **continue / break** — element ötür / loop-dan çıx
- **Label (outer:) / goto** — nested nəzarət (nadir) / birbaşa jump (avoid)

## Funksiyalar (Ch13)
- **Signature** — parametr + result tip kombinasiyası
- **Functions Are Values** — assign / ötür / qaytar bilinən dəyərlər
- **func(float64, float64) float64** — funksiya TİPİ yazılışı
- **Function Literal** — adsız, on-the-fly funksiya
- **Closure** — definisiya skopunu görən literal ("closure over X")
- **Bubble Analogy** — closure dəyişənləri həbs edib gəzidir
- **Loop Variable Tələsi** — closure ÇAĞIRIŞ vaxtını görür → 3,3,3
- **defer** — "çıxanda icra et, necə çıxırsa"
- **Stacking Defers (LIFO)** — son defer əvvəl icra
- **Named Result Parameters** — latitude/longitude; sənədləşdirmə dəyəri
- **Naked Return** — boş return; LEGAL amma harmful — explicit yaz
- **Deferred Closure** — err = closeErr; nəticəni çıxışdan sonra dəyiş
- **Variadic (...T)** — istənilən sayda arqument; daxildə slice

## Proqram Yaşam Dövrü (Ch14)
- **Executable Binary** — machine code + OS formatı; ~2 MiB asılılıqsız
- **package main / func main()** — executable tələbi / giriş nöqtəsi
- **init** — main-dən əvvəl magic funksiya — POOR STYLE
- **Package-Level var** — `var x = initFn()` — init əvəzi
- **Compiler** — mənbə → machine code
- **Mach-O / PE32+ / ELF** — macOS / Windows / Linux binary formatları
- **log.Fatal / panic** — dərhal ləğv; dərinlərdə QADAĞA
- **"Yalnız main-dən Exit"** — görünən çıxış qaydası

## Tao Fəlsəfəsi (Ch15)
- **Tao** — şeylərin daxili təbiəti (su aşağı axar)
- **"Thrashing vs Surfing"** — mübarizə / dalğadan istifadə
- **Kindness** — insanlar üçün kod (4 qrup)
- **Deep Abstraction** — kiçik API, güclü maşın
- **Micro-Improvements** — hər ziyarətdə kiçik refaktor
- **"Spaghetti Təmirsizdir"** — Ousterhout qanunu
- **Simplicity / Frugality** — azla çox; "Simplify, simplify!"
- **Extensibility Tələsi** — lazımsız genişlənmə üçün sadəliyi pozmaq
- **Humility** — óbvio > clever; pre-engineering YOX
- **Liddell Hart Qanunu** — öz xəta meylimizi görməmək = ən təhlükəli
- **Wúwéi (Not Striving)** — zorlamamaq (tənbəllik DEYİL)
- **Problem-Eliminating** — problemi ləğv etmək > həll etmək
- **"Ən Yaxşı Optimallaşdırma"** — işi HEÇ ETMƏMƏK
- **Programming ≠ Typing** — fikir klaviaturadan əvvəl
- **Buffalo Qanunu (Weinberg)** — istiqaməti istək müəyyən edir
- **Dàodé Jīng** — üç xəzinə mənbəyi
