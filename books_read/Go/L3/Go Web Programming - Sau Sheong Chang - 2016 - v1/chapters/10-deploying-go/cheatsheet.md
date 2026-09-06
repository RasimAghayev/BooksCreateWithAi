# Chapter 10 — Deploying Go Cheatsheet

## `go build` — static binary yaratmaq

**Nə edir:** Tək icra edilə bilən fayl yaradır. Heç bir xarici kitabxana asılılığı yoxdur (CGO olmadıqda).

```bash
go build
# → ws-s binary yaranır (statik link)
scp ws-s user@server:/path/
```

**Mənbə:** Chapter 10, page 261

---

## `nohup ./ws-s &` — prosesi fonda saxlamaq

**Nə edir:** HUP (hangup) signal-ı bloklayır, logout olduqda da proses davam edir.

```bash
nohup ./ws-s &
```

**Çatışmazlıqları:** Crash olarsa restart olunmur — Upstart/systemd lazımdır.

**Mənbə:** Chapter 10, page 261

---

## Upstart job konfiqurasiyası (`/etc/init/ws.conf`)

**Nə edir:** Init daemon vasitəsilə prosesi idarə edir — auto-start, auto-respawn.

```bash
# /etc/init/ws.conf
respawn
respawn limit 10 5
setuid sausheong
setgid sausheong
exec /go/src/github.com/sausheong/ws-s/ws-s

sudo start ws     # başlat
sudo stop ws      # dayandır
sudo kill -0 PID  # test (proses ölsə Upstart respawn edir)
```

**Parametrlər:**
- `respawn` — crash olduqda yenidən başlat
- `respawn limit 10 5` — maksimum 10 dəfə, 5 san interval
- `setuid`/`setgid` — hansı user ilə işləsin
- `exec` — icra ediləcək əmr

**Mənbə:** Chapter 10, page 262

---

## Heroku: `os.Getenv("PORT")` ilə port təyini

**Nə edir:** Heroku porta nəzarət etmir — kodunuz environment variable-dan oxumalıdır.

```go
server := http.Server{
    Addr: ":" + os.Getenv("PORT"),
}
```

**Mənbə:** Chapter 10, page 264

---

## `godep save` — Go dependency idarəsi

**Nə edir:** Bütün asılılıqların source code-unu `Godeps/_workspace`-ə köçürür, `Godeps.json` yaradır. Heroku üçün vacibdir.

```bash
go get github.com/tools/godep
godep save
# → Godeps/ + Godeps.json yaranır
```

**Mənbə:** Chapter 10, page 265

---

## `Procfile` — Heroku process təyini

**Nə edir:** Hansı executable-ın işə düşəcəyini bildirir.

```bash
# Procfile
web: ws-h
```

**Mənbə:** Chapter 10, page 265

---

## `heroku create` + `git push heroku master`

**Nə edir:** Heroku tətbiqi yaradır və Git push ilə deploy edir (avtomatik build + dyno run).

```bash
heroku login
heroku create ws-h
git push heroku master
# → Build + deploy avtomatik
# Tətbiq: https://ws-h.herokuapp.com
```

**Mənbə:** Chapter 10, page 266

---

## GAE: `package` adı dəyişikliyi və `init()` funksiyası

**Nə edir:** GAE özü process-i idarə edir — `main()` yox, `init()` lazımdır, package adı `main` olmamalıdır.

```go
// main.go deyil — məs: package app
package app

func init() {
    http.HandleFunc("/post/", handlePost)
}
// ListenAndServe() — yoxdur, GAE idarə edir
```

**Mənbə:** Chapter 10, page 269

---

## GAE: MySQL driver + Cloud SQL

**Nə edir:** PostgreSQL əvəzinə MySQL (Cloud SQL — managed, IPv4/IPv6).

```go
import (
    "database/sql"
    _ "github.com/ziutek/mymysql/godrv"
)

func init() {
    Db, _ = sql.Open("mymysql",
        "cloudsql:<app ID>:<instance>*<db>/<user>/<password>")
}

// SQL: $1, $2 → ? (MySQL placeholder)
Db.QueryRow("select id, content, author from posts where id = ?", id)
```

**Mənbə:** Chapter 10, page 270

---

## `app.yaml` — GAE konfiqurasiya faylı

**Nə edir:** GAE tətbiq parametrləri: app adı, runtime versiyası, URL routing.

```yaml
application: ws-g-1234
version: 1
runtime: go
api_version: go1
handlers:
  - url: /.*
    script: _go_app
```

**Mənbə:** Chapter 10, page 270

---

## `goapp serve` / `goapp deploy` — GAE SDK əmrləri

**Nə edir:** `goapp serve` lokal run + admin UI (localhost:8000); `goapp deploy` Google serverlərinə push.

```bash
goapp serve     # lokal test (Cloud SQL lokal işləmir)
goapp deploy    # production deploy
# → http://ws-g-1234.appspot.com
```

**Mənbə:** Chapter 10, page 270

---

## `Dockerfile` — Docker image təyini

**Nə edir:** Image build etmək üçün instructionlar. Hər sətir image layer yaradır.

```dockerfile
FROM golang                                     # base image (Debian + Go)
ADD . /go/src/github.com/sausheong/ws-d        # kodu köçür
WORKDIR /go/src/github.com/sausheong/ws-d      # iş qovluğu
RUN go get github.com/lib/pq                    # dependency yüklə
RUN go install github.com/sausheong/ws-d        # binary quraşdır
ENTRYPOINT /go/bin/ws-d                         # container start olanda run
EXPOSE 8080                                     # port 8080 digər container-lərə açıq
```

**Mənbə:** Chapter 10, page 275

---

## `docker build` / `docker run`

**Nə edir:** Dockerfile-dən image yaratmaq və image-dən container işə salmaq.

```bash
# Build
docker build -t ws-d .

# Lokal run
docker run --publish 80:8080 --name simple_web_service --rm ws-d
# --publish 80:8080: host port 80 → container port 8080
# --name: container adı
# --rm: dayandıqda avtomatik sil

# Status
docker ps        # aktiv container-lər
docker images    # mövcud image-lər
```

**Mənbə:** Chapter 10, page 275

---

## `docker-machine create` — cloud-da Docker host yaratmaq

**Nə edir:** AWS, DigitalOcean, Google Compute Engine və s. cloud-larda remote Docker host yaradır.

```bash
# DigitalOcean nümunəsi
docker-machine create --driver digitalocean \
    --digitalocean-access-token <token> wsd

# Remote host-a qoşul
docker-machine env wsd
eval "$(docker-machine env wsd)"

# Eyni build/run remote-da
docker build -t ws-d .
docker run --publish 80:8080 --name simple_web_service --rm ws-d

# Test
curl http://<remote-ip>/post/1
```

**Mənbə:** Chapter 10, page 277
