# Build systems with Go: Everything a Gopher must know — Xülasə (AZ)

**Juan M. Tirado, 2021, 472 səh., L2 (Elementar)**

## 📖 Kitab deyir

Kitab iki hissədən ibarətdir. **Hissə I** Go dilinin özünü əhatə edir:
ilk proqram və `go build` (ch1); dilin əsasları — dəyişənlər, funksiyalar,
pointerlər, defer/panic/recover, init (ch2); array/slice/map (ch3); struct,
metod, interfeys (ch4); reflection (ch5); konkurensiya — goroutine, kanal,
select, WaitGroup, Timer/Ticker, Context, sync (ch6); I/O — Reader/Writer,
fayllar, bufio (ch7); encoding-lər — CSV, JSON, XML, YAML (ch8); HTTP
klient/server, cookie, middleware (ch9); şablonlar (ch10); test, benchmark,
coverage, profiling (ch11); modullar və sənəd (ch12).

**Hissə II** sistem qurulması: Protocol Buffers (ch13); gRPC — server/client,
streaming, transcoding, interceptorlar (ch14); Zerolog ilə logging (ch15);
Cobra ilə CLI (ch16); `database/sql` + GORM (ch17); Cassandra/GoCQL (ch18);
Kafka — Confluent, Segmentio, REST Proxy klientləri (ch19).

Kitabın fəlsəfəsi: hər mövzu kiçik, müstəqil nümunələrlə izah olunur
("Save the world with Go!!!" motivi), böyük layihə yoxdur — dilin və
ekosistemin alətlərinə praktiki giriş.

## 👨‍🏫 Müəllim qeydi

Bu kitab "Everything a Gopher must know" adına uyğun gəlir — Go-nun demək
olaraz bütün istehsal səthinə toxunur. Praktikada belə tətbiq etmək daha
yaxşı olar:

1. **Concurrency (ch6) əsasdır** — kitab Context-in API sorğularında
   istifadəsini göstərir, amma real layihədə hər xarici çağırışın öz
   timeout-lu context-i olmalıdır; errgroup və worker-pool pattern-lərini
   kitab əhatə etmir — onları əlavə öyrənin.
2. **ioutil (ch7) köhnəlib** — müasir Go-da `os.ReadFile`/`os.WriteFile`
   və `io.ReadAll` istifadə edin; kitabın nümunələri Go 1.15 dövrünə aiddir.
3. **gRPC auth nümunəsi (ch14)** tədris məqsədlidir — production-da TLS +
   token/interceptor zənciri və ya istifadəçi doğrulama mərkəzi lazımdır.
4. **GORM-də raw SQL bilmək yenə də vacibdir** — ORM-in generasiya etdiyi
   sorğuları `Debug()` ilə izləyin; N+1 probleminə Preload diqqət yetirin.
5. Kitab generics-i əhatə etmir (yazıldığı vaxt Go-da yox idi) — müasir
   layihələrdə type-parameter-ləri (Go 1.18+) reflection (ch5-dəki
   MakeFunc işarları) əvəz edə bilər.

## Ən vacib 5 fikir

1. **Error dəyəri Go-nun sütunudur** — `(dəyər, error)` konvensiyası +
   `if err != nil`; panic yalnız bərpaolunmaz hallarda (ch2, 7).
2. **Concurrency = kanallar + select + Context** — vaxt idarəetməsi və
   ləğvetmə üçün; `defer cancel` şərtdir (ch6).
3. **Hər I/O interfeysə dayanır** — Reader/Writer anlayışı fayl, şəbəkə,
   bufer, HTTP body hamısında eynidir (ch7).
4. **IDL + kod generasiyası** (PB/gRPC) — sorğu/cavab strukturlarını
   tək mənbədən bir neçə dil üçün törədir (ch13-14).
5. **Test/benchmark/profiling inteqrasiya olunub** — `go test -bench`,
   `-cover`, `go tool pprof` xətti xarici alət tələb etmir (ch11).

## Kitabın ən dəyərli hissəsi

Chapter 6 (Concurrency) — səh. 102-142: goroutine-lərdən Context-ə qədər
tam konkurensiya sistemi; kitabın digər fəsillərindəki nümunələrin hamısı
bu biliyə dayanır. Həmçinin Chapter 14 (gRPC) — server, streaming,
transcoding, interceptor kombinasiyası real mikroservis əlaqəsinə ən yaxın
əhatədir.
