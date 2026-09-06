# Effective Go (RU) — Müəllim Qeydləri (Teacher Notes)

> Bu sənəd Effective Go-nu öyrədərkən istifadə üçün metodiki qeydləri, çətin anları, müzakirə suallarını və tapşırıqları birləşdirir.

---

## 1. Auditoriya və ön şərtlər

- **Hədəf auditoriya:** Go sintaksisini bilən, amma "idiomatik necə" sualı ilə üz-üzə olan developer-lər. Rus dilini oxuyan AZ-tələbələr üçün ikili fayda: dil + texniki rus.
- **Ön şərtlər:** A Tour of Go + How to Write Go Code (sənədin öz ön-şərti).
- **Format:** 45 səhifə — 4-6 saatlıq intensiv və ya 9 module-luq kurs. 100% practical kod nümunəli.

---

## 2. Tədris axını

| Module | Mövzu | Vaxt |
|--------|-------|------|
| 1 | Format + adlar + semicolon | 10% |
| 2 | if/for/switch/type switch | 12% |
| 3 | Çoxlu qayıdış + named + defer | 12% |
| 4 | new/make/slice/map/çap | 16% |
| 5 | const/init + metodlar | 12% |
| 6 | İnterfeyslər + çevirmə | 15% |
| 7 | `_` + embedding | 8% |
| 8 | Kanallar + şüar | 15% |
| 9 | Xətalar + panic/recover + veb | 10% |

---

## 3. Çətin anlar və izah üsulları

### a) new vs make (ən klassik qarışıqlıq)
Dialoq: `new([]int)` nə qaytarır? → nil-slice-a pointer — İSTİFADƏSİZ! `make([]int, 100)` → işlək slice. Cədvəl çək: new=zero+`*T`, make=init+`T`, make yalnız 3 tip. "Zero value hazırdır" fəlsəfəsini SyncedBuffer ilə göstər — amma diqqət: interfeyslərin zero value-su nil-dir (Olan mövzu burada yoxdur, amma sual gəlir).

### b) Lekserin `;` qaydası
"Niya `{` növbəti sətirdə OLMUR?" — C-dən gələnlərin 1 nömrəli xətası. Whiteboard: `if i < f()` sətir sonu → `;` daxil → `{` artıq if-in DEYİL. Rule: "sətir sonundaki token statement BİTİRƏ BİLİRSƏ → `;`".

### c) Massiv = VALUE
C gözləntisini qır: funksiyaya massiv KOPYA gedir. `Sum(a *[3]float64)` pointer həlli — amma sənəd özü deyir: idiomatik deyil, SLICE istifadə et. 2 dərs bir yerdə: massiv → slice keçidi.

### d) `:=` redeclaration
"d, err := f.Stat()" — err haradadır? Qaydaları 3 bənddə ver: eyni scope + təyin edilə bilən + YENİ var da var. Bu CH2 #1 (shadowing) səhvinin MÜSBƏT tərəfi — eyni scope-da shadow YOXDUR.

### e) Embedding vs inheritance
Receiver FƏRQİ ən çox sual yaradan yer: `bufio.ReadWriter.Read` çağrılınca `rw.reader` receiver-dir, rw YOX. Java-this zehniyyətindən qopma anı. Job.Logger sahə-çıxışı göstər.

### f) Kanal sinxronizasiyası
"Unbuffered = kommunikasiya + sinxron" — mutexsiz dünyanın açarı. Semafor + worker pool + resultChan 3-lüsünü ardıcıl qur: (1) limit, (2) resurs, (3) cavab axını. Serve loop-var bug-unu mütləq qeyd et (Go <1.22).

### g) Panic/recover paket idiomu
Compile nümunəsi defer DAXİLİNDƏ named result dəyişir — bir neçə konsepsiya bir yerdə (defer + named + recover + assertion + re-panic). Addım-addım izlə: parse error → panic(Error) → unwind → defer → recover → err qayıdışı.

---

## 4. Müzakirə sualları

1. "Zero value hazırdır" nə vaxt OLA BİLMƏZ? (handle/ID sahələr — zero işlək deyil)
2. `ring.New` gözəldir — amma 2 tip olsa? (paket bölgüsü prinsipi)
3. Fmt.Println(seq) nə çap edir və NİYƏ? (Stringer — Sequence nümunəsi)
4. Semafor worker-pool-dan nə vaxt üstündür? (request ölçüsü bəlli deyilsə)
5. Panic-i kitabxanada NEÇƏƏ yerdə görmək olar? (init asılılığı — CubeRoot toy nümunəsi real kitabxanada olmaz)
6. `var _ json.Marshaler = (*T)(nil)` niyə VACİBDİR və niyə hər tip üçün YOX?
7. Job.Logger embed-i hansı hallarda təhlükəlidir? (ad toqquşması + Logger API açıq)

---

## 5. Praktik tapşırıqlar

**Tapşırıq 1:** 3 fərqli tipdə (struct, int, func) Handler yazın — `/counter`-ı int-receiver ilə qurun.

**Tapşırıq 2:** SyncedBuffer kimi zero-value-hazır tip DİZAYN edin (id generator, cache statləri) — new/var-sonsuz işləsin.

**Tapşırıq 3:** Effective Append: manual Append yazın → builtin append ilə müqayisə → `x = append(x, y...)` sınayın.

**Tapşırıq 4:** ByteSize kimi tip + String metodu (rekursiya tələsini bilərəkdən yaradın → düzəldin); 2D slice iki üsulla (sətir-sətir / tək array).

**Tapşırıq 5:** ReadFull named-return; nextInt çoxlu-qayıdış; trace/un defer patterni.

**Tapşırıq 6:** Serve: semafor versiyası → worker pool → resultChan RPC — 3 mərhələli inkişaf, hər birində `-race` ilə test.

**Tapşırıq 7:** `var _ Iface = (*T)(nil)` — json.Marshaler tipli layihədə kompayl yoxlaması qırın və xətanı görün.

**Tapşırıq 8:** SafelyDo + Compile panic/recover patternlərini yazın; panic-i故意 qaldırıb izləyin (unwind sırası).

**Tapşırıq 9:** QR-serveri (chart API ölçüb — lokal imgservis və ya bəsit echo ilə əvəz edin) — flag + template + HandlerFunc tam zənciri.

---

## 6. Sınav sualları

**Asan:**
1. `make([]int, 5)` vs `new([]int)` fərqi?
2. `for pos, char := range "日"` pos nədir? (0 — bayt pozisiyası)
3. Getter adı `owner` sahəsi üçün? (Owner)

**Orta:**
4. `d, err := f.Stat()` — err yeni elandır mı? (Xeyr — reassign)
5. Unbuffered kanalın "sinxronizasiya" təbiəti nə deməkdir?
6. Value-metodu pointer-də necə çağrılır? (avtomatik dereference)

**Çətin:**
7. Embed-də metodun receiver-i kimdir və bu irsdən fərqin harasıdır?
8. Compile funksiyasında defer daxilində `err = e.(Error)` uğursuz olsa nə baş verir? (re-panic → unwind davam)
9. Niyə `serve`-in her-request-goroutine versiyası resurs baxımından problemli idi və worker pool bunu necə həll etdi?

---

## 7. Kollektiv layihə ideyası

**"Effective Idioms Linter"** — tələbələr Effective Go qaydalarından 10 maddəlik "idiom linter" siyahısı hazırlayır (GetOwner ✓/✗, _ = err ✓/✗, `;` növbəti sətir `{|` ✓/✗, new([]T) ✓/✗ ...) və mövcud open-source Go repo-larında manual audit aparır. Nəticə: hər tapıntı üçün sənəd referansı + düzəliş.

---

## 8. Əlavə resurslar

- Orijinal (ingilis): https://go.dev/doc/effective_go
- Issue 28782 (yenilənmə müzakirəsi): golang/go
- A Tour of Go: https://go.dev/tour
- How to Write Go Code: https://go.dev/doc/code
- Go wiki — loop variable pitfall (Serve nümunəsi)
- Go Code Review Comments (birbaşa davamçı): https://go.dev/wiki/CodeReviewComments
- The Go Programming Language (Donovan/Kernighan) — dərinləşdirmə
