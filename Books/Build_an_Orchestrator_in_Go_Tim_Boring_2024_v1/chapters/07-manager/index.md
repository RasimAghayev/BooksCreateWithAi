# Chapters 7-9 — Manager, Manager API, Failures (səh. 178-259)

## Bu fəsillər nədən bəhs edir?

(7) Manager komponenti: naive round-robin scheduler, SendWork, UpdateTasks,
AddTask; control vs data plane. (8) Manager REST API (worker API-yə paralel).
(9) Xəta senariləri: app/task/worker/manager səviyyələri, Docker restart
policy vs orkestrator health check-ləri, RestartCount limiti.

## Əsas fikirlər

### 1. Manager strukturu (Ch7)
```go
type Manager struct {
    Pending        queue.Queue
    TaskDb         map[uuid.UUID]*task.Task
    EventDb        map[uuid.UUID]*task.TaskEvent
    Workers        []string
    WorkerTaskMap  map[string][]uuid.UUID
    TaskWorkerMap  map[uuid.UUID]string
}
```
- **Control plane** (manager — idarə) vs **data plane** (worker — icra):
  şəbəkə dünyasından gələn ayırma

### 2. Naive scheduler — round-robin
```go
func (m *Manager) SelectWorker() string {
    var newWorker int
    if m.LastWorker < len(m.Workers)-1 {
        newWorker = m.LastWorker + 1
    }
    m.LastWorker = newWorker
    return m.Workers[newWorker]
}
```
- Sadə amma kor: worker-də resurs YOXDURSA nə olacaq? (→ Ch10 həll)

### 3. SendWork — manager-in "iş atı"
```
Pending queue → dequeue → SelectWorker → EventDb/Map-lər yaz →
task.Scheduled → POST http://worker/tasks (JSON) → 201 yoxla → xəta halında
ErrResponse decode
```

### 4. UpdateTasks — worker-dən dövlət sinxronu
- Hər worker-in /tasks GET → TaskDb state/Start/Finish/ContainerID yenilə
- Manager periodik döngələrlə: ProcessTasks (10s), UpdateTasks, sonra
  DoHealthChecks (60s)

### 5. Manager API (Ch8)
- Eyni qəlib: chi router, POST/GET/DELETE /tasks — amma fərq:
  StartTaskHandler → AddTask (pending queue-yə) + 201
  StopTaskHandler → nüsxə task Completed + AddTask (worker stop əmri kimi)
- RunTasks funksiyası worker-ə köçürüldü (encapsulation refactor)

### 6. Xəta təsnifatı (Ch9)
| Səviyyə | Nümunə | Bərpa |
|---|---|---|
| App start | DB çatışmır | retry + exponential backoff |
| Task crash | proses exit | orkestrator restart (Docker policy YOX) |
| Worker | API ölü / machine crash | manager yenidən başladır |
| Manager | proses ölür | systemd/restart; datastore backup |

- **Docker --restart NİYƏ YARAMIR:** məsuliyyət bulanıqlaşır — Docker-mı,
  orkestratormu? Orkestrator sistemdə restart məntiqi MANAGER-də olmalıdır

### 7. Health check mexanizmi
```go
// Task configə HealthCheck sahəsi: "/health"
func (m *Manager) checkTaskHealth(t task.Task) error {
    url := fmt.Sprintf("http://%s:%s%s", host, port, t.HealthCheck)
    resp, err := http.Get(url)
    if err != nil { return err }
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("healthcheck failed: %d", resp.StatusCode)
    }
    return nil
}

func (m *Manager) doHealthChecks() {
    for _, t := range m.GetTasks() {
        if t.State == task.Running && t.RestartCount < 3 {
            if err := m.checkTaskHealth(*t); err != nil {
                m.restartTask(t)   // RestartCount++ ilə
            }
        }
    }
}
```
- Worker UpdateTasks: `ContainerInspect` → `State.Status == "exited"` →
  task.Failed; HostPorts network-dən çəkilir
- **RestartCount < 3 qanunu:** sonsuz restart loop-un qarşısı — 3 cəhddən
  sonra task Failed saxlanılır (insan müdaxiləsi gözlənilir)

## Əsas terminlər

- Control Plane / Data Plane
- Round-robin scheduling
- SendWork / UpdateTasks
- Exponential Backoff
- Health Check endpoint
- RestartCount limiti
- ContainerInspect

## Praktik nətiqə

- Manager: qərar verən (hansı worker, restart etmək?), worker: icra edən
- Restart məsuliyyəti orkestratorda — Docker policy-ə buraxma
- Health checklər = sadə GET /health + status 200 + cəhd limiti
- API hər komponentdə eyni qəlibdə (chi + handler + queue) — bir dəfə öyrən,
  hər yerdə tanı

## Mənbə

Pages: 178-259 (Chapters 7-9, Build an Orchestrator in Go)
