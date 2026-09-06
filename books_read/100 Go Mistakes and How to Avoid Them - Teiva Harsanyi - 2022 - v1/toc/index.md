# 100 Go Mistakes and How to Avoid Them — Mündəricat (TOC)

**Müəllif:** Teiva Harsanyi
**Nəşriyyat:** Manning Publications, 2022 · **ISBN:** 9781617299599
**Səhifə:** 364 + ön söz (385 PDF səh.) · **12 chapter · 100 səhv (mistake)**

> PDF offset: book page N → PDF page N+20
> Kitab 100 ümumi səhvi 6 kateqoriyada qruplaşdırır: Bugs, Needless complexity, Weaker readability, Suboptimal organization, Lack of API convenience, Under-optimized code, Lack of productivity.

## Chapter-lər

| Ch | Başlıq | Səhvlər | Book səh. | PDF səh. |
|----|--------|---------|-----------|----------|
| 1 | Go: Simple to learn but hard to master | — | 1-6 | 21-26 |
| 2 | Code and project organization | #1-#16 | 7-55 | 27-75 |
| 3 | Data types | #17-#29 | 56-94 | 76-114 |
| 4 | Control structures | #30-#35 | 95-112 | 115-132 |
| 5 | Strings | #36-#41 | 113-125 | 133-145 |
| 6 | Functions and methods | #42-#47 | 126-142 | 146-162 |
| 7 | Error management | #48-#54 | 143-161 | 163-181 |
| 8 | Concurrency: Foundations | #55-#60 | 162-192 | 182-212 |
| 9 | Concurrency: Practice | #61-#74 | 193-233 | 213-253 |
| 10 | The standard library | #75-#81 | 234-261 | 254-281 |
| 11 | Testing | #82-#90 | 262-298 | 282-318 |
| 12 | Optimizations | #91-#100 | 299-354 | 319-374 |

## Section-lər (səhv siyahısı)

### Ch 2 — Code and project organization (#1-#16)
#1 Unintended variable shadowing (dəyişən kölgələnməsi) · #2 Unnecessary nested code · #3 Misusing init functions · #4 Overusing getters and setters · #5 Interface pollution · #6 Interface on the producer side · #7 Returning interfaces · #8 any says nothing · #9 Being confused about when to use generics · #10 Not being aware of the possible problems with type embedding · #11 Not using the functional options pattern · #12 Project misorganization · #13 Creating utility packages · #14 Ignoring package name collisions · #15 Missing code documentation · #16 Not using linters

### Ch 3 — Data types (#17-#29)
#17 Creating confusion with octal literals · #18 Neglecting integer overflows · #19 Not understanding floating points · #20 Not understanding slice length and capacity · #21 Inefficient slice initialization · #22 Being confused about nil vs. empty slices · #23 Not properly checking if a slice is empty · #24 Not making slice copies correctly · #25 Unexpected side effects using slice append · #26 Slices and memory leaks · #27 Inefficient map initialization · #28 Maps and memory leaks · #29 Comparing values incorrectly

### Ch 4 — Control structures (#30-#35)
#30 Ignoring the fact that elements are copied in range loops · #31 Ignoring how arguments are evaluated in range loops · #32 Ignoring the impact of using pointer elements in range loops · #33 Making wrong assumptions during map iterations · #34 Ignoring how the break statement works · #35 Using defer inside a loop

### Ch 5 — Strings (#36-#41)
#36 Not understanding the concept of a rune · #37 Inaccurate string iteration · #38 Misusing trim functions · #39 Under-optimized string concatenation · #40 Useless string conversions · #41 Substrings and memory leaks

### Ch 6 — Functions and methods (#42-#47)
#42 Not knowing which type of receiver to use · #43 Never using named result parameters · #44 Unintended side effects with named result parameters · #45 Returning a nil receiver · #46 Using a filename as a function input · #47 Ignoring how defer arguments and receivers are evaluated

### Ch 7 — Error management (#48-#54)
#48 Panicking · #49 Ignoring when to wrap an error · #50 Checking an error type inaccurately · #51 Checking an error value inaccurately · #52 Handling an error twice · #53 Not handling an error · #54 Not handling defer errors

### Ch 8 — Concurrency: Foundations (#55-#60)
#55 Mixing up concurrency and parallelism · #56 Thinking concurrency is always faster · #57 Being puzzled about when to use channels or mutexes · #58 Not understanding race problems · #59 Not understanding the concurrency impacts of a workload type · #60 Misunderstanding Go contexts

### Ch 9 — Concurrency: Practice (#61-#74)
#61 Propagating an inappropriate context · #62 Starting a goroutine without knowing when to stop it · #63 Not being careful with goroutines and loop variables · #64 Expecting deterministic behavior using select and channels · #65 Not using notification channels · #66 Not using nil channels · #67 Being puzzled about channel size · #68 Forgetting about possible side effects with string formatting (etcd data race, deadlock) · #69 Creating data races with append · #70 Using mutexes inaccurately with slices and maps · #71 Misusing sync.WaitGroup · #72 Forgetting about sync.Cond · #73 Not using errgroup · #74 Copying a sync type

### Ch 10 — The standard library (#75-#81)
#75 Providing a wrong time duration · #76 time.After and memory leaks · #77 Common JSON-handling mistakes (type embedding, monotonic clock, map of any) · #78 Common SQL mistakes (sql.Open, pooling, prepared statements, null values, row iteration) · #79 Not closing transient resources (HTTP body, sql.Rows, os.File) · #80 Forgetting the return statement after replying to an HTTP request · #81 Using the default HTTP client and server

### Ch 11 — Testing (#82-#90)
#82 Not categorizing tests (build tags, env vars, short mode) · #83 Not enabling the -race flag · #84 Not using test execution modes (parallel, -shuffle) · #85 Not using table-driven tests · #86 Sleeping in unit tests · #87 Not dealing with the time API efficiently · #88 Not using testing utility packages (httptest, iotest) · #89 Writing inaccurate benchmarks · #90 Not exploring all the Go testing features

### Ch 12 — Optimizations (#91-#100)
#91 Not understanding CPU caches (cache line, slice of structs vs struct of slices, predictability, cache placement) · #92 Writing concurrent code that leads to false sharing · #93 Not taking into account instruction-level parallelism · #94 Not being aware of data alignment · #95 Not understanding stack vs. heap (escape analysis) · #96 Not knowing how to reduce allocations (API changes, compiler optimizations, sync.Pool) · #97 Not relying on inlining · #98 Not using Go diagnostics tooling (profiling, execution tracer) · #99 Not understanding how the GC works · #100 Not understanding the impacts of running Go in Docker and Kubernetes

## Kitab strukturu qeydi

Kitab tək layihə üzərində deyil — hər səhv müstəqil kontekstdə izah olunur (scenario + yanlış kod + düzəliş + müzakirə). project-state/ lazım deyil.
