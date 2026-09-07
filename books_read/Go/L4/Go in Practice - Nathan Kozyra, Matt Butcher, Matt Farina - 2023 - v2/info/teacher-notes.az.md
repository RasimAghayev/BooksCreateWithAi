# Go in Practice, Second Edition — Müəllim Qeydləri (Azərbaycanca)

## 📖 Kitab deyir:

Kitab "Problem → Solution → Discussion" triadası ilə yazılıb: hər texnika real dünya
problemi ilə açılır, idiomatik həll göstərilir, sonra nüanslar müzakirə olunur. 4 hissə:
fundamentallar (CLI, struct, error), möhkəm tətbiqlər (concurrency, testlər, fayl/şəbəkə),
web tətbiqləri (HTTP server, template, upload, xarici servis), bulud + advanced
(mikroservis, deploy, refleksiya, codegen).

Kitabın ən məqamlı tövsiyələri:
- "Returning nil results along with errors isn't always the best practice" — işlənəbilən
  dəyər + error birgə qaytar
- "log.Fatal* calls os.Exit — the deferred function isn't called" (şəbəkə log-da!)
- "The go toolchain provides the ability to cross-compile out of the box"
- "Code generation is a somewhat common approach... generated code is much faster
  at runtime than reflection-based code"
- "Close your channels and return from your goroutines" (done channel pattern)

## 👨‍🏫 Müəllim qeydi:

1. **Kitabın ən böyük dəyəri formatındadır.** Problem/Solution/Discussion hər texnikanı
   "niyə?" sualı ilə birləşdirir. Tədrisdə bu strukturu saxlamağa çalış — sadəcə sintaksis
   deyil, hər həllin motivasiyası əsas mövzudur. Məsələn graceful shutdown tək funksiya
   yox, "deploy zamanı istifadəçi təcrübəsi" problemi kimi təqdim edilməlidir.

2. **1.22 tarixçəsinə diqqət.** Kitab Go 1.22-dən ƏVVƏL yazılıb, amma ikinci nəşr bəzi
   1.22 xüsusiyyətlərini (metod-prefiksli route, PathValue) əhatə edir. Oxucular bu
   versiya sərhədini bilsin: köhnə resurslarda "bun üçün üçüncü tərəf mux lazımdır"
   deyəndə, müasir Go-da bu ARTIQ YOXDUR — standart ServeMux kifayətdir.

3. **Chapter 4 error fəlsəfəsi digər kitabların ən çox yanlış çatdırdığı yerdir.** Burada
   panic/error ayrımı "bacarıqsızlıq" kimi YOX, "kontrol axını sahibliyi" kimi izah olunur:
   error = sənə sənədləşdirilə, gözlənilən problem; panic = davam MÜMKÜNSÜZ hal. Fatal-ın
   defer-ləri öldürməsi (os.Exit) network logging nümunəsində çox yadda qalıcı dərsdir.

4. **Generics fəsli (Ch3) minimal, amma düzgün.** Kitab generics-i "çox tip üçün eyni
   alqoritm" kimi təqdim edir və ~ təxmini + constraints paketini göstərir. Amma real
   layihələrdə generics-in ƏSAS istifadə sahəsi bu gündə standart slices/maps paketləridir
   — onlara ayrıca vaxt ayır (book.Comparable, cmp.Ordered).

5. **Slog (6.2.2) və 1.21+ xüsusiyyətləri ikinci nəşrin "təzə" hissəsidir.** JSONS axını
   + jq birləşməsi production observability üçün minimum standartdır; köhnə log.Print
   yanaşması ilə müqayisə etdirərək öyrət.

## Ən vacib 5 fikir

1. Production Go = explicit error axını + defer zəmanəti + `-race`/`--fuzz` CI intizamı.
2. 1.22 routeri + go:embed + slog — modern standart kitabxana köhnə "lazımsız paket"
   siyahısını ləğv edib.
3. Concurrency ədəbiyyatı: parametr ötür, sender bağlayır, done siqnalı ver — 3 qızıl qayda.
4. Performans: parse/template keş, KeepAlive transport, dərhal Body.Close, JSON codegen.
5. go:generate = dev-time metaproqramlaşdırma; nəticə kod VCS-də — Go-nun makro cavabı.

## Kitabın ən dəyərli hissəsi

Chapter 4 (errors) — panic/recover/defer/named-return interaksiyasının ən dolğun
praktik təhlili; Chapter 12.3 (REST sürəti) — keep-alive və connection reuse detalı,
başqa heç bir giriş kitabında bu dərinlikdə yoxdur.
