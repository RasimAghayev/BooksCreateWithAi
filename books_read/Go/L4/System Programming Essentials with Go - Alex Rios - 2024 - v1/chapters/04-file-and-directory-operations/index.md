# Chapter 4 — File and Directory Operations (Fayl və Kataloq Əməliyyatları)

## Bu chapter nədən bəhs edir?
Fayl icazələrinin təhlilinə (os.Stat/FileMode/Perm), path/filepath paketinə (Join/Clean/Split/WalkDir), simvolik keçidlərə (Symlink/Readlink/broken-link detector), kataloq ölçüsü hesabına, dublikat fayl tapıcısına (MD5 hash) və mmap ilə fayl optimizasiyasına.

## Əsas fikirlər

### 1. os.Stat — metadata oxunuşu
```go
func Stat(name string) (FileInfo, error)

info, err := os.Stat("example.txt")
if err != nil {
    if os.IsNotExist(err) {       // mövcud deyil ayrıca ele alınmalı
        fmt.Println("File does not exist")
    } else { panic(err) }
}
fmt.Printf("Name: %s, Size: %d, Mode: %s, ModTime: %s",
    info.Name(), info.Size(), info.Mode(), info.ModTime())
```
**FileInfo interfeysi:** Name(), Size(), Mode(), ModTime() + IsDir().

### 2. Linux fayl tipləri ↔ FileMode bitləri
| Simvol | Tip | Go biti |
|--------|-----|---------|
| `-` | Regular file | (bit yoxluğu; IsRegular()) |
| `d` | Directory | os.ModeDir; IsDir() |
| `l` | Symbolic link | os.ModeSymlink (mode&os.ModeSymlink != 0) |
| `p` | Named pipe (FIFO) | os.ModeNamedPipe |
| `c` | Character device | os.ModeCharDevice |
| `b` | Block device | (birbaşa bit yoxdur) |
| `s` | Socket | os.ModeSocket |

### 3. İcazələr (permissions) sistemi
**RWX üçlük:** owner / group / others — hər biri read(4) + write(2) + execute(1).
`-rw-r--r--` = 644; `drwxr-xr-x` = **755** (oktal — Unix-in köhnə konvensiyası, onilliklərdir qorunur).

```go
fileInfo, _ := os.Stat("example.txt")
permissions := fileInfo.Mode().Perm()        // yalnız icazə bitləri
fmt.Sprintf("%o", permissions)               // "644" kimi oktal string
```

### 4. path/filepath paketi
**Path separator problemi:** Unix `/`, Windows `\` — filepath platformadan asılılığı aradan qaldırır.

```go
filepath.Join("/home/user", "document.txt")    // → /home/user/document.txt
filepath.Clean("/home/user/../documents/file.txt") // → /home/documents/file.txt
dir, file := filepath.Split("/home/user/documents/myfile.txt")
// dir: /home/user/documents/ | file: myfile.txt
```

### 5. filepath.WalkDir — kataloq ağacı gəzintisi
```go
func WalkDir(root string, fn fs.WalkDirFunc) error
type WalkDirFunc func(path string, d DirEntry, err error) error
```
**DirEntry 4 metodu:** Name() (base name), IsDir(), Type() (FileMode tiplərinin alt çoxluğu), Info() (FileInfo; dəyişmiş/silinmiş fayllarda ErrNotExist qaytara bilər; symlink-in ÖZÜ haqqında məlumat — hədəfi yox!).

**Nəticə semantikası:**
| Callback qaytarır | Nə olur |
|---|---|
| `filepath.SkipDir` | bu kataloqu ötür |
| `filepath.SkipAll` | hamısını ötür, dayan |
| `nil` | davam et |
| xəta | WalkDir tam dayanır, xətanı qaytarır |

**Praktik nümunə — kataloq siyahılayıcı (CLI genişlənməsi):**
```go
flag.StringVar(&outputFileName, "f", "", "Output file (default: stdout)")
flag.Parse()

// MultiWriter: həm stdout, həm fayl:
if cfg.OutputFile != "" {
    outputFile, err := os.Create(cfg.OutputFile)
    defer outputFile.Close()
    outputWriter = io.MultiWriter(cfg.OutStream, outputFile)
} else {
    outputWriter = cfg.OutStream
}

for _, directory := range directories {
    err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
        if path == ".git" {
            return filepath.SkipDir        // .git-i ötür
        }
        if d.IsDir() {
            fmt.Fprintf(outputWriter, "%s\n", path)
        }
        return nil
    })
    if err != nil { /* log + continue */ }
}
```
**io.MultiWriter:** bir yazı — bir neçə hədəfə (stdout + fayl eyni anda).

### 6. Simvolik keçidlər (symlink)
**Yaratma (CLI):**
```bash
ln -s /home/user/documents/important_document.txt /home/user/desktop/shortcut.txt
```
**Yaratma (Go):**
```go
err := os.Symlink(sourcePath, symlinkPath)
// ls -l: lrwxrwxrwx ... shortcut -> source (l = symlink)
```

**Silmə (unlink):**
```bash
unlink /home/user/desktop/shortcut_to_document.txt
# və ya: rm /home/user/desktop/shortcut_to_document.txt
```
```go
err := os.Remove(filePath)    // link-i silir — HƏDƏF FAYLA TOXUNMUR
```

**Broken symlink detector CLI:**
```go
if info.Mode()&os.ModeSymlink != 0 {          // symlink-dirmi?
    target, err := os.Readlink(path)          // hədəfi oxu
    if err != nil { /* read xətası */ } else {
        _, err := os.Stat(target)              // hədəf mövcudmu?
        if err != nil && os.IsNotExist(err) {
            fmt.Fprintf(outputWriter, "Broken symlink found: %s -> %s\n", path, target)
        }
    }
}
```
**Link/unlink fəlsəfəsi:** link = yeni ad əlaqəsi; unlink = adın silinməsi (son ad silinəndə fayl da "yox olur" — inode referens sayğacı).

### 7. Kataloq ölçüsü hesabı
```go
func calculateDirSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
        if err != nil { return err }
        if !fileInfo.IsDir() {
            size += fileInfo.Size()          // yalnız FAYL ölçüləri cəmlənir
        }
        return nil
    })
    return size, err
}

// İnsan-oxunarlı vahid çevirmə:
switch {
case size < 1024:            unit = "B"
case size < 1024*1024:       size /= 1024; unit = "KB"
case size < 1024*1024*1024:  size /= 1024*1024; unit = "MB"
default:                     size /= 1024*1024*1024; unit = "GB"
}
```

### 8. Dublikat fayl tapıcısı — MD5 hash
```go
func computeFileHash(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil { return "", err }
    defer file.Close()
    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil { return "", err }  // fayl → hash
    return fmt.Sprintf("%x", hash.Sum(nil)), nil                     // hex string
}

func findDuplicateFiles(rootDir string) (map[string][]string, error) {
    duplicates := make(map[string][]string)
    err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
        if err != nil { return err }
        if !info.IsDir() {
            hash, err := computeFileHash(path)
            if err != nil { return err }
            duplicates[hash] = append(duplicates[hash], path)   // hash → fayl siyahısı
        }
        return nil
    })
    return duplicates, err
}

// Nümayiş: yalnız >1 faylı olan hash-lər:
for hash, files := range duplicates {
    if len(files) > 1 {
        fmt.Printf("Duplicate Hash: %s\n", hash)
        for _, file := range files { fmt.Fprintln(outputWriter, "  -", file) }
    }
}
```

### 9. mmap — memory-mapped files
**Problem:** Yaddaş həddini aşan böyük fayllar.
**Həll:** Faylı birbaşa yaddaşa xəritələ — OS disk yazılarını idarə edir, proqram yaddaşdakı data ilə işləyir.

```go
file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
defer file.Close()

fileInfo, _ := file.Stat()
fileSize := fileInfo.Size()

data, err := syscall.Mmap(int(file.Fd()), 0, int(fileSize),
    syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
if err != nil { /* handle */ }
defer syscall.Munmap(data)          // mütləq təmizlə!

// Yaddaşda birbaşa redaktə:
newContent := []byte("Hello, mmap!")
copy(data, newContent)               // dəyişiklik fayla əks olunur (MAP_SHARED)
```
**Arqumentlər:**
- `int(file.Fd())` — fayl deskriptoru
- `0` — offset (faylın başlanğıcı)
- `int(fileSize)` — xəritələnən uzunluq
- `PROT_READ|PROT_WRITE` — qoruma rejimləri
- `MAP_SHARED` — dəyişikliklər fayla sinxronlaşır

**OOM təhlükəsizliyi:** File-backed+shared mapping 64-bit mühitdə OOM killer-dən azaddır; 32-bitdə yalnız ünvan sahəsi limiti (mmap graseful fail edir, killer yox).

## Əsas terminlər
- os.Stat / FileInfo — metadata interfeysi
- os.IsNotExist — "yoxdur" xətasının ayrılması
- FileMode / Mode.Perm() — tip bitləri / icazə bitləri
- RWX üçlüyü (owner/group/others), oktal (755)
- path/filepath: Join/Clean/Split/WalkDir
- DirEntry vs FileInfo — gəzinti vəziyyəti
- SkipDir / SkipAll — gəzinti nəzarəti
- io.MultiWriter — çoxhədəfli yazı
- os.Symlink / os.Readlink / os.Remove
- Broken/dangling symlink — hədəfsiz keçid
- MD5 hash / io.Copy(hash, file)
- mmap / Munmap, PROT_READ/WRITE, MAP_SHARED
- File-backed vs anonymous mapping
- OOM killer

## Praktik nəticə
1. Stat xətalarını böl: IsNotExist ayrıca — "yoxdur" normal ssenaridir, panic deyil.
2. Kataloq gəzintisində WalkDir (DirEntry — stat-sız) > Walk (FileInfo); .git kimi qovluqları SkipDir ilə ötür.
3. Symlink yoxlaması: mode&ModeSymlink + Readlink + Stat(target) üçlüyü — broken-link təmizləyici skriptlərin əsası.
4. Dublikat axtarışı = hash map qur + len(files)>1 filtri; MD5 dedup üçün kifayət, kriptoqrafik bütövlük üçün SHA-256.
5. Böyük fayl redaktəsində mmap: read/write syscall-ları aradan qalxır, OS virtual yaddaşla idarə edir; Munmap defer ilə.
6. CLI output-a `-f` flag + io.MultiWriter — istifadəçi ekranda VƏ faylda eyni nəticəni görür.

## Mənbə
Pages: 61-81 (PDF səh. 82-103)
