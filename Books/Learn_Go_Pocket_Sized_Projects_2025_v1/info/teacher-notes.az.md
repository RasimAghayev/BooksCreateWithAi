# Learn Go with Pocket-Sized Projects — Müəllim Qeydləri (AZ)

📖 Kitab deyir: 12 kiçik layihə ilə Go-nun demək olar bütün əsaslarını —
testdən gRPC-yə, generic-lərdən TinyGo-ya qədər — öyrə.

👨‍🏫 Müəllim qeydi: bu, gördüyüm ən balanslı "Go-ya giriş" kitablarından
biridir, amma 3 praktik qeyd:

1. **`crypto/rand` vs `math/rand` (Ch5):** kitab `pickWord` üçün crypto/rand
   göstərir — tədris baxımından doğrudur (determinizm problemindən qaçır),
   amma oxucu başa düşməlidir: oyun üçün `math/rand` və ya `rand/v2` kifayətdir
   və DAHA SÜRƏTLİDİR. Fərq "təhlükəsizlik ehtiyacı" sualındadır, sürədə deyil.

2. **Loop dəyişəni (Ch9):** kitab Go 1.22-dən əvvəlki davranışı düzgün izah
   edir və `p := p` fəndini göstərir. Amma 2025-dən etibarən Go 1.22+ həmin
   kölgəni AVTOMATİK edir — yeni kodda bu sətir artıq lazım deyil. Amma köhnə
   modullarla işləyəndə bilmək lazımdır. Kitabın yanaşması "hər halda yaz,
   zərər verməz" — mən də buna qoşuluram (backward compatibility).

3. **Xəta API-nin görünməsi (Ch10):** kitab typed error + errors.As göstərir —
   dəqiqdir, amma real layihələrdə 90% halda SADƏCƏ `errors.Is` + wrap kifayətdir.
   Typed error yalnız zəngin kontekst lazım olanda (API cavabı, DB idarəsi)
   qiymətlidir. Yeni başlayanlar bunu "hər yerdə typed error" kimi oxumamalıdır.

## Ən vacib 5 fikir

1. TDD gözlənilən davranışı SƏNİNLƏ yaşayır — refactoring cəsarəti verir
2. Xətalar = dəyərlər: wrap + sentinel + (bəzən) typed
3. Asılılıqlar kiçik interfeyslərlə injekt — io.Writer-dan gRPC client-ə
4. Paralellik: `go test -race` olmadan heç bir paralel kod tamamlanmayıb sayılır
5. Məzmunu domain-də saxla, protokolları (JSON/proto) kənarda tut

## Kitabın ən dəyərli hissəsi

Chapter 4 (pocketlog) — kitabxana YAZMAQ prosesini (API sabitliyi, doc.go,
external test) bir dəfəyə öyrədir; bu, tək-tək layihələrdə görünməyən
" craftsmenlik" tərəfini göstərir.

⚠️ Uyğunsuzluq yoxdur — kitabın kodu və mətni tutarlıdır; yalnız versiya
zamani (Go 1.22/1.23 loop dəyişəni, `any` alias) qeydləri aktualdır.
