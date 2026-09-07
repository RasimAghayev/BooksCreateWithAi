# Everyday Golang — Mündəricat

**Müəllif:** Alex Ellis (OpenFaaS qurucusu) | **İl:** 2021 | **Səviyyə:** L2 (Elementary)

Production-dan götürülmüş gündəlik Go pattern-ləri — 18 fəsil, tam praktik:
quraşdırmadan OpenFaaS funksiyalarınadək. Bilik faylları 4 qrupa birləşdirilib
(kitabın fəsil strukturu qorunmaqla).

## Bilik faylları

| Fayl | Əhatə | Orijinal fəsillər |
|---|---|---|
| [01-05-foundations-json-cli](../chapters/01-05-foundations-json-cli/index.md) | Go modules, vendoring, HMAC, flags, çoxpaketli layihə, cross-compile, HTTP+JSON, CLI prinsipləri | 1-5 |
| [06-08-testing-concurrency-database](../chapters/06-08-testing-concurrency-database/index.md) | Unit-testlər (test-table, coverage, parallel), httptest, goroutine/WaitGroup/errgroup/singleflight, Mutex/channels/context/worker pool, PostgreSQL | 6-8 |
| [09-17-config-http-metrics-release](../chapters/09-17-config-http-metrics-release/index.md) | YAML+mergo, ldflags version, embed, templates, HTTP server/mux/middleware, Prometheus RED+Collector, GitHub Actions, multi-arch Docker | 9-17 |
| [18-openfaas-and-go](../chapters/18-openfaas-and-go/index.md) | OpenFaaS, funksiya şablonları, faas-cli, secrets, DB init(), REST routing | 18 |

## Oxu ardıcıllığı

Ardıcıl 1→4: hər fəsil əvvəlkinin kodunu davam etdirir (eyni nümunələr
təkmilləşir: hash-browns → get-json → isolate-astros → todo → todo-fn).
