# Chapter 3 — A bookworm's digest: loops və maps (səh. 79-121)

## Bu fəsil nədən bəhs edir?

İlk "real" layihə — bookworms CLI: JSON faylından kitab həvəskarlarını yüklə,
ortaq kitabları tap, çap et. JSON formatı (dəyər tipləri; obyekt sahələri
SIRASIZ, array-lər SIRALI), testdata/ qovluq konvensiyası (go tool İGNOR
edir), iki-fayllı split (main.go = terminal/UI; bookworm.go = biznes
lojikası — reuse üçün), os.Open (yalnız READ; *os.PathError; Create 0666
truncate EDİR; OpenFile = O_APPEND|O_CREATE|O_WRONLY + xüsusi hüquqlar),
**defer** (Close() Open ilə EYNİ blokda; LIFO; Windows lock təhlükəsi;
"Open görürsənsə Close görməlisən"), encoding/json (struct tag-lər
\`json:"name"\`; sahələr EXPOSE olmalıdır — decoder yazır!; NewDecoder(f).
Decode(&bookworms) — pointer COPY YOX; io.Reader-dan oxuyur = streaming),
go run . (package-in BÜTÜN faylları; go run main.go tək fayl), test:
TestCase-də wantErr bool + "file doesn't exist" / "invalid JSON"
unhappy-path halları; global test dəyişənləri (*_test.go xaricə görünmür);
**bütün loop-lar for** (klasik; while-vari şərt; sonsuz for{}; range);
i++ l-value DEYİL (i++ < 5 YOX; ++i YOXDUR), range copy Problemi (bw
dəyişməsi slice-a TƏSİR ETMİR; 1.20-dən əvvəl eyni dəyişən!), map nəzəriyyəsi
(key comparable olmalı: slice/map/func YOX; "key1 == key2 yaza bilirikmi?"),
zero value gözəlliyi (count[book]++ — 0-dan başlayır!), slice vs array
([n]T fiks; []T = pointer+len+cap; make ölçülü; append yeni slice qaytarır
— həmişə books = append(books, b)), map-dən determinizm problemi (range
SIRA ZƏMANƏT ETMİR → sort.Slice anonymous less funksiyası ilə; UTF sırası —
Yunan hərfləri Latindən SONRA), sort.Interface (byAuthor custom tip:
Len/Swap/Less; sort.Sort IN-PLACE, copy YOX), reflect.DeepEqual (test üçün
OK, production YOX — performans), displayBooks + Example_main, bonus:
recommendOtherBooks (set map[Book]struct{} — 0 bayt value! bool YOX),
**bufio** (system call sayını azalt: NewReaderSize(f, 1MiB); Writer üçün
Flush() MÜTLƏQ — son chunk itir!).

## Əsas fikirlər

### 1. JSON + struct Decoding
```go
type Bookworm struct {
    Name  string  `json:"name"`
    Books []Book  `json:"books"`     // slice sahələri ÇOĞUL adlanır
}
type Book struct {
    Author string `json:"author"`
    Title  string `json:"title"`
}
func loadBookworms(filePath string) ([]Bookworm, error) {
    f, err := os.Open(filePath)
    if err != nil { return nil, err }
    defer f.Close()                          // Open↔Close EYNİ blokda
    var bookworms []Bookworm
    err = json.NewDecoder(f).Decode(&bookworms)   // POINTER!
    if err != nil { return nil, err }
    return bookworms, nil
}
```
- **Decoder vs Unmarshal:** Decoder io.Reader-dan OXUYUR (streaming);
  Unmarshal []byte tələb edir — böyük fayllar üçün fərq
- **EXPOSE QAYDASI:** decoder sahələrə YAZIR → kiçik hərf = həmişə BOŞ
  ("saatlarla debug" klassikası!); tag adı sahə adı ilə EYNİ olmaq ZƏRURİ
  deyil (konvensiya)
- **JSON sintaksisi:** "key":value; obyekt sahələri SIRASIZ, array SIRALI

### 2. os.Open / Create / OpenFile
| Funksiya | Rejim | Xatırlatma |
|---|---|---|
| os.Open | yalnız READ | *PathError; descriptor read-only |
| os.Create | RW 0666 | MÖVCUD FAYLI TRUNCATE EDİR! |
| os.OpenFile | istəyənə görə | O_APPEND\|O_CREATE\|O_WRONLY; xüsusi permissions |

### 3. defer — Struktur Sazişi
```go
f, err := os.Open(filePath)
if err != nil { return nil, err }
defer f.Close()      // YALNIZ uğurdan sonra! err halında f == nil
```
- **LIFO:** son defer əvvəl icra; hər return yolunda QAÇINMAZ
- **Niyə manual Close:** GC "bir gün" bağlayar; Windows-da lock =
  özünü blokla; "be polite"
- **Niyə err-check-dən SONRA:** f == nil olanda Close panic edər

### 4. Bütün Looplar = for
```go
for i := 0; i < 5; i++ { }          // klassik
for line != lastLine { }            // while-vari
for { }                              // sonsuz (return/break ilə)
for i, bw := range bookworms { }     // slice: index + COPY
for name, tc := range tests { }      // map: key + value
for _, bw := range bookworms { }     // index lazım deyil
for i := range bookworms { }         // yalnız index (copy YOX!)
```
- **i++ məhdudiyyəti:** l-value DEYİL → i++ < 5, fmt.Println(i++), ++i —
  HAMISI compile xətası; i++ yalnız STATEMENT
- **range COPY verir:** bw-ni dəyişmək slice-a təsir ETMİR (for i :=
  range + books[i] = ... lazımdırsa)

### 5. map — Unikal Açarlar Kolleksiyası
```go
count := make(map[Book]uint)      // 451 bilinirsə: make(map[Book]uint, 451)
for _, bookworm := range bookworms {
    for _, book := range bookworm.Books {
        count[book]++             // ZERO VALUE = 0 → ++ birbaşa işləyir!
    }
}
```
- **Açar comparable olmalı:** "key1 == key2 yazmaq olarmı?" — slice/map/
  func YOX; onları SAXLAYAN struct-lar da YOX ("invalid map key type")
- **Hashable:** Java-dan fərqli Hash() funksiyası LAZIM DEYİL — compiler
  bilir
- **Comma-ok + if scope:**
```go
if v, ok := mapped[3]; ok { ... }   // v yalnız if daxilində YAŞAYIR
```

### 6. slice vs array
| array [5]string | slice []string |
|---|---|
| uzunluq TİPİN bir hissəsi | pointer + len + cap |
| resize OLMAZ | dinamik |
| praktikada az istifadə | gündəlik standart |
- **append:** yeni slice QAYTARIR → həmişə `books = append(books, b)`;
  0-cap → 1-cap → böyüyərək realloksiya
- **make variantları:** make([]Book, 5) = 5 sıfır-element DOLU;
  make([]Book, 0, 5) = boş + 5 cap (append üçün ideal)

### 7. Determinizm + sort
```go
func sortBooks(books []Book) []Book {
    sort.Slice(books, func(i, j int) bool {    // anonymous less
        if books[i].Author != books[j].Author {
            return books[i].Author < books[j].Author
        }
        return books[i].Title < books[j].Title
    })
    return books
}
```
- **Problem:** map range SIRA ZƏMANƏT ETMİR → test FLAKY; həll = sort
- **UTF qeydi:** < müqayisəsi BYTE sırası — Yunan hərfləri Latindən
  həmişə SONRA
- **sort.Slice IN-PLACE:** orijinal slice DƏYİŞİR (copy yox)

### 8. sort.Interface — Obyekt Yönümlü Variant
```go
type byAuthor []Book
func (b byAuthor) Len() int           { return len(b) }
func (b byAuthor) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byAuthor) Less(i, j int) bool {
    if b[i].Author != b[j].Author { return b[i].Author < b[j].Author }
    return b[i].Title < b[j].Title
}
sort.Sort(byAuthor(books))            // tip çevirməsi ilə
```
- **Hər 3 metod MÜTLƏQ** (istifadə olunmasa belə); custom tip = sıralama
  strategiyasının ADI (byAuthor, byTitle...)

### 9. Test Helper-ləri
```go
func equalBooks(t *testing.T, books, target []Book) bool {
    t.Helper()                     // xəta sətri HELPER-i yox TEST-i göstərsin
    if len(books) != len(target) { return false }
    for i := range books {
        if books[i] != target[i] { return false }
    }
    return true
}
```
- **t.Helper() MÜTLƏQ:** fail line nömrəsi düzgün testə düşsün
- **Early exit:** uzunluq müqayisəsi əvvəl — lazımsız döngədən qaç
- **reflect.DeepEqual:** qısa amma YAVAŞ — testdə OK, productionda YOX;
  manual helper = oxunaqlı + performans (əslən emptiness/equality nil/empty
  fərqlərinə diqqət)

### 10. Unhappy-Path Test Halları
```go
tests := map[string]struct {
    bookwormsFile string
    want          []Bookworm
    wantErr       bool                 // err VAR/YOXDUR (tipi yox)
}{
    "file exists":       {..., wantErr: false},
    "file doesn't exist": {bookwormsFile: "testdata/no_file.json", want: nil, wantErr: true},
    "invalid JSON":       {bookwormsFile: "testdata/invalid.json", want: nil, wantErr: true},
}
```
- İki yoxlama: err != nil && !wantErr → FAİL; err == nil && wantErr → FAİL
- **testdata/ faylları:** no_file.json (mövcud olmayan ad), invalid.json
  (kəsilmiş — bağlanmayan mötərizə)

### 11. set — Unikal Dəyərlər İdiomu
```go
type set map[Book]struct{}          // struct{} = 0 BAYT (bool = 1 bit)
func (s set) Contains(b Book) bool {
    _, ok := s[b]
    return ok
}
// membership çox soruşulacaqsa: map[X]bool + found := myMap[element]
```

### 12. bufio — System Call Optimallaşdırması
```go
// READ: 10MiB fayl = default bufer ölçüsünə görə ÇOX syscall
buffedReader := bufio.NewReaderSize(f, 1024*1024)   // 1MiB chunks
decoder := json.NewDecoder(buffedReader)

// WRITE: son chunk bufferdə QALIR!
buffedWriter := bufio.NewWriter(f, 1024*1024)      // (NewWriter ölçüsüz də var)
for _, data := range contents {
    buffedWriter.Write(data)
}
buffedWriter.Flush()      // MÜTLƏQ — yoxsa son data İTİR!
```
- **Niyə:** syscall bahalıdır; 1MiB chunk = syscall sayı /1MiB
- **Trade-off:** 1GiB bufer = tez-tez syscall YOX amma 1GiB RAM

## Əsas terminlər
- testdata/ — go tool-un ignore etdiyi fixture qovluğu
- Struct tag (`json:"..."`) — sahə↔JSON açar xəritəsi
- Decoder vs Unmarshal — io.Reader vs []byte
- defer (LIFO) — gecikmiş icra; hər return yolunda
- *os.PathError — os əməliyyatlarının xəta tipi
- Zero value (map üçün) — count[book]++ = 0-dan başlayır
- Comparable/hashable key — slice/map/func ehtiva edənlər YOX
- slice (ptr+len+cap) vs array ([n]T fiks)
- append — yeni slice qaytaran built-in
- Determinizm — map range sırasız → sort zərurəti
- sort.Slice vs sort.Interface — callback vs custom tip (Len/Swap/Less)
- Anonymous function — adsız inline funksiya (less kimi)
- t.Helper() — fail sətrini testə yönləndirən marker
- reflect.DeepEqual — universal yavaş müqayisə (test-only)
- set (map[T]struct{}) — 0-bayt-value unikal kolleksiya
- bufio.Reader/Writer — syscall azaldan buferlər
- Flush() — Writer-dan çıxarışın SON addımı

## Praktik nəticə

1. **Fayl açma ritualı:** Open → err check → defer Close → istifadə;
   Create = TRUNCATE olduğunu UNUTMA; append = OpenFile + flags.
2. **JSON decoding şablonu:** expose edilmiş struct + tag + NewDecoder(r)
   + Decode(&v) — pointer qeydini və kiçik-hərf tələsini nəzərdən qaçırma.
3. **Deterministik nəticə üçün:** map-dən çıxan məlumatı SORT et —
   sort.Slice + anonymous less (və ya sort.Interface + adlandırılmış tip).
4. **Counter idiomu:** make(map[K]uint) + count[k]++ — zero value sayəsində
   init YOXDUR.
5. **Test helper-ləri:** t.Helper() + early-exit + manual müqayisə
   (DeepEqual-dən davamlı); global test Book-ları *_test.go daxilində
   təhlükəsizdir.
6. **Unhappy pathlər:** olmayan fayl + yanlış JSON hər loader testində.
7. **Böyük fayllar:** bufio.Reader 1MiB ilə sar; Writer-də Flush()!
8. **go run . (nöqtə!)** — çoxfayllı packaegdə tək fayl adı İŞLƏMƏZ.

## Mənbə
Pages: 79-121 (PDF 80-122)
