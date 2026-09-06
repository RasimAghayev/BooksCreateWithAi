# Chapter 7 — Concurrency patterns

## Bu chapter nədən bəhs edir?

Üç real concurrency pattern-in tam implementasiyası: **Runner** (proqramın ömür müddətini idarəetmə — timeout + OS interrupt), **Pooling** (buffered channel ilə resurs hovuzu — DB connection pool), **Work** (unbuffered channel ilə goroutine hovuzu — işçi pool). Hər biri channel + select + WaitGroup kombinasiyasının praktik tətbiqidir.

## Əsas fikirlər

### Pattern 1: Runner — proqram ömrünün idarəsi
**Məqsəd:** Cron job və ya bulud worker kimi nəzarətsiz işləyən proqramlar üçün: tapşırıqları vaxtında bitsə normal, bitməsə timeout, Ctrl+C olsa təmiz shutdown.

**3 sonlanma ssenarisi:**
1. İş vaxtında bitir → normal qayıdış (nil)
2. Vaxt bitir → `ErrTimeout`
3. OS interrupt (Ctrl+C) → `ErrInterrupt` + təmiz söndürmə

**Kitabdan kod nümunəsi (runner.go — tam paket):**
```go
// Package runner manages the running and lifetime of a process.
package runner

import (
    "errors"
    "os"
    "os/signal"
    "time"
)

// Runner runs a set of tasks within a given timeout and can be
// shut down on an operating system interrupt.
type Runner struct {
    // interrupt channel reports a signal from the operating system.
    interrupt chan os.Signal

    // complete channel reports that processing is done.
    complete chan error

    // timeout reports that time has run out.
    timeout <-chan time.Time

    // tasks holds a set of functions that are executed
    // synchronously in index order.
    tasks []func(int)
}

// ErrTimeout is returned when a value is received on the timeout.
var ErrTimeout = errors.New("received timeout")

// ErrInterrupt is returned when an event from the OS is received.
var ErrInterrupt = errors.New("received interrupt")

// New returns a new ready-to-use Runner.
func New(d time.Duration) *Runner {
    return &Runner{
        interrupt: make(chan os.Signal, 1),
        complete:  make(chan error),
        timeout:   time.After(d),
    }
}

// Add attaches tasks to the Runner.
func (r *Runner) Add(tasks ...func(int)) {
    r.tasks = append(r.tasks, tasks...)
}

// Start runs all tasks and monitors channel events.
func (r *Runner) Start() error {
    // We want to receive all interrupt based signals.
    signal.Notify(r.interrupt, os.Interrupt)

    // Run the different tasks on a different goroutine.
    go func() {
        r.complete <- r.run()
    }()

    select {
    // Signaled when processing is done.
    case err := <-r.complete:
        return err

    // Signaled when we run out of time.
    case <-r.timeout:
        return ErrTimeout
    }
}

// run executes each registered task.
func (r *Runner) run() error {
    for id, task := range r.tasks {
        // Check for an interrupt signal from the OS.
        if r.gotInterrupt() {
            return ErrInterrupt
        }

        // Execute the registered task.
        task(id)
    }

    return nil
}

// gotInterrupt verifies if the interrupt signal has been issued.
func (r *Runner) gotInterrupt() bool {
    select {
    // Signaled when an interrupt event is sent.
    case <-r.interrupt:
        // Stop receiving any further signals.
        signal.Stop(r.interrupt)
        return true

    // Continue running as normal.
    default:
        return false
    }
}
```

**Sub-kod izahı:**
- `interrupt chan os.Signal` → OS siqnallarını (Ctrl+C) qəbul edir; **buffered (1)** — runtime siqnalı nonblocking göndərir, alan yoxdursa atılır (Ctrl+C bir neçə dəfə basılsa da 1 dənə qəbul olunur)
- `complete chan error` → unbuffered; tapşırıq goroutine-i bitəndə error/nil göndərir və main qəbul edənə qədər gözləyir — təmiz sonlanma
- `timeout <-chan time.Time` → `time.After(d)` factory — d müddətindən sonra time.Time göndərən kanal (receive-only elan!)
- `signal.Notify(r.interrupt, os.Interrupt)` → OS-dən interrupt siqnallarını kanala yönləndir
- `select { case ...: case ...: }` → 2 hadisədən birini gözləyir: complete və ya timeout — hansı birincı gəlirsə
- `select + default` (`gotInterrupt`-də) → **nonblocking peek**: siqnal yoxdursa dərhal `false` qaytarır, bloklamır
- `signal.Stop(r.interrupt)` → bir dəfə interrupt alındıqdan sonra sonrakı siqnalları qəbul etməyi dayandırır
- `Add(tasks ...func(int))` → variadic — istənilən sayda `func(int)` qəbul edir, slice-a append edir

**İstifadə (main.go):**
```go
const timeout = 3 * time.Second

r := runner.New(timeout)
r.Add(createTask(), createTask(), createTask())

if err := r.Start(); err != nil {
    switch err {
    case runner.ErrTimeout:
        log.Println("Terminating due to timeout.")
        os.Exit(1)
    case runner.ErrInterrupt:
        log.Println("Terminating due to interrupt.")
        os.Exit(2)
    }
}
log.Println("Process ended.")
```
- Error dəyişənləri (`ErrTimeout`, `ErrInterrupt`) `switch` ilə dəqiq identifikasiya olunur — Go-nun error dəyişəni konvensiyası
- Exit kodları: 0 normal, 1 timeout, 2 interrupt

### Pattern 2: Pooling — resurs hovuzu (buffered channel)
**Məqsəd:** DB connection, memory buffer kimi məhdud resursları çox goroutine arasında paylaşmaq: al → istifadə et → geri qaytar.

**Kitabdan kod nümunəsi (pool.go — tam paket):**
```go
// Package pool manages a user defined set of resources.
package pool

import (
    "errors"
    "log"
    "io"
    "sync"
)

// Pool manages a set of resources that can be shared safely by
// multiple goroutines. The resource being managed must implement
// the io.Closer interface.
type Pool struct {
    m         sync.Mutex
    resources chan io.Closer
    factory   func() (io.Closer, error)
    closed    bool
}

// ErrPoolClosed is returned when an Acquire returns on a closed pool.
var ErrPoolClosed = errors.New("Pool has been closed.")

// New creates a pool that manages resources.
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
    if size <= 0 {
        return nil, errors.New("Size value too small.")
    }

    return &Pool{
        factory:   fn,
        resources: make(chan io.Closer, size),
    }, nil
}

// Acquire retrieves a resource from the pool.
func (p *Pool) Acquire() (io.Closer, error) {
    select {
    // Check for a free resource.
    case r, ok := <-p.resources:
        log.Println("Acquire:", "Shared Resource")
        if !ok {
            return nil, ErrPoolClosed
        }
        return r, nil

    // Provide a new resource since there are none available.
    default:
        log.Println("Acquire:", "New Resource")
        return p.factory()
    }
}

// Release places a new resource onto the pool.
func (p *Pool) Release(r io.Closer) {
    // Secure this operation with the Close operation.
    p.m.Lock()
    defer p.m.Unlock()

    // If the pool is closed, discard the resource.
    if p.closed {
        r.Close()
        return
    }

    select {
    // Attempt to place the new resource on the queue.
    case p.resources <- r:
        log.Println("Release:", "In Queue")

    // If the queue is already at capacity we close the resource.
    default:
        log.Println("Release:", "Closing")
        r.Close()
    }
}

// Close will shutdown the pool and close all existing resources.
func (p *Pool) Close() {
    // Secure this operation with the Release operation.
    p.m.Lock()
    defer p.m.Unlock()

    // If the pool is already closed, don't do anything.
    if p.closed {
        return
    }

    // Set the pool as closed.
    p.closed = true

    // Close the channel before we drain the channel of its
    // resources. If we don't do this, we will have a deadlock.
    close(p.resources)

    // Close the resources
    for r := range p.resources {
        r.Close()
    }
}
```

**Sub-kod izahı:**
- `resources chan io.Closer` → **buffered channel = hovuzun özü**; interfeys tipi sayəsində istənilən io.Closer resursu idarə olunur
- `factory func() (io.Closer, error)` → istifadəçi tərəfindən təchiz edilən yeni resurs yaradıcı funksiya
- `Acquire`: `select + default` — hovuzda boş resurs varsa al (`case`), yoxdursa **factory ilə yeni yarat** (`default`)
- `Release`: `select + default` — hovuzda yer varsa geri qoy, doludursa resursu **bağla və at** (`r.Close()`)
- `Close`: mutex altında `closed` flag-i set et → **əvvəl kanalı bağla, sonra drain et** (əks halda deadlock!) → qalan resursları bağla
- **Mutex+flag zərurəti:** qapalı kanala send **panic** verir — `Release` və `Close` eyni mutex ilə qorunur; `closed` oxu/yazmaı sinxronlaşdırılır
- `r, ok := <-p.resources` → iki dəyərli receive: qapalı hovuzda `ok=false` → `ErrPoolClosed`

**İstifadə (main.go — simulyasiya edilmiş DB connection):**
```go
const (
    maxGoroutines   = 25 // the number of routines to use.
    pooledResources = 2 // number of resources in the pool
)

type dbConnection struct {
    ID int32
}

func (dbConn *dbConnection) Close() error {
    log.Println("Close: Connection", dbConn.ID)
    return nil
}

var idCounter int32

func createConnection() (io.Closer, error) {
    id := atomic.AddInt32(&idCounter, 1)
    log.Println("Create: New Connection", id)
    return &dbConnection{id}, nil
}

func main() {
    var wg sync.WaitGroup
    wg.Add(maxGoroutines)

    p, err := pool.New(createConnection, pooledResources)

    for query := 0; query < maxGoroutines; query++ {
        go func(q int) {
            performQueries(q, p)
            wg.Done()
        }(query)
    }

    wg.Wait()
    p.Close()
}

func performQueries(query int, p *pool.Pool) {
    conn, err := p.Acquire()
    if err != nil {
        log.Println(err)
        return
    }

    // Release the connection back to the pool.
    defer p.Release(conn)

    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
    log.Printf("QID[%d] CID[%d]\n", query, conn.(*dbConnection).ID)
}
```

**Sub-kod izahı:**
- `atomic.AddInt32(&idCounter, 1)` → unikal ID generator — atomik artırma
- `defer p.Release(conn)` → funksiya bitəndə connection hovuza avtomatik qayıdır
- `conn.(*dbConnection).ID` → **type assertion** — io.Closer-dən konkret tipə çevrilib ID-yə çıxış
- 25 goroutine 2 bağlamanı paylaşır — pool-log-da "Shared Resource" / "New Resource" nüansları görünür

### Pattern 3: Work — goroutine hovuzu (unbuffered channel)
**Məqsəd:** Ştat sayda goroutine-dən ibarət işçi pool; unbuffered channel **qarantiya** verir ki, `Run` qayıdanda iş mütləq bir goroutine tərəfindən **gözlənilir** — iş itmir, növbədə ilişmir.

**Kitabdan kod nümunəsi (work.go — tam paket):**
```go
// Package work manages a pool of goroutines to perform work.
package work

import "sync"

// Worker must be implemented by types that want to use the work pool.
type Worker interface {
    Task()
}

// Pool provides a pool of goroutines that can execute any Worker
// tasks that are submitted.
type Pool struct {
    work chan Worker
    wg   sync.WaitGroup
}

// New creates a new work pool.
func New(maxGoroutines int) *Pool {
    p := Pool{
        work: make(chan Worker),
    }

    p.wg.Add(maxGoroutines)
    for i := 0; i < maxGoroutines; i++ {
        go func() {
            for w := range p.work {
                w.Task()
            }
            p.wg.Done()
        }()
    }

    return &p
}

// Run submits work to the pool.
func (p *Pool) Run(w Worker) {
    p.work <- w
}

// Shutdown waits for all the goroutines to shutdown.
func (p *Pool) Shutdown() {
    close(p.work)
    p.wg.Wait()
}
```

**Sub-kod izahı:**
- `Worker` interfeysi → tək `Task()` metodu — hər hansı iş bu interfeyslə təqdim olunur
- `make(chan Worker)` → **unbuffered!** — `Run` yalnız bir işçi qəbul edəndə qayıdır: işin götürüldüyü qarantiyası
- `for w := range p.work` → worker goroutine-lər kanalı dinləyir; kanal bağlananda loop bitir → `wg.Done()` → goroutine ölür
- `Run(w)` → tapşırığı kanala göndərir — işçi hazır olana qədər bloklanır (backpressure!)
- `Shutdown()` → `close` + `Wait` — bütün işlər bitənə qədər gözlə

**İstifadə (main.go):**
```go
var names = []string{
    "steve", "bob", "mary", "therese", "jason",
}

type namePrinter struct {
    name string
}

// Task implements the Worker interface.
func (m *namePrinter) Task() {
    log.Println(m.name)
    time.Sleep(time.Second)
}

func main() {
    // Create a work pool with 2 goroutines.
    p := work.New(2)

    var wg sync.WaitGroup
    wg.Add(100 * len(names))

    for i := 0; i < 100; i++ {
        for _, name := range names {
            np := namePrinter{
                name: name,
            }

            go func() {
                // Submit the task to be worked on.
                p.Run(&np)
                wg.Done()
            }()
        }
    }

    wg.Wait()
    p.Shutdown()
}
```
- 500 goroutine yaradılır, hamısı `p.Run(&np)` üçün yarışır; 2 işçi tədricən emal edir
- Unbuffered olduğundan `Run` qayıdanda iş **mütləq götürülüb** — "növbədə itən iş" anlayışı yoxdur

## Üç pattern-in müqayisəsi
| Pattern | Channel növü | İstifadə sahəsi | Əsas mexanizm |
|---|---|---|---|
| Runner | buffered(1) + unbuffered + time.After | Proqram lifetime, cron, timeout | select, signal.Notify |
| Pool | buffered (size) | Resurs paylaşımı (DB conn) | select+default, mutex+flag |
| Work | unbuffered | İş paylama (worker pool) | for range, backpressure |

## Əsas terminlər
- select statement (çoxşaxəli seçim)
- default case (nonblocking əməliyyat üçün)
- signal.Notify / signal.Stop (OS siqnal idarəsi)
- time.After (vaxt bitən kanalı)
- Variadic Parameter (`...T` — dəyişən saylı arqument)
- Factory Function (fabrik funksiyası)
- Resource Pool (resurs hovuzu)
- Worker Pool (işçi hovuzu)
- Backpressure (geri təzyiq — kanalın doluluq siqnalı)
- Type Assertion (`x.(T)` — tip müəyyənləşdirmə)
- Error Variable (identifikasiya oluna bilən xəta dəyişəni)
- Deadlock (qarşılıqlı bloklanma)

## Praktik nəticə
- Nəzarətsiz işləyən proqramlarda (cron, cloud worker) Runner pattern ilə timeout + təmiz interrupt söndürmə qur.
- `select`-in `default`-u kanalı **peek** etmək üçün — bloklamayan yoxlama (gotInterrupt nümunəsi).
- Resurs hovuzu: buffered channel + `select/default` + mutex-qorunan closed flag; kanalı **əvvəl bağla, sonra drain et** — əks halda deadlock.
- İşçi hovuzu üçün unbuffered channel seç — `Run` qayıdanda işin götürüldüyü qarantiyası itən növbə yaratmır.
- Xətaları `errors.New` ilə paket səviyyəli dəyişənlər kimi elan et — çağırıcı `switch err` ilə dəqiq identifikasiya edə bilsin.
- `defer`-i resurs əldə etdikdən dərhal sonra `Release` üçün yaz — ac 3 pattern-də də bu vərdişdir.

## Mənbə
Pages: 179-204 (PDF), book pages 158-183
