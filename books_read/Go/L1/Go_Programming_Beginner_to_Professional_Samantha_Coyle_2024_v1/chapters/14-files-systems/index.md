# Chapter 14 — File and Systems (səh. 440-469)

## Bu fəsil nədən bəhs edir?

Fayl sistemi: icazələr (simvolik/oktal), flag ilə arqumentlər (Int/Var
variantları), siqnallar (SIGINT/SIGTERM/SIGKILL, signal.Notify + cleanup),
fayl yaratma/yazma (Create, Write/WriteString, WriteFile), Stat ilə
mövcudluq yoxlaması, oxuma (ReadFile, bayt-bayt Read), fayl backup,
CSV (encoding/csv) və embed (//go:embed).

## Əsas fikirlər

### 1. Fayl icazələri
**İcazə tipləri:** Read (r=4), Write (w=2), Execute (x=1).

**Kimlər üçün:** Owner (sahib) / Group (qrup) / Others (digərləri).

**Simvolik:** `-rw-r--r--` (birinci simvol: `-` = fayl, `d` = qovluq).

**Oktal:** 3 rəqəm = owner/group/others; 6 = 4+2 = rw-; 7 = rwx.

**Ən vacib qayda:** Go-da `0777` (oktal, baş sıfır!) ≠ `777` (decimal
511). Baş sıfır kompilyatora oktal bildirir:
```go
os.WriteFile("test.txt", data, 0644)   // 0644 — DOĞRU oktal
```

### 2. flag — Int/Var variantları
```go
// Pointer qaytaran (function variyantları):
v := flag.Int("value", -1, "Needs a value for the flag.")
flag.Parse()
fmt.Println(*v)               // dereference lazımdır

// Öz dəyişəninə yaza (Var variyantları):
var v int
flag.IntVar(&v, "value", -1, "Needs a value for the flag.")
flag.Parse()
fmt.Println(v)                // birbaşa int
```
```bash
flagapp -value=5     # 5
flagapp              # -1 (default)
```

### 3. Siqnallar (signals)
| Siqnal | Mənbə | Xassə |
|---|---|---|
| SIGINT | Ctrl+C | tutula bilər; graceful |
| SIGTERM | kill | tutula bilər; graceful |
| SIGKILL | kill -9 | TUTULA BİLMƏZ; dərhal öldürür |

**Niyə vacib:** `defer` həmişə icra olunmur — os.Exit(1), Ctrl+C
zamanı defer-lər çağırılmır! Siqnal tutub ÖZƏL cleanup aparmaq olar.

**Signal tutma pattern:**
```go
func main() {
    sigs := make(chan os.Signal, 1)        // 1-buferli siqnal kanalı
    done := make(chan struct{})             // bitmə kanalı
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTSTP)

    go func() {                             // dinləyici goroutine
        for {
            s := <-sigs
            switch s {
            case syscall.SIGINT:
                fmt.Println("My process has been interrupted. Someone might have pressed CTRL-C")
                fmt.Println("Some clean up is occuring")
                cleanUp()
                done <- struct{}{}
            case syscall.SIGTSTP:
                fmt.Println("Someone pressed CTRL-Z")
                cleanUp()
                done <- struct{}{}
            }
        }
    }()

    fmt.Println("Program is blocked until a signal is caught(ctrl-z, ctrl-c)")
    done <- struct{}{}                      // gözlə
    fmt.Println("Out of here")
}

func cleanUp() {
    fmt.Println("Simulating clean up")
    for i := 0; i <= 10; i++ {
        fmt.Println("Deleting Files.. Not really.", i)
        time.Sleep(1 * time.Second)
    }
}
```
- signal.Notify(c, sig...) — kanala göndəriləcək siqnallar
- `struct{}` tipi — sıfır yaddaşlı siqnal (yalnız "baş verdi" mənası)

### 4. Fayl yaratma və yazma
```go
// Create — touch kimi; VARSA truncate!
f, err := os.Create("test.txt")
if err != nil {
    panic(err)
}
defer f.Close()                     // həmişə bağla

f.Write([]byte("Using Write function.\n"))    // bayt slice
f.WriteString("Using WriteString function.\n") // string birbaşa

// WriteFile — bir addımda yarat+yaz:
err = os.WriteFile("test.txt", []byte("Look!"), 0644)

// OpenFile — APPEND rejimi:
f, err := os.OpenFile(
    workingFile,
    os.O_APPEND|os.O_CREATE|os.O_WRONLY,   // bayraqlar birləşdirilir
    0644,
)
if err != nil {
    return err
}
defer f.Close()
f.Write([]byte(notes))
```

### 5. Mövcudluq yoxlaması (Stat)
```go
flag.StringVar(&name, "name", "", "File name")
flag.Parse()

file, err := os.Stat(name)
if err != nil {
    if os.IsNotExist(err) {
        fmt.Printf("%s: File does not exist!\n", name)
        return
    }
    fmt.Println(err)
    return
}
// FileInfo — ad, ölçü, mode, ModTime, IsDir:
fmt.Printf("file name: %s\nIsDir: %t\nModTime: %v\nMode: %v\nSize: %d\n",
    file.Name(), file.IsDir(), file.ModTime(), file.Mode(), file.Size())
```

### 6. Fayl oxuma
```go
// Tam oxu (kiçik fayllar üçün!):
content, err := os.ReadFile("test.txt")
fmt.Println(string(content))          // []byte → string

// Bayt-bayt (böyük fayl üçün yaddaş qənaəti):
f, err := os.Open("test.txt")
if err != nil {
    log.Fatalf("unable to read file: %v", err)
}
buf := make([]byte, 1)                // bufer ölçüsü — 1 = simvol-simvol
for {
    n, err := f.Read(buf)
    if err == io.EOF {
        break
    }
    if err != nil {
        fmt.Println(err)
        continue
    }
    if n > 0 {
        fmt.Print(string(buf[:n]))
    }
}
```
- ReadFile-də EOF error YOXDUR (hamısı bir dəfə oxunur)

### 7. Fayl backup (Exercise 14.02)
```go
var ErrWorkingFileNotFound = errors.New("The working file is not found.")

func createBackup(working, backup string) error {
    // 1) Mövcudluq:
    _, err := os.Stat(working)
    if err != nil {
        if os.IsNotExist(err) {
            return ErrWorkingFileNotFound
        }
        return err
    }
    // 2) Oxu:
    workFile, err := os.Open(working)
    if err != nil {
        return err
    }
    content, err := io.ReadAll(workFile)   // *os.File = io.Reader!
    if err != nil {
        return err
    }
    // 3) Yaz:
    return os.WriteFile(backup, content, 0644)
}

func addNotes(workingFile, notes string) error {
    notes += "\n"
    f, err := os.OpenFile(workingFile,
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = f.Write([]byte(notes))
    return err
}

// main: backup → 10 qeyd əlavə et
for i := 1; i <= 10; i++ {
    note := data + " " + strconv.Itoa(i)
    addNotes(workingFile, note)
}
```

### 8. CSV
```go
in := `firstName, lastName, age
Celina, Jones, 18
Cailyn, Henderson, 13
Cayden, Smith, 42
`
r := csv.NewReader(strings.NewReader(in))   // io.Reader tələb edir
for {
    record, err := r.Read()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(record)        // []string — sütunlar!
    // record[0] = firstName
}
```
- Sətir = record ([]string); sütun = record[i]

### 9. Embedding
**Nədir:** faylları BINARY-yə daxil et — xarici asılılıq daşıma.

```
embedding_example/
├── main.go
└── templates/
    └── template.txt      # "Hello {{.Name}}"
```

```go
package main

import (
    "embed"
    "os"
    "text/template"
)

type Person struct {
    Name string
}

var (
    //go:embed templates            // DİREKTİV — dəyişəndən DƏRHAL əvvəl
    f embed.FS                       // embedded virtual filesystem
)

func main() {
    p := Person{"John"}
    tmpl, err := template.ParseFS(f, "templates/template.txt")
    if err != nil {
        panic(err)
    }
    err = tmpl.Execute(os.Stdout, p)     // "Hello John"
    if err != nil {
        panic(err)
    }
}
```
- Binary başqa yerə köçürülsə də template İŞLƏYİR — içindədir!
- Ehtiyat: böyük fayllar binary-ni şişirdir

## Əsas terminlər
- rwx / oktal (4+2+1) / 0644 baş sıfır
- Owner/Group/Others
- flag.Int vs IntVar — pointer qaytarma vs öz dəyişənə yazma
- SIGINT/SIGTERM/SIGKILL — tutula bilən/tutula bilən/qeyri-tutulan
- signal.Notify(c, sig...)
- struct{} kanalı — sıfır yaddaşlı bitmə siqnalı
- os.Create (truncate!) vs OpenFile(O_APPEND)
- os.Stat / os.IsNotExist — mövcudluq
- os.ReadFile vs f.Read(buf) — hamısı/bölməli
- io.ReadAll — Reader-dan hər şeyi oxu
- encoding/csv — Reader + []string record-lar
- //go:embed + embed.FS — binary-yə daxil etmə
- template.ParseFS — embedded FS-dən oxu

## Praktik nəticə
İcazələrdə həmişə oktal + baş sıfır (0644). Yazmaq: bir dəfəlik →
os.WriteFile; append → OpenFile(O_APPEND|O_CREATE|O_WRONLY). Mövcudluq:
os.Stat + IsNotExist. Böyük faylları ReadFile-la YOX — Read(buf) loop
və ya Scanner. Ctrl+C-də defer çağırılmır — signal.Notify + goroutine
dinləyici + done kanalı ilə öz cleanup-ını yaz. CSV üçün csv.NewReader
(io.Reader!) — record []string. Assetləri binary-yə qoymaq üçün
//go:embed + embed.FS — tək fayl paylanması.

## Mənbə
Pages: 440-469 (PDF 440-469)
