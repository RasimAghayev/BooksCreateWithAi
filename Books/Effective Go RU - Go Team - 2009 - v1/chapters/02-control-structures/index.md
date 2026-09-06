# Chapter 2 — İdarə konstruksiyaları (səh. 6-10)

## Bu chapter nədən bəhs edir?

Go-nun if/for/switch/type switch konstruksiyalarının C-dən fərqləri: init ifadəsi, else buraxma stili, `:=` təkrar mənimsətmə qaydası, range loop xüsusiyyətləri (rune parse), fallthrough-unsuz switch, label break/continue və type switch idiomu.

---

## Əsas fikirlər

### 1. Ümumi fərqlər (C-dən)

- `do`/`while` YOX — yalnız ümumiləşdirilmiş `for`
- `switch` daha çevik; `if`/`switch` init ifadəsi qəbul edir (for kimi)
- `break`/`continue` + **label** — hansı konstruksiyanın bitdiyini seçə bilər
- Yeni: **type switch**, **select**
- Mötərizə YOX; `{ }` məcburi

### 2. If — init + else buraxma

```go
if x > 0 {
    return y
}
// init ifadəsi ilə lokal dəyişən:
if err := file.Chmod(0664); err != nil {
    log.Print(err)
    return err
}
```

**Kitabxana stili:** if gövdəsi `break/continue/goto/return` ilə bitirsə — **else BURAX**; xəta halları səhifə boyu təmizlənir, uğur axını aşağı axır:

```go
f, err := os.Open(name)
if err != nil {
    return err
}
d, err := f.Stat()
if err != nil {
    f.Close()
    return err
}
codeUsing(f, d)
```

### 3. `:=` — təkrar mənimsətmə (redeclaration)

`d, err := f.Stat()` — err artıq elan olunmuşdu! Bu **qanunidir**: `:=` mənimsətmə şərtləri:
1. Dəyişən EYNİ scope-da artıq mövcuddursa → YENİDƏN ELAN OLUNMUR, sadəcə mənimsədilir (xarici scope-da isə YENİ yaradılır)
2. Dəyər təyin edilə bilən tiptədir
3. Ən azı bir DİGƏR dəyişən yeni yaradılır

Bu praqmatizm — uzun if-else zəncirində tək `err` dəyişəni. (Funksiya parametrlərinin scope-u gövdə ilə eynidir.)

### 4. For — 3 forma + range

```go
for init; condition; post { }   // C-for
for condition { }                // C-while
for { }                          // for(;;)
```

**Range** — array/slice/string/map/channel üzrə:

```go
for key, value := range oldMap { newMap[key] = value }
for key := range m { if key.expired() { delete(m, key) } }   // yalnız key
for _, value := range array { sum += value }                   // yalnız value (_)
```

**String range — UTF-8 parse:** hər iterasiya 1 RUNE verir; yanlış enkodiya 1 bayt + U+FFFD (dəyişdirmə simvolu):

```go
for pos, char := range "日本\x80語" {
    fmt.Printf("character %#U starts at byte position %d\n", char, pos)
}
// U+65E5 '日' @ 0 · U+672C '本' @ 3 · U+FFFD @ 6 · U+8A9E '語' @ 7
```

**Vergül operatoru YOX**; `++/--` ifadə deyil → çoxdəyişənli idarə **paralel mənimsətmə** ilə:

```go
// massivi çevir:
for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
    a[i], a[j] = a[j], a[i]
}
```

### 5. Switch — qiymətsiz, case-sıralı, yarıömürlü

- ifadə KONSTANT/İNTEGER OLMAQ MƏCBURİ DEYİL
- case-lər YUXARIDAN AŞAĞI yoxlanılır
- ifadəsiz switch = `switch true` — **if-else-if zəncirinin idiomu**:

```go
func unhex(c byte) byte {
    switch {
    case '0' <= c && c <= '9':
        return c - '0'
    case 'a' <= c && c <= 'f':
        return c - 'a' + 10
    case 'A' <= c && c <= 'F':
        return c - 'A' + 10
    }
    return 0
}
```

- **Fallthrough YOX** (avtomatik); vergüllə çox dəyər: `case ' ', '?', '&', '=', '#', '+', '%':`
- switch-i erkən bitirmək üçün `break`; LOOP-u bitirmək üçün **label**:

```go
Loop:
    for n := 0; n < len(src); n += size {
        switch {
        case src[n] < sizeOne:
            if validateOnly {
                break            // switch-i qırır
            }
            size = 1
            update(src[n])
        case src[n] < sizeTwo:
            if n+1 >= len(src) {
                err = errShortInput
                break Loop       // FOR-u qırır!
            }
            // ...
        }
    }
```

- `continue` də label qəbul edir (yalnız looplar üçün).

### 6. Type switch — interfeysin dinamik tipi

```go
var t interface{}
t = functionOfSomeType()
switch t := t.(type) {
default:
    fmt.Printf("неожиданный тип %T\n", t)
case bool:
    fmt.Printf("логическое %t\n", t)    // t burada bool
case int:
    fmt.Printf("целое %d\n", t)          // t burada int
case *bool:
    fmt.Printf("указатель на bool %t\n", *t)
case *int:
    fmt.Printf("указатель на int %d\n", *t)
}
```

İdiom: hər case-də eyni adda YENİ dəyişən — hər case-də FƏRQLİ TİPLƏ.

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Init ifadəsi | if/switch/for-da `;`-dan əvvəlki lokal elan |
| Else buraxma | Xəta return-lə bitirsə else lazım deyil — uğur axını aşağı |
| Redeclaration (`:=`) | Eyni scope-da mövcud dəyişən + yeni dəyişən → mənimsətmə |
| Range rune parse | String range-də hər iterasiya = 1 rune (U+FFFD xətalı üçün) |
| Paralel mənimsətmə | `i, j = i+1, j-1` — vergül operatorunu əvəz edir |
| Qiymətsiz switch | if-else-if zənciri kimi işləyən idiom |
| Label break/continue | `break Loop` — nested switch-dən loop qırma |
| Type switch | `switch t := t.(type)` — dinamik tipə görə budaqlanma |

---

## Praktik nəticə

1. **Xəta yoxlamalarında else YOX** — `if err != nil { return }` + davam aşağı; oxunuş səhifə boyu düz axır.
2. **`:=` err təkrarı qanunidir** — eyni scope-da; yeni dəyişən də olmalıdır.
3. **İf-else-if → qiymətsiz switch** — daha təmiz budaqlanma.
4. **Fallthrough lazımdırsa** — case-lərə vergüllə çoxlu dəyər ver.
5. **Loop içində switch-dən loop qırma:** `Loop:` label + `break Loop`.
6. **Type switch-də eyni ad** — hər case öz tipi ilə `t`.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi), 2009
- PDF səhifələri: 6-10
