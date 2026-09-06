# Chapter 5 — Конструирование облачной службы

## Bu chapter nədən bəhs edir?

Chapter 4-də öyrənilən şablonları real bir bulud xidmətinin tətbiqində göstərir: dağıtılmış açar/dəyər (key/value) mağazasından başlayaraq RESTful API, əməliyyat qeydiyyatı (transaction logging), dayanıqlı yazma (durability) və Docker konteynerizasiyası qədər. Bu chapter çox praktikdir və təxminən 56 səhifəyə yayılır.

## Əsas fikirlər

### 1. Minimal miqyaslana bilən açar/dəyər (key/value) mağazası
**Nədir:** HTTP üzərindən çalışan, əsas GET, PUT, DELETE əməliyyatlarını dəstəkləyən minimalist bulud xidməti.

**Necə işləyir:**
- `net/http` standart kitabxanası ilə RESTful endpoint-lər qurulur
- `gorilla/mux` marşrutlaşdırıcısı (router) dəyişən yol seqmentləri (`{key}`, `{category}/{id:[0-9]+}`) və metod şərtləri (methods, schemes) ilə
- Sorğular `http.Request` vasitəsilə marşrutlaşdırılır, cavablar `http.ResponseWriter` ilə qaytarılır
- Keçid (transition) qeydiyyatçısı (transaction logger) əvəzinə birinci versiyada sadəcə yaddaşda saxlanır

### 2. Əməliyyat qeydiyyatı (Transaction Logging)
**Nədir:** Dayanıqlılığı (durability) təmin etmək üçün hər dəyişikliyi diskə əlavə edən (append-only) qeydiyyat mexanizmi.

**Necə işləyir:**
- Hər qeyd ardıcıl (sequential) nömrə, hadisə növü (PUT/DELETE), açar və dəyər mətnləri ilə yazılır
- Fayl əsaslı realizasiya: `os.File`, `bufio.Scanner` ilə oxuma, `fmt.Fprintf` ilə yazma
- PostgreSQL əsaslı realizasiya: `database/sql` və `github.com/lib/pq` driver-i ilə parametrizə olunmuş SQL (`$1, $2, $3`) əməliyyatları
- `TransactionLogger` interfeysi yaradılır — iki metod: `WritePut(key, value)` və `WriteDelete(key)`

**Kitabdan kod nümunəsi (fayl əsaslı qeydiyyatçı):**
```go
type FileTransactionLogger struct {
    events chan Event
    errors chan error
    file   *os.File
}

func (l *FileTransactionLogger) Run() {
    for e := range l.events {
        _, err := fmt.Fprintf(l.file, "%d\t%s\t%s\t%s\n", e.Sequence, e.EventType, e.Key, e.Value)
        if err != nil {
            l.errors <- err
        }
    }
}
```

### 3. RESTful marşrutlaşdırma və xidmət strukturu
**Nədir:** `gorilla/mux` ilə dəyişən yollar, metod filtrasiyası və host-based marşrutlaşdırma.

**Necə işləyir:**
- `{key}` — dəyişən yol seqmenti
- `{category}/{id:[0-9]+}` — regex ilə məhdudlaşdırılmış seqment
- `r.Methods("GET").HandlerFunc(getHandler)` — yalnız GET sorğuları
- `r.Host("example.com")` — host əsaslı marşrut

### 4. Dayanıqlı yazma və diskdən bərpa
**Nədir:** Xidmət yenidən başladıqda vəziyyəti diskdən qeydiyyat faylından oxuyaraq bərpa etmə.

**Necə işləyir:**
- Qeydiyyat faylı təkrar oxunur və hər sətir təkrar tətbiq edilir (replay)
- Sıra nömrəsi ilə təkrarlar qarışdırılır və ardıcıllıq təmin edilir
- `key/value` mağazası qeydiyyat faylını oxuyaraq əvvəlki vəziyyətə qayıdır

### 5. Docker konteynerizasiyası
**Nədir:** Go-nun statik binar (static binary) kompilyasiya xüsusiyyəti sayəsinde Docker konteynerlərdə ideal şəkildə işləməsi.

**Necə işləyir:**
- Çox mərhəmətli (multi-stage) build: birinci mərhələdə binar build edilir, ikinci mərhələdə yalnız binary kopyalanır
- Nəticə: ~2 MB-lik minimal image, heç bir xarici runtime tələb etmir
- Scratch və ya Alpine base image-lərdən asılılıq yoxdur

## Əsas terminlər
- Transaction Log (əməliyyat qeydiyyatı) — append-only dayanıqlı qeyd
- RESTful API — uniform interfeys üzərindən HTTP sorğuları
- gorilla/mux — Go marşrutlaşdırıcısı (router)
- key/value store — açar/dəyər mağazası
- Durability (dayanıqlılıq) — diskdə qeydiyyat vasitəsilə məlumat itkisinin qarşısı alınması
- Multi-stage Docker build — çox mərhəmətli Docker qurulumu
- database/sql — Go standart verilənlər bazası interfeysi
- Replay — qeydiyyat faylını yenidən icra edərək vəziyyəti bərpa etmə

## Praktik nəticə
Bu chapter 4-cü chapter-dakı şablonları real tətbiqdə göstərir: Go xidməti `net/http` və `gorilla/mux` ilə qurulur, əməliyyat qeydiyyatı interfeysi ilə dayanıqlılıq təmin edilir və statik binar özəlliyi sayəsində Docker konteynerə asanlıqla yerləşdirilir. Bu, "kitabxanalardan kod yazmaq" anlayışına əsaslanır.

## Mənbə
Pages: 125-180 (PDF səh. 125-180)
