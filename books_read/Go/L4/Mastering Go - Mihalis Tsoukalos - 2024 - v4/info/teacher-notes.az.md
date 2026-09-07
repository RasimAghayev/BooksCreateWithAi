# Mastering Go, Fourth Edition — Müəllim Qeydləri (Azərbaycanca)

## 📖 Kitab deyir:

Go-nu "gofers" üçün tampectrum kitabı kimi təqdim edir: əsas tiplərdən başlayaraq
generics, reflection, sistem proqramlaşdırma, concurrency, web/REST, test/fuzz,
observability və performansa qədər — hamısı bir statistika tətbiqinin təkamülü
üzərində. Kitabın mərkəzi mesajları:

- "Make it work, then make it beautiful, then if you really have to, make it fast"
- "Do not use any of its features prematurely" — paketi, interfeysi, generic-si
  ehtiyac yarananda yarat
- "Just because we can use channels, it does not mean that we should" — sadə həllər
  (unikal index + WaitGroup) üstündür
- "Make sure your program benefits from concurrency" — paralelliyi əvvəlcədən YOX,
  ehtiyacında əlavə et (stats.go 3x sürətlənməsi ilə isbat olundu)
- "Logging is for application code, not library code" — paketlər error qaytarsın

## 👨‍🏫 Müəllim qeydi:

Kitab əla "lens"dir, amma praktikada bir neçə incəliyi vurğulayardım:

1. **Go 1.22 loop semantikası kitabın Ch8 mövzusunu arxa plana salır.** Kitab closure
   tələsini parametr ötürmə ilə həll edir (d3 nümunəsi) — texniki olaraq hələ də good
   practice-dir, amma yeni Go versiyalarında RACE ARTIQ yaranmır. Tədrisdə "niyə bu
   pattern var" tarixi kontekstində izah et, yoxsa tələbə parametr ötürməni məcburiyyət
   kimi əzbərləyir.

2. **Kitabın Go 1.22 ServeMux bölməsi çox qısadır — amma ən praktik yenilikdir.**
   `GET /users/{id}` pattern-ləri standart routerdə işləyir. Yeni layihələrdə gorilla/mux
   asılılığına ehtiyac əhəmiyyətli azalıb; kitabın Ch11 subrouter nümunələrini müasir
   net/http ilə müqayisə edərək öyrət.

3. **slog-i kitab gec tanıdır (Ch7).** Production kodunda Ch1-dən etibarən fmt.Println
   yerinə slog istifadə et — strukturlaşdırılmış log onuCh13-də Prometheus-a körpü olur;
   gec öyrəşən developer JSON handler-in dəyərini gec görür.

4. **Tri-color GC izahı əladır, amma "generational" xarakteristikası qarışıqdır.**
   Kitabın öz rəsmi sitasında "non-generational" deyilir, sonrakı bölmədə isə
   "generational collector" başlığı var. Go GC gerçəkdə non-generational-dır — tədrisdə
   bu ziddiyyəti açıq dey.

5. **Error handling-də kitabın `err.Error() == "..."` müqayisəsi (Ch2) bad practice
   kimi özü etiraf edir, amma nümunədə saxlayır.** Müasir standart: %w wrap + errors.Is/
   As (kitabın öz Go 1.13 resursu bunu göstərir). Tələbəyə Ch2-də YOX, Ch5+ formatında
   öyrət.

6. **RabbitMQ seçimi kontekstlidir.** Kitab "sürət → Kafka" deyir — amma məqsəd
   sadəliqdursa, cloud-native Go sistemlərində Go-nun öz kanal goroutine-ləri və ya
   NATS daha yüngül ola bilər; RabbitMQ yalnız broker MÜTLƏQ olanda.

7. **stats layihəsinin təkamülü kitabın ən güclü tədris alətidir.** Hər fəsildə EYNİ
   problemin daha abstrakt həllini görürük — bu, "refactoring yol xəritəsi" kimi
   oxunmalıdır. Tələbələrə eyni öz layihələrində bu pilləli yolu təkrar etməyi tapşır.

## Ən vacib 5 fikir

1. Sadəliyi qorumaq aktiv qərardır — hər yeni abstraksiya (interface, generic, channel)
   haqq qazanmalıdır.
2. Concurrency-də sıra: əvvəl düzgün, sonra paralel; ölçmə olmadan optimizasiya YOX.
3. Data strukturu = performans: slice-of-structs > map-of-pointers (40x); pointer az,
   dəyər çox.
4. Test piramidası: unit + table + fuzz + race + coverage — hamısı bir ekosistemdir.
5. Production Go = observability: ölçülməyən sistem idarə oluna bilməz.

## Kitabın ən dəyərli hissəsi

Chapter 8 (Go Concurrency) — scheduler-dən monitor goroutine-a qədər tam spektr; və
Appendix (GC) — tri-color alqoritminin yazılış tərzi ilə izahı. Bu ikili birlikdə Go
runtime-ın "maşınxanasını" açır. Praktik ikili isə Ch9+Ch11 (web→REST təkamülü) —
real servisin necə böyüdüyünü göstərir.

⚠️ Uyğunsuzluq qeydi: Ch2-də `math.Abs(float64(AnotherGlobal))` nümunəsi ilkin nəşrlərdə
tip çevirmə sintaksisi fərqli izah olunub; 4-cü nəşrdə çevirmə doğrudur. Həmçinin kitabın
R1() funksiyasında (Ch13, reverse.go) `sAr := []byte(sAr)` çap səhfi var — mənbə kodda
`sAr := []byte(s)` olmalıdır; correct.go versiyasında artıq düzgündür.
