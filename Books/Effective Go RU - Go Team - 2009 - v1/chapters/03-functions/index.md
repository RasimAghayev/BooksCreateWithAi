# Chapter 3 — Funksiyalar (səh. 10-13)

## Bu chapter nədən bəhs edir?

Go funksiyalarının 3 idiomatik özəlliyi: çoxlu qaytarma dəyərləri (C-nin EOF/-1 və pointer-out idiomlarından qurtuluş), adlandırılmış nəticə parametrləri (sənədləşdirmə + naked return) və defer (resurs idarəsi, LIFO, arqument dərhal qiymətlənir, trace/recover patternləri).

---

## əsas fikirlər

### 1. Çoxlu qaytarma dəyərləri

**Problem (C-də):** xəta = xüsusi dəyər (-1 EOF) və ya pointer-out parametr. **Go həlli:** `Write` həm sayı həm xətanı qaytarır:

```go
func (file *File) Write(b []byte) (n int, err error)
// "Bəzi baytlar yazıldı, amma hamısı YOX, çünki disk doldu"
```

**Pointer-simulyasiyasız birləşmə:** funksiya ədədi və növbəti mövqeyi qaytarır:

```go
func nextInt(b []byte, i int) (int, int) {
    for ; i < len(b) && !isDigit(b[i]); i++ {
    }
    x := 0
    for ; i < len(b) && isDigit(b[i]); i++ {
        x = x*10 + int(b[i]) - '0'
    }
    return x, i
}

for i := 0; i < len(b); {
    x, i = nextInt(b, i)   // i yenilənir — pointer YOX
    fmt.Println(x)
}
```

### 2. Adlandırılmış nəticə parametrləri

Nəticə parametrləri ad ala bilər → funksiya başında zero value ilə initializə; `return` (naked) cari dəyərləri qaytarır. Adlar **sənədləşdirmə** rolunu oynayır:

```go
func nextInt(b []byte, pos int) (value, nextPos int)
```

**io.ReadFull — naked return idiomu:**

```go
func ReadFull(r Reader, buf []byte) (n int, err error) {
    for len(buf) > 0 && err == nil {
        var nr int
        nr, err = r.Read(buf)
        n += nr
        buf = buf[nr:]
    }
    return    // n və err cari dəyərləri
}
```

### 3. Defer (təxirə salınmış çağırış)

Defer funksiyanın (əhatəedən) bitməsində icra olunur — **resursların hər exit path-da azad edilməsi** üçün:

```go
// Contents faylın məzmununu string kimi qaytarır.
func Contents(filename string) (string, error) {
    f, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()    // hər return-da bağlanacaq — unutmaq mümkünsüz!
    var result []byte
    buf := make([]byte, 100)
    for {
        n, err := f.Read(buf[0:])
        result = append(result, buf[0:n]...)
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err    // f burada da bağlanır
        }
    }
    return string(result), nil    // və burada da
}
```

**2 üstünlük:** (1) close-u unutmaq mümkünsüz — yeni return yolu əlavə olunsa belə; (2) close-un açılışın YANINDA olması — kod sonda yerləşən qədər oxunaqlıdır.

**Vacib qayda — arqumentlər defer ZAMANINDA qiymətlənir** (çağırış anında YOX):

```go
for i := 0; i < 5; i++ {
    defer fmt.Printf("%d ", i)
}
// Çap: 4 3 2 1 0  (LIFO sıra!)
```

**Trace pattern** — arqumentin dərhal qiymətlənməsindən istifadə:

```go
func trace(s string) string {
    fmt.Println("entering:", s)
    return s
}
func un(s string) { fmt.Println("leaving:", s) }

func a() {
    defer un(trace("a"))    // trace İNDİ çağrılır, un sonra!
    fmt.Println("in a")
}
func b() {
    defer un(trace("b"))
    fmt.Println("in b")
    a()
}
func main() { b() }
// entering: b / in b / entering: a / in a / leaving: a / leaving: b
```

Defer **bloq deyil, FUNKSİYA səviyyəsində işləyir** — ən güclü tətbiqləri panic/recover ilə (ch9-da).

---

## Əsas terminlər

| Termin | İzah |
|--------|------|
| Çoxlu qaytarma | `(int, error)` — C pointer-out/sentinel idiomlarının əvəzi |
| Named result parameter | Adlı nəticə — zero value başlanğıc + sənədləşmə |
| Naked return | Arqumentsiz `return` — cari named dəyərlər |
| defer | Əhatəedən funksiyanın return-unda icra |
| LIFO | Defer-lərin son-giren-ilk-çıxan sırası |
| Dərhal arqument qiymətlənməsi | Defer anında arqumentlər hesablanır — icrada YOX |
| Trace pattern | `defer un(trace("a"))` — giriş/çıxış cütü |

---

## Praktik nəticə

1. **Çoxlu qaytarma C-idiomlarını öldürür:** -1 EOF yox, `(n, err)`; pointer-out yox, dəyər qaytar.
2. **Nəticə adları imzanın sənədidir:** `(value, nextPos int)` — hansı int nədir aydın; naked return qısa funksiyalarda.
3. **Resurs açılanda defer bağla:** `defer f.Close()` açılışın yanında — hər return path-i avtomatik qorunur.
4. **Defer arqumenti dərhal bağlanır:** dəyişən cari dəyəri istəyirsən — closure/ref-pattern (müasir Go-da dəyişən tut).
5. **LIFO sırası yadda saxla** — tərs resurs sifarişi.
6. **Funksiya-səviyyəli abstraksiya:** block-scope dillərindən fərqli — bu defer-in GÜCÜDÜR.

---

## Mənbə

- Sənəd: *Effective Go* (rus tərcüməsi), 2009
- PDF səhifələri: 10-13
