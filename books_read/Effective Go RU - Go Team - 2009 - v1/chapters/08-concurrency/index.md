# Chapter 8 — Konkurrentlik (səh. 36-41)

## Bu chapter nədən bəhs edir?

Go-nun konkurrentlik fəlsəfəsi ("Do not communicate by sharing memory; share memory by communicating"), goroutine-lər, kanalın sinxronizasiya təbiəti, semafor kimi buffered channel, kanal-kanal pattern (paralel RPC), multi-core paralelləşdirmə (Vector nümunəsi), concurrency vs parallelism fərqi və leaky buffer pool patterni.

---

## Əsas fikirlər

### 1. Əsas şüar

> **"Не обменивайтесь данными через общую память; вместо этого обменивайтесь памятью через коммуникации."**
> (Ümumi yaddaşla məlumat paylaşma; yaddaşı KOMMUNİKASİYA ilə paylaş.)

Dəyər kanaldan keçir və hər anda YALNIZ bir goroutine sahiblikdədir → **data race dizayn səviyyəsində mümkünsüzləşir.** Kökü: CSP (Communicating Sequential Processes) modeli — Unix kanallarının tip-təhlükəsiz ümumiləşməsi.

Ekstrem hallarda mutex daha yaxşıdır (məs., integer reference count) — amma yüksək səviyyə yanaşma olaraq kanallar aydın və düzgün proqramlar yazdırır.

### 2. Goroutine-lər

- Adlar mövcud terminlər (thread/coroutine/process) səhv assosiasiya yaratdığı üçün yeni
- Eyni adres fəzasında paralel icra olunan FUNKSİYA
- Yüngül: yalnız stack təyinatı; kiçik başlayır, lazımcas böyüyüb kiçilir (heaptan)
- OS thread-lərə MULTİPLEKS olunur — biri bloklananda (I/O) digərləri işləyir
- Bitəndə siqnal YOXDUR (Unix `&` fon rejimi kimi)

```go
go list.Sort()    // paralel; gözləmə YOX

// Closure ilə:
func Announce(message string, delay time.Duration) {
    go func() {
        time.Sleep(delay)
        fmt.Println(message)
    }()    // MÖTƏRİZƏ ŞƏRT — çağırış!
}
```

Funksional literallar Go-da CLOSURE-dir — istinad olunan dəyişənlər aktiv olduqca yaşayır. Amma bu nümunələr praktik deyil — bitmə siqnalı YOXDUR → **kanallar.**

### 3. Kanallar — kommunikasiya + sinxronizasiya

```go
ci := make(chan int)           // unbuffered = SİNXRON
cj := make(chan int, 0)         // eyni
cs := make(chan *os.File, 100)  // buffered
```

**Unbuffered = kommunikasiya + sinxronizasiya birlikdə** — 2 goroutine məlum vəziyyətə gəlir.

**Bitmə gözləməsi idiomu:**

```go
c := make(chan int)
go func() {
    list.Sort()
    c <- 1        // siqnal; DƏYƏR ƏHƏMİYYƏTSİZ
}()
doSomethingForAWhile()
<-c              // sort bitənə qədər gözlə (dəyər atılır)
```

Qəbul edən dəyər gələnə qədər bloklanır; unbuffered göndərən qəbul edilənə qədər; buffered — bufferə kopyalanana qədər.

### 4. Buffered channel = semafor

```go
var sem = make(chan int, MaxOutstanding)

func handle(r *Request) {
    sem <- 1        // "boş yer gözlə" — doludursa BLOK
    process(r)       // uzun sürə bilər
    <-sem            // yer azad et
}
```

**Xətalı Serve** — hər request üçün YENİ goroutine (resurssuz böyümə!):

```go
func Serve(queue chan *Request) {
    for req := range queue {
        sem <- 1
        go func() {
            process(req)      // Go <1.22: loop var hamıda ORTAQ — BUG!
            <-sem
        }()
    }
}
```

**Düzgün — fiksiləşmiş handler sayı (worker pool):**

```go
func handle(queue chan *Request) {
    for r := range queue {
        process(r)
    }
}
func Serve(clientRequests chan *Request, quit chan bool) {
    for i := 0; i < MaxOutstanding; i++ {
        go handle(clientRequests)    // N goroutine = N paralel process
    }
    <-quit    // söndürmə siqnalı
}
```

### 5. Kanal-kanal — paralel, bloklanmayan RPC

Kanal — birinci dərəcəli DƏYƏR: ötürülə, saxlanıla bilər:

```go
type Request struct {
    args       []int
    f          func([]int) int
    resultChan chan int    // hər klient ÖZ cavab kanalı!
}

request := &Request{[]int{3, 4, 5}, sum, make(chan int)}
clientRequests <- request                          // göndər
fmt.Printf("answer: %d\n", <-request.resultChan)   // cavab beklə — mutexsiz RPC!

func handle(queue chan *Request) {
    for req := range queue {
        req.resultChan <- req.f(req.args)
    }
}
```

### 6. Paralelləşdirmə — multi-core

```go
type Vector []float64

func (v Vector) DoSome(i, n int, u Vector, c chan int) {
    for ; i < n; i++ {
        v[i] += u.Op(v[i])
    }
    c <- 1    // hissə bitdi siqnalı
}

var numCPU = runtime.GOMAXPROCS(0)   // NumCPU default; istifadəçi təyin edə bilər

func (v Vector) DoAll(u Vector) {
    c := make(chan int, numCPU)
    for i := 0; i < numCPU; i++ {
        go v.DoSome(i*len(v)/numCPU, (i+1)*len(v)/numCPU, u, c)
    }
    for i := 0; i < numCPU; i++ {
        <-c    // hamısı bitənə qədər say
    }
}
```

**Concurrency ≠ parallelism:** Konkurrentlik = struktur (müstəqil komponentlər); paralellik = hesabın çoxprosessorlu icrası. Go — konkurrent dil; bütün paralelləşdirmə taskları modele uymur.

### 7. Leaky buffer — free-list pool

RPC paketindən real pattern — buffered channel + GC ilə buffer reuse (select+default):

```go
var freeList = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)

func client() {
    for {
        var b *Buffer
        select {
        case b = <-freeList:      // boş buffer varsa AL
        default:                  // YOXDURSA
            b = new(Buffer)       // yarat
        }
        load(b)             // şəbəkədən oxu
        serverChan <- b     // serverə göndər
    }
}

func server() {
    for {
        b := <-serverChan
        process(b)
        select {
        case freeList <- b:      // geri qaytar (yer varsa)
        default:                  // doludursa — GC götürəcək
        }
    }
}
```

`select`-in `default`-u: heç bir case hazır deyilsə — bloklamaMADAN. Bir neçə sətirlə pool idarəsi — yaddaş GC-yə həvalə.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| CSP | Communicating Sequential Processes — Go konkurrentliyin mənbəyi |
| Goroutine | Eyni adres fəzasında yüngül paralel funksiya; stack heapdan böyüyür |
| Unbuffered channel | Göndərən = qəbul edənə qədər blok — SİNXRON |
| Buffered semafor | `sem <- 1` / `<-sem` — throughput limiti |
| Worker pool | Fiksiləşmiş N handler goroutine — resurs limiti |
| Loop var bug | Serve nümunəsi — Go <1.22-də closure req hamıya ortaq |
| resultChan | Request DAXİLİNDƏ cavab kanalı — mutexsiz RPC |
| GOMAXPROCS(0) | İstifadəçi üstünlüyünə hörmət edən core sayı |
| Leaky buffer | `freeList` channel + default — GC idarəli pool |
| select default | Bloklamayan yoxlama — hazır case yoxdursa |

---

## Praktik nəticə

1. **Şüarı daxililəşdir:** Dəyəri kanaldan keçir, sahibliyi tək goroutine-də saxla — race-lər dizaynla ölür. Mutex yalnız ağır kontendlərdə (ref count).
2. **Goroutine bitiş siqnalı ver:** `c <- 1` + `<-c` idiomu; closure-lar dəyişənləri yaşadır.
3. **Semafor = buffered channel:** MaxOutstanding limiti `sem <- 1` ilə; amma hər request-ə goroutine YOX — worker pool (fiksiləşmiş N) qənaətcildir.
4. **Serve loop-var tələsinə diqqət:** Go <1.22-də closure req-i kopyalamalı (`req := req`).
5. **Cavabı request-ə daxilə qoy:** `resultChan` — asinxron RPC, sıraya gözləmədən.
6. **numCPU = GOMAXPROCS(0)** — NumCPU istifadəçi üstünlüyünü (env) hörmət etmir.
7. **Pool patterni:** `select { case freeList <- b: default: }` — dolu haldan GC-yə keçiş; bir neçə sətirlik idarə.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi), 2009
- PDF səhifələri: 36-41
- İstinad: Go wiki — loop variable pitfall; CSP (Hoare)
