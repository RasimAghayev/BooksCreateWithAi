# Production Go — Mündəricat (TOC)

**Müəlliflər:** Herman Schaaf, Shawn Smith
**Nəşriyyat:** Leanpub, 2018-12-20 · **Səhifə:** 141 + lisenziyalar (147 PDF səh.)
**Struktur:** 14 bölmə — basics-dən production-unvaryinglarına (CI/deploy/monitoring)

> PDF offset: book page N → PDF page N+6

## Bölmələr

| # | Başlıq | Book səh. | PDF səh. |
|---|--------|-----------|----------|
| 1 | Introduction | i | 6 |
| 2 | Getting Started (Installing Go, Editors, Linters) | 1-6 | 7-12 |
| 3 | Basics (sintaksis icmalı) | 7-51 | 13-57 |
| 4 | Style and Error Handling | 46-51 | 52-57 |
| 5 | Strings | 52-62 | 58-68 |
| 6 | Supporting Unicode | 62-67 | 68-73 |
| 7 | Concurrency | 68-81 | 74-87 |
| 8 | Testing | 82-100 | 88-106 |
| 9 | Benchmarks | 101-112 | 107-118 |
| 10 | Tooling (godoc, Go Guru, Race Detector, Report Card) | 113-127 | 119-133 |
| 11 | Security (CSRF, HSTS, CSP) | 128-131 | 134-137 |
| 12 | Continuous Integration (Travis, Drone) | 132-134 | 138-140 |
| 13 | Deployment | 135-136 | 141-142 |
| 14 | Monitoring / Optimization / Gotchas / Further Reading | 136-139 | 142-145 |

## Bölmə daxili strukturlar

### Ch 3 — Basics
Program Structure / Variables & Constants / Basic Data Types (underflow!) / Structs / Operators / Conditional / Arrays / Slices / Maps / Loops / Functions / Exported Names / Pointers / Goroutines / Channels / Interfaces / Error Handling / Reading Input / Writing Output

### Ch 4 — Style & Error Handling
İdiomatik üslub, error handling fəlsəfəsi (error = dəyər)

### Ch 5 — Strings
Append (strings.Builder) / Split / Count & Find / Advanced string functions / Range

### Ch 6 — Unicode
Encoding tarixçəsi / Strings are byte slices / Printing / Runes / RTL dillər

### Ch 7 — Concurrency
sync.WaitGroup / Channels / Goroutines in web handlers / **Pollers** / **Race conditions** (detector ilə)

### Ch 8 — Testing
Table-driven tests / HTTP handler testləri / Mocking / Coverage / Examples

### Ch 9 — Benchmarks
Sadə benchmark / Müqayisə / Timer reset / Memory allocations (B/op, allocs/op) / **Modulo vs Bitwise-and**

### Ch 10 — Tooling
Godoc / Go Guru / **Race Detector** (DATA RACE analiz) / Go Report Card

### Ch 11 — Security
CSRF / HSTS / CSP

### Ch 12 — CI
Travis CI / Drone / General tips

### Ch 14 — Monitoring & Gotchas
Prometheus (work-in-progress) / Optimization (TODO) / **Common Gotchas: Nil interface**

## Layihə davamlılığı

```json
{
  "project_continuity": false,
  "note": "Kitab konsept-bölməlidir; hər hissə müstəqil nümunələrlə. Son 3 bölmə (Monitoring, Optimization) yarımçıqdır — müəlliflərin 'work in progress' qeydi var. project-state/ lazım deyil."
}
```
