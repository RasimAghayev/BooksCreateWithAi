# Chapter 2 — Supporting Patterns in Brief (səh. 40-61)

## Bu chapter nədən bəhs edir?

EDA-nı dəstəkləyən pattern-lər: Domain-Driven Design (subdomain, bounded
context, ubiquitous language), domain-centric arxitekturalar (onion, clean,
hexagonal), CQRS, monolith vs modular monolith vs microservices.

## Əsas fikirlər

### 1. DDD (Domain-Driven Design)
**Məqsəd:** Mürəkkəb business ideyasını software-ə model etmək üçün problemi
dərindən anlamaq.

**Anlayışlar:**
- **Subdomain-lərə bölmə:** Core (ən yüksək dəyər — MallBots-da Depot),
  Supporting (Orders, Payments), Generic (auth, notifications — hazır həll almağa
  dəyər)
- **Bounded Context (sərhədli kontekst):** Hər model müəyyən kontekstdə yaşamalıdır
- **Ubiquitous Language (ümumvahid dil):** Kontekst daxilində eyni terminlər —
  "Product" Stores-də kataloq elementi, Orders-də fərqli mənada OLA BİLƏR; hər
  kontekst öz UL-unu daşıyır

**Context münasibətləri:** Separate Ways (əlaqəsiz), Partnership (birgə
idarə), Customer-Supplier, Conformist, ACL (anti-corruption layer)...

**EDA bağlantısı:** İnteqrasiya hadisələrinin adları UL əsasında seçilir.

### 2. Domain-centric arxitekturalar (təkamül)
1. **Layered (tradisional):** Presentation → Business → Data — business data
   modelinə asılıdır, dəyişikliklər qatlara yayılır
2. **Onion (Palermo, 2008):** Asılılıqlar daxildən xaricə; domain mərkəzdə
3. **Clean (Martin):** "Source code dependencies can only point inwards" —
   framework-ə asılılıq yoxdur
4. **Hexagonal (Cockburn):** Ports & Adapters

**Hexagonal tətbiqi:**
```
[REST adapter]                    [gRPC adapter]
       \                               /
        → Domain model + domain services
       /                               \
[NATS adapter]                   [PostgreSQL adapter]
```
- **Domain** — model, domain service-lər, heç bir texnologiyaya asılı deyil
- **Ports (interfaces)** — domain-in təyin etdiyi "giriş/çıxış" müqavilələri
- **Adapters** — portları texnologiya ilə realləşdirir (NATS, REST, Postgres)
- Nəticə: kiçik, test oluna bilən komponentlər; "rulebook, not guidebook" —
  qayda kitabı, məcburiyyət deyil

### 3. CQRS (Command and Query Responsibility Segregation)
Obyektlər **Command** (yazma, dəyişmə) və **Query** (oxuma) olmaqla ikəyə bölünür.

**Lentin analogiyası:** Tətbiq = lent; sol tərəf komandalar (dəyişdirir), sağ
tərəf sorğular (oxuyur). Uzun müddət eyni model hər ikisini daşıyırsa:

- **Model səviyyəsində:** ayrı Command/Query modelləri (seliçilərin qarışması azalır)
- **Verilənlər bazası səviyyəsində:** yazma DB-si + ayrıca optimize olunmuş oxu
  modeli (projections); fine-tuned SQL-in həddi var — oxu nüsxələri sərbəst
  formalarda qurulur

**Nə vaxt:** task-based UI, event sourcing (hadisələrdən istənilən sayda
projection yaradılır), mürəkkəb oxu tələbləri.

### 4. Tətbiq arxitekturaları
| Arxitektura | Təsvir | Nə vaxt |
|---|---|---|
| **Monolith** | Tək kod bazası, tək deploy; sadə, amma böyüyəndə komanda toqquşması | Kiçik başlanğıc |
| **Modular Monolith** | Tək deploy, amma daxildə module-lər (MallBots yanaşması) | Orta miqyas — module-lər sonradan mikro-servislərə çıxır |
| **Microservices** | Hər servis ayrıca deploy/skal | Böyük komandalar, müstəqil buraxılış |

**Green field tövsiyəsi:** modular monolith-dən başla; sərhədləri (module) doğru
çək, sonra lazım olandan ayır.

## Əsas terminlər

- Domain-Driven Design (DDD)
- Core / Supporting / Generic Subdomain
- Bounded Context / Ubiquitous Language
- Onion / Clean / Hexagonal Architecture
- Ports & Adapters (portlar və adapterlər)
- CQRS — Command & Query Separation
- Modular Monolith
- Anti-Corruption Layer (ACL)

## Praktik nəticə

- Core domain-ə sərf et; generic-i hazır al
- Hər bounded context öz dilini və modellərini daşıyır — "Product" hər yerdə
  eyni deyil
- Domain mərkəzdə, texnologiya kənarda (hexagonal) — dəyişiklik lokal qalsın
- Oxu/yazma asimmetrik böyüyəndə CQRS; ilk gün deyil
- Modular monolith → microservices təkamül yolunu açıq saxla

## Mənbə

Pages: 40-61 (Chapter 2, Event-Driven Architecture in Golang)
