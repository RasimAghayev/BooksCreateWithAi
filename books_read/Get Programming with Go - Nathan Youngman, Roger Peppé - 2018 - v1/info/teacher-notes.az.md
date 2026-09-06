# Get Programming with Go — Müəllim qeydləri (Teacher Mode)

📖 Kitab deyir:
Go-nu sıfırdan öyrən: dəyişənlər, tiplər, funksiyalar, kolleksiyalar, struct/metod/interfeys, pointer/nil/error və nəhayət goroutine/channel. Hər addımda Go-nun sadəliyi (25 açar söz, vahid format, koercion yox) və təhlükəsizliyi (bounds check, race detector, explicit conversion) vurğulanır. Mars/gopher hekayəsi ilə motivasiya.

👨‍🏫 Müəllim qeydi:
Məncə bunu praktikada belə tətbiq etmək daha yaxşı olar:

1. **Kitab 2018-də yazılıb (Go 1.9-1.10 dövrü) — bir sıra detallar yenilənib:**
   - `ioutil.ReadDir/ReadFile/WriteFile/Discard` → artıq deprecated: `os.ReadDir`, `os.ReadFile`, `io.Discard` istifadə et.
   - Error wrapping: kitabda YOX — Go 1.13-dən `fmt.Errorf("...: %w", err)` + `errors.Is/As` standartdır. `switch err { case ErrBounds: }` müqayisəsi wrapped error-larla işləmir — müasir kodda `errors.Is(err, ErrBounds)`.
   - Generics: kitabda YOX (Go 1.18-dən var) — "make zero value useful" proverbs ilə birlikdə oxu; `interface{}`-in bir çox istifadəsi artıq tip parametri ilə dəqiq yazılır.
   - Loop variable semantics: kitabın closure dərslərində (L14) ehtiyat hissəsi var — Go 1.22-dən for loop dəyişəni hər iterasiyada YENİDIR (köhnə tələsən aradan qalxdı).
   - `rand.Seed`/`rand.Intn` — Go 1.20-dən avtomatik seed; `math/rand/v2` (Go 1.22+) modern seçim.

2. **Module sistemi keçirilib:** kitabın GOPATH əsrində yazılıbsa da, bu mövzuda praktik olaraq heç nə demir — bu, kitabın güclü tərəfidir: sintaksis və semantika dərsləri yaşayır. Sadəcə öz layihələrində `go mod init` ilə başla.

3. **Test dərsləri YOXDUR** — kitabın ən böyük boşluğu. Testing-i sonradan öyrənmək lazımdır (Go in Action Ch 9 və ya Go Web Programming Ch 8 bu boşluğu doldurur). Hər capstone üçün _test.go yazmağı tapşırıq kimi verərdim.

4. **Concurrency bölməsi dərinliyi:** kitab worker/mutex/pipeline pattern-lərini beginner səviyyəsində MÜKƏMMƏL izah edir. Amma production üçün əlavələr: (a) `context.Context` — cancel/timeout üçün (kitabda yox); (b) `sync.WaitGroup` (kitabda yox — Go in Action-da var); (c) `errgroup` — paralel error yığılması. "Hər worker öz select loop-u" modeli real layihələrdə hələ də qızıl standardır.

5. **fmt.Stringer dərsi** — qiymətli; əlavə edərdim: `strconv`-in `Stringer` tətbiqi və `MarshalJSON`/`MarshalText` pattern-ləri (kitab L24-də qeyd edir amma dərinləşmir). Modern JSON: `json:"name,omitempty"` tag-ləri.

6. **Pointer dərsində "map = pointer" izahı** — texniki olaraq dəqiq deyil (map = pointer-hərf runtime hchan/hamt strukturuna), amma tədris üçün düzgün səlahiyyətləşdirmədir. Real qayda: map header kopyalanır, underlying bucket-lər paylaşılır.

7. **Sudoku capstone (L29)** — müəllimlik baxımından ən yaxşı hissələrdən biri: array + pointer + error + validation bir tapşırıqda. Həllərdə (Appendix) bir neçə düzgün yanaşma göstərilir — müqayisə etmək tələbə üçün faydalıdır.

## Ən vacib 5 fikir

1. **"Errors are values"** — kitabın L28-i bu prinsipi safeWriter ilə gətirir; müasir Go-da həmən pattern `multierr`/`errors.Join` kitabxanalarında yaşayır — amma konsepti anlamaq üçün ən yaxşı dərs budur.
2. **3-indeks slicing** — Go-nun slice modelini (pointer+len+cap) dərindən anlamağın ən qısa yolu; hər senior müsahibədə bu sualı verir.
3. **Comma-ok + nil interface tələsi (L27)** — production bug-larının 2 ən məşhur mənbəyi; kitab hər ikisini "wat?" anları ilə izah edir.
4. **Embedding = forwarding, inheritance DEYİL** — receiverin tip dəyişməməsi (L23) çoxlarının Go-ya keçərkən düştüyü first misconception-dır; kitab bunu açık aydın göstərir.
5. **Worker pattern + command channel (L31)** — mutex-siz concurrency apellyasiyası: "Don't communicate by sharing memory; share memory by communicating" proverbi kodlaşdırılıb.

## Kitabın ən dəyərli hissəsi

**Unit 4-5 (Lesson 16-24), pages 121-195** — çünki:
- Slice-in daxili modeli (len/cap, 3-indeks, append davranışı) — Go-nun ən çox istifadə olunan data strukturunun tam məntiqi.
- Composition/embedding/interface üçlüyü — Go-nun OOP alternativinin tam fəlsəfəsi; hər digər kitab bunu səthi keçir, bu kitab "bu inheritance deyil" fərqini texniki dəqiqliklə göstərir.
- Bu 2 unit Go developer-in gündəlik kodunun 90%-nin əsasıdır.

## Uyğunsuzluq / diqqət qeydləri

⚠️ **Listing 27.13 (number struct):** "valid" pattern üçün gözəl nümunədir — amma müasir Go-da `sql.Null*` tipləri və ya generic `Optional[T]` (kitabxanalar) eyni məqsədə xidmət edir; yalnız öz valid-ını yazmaq yeganə yol deyil.

⚠️ **Lesson 30 pipeline:** sentinel-dən close-a keçid mümkün dərslər ardıcıllığı kimi göstərilir — real kodda HƏMİŞƏ close istifadə et; sentinel variantı yalnız tarixi maraqdır.

⚠️ **Lesson 31 "nil channel blocks forever"** — doğru; amma `close(nil channel)` → panic də əlavə edilib, tamamdır. Çox tutorial bunu unudur — kitab yadırkən qeyd edir. ✓

⚠️ **Bütün rəqəmlər Playground-la bağlı:** "same pseudorandom numbers" — Playground-da time dondurulub; lokalda `rand.Seed(time.Now().UnixNano())` (kitab dövrü) / `rand.New(rand.NewSource(...))` lazım idi — Go 1.20+ avtomatikdir. Chapter oxuyarkən bunu bilmək çətinliyi aradan qaldırır.

⚠️ **Go version-sensitive məqamların siyahısı** (müəllim üçün): loop var (1.22), rand seed (1.20), errors.Is (1.13), generics (1.18), ioutil deprecation (1.16). Tələbəyə müasir ekvivalentləri paralel göstər.
