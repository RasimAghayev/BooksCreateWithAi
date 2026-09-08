# 100 Go Mistakes — Müəllim Qeydləri (AZ)

📖 Kitab deyir: 100 səhv kateqoriyaları ilə — bug-dan optimallaşdırmaya qədər;
hər səhv konteksti ilə, "her yerde tətbiq olunan" dogma yox.

👨‍🏫 Müəllim qeydi: bu kitab Go-nun "ikinci dərəcəli biliyi"nin ən yaxşı
toplusudur; amma oxucu üçün 3 qeyd:

1. **Səhv nömrələrinə görə deyil, KÖKƏ görə öyrəş:** #30 (range kopyası),
   #32 (pointer tələsi), #63 (loop closure) — hamısı EYNİ kökdür: element
   kopyası + referens semantikası. Bu 3-ü birlikdə "slice iteration modeli"
   dərsi kimi təqdim etmək daha effektivdir. Kitab onları ayrı-ayrı verir
   (2 fəsildə yayılmış) — tələbə üçün sintez lazımdır.

2. **Optimallaşdırma fəsli "hazır komanda" deyil, "hazır olma" dərsidir:**
   #91-94 (cache/ILP/alignment) mikro-səviyyədə nadir hallarda lazım olur —
   real dünyada #95-98 (escape, pool, pprof, GC) 90% dəyəri verir. Yeni
   oxucu 91-94-ü "mən bunu hər gün yazmalıyam" kimi oxumamalıdır; benchmark
   GÖSTƏRMƏDƏN heç biri tətbiq olunmaz.

3. **Beyin səhvlərdən öyrənir (giriş tezisi) — amma yalnız SƏN ETDİYİN
   səhvlərdən:** kitabı xəttli oxumaqdansa, öz kodunda `go vet`, `-race`,
   `staticcheck` ilə taramaq + tapılan səhvlərə görə kitabda müvafiq bölməni
   oxumaq daha effektiv işləyir. Kitabın özü də linter siyahısı ilə bitirir
   (#16) — bu, metodoloji nöqtədir.

## Ən vacib 5 fikir

1. Görünməz semantika: range kopyası, map qaydasızlığı, select randomluğu
2. Yaddaş: backing paylaşımı, capacity, escape — statik analizlə görünən
3. Paralellik: lifecycle + race + scheduler bilgisi bir paketdir
4. Xətalar: wrap + Is/As + bir dəfə handle
5. Ölçmə → dəyiş: benchmark+benchstat+pprof zənciri

## Kitabın ən dəyərli hissəsi

Chapter 9 (Concurrency: Practice) — ən çox istehsal bugının mənbəyi olan
fəsil; #62 (lifecycle) və #66 (nil kanal) tək başına kitabı dəyərli edir.

⚠️ Uyğunsuzluq: kitabın bölgüsündə #17/#18 (overflow/float) fəsil
sərhədləri TOC-da fərqli göstərilir (2nd ed. yenidən qruplaşdırma) — mətndə
data-tiplər fəslindədir; nömrə ardıcıllığı tamamdır, yalnız bölmə sayları
fərqlidir. Oxucu çaşmasın deyə biz mətn təqib etdik.
