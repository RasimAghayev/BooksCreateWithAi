# Know Go — Müəllim Qeydləri (Azərbaycanca)

📖 **Kitab deyir:**

Generics Go-ya 2022-də (1.18) gəldi və kitab onu interfeyslərdən başlayaraq
təkamül yolu ilə izah edir: konkret tip → interfeys → any → type parameter.
Mərkəzi ideya: constraint = tip dəsti, operator = zəmanət mübadiləsi.
Kitab Set/Stack kimi container-lər, onların concurrency-safe versiyaları və
slices/maps/cmp/iter paketləri ilə praktiki tərəfi göstərir. Sonda generics-in
həll etmədiyi şeylər (option type, enum, makroslar) və "Go 2 olmayacaq"
mövzusunda realçı mövqe saxlayır. Performans sualına: runtime-da fərq YOX,
interface/reflection əvəz olunsa 2x qədər QAZANC mümkündür.

👨‍🏫 **Müəllim qeydi:**

Mənim praktikada gördüyüm ən çox generics səhvi — yersiz istifadədir.
Kitab da "move cautiously" deyir, amma vurğulayım: 90% hallda interface
sufisentdir. Generic yazmağın LEGİTIM səbəbləri:

1. **Container/utility kitabxanası** (slices paketi dəqiq belə doğulub) —
   eyni alqoritm N tip üçün.
2. **Performance-kritik yolda interface indirection-in aradan qaldırılması**
   (btree nümunəsi 2x).
3. **Type safety:** `(T, error)` qaytaran funksiyalar any/assertionsız.

Generics yazmağın QADAĞAN olmalı olduğu hallar: bir dəfə istifadə olunacaq
kod, interface artıq işləyirsə, "cool" olsun deyə. Arundel-in "cow and
chicken" (Farm[Animal]) nümunəsi dərsin əsas prizmasıdır: constraint TİP
DEYİL — class sistemi refleksindən qorun.

Sərbəst tövsiyələr:

- `~[]E` idiomunu sadəcə slice qaytaran funksiyalarda YOX, İCTİMAİ API-lərdə
  istifadə et — daxili helperlərdə `[]E` oxunaqlılıq qazandırır.
- `-race` bayrağını CI-da DEFAULT et — kitabın "smoke test" dərsi müasir
  Go praktikasında əskiksizdir.
- Iterator API layihələşdirəndə Seq2[value, error] patternini nəzərdən
  keçir — Go 1.23-dən sonra "slice qaytarıram" deyə API bağlama.
- Kitabın məşq reposunu gərçəkdən yazın (GOAL/HINT/SOLUTION formatı
  bilik möhkəmlənməsinin ən effektiv yolu).

## Ən vacib 5 fikir

1. **Constraint = zəmanət mübadiləsi:** operator istəyirsənsə tip dəstini
   daralt — "bigger the interface, weaker the abstraction" (Ch3).
2. **Runtime-da generic YOXDUR:** hər şey compile-da konkret tipə
   instantiate olunur — bu, performans (indirection YOX) və []any fərqinin
   açarıdır (Ch2, Ch10).
3. **~ approximation derived tiplər üçün:** [S ~[]E, E any] — Go-nun öz
   slices paketi belə imzalanır (Ch4).
4. **Concurrency safety kompozisiyadır:** bir dəfə mutex qoy — bütün
   metodlar (String, Union — All/Add üzərindən) avtomatik safe alır (Ch8).
5. **Abstraksiya PULSUZ deyil:** Griesemer-in "needless abstraction
   introduces complexity" sitatı — generics də istisna deyil (Ch10).

## Kitabın ən dəyərli hissəsi

Chapter 3 (Constraints) + Chapter 4 (Operations) — səh. 44-90. Bu 47 səhifə
generics-in SEMANTİK nüvəsidir: type set/union/intersection/approximation
və operator-constraint mübadiləsi. Əksər "generics bilirəm" deyənlərin
zəif nöqtəsi məhz budur. Chapter 8-in race detector dərsi (smoke test +
-race) isə təkcə generics üçün deyil, hər konkurrent Go kodu üçün əvəzolunmaz
praktikadır.

## Kitabdakı püşkələr haqqında

Kitabın sonunda müəllif öz kitabları/mentorluğu barədə yazır (For the Love
of Go, The Power of Go: Tools) — bilavasitə texniki məzmun deyil, amma
"Know Go" seriyasın "Know Go: Generics" kimi bir hissə kimi anlamağa kömək
edir. Bu kitab 210 səhifədir və interfeys bilmədiyi fərziyyəsi ilə
başlayır — ardıcıllıq sırf başlanğıcdan L3-ə qədər. Tələbə üçün ideal yol:
bütün GOAL məşqlərini REAL yazmaq — həllərə baxmadan.
