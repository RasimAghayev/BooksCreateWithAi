# Chapters 3-4 — Data Types, Control Structures (#17-#38) (səh. 76-132)

## Bu fəsillər nədən bəhs edir?

(3) Tam ədəd overflow-ləri, floating point dəqiqliyi, slice len/cap mexanikası,
nil vs empty slice, append yan təsirləri, slice/map memory leak-lər, müqayisə
üsulları. (4) range loop kopyaları, ifadənin bir dəfə qiymətləndirilməsi,
pointer element tələsi, map iterasiya qaydasızlığı, break/switch tələsi, loop
daxilində defer.

## Əsas səhvlər və həllər

### #17: Integer overflow (susuz!)
- `math.MaxInt` yoxlaması: `if counter == math.MaxInt { panic(...) }`
- Vurma yoxlaması: `(a*b)/b != a` → overflow olub
- 0-la başlayan literal = OCTAL: `100 + 010 = 108`!

### #18: Floating point
- float32/64 = approksimasiya (sign+exponent+mantissa)
- **Əməliyyat SIRASI dəqiqliyi dəyişir:** `a*(b+c)` ≠ `a*b + a*c` çox vaxt
- Böyük+kiçik toplama: əvvəl KİÇİKLƏRİ topla

### #19-#23: Slice mexanikası
- **make(len, cap):** ptr+len+cap; dolanda → kapasite 2× → KOPYA
- **Slicing paylaşır:** `s2 := s1[1:]` — EYNİ backing array! `append(s2, ...)` →
  s1-in elementləri DƏYİŞƏ BİLƏR (əgər s2 dolmayıbsa)
- **Təkmil init:** uzunluq məlum → `make([]T, 0, n)` — 400%-dək sürət fərqi
- **nil vs empty:** `[]T{}` (nil=false) vs `var s []T` (nil=true); JSON-da
  `null` vs `[]` FƏRQİ! Boşluq yoxlaması: `len(s) == 0` (hər ikisini tutur)

### #24: copy istiqaməti
- `copy(dst, src)` — dst BİRİNCİ; `append` slice-a bağlı mütərəqqi

### #25: append yan təsiri
```go
s := []int{1, 2, 3}
f(s[:2])          // f daxilində append → s[2] silinə bilər!
// həll: kopya ilə: sCopy := make([]int, 2); copy(sCopy, s)
```

### #26: Slice memory leak-ləri
- **Capacity leak:** 1M elementli slice-dən 5KB götürsən `s[:5000]` → backing
  array 1M QALIR. Həll: kopya (və ya `slices.Clip`)
- **Pointer leak:** pointer-sahəli struct elementləri slice-də qalır — 0-a
  endirmək GC-ni azad etmir; elementləri nil et

### #27-#28: Map
- Bucket mexanikası; `make(map[K]V, n)` — n element üçün bucket hazırlığı
- **Map YALNIZ böyüyür, kiçilmir ZƏNCİRLİ — recreating (kopya-yeni map) və ya pointer dəyərlər

### #29: Müqayisə
- `==` yalnız comparable tiplərdə (string, chan, interface, pointer sahəsiz struct)
- `reflect.DeepEqual` — universal amma YAVAŞ; `bytes.Compare` kimi
  optimallar var; performance kritik → custom metod

### #30: range kopyası
```go
for _, account := range accounts {
    account.balance += 100  // KOPYA dəyişir — orijinal YOX!
}
// düzgün: indeks ilə: accounts[i].balance += 100
// və ya: for _, a := range accountsPointer
```

### #31: range ifadəsi BİR DƏFƏ qiymətlənir
```go
s := []int{0,1,2}
for range s {
    s = append(s, 10)   // sonsuz DEYİL — len 3 kimi saxlanılır
}
// ARRAY: kopyası üzərində iterasiya → dəyişiklik GÖRÜNMÜR
// &a ilə range → orijinala baxır
```

### #32: pointer element + range
```go
type Customer struct{ ID string }
customers := []Customer{...}
ids := map[string]*Customer{}
for _, customer := range customers {
    ids[customer.ID] = &customer  // HAMISI SON elementə işarə edir!
}
// həll: customer := customer (kölgə) və ya indeks: &customers[i]
```

### #33: map iterasiyası
- **Sıra QAYDASIZDIR — qətiliklə güvənmə** (Go şüurlu randomlaşdırır)
- Map daxilində INSERT/DELETE iterasiya zamanı: yeni element gələ bilər
  (0 və ya 1 dəfə) — silinən "gözlərilə bilər"

### #34: break nəyi break edir?
```go
for {
    switch i {
    case 2:
        break  // SWITCH-i break edir, loop-u YOX!
    }
}
// select + ctx.Done() eyni tələ
// həll: label  və ya return / goto
```

### #35: defer loop daxilində
```go
// PIS: hər iterationda defer YIĞILIR — fayl yalnız funksiya sonunda bağlanır
for path := range ch {
    f, _ := os.Open(path)
    defer f.Close()   // LEAK axını sona qədər!
}
// həll 1: helper funksiya (hər fayl üçün scope)
// həll 2: closure func() error { defer f.Close(); ... }()
```

## Əsas terminlər

- Backing Array (arxa massiv)
- Slice Header (ptr/len/cap)
- Reslicing vs Copy
- Capacity Leak
- Octal Literal (0-prefix)
- Named Result Parameter
- Loop Label

## Praktik nəticə

- Slice paylaşımını xəritədə tuta: kəsəndə APPEND yan təsirini düşün
- Böyük slice-dən kiçik lazımsa — CLONE/kopya et
- range: element kopyasıdır; struct dəyişmək → indeks/pointer
- map sırasına güvənma; defer loop-da YOX

## Mənbə

Pages: 76-132 (Chapters 3-4, 100 Go Mistakes 2nd ed.)
