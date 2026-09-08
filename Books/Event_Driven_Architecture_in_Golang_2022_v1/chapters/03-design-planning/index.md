# Chapter 3 — Design and Planning (səh. 62-93)

## Bu chapter nədən bəhs edir?

MallBots-un dizayn mərhələsi: EventStorming workshop-u (Big Picture →
Design-level), kağız üzərində hadisə axınlarının kəşfi, storytelling ilə implicit
hadisələrin aşkarlanması, people/system identifikasiyası, hotspot analizi və
Architecture Decision Record (ADR).

## Əsas fikirlər

### 1. EventStorming nədir?
**Rəngli sticky notes** ilə tətbiqin hadisə axınını vizuallaşdıran, asan və
cəlbedici workshop. Şirkət üzrə müxtəlif bölmələrdən insanlar bir masaya toplanır
(developer, product owner, business, dəstək...).

**Anlayışlar:**
- **Domain Event (narıncı)** — "OrderSubmitted" — business dilində keçmiş hadisə
- **Command (mavi)** — hadisəni tetikleyən əmr (SubmitOrder)
- **Aggregate (sarı)** — komandanın idarə etdiyi model (Order)
- **Policy / Reaction (bənövşəyi)** — "whenever X, always Y" avtomatik reaksiya
- **Hotspot (qırmızı/bənövşəyi)** — problem/şübhə sahəsi
- **Opportunity (yaşıl)** — təkmilləşmə şansı
- Ubiquitous Language söz lüğəti də burada formalaşır

### 2. Big Picture mərhələləri
1. **Chaotic exploration** — hər kəs bildiyi hadisələri qeyd edir; heç nə
   silinmir, ikinci düşünmələr qadağandır ("fight the urge to crumple")
2. **Enforce the timeline** — hadisələr xronoloji sıraya düzülür
3. **Pivotal events** — mənalı dəyişiklik nöqtələri (OrderAccepted) + **swim
   lanes** (axın zolaqları) ilə qruplaşdırma
4. **People & systems** — xarici/iç istifadəçilər (Store owner, Customer, Admin)
   və sistemlər (Payment gateway) ayrıca rənglə tanınır
5. **Storytelling** — iştirakçı axını "hekayə" kimi danışır; pəncərə açıq
   qalmış yerlər (görünməyən hadisələr) üzə çıxır: "Order Canceled nə vaxt olur?
   Müştəri nə ilə ləğv edir?" → yeni hadisələr kəşf olunur
6. **Hotspots & opportunities** — problemlər qeyd olunur

**Real kəşflər (nümunə):** BotAvailable = "on və idle" (yalnız "on" deyil);
müştərinin "magik" order ləğvi → OrderCancellationRequested hadisəsi əlavə
olundu; mağazanın müvəqqəti bağlanması ayrı hadisə tələb etdi.

### 3. Design-level EventStorming
Big Picture-dən sonra **tək core bounded context** dərinləşir:
- Yalnız həmin kontekstin komandası iştirak edir
- Nəticə **BDD senarilərinə** (Gherkin) çevrilir:

```gherkin
Feature: Create Store
  Scenario: Creating a new store
    Given a valid store owner is logged into the kiosk
    When they create a store
    Then the store is created
```

Go-da godog ilə addımlar implement olunur:
```go
func InitializeScenario(ctx *godog.ScenarioContext) {
    ctx.Step(`^a store called "([^"]*)" exists$`, aStoreExists)
    ctx.Step(`^a valid store owner is logged into the kiosk$`, aValidOwnerIsLoggedIn)
    // ...
}
```

### 4. ADR (Architecture Decision Record)
Michael Nygard formatı — qərarlar sənədlənir:
- **Title, Status (proposed/accepted), Context, Decision, Consequences**
- İlk qərarlar: "ADR log saxlanmalıdır" (meta-ADR), "modular monolith ilə
  başla" 
- Hər mühüm qərar (NATS seçimi, CQRS tətbiqi) üçün qeyd — gələcəkdə "niyə belə
  etdik" sualına cavab

## Əsas terminlər

- EventStorming (hadisə fırtınası)
- Sticky Note mədəniyyəti (domain event, command, aggregate, policy)
- Chaotic Exploration (xaotik kəşf)
- Pivotal Event / Swim Lane
- Storytelling (hekayə danışma)
- Hotspot / Opportunity
- Gherkin / BDD / godog
- ADR (Architecture Decision Record)

## Praktik nəticə

- Dizayn işə workshop ilə başla — hadisələr business-in dilindən gəlir
- "Heç nə silinməz" qaydası: xaotik mərhələdə tənqid YOX
- Storytelling ən güclü kəşf alətidir — axını danış, boşluqlar görünsün
- BDD senariləri design-level nəticənin birbaşa test-ə çevrilir
- Hər arxitektura qərarını ADR ilə qeyd et — amma qısa

## Mənbə

Pages: 62-93 (Chapter 3, Event-Driven Architecture in Golang)
