# Everyday Golang — Müəllim Qeydləri (Azərbaycanca)

📖 **Kitab deyir:**

Go-nu gündəlik pattern-lərlə öyrə: ən sadəsindən başla (standart flag,
WaitGroup), mürəkkəb alətə yalnız həddində keç (Cobra, kanallar). Hər nümunə
real production alətindən (OpenFaaS CLI, arkade) gəlir: HMAC webhook
doğrulaması, ISS pozisiya API-si, todo DB, bcrypt funksiyası. Testlərdə
assertion YOX — konsistenslik Go-nun dəyəridir. Concurrency pilləli təqdim
olunur: WaitGroup → errgroup → kanallar → worker pool → singleflight.
Sonda: Prometheus RED metrikaları, GitHub Actions release, multi-arch
Docker, OpenFaaS funksiyaları.

👨‍🏫 **Müəllim qeydi:**

Bu kitab "Go öyrənmək" kitabları içində ən practice-oriented olanlardan
biridir və müəllifin OpenFaaS infrastruktur təcrübəsi hər səhifədə hiss
olunur. Qiymətləndirmələrim:

1. **Ən dəyərli dərs — pilləli sadəlik:** hər texnologiya üçün "əvvəl sadə,
   sonra güclü" prinsipi (flag→Cobra, WaitGroup→errgroup). Tələbələr əksər
   hallarda ƏKSLİNİ edir — dərhal Cobra/kanal ilə başlayırlar. Kitabın bu
   mövqeyi KONGO-ya yaxındır: yalnız lazım olanda mürəkkəbləşdir.
2. **Loop-capture tələsi (Ch7)** — Go 1.22-dən ƏVVƏLKi semantika ilə
   izah olunur. 1.22+ versiyalarında for dəyişəni per-iteration olur; amma
   köhnə kod bazaları və TC müsahibələri üçün bu bilik HƏLƏ vacibdir —
   kitabın izahı ən aydın olanlardan biridir.
3. **2021 tarixi haqqında:** ioutil hələ istifadə olunur ( indi os.ReadFile/
   io.ReadAll); GO111MODULE build-arg-ları artıq lazımsız; AWS Lambda Go
   dəstəyi verir. Amma pattern-lərin hamısı dəyişməzdir.
4. **OpenFaaS fəsli bonus kimi:** serverless marağı olmayan oxucu üçün əsas
   dərs oradadır: init()-də DB bağlantısı, secrets fayl mount, handler
   testlərinin build-də icrası — bunlar K8s/Cloud Run istənilən Go
   konteynerinə aiddir.

## Ən vacib 5 fikir

1. **Sadədən başla:** standart kitabxana 90% hallarda bəsdir; həddi
   aşanda dəyiş (Ch2, 5).
2. **İlk testi SİNDİR:** yaşıl test heç nə sübut etmir (Ch6).
3. **Concurrency hər alətin öz problemi:** WaitGroup=bitmə, errgroup=xəta,
   pool=limit, singleflight=dedupe (Ch7) — universal həll YOXDUR.
4. **RED metrikanı 3 sətirlə əldə et:** InstrumentHandlerCounter/Duration —
   mikro servisdə ölçmə bəhanəsi olmasın (Ch14).
5. **Pipeline-ı tag-ə bağla:** push → test+build+upload+push-image —
   paylanma əl işi olmasın (Ch15-16).

## Kitabın ən dəyərli hissəsi

Ch7 (Goroutines) — loop-capture + errgroup + singleflight + worker pool bir
fəsildə; və Ch6 — httptest və izolyasiya pattern-ləri. Bu iki fəsil Go
müsahibə suallarının 80%-ini əhatə edir.

## Kim üçündür

Başlanğıc→orta (L2) Go öyrənənlər; xüsusən DevOps/cloud yönümlü
developer-lər (Prometheus, Docker, serverless fəsilləri onlara aid).
"Learning Go" (Bodner) ilə paralel oxuna bilər: bu praktik, o nəzəri
tamamlıq verir.
