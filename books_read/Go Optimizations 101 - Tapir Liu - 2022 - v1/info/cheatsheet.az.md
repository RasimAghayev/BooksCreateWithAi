# Go Optimizations 101 — Cheat Sheet (Azərbaycanca)

> **Kitab:** Go Optimizations 101 — Tapir Liu, go101.org, 2022 (v1.19 istinad) · 🧠 Expert (5/5)
> Ən vacib optimizasiya qaydaları — bir baxışda. Benchmark-sız optimallaşdırma YOXDUR.

---

## 1. Ölçü və kopya (Ch 2)
| Qayda | Nümunə |
|---|---|
| Sahə sırası padding-i azaldır | `{int8,int64,int16}`=24B → `{int8,int16,int64}`=16B |
| ≤4 word struct kopyası xüsusi sürətlidir | 9→10 sahə: 2-7× yavaşlama — hot struct-u kiçik saxla |
| `for _, v := range` böyük elementi İKİ dəfə kopyalayır | `for i := range s` + `s[i]`, ya `&s[i]` |
| Böyük array-i value ötürmə = tam kopya | pointer `*[N]T` və ya slice `[]T` |

## 2. Allocation azaltma (Ch 3)
| Qayda | Təsir |
|---|---|
| `make([]T, 0, n)` ilə pre-allocate | 4 alloc → 1 alloc, 2.4× sürət |
| In-place filtrləmə `data[:k]` qaytar | 8× sürət, 0 allocation |
| 100 obyekti bir `make([]T,100)`-da + `&books[i]` | 101→2 alloc |
| Pool (sync.Pool / custom) spawn/destroy dövrələrində | GC təzyiqi ↓ |
| Size class israfı: 33B→48B blok; 32769B→40960B | strukturu 32B-a sığdır |

## 3. Stack-də qal / escape (Ch 4)
| Qayda | Nümunə |
|---|---|
| `fmt.Println` reflect → heap; builtin `println` → stack | hot path-də fmt-dən qaç |
| Interface metoda pointer ötürmə → escape; konkret tip → stack | `t.M(&x)` OK, `i.M(&y)` escape |
| `make([]T, runtimeDəyişən)` həmişə heap | sabit N stack şansı verir |
| `var a [10<<20]byte; s := a[:]` | 64KB həddini 10MB-a qaldırır |
| Deep rekursiya: dummy 64MB frame funksiyası ilə stack-i əvvəlcədən böyüt | 42ms→4.7ms |
| Yoxlama: `go run -gcflags=-m` | "moved to heap" izlə |

## 4. GC (Ch 5)
| Alət | Nə edir |
|---|---|
| `GODEBUG=gctrace=1` | GC% və heap hədəflərini göstər |
| `GOGC=1000` və ya `SetGCPercent` | cycle-ları seyrlət (heap bahasına) |
| Memory ballast `make([]byte, 150<<20)` + `KeepAlive` | v1.18-ə qədər; virtual RAM pulsuz |
| `GOMEMLIMIT` (v1.19+) | müasir ballast alternativi, container OOM qoruması |
| Uzunömürlü kiçik substring/subslice → kopyala | böyük blok 1 bayta görə yığılmır |
| Qısaömürlü zibil: string concat, []byte↔string, boxing | hot path-də azalt |

## 5. Loop-lar (Ch 6-7)
| Qayda | Təsir |
|---|---|
| Loop-da pointer/struct-sahə yığma: `n := *sum; loop; *sum = n` | 5× sürət (aliasing yoxdursa!) |
| Array pointer loop: `a := t.a; _ = *a` və ya `t.a[:]` | nil-check loop-dan çıxır, 40% ↓ |
| `for i := range` > `for _, v := range` | hətta kiçik elementlərdə 12% ↓ |

## 6. Array/Slice əməliyyatları (Ch 8)
```go
// Clone (standart):         make+copy, amma 4 ŞƏRTLİ (təmiz id + 2-arg make + copy ifadə-siz)
y = make([]T, len(s)); copy(y, s)
// ≤64B kopya (Go 1.17+):    2-3× sürətli
*(*[64]byte)(d) = *(*[64]byte)(s)
// Sıfırla (memclr):         for i := range s { s[i] = zero } — vectorized
// Böyüt (bir addımda):      n hesabla → make([]T, 0, n) və ya Grow
// Subslice 3-indeks:        s[i:i+4:i+4] — hot loop-da yoxlamasız
// Artıq append olmayacaqsa: append(x[:len(x):len(x)], ...)
// Insert (allocation-sız):  s = s[:len(s)+len(vs)]; copy(s[i+len(vs):], s[i:]); copy(s[i:], vs)
```
- append böyüməsi: <256 cap → 2×; sonrası ~1.25× (1.18+ monoton)

## 7. String/Byte (Ch 9)
| Hall | Qayda |
|---|---|
| Bir statement concat | `+` (stack ≤32B mümkün) |
| Döngülü concat | `strings.Builder` + `Grow(n)` mütləq |
| map OXU `m[string(b)]` | allocation-sız |
| map SET/`++` `m[string(b)]++` | 1 allocation — sıx modifikasiyada `map[K]*V` və ya index-table |
| Case-insensitive | `strings.EqualFold` (5.6× ↓ ToLower-dən) |
| Çoxhissəli açar | `[2]string{a,b}` array açarı, concat YOX (3.5×) |
| Writer-ə string | bufferli `WriteString` wrapper / `bufio.Writer` |
| >32B nəticə | həmişə heap — sayını azalt |

## 8. BCE (Ch 10)
```bash
go run -gcflags="-d=ssa/check_bce" main.go   # qalan yoxlamaları gör
```
| Qayda | Nümunə |
|---|---|
| Ən böyük indeks əvvəl | `s[3]\|s[2]\|s[1]\|s[0]` — 1 yoxlama |
| Slice-ı YENİDƏN təyin et | `is = is[:256]` (işləyir); `_ = is[:256]` (İŞLƏMİR) |
| Loop şərti `len(buf)` | `i < len(buf)` BCE; `i <= n` YOX |
| Qlobal slice → lokal | `s := s; for i := range s` |
| Array > slice (BCE üçün) | `a[100]`, `a[n byte]` — yoxlamasız |

## 9. Map (Ch 11)
| Qayda | Təsir |
|---|---|
| `m[k]++` ≥ `m[k] = m[k]+1` | 43% fərq |
| Açar/element pointer-siz → GC scan yox | `[32]byte` açar string əvəzinə (65K entry+) |
| Söz-sayğac: map+`[]int` index-table | 7× sürət, 0 alloc |
| Bool/enum map əvəzinə `[2]T` + `b2i()` | map-switch 11× yavaş; index-table = if-else |
| `make(map, n)` pre-allocate | rehash-lardan qaç |
| Təmizlənən map-in backing array-i QALIR | tam azad: `m = nil` |

## 10. Channel (Ch 12)
NoSync(2ns) < Atomic(7ns) < Mutex(14ns) < **Channel(61ns)**
- Sadə sinxronizasiyada channel YOX
- Multi-case select: kanalları struct-elementli bir kanala birləşdir (1295→851ns)
- Try-send/receive (`select+default`): 5.3-5.6ns — xüsusi optimizasiya

## 11. Funksiya (Ch 13)
| Qayda | Nümunə |
|---|---|
| Hot funksiya inline-olsun: cost<80 | `-gcflags="-m -m"` ilə yoxla |
| Soyuq yolu ayır | hot `concat` + `//go:noinline concatSlow` |
| ≤4 word param value, böyük pointer | Add4: value 2.7ns; Add5: pointer 12ns |
| defer loop-da = 10× | anonim funksiya ilə Loop-dan çıxar (veyə defer-siz) |
| Arqument həmişə qiymətlənir | `debugOn && debugPrint(h+w)` — short-circuit |
| defer/recover/go/rekursiya inline POZULUR | hot path-dən kənarlaşdır |
| for-range plain-for-dan ucuz (1.18+) | inline cost 10 vs 18 |

## 12. Interface (Ch 14)
| Boxing | Xərc |
|---|---|
| Pointer / map / channel / func / sabit / zero-size / bool / int8 | ~1.2ns, 0 alloc |
| [0,255] kiçik int, zero float/string/slice | ~3.4ns |
| Kənar int, non-zero float | ~20× pointer |
| Non-constant string, non-nil slice | **~50×** — qaçın! |
| Böyük struct (məs. [100]int) | 592ns — pointer box et |

```go
// İki interfaceə: bir box-la — y = x (assign allocation-sız)
// fmt çox-eyni-dəyər: var i interface{} = x; fmt.Fprint(w, i, i, i)
// Kiçik dəyər çoxluğu: var values [N]uint16; box(&values[v]) — 20× ↓
```
- Interface çağırış ~8× (vtable 2.5ns + inline-itkisi) — hot loop-da konkret tip
- API sərhəddində konkret tiplər (image.RGBA64Image modeli)

---

## Qızal qaydalar (hamısı üçün)
1. **Ölç, sonra optimallaşdır** — `pprof`, `benchmark`, `AllocsPerRun`
2. **Hot path-i tap** — kodun 90%-inə toxunma
3. **Compiler və runtime-a güvə, amma yoxla** — `-m`, `-d=ssa/check_bce`, `gctrace`
4. **Versiya bağlılığını** yadda saxla — fəndlər v1.19 əsasıdır
5. **Oxunaqlıq > mikro-optimizasiya** — trade-off (kitabın öz prinsipi!)
