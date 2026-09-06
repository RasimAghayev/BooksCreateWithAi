# 100 Go Mistakes and How to Avoid Them — Xarimi References

## Kitabın öz resursları
- Müəllifin blogu (100 mistakes mənbəyi): https://teivah.medium.com
- "The Top 10 Most Common Mistakes I've Seen in Go Projects" (2019): teivah.medium.com

## Go rəsmi sənədləri
- Go FAQ: https://go.dev/doc/faq
- Effective Go: https://go.dev/doc/effective_go
- Go memory model: https://go.dev/ref/mem
- Go spec (select random selection): https://go.dev/ref/spec
- Diagnostics: https://go.dev/doc/diagnostics
- Profiling (pprof): https://go.dev/blog/pprof
- Race detector: https://go.dev/doc/articles/race_detector
- context paketi: https://pkg.go.dev/context
- sync paketi: https://pkg.go.dev/sync
- database/sql: https://pkg.go.dev/database/sql
- net/http (Response error semantikası): https://pkg.go.dev/net/http
- httptest: https://pkg.go.dev/net/http/httptest
- iotest: https://pkg.go.dev/testing/iotest
- go command (vet, test): https://golang.org/cmd/vet/, https://golang.org/cmd/go/
- Go install: https://golang.org/doc/install

## Go blogları və məqalələr
- "Defer, Panic and Recover": https://go.dev/blog/defer-panic-and-recover
- "Error handling and Go": https://go.dev/blog/error-handling-and-go
- Go map iteration randomization: http://mng.bz/M2JW
- Go scheduler (design doc): http://mng.bz/N611, http://mng.bz/lxY8
- OS thread size: http://mng.bz/DgMw
- CPU cache latencies: http://mng.bz/o29v
- Go 2021 developer survey: https://go.dev/blog/survey2021-results

## Alətlər və kitabxanalar
- project-layout: https://github.com/golang-standards/project-layout
- golangci-lint: https://github.com/golangci/golangci-lint
- errcheck: https://github.com/kisielk/errcheck
- gocyclo: https://github.com/fzipp/gocyclo
- goconst: https://github.com/jgautheron/goconst
- goimports: https://godoc.org/golang.org/x/tools/cmd/goimports
- shadow analyzer: golang.org/x/tools/go/analysis/passes/shadow
- testify: https://github.com/stretchr/testify
- go-cmp: https://github.com/google/go-cmp
- GoDS (sıralı strukturlar): https://github.com/emirpasic/gods
- Pebble (capacity+append nümunəsi): https://github.com/cockroachdb/pebble
- errgroup: https://pkg.go.dev/golang.org/x/sync/errgroup
- automaxprocs: https://github.com/uber-go/automaxprocs
- go-sqlmock: https://github.com/DATA-DOG/go-sqlmock
- benchstat: https://golang.org/x/perf
- perflock: https://github.com/aclements/perflock

## Real-world case-lər
- etcd data race fix (PR 7816): https://github.com/etcd-io/etcd/pull/7816
- Echo 405 PR (channel müzakirəsi üçün): https://github.com/labstack/echo/pull/205
- go-sql-driver/mysql: https://github.com/go-sql-driver/mysql

## Go issues
- #14813 (benchmark/compiler optimization): https://github.com/golang/go/issues/14813
- #6104 (threadcreate profili pozulub): https://github.com/golang/go/issues/6104
- #33803 (CFS-aware GOMAXPROCS): https://github.com/golang/go/issues/33803

## Elmi istinadlar
- T. Tu, X. Liu et al., "Understanding Real-World Concurrency Bugs in Go" (ASPLOS 2019)
- J. S. Moser, H. S. Schroder et al., "Mind Your Errors" (Psychological Science, 2011)
- J. Metcalfe, "Learning from Errors" (Annual Review of Psychology, 2017)
- Synopsys, "The Cost of Poor Software Quality in the US" (2020)
- R. C. Martin, "Clean Code" (2008) — 10:1 oxuma/yazma nisbəti
- LMAX Disruptor white paper (Martin Thompson et al.): https://lmax-exchange.github.io/disruptor/files/Disruptor-1.0.pdf
- Ariane 5 disaster: https://www.bugsnag.com/blog/bug-day-ariane-5-disaster
- Wes Dyer sitatı: "Make it correct, make it clear, make it concise, make it fast, in that order."
