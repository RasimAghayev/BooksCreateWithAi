# Effective Go (RU) — Xülasə (Azərbaycanca)

> **Sənəd:** Effective Go (rus tərcüməsi) — Go Team, go.dev rəsmi dokumentasiya (2009; 2022-də "yenilənməyəcək" qeydi — issue 28782)
> **Səviyyə:** 🎯 Intermediate (3/5) · **Dil:** Go · **45 səh., 9 emal bölməsi**

---

## Sənədin mahiyyəti

Effective Go — idiomatik Go yazımının **kanonik bələdçisi**: dil spesifikasiyasını tamamlayan, "Go-da bu necə ELEGANT yazılır" sualına cavab verən rəsmi mətn. Əsas tezis: C++/Java proqramını birbaşa tərcümə etmək YOX — problemi Go prizmasından yenidən düşünmək lazımdır. Sənəd 2009-dan bəri dəyişməyib: dil stabil qaldığından ÖZÜ aktuallığını saxlayır, amma ekosistem (modullar, test, generics) kənarında qalır.

## Bölmə-bölmə xülasə

### Formatlaşdırma və adlar
gofmt format müzakirələrini bitirir (tab, `{` eyni sətir — lekserin avtomatik `;` qaydası üzündən). Adlar semantikdirlər: böyük hərf = export; paket adı + export təkrarsızlığı (`ring.New`, `bufio.Reader`); getter `Owner` (Get YOX); 1-metodlu interfeys `-er` suffiksi; MixedCaps.

### İdarə konstruksiyaları
if/switch init ifadəsi qəbul edir; xəta-return-lərdə else buraxılır; `:=` eyni scope-da mövcud dəyişəni təkrar ELAN ETMİR (err zənciri). Qiymətsiz switch = if-else-if idiomu; fallthrough yoxdur; `break Loop` label patterni. String range rune-ları pars edir (yanlış → U+FFFD). Massiv çevirmə paralel mənimsətmə ilə.

### Funksiyalar
Çoxlu qayıdış C-nin -1/pointer-out idiomlarını öldürür; named nəticələr sənədləşmə + naked return; **defer** resurs açılışının yanında (unudulmaz close), arqumentləri dərhal, LIFO — trace patterni `defer un(trace("a"))`.

### Data: new/make/slice/map
`new(T)` = zero-value `*T`; `make` = slice/map/channel üçün initializə OLUNMUŞ T. "Zero value hazırdır" dizaynı (Buffer/Mutex — transitiv). Composite literal = konstruktor əvəzi (`&File{fd: fd}`). Massiv = value (kopya); slice = ptr+len+cap referans strukturu — append qayıtmalı. Map: comma-ok mövcudluq, delete təhlükəsiz. Çap: %v ailəsi; String metodu — %s rekursiya tələsi (əsas tipə çevir).

### İlkinləşdirmə və metodlar
Const = compile-time ifadələr; iota + implicit təkrar (KB..YB). var runtime; init determinist sıra (import→var→init) + state yoxlaması. Metod hər adlı tipə; pointer receiver caller-ı mutasiya edir (qayıtışsız Append); `*ByteSlice` Write ilə io.Writer — Fprintf hədəfi. Value metodu hər ikisində, pointer metodu yalnız pointer-də (lokalda avtomatik `&`).

### İnterfeyslər
İmplicit təmin; 1-2 metodlu kiçik interfeyslər normadır. Çevirmə (`sort.IntSlice(s).Sort()`) metod dəsti dəyişir. Type assertion/switch — comma-ok. "Yalnız interfeys export et": konstruktor interfeys qaytarsın (hash.Hash32, crypto/cipher Block→Stream). Handler interfeysini struct, int, channel, FUNKSİYA (HandlerFunc adapter) təmin edir — metod hər tipə yaxındır.

### Boş identifikator və embedding
`_`: istifadəsiz dəyərlər, dev-import, side-effect import (`_ "net/http/pprof"`), runtime yoxlama. **ƏSAS pattern: `var _ json.Marshaler = (*T)(nil)`** — kompayl-time interfeş şərtlənməsi. Embedding: interfeys birləşmə (ReadWriter) və struct promote (Job.Logger); irsdən fərq — receiver daxili tip; ad gizlətmə qaydaları.

### Konkurrentlik
Şüar: **"Yaddaşı kommunikasiya ilə paylaş"** (CSP). Goroutine — bitiş siqnalı yoxdur → kanallar; unbuffered = kommunikasiya+SİNXRON. Buffered channel semafordinu; hər-request-goroutine (resurs qeyri-limitli!) əvəzinə worker pool; Request.resultChan — mutexsiz RPC; GOMAXPROCS(0) ilə paralelləşdirmə; leaky buffer — select/default + GC ilə pool.

### Xətalar / panic / recover
error interfeysi zəngin strukturlarla (PathError: op+path+sistem xətası; prefiks "image: ..."). Caller type assertion ilə xətanı incə analiz edir (ENOSPC → retry). Panic yalnız davam-mümkünsüzdür; recover YALNIZ defer-də — safelyDo goroutine xilası; parserlərdə paket-daxili panic idiomu (defer named-result mutasiyası + re-panic). Yekun: flag + template.Must + HandlerFunc ilə bir neçə sətirlik QR veb-serveri.

---

## Sənədin əsas mesajları

1. **Birbaşa tərcümə YOX** — hər problemi Go idiomları ilə yenidən qur.
2. **Adlar və format dilin hissəsidir** — export semantikası, gofmt, kanonik imzalar.
3. **Zero value hazırlığı** — initializasiyasız işlək tiplər dizayn et.
4. **Kiçik interfeyslər + implicit təmin + çevirmə** — polimorfizmin Go yolu.
5. **Kommunikasiya idarəetməsi** — kanallar data race-i dizayn səviyyəsində öldürür.
6. **Defer funksiya-səviyyəli abstraksiyadır** — resurs + recover patternlərinin açarı.
7. **Panic paket sərhədini keçməz** — xaricə error, daxilə panic.

## Sənəddən sonra

- Go 1.18+ generics (sənədin "append niyə builtin" müzakirəsi tarixi kontekst verir)
- Go 1.22 loop var semantikası (Serve nümunəsindəki bug aradan qalxdı)
- Modern: modullar, test, race detector — sənəddə YOXDUR (2022 qeydi)
- A Tour of Go + How to Write Go Code — sənədin göstərdiyi ön-şərtlər
