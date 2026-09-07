# Введение — Giriş (səh. 7-12)

## Bu bölmə nədən bəhs edir?

Kitabın məqsədi (bir mikroservisdən tam mikroservis qrupuna qədər praktik
yol), Golang-ın ümumi icmalı (nə üçün mikroservislər üçün ideal: goroutine +
kanal modeli, sürətli compile, minimализм, ekosistem), monolit vs
mikroservis müqayisəsi (monolitin 5 problemi: dəstə çətinliyi, təhlükəsizlik,
anlaşılmazlıq, gec çatdırma, etibarsızlıq; mikroservislərin öz çatışmamazlıqları:
düzgün modul bölgüsü tələbi, paylanmış sistem mürəkkəbliyi — distributed
monolith təhlükəsi, 2PC/Saga mexanizmləri), fayl arxivi və mənbələr.

## Əsas fikirlər

### 1. Kitabın Konsepsiyası və Texnologiya Stəki
- **Yol:** tək mikro servis (User) → tam qrup (Auth, Gateway, Transaction) →
  virtual serverlərə çatdırılma
- **Stek:** Golang, PostgreSQL, Kafka, Redis, Docker, Docker-Compose,
  Kubernetes + gRPC, Protobuf, Swagger, ORМ
- **Tələb:** Golang dərin bilik tələb olunmur — dilə deyil, MİKROSERVİS
  problemlərinə fokus

### 2. Niyə Golang? (Goroutine + Kanal Modeli)
- **Nədir:** Google-də yaradılmış, server üzrə kompilyasiya olunan dil
- **Güclü tərəfləri:**
  - Goroutine-lər — bir sətirdə yaradılan yüngül axınlar; sistem thread-lərindən
    qat-qat ucuz, on minlərlə paralel işləyə bilər
  - Kanallar — goroutine-lər arasında data mübadiləsi; race condition
    tələlərindən uzaq, təhlükəsiz paralellik
  - Compile → maşın kodu: saniyələr çəkir; JS/Python-dan sürətli icra
  - Minimализм: miras YOX, overloading YOX — amma güclü standart kitabxana
  - go.mod modul sistemi — dəqiq versiya fiksliyi + pkg.go.dev repozitoriyası
- **Çatışmamazlıqları:** sintaksis "çox sadə"; OOP yoxdur (kompozisiya);
  generics gec gəlib, funksionallığı Java/C#-dan zəif

### 3. Monolit vs Mikroservis
**Monolitin problemləri:**
1. **20 nəfərlik komanda = versiya konfliktləri** — bir kod bazasında
   iş çətindir; legacy-yə heç kim getmək istəmir
2. **Təhlükəsizlik** — köhnəlmiş kitabxanalar = yığılan zəifliklər
3. **Mürəkkəblik artımı** — IDE gecikir, build gecikir, yeni üzvlər
   itirir
4. **Çatdırılma tempi** — 2 həftədə bir reliz; SaaS dövründə gündə
   bir neçə dəfə production-a çıxış gözlənilir
5. **Etibarsızlıq** — monolit çökürsə, BÜTÜN sistem çökür; funksiyanın
   bir hissəsini söndürmək mümkün deyil

**Mikroservisin tələbləri:**
- Düzgün modul bölgüsü — əks halda **paylanmış monolit** (distributed
  monolith): bir-birinə kəskin bağlı, ayrı yaşaya bilməyən servis dəsti —
  İKİ dünyanın ən pis cəhətləri bir yerdə
- Hər servisin öz DB-si → data konsistensiyası ayrıca mexanizmlə:
  **iki fazalı commit (2PC)** və ya **Saga** (bölmə 4-də ətraflı)
- Dağıtım sürəti üçün: Docker, Kubernetes, CI/CD, ArgoCD, Vault
  (sekrete/gizli açarlar üçün)

### 4. Seçim Qaydası
- **Mikroservis:** illər boyu yaşayacaq, komandalı, inkişaf planlı →
  başdan mikro servis imkanını qoy
- **Monolit:** bir dəfə yayımlanacaq mini-sayt, dəyişiklik olmayacaq →
  monolit kifayətdir; mikro servis = lazımsız mürəkkəblik

## Əsas terminlər
- Goroutine (qısa yüngül axın) — Go-nun paralellik vahidi
- Kanal (channel) — goroutine-lər arası təhlükəsiz data mübadiləsi
- Monolit (bütün-əhatəli tək proqram) — hamısı bir kod bazasında
- Distributed monolith (paylanmış monolit) — yanlış bölgü nəticəsi:
  mikro servis görünen, amma bir-birindən asılı servis dəsti
- Saga (saga pattern) — paylanmış transaksiyalar üçün addım-addım
  kompensasiya modeli
- 2PC (iki fazalı commit) — paylanmış sistemdə data uyğunluğu üçün
  hazırlıq + təsdiq protokolu
- SaaS (Software-as-a-Service) — fasiləsiz yayım modeli; gündə bir neçə
  dəfə production yeniləmə
- Legacy code (köhnəlmiş kod) — dəstəyi çətin, inkişafdan qalan kod

## Praktik nəticə

1. **Mikroservisə keçmə qərarı texniki deyil, BİZNES qərarıdır:** layihənin
   ömrü, komanda ölçüsü, inkişaf tempi — bunlar arxitekturanı müəyyən edir.
2. **Goroutine + kanal:** Go-nu mikro servis dili edən əsas — thread-lər
   ucuza düşür, IPC təhlükəsizləşir.
3. **Bölgü dizaynı ƏSAS riskdir:** yanlış bölgü paylanmış monolit yaradır —
   ən pis variant. Domen sərhədlərini (User, Auth, Transaction) əvvəlcədən
   dəqiq müəyyən et.
4. **Paylanmış transaksiya problemini BAŞDAN qəbul et:** hər servisin öz
   DB-si olanda konsistensiya pulsuz olmur — 2PC/Saga (bölmə 4).
5. **Fayl arxivi:** https://zip.bhv.ru/9785977521208.zip — kitabın final
   layihə kodu.

## Mənbə
Pages: 7-12 (PDF 8-13)
