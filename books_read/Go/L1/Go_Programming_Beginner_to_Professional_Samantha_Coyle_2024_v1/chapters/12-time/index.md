# Chapter 12 — About Time (səh. 398-419)

## Bu fəsil nədən bəhs edir?

time paketi: vaxt yaratma (Now, Date), müqayisə (After/Before/Equal),
Duration hesablanması (Sub, SLA yoxlaması), vaxt idarəetməsi (Add,
mənfi duration), Parse (RFC3339/UnixDate/ANSIC), Format, LoadLocation
+ In (saat qurşaqları) və AddDate.

## Əsas fikirlər

### 1. Vaxt yaratma
```go
start := time.Now()
fmt.Println("The script has started at: ", start)
time.Sleep(2 * time.Second)          // 2 saniyə gözlə
end := time.Now()
// Çıxışda: 2023-09-27 08:19:33.8358274 +0200 CEST m=+0.001998701
```

**Həftə günü + saat idarəetməsi (test strategiyası nümunəsi):**
```go
day := time.Now().Weekday()
hour := time.Now().Hour()
if day.String() == "Monday" {
    if hour >= 1 {
        fmt.Println("Performing full blown test!")
    } else {
        fmt.Println("Performing hit-n-run test!")
    }
} else {
    fmt.Println("Performing hit-n-run test!")
}
```

**Log fayl adı (gündəlik):**
```go
appName := "HTTPCHECKER"
action := "BASIC"
date := time.Now()
logFileName := appName + "_" + action + "_" +
    strconv.Itoa(date.Year()) + "_" + date.Month().String() + "_" +
    strconv.Itoa(date.Day()) + ".log"
// HTTPCHECKER_BASIC_2024_March_16.log
```
- Year()/Day() → int — stringlə birləşdirmək üçün strconv.Itoa lazımdır

### 2. Vaxt müqayisəsi
```go
// Parse ilə deadline:
now := time.Now()
onlyAfter, err := time.Parse(time.RFC3339, "2020-11-01T22:08:41+00:00")
if err != nil {
    fmt.Println(err)
}

if now.After(onlyAfter) {        // After — sonra?
    fmt.Println("Executing actions!")
} else {
    fmt.Println("Now is not the time yet!!")
}

// Before — əvvəl?; Equal — bərabər?
nowToo := now
time.Sleep(2 * time.Second)
later := time.Now()
now.Equal(nowToo)    // true
now.Equal(later)     // false
```
- **Wall clock vs monotonic clock:** wall — NTP ilə sinxronlaşan
  (dəyişə bilər); monotonic — yalnız irəli, ölçmə üçün etibarlı

### 3. Duration hesablanması
```go
start := time.Now()
sum := 0
for i := 1; i < 10000000000; i++ {
    sum += i
}
end := time.Now()
duration := end.Sub(start)                  // Duration tipi

fmt.Println(duration.Hours())               // 4 resolution:
fmt.Println(duration.Minutes())             //   Hours, Minutes,
fmt.Println(duration.Seconds())             //   Seconds,
fmt.Println(duration.Nanoseconds())         //   Nanoseconds
// Gün/həftə/ay üçün özün hesabla
```

**SLA/deadline yoxlaması:**
```go
deadlineSeconds := time.Duration((600 * 10) * time.Millisecond)   // 6s
start := time.Now()
// ... iş ...
end := time.Now()
duration := end.Sub(start)
transactionTime := time.Duration(duration.Nanoseconds()) * time.Nanosecond

if transactionTime <= deadlineSeconds {
    fmt.Println("Performance is OK transaction completed in", transactionTime)
} else {
    fmt.Println("Performance problem!")
}
```
- 6 resolution: Hour/Minute/Second/Millisecond/Microsecond/Nanosecond
- Duration-ları birbaşa müqayisə üçün EYNİ rezolyutsiyaya çevir

### 4. Vaxt idarəetməsi (Add)
```go
timeToManipulate := time.Now()
toBeAdded := time.Duration(10 * time.Second)
fmt.Printf("%v duration later %v", toBeAdded, timeToManipulate.Add(toBeAdded))
// 10s sonra

// Çıxarmaq üçün MƏNFİ duration:
toBeAdded := time.Duration(-10 * time.Minute)
timeToManipulate.Add(toBeAdded)     // 10 dəqiqə ƏVVƏL
```

**elapsedTime funksiyası (Exercise 12.02):**
```go
func elapsedTime(start time.Time, end time.Time) string {
    elapsed := end.Sub(start)
    hours := strconv.Itoa(int(elapsed.Hours()))
    minutes := strconv.Itoa(int(elapsed.Minutes()))
    seconds := strconv.Itoa(int(elapsed.Seconds()))
    return "The total execution time elapsed is: " + hours +
        " hour(s) and " + minutes + " minute(s) and " + seconds + " second(s)!"
}

start := time.Now()
time.Sleep(2 * time.Second)
end := time.Now()
fmt.Println(elapsedTime(start, end))
```

### 5. Parse — stringdən vaxta
**Referens format:** `Mon Jan 2 15:04:05 -0700 MST 2006` (Go-nun
məşhur "1 2 3 4 5 6" sxemi — POSIX əsaslı).

**Əsas formatlar:**
- RFC3339: `2019-09-27T22:18:11+00:00`
- UnixDate: `Mon Sep 27 18:24:05 2019`
- ANSIC: `Thu Oct 17 13:56:03 2023`

```go
t1, err := time.Parse(time.RFC3339, "2019-09-27T22:18:11+00:00")   // ✓
t2, err := time.Parse(time.UnixDate, "2019-09-27T22:18:11+00:00")  // XƏTA — format uyğunsuz
// cannot parse ... as "Mon"
```
- Format uyğun gəlməzsə error; zero-time qaytarır
- **Epoch:** 1 yanvar 1970 — Unix vaxt başlanğıcı

### 6. Format + Date yaratma
```go
// Öz tarixini yarat:
date := time.Date(2019, 9, 27, 18, 50, 48, 324359102, time.UTC)
// func Date(year, month, day, hour, min, sec, nsec int, loc *Location) Time

// AddDate — il/ay/gün əlavə et:
nextDate := date.AddDate(1, 2, 3)     // 2020-11-30 ...

// Format — istədiyin şəkildə çap:
date.Format(time.ANSIC)               // "Thu Oct 17 13:56:03 2023"
```

### 7. Saat qurşaqları (LoadLocation + In)
```go
current := time.Now()
losAngeles, err := time.LoadLocation("America/Los_Angeles")
if err != nil {
    fmt.Println(err)
}
fmt.Println("The local current time is:", current.Format(time.ANSIC))
fmt.Println("The time in Los Angeles is:", current.In(losAngeles).Format(time.ANSIC))
// Fri Oct 18 08:14:48 2019 → Thu Oct 17 23:14:48 2019
```

**timeDiff funksiyası (Exercise 12.03):**
```go
func timeDiff(timezone string) (string, string) {
    current := time.Now()
    remoteZone, err := time.LoadLocation(timezone)
    if err != nil {
        fmt.Println(err)
    }
    remoteTime := current.In(remoteZone)
    return current.Format(time.ANSIC), remoteTime.Format(time.ANSIC)
}
fmt.Println(timeDiff("America/Los_Angeles"))
```

## Activity icmalları
- **12.01-12.02:** Custom format `15:32:30 2023/10/17` — Now() +
  strconv.Itoa hissə-hissə (Month() AD qaytarır — int çevir!)
- **12.03:** Sleep(2s) + Sub + Seconds → "took exactly 2.0016895 seconds"
- **12.04:** 6h 6m 6s sonra — Duration((6*3600+6*60+6)*Second) + Add
- **12.05:** NY və LA vaxtı — LoadLocation("America/New_York") +
  In() + ANSIC format

## Əsas terminlər
- time.Now / time.Date — indiki / xüsusi vaxt
- Weekday()/Hour()/Year()/Month()/Day() — komponentlər
- time.Sleep — gözləmə
- After/Before/Equal — müqayisələr
- time.Sub — fərq → Duration
- Duration resolutions — H/M/S/Ms/µs/ns
- Add / mənfi duration — irəli/geri
- AddDate — il/ay/gün əlavə
- Parse/Format — string↔Time (RFC3339, UnixDate, ANSIC)
- Referens vaxtı — "Mon Jan 2 15:04:05 -0700 MST 2006"
- Epoch — 1970-01-01
- LoadLocation / In — saat qurşağı çevrilməsi
- Wall vs monotonic clock — NTP/dəyişən vs ölçmə etibarlı

## Praktik nəticə
İcra müddəti: start := Now() → iş → Sub → Hours/Minutes/Seconds.
Deadline yoxlaması: hər iki tərəfi EYNİ rezolyutsiyada Duration-a
çevir, sonra müqayisə et. Gələcək/keçmiş: Add(+/-Duration). String
tarixlər: Parse (format KOMPATİB olmalıdır — error yoxla!); çıxış üçün
Format. Qurşaqlar: LoadLocation + In. Ay adı Month() stringdir —
rəqəm lazımdırsa int(Month()). Şəxsi formatda Year/Day int-dir —
strconv.Itoa ilə birləşdir.

## Mənbə
Pages: 398-419 (PDF 398-419)
