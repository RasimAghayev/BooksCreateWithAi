# The Go Workshop — Xülasə (Azərbaycanca)

**Müəlliflər:** Delio D'Anna, Andrew Hayes, Sam Hennessy, Jeremy Leasor, Gobin Sougrakpam, Dániel Szabó | **Nəşriyyat:** Packt | **İl:** 2019 | **Səviyyə:** L2 (Elementary)

## Kitabın ümumi məqsədi

"The Go Workshop" — interaktiv workshop-formatlı, 19 fəsillik TAM GO dərsliyidir.
Əsas fərqi: hər fəsildə Exercise + Activity strukturudur — oxucu KOD YAZARAQ öyrənir.
Sıfırdan başlayır (dəyişənlər, tiplər), web proqramlaşdırmaya gətirir (HTTP
client/server, JSON, DB), concurrency və security ilə bitir. Kitab boyu ~40 praktiki
məşq + ~20 activity — real layihə təcrübəsi.

## Fəsil-fəsil xülasə

1. **Dəyişənlər/Operatorlar:** var/:= bəyanları, UTF-8 adlar, pointer-lər (value
   vs pointer, stack/heap, escape), iota, scope/shadowing.
2. **Məntiq/Döngələr:** initial if (err := f(); err != nil — GO İDİOMASI), switch,
   for-un 3 forması, range, break/continue, FizzBuzz, Bubble Sort.
3. **Əsas Tiqlər:** int/uint sinifləri, wraparound, math/big, string/rune/[]byte
   (len bayt sayı!), raw vs interpreted literal.
4. **Kompozit Tiqlər:** array ölçü=tip, slice hidden array (append cap qəlibi,
   bağlı/kopiyalı ssenarilər), map comma-ok/delete, interface{} + type switch.
5. **Funksiyalar:** çoxlu qaytarma, naked return (shadowing tələsi!), variadic
   (pack/unpack), closure (qorunmuş factory), funksiya tipləri, defer LIFO.
6. **Xətalar:** syntax/runtime/semantic, Err prefiksi, panic/recover/defer, tri-color
   istifadə halları.
7. **İnterfeyslər:** implicit satisfaction, duck typing, polimorfizm, io.Reader
   ("accept interfaces, return structs"), type assertion.
8. **Paketlər:** DRY, görünürlük, GOPATH/GOROOT, init sırası (import→var→init→main).
9. **Debuginq:** incremental test, unit test, fmt verb-lər (%T, %#v), log paketi
   (SetFlags), Fatal vs Panic.
10. **Vaxt:** Now/Sub/Add, Duration 6 rezolyusiya, Parse/Format (RFC3339/ANSIC),
    LoadLocation/In timezone.
11. **JSON:** Unmarshal/Marshal, tag-lər (omitempty, -), naməlum strukturu
    (map[string]interface{} + type switch), GOB binary protokolu.
12. **Fayllar:** rwx icazələri, flag paketi (CLI), signal tutma, os.Create/OpenFile
    (O_APPEND), os.Stat/IsNotExist, ReadFile/ReadAll, CSV reader.
13. **SQL:** database/sql + driver arxitekturası, sql.Open/Ping, Prepare
    (injection qorunması), Query/Scan, Exec/RowsAffected, TRUNCATE/DROP.
14. **HTTP Klient:** http.Get/Post, JSON decode/encode, multipart fayl upload,
    Authorization header, custom Client{Timeout}.
15. **HTTP Server:** ServeHTTP, Handle/HandleFunc, sayğac handler, querystring,
    html/template ({{if}}/{{.Field}}), FileServer/StripPrefix, ParseForm, JSON API.
16. **Concurrency:** goroutine, WaitGroup, race + `-race`, atomic, mutex, channel
    (buffered/unbuffered, close, range, worker pool, done channel), context.
17. **Alətlər:** go build/run, gofmt/goimports, go vet, race detector, go doc, go get.
18. **Təhlükəsizlik:** SQL/command injection, XSS, hash (SHA/BLAKE), AES-GCM
    simmetrik, RSA asimmetrik, crypto vs math/rand, TLS/x509, bcrypt parollar.
19. **Xüsusiyyətlər:** build tags/filename suffix, GOOS/GOARCH, reflect
    (TypeOf/ValueOf/DeepEqual), wildcard `./...`, unsafe/cgo.

## Ən vacib 5 fikir

1. **Workshop metodu:** hər konsept Exercise ilə əyani, Activity ilə müstəqil
   sınanır — passiv oxu YOX, aktiv yazma.
2. **Go error fəlsəfəsi:** error = dəyər; `if err != nil` — hər funksiyada explicit
   yoxlama; panic yalnız bütövlük üçün, recover defer-də.
3. **Slice/Channel daxili mexanika:** hidden array + cap (append bağlantı qirir);
   buffered/unbuffered channel bloklama semantikası — Go-nun ən bug-yaradan
   yerləri AÇIQLANIR.
4. **"Accept interfaces, return structs":** io.Reader qəbul edən funksiya string/
   fayl/HTTP üçün İŞLƏYİR — 3 fərqli funksiyanı 1-ə endirir.
5. **Alət mədəniyyəti:** gofmt+goimports+vet+`-race` — build prosesinin məcburi
   hissəsi; security: Prepare placeholder, html/template, bcrypt, TLS.

## Kitabın ən dəyərli hissəsi

Chapter 4 (slice daxili mexanizmi — linked/noLink/capLink ssenariləri) və Chapter
16 (Concurrency — atomic/mutex/channel/worker pool tam spektri) — bu iki fəsil
kitabın ən dərin texniki payıdır. Chapter 7 (interfaces) isə Go-nun fəlsəfi nüvəsini
(io.Reader transformasiya nümunəsi ilə) açır.
