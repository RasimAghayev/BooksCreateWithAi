# Chapter 10 — About Time (Vaxt Haqqında)

## Bu fəsil nədən bəhs edir?

time paketi — vaxt yaratma (Now, Date, Unix), müqayisə (After/Before/Equal), duration
hesablama (Sub, Hours/Minutes/Seconds/Nanoseconds), SLA deadline patterni, vaxt
manipulyasiyası (Add, mənfi duration, AddDate), Parse (RFC3339/UnixDate/ANSIC),
Format, LoadLocation/In ilə timezone keçidi və 5 activity (format, ölçmə, gələcək
tarix, çoxzonalı çap).

## Əsas fikirlər

### 1. Vaxt Yaratma
**İcra vaxtının ölçülməsi (ən common pattern):**
```go
start := time.Now()
time.Sleep(2 * time.Second)       // gecikmə simulyasiyası
end := time.Now()
```

**Həftə günü/saata görə davranış (kitabdan):**
```go
Day := time.Now().Weekday()
Hour := time.Now().Hour()
if Day.String() == "Monday" {      // release günü
    if Hour >= 1 {
        fmt.Println("Performing full blown test!")
    } else { ... }
}
```

**Log fayl adı patterni:**
```go
AppName := "HTTPCHECKER"
Date := time.Now()
LogFileName := AppName + "_" + Action + "_" +
    strconv.Itoa(Date.Year()) + "_" + Date.Month().String() + "_" +
    strconv.Itoa(Date.Day()) + ".log"
// HTTPCHECKER_BASIC_2019_September_27.log
```
**Qeyd:** time tipləri string-ə implicit çevrilmir — int hissələr üçün strconv.Itoa;
Month() AD qaytarır (September), rəqəm üçün int(Date.Month()).

### 2. Müqayisə — After/Before/Equal
```go
now := time.Now()
only_after, _ := time.Parse(time.RFC3339, "2020-11-01T22:08:41+00:00")
if now.After(only_after) {        // ŞƏRTLİ İCRA
    fmt.Println("Executing actions!")
} else {
    fmt.Println("Now is not the time yet!!")
}

now_too := now                    // kopya — EYNI an
time.Sleep(2 * time.Second)
later := time.Now()
now.Equal(now_too)   // true
now.Equal(later)     // false
```
3 müqayisə funksiyası: `After(t)`, `Before(t)`, `Equal(t)` — bool qaytarır.

### 3. Duration Hesablaması — Sub
```go
Start := time.Now()
// ... uzun əməliyyat ...
End := time.Now()
Duration := End.Sub(Start)        // time.Duration — int64 əsaslı

Duration.Hours()       // 4 rezolyusiya:
Duration.Minutes()     //   Hours, Minutes,
Duration.Seconds()     //   Seconds, Nanoseconds
Duration.Nanoseconds() // (gün/həftə hesabla özün)
```

**elapsedTime helper (kitabdan):**
```go
func elapsedTime(start time.Time, end time.Time) string {
    Elapsed := end.Sub(start)
    Hours := strconv.Itoa(int(Elapsed.Hours()))
    Minutes := strconv.Itoa(int(Elapsed.Minutes()))
    Seconds := strconv.Itoa(int(Elapsed.Seconds()))
    return "The total execution time elapsed is: " + Hours +
        " hour(s) and " + Minutes + " minute(s) and " + Seconds + " second(s)!"
}
```

### 4. SLA/Deadline Patterni
```go
deadline := time.Duration((600 * 10) * time.Millisecond)   // 6 saniyə
Start := time.Now()
// ... əməliyyat ...
End := time.Now()
Duration := End.Sub(Start)
TransactionTime := time.Duration(Duration.Nanoseconds()) * time.Nanosecond
if TransactionTime <= deadline {
    fmt.Println("Performance is OK")
} else {
    fmt.Println("Performance problem!")
}
```
**6 rezolyusiya:** Hour, Minute, Second, Millisecond, Microsecond, Nanosecond.
Müqayisə üçün duration-lar eyni vahiddə normalize olmalıdır.

### 5. Vaxt Manipulyasiyası — Add/AddDate
```go
TimeToManipulate := time.Now()
ToBeAdded := time.Duration(10 * time.Second)
TimeToManipulate.Add(ToBeAdded)                  // +10 saniyə

ToBeAdded := time.Duration(-10 * time.Minute)   // MƏNFİ = geri get!
TimeToManipulate.Add(ToBeAdded)                  // 10 dəqiqə əvvəl

date := time.Date(2019, 9, 27, 18, 50, 48, 324359102, time.UTC)
next_date := date.AddDate(1, 2, 3)               // +1 il, 2 ay, 3 gün
// 2020-11-30 18:50:48.324359102 +0000 UTC
```

### 6. Parse — String-dən time-a
```go
t1, _ := time.Parse(time.RFC3339, "2019-09-27T22:18:11+00:00")  // OK
t2, _ := time.Parse(time.UnixDate, "2019-09-27T22:18:11+00:00") // SƏHV format!
t3, _ := time.Parse(time.ANSIC, "2019-09-27T22:18:11+00:00")    // SƏHV format!
```
**Qayda:** parse STRING format sabitinə UYĞUN olmalıdır — yoxsa zero-time qaytarır.
**3 standart:** RFC3339 (`2019-09-27T22:18:11+00:00`), UnixDate (`Mon Sep 27...`),
ANSIC (`Thu Oct 17 13:56:03 2019`).

**Go-nun referens vaxtı:** `Mon Jan 2 15:04:05 -0700 MST 2006` — 01/02 03:04:05
05/06 06:09 patterni; custom format üçün bu şablondan istifadə olunur.

### 7. Format — time-dən string-ə
```go
whatstheclock := func() string {
    return time.Now().Format(time.ANSIC)   // "Thu Oct 17 13:56:03 2019"
}
```

### 8. Timezone — LoadLocation/In
```go
func timeDiff(timezone string) (string, string) {
    Current := time.Now()
    RemoteZone, _ := time.LoadLocation(timezone)   // (time, error) qaytarır!
    RemoteTime := Current.In(RemoteZone)          // həmin zonada SƏRHƏDLİ nüsxə
    return Current.Format(time.ANSIC), RemoteTime.Format(time.ANSIC)
}
timeDiff("America/Los_Angeles")
```

### 9. Activity-lər
- **10.01/10.02:** custom format `15:32:30 2019/10/17` — hissə-hissə ayır + Itoa +
  concat; Month AD verir → int(Month()) lazım olur
- **10.03:** Sleep(2s) ölçmə — Sub().Seconds() → "2.0016895 seconds" (bir az çox —
  overhead normaldır)
- **10.04:** +6 saat 6 dəq 6 san — Add(duration) ANSIC çapı
- **10.05:** NY + LA vaxtı — 2 LoadLocation + In()

## Əsas terminlələr
- time.Now() — cari an; monotonic + wall clock
- Wall Clock — NTP-sinxronlaşdırıla bilən, dəyişən saat
- Monotonic Clock — geri getməyən, yalnız irəli ölçən
- time.Duration — vaxt fərqi tipi (int64 əsaslı)
- Sub() — end.Sub(start) → Duration
- Hours/Minutes/Seconds/Nanoseconds — Duration konvertorları
- After/Before/Equal — vaxt müqayisə funksiyaları
- Add(duration) — vaxta duration əlavə; MƏNFİ = geri
- AddDate(y, m, d) — il/ay/gün əlavə
- time.Date(...) — xüsusi an yaratma (8 komponent + Location)
- Parse(format, s) — string → time; format uyğunsuzluğu = zero-time
- Format(format) — time → string
- RFC3339/UnixDate/ANSIC — hazır format sabitləri
- Referens Time `01/02 03:04:05 06` — custom format şablonu
- LoadLocation(zone) — timezone obyekti (America/Los_Angeles)
- In(loc) — həmin zonada təqdimat
- Epoch — UNIX vaxtın başlanğıcı (1 yanvar 1970)

## Praktik nətidə

(1) Start/End + Sub() — performans ölçmənin standart üçlüyü. (2) Duration 4 əsas
rezolyusiyada qaytarır; gün/həftə özün hesabla. (3) Deadline müqayisəsi üçün hər iki
tərəfi EYNİ rezolyusiyaya normalize et. (4) Add() yalnız Duration qəbul edir —
`time.Duration(10 * time.Second)` ilə hazırla. (5) Mənfi duration — keçmişə getməyin
ən sadə yolu (Sub qəlizdir). (6) Month() AD, Day()/Year() rəqəm — İtoa concat-dan
əvvəl. (7) Parse-də string formatı SABİTLƏ uyğun olmalıdır — yoxsa zero-time;
_ error-u yudsan da səhv nəticə səni çaşdıra bilər. (8) LoadLocation 2 dəyər qaytarır
— error-u yoxla. (9) SLA proqramlaşdırması: Duration əsaslı müqayisə — 6 rezolyusiya
seçiminin praktik tətbiqi. (10) Sleep ilə simulyasiya + Seconds() — ölçmə testlərinin
ən sadə formaları.

## Mənbə
Pages: 353-373 (PDF 386-407)
