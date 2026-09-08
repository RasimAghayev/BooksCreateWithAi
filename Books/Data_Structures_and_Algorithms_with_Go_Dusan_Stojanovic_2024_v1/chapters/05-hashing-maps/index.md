# Chapter 5 — Hashing and Maps (səh. 248-278)

## Bu fəsil nədən bəhs edir?

Hashing prosesi: hash funksiyası seçimi (division, multiplication, mid-square,
digit folding, radix conversion, perfect, digit analysis), toqquşma (collision)
həlli (load factor, open addressing — linear/random/quadratic/double probing,
separate chaining) və maps — key-value cütlükləri saxlayan kolleksiya Go-da.

## Əsas fikirlər

### 1. Hashing nədir?
- Açarı (key) dəyərə çevirmə prosesi; dəyərlər **array (hash table)**-də
  saxlanılır, tutum = fərqli açarların sayı
- Hər dəyərin açar dəyəri ilə birbaşa əldə edilən **unikal mövqeyi** var
- **Hash funksiyası (h()):** açarı table index aralığına aid ədədə çevirən
  funksiya; ideal halda hər açar unikal mövqe verir
- **Hash table (map):** bu yanaşma ilə indekslənən array

**Toqquşma (collision):** 2+ açar eyni mövqeyə düşür → bu açarlar
**synonyms** (sinonimlər); sinonimlər dəsti — **equivalence class**
(ekvivalentlik sinfi)

### 2. Hash funksiyası — meyarlar
1. **Mümkün qədər sadə** — hər çıxışda icra olunur, hesablama sürətli olmalıdır
2. **Mümkün qədər vahid (uniform)** — toqquşmalardan qaçınmaq üçün

**Açar paylanmasından asılı OLMAYAN metodlar** (paylanma əvvəlcədən məlum
deyil, yalnız bəzi əsas xassələr bilinir):

| Metod | Formula / Mexanizm | Nümunə |
|---|---|---|
| **Division** | h(k) = k mod n (n ≤ table ölçüsü) | h(27) = 27 mod 10 = 7; n cüt olsa cüt açarlar cüt indexlərə düşür — yalnız yarısı yüklənir! |
| **Multiplication** | h(k) = floor(n * (c·k mod 1)), 0<c<1 | h(27) = floor(10 * 0.986) = 9; c qızıl nisbətə yaxın seçilməli |
| **Mid-square** | k² -nin ORTASINDAN lazımi rəqəm sayı götürülür | 100 slot üçün 2 rəqəm, 1000 üçün 3 |
| **Digit folding** | k parçalara bölünür (son parça qısa ola bilər), parçalar toplanır; table-dan böyükdürsə mod | 27051989 → parçalar cəmi |
| **Radix conversion** | k, p əsasında ədəm kimi qəbul edilir, q əsasına (q>p) "çevrilir" | 275₁₀ (q=12) → 2·12²+7·12+5 = 377; 100 slot üçün son 2 rəqəm: 77 |
| **Perfect** | minimal funksiya, toqquşma MÜMKÜN DEYİL; n açarı n unikal indexə | tapılması çətindir |

**Açar paylanmasından asılı metodlar** (açar dəsti əvvəlcədən məlumdur):
- **Digit analysis ( heuristic)** — ən populyar:
  1. Bütün açar dəyərlərinin rəqəmsal təsvirləri təhlil edilir
  2. Hər rəqəm mövqeyində rəqəmlərin görünmə sayı cədvəli qurulur
  3. İndeksləməyə lazımi qədər mövqe seçilir — **rəqəm müxtəlifliyi ən az
     dəyişən** sütunlar seçilir

- Perfect-dən başqa BÜTÜN metodlar toqquşma yarada bilər

### 3. Toqquşma həlli
**Ən sadə yol:** table ölçüsünü artırmaq (daha çox slot = daha az toqquşma
ehtimalı)

**Load factor (α):** hash table-in doluluq dərəcəsi
```
α = element sayı / table ölçüsü
```

**2 standart üsul: open addressing + separate chaining**

**Open addressing:** toqquşma olanda açar üçün başqa boş ünvan axtarılır.
Hər açar üçün **probe sequence** (ünvanlar ardıcıllığı) yaradılır; insert
zamanı boş slot tapılanadək nəzərdən keçirilir. Yeni ünvanın hesablanması =
**rehashing**:

| Üsul | Formulası | Qeyd |
|---|---|---|
| **Linear probing** | h(k, i) = (h(k) + i) mod n | növbəti boş slot; table sonuna çatanda başdan davam |
| **Random probing** | 0..n-1 arası pseudo-random ardıcıllıq | probe sequence kimi istifadə |
| **Quadratic probing** | h(k, i) = (h(k) + i²) mod n | addımlar kvadratik artır: h=2 nümunə: 2, 3, 6, 11, 18 |
| **Double hashing** | h(k, i) = (h1(k) + i·h2(k)) mod n | 2 müstəqil hash funksiyası |

Nümunə (linear probing): 27 və 21 eyni h() dəyəri verir → 21 növbəti boş
slota düşür.

**Separate chaining:** sinonimlər **list**-ə zəncirlənir; list başları hash
table slot-larındadır. Nümunə: 5 slot, h(k) = k mod 5.

### 4. Maps
- **Key → value** xidməti göstərən data struktur; **key-value cütlükləri**
  kolleksiyası
- Hər açar map-də YALNIZ BİR DƏFƏ rast gəlinir
- Hashing üçün istifadə olunur: developer mapping qaydalarını (hash
  funksiyasını) müəyyən edir; toqquşma olanda saxlanan dəyər (adətən) yeni
  dəyərlə ƏVƏZ OLUNUR

**Əməliyyatlar (hamısı açarla):**
1. Dəyərin daxil edilməsi (insert)
2. Mövcud dəyərin yenilənməsi (update)
3. Dəyərin alınması (get)
4. Dəyərin silinməsi (remove)
5. Açarın mövcudluğunun yoxlanması (contains)

### 5. Maps in Go
- Go daxili data strukturudur; **zero value = nil map** (açar ƏLAVƏ OLUNA
  BİLMƏZ, oxuma sıfır dəyər verir)
```go
var m map[string]int          // nil map
m := make(map[string]int)      // istifadəyə hazırlanmış map
var m = map[string]int{        // initializer
    "Monday": 1, "Tuesday": 2, "Wednesday": 3,
    "Thursday": 4, "Friday": 5, "Saturday": 6, "Sunday": 7,
}
```
- `[keyType]valueType`; struct custom tipləri də açar/dəyər ola bilər

**Insert / Update (bir operatorla):**
```go
m["Monday"] = 0   // açar varsa update, yoxdursa insert
```

**Get:**
```go
day := m["Monday"]
// açar yoxdursa element tipinin ZERO dəyəri qaytarılır (int → 0)
```

**Delete:**
```go
delete(m, "Thursday")
```

**Açar mövcudluq yoxlaması — "comma ok" idiomu:**
```go
day, ok := m["Monday"]   // ok=true açar varsa; yoxdursa false + zero dəyər
```

**İterasiya (for range):**
```go
for key, element := range m {
    fmt.Println(key, element)   // ilk dəyər = key, ikinci = element
}
```

## Termindirmə (AZ)
- Hashing — Hashing (açarı dəyərə çevirmə prosesi)
- Hash function — Hash funksiyası
- Hash table — Hash cədvəli
- Collision — Toqquşma (eyn indexə düşən açarlar)
- Synonyms — Sinonimlər (eyn mövqeyə xəritələnən açarlar)
- Equivalence class — Ekvivalentlik sinfi (sinonimlər dəsti)
- Load factor (α) — Yükləmə amili (element sayı / table ölçüsü)
- Open addressing — Açıq ünvanlama (boş slot axtarışı)
- Probe sequence — Zond ardıcıllığı (açara aid ünvanlar massivi)
- Rehashing — Yenidən hash-ləmə (yeni ünvan hesablanması)
- Linear probing — Xətti zondlama
- Quadratic probing — Kvadratik zondlama
- Double hashing — İkili hash-ləmə
- Separate chaining — Ayrı zəncirləmə (synonym-lər list-də)
- Division / Multiplication / Mid-square / Digit folding / Radix conversion /
  Perfect hash — bölmə / vurma / orta-kvadrat / rəqəm qatlaması / əsas
  çevirmə / mükəmməl hash metodları
- Digit analysis — Rəqəm təhlili (paylanmadan asılı heuristic)
- Golden ratio — Qızıl nisbət (multiplication metodunda c)
- Map — Xəritə (key-value kolleksiyası)
- "comma ok" idiom — açar mövcudluğu yoxlaması (day, ok := m[k])

## Kviz sualları
1. Hash funksiya hansı 2 meyarla qiymətləndirilir? (sadəlik — hər çıxışda
   hesablanır; vahidlik — toqquşma az)
2. Division metodunda n cüt seçilsə nə baş verir? (cüt açarlar cüt
   indexlərə — table-in yalnız yarısı yüklənir)
3. Multiplication-da c niyə qızıl nisbətə yaxın olmalıdır? (yaxın paylanma —
   toqquşmaları azaldır)
4. Load factor necə hesablanır? (α = element sayı / table ölçüsü)
5. Quadratic probing-in probe sequence-i nümunə (h=2): (2, 3, 6, 11, 18 —
   h(k)+i² addımlarla)
6. Separate chaining-də sinonimlər harada saxlanılır? (slot-da başlanan
   list-də zəncirlənir)
7. nil map-ə açar əlavə etmək olarmı? (xeyr — nil map açar qəbul etmir;
   make lazımdır)
8. `v, ok := m[k]` nə verir? (açar yoxdursa ok=false, v=zero dəyər)
