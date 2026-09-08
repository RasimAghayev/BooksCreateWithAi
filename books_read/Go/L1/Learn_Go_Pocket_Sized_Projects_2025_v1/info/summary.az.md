# Learn Go with Pocket-Sized Projects — Xülasə (AZ)

## Kitab kimin üçündür?

Proqramlaşdırmaya yeni başlayanlar və ya başqa dildən Go-ya keçənlər üçün
"cepli layihə" metodu: 12 müstəqil, kiçik, əyləncəli layihə — hər biri bir-iki
yeni konsept qrupu öyrədir. TDD hər layihənin nüvəsindədir.

## 12 layihə xətti

| # | Layihə | Öyrənilən əsas konseptlər |
|---|---|---|
| 1 | Meet Go | Dilin tarixi, fəlsəfəsi, alət dəsti |
| 2 | Hello, earth! | Example test, custom tip, switch, map, TDT, flag |
| 3 | Bookworm's digest | JSON, defer, range, map-counter, sort.Interface |
| 4 | pocketlog kitabxanası | iota enum, functional options, io.Writer, doc.go |
| 5 | Gordle oyunu | runes, sentinel error, %w, Stringer, strings.Builder |
| 6 | Money converter | fixed-point Decimal, XML, HTTP, DI, httptest, timeout |
| 7 | Generic cache | generics, goroutine, channel, race, mutex, TTL |
| 8 | Gordle service | REST, ServeMux, repository, status kodları |
| 9 | Maze solver | PNG, linked list, select, quit kanal, WaitGroup, GIF |
| 10 | Habits tracker | protobuf, gRPC, minimock, context, integration test |
| 11 | HTML UI | go:embed, template, form, gRPC client |
| 12 | Wasm/TinyGo | js.FuncOf, DOM, mikrokontroler, build tag-lər |

## Ən vacib 5 fikir

1. **"Errors are values"** — xətalar adi dəyərlər kimi qaytarılır, wrap olunur
   (`%w`), sentinel + `errors.Is` / typed + `errors.As`
2. **Test-first mədəniyyəti** — Example test (stdout), TDT (cədvəl),
   `t.Parallel()`, `go test -race`, golden files, `-short` integration skip
3. **Asılılıqlar interfeyslə injekt olunur** — io.Writer, io.Reader, kiçik lokal
   interfeyslər ("discovered, not designed"); httptest/minimock ilə mock
4. **Sadəliyin praqmatikası** — "make it work, make it clean, maybe fast";
   kanal ehtiyac yoxdansa mutex; domain ≠ protokol
5. **Go-nun ekosistemi genişdir** — Wasm (brauzer), TinyGo (mikrokontroler),
   generics, embed — hamısı standart/official alətlərlə

## Kitabın ən dəyərli hissəsi

Chapter 6 (money converter) float təhlükəsini int-arifmetika ilə, Chapter 7
(generic cache) race-in `go test -race` ilə AŞKAR EDİLMƏSİNİ göstərir — iki
dərs də praktikada ən çox burulğan yaradan yerlərdir.

## Metodologiya

- Hər fəsil: məqsəd → TDD ilə addım-addım → refactor → "side quests" (əlavə
  tapşırıqlar) → summary
- Manning livebook formatı; kod listinqləri tam izahlı
- Kod: github.com/ThreeDotsTech/... deyil — kitabın öz repo strukturu
  (learngo-pockets/*)
