# Chapter 1 — Fundamentals of Data Structures and Algorithms (səh. 34-84)

## Bu fəsil nədən bəhs edir?

Data strukturlarının və alqoritmlərin əsas konseptləri: data strukturunun
tərifi, xarakteristikaları (linear/qeyri-linear, statik/dinamik, homogen/
heterogen), yaddaş təmsili (sequential/linked), pointer-lər, Go struct-ları,
alqoritm tarixi, təsnifatları, O-notasiyası və Go funksiyaları.

## Əsas fikirlər

### 1. Data struktur nədir?
- **Data structure (data strukturu):** datanın təşkilinin təsviri
- Nümunə: 2D nöqtə (x, y koordinatları); ünvan strukturu (küçə, şəhər,
  poçt kodu)
- Hər data struktur üzərində **spesifik əməliyyatlar** icra olunur
  (məs. array-də sort)
- Məşhur strukturlar: arrays, graphs, maps — kitab boyu hamısı əhatə olunacaq

### 2. Data strukturlarının xarakteristikaları
**Elementlərin münasibətinə görə:**
| Sinif | Təsvir | Nümunə |
|---|---|---|
| Linear (xətti) | hər element yalnız 2 qonşusu ilə əlaqədə | array, list, stack, queue |
| Non-linear (qeyri-xətti) | bir element çox elementlə əlaqədə | tree, graph |

**Ölçü dəyişkənliyinə görə:**
| Sinif | Təsvir |
|---|---|
| Static (statik) | ölçü proqram icrası zamanı dəyişmir; yaddaş istifadəsi çox vaxt optimal deyil |
| Dynamic (dinamik) | icra zamanı ölçü artır/azalır |

**Element tipinə görə:**
| Sinif | Təsvir |
|---|---|
| Homogenous (homogen) | bütün elementlər eyni tipdə (array) |
| Heterogenous (heterogen) | elementlər fərqli tiplər (Person struct: string ad + int yaş) |

### 3. Yaddaş təmsili
**Sequential (ardıcıl):**
- Elementlər **fasiləsiz** yaddaş sahəsində bir-birinin ardınca
- Fiziki və məntiqi sıra EYNİDİR
- Bir element bir neçə memory location tutə bilər (32-bit element + 16-bit
  location ⇒ 2 location; 5 elementlik array ⇒ 10 location)

**Linked (zəncirvari):**
- Elementlər yaddaşın təsadüfi yerlərində, **qeyri-fasiləsiz**
- Fiziki ≠ məntiqi sıra; elementlər **pointer**-lə birləşir
- Hər elementə pointer sahəsi əlavə olunduğundan daha çox yer tutur

**Pointer-lər (Go):**
```go
var pi *int   // int pointer tərifi
i := 27
pi = &i       // & — ünvan operatoru
*pi = 18      // * — göstərilən dəyəri oxu/yaz
fmt.Println(*pi)
```

### 4. Go struct-ları
```go
type Point struct {
    X int
    Y int
}

var (
    p1 = Point{27, 5}   // bütün sahələr
    p2 = Point{X: 18}   // Y ötürülür → zero value (0)
    p3 = Point{}        // hər ikisi ötürülür → {0, 0}
)

p := Point{27, 5}
p.X = 18                // . operatoru ilə giriş
pp := &p                // struct pointer
```
- Ötürülməyən sahələrə **default (zero) value** təyin olunur
- Məşhur strukturlar (array, map) Go-da daxili tip kimi; qalanları üçün
  custom struct yazılır

### 5. Alqoritmlərin tarixi
- **Muhammad Al-Khwarizmi** (fars riyaziyyatçısı): "alqoritm" termini onun
  adından; hind rəqəmlərini və onluq sistemi ərəb riyaziyyatına gətirdi
- XII əsr latın tərcüməsi "Algorismi de numero indorum" — adın yanlış
  tərcüməsi "algorithm" terminini doğurdu
- **Ada Lovelace** (1842): maşın üçün ilk alqoritm — Bernoulli ədədləri,
  Babbage-in analitik maşını üçün → ilk proqramçı; Ada dili onun şərəfinə
- **Alan Turing:** alqoritmləri riyaziyyata gətirdi; hər yaxşı tərif olunmuş
  alqoritm **Turing maşını**nda icra oluna bilər (sonsuz lent + simvollar +
  qaydalar cədvəli — hesablamanın riyazi modeli)

### 6. Alqoritm nədir?
- Giriş dəyər(lər)ini çıxış dəyər(lər)inə çevirən **proses**
- Düzgün alqoritm: hər giriş üçün etibarlı çıxış. Nümunə (yanlış alqoritm):
  x² + x + 41 → 39-dan böyük tam ədədlərdə bəzən mürəkkəb ədəd verir
- **Proof of correctness:** predikat hesabı ilə teorem kimi sübut;
  **invariant** — alqoritm icrası boyu dəyişməyən ifadə; sonlu vaxtda
  bitmə yoxlanılır

Düzgün alqoritmin xassələri:
1. Diskret əməliyyatlar ayrı-ayrı addımlarla nəticəyə aparır
2. Sonlu sayda addımdan sonra çıxış verir
3. Eyni giriş → həmişə eyni çıxış (kvadratın tərəfi 2 → sahə həmişə 4)
4. Yalnız müəyyən giriş dəsti üçün deyil, bütün tiplər üçün tətbiq olunur

### 7. Alqoritmlərin təmsili
| Üsul | Xüsusiyyət |
|---|---|
| **Pseudocode** | təbii dil + proqramlaşdırma konvensiyaları (loop, if, funksiya); mahiyyətə fokus; istənilən dilə asan çevrilir |
| **Flowchart** | qrafik təmsil; icra axını izləmək üçün əla; mürəkkəb alqoritmlərdə qarışıq |
| Proqram dili | real implementasiya (bu kitabda Go) |

Flowchart blokları: **rectangle** = proses, **diamond** = qərar,
paralleloqram = giriş/çıxış, oval = baş/aşaqı...

### 8. Təsnifatlar
**Dizayn paradiqmasına görə:**
| Paradiqma | Mahiyyət | Nümunə |
|---|---|---|
| Divide and conquer | problem eyni tipin kiçik problemlərinə bölünür | binary sort |
| Dynamic programming | alt-problemlərin optimal həllərindən optimal həll qurulur (üst-üstə düşən alt-problemlər) | ən qısa yol |
| Greedy | hər addımda "o an ən yaxşı" seçim (həmişə optimal deyil) | Kruskal (MST) |
| Linear programming | xətti funksiyanın maks/min-i xətti şərtlər altında | — |
| Randomized | həll üçün təsadüfülik istifadə olunur | — |
| Genetic | təkamül prosesini imitasiya edir | — |
| Heuristic | optimala yaxın həll; optimal həll qeyri-mümkün olanda (məs. yaddaş çatışmazlığı) | — |

**Implementasiyaya görə:**
| Sinif | Təsvir |
|---|---|
| Iterative | loop-larla; hər recursion iterativ həllə çevrilə bilər (çox vaxt daha mürəkkəb) |
| Recursive | özünü çağırır (şərt ödənilənə qədər); nümunə: Fibonacci |
| Sequential | əməliyyatlar bir-birinin ardınca; köhnə tək-CPU kompüterlər üçün |
| Parallel | alt-problemlər paralel icra olunur, nəticələr birləşdirilir; çox-CPU |
| Distributed | paralel kimidir, amma şəbəkə ilə birləşən bir neçə kompüterdə |

**Sahəyə görə:** search, sorting, merge, numerical, graph, string,
combinatorial, computational geometry, machine learning, cryptography,
compression, parsing

### 9. Komplekslik və O-notasiyası
- **Space** (yaddaş) + **time** (CPU vaxtı) — iki resurs
- **Big O** — "order of approximation" (yaxınlaşma dərəcəsi)
- Artan komplekslik sırası:
```
O(1) < O(log n) < O(n) < O(n log n) < O(n^k), k > 2 (polinom) < O(k^n), k > 2 (eksponensial)
```

| Sinif | Xassə | Nümunə |
|---|---|---|
| O(1) constant | girişdən asılı deyil; ən arzuolunan, ən nadir | listin əvvəlinə element əlavəsi |
| O(log n) logarithmic | yavaş artım; problemi trivial qalana qədər kiçilt | sıralı array-də binary search |
| O(n) linear | n dəfə icra olunan loop | sıralanmamış array-də sequential search |
| O(n log n) log-linear | ən çox rast gəlinən; böl + hamısını emal et | merge sort |
| O(n²) square | iki iç-içə loop | birbaşa sort üsulları (bubble/selection/insertion) |
| O(n^k) polynomial | k iç-içə loop | square onun xüsusi halıdır |
| O(k^n) exponential | arzuolunmaz; sürətli artım | exhaustive search |

- Seçim varsa **ən aşağı komplekslikli** alqoritm götürülməlidir

### 10. Go funksiyaları
```go
func hello() {                    // arqumentsiz, qaytarmır
    fmt.Println("Hello World")
}
func inc(i int) int {             // bir arqument
    return i + 1
}
func sum(i, j int) int {          // eyni tip → sonuncuda yazılır
    return i + j
}
func calc(i int) (int, int) {     // çoxlu nəticə
    return i*i, i+i
}
func inc(i int) (res int) {       // adlandırılmış nəticə (naked return)
    res = i + 1
    return
}
func inc(i *int) {                // pointer arqumenti
    *i = *i + 1                   // return-siz dəyişiklik görünür
}
```
- Arqument tipi addan SONRA; eyni tip qruplaşdırıla bilər
- Naked return: adlı nəticələr funksiya başında təyin olunmuş kimi davranır
- Pointer arqumenti: funksiya daxilində dəyişiklik çağırıcıda görünür

### 11. Data strukturları + alqoritmlər
- Proqramın 2 əsas bina bloğu; biri digərisiz mənasızdır
- Ümumiyyətlə data strukturu alqoritmdən ƏVVƏL gəlir (data olmadan emal yoxdur)
- Struktur və alqoritmin mürəkkəbliyi bir-birindən asılı deyil: sadə struktur
  üzərində mürəkkəb alqoritm, və əksinə mümkündür
- Düzgün seçim proqramın effektivliyini müəyyən edir

## Termindirmə (AZ)
- Data Structure — Data Strukturu (datanın təşkili təsviri)
- Linear / Non-linear — Xətti / Qeyri-xətti
- Static / Dynamic — Statik / Dinamik (ölçü dəyişkənliyi)
- Homogenous / Heterogenous — Homogen / Heterogen
- Sequential / Linked representation — Ardıcıl / Zəncirvari yaddaş təmsili
- Pointer — Göstərici (başqa dəyişənin yaddaş ünvanını saxlayan dəyişən)
- Algorithm — Alqoritm (girişi çıxışa çevirən proses)
- Invariant — İnvariant (icra boyu dəyişməyən ifadə)
- Pseudocode — Pseudokod
- Flowchart — Axın diaqramı
- Proof of correctness — Düzgünlüyün sübutu
- Turing machine — Turing maşını
- Big O notation — Böyük O notasiyası
- Divide and Conquer — Böl və Hökm Et
- Dynamic Programming — Dinamik Proqramlaşdırma
- Greedy — Acgöz (seçim)
- Heuristic — Evristik (optimala yaxın)
- Naked return — Çılpaq qaytarış (adlı nəticələrlə)

## Kviz sualları
1. Array hansı 3 xarakteristikaya görə təsnif olunur? (linear, statik —
   ölçü dəyişmir, homogen elementlər)
2. Linked representation niyə daha çox yer tutur? (pointer sahəsi datadan
   əlavə saxlanılır)
3. Ada Lovelace-in alqoritmi nə üçün ilk idi? (maşın üçün yazılmış ilk
   alqoritm — Bernoulli ədədləri, Babbage analitik maşını)
4. O(1), O(log n), O(n log n), O(n²) — hansı sırayla artır? (məhz bu sıra;
   constant < log < linear < log-linear < square < exponential)
5. Hər recursive alqoritm nəyə çevrilə bilər? (iterativ — çox vaxt daha
   mürəkkəb kodla)
6. `func inc(i, j int) int` — nə qaytarır? (i və j arqumentləri, int nəticə;
   qaytarış dəyəri funksiya bədənindən asılıdır)
