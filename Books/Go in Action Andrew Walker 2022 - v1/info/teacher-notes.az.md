# Müəllim Qeydi — Go in Action, Second Edition

📖 Kitab deyir:
- Go ilk növbədə praktik dildir. Kompilyasiya sürəti, concurrency və error handling əsas üstünlükləridir.
- `go fmt`, `go mod`, `go doc` alətləri developer workflow (iş axını)nda mərkəzdədir.
- Slice-lər arrays-dən daha çox işlənir, nil/empty slice fərqi vacibdir.
- Interface-lər implicit implementasiya tələb edir — inheritance (miras) yoxdur.
- Map-lər orderless (sırasız)dir, nil map panic verir.

👨‍🏫 Müəllim qeydi:
Məncə bu kitabın ən güclü tərəfi onun praktik yanaşmasıdır. İlk chapter-dan başlayaraq real tətbiq (word counter) qurulur və bu developer-ə dərhal tətbiq edə biləcəyi bilik verir. Lakin bir neçə məsləhət:

## Ən vacib 5 fikir

1. **Concurrency (paralellik) əvvəlcə praktik edin** — Kitab goroutine və channel-ləri izah edir, lakin Chapter 8-də (hələ yoxdur) daha dərindən əhatə olunur. Chapter 1-dəki nümunələri kompayl edib işlətmək tövsiyə olunur.

2. **Error handling pattern-ini ənənnən edin** — `if err != nil` təkrarlanması ilkin olaraq verbose (çox sözlü) görünə bilər, lakin bu Go-nun explicitly (açıq-aydın) yanaşmasıdır. Məncə bu pattern tezliklə alışıla bilər və debug-ı asanlaşdırır.

3. **Slice append davranışını başa düşün** — `append` həmişə yeni slice qaytarır, lakin capacity dolduqda backing array dəyişir. Paylaşılan data üzərində dəyişiklik edərkən diqqətli olun.

4. **Generics yeni versiyadır** — Kitab Type Parameter-ləri izah edir, lakin real dünya tətbiqlərində generics istifadəsi hələ inkişaf etməkdədir. Məncə başlanğıcda interface-lərdən istifadə etmək daha təhlükəsizdir.

5. **Go Modules defaultdir** — `GOPATH` artıq keçmişdə qaldı. Hər layihə üçün `go mod init` ilə başlayın.

## Kitabın ən dəyərli hissəsi

Chapter 2, pages 25-66 — çünki burada ilk real Go tətbiqi adım-addım qurulur. `go build`, `go run`, `gofmt`, error handling, `bufio.Scanner` və `io.Reader` konsepsiyaları praktik nümunələrlə öyrədilir. Bu chapter oxumadan sonra developer real Go layihəsi qurmağa başlaya bilir.

## Təlimatlar

- **Yeni başlayanlar:** Chapter 1 və 2-ni ardıcıl oxuyun. Word counter tətbiqini özünüz yazın.
- **Orta səviyyəli:** Chapter 3 və 4-ü tiplər üzərində praktik edin. Slice və map emalı üçün mini layihə qurun.
- **Tədris məqsədləri üçün:** Chapter 1-dəki concurrency və type system bölmələri diskussiya üçün ən yaxşı materialdır.
