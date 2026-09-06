# Get Programming with Go — Xülasə (Azərbaycan dilində)

**Müəlliflər:** Nathan Youngman, Roger Peppé
**Nəşriyyat:** Manning Publications Co., 2018 · **ISBN:** 9781617293092
**Səhifə:** 284 (+həllər + indeks) · **7 Unit · 32 Lesson** (capstone-larla)

## Kitabın ümumi məqsədi

Proqramlaşdırmaya tam yeni başlayanlardan başlayaraq Go-nun bütün əsasını addım-addım öyrədir: dəyişənlərdən concurrency-yə qədər. Hər lesson "Consider this" düşündürücü sualı, "Quick check" mini-quizləri, "Experiment" tapşırıqları və hər unit-in sonunda capstone layihəsi ilə tədris strukturu daşıyır. Mars və gopher tematikası — kitabın davamlı motivasiya çərçivəsidir.

## Unit-by-unit xülasə

### Unit 0 — Getting Started (L1)
Go nədir: kompilyasiya olunan + sadəlik + efficiency (Rob Pike sitatları). Go Playground (play.golang.org) — Run/Format/Share. package/import/func — ilk proqram. One true brace style (kompilyator avtomatik `;` əlavə edir). 25 açar sözündən ilk 3-ü.

### Unit 1 — Imperative Programming (L2-5)
Print/Println/Printf + format verb-lər (%v, align, padding). const/var, modulus, increment operatorları. rand.Intn + off-by-one. bool (yalnız true/false — "" və 0 false DEYİL!), müqayisə operatorları. if/else if/else, logical operatorlar (short-circuit), switch (vergüllü case-lər, şərtsiz forma, eksplisit fallthrough), for (Go-da yeganə loop!). Scope (package/function/block; case-lər brace-siz scope açır!), short declaration `:=` (if/switch/for daxilində), dar-geniş scope trade-off + refactoring. **Capstone:** Mars biletləri cədvəli.

### Unit 2 — Types (L6-11)
float32/64, zero value, %f width/precision. Float dəqiqliyi: 0.1+0.2≠0.3, müqayisə tolerance (math.Abs), vurma-bölmə sırası, pul üçün int sent-lər. 10 integer tipi, wrap-around, uint8 CSS rənglər, hex/binar, int64 ilə 2038 problemi. big package (Int/Float/Rat, SetString, Div). Untyped constants — compile-time big hesablamalar (24 kvintillion!). String-lər (immutable, raw string, rune/byte, UTF-8, range-ile-decode), Caesar/ROT13 şifrələri. Konvertasiya: eksplisit, truncate, Arianne 5 faciəsi (range yoxlama!), strconv Itoa/Atoi, bool-un yalnız if ilə. **Capstone:** Vigenère cipher (keyword + modulus).

### Unit 3 — Building Blocks (L12-15)
Funksiya anatomiyası: parametr/arqument, export, çoxlu nəticə (adlı/adsız), variadic + empty interface (Println-in sirri). Pass by value — funksiyalar izolyasiya olunur. Yeni tiplər (`type celsius float64`) — vahid qarışdırmamaq üçün (celsius≠fahrenheit). Metodlar (receiver, dot notation, simmetrik adlar: hər tipin öz celsius() metodu). First-class funksiyalar: dəyişənə mənimsətmə (sensor), parametr (measureTemperature), funksiya tipləri (type sensor), anonymous funksiyalar + closure (calibrate — s və offset closure-da yaşayır). **Capstone:** temperatur cədvəlləri (drawTable + first-class funksiya).

### Unit 4 — Collections (L16-20)
Array: sabit uzunluq = tip, kopyalanır, [5]≠[8], bounds panic, composite literal, range. Slice: array-ə pəncərə, paylaşır, default indekslər, string slicing (bayt!). append + len/cap (2x böyümə), **3-indeks slicing** (Jupiter-i qoru!), make preallokasiya, variadic funksiyalar (`planets...`). Map: key-value, zero-value lookup, comma-ok, delete, make; **kopyalanmır**; sayğac/qruplama (map of slices)/set (map[T]bool) patternləri; sıra qeyri-müəyyən. **Capstone:** Conway's Game of Life (Universe [][]bool, wrap-around modulus, double buffering).

### Unit 5 — State and Behavior (L21-25)
Struct: sahələr, 2 literal forması (field-value davamlı/values-only kövrək), %+v, kopya semantikası, slice-of-structs, JSON + struct tag (sahələr EXPORTED!). Metodlar struct-larla: DMS→decimal (coordinate), constructor konvensiyası (newLocation, errors.New — tək New). Class alternativi: world + distance metodu (radius sahə). Kompozisiya (report = sol + temperature + location), embedding (sahə adı YOX — metodlar+sahələr promote, name collision → ambiguous selector → öz metod priority). **Bu inheritance DEYİL** — receiver həmişə daxili tip; Gang of Four "favor composition". Interfeyslər: talk (martian+laser), implicit satisfaction, `-er` konvensiya, shout(t talker); embedding interfeysi təmin edir (starship→laser→talker); **stardater**: time.Time-ə sonradan interfeys tətbiqi; fmt.Stringer, io.Reader/Writer, json.Marshaler. **Capstone:** Mars heyvan sığınacağı (day/night cycle, Stringer, move/eat).

### Unit 6 — Down the Gopher Hole (L26-29)
Pointer: `&`/`*`, pointer tipləri, NASA administrator misal dəyişənliyi, struct-larla ergonomiya (avtomatik dereference, `&person{...}` literal), interior pointer (`&player.stats`), mutation (pointer parametr/receiver; nathan.birthday() — Go avtomatik `&` edir), time.Time immutable counter-örnek. Gizli pointer-lər: map = pointer; slice = ptr+len+cap; *[]string yalnız slice-in ÖZÜNÜ dəyişmək üçün. Pointer receiver + interface → yalnız pointer satisfy edir. Nil: billion dollar mistake, panic, guard clause (nil receiver metod çağrılır!), nil funksiya (default pattern), nil slice (append OK), nil map (oxu OK yazı panic), **nil interface tələsi** (`(*int)(nil)` — tip+dəyər), nil-alternativ (valid bool struct). Error: çoxqayıdış konvensiyası, defer, **safeWriter** ("errors are values" — Rob Pike), errors.New + Err dəyişənləri (ünvan müqayisəsi), custom error (SudokuError []error, Error() string), type assertion, panic/recover (yalnız defer-də; panic > os.Exit). **Capstone:** Sudoku qaydaları (array pointer + validasiya).

### Unit 7 — Concurrent Programming (L30-32)
Goroutine (`go`), time sharing, ixtiyari sıra, arqumentlər kopya. Channel: make/send/receive, blocking, deadlock, nil channel. select + time.After (timeout pattern). Pipeline: source→filter→print, sentinel → **close + range** (idiomatik). Concurrent state: race condition (shared phone), race detector, mutex (Lock/defer Unlock, struct sahəsi, sadə kod tək state), pitfalls (deadlock təkrar Lock, unlock panic). Long-lived worker: `for { select {} }`, RoverDriver (commandc kanal metodlar arxasında; direction TƏK goroutine-də → mutexsiz safe!). **Capstone:** Mars grid + rover-lər + life axtarışı + radio buffer.

## Kitabın əsas mesajları

1. **Go sadəliyi dəbdən deyil, dizayndandır** — 25 açar söz, vahid format (gofmt), minimal konsept sayı.
2. **Statik tiplər səni qoruyur** — vahid qarışdırma (celsius/fahrenheit), koercion yoxdur, xətalar kompilyasiyada.
3. **Funksiyalar/struct/interfeys klassik OOP-nin əvəzidir** — inheritance YOX, kompozisiya + implicit interfeys.
4. **Errors are values** — exception yox, dəyərlər; creative pattern-lər (safeWriter) var.
5. **Concurrency = go + channel + select** — pipeline, worker, timeout vahid primitivlərdən qurulur.
6. **Test edilən kod = dizayn edilmiş kod** — hər lesson "Consider this" ilə düşündürür.

## Ən dəyərli 5 fikir (az oxucu üçün)

1. **3-indeks slicing** (`s[i:j:j]`) — append-in alt array-i səmimi korlamamasının qarantiyası; kitabın ən praktik "aha"-sı.
2. **Comma-ok idiomu** — map-də "yoxdur" ilə "dəyəri 0-dur" fərqi; hər lookup-da ağılda saxla.
3. **Untyped constants compile-time big-dir** — `const d = 24e18` işləyir, çünki kompilyator big package istifadə edir.
4. **Nil interface tələsi** — interface tip+dəyər cütüdür; `(*int)(nil)` nil-ə bərabər DEYİL. Bugün belə senior-ları yıxır.
5. **Mutex yerinə worker** — dəyəri tək goroutine-ə ver, komandaları kanalla göndər (RoverDriver) — race olana yoxdur, mutexə ehtiyac yoxdur.

## Kitabın ən dəyərli hissəsi

**Unit 6 (Lesson 26-28), pages 201-247** — çünki:
- Pointer-ları təhlükəsiz şəkildə (dangling yoxdur, aritmetika yoxdur) RAM konsepti ilə birləşdirir — yeni başlayan üçün ən dəqiq pointer dərsi.
- Nil semantikası (nil slice/map/interface fərqləri) — praktikada ən çox bug yaradan sahəni ən yaxşı izah edən mənbələrdən biridir.
- Error handling fəlsəfəsi (safeWriter, sentinel errors, custom types) — "errors are values" prinsipini əsaslandırır, sadəcə sintaksis yox.

Bu 3 lesson birlikdə "orta səviyyəli Go developer"-in gündəlik problem yaradan 80%-ni əhatə edir.
