# Chapter 6 — Pointers (Pointer-lər)

## Bu chapter nədən bəhs edir?
Loop daxilində lazımsız nil pointer yoxlamalarından və pointer dereference-lərindən qaçmaqla loop performansını artırmaq texnikalarına. Compiler-ın (v1.19) iki konkret boşluğu və onların kod-səviyyəli həlləri.

## Əsas fikirlər

### 1. Loop daxilində lazımsız nil array pointer yoxlamasından qaçın
**Nədir:** Compiler hər `a[i]` yazısında `a`-nın nil olmasını loop İÇINDƏ yoxlayır (TESTB instruksiyası) — baxmayaraq ki, bir dəfə yoxlama kifayətdir.

**Kitabdan kod nümunəsi:**
```go
//go:noinline
func g0(a *[N]int) {            // 517.6 ns/op — TESTB loop-daxili
    for i := range a { a[i] = i }
}
//go:noinline
func g1(a *[N]int) {            // 398.1 ns/op — TESTB loop-dan kənarda
    _ = *a                      // bir dəfə dereference → nil yoxlaması məcbur edir
    for i := range a { a[i] = i }
}
//go:noinline
func g2(x *[N]int) {            // slice törət — g1 qədər sürətli
    a := x[:]
    for i := range a { a[i] = i }
}
```
**Sub-kod izahı:**
- `_ = *a` → nil pointer-də panic baş verəcəyi üçün compiler yoxlamanı bir dəfə, loop-dan ƏVVƏL edir
- Assembly sübutu: g0-da TESTB hər iterasiyada, g1-də bir dəfə
- `g2` — daha təmiz həll: array pointer-dən slice törədib slice üzərində loop

### 2. Array pointer struct SAHƏSİDİRSə — lokal kopya çıxarın
**Nədir:** `t.a` sahəsi üzərində loop edəndə `_ = *t.a` kömək ETMİR (yoxlama yenə loop içində qalır).

**Kitabdan kod nümunəsi:**
```go
type T struct { a *[N]int }

//go:noinline
func f3(t *T) {            // 390 ns — ƏN SÜRƏTLİ
    a := t.a               // sahəni LOKAL dəyişənə kopyala
    _ = *a                 // nil yoxlamasını loop-dan çıxar
    for i := range a { a[i] = i }
}
//go:noinline
func f4(t *T) {            // 388 ns — eyni sürətli, daha təmiz
    a := t.a[:]            // slice törət
    for i := range a { a[i] = i }
}
// f0 (birbaşa t.a): 623 ns, f1 (_=*t.a): 637 ns — fərq yoxdur
```
**Sub-kod izahı:**
- `f1`-in `_ = *t.a` sətri faydasızdır — compiler sahənin loop daxilində dəyişə biləcəyini ehtimal edir
- `a := t.a` lokal kopya → həm nil yoxlaması, həm load loop-dan çıxır
- **Müəllifin tövsiyəsi:** slice variantını (f4) seç — compiler slice-ləri array-lərdən daha yaxşı optimizasiya edir
- İstisna: loop yalnız OXUYURSA (yazmırsa) — f1 də kifayət edir

### 3. Loop daxilində pointer dereference-dən qaçın — register > memory
**Nədir:** CPU register-lər üzərində işləmək yaddaş (memory) üzərində işləməkdən qat-qat sürətlidir; loop daxili hər dereference = hər iterasiyada memory əməliyyatı.

**Kitabdan kod nümunəsi:**
```go
//go:noinline
func f(sum *int, s []int) {      // 3024 ns/op — *sum hər iterasiyada memory-də
    for _, v := range s { *sum += v }
}
//go:noinline
func g(sum *int, s []int) {      // 566.6 ns/op — 5.3× sürətli!
    var n = *sum                 // bir dəfə oxu → register-də saxla
    for _, v := range s { n += v }
    *sum = n                     // bir dəfə yaz
}
```
**Sub-kod izahı:**
- f: MOVQ (AX),SI + ADDQ + MOVQ SI,(AX) — hər iterasiyada load+store
- g: loop yalnız register-lərlə (ADDQ DI,SI) — memory yoxdur
- **MÜHÜM EHTİYAT:** f və g EKVİVALENT DEYİL! Əgər `sum` `s`-in öz elementinə işarə edirsə (`sum = &s[2]`), nəticələr fərqlənir (f: 6, g: 4) — aliasing varsa g istifadə ETMƏ
- Təhlükəsiz alternativ: `h(s []int) int` — dəyəri qaytar, caller-da `*sum += h(s)`

## Əsas terminlər
- Nil pointer check (nil yoxlaması) — TESTB instruksiyası
- Pointer dereference (istiqamətləndirmə) — memory-dən oxu/yaz
- Register vs memory — CPU-nun işləmə rejimləri
- Aliasing (lüğəti üst-üstə düşmə) — iki pointer eyni yaddaşa işarə edir
- `//go:noinline` — benchmark təcridi üçün directive
- Slice derivation (slice törədilməsi) — `x[:]` idiomu

## Praktik nəticə
1. Array pointer üzərində loop edərkən: lokal kopya + `_ = *a`, ya da slice törət (`t.a[:]`) — 25-40% sürət qazancı.
2. Loop daxilində hər hansı pointer/struct-sahə yazışı: lokal dəyişəndə yığ, sonda bir dəfə yaz — 5× sürət.
3. Aliasing ehtimalı varsa (pointer loop datasına işarə edə bilər) — lokal-yığma transformasiyası DÜZGÜNLÜYÜ POZUR, istifadə etmə.
4. Bu compiler boşluqları gələcək versiyalarda düzəilə bilər — hər halda benchmark ilə yoxla.

## Mənbə
Pages: 68-73 (PDF səh. 68-73)
