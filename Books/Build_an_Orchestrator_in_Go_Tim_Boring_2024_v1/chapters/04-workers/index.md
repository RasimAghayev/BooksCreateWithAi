# Chapters 4-6 — Worker, Worker API, Metrics (səh. 104-177)

## Bu fəsillər nədən bəhs edir?

(4) Worker-in daxili: queue + DB + runtime, task state machine keçid cədvəli,
RunTask/StartTask/StopTask implementasiyası. (5) chi router ilə REST API:
handlers, routes, JSON decode (DisallowUnknownFields), curl ilə test. (6)
/proc filesystem-dən metrics: goprocinfo ilə memory/disk/CPU, /stats
endpoint-i.

## Əsas fikirlər

### 1. Worker strukturu və axını (Ch4)
```
API → Queue (müvəqqəti tutum) → RunTask (döngü) → StartTask/StopTask → DB (state)
```
- Worker = runtime (konteyner işlətmə) + queue + DB + metrics + API
- Task ↔ konteyner BİRBİRƏ (1:1)

### 2. State transition map — keçidlərin qanunu
```go
var stateTransitionMap = map[State][]State{
    Pending:   {Scheduled},
    Scheduled: {Scheduled, Running, Failed},
    Running:   {Running, Completed, Failed},
    Completed: {},
    Failed:    {},
}

func Contains(list []State, s State) bool { ... }
func ValidStateTransition(src, dst State) bool {
    return Contains(stateTransitionMap[src], dst)
}
```
- Cədvəl YOXSA sizin valid keçidlər haqqında intuksiyanız hər yerdə səpələnir;
  birtərəfli mənbə = yoxlanıla bilən kontrakt

### 3. RunTask — dispatch
```go
func (w *Worker) RunTask() task.DockerResult {
    t := w.Queue.Dequeue()
    if t == nil { return task.DockerResult{} }
    taskQueued := t.(task.Task)   // type assertion
    switch taskQueued.State {
    case Pending:  return w.StartTask(taskQueued)
    case Running:  return w.StopTask(taskQueued)
    default:       fmt.Println(...)  // ya restart mantığı (sonra)
    }
}
```
- StartTask: StartTime yaz → Docker.Run() → DB-yə state Running
- StopTask: Docker.Stop() → FinishTime → DB-yə Completed

### 4. chi router API (Ch5)
```go
r := chi.NewRouter()
r.Route("/tasks", func(r chi.Router) {
    r.Post("/", a.StartTaskHandler)   // 201 + yaradılmış task
    r.Get("/", a.GetTasksHandler)     // 200 + siyahı
    r.Route("/{taskID}", func(r chi.Router) {
        r.Delete("/", a.StopTaskHandler)  // 204
    })
})
```
- Handler-lər task ƏMƏLİYYATI TİRMİRLƏR — sadəcə queue-ya qoyur + cavab;
  ayrılan concern: request qəbulu vs icra
- Decode: `d := json.NewDecoder(r.Body); d.DisallowUnknownFields()` — sərt
  parsinq (400 + ErrResponse)
- `{taskID}` → `chi.URLParam(r, "taskID")` → `uuid.Parse()`

### 5. Metrics — /proc (Ch6)
| Fayl | Məzmun |
|---|---|
| /proc/meminfo | MemTotal, MemFree, MemAvailable (kB) |
| /proc/diskstats | disk əməliyyat statistikası |
| /proc/stat | cpu user/nice/system/idle... |
| /proc/loadavg | 1/5/15 dəq yük ortalaması |

```go
type Stats struct {
    MemStats  *linux.MemInfo
    DiskStats *disk.UsageStats
    CpuStats  *linux.CPUStat
    LoadStats *loadavg.LoadAvg
}

func (s *Stats) MemAvailableKb() uint64 { return s.MemStats.MemAvailable }
func (s *Stats) CpuPercentUsed() float64 {
    idle := s.CpuStats.Idle + s.CpuStats.IOWait
    total := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System +
             idle + s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
    return (float64(total) - float64(idle)) / float64(total)
}
```
- goprocinfo + syscall.DiskUsage — /proc-un Go qatışı
- Worker dövrü: `go w.CollectStats()` — periyodik toplayır
- **/stats endpoint:** scheduler "ən yaxşı node" seçimi üçün bunu oxuyur

## Əsas terminlər

- State Transition Table (vəziyyət keçid cədvəli)
- Queue (FIFO) / Dequeue
- chi router / URLParam
- DisallowUnknownFields (sərt JSON)
- /proc filesystem
- load average (1/5/15 dəq)
- CpuPercentUsed (idle əsaslı formula)

## Praktik nəticə

- State machine KEÇİDLƏRİ cədvəldə tək mənbədə saxla
- API handler-ləri icra etməsin — queue + dərhal cavab (ayrılan concern)
- Metrics: sadə /proc oxunuşu + /stats JSON — scheduler-in qida mənbəyi

## Mənbə

Pages: 104-177 (Chapters 4-6, Build an Orchestrator in Go)
