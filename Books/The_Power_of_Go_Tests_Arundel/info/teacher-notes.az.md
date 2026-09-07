# The Power of Go: Tests — Müəllim Qeydləri (Azərbaycanca)

## 📖 Kitab deyir:

Testlərə davranış kimi baxın (funksiyaların daxilinə yox), test-first yazın, kiçik
addımlarla gedərək inam qazanın, fail mesajlarını gələcək özünüz üçün yazın. Kitabın
əsas mesajları:

- "Testing shows the presence, not the absence of bugs" — Dijkstra; test inamı
  ARTIRIR, sübut ETMİR
- "What are we really testing here?" — hər fəslin təkrarlanan açar sualı
- "We simply can't write untestable functions when the test comes first" — test-first
  yazsan untestable kod YARANA BİLMƏZ
- "The tests are a canary in a coal mine revealing by their distress the presence of
  evil design vapors" — Kent Beck; test çətinliyi = dizayn qoxusu
- "Just because a statement is executed does not mean it is bug-free" — Donovan &
  Kernighan; coverage = siqnal, hədəf DEYİL
- "Tests are not a substitute for thinking" — Rich Hickey
- "Bad tests are worse than no tests" — flaky testlərə toleranssızlıq
- "Thou shalt not suffer a flaky test to live" — sil və ya düzəlt, dözme YOX

## 👨‍🏫 Müəllim qeydi:

Kitab Go test mədəniyyətinin ən yaxşı praktikalarını bir axında toplayır. Amma
tədrisdə bir neçə incəliyi vurğulayardım:

1. **Mock-lara münasibət kontekstə bağlıdır.** Kitab mock-lara kəskin mənfidir
   ("son option", "interface pollution"). Bu, DEFAULT mövqe kimi düzgündür — amma
   sqlmock nümunəsinin özində müəllif etiraf edir: bəzən lesser evil var. Tələbəyə
   qayda kimi deyil, DÜŞÜNMƏ alqoritmi kimi öyrədin: (1) dependency-ni ləğv edə
   bilərəm? (2) scope-u azalda? (3) chunk? (4) adapter? (5) ancaq sonra fake/mock.

2. **Fuzz fəsli syntaktik olaraq asan, konseptual olaraq çətindir.** Tələbələr
   f.Fuzz sintaksisini 5 dəqiqədə öyrənir, amma "property nədir?" sualında ilişir.
   Square non-negative nümunəsindən başlamaq yaxşıdır — sonra öz kodlarında
   "eyni input → eyni output", "roundtrip" kimi invariant axtarmağa yönəldin.

3. **testscript fəsli real dünyada ən çox İHMAL edilən kameradir.** Bir çox komanda
   CLI alətlərini heç vaxt binary səviyyəsində test etmir. Bu fəslin dəyəri birbaşa
   praktikadır: RunMain + txtar kombinasiyası tələbələrin modul layihələrinə
   birbaşa tətbiq oluna bilər.

4. **Mutation testing-in "passed/failed" terminologiyası TƏRSİdir:** "1 passed" =
   mutant ÖLDÜRÜLDÜ (yaxşı!), "4 failed" = mutant SAĞ QALDI (pis!). Test çıxışını
   oxuyarkən tələbələr bunu tez səhv oxuyur — xüsusi vurğulayın.

5. **`var Now = time.Now` seam-i paralel testlərlə DAVAMISHdır.** Kitab bəyan edir
   ("sequential olmalıdır"), amma tələbələr bunu buraxır → flaky race. Alternative
   olaraq parametr variantını da göstərin: `OneHourAgo(now time.Time)` — az paperwork
   ilə paralel-güvenli həll.

6. **Ch11 fəsli mühəndis MƏDƏNİYYƏTİ haqqındadır** — kod texnologiyasından çox.
   "Guerilla testing", zero defects, bug bankruptcy kimi mövzular yeni başlayanlara
   radikal gəlir; müzakirə formatında öyrətmək daha effektli olur.

## 🎯 Tədris planı (draft)

| Fəsil | Müddət | Fokus |
|-------|--------|-------|
| 1-2 | 1 həftə | TDD tsikli, go test/go-cmp/vet alətləri |
| 3 | 1 həftə | Table-driven + subtest + paralel; fail mesajları |
| 4-5 | 1 həftə | Error testləri; invalid input tabloları |
| 6-7 | 1 həftə | Fuzz + invariantlar; coverage + mutation |
| 8 | 1-2 həftə | "Untestable" refaktorları (ən praktik hisə) |
| 9 | 1 həftə | testscript + txtar |
| 10 | 1-2 həftə | Dependency strategiyası; adapter/interface |
| 11 | 1 həftə | Suite review, müzakirə |

## ⚠️ Diqqət çilənləri

- **double quotes testscript-də escape ETMİR** — literal çap olunur (shell refleksi təhlükəli)
- **TestMain-də os.Exit(testscript.RunMain(...)) unutma** — yoxsa custom komanda
  həmişə "uğurlu" görünür
- **assert.Equal frameworklərdə** "1 ≠ 0" deyir — kontekst YOX; standart testing +
  cmp.Diff daha informativdir
- **map iteration sırası müəyyənsizdir** — for range ilə müqayisə YOX, cmp.Equal işlət
- **time.Time == müqayisəsi** — Sub().Abs() + delta toleransı ilə
