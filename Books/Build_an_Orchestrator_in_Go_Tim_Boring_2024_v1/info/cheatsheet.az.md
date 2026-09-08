# Build an Orchestrator in Go — Cheat Sheet (AZ)

## Task state machine

```go
type State int
const (
    Pending State = iota
    Scheduled
    Running
    Completed
    Failed
)

var stateTransitionMap = map[State][]State{
    Pending:   {Scheduled},
    Scheduled: {Scheduled, Running, Failed},
    Running:   {Running, Completed, Failed},
    Completed: {},
    Failed:    {},
}

func ValidStateTransition(src, dst State) bool {
    return Contains(stateTransitionMap[src], dst)
}
```

## Task strukturu

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
    HealthCheck   string      // "/health"
    RestartCount  int         // limit: 3
    StartTime     time.Time
    FinishTime    time.Time
}

type TaskEvent struct {
    ID        uuid.UUID
    State     State
    Timestamp time.Time
    Task      Task
}
```

## Docker SDK

```go
cli, _ := client.NewClientWithOpts(client.FromEnv)
out, _ := cli.ImagePull(ctx, image, types.ImagePullOptions{})
io.Copy(os.Stdout, out)

resp, _ := cli.ContainerCreate(ctx,
    &container.Config{Image: img, Env: env, ExposedPorts: ports},
    &container.HostConfig{
        RestartPolicy: container.RestartPolicy{Name: "none"},
        Resources: container.Resources{Memory: mem, NanoCPUs: cpu},
        PortBindings: bindings,
    }, nil, nil, name)
cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
// Stop: cli.ContainerStop(ctx, id, timeout); cli.ContainerRemove(ctx, id, opts{})
// Inspect: cli.ContainerInspect(ctx, id) → .State.Status, .NetworkSettings...
```

## Worker API (chi)

```go
r := chi.NewRouter()
r.Route("/tasks", func(r chi.Router) {
    r.Post("/", a.StartTaskHandler)   // 201
    r.Get("/", a.GetTasksHandler)     // 200
    r.Route("/{taskID}", func(r chi.Router) {
        r.Delete("/", a.StopTaskHandler)  // 204
    })
})

// sərt JSON:
d := json.NewDecoder(r.Body)
d.DisallowUnknownFields()
err := d.Decode(&te)
// path param:
taskID := chi.URLParam(r, "taskID")
tID, _ := uuid.Parse(taskID)
```

## Metrics (/proc)

```go
type Stats struct {
    MemStats  *linux.MemInfo
    DiskStats *disk.UsageStats
    CpuStats  *linux.CPUStat
    LoadStats *loadavg.LoadAvg
}
// CPU% (2 snapshot, 3s arayla):
total := stat2Total - stat1Total
idle  := stat2Idle - stat1Idle
usage := (total - idle) / total
```

## Scheduler interfeysi + E-PVM

```go
type Scheduler interface {
    SelectCandidateNodes(t task.Task) []*node.Node
    Score(t task.Task, nodes []*node.Node) map[string]float64
    PickNode(scores map[string]float64) string
}

// E-PVM score:
cpuLoad := calculateLoad(cpuUsage, math.Pow(2, 0.8))
memPercent := (memUsageNorm + memAllocatedNorm) / 2
scores[node.Name] = 0.6*cpuLoad + 0.4*memPercent  // ağırlıqlı ortaalı mənada
// Pick: minimum score
```

## Health check dövrü

```go
func (m *Manager) doHealthChecks() {
    for _, t := range m.GetTasks() {
        if t.State == Running && t.RestartCount < 3 {
            resp, err := http.Get(healthURL(t))
            if err != nil || resp.StatusCode != 200 {
                m.restartTask(t)   // RestartCount++
            }
        }
    }
}
// Periodik: DoHealthChecks (60s), UpdateTasks (15s), ProcessTasks (10s)
```

## Store interfeysi + BoltDB

```go
type Store interface {
    Put(key string, value interface{}) error
    Get(key string) (interface{}, error)
    List() (interface{}, error)
    Count() (int, error)
}

// BoltDB:
db, _ := bolt.Open("tasks.db", 0600, nil)
db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte("tasks"))
    data, _ := json.Marshal(t)
    return b.Put([]byte(key), data)
})
db.View(func(tx *bolt.Tx) error {
    v := tx.Bucket([]byte("tasks")).Get([]byte(key))
    return json.Unmarshal(v, &t)
})
// switch: "memory" → InMemory...; "persistent" → BoltDB TaskStore
```

## Cobra CLI

```go
var runCmd = &cobra.Command{
    Use: "run",
    Run: func(cmd *cobra.Command, args []string) {
        manager, _ := cmd.Flags().GetString("manager")
        file, _ := cmd.Flags().GetString("filename")
        data, _ := os.ReadFile(file)
        resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
        // 201 yoxla
    },
}
func init() {
    rootCmd.AddCommand(runCmd)
    runCmd.Flags().StringP("manager", "m", "localhost:5555", "...")
    runCmd.Flags().StringP("filename", "f", "task.json", "...")
}
```

## Əsas dövrlər (hamısı goroutine)

```bash
go w.RunTasks()       # queue → Start/Stop
go w.CollectStats()    # /proc → memory dövrü
go w.UpdateTasks()     # ContainerInspect → state sync (15s)
go api.Start()         # chi server
m.ProcessTasks()       # SendWork (10s)
m.DoHealthChecks()     # sağlamlıq (60s)
```

## Manager strukturu (yekun)

```go
type Manager struct {
    Pending       queue.Queue
    TaskDb        store.Store     // interfeys!
    EventDb       store.Store
    Workers       []string
    WorkerTaskMap map[string][]uuid.UUID
    TaskWorkerMap map[uuid.UUID]string
    WorkerNodes   []*node.Node
    Scheduler     scheduler.Scheduler  // interfeys!
}
```
