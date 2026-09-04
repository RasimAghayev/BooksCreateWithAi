# Chapter 0 — Введение

Pages: 7-12

## Bu chapter nədən bəhs edir?
Bu chapter kitabın giriş hissəsidir. Kitabın məqsədi, hədəf auditoriyası,
Go dilinin ümumi baxışı, monolit və mikroservis arxitekturasının müqayisəsi,
kitabda işlənən texnologiyalar və fayl arxivi haqqında məlumat verir.

## Əsas fikirlər

### 1. Kitabın məqsədi
Bu kitab praktiki olaraq Go dilində mikroservis tətbiqləri necə
hazırlanır bunu göstərir. İş bir kiçik mikroservisdən başlayaraq,
müxtəlif texnologiyalar tətbiq edən mikroservis qrupuna və virtual
serverlərə yerləşdirilmiş hazır tətbiqə qədər davam edir.

### 2. Hədəf auditoriya
Kitab praqti bacarıqlarına üstünlük verir. Təkrarlanan Oxford (Oxford / Oksford
üniversiteti) Go dili haqqında dərin bilik tələb edir. Təkrar edilmir:
kitab daxilində Go dilinin özünə dərin dəriləriş edilmir, əsas fokus
mikroservislər, onların hazırlanması zamanı qarşılaşılan problemlər və
həll yollarıdır.

İşlənən texnologiya stack-i:
- Golang
- PostgreSQL
- Kafka
- Redis
- Docker
- Docker-Compose
- Kubernetes

### 3. Go dilinin ümumi baxışı
Go (Golang) Google tərəfindən yaradılmış, server development (server development
/ server tərəfi inkişafı) üçün nəzərdə tutulmuş compiled (compiled / kompilyasiya
edilən) proqramlaşdırma dilidir.

Xüsusiyyətləri:
- **Goroutines (goroutines / qorutinlər)** — yüngül execution thread-ləri
  (execution thread / icra növü), sistem thread-lərindən daha ucuzdur, on min
  sayda eyni anda işləyə bilərlər.
- **Channels (channels / kanallar)** — goroutine-lər arasında data exchange
  (data exchange / məlumat mübadiləsi) üçün istifadə olunur, race condition
  (race condition / yarış vəziyyəti) kimi klasik problemlərsiz təhlükəsiz
  multithreaded (multithreaded / çoxnüvuli) proqramlar yazmağa imkan verir.
- **Sürətli kompilyasiya** — maşın koduna kompilyasiya edir, interpretasiya
  dilləri (JS, Python) ilə müqayisədə daha sürətlidir. Build prosesi saniyələrlə
  gedir.
- **Minimalizm** — kiçik standard sintaksis, function overloading (function
  overloading / funksiya üstünə yükləmə) yoxdur, class inheritance (class
  inheritance / sinif irsi) yoxdur, komplex abstraksiyalar yoxdur. Lakin güclü
  standard library (standard library / standart kitabxana) var: şəbəkə, HTTP,
  fayl sistemi, proseslər.
- **Ekosistem** — go.mod ilə asılılıqların idarə edilməsi, pkg.go.dev
  açıq kitabxana repozitoriyası.

Mənfi cəhətləri:
- Sintaksis məhdudiyyətləri (bəzi developer-lər "çox sadə" və "sıxıcı" hesab edir)
- OOP yoxdur — hər şey composition (composition / tərkib) vasitəsilə həyata
  keçirilir
- Generics (Generics / cəmiçi tiplər) yeni əlavə edilib və funksionallığı
  Java/C#'dan geridədir

Tətbiq sahələri:
- Mikroservislər
- API xidmətləri
- Yüksək yüklü web server-lər
- Message broker-lər
- Paylanmış sistemlər
- Docker, Kubernetes, Prometheus kimi bulud alətləri Go ilə yazılıb

### 4. Monolit vs Mikroservis arxitekturası

**Monolitin mənfi cəhətləri:**
- Kod bazası böyüdükcə saxlanması çətinləşir, 20 nəfərdən çox developer
  komandası version conflict-lərə (version conflict / versiya ziddiyyətləri)
  səbəb olur.
- Tool-ların yenilənməsi və yeni texnologiyaların introducelarını (introducelarını
  / tətbiqini) çətinləşdirir.
- Təhlükəsizlik: köhnə kitabxanalar vulnerability (vulnerability / zəiflik)
  ehtimalı artırır.
- Yeni developer-lər üçün sistemə daxil olmaq çətinləşir.
- Deployment (deployment / yerləşdirmə) tezliyi aşağıdır — həftə ilə və ya
  aylarla, real dünya isə SaaS (SaaS / proqram kimi xidmət) və continuous
  deployment (continuous deployment / davamlı yerləşdirmə) tələb edir.
- Test və bug fix dövrü uzundur.
- Reliability (reliability / etibarlılıq) problemi: monolit iflas edərsə, bütün
  sistem iflas edir. Mikroservislərdə isə problemli hissəni ayırıb
  dayandırmaq mümkündür.

**Mikroservislərin mənfi cəhətləri:**
- Modullara ayrılma məsələsi: yanlış bölünmə distributed monolit
  (distributed monolit / paylanmış monolit) yaradır — bir-birindən
  asılılığı yüksək servislər qrupudur.
- Distributed systems (distributed systems / paylanmış sistemlər) mürəkkəbliyi:
  data consistency (data consistency / məlumatların ardıcıllığı) üçün iki
  mərhəmətli fiqaslama (2PC) və ya Saga (Saga / saga) pattern-lərindən istifadə
  etmək lazımdır.
- Bir çox kompleks alətlər tələb edir: Docker, Kubernetes, CI/CD, ArgoCD,
  Vault və s.

**Seçim qaydası:**
- Uzunmüddətli, böyük komandalı layihələr → mikroservislər
- Sadə, bir dəfə yerləşdiriləcək tətbiqlər → monolit (mikroservislər yalnış
  mürəkkəbliyə səbəb olar)

### 5. Fayl Arxivi
Kitabın sonuncu layihəsi (final project) fayl arxivində verilib:
https://zip.bhv.ru/9785977521208.zip

### 6. Əsas istinadlar
1. https://habr.com/ru/companies/kryptonite/articles/798703/
2. Kris Riçardson. Mikroservislər. İnkişaf və refaktorinq patternləri

## Əsas terminlər
- Microservices (mikroservislər)
- Monolith (monolit arxitektura)
- Golang (Go proqramlaşdırma dili)
- Goroutine (qorutin / yüngül icra növü)
- Channel (kanal / goroutine-lər arasında məlumat mübadiləsi)
- Race Condition (yarış vəziyyəti / eyni anda eyni resursa müraciət)
- Compilation (kompilyasiya)
- Standard Library (standart kitabxana)
- go.mod (Go modul faylı)
- Distributed Systems (paylanmış sistemlər)
- Version Conflict (versiya ziddiyyəti)
- Legacy Code (köhnə kod)
- Vulnerability (zəiflik)
- Continuous Deployment (davamlı yerləşdirmə)
- SaaS (proqram kimi xidmət)
- Docker (konteynerləşdirmə platforması)
- Kubernetes (konteyner orkestrasiya platforması)
- CI/CD (davamlı inteqrasiya/davamlı yerləşdirmə)
- 2PC (Two-Phase Commit / iki mərhəmətli fiqaslama)
- Saga (kompensasiya pattern-i)

## Mənbə
Pages: 7-12
