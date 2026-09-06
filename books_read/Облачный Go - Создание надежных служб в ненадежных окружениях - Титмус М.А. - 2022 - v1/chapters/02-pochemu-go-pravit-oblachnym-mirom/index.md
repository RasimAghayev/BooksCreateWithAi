# Chapter 2 — Почему Go правит облачным миром

## Bu chapter nədən bəhs edir?

Go dilinin bulud hesablamaları üçün nəyə görə xüsusi seçildiyini izah edir. 2007-ci ildə Google-da Robert Griesemer, Rob Pike və Ken Thompson tərəfindən yaradılan Go, dövrün server dillərinin qüsurlarını (gec kompilasiya, çox qarmaşıq OOP, yaddaş təhlükəsizliyi problemləri) həll etmək üçün dizayn edilib. Bu chapter konseptualdir, kod nümunəsi yoxdur.

## Əsas fikirlər

### 1. Go-nun yaranma tarixi
**Nədir:** 2007-ci il, Google. O dövrün dilləri (C++, Java, Python və s.) paylılıq (concurrency), yaddaş təhlükəsizliyi və sürətli kompilasiya tələblərini qarşılamırdı.

**Necə işləyir:** Üç senior engineer (Griesemer, Pike, Thompson) yeni bir dil yaratdılar. Məqsəd: kompilyasiya sürətli, concurrency daxili, yaddaş təhlükəsiz, sadə sintaksis.

### 2. Bulud dünyasının xüsusiyyətləri
**Nədir:** Bulud tətbiqləri serverlər, şəbəkələr və paralel işləmə tələb edir. Əvvəlki dillər bu tələbləri qarşılamırdı.

**Çatışmamazlıqları:**
- Qarmaşıq kod, çox oxunmaz
- Gec kompilasiya (saatlarla, böyük klasterlərdə belə)
- Yaddaş idarəsi problemləri (C/C++-da malloc/free)
- Versiya uyğunsuzluğu və bahalı yeniləmələr

### 3. Kompozisiya irsiyyət əvəzinə (Composition over Inheritance)
**Nədir:** Go-da klassik OOP inheritance (miras) yoxdur. Əvəzinə kiçik tiplər böyük tiplərə daxil edilir (embedding).

**Necə işləyir:** `Car has Engine` münasibəti. `Engine` struktur `Car`-ə daxil edilir. Interface-lər davranış müqaviləsidir, tip xassələri deyil.

**Üstünlükləri:**
- Inheritance ağacı ("diamond problem") yoxdur
- `Shape` interface → `Rectangle` və `Circle` onu avtomatik implement edir, `implements` elan etməyə ehtiyac yoxdur

### 4. Sadəlik və qüvvə (Simplicity and Power)
**Nədir:** Go kodunun oxunması, yoxlanması və saxlanması asandır.

**Necə işləyir:** Az, amma çox faydalı standart kitabxana. Görüşməli (readable) kod. "Clever" kod yerinə aydın kod üstün tutulur.

### 5. CSP modeli və Goroutines
**Nədir:** Communicating Sequential Processes (CSP) — Tony Hoare tərəfindən təklif olunmuş model.

**Necə işləyir:** Goroutine-lər və kanallar (channels) vasitəsilə mesajlaşma. Shared memory və mütləq məxaric (mutex) əvəzinə "data göndər, döyüşdürmə" modeli.

**Konkurensiya vs Parallelizm:**
- **Konkurensiya:** Müstəqil proseslərin birlikdə irəliləməsi (zamanlanması müəyyən deyil)
- **Parallelizm:** Eyni anda bir neçə prosesin işləməsi

### 6. Sürətli kompilasiya
**Nədir:** Go kompilyatoru bir neçə saniyədə böyük layihələri kompilasiya edir.

**Necə işləyir:** Sadə asılılıq resolüsiya modeli. C/C++-dakı kimi include faylları və dərin asılılıq zəncirləri yoxdur.

**Mərhələ:** 1.8 milyon sətirli Kubernetes v1.20.2, MacBook Pro (Intel i9, 8 core, 32 GB RAM) → **45 saniyə**.

### 7. Dilin sabitliyi (Stability)
**Nədir:** Go 1 mart 2012-də buraxıldı. Go 1 proqramları gələcək Go 1.x buraxışlarında işləməyə davam edəcək.

**Necə işləyir:** Backward compatibility zəmanəti. Komanda kritik dəyişiklikləri `go fix` vasitəsilə avtomatik yeniləyir.

### 8. Yaddaş təhlükəsizliyi (Memory Safety)
**Nədir:** Manual memory management (malloc/free), pointer aritmetikası və buffer overflow yoxdur.

**Necə işləyir:** Garbage Collector (GC) avtomatik yaddaş idarə edir. Null pointer-dən bəzi hallar xaricində, Go runtime panik edir və proqramı dayandırır.

### 9. Performans
**Nədir:** Kompilyasiya dilinə baxən dinamik dillər (Python, Ruby) ilə müqayisədə yüksək performans.

**Nəyə lazımdır:** Bulud tətbiqləri tez-tez yük dalğaları (traffic spikes) alır — vertikal miqyaslama qənaət etmək üçün yüksək performans lazımdır.

### 10. Statik kompozisiya (Static Linking)
**Nədir:** Go binary-ləri bütün lazımi kitabxanaları özünə daxil edən tək executable fayldır.

**Necə işləyir:** ~2 MB ölçülü "Hello World" binary yaradır, amma heç bir xarici runtime və ya kitabxana tələb etmir.

**Üstünlükləri:**
- Konteynerlərdə (Docker) ideal — base image-lərdən asılılıq yoxdur
- Conflict riski yoxdur

### 11. Statik tip sistemi (Static Typing)
**Nədir:** Dəyişənlərin tipləri kompilasiya vaxtı yoxlanılır.

**Necə işləyir:** Python-da `my_varaible = my_variable + 1` (typo) kimi xətalar kompilasiya vaxtında yaxalanır. Dynamic dillərdə isə bu runtime-da, bəzən production-da aşkar olunur.

**Üstünlükləri:**
- Bug-ların erkən aşkar edilməsi
- Daha az debug vaxtı
- Kod özü sənəddir (self-documenting)

## Əsas terminlər
- CSP (Communicating Sequential Processes)
- Goroutines (qorutinlər)
- Channels (kanallar)
- Composition (kompozisiya)
- Structural typing (struktural tipləmə)
- Duck typing (ördək tipi)
- Static linking (statik bağlantı)
- Memory safety (yaddaş təhlükəsizliği)
- Garbage Collector (GC / zibil toplayıcı)

## Praktik nəticə
Go seçimi "sürətli yaz" (Python) vs "sürətli icra" (C++) dilemma-sını həll edir: kompilasiya sürətli, concurrency daxili, yaddaş təhlükəsiz və sabit API ilə. Bulud dünyasında bu dördü eyni vaxtda tələb olunur.

## Mənbə
Pages: 36-46 (PDF səh. 36-46)
