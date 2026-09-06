# Chapter 7 — Structs (Strukturlar)

## Bu chapter nədən bəhs edir?
Struct sahələrinə pointer üzərindən loop daxilində erişməyin xərcinə və kiçik struct-ların xüsusi optimizasiyasına. Qısa fəsil — Ch 2 və Ch 6-nın struct kontekstinə tətbiqi.

## Əsas fikirlər

### 1. Struct pointer-inin sahələrinə loop daxilində erişməyin
**Nədir:** `t.x += i` loop daxilində — hər iterasiyada pointer dereference = memory əməliyyatı.

**Necə işləyir:** CPU register-lə yaddaş arasında fərq — register-də cəm yığ, sonda bir dəfə sahəyə yaz.

**Kitabdan kod nümunəsi:**
```go
const N = 1000
type T struct { x int }

//go:noinline
func f(t *T) {                // 2402 ns/op — memory hər iterasiyada
    t.x = 0
    for i := 0; i < N; i++ { t.x += i }
}
//go:noinline
func g(t *T) {                // 461.3 ns/op — 5.2× sürətli!
    var x = 0                 // lokal register dəyişəni
    for i := 0; i < N; i++ { x += i }
    t.x = x                   // bir dəfə yaz
}
```
**Sub-kod izahı:**
- `f` ≡ `h` — `x := &t.x; for { *x += i }` — pointer-in loop daxilində dereference-i
- `g` — cəmi loop boyu register-də, yalnız 2 memory əməliyyatı (başda ox, sonda yaz)
- Bu, Ch 6.2-dəki "avoid pointer dereferences in a loop" qaydasının struct variantıdır

**Nəyə lazımdır:** Hot loop-larda struct sahəsi (və ya map value, slice element pointer-ı) yığılan kimi — lokal dəyişənə köçür, sonda geri yaz.

**Aliasing xəbərdarlığı (Ch 6-dan):** Əgər `t.x` loop daxilində başqa pointer ilə dəyişilə bilirsə, bu transformasiya semantikanı dəyişir — əmin deyilsənsə tətbiq etmə.

### 2. Kiçik struct-lar xüsusi optimizasiya olunur
Ch 2.6-ya istinad: ≤4 native-word sahəli struct-lar (və ≤4 elementli array-lər) kopyalanması sürətli yoldan icra olunur; 9→10 sahə keçidi 2-7× yavaşlama yaradır. Böyük struct kopyalarından hot path-də qaçın — pointer ötür.

### 3. Sahə sırası ilə struct-u kiçilt
Ch 2.5-ə istinad: sahələri ölçüyə görə qruplaşdırmaq padding-i azaldır (T1: 24 bayt → T2: 16 bayt). Milyonlarla instans yaratdıqda real qənaət.

## Əsas terminlər
- Field access through pointer (pointer ilə sahə erişimi)
- Register vs memory (reqistr vs yaddaş)
- Small-size struct (kiçik ölçülü struktur) — ≤4 word
- Struct padding (struktur doldurulması)

## Praktik nəticə
1. Loop daxilində struct sahəsi yığırsansa: `var x = t.x` → loop → `t.x = x` — 5× sürət.
2. Hot struct-ları ≤4 word saxla; böyükləri pointer ilə ötür.
3. Massiv instanslı struct-larda sahə sırasını ölçüyə görə düz.

## Mənbə
Pages: 74-75 (PDF səh. 74-75)
