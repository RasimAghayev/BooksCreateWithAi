# Chapter 12 — Channels (Kanallar)

## Bu chapter nədən bəhs edir?
Channel-lərin sinxronizasiya aləti kimi xərcinə, select bloklarının case-sayıyla xərcinə, çoxkanal → birmənalı birləşdirmə strategiyalarına və try-send/try-receive xüsusi optimizasiyasına.

## Əsas fikirlər

### 1. Channel ən yavaş sinxronizasiya yoludur (v1.19)
**Kitabdan benchmark:**
```go
Benchmark_NoSync:    2.25 ns    // paylaşma yoxdur — ideal
Benchmark_Atomic:    7.11 ns    // atomic.AddInt32
Benchmark_Mutex:     14.25 ns   // sync.Mutex
Benchmark_Channel:   61.44 ns   // ch <- struct{}{}; g++; <-ch — 27× yavaş!
```
**Praktik qayda:** Sadə sayğaç/qiymət sinxronizasiyası üçün kanal İŞLƏTMƏ — ən yaxşısı paylaşmadan qaçmaq, sonra atomic, sonra mutex. Channel — məntiqi axın üçün (mesaj, hadisə, pipeline), sadə mutex-konflikt üçün yox.

### 2. Select bloku — case sayı artdıqca xərc artır
```go
Benchmark_Select_OneCase:   58.90 ns
Benchmark_Select_TwoCases:  115.3 ns   // 2× — hər case üçün iterativ yoxlama
```
- 1 case-li select (defaultsuz) adi channel əməliyyatı kimi compilə olunur
- Case-lərin sayını minimallaşdır

### 3. Bir neçə kanalı BİRİNƏ birləşdir — select-dən yayın
**Nədir:** Fərqli tip kanalların əvəzinə vahid element tipili (interface və ya struct) bir kanal + növ ayrımı.

**Kitabdan kod nümunəsi:**
```go
// Sürət sırası (2 kanal → select ən yavaş):
Benchmark_TwoChannels:            1295 ns  // select { <-x | <-y }
Benchmark_OneChannel_Interface:   940.9 ns // type switch ilə ayırma
Benchmark_OneChannel_Struct:      851.0 ns // sahə yoxlaması ilə ayırma — ƏN SÜRƏTLİ

// Struct variantı:
type T struct { x int; y string }
var x = make(chan T)
go func() { for { x <- T{x: 1} } }()
go func() { for { x <- T{y: "hello"} } }()
for {
    v := <-x
    if v.y != "" { vy = v.y } else { vx = v.x }   // mesaj növü sahə ilə
}
```
**Sub-kod izahı:**
- Interface variantı: `switch v := v.(type)` — boxing allocation-ları gətirir
- Struct variantı: hər növ öz sahəsində — boxing yoxdur, sıfır-dəyər ilə ayrım

### 4. Try-send / try-receive — xüsusi optimizasiya
**Nədir:** 1 case + 1 default olan select bloku — compiler tərəfindən xüsusi, sürətli yolla compilə olunur:
```go
select {
case <-c:      // try-receive
default:
}
select {
case c <- v:   // try-send
default:
}
```
**Nəticə:** Bu forma adi multi-case select-dən qat-qat sürətlidir — non-blocking yoxlama lazım olanda üstünlük ver.

## Əsas terminlər
- Synchronization (sinxronizasiya) — paylaşılan vəziyyət qoruması
- Atomic operation (atomar əməliyyat) — `sync/atomic` paketi
- select block — çoxkanal gözləmə konstruksiyası
- Try-send/try-receive — default-lu non-blocking əməliyyat
- Channel merging (kanal birləşdirmə) — vahid kanal + növ sahəsi
- Boxing penalty — interface elementin allocation xərci

## Praktik nəticə
1. Sadə sayğaclarda channel yox: atomic (7ns) > mutex (14ns) > channel (61ns).
2. Ən yaxşısı paylaşmamaktır — hər goroutine öz lokal nüsxəsini saxlasın, nəticəni birləşdirsin.
3. Multi-case select-ləri struct-elementli tək kanala çevir — 1.5× sürət + daha sadə axın.
4. Non-blocking poll lazımdırsa `select+default` işlət — xüsusi optimizasiyalı formadır.

## Mənbə
Pages: 125-129 (PDF səh. 125-129)
