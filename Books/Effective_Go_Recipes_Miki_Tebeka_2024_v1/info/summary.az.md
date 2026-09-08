# Effective Go Recipes — Xülasə (AZ)

**Miki Tebeka, 2024 (Pragmatic Bookshelf), 268 səh., 84 resept, L3 (Intermediate)**

## 📖 Kitab deyir

Kitab 84 qısa reseptdən ibarətdir — hər biri Task → Solution → Discussion
quruluşundadır. Mövzular:

- **I/O (ch1):** bytes.NewReader, gzip janitor, bytes.Buffer ilə SQL,
  şərtli dekompressiya, öz io.Writer (Benford), os.Pipe, mmap axtarışı.
- **Serializasiya (ch2):** gob eventlər, anonim struct JSON, streaming
  JSON, pointer-lərlə itkin sahə, MarshalJSON/UnmarshalJSON, mapstructure
  dinamik tiplər, reflect ilə struct tag parse.
- **HTTP (ch3):** pagination-li GET, chunked POST, middleware, server
  timeout-ları, paralel API versiyaları.
- **Mətn (ch4):** %#v debugging, Stringer, encoding aşkarı, regexp
  camelCase→snake, EqualFold, Unicode normalizasiyası (NFKC).
- **Funksiyalar (ch5):** registry map, functional options (WithX),
  closure arqumentlər, watcher-lər, go:linkname (dirty hack).
- **Əsas tiplər (ch6):** comma-ok, slice-Stack (memory leak qoruması),
  make-cap benchmark optimallaşdırması (19→1 alloc), JSON vaxt formatı,
  composite açarlar, təbii dil vaxt parse (timezone).
- **Struct/Interface (ch7):** ad hoc syncer, errWriter (ResponseWriter
  sargısı), generics — Max, Ring buffer, pointer-təhlükəsiz UnmarshalJSON.
- **Xətalar (ch8):** %w wrap, panic→error (named returns), safelyGo,
  errors.Is, runtime.Callers ilə stack Wrap.
- **Konkurensiya (ch9-10):** fan-out/fan-in, buffered-channel semafor,
  worker pool, Context timeout (select+Done), ctx logger; sync.Once,
  WaitGroup, RWMutex, race detector, atomic.Value ilə 21x Now().
- **Socket/cgo (ch11-12):** TCP fayl transferi, Unix-socket JSON RPC, NTP
  UDP binary; os/exec (ping, bc pipe), cgo (ioctl, snowball stemmer).
- **Test (ch13):** CI-yalnız testlər, YAML table tests, fuzzing,
  RoundTripper mock, TestMain, end-to-end server test, go/analysis linter.
- **Build/Ship (ch14-15):** embed, ldflags versiya, CGO_ENABLED=0 statik,
  build tags, goreleaser, go:generate; ardanlabs/conf konfiqurasiya,
  replace patch, Docker multistage, graceful shutdown, zap Check, expvar
  metrics, Delve attach.

## 👨‍🏫 Müəllim qeydi

Bu, "problem → həll" formatlı praktik sorğu kitabıdır — sistematik dərslik
deyil. Belə istifadə etmək daha yaxşı olar:

1. **Sıra əhəmiyyətlidir:** ch1-2 (I/O + serializasiya) hər şeyin təməlidir;
   yeni başlayan ch9-10-a (konkurensiya) dərhal atlamasın — kanal
   semantikası ch1-2-dəki Reader/Writer zəncirləri ilə möhkəmlənir.
2. **Reseptlər təcrididir** — amma cross-referanslar var (mmap → build
   tags; errWriter → metrics middleware). İkinci oxunuşda bu bağlantıları
   qeyd edin.
3. **Generics (ch7) müasir Go üçün əvəzsizdir** — kitab Go 1.20 dövründə
   yazılıb; golang.org/x/exp/slices/maps artıq standart kitabxanada
   `slices`/`maps` paketləridir.
4. **zap.Check (ch15)** hər logging sahəsündə yadda saxlanmalı idiomdur —
   bahalı funksiya çağırışları parametr qiymətlənməsində gizlidir.
5. Kitab SQLite testlərinə işarə edir (640 test sətri/source sətir) —
   "pain vs gain" balansı hər layihədə fərdi qərardır.

## Ən vacib 5 fikir

1. **io.Reader/io.Writer hər yerdədir** — funksiyalara konkret tip yox,
   interfeys qəbul etdirin (Recipe 4, 22).
2. **Comma-ok üç sahədə həyat qurtarır** — map, kanal, type assertion
   (Recipe 31).
3. **make(0, n) ilə append growth-u gələcəkdə öldürün** — benchmark
   3.7x sürət + 19x az alloc (Recipe 33).
4. **Context ilk parametr** — timeout/cancel üçün select+Done; loop-daxili
   cancel-i defer etməyin (Recipe 50).
5. **API-ni pozmadan genişləndirin** — functional options, ad hoc
   interfeyslər, build tags, API versiyalama (Recipe 27, 37, 75, 19).

## Kitabın ən dəyərli hissəsi

Chapter 6 (Basic Types) və Chapter 13 (Testing) — birincisi zero-value
tələləri, slice daxili mexanikası və vaxt mürəkkəbliyini real bug-larla
izah edir; ikincisi YAML table tests + fuzzing + mock + e2e + linter
tam test arsenalını bir fəsildə verir.
