# Chapter 6 — Pointers (Göstəricilər)

## Bu chapter nədən bəhs edir?

Pointer sintaksisi və semantikası, Java/Python class-ları ilə müqayisəsi, pointer-lərin
mutable parametr göstəricisi kimi istifadəsi, map/slice-in pointer təbiəti, buffer
pattern-i, stack/heap, escape analysis və garbage collector iş yükünün azaldılması.

## Əsas fikirlər

### 1. Pointer Praymer
**Nədir:** Başqa dəyişənin yaddaş ünvanını saxlayan dəyişən.

**Necə işləyir:**
- `&x` — address operator (ünvan); `*p` — indirection/dereference (göstərilən dəyər)
- Hər pointer eyni ölçüdədir (ünvan sayı); zero value — `nil` (C-dəki NULL-dan fərqli:
  0-la conversion YOXDUR; universe block-da olduğundan kölgələnə bilər — `nil` adını
  dəyişən kimi istifadə ETMƏ).
- Nil pointer-in dereference-i **panic**.
- `new(T)` — T-nin zero value-suna pointer qaytarır (az işlənir); struct üçün `&Foo{}`
  işlət. Primitiv literal/constant-a `&` apply oluna BİLMƏZ (`&"Perry"` — compile xətası).

**Kitabdan kod nümunəsi:**
```go
x := 10
pointerToX := &x
fmt.Println(pointerToX)  // ünvan
fmt.Println(*pointerToX) // 10
z := 5 + *pointerToX     // 15
var x *int               // nil — *x panic edir
```

**Pointer-to-primitive helper pattern (struct sahələri üçün):**
```go
func stringp(s string) *string {
    return &s   // parametr dəyişəndir → ünvanı var
}
p := person{MiddleName: stringp("Perry")}
```

### 2. Pointerlərdən Qorxma — Bu, Sənin Tanıdığın Davranışdır
**Nədir:** Java/JS/Python/Ruby class instansiyaları = pointer davranışı; Go-da bu,
yalnız struct-lar üçün seçimdir.

**Müqayisə fərqləri:** Həmin dillərdə class sahə dəyişikliyi çağırana görünür,
parametr-in yenidən təyinatı yox, null parametr doldurula bilmir — Go-da pointer ilə
DAÇEYİ davranış. O dillər də pass-by-value-dir — kopyalanan şey pointer-in özüdür.
Fərq: Go həm primitiv, həm struct üçün value/pointer SEÇİMİ verir.

**Default seçim:** **Value işlət** — data axını daha aydın, GC işi azdır. Pointer
yalnız səbəb olduqda.

### 3. Pointer = Mutable Parametr Siqnalı
**Nədir:** Go-da immutable bəyanı yoxdur; pointer parametr "bu funksiya dəyəri dəyişə
bilər" mesajıdır. Nonpointer parametr — funksiya orijinalı dəyişə BİLMİR (kopya alır).

**İki vacib nüans:**
1. **Nil pointer-i funksiya daxilində non-nil etmək mümkün deyil** — kopyalanan ünvan
   dəyişdirilə bilməz:
```go
func failedUpdate(g *int) { x := 10; g = &x }  // f main-də nil qalır
```
2. **Dəyəri dəyişmək üçün dereference lazımdır**:
```go
func failedUpdate(px *int) { x2 := 20; px = &x2 }  // x dəyişmir
func update(px *int) { *px = 20 }                    // x = 20 ✓
```

### 4. Pointer — Son Çarə
**Qaydalar:**
- Struct-u doldurmaq üçün pointer parametr YOX — funksiya struct **qaytarsın**:
```go
// YOX:
func MakeFoo(f *Foo) error { f.Field1 = "val"; ... }
// BELƏ:
func MakeFoo() (Foo, error) { f := Foo{Field1: "val", Field2: 20}; return f, nil }
```
- İstisna: interface gözləyən funksiyalar — `json.Unmarshal(data, &f)` (generics
  olmadığından tip ötürülməsi belə aparılır). Bu istisnadır, ümumi qayda DEYİL.
- Qaytarmada value type üstünlük; pointer qaytar — yalnız daxili state modifikasiyası
  lazımdırsa (I/O buffer-ləri, concurrency tipləri).

**Performans rəqəmləri (i7-8700, 32GB):** Pointer ötürməsi ~1ns (sabit); 10MB value
ötürməsi ~1ms. Qaytarmada əks nəticə: <1MB data üçün pointer qaytarmaq DAHA YAVAŞ
(100 bayt: value ~10ns, pointer ~30ns — heap escape + GC xərci); >1MB-da pointer
qalib gəlir (10MB: value ~2ms, pointer ~0.5ms). Praktikada fərq əhəmiyyətsizdir —
yalnız megabaytlarla işləyəndə pointer düşün.

### 5. Zero Value Yoxdursa Pointer (comma-ok üstünlüyü)
**Nədir:** "0 təyin edilib" ilə "heç təyin edilməyib" fərqi.

**Qaydalar:** Bu fərq önəmlidirsə nil pointer istifadə et; amma funksiyadan nil pointer
yerinə **comma-ok idiomu** qaytar (value + bool). JSON istisna: nullable sahələr üçün
pointer field-lər doğrudur (`*string` sahə — "yoxdur" vs boş string fərqi).

### 6. Map və Slice-in Pointer Təbiəti
**Map:** runtime-da pointer-to-struct kimi implement olunub → kopya üzərindən dəyişikliklər
orijinala görünür. **API-də map parametr/qaytarmağa QADAĞA yaxın**: (1) məzmunu
self-documenting deyil — kod izləmədən başqa yolu yox; (2) nə qoyulduğunu izləmək
mümkün deyil. Güclü tip dili olan Go-da map yerinə **struct** işlət.

**Slice:** 3 sahəli struct — `len` (int), `cap` (int), data pointer. Kopyada len/cap/pointer
kopyalanır: element dəyişikliyi ortaq yaddaşda → görünür; `append` (uzunluq dəyişikliyi)
kopyanın metadata-sında → görünMÜR (orijinal len-i keçən dəyərlər "gizli" qalır; capacity
yetərsə olsa yeni blok ayrılır). Bu səbəbdən funksiyaya ötürülən slice-in **məzmunu dəyişə
bilər, ölçüsü dəyişə bilməz**. Konvensiya: funksiya slice-i dəyişmirsə docs-da yaz.

**Niyə slice hər ölçüdə funksiyaya keçir:** ötürülən data sabitdir (2 int + pointer);
array-da isə bütün data kopyalanır — ölçü tipin hissəsi olduğu üçün hər ölçü ayrı tip.

### 7. Buffer kimi Slice
**Nədir:** Xarici mənbədən oxuyarkən hər iterasiyada yeni yaddaş ayırmamaq üçün bir
dəfə yaradılan slice.

**Kitabdan kod nümunəsi:**
```go
file, err := os.Open(fileName)
if err != nil { return err }
defer file.Close()
data := make([]byte, 100)     // buffer — bir dəfə ayrılır
for {
    count, err := file.Read(data)
    if err != nil { return err }
    if count == 0 { return nil }
    process(data[:count])     // yalnız dolu hissə
}
```

**Sub-kod izahı:**
- `make([]byte, 100)` — nonzero len buffer (funksiyaya kəsərək ötürülən oxu buffer üçün
  düzgün formadır)
- `data[:count]` — oxunmuş hissə; `count == 0` — EOF

### 8. GC İş Yükünü Azaltma — Stack, Heap, Escape Analysis
**Anlayışlar:**
- **Garbage** — heç bir pointer-in işarə etmədiyi data; GC-nin işi onu tapıb yaddaşı
  geri qaytarmaqdır.
- **Stack** — funksiya çağırışlarının paylaşılan ardıcıl blok; ayırmaq = stack pointer-i
  sürüşdürmək (çox sürətli); funksiya çıxanda avtomatik deallokasiya. Go-nun xüsusiyyəti:
  hər goroutine-in öz stack-i var və runtime iş vaxtında stack-i BÖYÜDƏ bilər (kiçik
  başlayır — amma böyümə zamanı kopyalama yavaşdır).
- **Heap** — GC-nin idarə etdiyi yaddaş. Pointer-in göstərdiyi data stack-də qala bilməzsə
  "escape" edir və heap-də yerləşir. Escape şərtləri: ölçü compile vaxtı bilinmir /
  pointer funksiyadan qaytarılır / təhlükəsizlik üçün kompilyator konservativdır.
  (C-dən fərqli: lokala pointer qaytarmaq Go-da leqaldır — data heap-ə keçir; C-də
  dangling pointer bug-u olardı.)
- Value tipləri (primitiv, array, struct) — ölçüsü compile vaxtı bilinir → stack-də.

**Performans arqumentləri:**
- GC hər dövrə < 500µs hədəfləyir (low latency, Jeff Dean "The Tail at Scale"
  tövsiyəsinə uyğun; Rick Hudson ISMM 2018 talk). Çox zibil → dövrələr uzanır, yaddaş artır.
- **Mechanical sympathy** (Martin Thompson, 2011): RAM ardıcıl oxunur — `[]struct`
  datanı ardıcıl saxlayır; `[]*struct` isə RAM boyu səpələyir (~100x yavaş, Forrest
  Smith ölçümləri). Java/Python-da hər obyekt heap-dədir; Go-da idiomatik yol həm də
  ən effektivli yol: az pointer, az zibil, struct-ları ardıcıl.

## Əsas terminlər
- Dereference (dəyərə keçid) — `*p` ilə pointer-in hədəfini oxumaq
- Address operator (ünvan operatoru) — `&x` ilə ünvan almaq
- Escape analysis (çıxış analizi) — pointer-in datasının stack/heap qərarı
- Mechanical sympathy (mexaniki uyğunluq) — hardware-i nəzərə alan kod yazma
- Buffer (tampon) — təkrar istifadə üçün bir dəfə ayrılmış slice
- Comma-ok idiom — value + bool qaytarma (nil pointer əvəzinə)

## Praktik nəticə

Pointer qərar ağacı: (1) default — value; (2) funksiya parametri dəyişməlidirsə — pointer
(oxu koduna "mutable" siqnalı); (3) struct dolduran funksiya — pointer götürməsin,
struct qaytarsın; (4) API-da map yox, struct; (5) nullable anlayış üçün pointer yerinə
comma-ok; JSON nullable sahə — istisna; (6) oxu dövrələrində buffer slice; (7) böyük
struct-ların `[]T`-si — RAM-ə sığışan ən sürətli forma. Nil pointer-i funksiya daxilində
"diriltmək" mümkün deyil — məsələ mütləq kənarda həll olunmalıdır.

## Mənbə
Pages: 159-188
