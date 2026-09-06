# Fəsil 10 — Go-nun deploy edilməsi

## Bu fəsil nədən bəhs edir?

Go web tətbiqlərinin 4 deploy üsulu: standalone server (nohup/Upstart), Heroku (PaaS — Godep + Procfile), Google App Engine (məhdud sandbox + Cloud SQL), Docker (Dockerfile + Docker Machine). Hər üsulun müqayisəsi: kod dəyişikliyi, sistem işi, bakım, scalability, platform bağlanqlılığı.

## Əsas fikirlər

### 1. Cloud computing anlayışları
**NIST 3 model:**
- **IaaS** (Infrastructure-as-a-Service): basic compute/storage/network — AWS EC2, Google Compute Engine, Digital Ocean Droplets
- **PaaS** (Platform-as-a-Service): tətbiq deploy platforması — Heroku, AWS Elastic Beanstalk, Google App Engine
- **SaaS** (Software-as-a-Service): hazır servis — Heroku Postgres, AWS RDS, Google Cloud SQL

**Go-nun deploy üstünlüyü:** tək **statically linked binary** — heç bir kitabxana asılılığı yoxdur; amma template/static fayllar da lazım olur.

### 2. Standalone server deploy
**Addımlar:**
```bash
go build                 # ws-s binary yarat
# serverə kopyala
./ws-s                   # foreground — terminal bağlıdır!
nohup ./ws-s &           # background + HUP siqnalına laqeyd
```

**nohup problemi:** crash olsa xəbərdarlıq yoxdur; server restart olsa manual restart lazımdır.

**Həll — init daemon (Upstart/systemd):** init = kernel tərəfindən başladılan ilk proses. **Upstart job faylı** (`/etc/init/ws.conf`):
```bash
respawn
respawn limit 10 5
setuid sausheong
setgid sausheong
exec /go/src/github.com/sausheong/ws-s/ws-s
```

**Stanza izahı:**
- `respawn` → job ölərsə yenidən başlat
- `respawn limit 10 5` → maksimum 10 cəhd, 5 saniyə arayla; 10-dan sonra job "failed" sayılır
- `setuid/setgid` → hansı user/group ilə işə düşsün
- `exec` → icra olunacaq binary

**İdarəetmə:**
```bash
sudo start ws            # job başlat
ps -ef | grep ws          # PID yoxla (2011)
sudo kill -0 2011         # öldür
ps -ef | grep ws          # YENİ PID (2030) — Upstart avtomatik respawn etdi!
```

**Port qeydi:** production-da 8080 → 80 portuna keç, yaxud proxy/redirect ilə yönləndir.

### 3. Heroku deploy
**Tələblər:** asılılıq tərif faylı + **Procfile** (nə icra olunacaq). Deploy Git push ilə; **slug** yığılır; **dyno**-larda (izolyasiya olunmuş yüngül Unix konteynerlər) işləyir.

**Kod dəyişikliyi (yalnız 1):** port env-dən oxunmalıdır:
```go
func main() {
    server := http.Server{
        Addr: ":" + os.Getenv("PORT"),   // PORT env variable!
    }
    http.HandleFunc("/post/", handlePost)
    server.ListenAndServe()
}
```

**Godep ilə asılılıqlar:**
```bash
go get github.com/tools/godep
godep save
```
→ `Godeps/` qovluğu + `Godeps/_workspace` (asılılıqların source-u) + `Godeps.json`:
```json
{
  "ImportPath": "github.com/sausheong/ws-h",
  "GoVersion": "go1.4.2",
  "Deps": [
    {
      "ImportPath": "github.com/lib/pq",
      "Comment": "go1.0-cutoff-31-ga33d605",
      "Rev": "a33d6053e025943d5dc89dfa1f35fe5500618df7"
    }
  ]
}
```

**Procfile:**
```
web: ws-h
```
→ `web` prosesi = `ws-h` binary — build bitəndə icra olunacaq.

**Deploy:**
```bash
heroku login                  # toolbelt ilə (https://toolbelt.heroku.com)
heroku create ws-h            # app yarat (https://ws-h.herokuapp.com)
git push heroku master         # push → build → deploy
```

### 4. Google App Engine (GAE) deploy
**Üstünlüklər:** avtomatik scale + load balance; Google servis inteqrasiyası (Accounts auth, mail, logs, images).
**Məhdudiyyətlər (sandbox):** fayl sistemi **read-only**; request maksimum **60 saniyə**; birbaşa network/sistem çağırışı YOX → PostgreSQL-ə birbaşa qoşulmaq OLMAZ.

**Bunun əvəzinə:** Google Cloud SQL (MySQL əsaslı) + `cloudsql` paketi.

**Lazımi kod dəyişiklikləri (5):**
1. Package adı `main`-dən başqa adına (main yalnız standalone üçün; GAE-də paket kimi deploy olunur)
2. `main()` sil → handler qeydiyyatını `init()`-ə köçür:
```go
func init() {
    http.HandleFunc("/post/", handlePost)
}
```
(Server struct + ListenAndServe lazım DEYİL — GAE idarə edir.)
3. MySQL driveri (`_ "github.com/ziutek/mymysql/godrv"`) + cloudsql DSN:
```go
Db, err = sql.Open("mymysql", "cloudsql:<app ID>:<instance name>*<database name>/<user name>/<password>")
```
4. SQL sorğularını MySQL formatına: `$1` → `?`
```go
// PostgreSQL:  where id = $1
// MySQL:       where id = ?
```
5. app.yaml yarat:
```yaml
application: ws-g-1234
version: 1
runtime: go
api_version: go1
handlers:
- url: /.*
  script: _go_app
```

**CLI:**
```bash
goapp serve      # lokal test (http://localhost:8000 — admin)
goapp deploy     # Google serverlərinə push + compile + deploy
```
→ http://ws-g-1234.appspot.com

**Qeydlər:** Cloud SQL lokal dev-də dəstəklənmir (lokal MySQL istifadə et); IPv6 default pulsuz, IPv4 pullu (təyin et); "Follow App Engine App" seçimi mütləq.

### 5. Docker — anlayış
**Nədir:** Konteynerlərdə tətbiq build/ship/run üçün açıq platforma (dotCloud, 2013). Konteyner texnologiyası YENİ DEYİL (Unix, LXC 2008; Heroku dyno-ları da konteynerdir).

**VM vs Container:**
| | VM | Container |
|---|---|---|
| Virtualizasiya | Tam sistem + Guest OS | **OS səviyyəsində** — izolyasiya olunmuş user space |
| Resurs | Ağır | Çox yüngül |
| Start | Yavaş | Sürətli |

**Docker komponentləri:**
- **Docker client** → CLI
- **Docker daemon** → host-da konteynerləri idarə edən proses (Docker host)
- **Docker image** → read-only şablon; konteyner burdan işə düşür
- **Dockerfile** → image yaratmaq üçün təlimat faylı
- **Docker registry** → image deposu (Docker Hub — public/private)
- **Docker container** → işləyən nüsxə

**Quraşdırma:** Linux üçün `wget -qO- https://get.docker.com/ | sh`; yoxlama: `sudo docker run hello-world`.

### 6. Go web servisinin dockerləşdirilməsi
**Dockerfile (listing 10.7):**
```dockerfile
FROM golang
ADD . /go/src/github.com/sausheong/ws-d
WORKDIR /go/src/github.com/sausheong/ws-d
RUN go get github.com/lib/pq
RUN go install github.com/sausheong/ws-d
ENTRYPOINT /go/bin/ws-d
EXPOSE 8080
```

**Sub-kod izahı:**
- `FROM golang` → bazası: Go quraşdırılmış Debian image, GOPATH=/go
- `ADD . /go/src/...` → lokal kodu konteynerin workspace-inə kopyala
- `WORKDIR` → iş qovluğu təyin et
- `RUN go get ...` → asılılıqları çək (build zamanı)
- `RUN go install ...` → binary-nı /go/bin-ə build et
- `ENTRYPOINT /go/bin/ws-d` → konteyner başlayanda icra olunacaq əmr
- `EXPOSE 8080` → portu digər KONTEYNERLƏRƏ açır (public-a yox!)

**Build və run:**
```bash
docker build -t ws-d .          # image yarat
docker images                   # yoxla: ws-d latest 65e8437fce6b 534.7 MB
docker run --publish 80:8080 --name simple_web_service --rm ws-d
docker ps                       # aktiv konteynerlər
```
- `--publish 80:8080` → host-un 80 portunu konteynerin 8080-ə map et (public giriş!)
- `--rm` → konteyner çıxanda silinsin (yoxdursa — restart mümkün qalır)

**Test:**
```bash
curl -i -X POST -H "Content-Type: application/json" \
  -d '{"content":"My first post","author":"Sau Sheong"}' http://127.0.0.1/post/
curl -i -X GET http://127.0.0.1/post/1
```

### 7. Docker Machine — bulutta host
**Nədir:** CLI — lokal və ya cloud-da Docker host (yəni daemon işləyən VM) yaradır. Dəstəklənənlər: AWS, Digital Ocean, Google, Azure, Rackspace, Softlayer, VMWare, OpenStack, Hyper-V...

**Quraşdırma (Linux):**
```bash
curl -L https://github.com/docker/machine/releases/download/v0.3.0/docker-machine_linux-amd64 /usr/local/bin/docker-machine
chmod +x /usr/local/bin/docker-machine
```

**Digital Ocean-da host (VPS provayder):**
1. Hesab + Applications & API → **Generate New Token** (Write seçimi ilə) — token bir də göstərilir, saxla!
2. Host yarat:
```bash
docker-machine create --driver digitalocean --digitalocean-access-token <token> wsd
```
3. Client-i remote host-a yönləndir:
```bash
docker-machine env wsd
# export DOCKER_HOST="tcp://104.236.0.57:2376" ...
eval "$(docker-machine env wsd)"
```
4. Remote-da image build + run (eyni əmrlər!):
```bash
docker build -t ws-d .
docker run --publish 80:8080 --name simple_web_service --rm ws-d
```
5. Test: `curl -i -X GET http://104.236.0.57/post/1`

**Güc:** bir də lokalsa — **eyni proses** istənilən cloud-da təkrarlana bilər.

### 8. Deploy üsulları müqayisəsi (cədvəl 10.1)
| | Standalone | Heroku | GAE | Docker |
|---|---|---|---|---|
| Tip | Public/private | Public | Public | Public/private |
| Kod dəyişikliyi | Yox | Aşağı | Orta | Yox |
| Sistem işi | Yüksək | Yox | Yox | Orta |
| Bakım | Yüksək | Yox | Yox | Orta |
| Deploy asanlığı | Aşağı | Yüksək | Orta | Aşağı |
| Platform dəstəyi | Yox | Aşağı | Yüksək | Aşağı |
| Platform bağlanqlılığı | Yox | Aşağı | Yüksək | Aşağı |
| Scalability | Yox | Orta | Yüksək | Yüksək |
| Qeyd | Hər şeyi özün edirsin | Liberal PaaS | Restrictive PaaS (Google qaydaları) | Yüksələn texnologiya |

## Əsas terminlər
- IaaS / PaaS / SaaS (NIST modelləri)
- VPS (Virtual Private Server)
- nohup (HUP siqnalını laqeyd etmə)
- Init Daemon / Upstart / systemd
- Respawn (avtomatik yenidən başlatma)
- Stanza (Upstart komanda bloku)
- Slug / Dyno (Heroku build nüsxəsi / konteyner)
- Procfile (Heroku proses tərifi)
- Godep / Godeps.json (asılılıq idarəsi)
- Sandbox (GAE məhdudiyyət mühiti)
- Google Cloud SQL (MySQL əsaslı)
- app.yaml (GAE konfiq)
- goapp serve/deploy (GAE CLI)
- Placeholder fərqi: `$1` (PostgreSQL) vs `?` (MySQL)
- Docker: client / daemon / host
- Image / Container (şablon / işləyən nüsxə)
- Dockerfile (FROM/ADD/WORKDIR/RUN/ENTRYPOINT/EXPOSE)
- Registry / Docker Hub (image deposu)
- --publish / --rm (port xəritəsi / avtomatik təmizləmə)
- Docker Machine (host yaratma CLI)
- eval docker-machine env (client yönləndirmə)
- Personal Access Token (API autentifikasiyası)

## Praktik nəticə
- Standalone üçün init daemon (systemd/Upstart) istifadə et — nohup crash-dan sonra xilas olmur.
- Heroku: `os.Getenv("PORT")` + Godep + Procfile → `git push` — ən az kod dəyişikliyi ilə ən sadə deploy.
- GAE güclü scale versə də sandbox ciddi bağlanqlılıq deməkdir: main→init, `$1`→`?`, PostgreSQL→Cloud SQL.
- Docker-da Go binary statically linked olduğundan scratch/alpine bazalı ultra-kiçik image-lər mümkündür (kitabın golang-base üsulundan daha optimal — müəllim qeydi).
- Docker Machine ilə lokal test → cloud deploy eyni əmrlər — bu Docker-ın ən böyük dəyəridir.
- Port: konteyner daxili 8080 → host 80 (`--publish`); `EXPOSE` yalnız konteynerlərarası sənədləşdirmədir.

## Mənbə
Pages: 277-301 (PDF), book pages 256-280
