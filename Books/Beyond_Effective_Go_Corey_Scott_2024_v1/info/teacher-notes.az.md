# Beyond Effective Go — Part 2 — Müəllim Qeydləri (Azərbaycanca)

📖 **Kitab deyir:**

Yüksək keyfiyyətli Go kodu: design prinsipləri Go ideologiyası ilə (DRY-nin
sərhədləri, composition, accept interfaces), Code UX üç sütunu (clarity/
consistency/predictability), behavior-yönümlü advanced testlər (TDT
kanonu, recorder-lər, concurrency testi), productivlik alətləri (gopls/
dlv/staticcheck, git axını), qeyri-adi pattern-lər (FP konseptləri,
middleware/decorator/futures, struct fndləri) və metaproqramlaşdırma
(API inteqrasiya, exec koordinasiya, AST+template kodgen).

👨‍🏫 **Müəllim qeydi:**

Corey Scott-un seriyası "Learning Go" (Bodner) ilə birlikdə müasir Go
kanonunun əsasını təşkil edir. Part 2-nin ən dəyərli tərəfi: HƏR texnika
üçün "nə vaxt YAZMA" sərhədi — bu, əksər texniki kitabda yoxdur. Tələbə
tövsiyələrim:

1. **Ch4-5 (design + UX) bütün junior/middle Go developer-lər üçün MÜTLƏQ
   oxunuşdur** — "accept interfaces, return structs" və DRY-nin
   kosmologiyası müsahibələrin əsas mövzularıdır. Kitabın ən qalıcı dəyəri
   bu iki fəsildədir.
2. **Ch6 (testing) — ən yaxşı praktik TDT dərsi** (Bodner-in nəzəriyyəsindən
   daha əlçatan); recorder patterni və t.Cleanup imkanları real layihələrdə
   dərhal tətbiq olunacaq.
3. **Ch8 (FP) — diqqətli oxun:** müəllif özü qeyd edir ki, currying və bəzi
   FP konstruksiyaları Go-nun ideologiyasına yadır. Öyrənin, amma gündəlik
   kode-yə KÖÇÜRMƏYİN — futures üçün sadə channel, options üçün sadə
   variantlar.
4. **Ch9 (metaproqramlaşdırma) — mühəndislik intizamının dərsi:** "alətə
   sərf olunan vaxt < qənaət olunan vaxt" meyarı bütün avtomatlaşdırma
   qərarlarına aiddir. AST bölməsi dar auditoriyaya hitab edir — dar
   kitabxana/framework müəllifləri üçün.

## Ən vacib 5 fikir

1. **DRY şərti deyil, KOMPROMİSdir:** yanlış abstraksiya dublikatdan
   bahadır (Ch4).
2. **Kod = insan ünsiyyəti:** clarity/consistency/predictability ölçüsüz
   görünür, amma hər biri konkret praktikadır (Ch5).
3. **Test implementation-i YOX, behavior-u yoxlayır:** recorder-lərin
   struktur yoxlaması bunun istisnasıdır — CALL ardıcıllığı davranışdır (Ch6).
4. **"Be lazy":** avtomatlaşdırma vaxt-hesabı ilə; staticcheck-dən
   gopls-ə qədər alətlər pulsuz sürətdir (Ch7).
5. **Go-da FP hissəvidir:** funksiyalar birinci sinif vətəndaşdır, amma
   currying/immutability tam DEYİL — idiomatik qal (Ch8).

## Kitabın ən dəyərli hissəsi

Ch5 (Code UX) — müəllifin "kod = kommunikasiya" çərçivəsi unikaldır; və
Ch6-nın recorder/TDT bölmələri. Ch4 klassikdir, amma "Effective Go" +
"Learning Go" oxuyanlar üçün tanış olacaq.

## Kim üçündür

L3+ developer-lər: real istehsal təcrübəsi olanlar üçün "niyə" suallarının
cavabı; junior üçün Ch4-5 oxuna bilər, amma Ch6+ təcrübə tələb edir.
Seriyadan əvvəl Part 1 və ya Effective Go + 1 il praktika tövsiyə olunur.
