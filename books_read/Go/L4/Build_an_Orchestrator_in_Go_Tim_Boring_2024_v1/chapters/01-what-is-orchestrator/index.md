# Chapters 1-3 — Konsept, Skelet, Task + Docker (səh. 24-103)

## Bu fəsillər nədən bəhs edir?

(1) Konteyner/orkestrator anlayışları, komponentlər (task/job/worker/manager/
scheduler/cluster), real sistemlər (Borg/Mesos/K8s/Nomad), Cube layihəsinin
mental modeli. (2) Task state machine, Task/TaskEvent/Worker/Manager/Node
skeletləri, scheduler interfeysi. (3) Docker Engine API ilə konteyner
yaradılma/başlatma/stop — Task-ın "əti".

## Əsas fikirlər

### 1. Konteyner vs VM (Ch1)
- **VM:** öz kernel + network stack + resurs nəzarəti
- **Konteyner:** YALNIZ konsept — kernel YOXDUR; namespace + cgroups üzərində
  proses izolyasiyası (Liz Rice "Containers from Scratch")
- Orkestrator = konteynerlərin avtomatik deploy/scale/idarəsi

### 2. Orkestrasiya komponentləri (universal)
| Komponent | Borc | K8s | Nomad | Borg |
|---|---|---|---|---|
| Manager | beyin, giriş nöqtəsi | control plane | server | BorgMaster |
| Worker | işi icra | kubelet | client | Borglet |
| Scheduler | task → node seçimi | kube-scheduler | (serverdə) | (BorgMasterdə) |
| Task | ən kiçik iş vahidi | pod | task | alloc |

- **Task spesifikasiyası:** resurs ehtiyacı (CPU/mem/disk), image, restart
  siyasəti, portlar
- **Job:** yüksək səviyyə — neçə instans, datacenter, update siyasəti

### 3. Cube mental modeli
- Scheduler: **feasibility (uyğunluq) → scoring (ballandırma) → picking
  (seçim)** mərhələləri
- Manager: pending queue → scheduler → worker-ə təhvil; task/event storage
- Worker: taskları run/stop; stats toplayır; API-si var
- Sadəlik: kanalsız, genericsiz — öyrənmə optimallaşdırması

### 4. Task state machine (Ch2)
```go
type State int

const (
    Pending State = iota
    Scheduled
    Running
    Completed
    Failed
)
```
- Hər keçid = TaskEvent (ID, State, Timestamp, Task) — audit izi

### 5. Task strukturu
```go
type Task struct {
    ID            uuid.UUID
    ContainerID   string
    Name          string
    State         State
    Image         string
    CPU           float64
    Memory        int64
    Disk          int64
    ExposedPorts  nat.PortSet
    PortBindings  map[string]string
    RestartPolicy string
    StartTime     time.Time
    FinishTime    time.Time
}
```

### 6. Worker / Manager / Node skeletləri
```go
// Worker: DB (task→container), Queue (FIFO event-lər), TaskCount
// Manager: pending queue, EventDb, WorkerTaskMap, TaskWorkerMap
//   selectWorker(), sendWork(), updateTasks(), healthcheck()...
// Node: Name, Ip, Cores, Memory/Allocated, Disk/Allocated, Role, TaskCount
```
- Manager map-ləri: hansı task HANSI worker-də, worker-də HANSI tasklar

### 7. Scheduler interfeysi (round-robin əvəzinə)
```go
type Scheduler interface {
    SelectCandidateNodes(t task.Task) []*Node  // feasibility
    Score(t task.Task, nodes []*Node) map[string]float64  // scoring
    PickNode(scores map[string]float64) string  // picking
}
```

### 8. Docker Engine API (Ch3)
- CLI əvəzinə SDK: `docker run -it -p 5432:5432 --name cube-book -e ... postgres`
  ↔ `NewClientWithOpts`, `ImagePull`, `ContainerCreate`, `ContainerStart`
- Hər SDK metodu `context.Context` qəbul edir
- Task config → Docker config çevrilməsi:
```go
func (d *Docker) Run() DockerResult {
    // ImagePull (io.Copy ilə stdout)
    rp := container.RestartPolicy{Name: d.Config.RestartPolicy}
    r := container.Resources{Memory: d.Config.Memory, NanoCPUs: ...}
    resp, err := dc.ContainerCreate(ctx, &container.Config{...}, &container.HostConfig{...}, ...)
    dc.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}
// Stop: ContainerStop + ContainerRemove
```
- **DockerResult:** {Error, Action, ContainerId, Result} — worker-in
  sorğu-cavab vahidi

## Əsas terminlər

- Namespace / cgroups (konteynerin əsası)
- Task / Job / Worker / Manager / Scheduler / Cluster
- Feasibility / Scoring / Picking
- TaskEvent (state keçid jurnalı)
- Docker Engine API / SDK
- ContainerCreate / ContainerStart / ContainerStop

## Praktik nəticə

- Bütün orkestratorlar eyni komponent dəstinə gəlir — adlar fərqlidir
- Task = state machine + config + Docker binding-ləri
- Skelet kodla başla (psixoloji + kompayl təsdiqi), sonra "ət" əlavə et

## Mənbə

Pages: 24-103 (Chapters 1-3, Build an Orchestrator in Go)
