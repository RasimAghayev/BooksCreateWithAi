# Go Programming - From Beginner to Professional — Xülasə (AZ)

**Samantha Coyle, 2024 (Packt, 2nd Edition), 680 səh., L1 (Beginner)**

## 📖 Kitab deyir

Kitab 6 hissəyə bölünmüş 21 fəsildə Go-nu sıfırdan production-a qədər aparır:

- **Part 1 Scripts (ch1-4):** dəyişənlər/operatorlar (shadowing, pointer,
  iota), axın idarəsi (if/switch/for/goto), core tiplər (wraparound, rune,
  nil), kompleks tiplər (slice daxili mexanika, map, struct embedding,
  any/type switch).
- **Part 2 Components (ch5-8):** funksiyalar (variadic, closure, defer
  FILO), xətalar (panic/recover, Err konvensiyası, %w wrapping),
  interfeyslər ("accept interfaces, return structs"), generics (constraint,
  comparable, inference).
- **Part 3 Modules (ch9-10):** go.mod/go.sum, go get, workspaces (go.work),
  paketlər (exported/unexported, init sırası, cmd+pkg strukturu).
- **Part 4 Applications (ch11-15):** debugging (%T, %#v, log.SetFlags,
  Fatal/Panic), time (Parse/Format, zone-lar), CLI (flag, signal-lar,
  Rot13 pipeline, Bubbletea TUI), fayllar (icazələr, OpenFile append, CSV,
  embed), PostgreSQL (Prepare + injection qoruması, GORM).
- **Part 5 Web (ch16-17):** HTTP server (Handler, middleware, html/template,
  FileServer), HTTP client (Get/Post, multipart upload, custom
  Client{Timeout}).
- **Part 6 Professional (ch18-21):** konkurensiya (goroutine, WaitGroup,
  atomic, Mutex, kanallar, context, sync.Cond, sync.Map), testing
  (table-driven, sqlmock, httptest, fuzz, benchmark, coverage ~80%), Go
  alətləri (gofmt, goimports, vet, --race, doc), cloud (Prometheus
  metrikləri, OTel tracing, multi-stage Docker, K8s konseptləri).

Kitabın fəlsəfəsi: hər mövzu kiçik Exercise/Activity-lərlə addım-addım
qurulur; kiçik nüanslar (wraparound, loop-variable capture, 0644 oktal)
real bug hekayələri ilə motivasiya olunur.

## 👨‍🏫 Müəllim qeydi

Packt-in "Beginner to Professional" seriyasına xas olaraq kitab həm
abzentinə, həm də sürətə ümid bağlayır. Belə istifadə etmək daha yaxşıdır:

1. **Exercise-ləri TƏK-TƏK yazın** — kitabın dəyəri oxuda deyil, barmaqda;
   bir çox incəlik (shadowing bug-ı, slice link semantikası) yalnız
   praktikada görünür.
2. **Sıra vacibdir:** ch1-7 ardıcıl oxunmalı; ch9-dan sonra istənilən
   tərtib mümkündür (modullar standalone-dur).
3. **Kitabın bug-nümunələri qızıl dəyərdir** — Exercise 18.03-ün -race
   demosu, Activity 1.04-ün shadowing tələsi: bunlar müsahibə suallarıdır.
4. **Versioning:** Go 1.21 hədəflənir — kitabda var (yeni), amma dərslərdə
   öz Go versiyanızı yoxlayın.
5. **2024-cü il üçün bir az sadələşdirmələr:** http.HandleFunc yeni
   layihələrdə 1.22+ router pattern-ləri ilə zənginləşir; amma kitabın
   yanaşması əsaslandırılmamış deyil — konsept baxımından təmizdir.

## Ən vacib 5 fikir

1. **Shadowing** — uşaq scope-də `:=` yeni dəyişən yaradır; köhnəni
   dəyişmək istəyirsən `=` (Activity 1.04).
2. **Slice = view** — təyinat kopya YOXDUR; müstəqil kopya üçün
   `append([]T{}, s...)` (ch4).
3. **Xəta dəyər kimi** — `(T, error)` + `Err...` konstantaları + panic
   yalnız bərpaolunmazda (ch6).
4. **Kanal + close + range** — worker pool-un üçlüyü; context.Done()
   sonsuz loop-u dayandırır (ch18).
5. **Test = hər səviyyədə** — unit table-driven, integration sqlmock,
   HTTP httptest, performans benchmark, coverage ~80% hədəf (ch19).

## Kitabın ən dəyərli hissəsi

Chapter 18 (Concurrent Work) və Chapter 19 (Testing) — birincisi Go-nun
konkurensiya modelini (WaitGroup → atomic → Mutex → kanal → context →
sync.Cond/sync.Map) təbii inkişaf xətti ilə izah edir; ikincisi test
piramidasını tam əhatə edir. Birlikdə "professional" addının əsl mənasını
daşıyır.

## Son söz

Kitabın öz üslubu ilə: "bu kitab sizi Go biliknizi professional Go
developer səviyyəsinə transformasiya etmək üçün alətlər və biliyi verir" —
düzdür, amma yalnız Exercise-lər yazılırsa.
