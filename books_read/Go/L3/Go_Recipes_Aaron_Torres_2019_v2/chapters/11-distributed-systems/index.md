# Chapter 11 — Распределенные системы (Distributiv sistemlər)

## Bu chapter nədən bəhs edir?

Consul service discovery, Raft konsensusu (in-memory klaster + FSM), Docker
konteynerləşdirməsi, Docker Compose orkestrasiyası, Prometheus monitorinqi
və go-metrics ilə custom metric toplanması.

## Əsas fikirlər

### 1. Service discovery — Consul
**Nədir:** Mikroservislərin ünvanlarını dinamik qeydiyyatı/axtarışı — IP/port
statik konfiqurasiyasını aradan qaldırır; network partition hallarını da
göstərir.

**Nəyə lazımdır:** Servis sayı artdıqda və mühitlər (staging/prod) dəyişəndə
ünvanları əl ilə saxlamaq mümkün deyil.

**Kitabdan kod nümunəsi:**
```go
// Client interfeysi — mock üçün:
type Client interface {
    Register(tags []string) error
    Service(service, tag string) ([]*api.ServiceEntry, *api.QueryMeta, error)
}

func NewClient(config *api.Config, address, name string, port int) (Client, error) {
    c, err := api.NewClient(config)
    // ...
    return &client{client: c, name: name, address: address, port: port}, nil
}

// Servisi qeyd et:
func (c *client) Register(tags []string) error {
    reg := &api.AgentServiceRegistration{
        ID:      c.name,
        Name:    c.name,
        Port:    c.port,
        Address: c.address,
        Tags:    tags,
    }
    return c.client.Agent().ServiceRegister(reg)
}

// Servisi tap (tag ilə filtrlə):
func (c *client) Service(service, tag string) ([]*api.ServiceEntry, *api.QueryMeta, error) {
    return c.client.Health().Service(service, tag, false, nil)
}

// İstifadə: consul agent -dev -node=localhost
cli.Register([]string{"Go", "Awesome"})
entries, _, _ := cli.Service("discovery", "Go")
```

**Sub-kod izahı:**
- Qeydiyyat servis başlayanda, silinmə shutdown-da
- Health().Service → sağlamlıq yoxlaması ilə axtarış
- Nəticələri cache-ləmək olar — hər sorğu üçün Consul-a getmə
- Consul agent yerli sorğuları tez edir

### 2. Raft konsensusu
**Nədir:** Distributiv sistemdə ortaq vəziyyət — leader seçkiləri + log
replikasiyası; quorum = (n/2)+1 node.

**Necə işləyir:** HSR-in hashicorp/raft kitabxanası: in-memory transport,
FSM (finite state machine) Apply metodunu çağırır.

**Kitabdan kod nümunəsi:**
```go
// State machine — icazəli keçidlər:
type state string
const (first state = "first"; second = "second"; third = "third")

var allowedState map[state][]state
func init() {
    allowedState = make(map[state][]state)
    allowedState[first] = []state{second, third}
    allowedState[second] = []state{third}
    allowedState[third] = []state{first}
}
func (s *state) CanTransition(next state) bool {
    for _, n := range allowedState[*s] {
        if n == next { return true }
    }
    return false
}
func (s *state) Transition(next state) {
    if s.CanTransition(next) { *s = next }
}

// FSM — raft interfeysi:
type FSM struct{ state state }
func (f *FSM) Apply(r *raft.Log) interface{} {
    f.state.Transition(state(r.Data))     // log datası = yeni state
    return string(f.state)
}
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) { return nil, nil }
func (f *FSM) Restore(io.ReadCloser) error { return nil }

// In-memory klaster (3 node):
func Config(num int) {
    rs := getRaftSet(num)               // in-memory store + transport
    for _, r1 := range rs {
        for _, r2 := range rs {
            r1.Transport.Connect(r2.Transport.LocalAddr(), r2.Transport)
        }
    }
    for _, r := range rs {
        raft.BootstrapCluster(r.Config, r.Store, r.Store,
            r.SnapShotStore, r.Transport, r.Configuration)
        r, err := raft.NewRaft(r.Config, r.FSM, r.Store, r.Store,
            r.SnapShotStore, r.Transport)
        rafts[r.Transport.LocalAddr()] = r
    }
}

// Handler — YALNIZ leader-ə Apply:
func Handler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    state := r.FormValue("next")
    for address, raft := range rafts {
        if address != raft.Leader() { continue }   // leader-i tap
        result := raft.Apply([]byte(state), 1*time.Second)
        newState, ok := result.Response().(string)
        if newState != state {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("invalid transition"))
            return
        }
        w.Write([]byte(newState))   // 200
        return
    }
}
```

**Sub-kod izahı:**
- Log daxil olma → FSM.Apply BÜTÜN node-larda icra olunur — ortaq vəziyyət
- `raft.Apply(data, timeout)` → leader-ə yazma; kvorum təsdiqindən sonra
- In-memory transport → test/POC üçün; produksiyada HTTP transport + service
  discovery lazımdır
- 3 node minimum — 1 node-ın ölümündən sağ qalmaq üçün (quorum 2/3)

### 3. Docker ilə konteynerləşdirmə
**Nədir:** Tətbiqi portable konteynerə qablaşdırmaq — VM-in faydaları,
yüngül konteynerdə.

**Kitabdan kod nümunəsi:**
```dockerfile
FROM alpine                    # minimal baza imici
ADD ./example/example /example # statik binary
EXPOSE 8000
ENTRYPOINT /example
```
```bash
# setup.sh — build-time dəyişənlərlə Linux binary:
env GOOS=linux go build -ldflags "-X main.version=1.0 -X main.builddate=$(date +%s)"
docker build . -t example
docker run -d -p 8000:8000 example
```
```go
// Versiya handler — build-time -X flag-ləri ilə:
var (
    version   string   // -X main.version
    builddate string   // -X main.builddate
)
type VersionInfo struct {
    Version   string
    BuildDate time.Time
    Uptime    time.Duration
}
func VersionHandler(v *VersionInfo) http.HandlerFunc {
    t := time.Now()
    return func(w http.ResponseWriter, r *http.Request) {
        v.Uptime = time.Since(t)
        vers, _ := json.Marshal(v)
        w.Write(vers)
    }
}
```

**Sub-kod izahı:**
- Go statik binary → alpine (çox kiçik) kifayətdir
- `-ldflags "-X main.var=dəyər"` → versiya/tarix binary-yə DAXİL EDİLİR
- `-p 8000:8000` → host↔container port xəritələndirməsi
- curl localhost:8000/version → `{"Version":"1.0","BuildDate":...}`

### 4. Orkestrasiya — Docker Compose
**Nədir:** Çoxkonteynerli lokal mühit — app + MongoDB bir komanda ilə.

**Kitabdan kod nümunəsi:**
```yaml
# docker-compose.yml
version: '2'
services:
  app:
    build: .           # Dockerfile-dan
  mongodb:
    image: "mongo:latest"   # hazır imic
```
```dockerfile
FROM golang:1.12.4-alpine3.9
ENV GOPATH /code/
ADD . /code/src/github.com/.../docker
WORKDIR /code/src/.../example
RUN GO111MODULE=on GOPROXY=off go build -mod=vendor   # vendor ilə offline build
ENTRYPOINT /code/src/.../example/example
```
```go
// App MongoDB-yə SERVİS ADI ilə qoşulur:
mongodb.Exec("mongodb://mongodb:27017")   // "mongodb" = compose service adı
```

**Sub-kod izahı:**
- `go mod vendor` → asılılıqları vendor/ qovluğuna; offline build mümkün
- Compose şəbəkəsində servislər bir-birini AD ilə görür (DNS)
- Lokal dev üçün ideal; produksiyada DB konteynerdə olmaya bilər
- Növbəti səviyyə: Docker Swarm (klaster, scale, load-balancing) /
  Kubernetes (Google-in Go-da yazdığı orkestrator)

### 5. Monitorinq — Prometheus
**Nədir:** Pull-model monitorinq — Prometheus müntəzəm /metrics endpoint-ini
sorğulayır (goroutine sayı, yaddaş və s.).

**Kitabdan kod nümunəsi:**
```go
// Go tətbiqi — 3 sətirlik exposition:
import "github.com/prometheus/client_golang/prometheus/promhttp"
http.Handle("/metrics", promhttp.Handler())
http.ListenAndServe(":80", nil)
```
```yaml
# docker-compose.yml
services:
  app:
    build: .
  prometheus:
    ports: ["9090:9090"]
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    image: "prom/prometheus"
```
```yaml
# prometheus.yml
global:
  scrape_interval: 15s
scrape_configs:
  - job_name: 'app'
    scrape_interval: 5s         # hər 5 saniyə
    static_configs:
      - targets: ['app:80']     # compose şəbəkəsindəki app
```

**Sub-kod izahı:**
- Pull model → bir neçə Prometheus server-i eyni app-i izləyə bilər, app
  yenidən deploy edilmir
- promhttp.Handler() → hazır Go runtime statistikası (yaddaş sızmalarının
  tapılması üçün dəyərli)
- Prometheus web UI: http://localhost:9090

### 6. Custom metriclər (go-metrics)
**Nədir:** Tətbiq-spesifik counter/timer metricləri — registry-də toplanır,
JSON endpoint-də çap olunur.

**Kitabdan kod nümunəsi:**
```go
// Counter — hər çağırışda artır:
func CounterHandler(w http.ResponseWriter, r *http.Request) {
    c := metrics.GetOrRegisterCounter("counterhandler.counter", nil)
    c.Inc(1)
    w.Write([]byte("success"))
}

// Timer — handler müddəti:
func TimerHandler(w http.ResponseWriter, r *http.Request) {
    currt := time.Now()
    t := metrics.GetOrRegisterTimer("timerhandler.timer", nil)
    w.Write([]byte("success"))
    t.UpdateSince(currt)           // başlanğıcdan bəri vaxt
}

// Report — bütün registry JSON:
func ReportHandler(w http.ResponseWriter, r *http.Request) {
    t := gometrics.GetOrRegisterTimer("reporthandler.writemetrics", nil)
    t.Time(func() {                // funksiyanın öz müddətini ölç
        gometrics.WriteJSONOnce(gometrics.DefaultRegistry, w)
    })
}
// curl /report →
// {"counterhandler.counter":{"count":1},
//  "timerhandler.timer":{"count":1,"max":60485,"mean":60485,...},
//  "reporthandler.writemetrics":{...percentile-lər...}}
```

**Sub-kod izahı:**
- `GetOrRegister*` → atomik "getir YA-da yarad" — thread-safe
- Timer → count, min/max/mean, percentile-lər (75%, 95%, 99%), rate-lər
- WriteJSONOnce → registry-nin snapshot-ı; exporter-lərlə Prometheus/Influx
  DB-yə də göndərmək olar

## Əsas terminlər

- Service Discovery (servis kəşfi)
- Consul Agent / Health Check
- Consensus (konsensus) / Raft
- Leader Election (lider seçkisi)
- Quorum ((n/2)+1)
- FSM (Finite State Machine — sonlu vəziyyət maşını)
- In-memory Transport
- Docker / Dockerfile / ENTRYPOINT
- ldflags -X (build-time dəyişənlər)
- Docker Compose / go mod vendor
- Prometheus / pull model / scrape_interval
- go-metrics / Counter / Timer / Registry / Percentile

## Praktik nəticə

- Consul: qeydiyyat startup-da, axtarış cache-lə; Client interfeysi mock üçün
- Raft: 3 node minimum; Apply yalnız leader-dən; in-memory transport test üçün
- Go + alpine = ~10MB imic; -ldflags -X ilə versiya build-də
- Compose: servis ADI = host adı; vendor ilə offline build
- Prometheus pull model: app yalnız /metrics verir; hədəflər yml-də
- Custom metric: GetOrRegister counter/timer; JSON/report və ya exporter

## Mənbə

Pages: 361-395 (Chapter 11, Go Programming Cookbook 2nd ed)
