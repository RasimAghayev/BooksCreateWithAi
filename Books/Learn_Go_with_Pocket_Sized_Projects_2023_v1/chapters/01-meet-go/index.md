# Chapter 1 — Meet Go (səh. 33-45)

## Bu chapter nədən bəhs edir?

Go-nun tarixi, dizayn fəlsəfəsi (sadəlik), əsas xüsusiyyətləri, digər dillərlə
müqayisəsi, built-in tooling (test, bench, fuzz) və "pocket-sized projects"
öyrənmə metodu.

## Əsas fikirlər

### 1. Go nədir?
Google-da böyük miqyaslı real problemləri həll etmək üçün yaradılmış proqramlaşdırma
dili. Dizaynı **sadəlik** ilə idarə olunur: cəmi 25 açar söz (2024-cü ilə qədər),
hər hansı anlaqı oxucu tez mənimsəyir.

### 2. Go harada güclüdür / zəifdir?
**Güclü:**
- Backend xidmətləri, API-lər, cloud computing — sadəlik + productivity
- Goroutine + channel ilə yüngül concurrency (paralellik)
- Cross-platform compiler — bir build çox platformaya
- Vast community + zəngin open source ekosistem

**Zəif / uyğun deyil:**
- OS yazmaq (garbage collector runtime-ın daxilindədir)
- Decades-lərə uzanan kitabxana tarixçəsi Java qədər deyil (amma böyüyür)

### 3. Generics (cinslər)
Go 1.18-dən əvvəl `slice`/`map` əməliyyatları hər tip üçün ayrı yazılırdı.
Generics sayəsində type-safe reusable funksiyalar — boilerplate azalır.

### 4. Dillərin müqayisəsi (Table 1.1)
| Ölçü | C++ | Python | Java | Go |
|---|---|---|---|---|
| Fəlsəfə | OOP, prosedual | Multi-paradigm | Yüksək səviyyə OOP | Yüksək səviyyə, çoxparadigmalı |
| Xətalar | İstisnalar (exceptions) | İstisnalar | İstisnalar | Error dəyərləri |
| Testing | Framework-lər | pytest və s. | JUnit | Native `test`, `bench`, `fuzz` |

### 5. Go tooling zənciri
- `go test` — unit testlər
- Benchmarking (`go test -bench`) — performansın commit-lərdə regresiyasını yoxla
- **Fuzzing** — sistemin random dəyərlərlə davranışını yoxlamaq (təhlükəsizlik
  üçün əvəzsiz)

### 6. Pocket-sized projects metodu
John Dewey-nin learning nəzəriyyəsinə əsaslanan yanaşma: **kiçik, tamamlanmış
layihələr** üzərindən öyrənmə — hər chapter müstəqil layihə (CLI aləti,
kitabxana, REST service, gRPC backend, maze solver...). Nəzəriyyə layihənin
ehtiyacı qədər verilir.

## Əsas terminlər

- Garbage Collector (zibil toplayıcı)
- Concurrency (paralellik) / Goroutine (qorutin)
- Generics (cins tiplər)
- Fuzzing (fuzz testi)
- Benchmark (performans ölçməsi)

## Praktik nəticə

- Go = sadəlik + productivity + native tooling üçlüyü
- Öyrənmə üsulu: kiçik tam layihələr — bu kitabın 12 layihəsi
- `go test` + benchmark + fuzz üçlüyünü commit rotasiyasına daxil et

## Mənbə

Pages: 33-45 (Chapter 1, Learn Go with Pocket-Sized Projects)
