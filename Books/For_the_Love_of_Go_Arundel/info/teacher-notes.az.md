# For the Love of Go — Müəllim Qeydləri (Azərbaycanca)

## 📖 Kitab deyir:

Go-nu TEST YAZARAQ öyrət: hər konsepsiya TDD tsikli ilə təqdim olunur — test
yaz → compile xətası gör → null implementation → FAIL → real kod → PASS.
Kitabın mərkəzi prinsipləri:

- "We simply can't write untestable functions when the test comes first"
- Null implementation: "Heç bilmədik — testin bug TUTA BİLDİYİNİ sübut etmədən
  ona inanma"
- "One behaviour, one test" — funksiya yox, davranış
- "Making it impossible to compile incorrect programs is the best kind of
  validation!"
- "Test behaviours, not functions"
- "The test is the definition of what's okay"
- "Left-align the happy path" — minimum indent
- Tao: kindness, simplicity, humility, not striving

## 👨‍🏫 Müəllim qeydi:

Kitab başlanğıc səviyyəsi üçün ƏN YAXŞI strukturlardan birinə malikdir
(təhsil dizaynı baxımından): kiçik addımlar, dərhal tətbiq, real ssenari,
compiler xətalarının "tərcüməsi". Tədrisdə vurğulanacaqlar:

1. **Tələbələr GOAL məşqlərini ATLAMAMALI.** Kitabın dəyəri nəzəri hissədə
   deyil — əl çəkmə məşqlərindədir. Hər fəsildə kodu KÖÇÜRüb İŞLƏDİN,
   sonra özün yaz.

2. **Null implementation anlayışı tələbələr üçün KOUNTER-İNTUİTİVDİR.**
   "Səhv kod yaz ki, testi yoxlayasan" — gəliri çətin başa düşür. Nümayiş
   edin: TestMultiply-ı `return 0` ilə KEÇİRİN → sonra düzgün Multiply
   yazın → YENİ KEÇİR. Fərqi soruşun: "Hansı halda testə İNAM edirsən?"

3. **Pointer-receiver tələsi (ch10) kitabın ƏN vacib praktik dərsidir.**
   SetPriceCents value-receiver ilə SƏSSİCƏ uğursuz olur (compile OLUR,
   amma dəyişməz!) — pass-by-value anlayışı üçün MÜKƏMMƏL təlimat. Bu
   fəsli tələbələr öz BAŞINA yaşasınlar; sonra izah edin.

4. **Map sırasının randomlığı (ch7) təcrübi göstərilMƏLİDİR:** GetAllBooks
   testini 10 dəfə ardıcıl işə salın — yarı fail görəcəklər. Bu, "flaky
   test" anlayışına ƏN yaxşı giriş qapısıdır.

5. **Closure loop tələsi (ch13: 3,3,3) müzakirə ilə gəlsin:** tələbələrdən
   NİYƏ soruşun; sonra "closure ÇAĞIRIŞ vaxtını görür" izahı. Go 1.22-dən
   əvvəl bu, RACE yaradırdı — loop var capture semantikası haqqında qeyd
   edin (müasir Go 1.22+ hər iterasiyada YENİ v yaradır — 3,3,3 artıq
   yaranmır, amma konsepsiya hələ vacibdir).

6. **Ch15 (Tao) mühazirə YOX, söhbət formatında:** fəlsəfə tələbələrin
   real təcrübələri ilə bağlansın — "hansı vaxt zorlamısan və alınmayıb?"

## 🎯 Tədris planı (draft)

| Fəsil | Müddət | Fokus |
|-------|--------|-------|
| 1-3 | 1 həftə | Kalkulyator: TDD tsikli, error pattern, go test/build |
| 4-5 | 1 həftə | Bookstore: struct, user stories, coverage |
| 6-7 | 1 həftə | Slice/map kolleksiyaları, go-cmp, random sıra |
| 8-9 | 1 həftə | Metodlar, pointer-lər, wrapping |
| 10 | 1 həftə | Validasiya sistemi (unexported/constructor/iota) |
| 11-12 | 1 həftə | Control flow, happy path, switch/for |
| 13-14 | 1 həftə | Closure/defer/variadic, binary/cross-compile |
| 15 | ½ həftə | Tao müzakirəsi |

## ⚠️ Diqqət çilənləri

- **kitab Go 1.19 dövründədir:** loop var semantikası (ch13), `sort.Slice`
  (müasir: `slices.SortFunc`), rand determinizmi — müasir alternativləri qeyd et
- **windows fərqi:** gofmt -d `diff` tələb edir (ch1) — öncədən xəbər ver
- **fail mesajlarının keyfiyyəti** kitabda dəfələrlə vurğulanır — tələbə
  layihələrində "want X, got Y" minimumu TƏLƏB ET
- **staticcheck quraşdırma:** SA4005 dərsləri real lint ilə təsdiqlənsin
