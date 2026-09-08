# System Programming Essentials with Go — Xarici References

## Kitabın repo-su
- Kitab kodu: https://github.com/PacktPublishing/System-Programming-Essentials-with-Go
- Müəllif: Alex Rios (github.com/alexrios — timer/v2, endpoints kitabxanaları)
- Kitab nümunə konfransları: ch11/otel docker-compose, ch13 cache, appendix-a/usb, appendix-a/bt

## Go rəsmi
- Go blog: https://blog.golang.org/
- Release notes: https://go.dev/doc/devel/release
- slog (Go 1.21): https://pkg.go.dev/log/slog
- runtime/trace: https://pkg.go.dev/runtime/trace
- go command (modules/work): https://go.dev/ref/mod · https://go.dev/doc/tutorial/workspaces
- MVS alqoritmi: https://go.dev/ref/mod#minimal-version-selection
- GOPRIVATE: https://go.dev/doc/faq#avoid-gopath · go.dev/reference/cmd/go/internal/modinit

## Sistem proqramlaşdırma
- APUE (kitab): W. Richard Stevens — Advanced Programming in the UNIX Environment
- Unix Network Programming (kitab): W. Richard Stevens
- Linux System Programming (kitab): Robert Love
- strace: https://strace.io/
- inotify: man7.org/linux/man-pages/man7/inotify.7.html
- Linux kernel /proc/mounts: man7.org/linux/man-pages/man5/proc.5.html (mounts, fstab format)
- Mkfifo: man7.org/linux/man-pages/man3/mkfifo.3.html
- UNIX domain sockets: man7.org/linux/man-pages/man7/unix.7.html
- lsof: https://github.com/lsof-org/lsof

## Yaddaş + performans
- Go GC guide: https://go.dev/doc/gc-guide
- GOGC/GOMEMLIMIT: go.dev/doc/gc-guide#GOGC · go.dev/doc/gc-guide#memory_limit
- Arena design doc: github.com/golang/proposal/blob/master/design/51317-arena-type.md
- Memory ballast (Dropbox): github.com/dropbox/godropops (blog: "Tuning Go's HTTP" post script)
- Escape analysis: go.dev/blog/cmd/go (compiler flags) · utcc.utoronto.ca escape yazıları
- pprof: https://github.com/google/pprof
- Benchmark: go.dev/doc/tutorial/add-a-test · pkg.go.dev/testing

## Şəbəkə
- net paketi: https://pkg.go.dev/net
- crypto/tls: https://pkg.go.dev/crypto/tls
- OpenSSL: https://www.openssl.org/
- gobwas/ws: https://github.com/gobwas/ws
- SACK: RFC 2018 (Selective Acknowledgment) · Go-Back-N: rfc-editor.org

## Telemetriya
- Prometheus: https://prometheus.io · metric types: prometheus.io/docs/concepts/metric_types/
- client_golang: https://github.com/prometheus/client_golang
- zap: https://github.com/uber-go/zap
- OpenTelemetry: https://opentelemetry.io · status: opentelemetry.io/status/
- OTel Go SDK: go.opentelemetry.io/otel · exporters: otlptracehttp
- otelhttp (contrib): github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/net/http/otelhttp
- Grafana "Where and why we use Go": grafana.com/blog/2015/08/21/where-and-why-we-use-go/

## Paylama
- Staticcheck: https://staticcheck.dev/ · action: github.com/dominikh/staticcheck-action
- GitHub Actions cache: docs.github.com/en/actions/using-workflows/caching-dependencies
- GoReleaser: https://goreleaser.com · action: github.com/goreleaser/goreleaser-action
- SemVer: https://semver.org/
- Go module proxy: go.dev/ref/mod

## Effektiv praktikalar
- sync.Pool: pkg.go.dev/sync#Pool
- sync.OnceValue/OnceValues (Go 1.21): go.dev/ref/spec (Go 1.21 release notes)
- singleflight (x/sync): pkg.go.dev/golang.org/x/sync/singleflight
- x/exp/mmap: pkg.go.dev/golang.org/x/exp/mmap
- mmap(2): man7.org/linux/man-pages/man2/mmap.2.html · msync(2)

## Capstone (distributed cache)
- hashicorp/golang-lru: https://github.com/hashicorp/golang-lru
- patrickmn/go-cache: https://github.com/patrickmn/go-cache
- dgraph-io/ristretto: https://github.com/dgraph-io/ristretto
- Consistent hashing (orijinal məqalə): David Karger et al., "Consistent Hashing and Random Trees"
- groupcache (singleflight mənşəi): https://github.com/golang/groupcache
- LIRS: www.cse.ohio-state.edu/~zhang/lirs.html · ARC: www.ibm.com ARC məqalələri

## Hardware avtomatlaşdırma
- freedesktop.org: https://www.freedesktop.org
- D-Bus: https://www.freedesktop.org/wiki/Software/dbus/
- godbus/dbus: https://github.com/godbus/dbus
- UDisks2: https://www.freedesktop.org/wiki/Software/udisks/
- Notification spec: specifications.freedesktop.org/notification-spec/notification-spec-latest.html
- Icon naming spec: specifications.freedesktop.org/icon-naming-spec/icon-naming-spec-latest.html
- go-bluetooth (muka): https://github.com/muka/go-bluetooth
- BlueZ: git.kernel.org/pub/scm/bluetooth/bluez.git
- XDG screensaver: freedesktop.org wiki · Wayland: wayland.freedesktop.org

## Real dünya (Ch 15)
- Go at Dropbox: youtube.com/watch?v=JOx9enktnUM
- HashiCorp Go seçimi (Nic Jackson): youtu.be/qlwp0mHFLHU
- SoundCloud Ruby→Go: SoundCloud engineering blog (developers.soundcloud.com arxiv)
- Docker/libcontainer: github.com/moby/moby (containerd/libcontainer tarixi)
- CNCF landscape: https://www.cncf.io · cncf.io/projects

## Şəbəkə alətləri (ch10 tls)
- mkcert (local CA): https://github.com/FiloSottile/mkcert
- Wireshark: https://www.wireshark.org
