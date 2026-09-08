# Chapter 3 — Design and Planning (səh. 62-93)

## Bu chapter nədən bəhs edir?

MallBots tətbiqinin planlaşdırılması: EventStorming workshop-u (chaotic
exploration → timeline enforcement → storytelling), design-level EventStorming,
bounded context-lərin identifikasiyası, Gherkin scenario-lər (godog) və
Architecture Decision Record (ADR) anlayışı.

## Əsas fikirlər

### 1. EventStorming nədir?
**Nədir:** Rəngli sticky note-larla tətbiqin axınlarını vizuallaşdıran, domain
ekspertləri və developerlərin bərabər iştirak etdiyi sürətli workshop.
Məqsəd — bir neçə insanın başındakı implicit bilikləri hamı ilə paylaşmaq.

**Diqqət:** Workshop-da tətbiqin NƏ edəcəyi müzakirə olunur, NECƏ deyil —
texnologiya qərarları (web server, DB) QƏTIYYƏN verilmir.

**Əsas konseptlər:**
- **Domain Event** (portağal sticky) — "sifariş verildi" kimi bitmiş fakt
- **Ubiquitous language** tərifləri — iş dünyasının sözlüyü
- **Hotspot** (qırmızı) — açıq sual/problem
- **Opportunity** — gələcək təkmilləşdirmə təklifi
- **Happy/unhappy path** — axın nəticələrinin etiketləri

### 2. Workshop addımları (Big Picture)
1. **Chaotic exploration:** Hər kəs domain event-ləri vaxt xəttinə qoyur —
   "hər şey qalır" prinsipi: birinci düşüncə çox vaxt doğru istiqamətdədir,
   sadəcə terminologiya düzəldilməli olur; tərəddüd edirsənsə sticky-ni
   təcrid olunmuş yerə qoy, silmə
2. **Timeline enforcement:** Event-ləri xronoloji sıraya düz, əlaqəli
   event-ləri **flow**-lara qruplaşdır (flow = domain-ə aid proses). Pivotal
   event-lər + swim lane-lər ilə sərhədlər görünür (müştəri interaksiyaları
   üst lane-də)
3. **Identify people and systems:** Actors — People (Store owner, Store
   administrator, Customer, Bot) və external systems (payment, depot);
   swim-lane-lərə yerləşdirilir
4. **Storytelling:** Hər axın üçün "hekayə danışan" seçilir — iştirakçılar
   eksik event-ləri, yanlış sıraları tapır. Nümunə kəşflər:
   - Səbətə item əlavə/çıxarılmasında müştəriyə YENİ TOTAL göstərilmirdi
   - "Müştəri düşüncə ilə sifarişi ləğv edirdi" — order selected event əlavə
     olundu (UI interaction eksik idi)
   - Müvəqqəti bağlama: add → remove → re-add ƏVƏZİNƏ reopen axını
   - Bot-a iş təyinatı: polling YOX — bot availability status əsas götürüldü
     ("Bot availability: A bot's readiness to be given work" — ubiquitous
     language-a əlavə olundu)
5. **Hotspots & opportunities:** Problem/ideya note-ları ilə bağlanış

### 3. Temporal event-lər
Vaxtla bağlı event-lər xüsusi işarə daşıyır: alarm saatı = gecikmiş event,
analoq saat = günün konkret vaxtı, təqvim = həftə/ay günləri.

### 4. Identifying the contexts
Big Picture nəticəsindən **bounded context-lər** çıxarılır: hər axın qrupu
domain-ə uyğun kontekstə çevrilir (ordering, store management, bot/depot).

### 5. Design-level EventStorming
- Yalnız **bir core bounded context** üzrə, daha kiçik qrupla
- Məqsəd dəyişir: nə etməli → **necə etməli**
- Event flow-a komanda (mavi sticky), policy (bənövşəyi), aggregate,
  read model kimi konseptlər əlavə olunur

### 6. Gherkin scenario-lər — godog
Domain ekspertinin yaza biləcəyi format — NƏ yazılır, NECƏ deyil:

```gherkin
Feature: Create Store
  Store owners need to manage their stores

  Scenario: Create a new store
    Given a store called "Waldorf Books" does not exist
    And a valid store owner is logged in
    When I create the store called "Waldorf Books"
    Then a store called "Waldorf Books" exists
```

**Go implementasiyası (godog):**
```go
func InitializeScenario(ctx *godog.ScenarioContext) {
    ctx.Step(`^a store called "([^"]*)" exists$`, aStoreExists)
    ctx.Step(`^a valid store owner is logged in$`, aValidStoreOwner)
    ctx.Step(`^I create the store called "([^"]*)"$`, iCreateTheStore)
    ctx.Step(`^no store called "([^"]*)" exists$`, noStoreExists)
}
```
Sub-kod izahı: Gherkin sətirləri regex ilə Go funksiyalarına bağlanır —
sənəd eyni zamanda icra olunan testə çevrilir.

### 7. Architecture Decision Record (ADR)
Michael Nygard formatı (Markdown):
```
# {RecordNum}. {Title}
## Context — bu qərarı motivasiya edən problem nədir?
## Decision — nə təklif/edinirik?
## Status — Proposed, Accepted, Deprecated, Superseded
## Consequences — nəticələr (müsbət və mənfi)
```

MallBots-in ilk 2 qərarı:
1. **Architecture decision log saxlanılır** — qərarların tarixçəsi app-ə böyük
   təsir edəcək
2. **Modular monolith istifadə olunur** — struktursuz monolith-in xaosundan və
   mikroservis deploy kompleksliyindən eyni anda qaçmaq

## Əsas terminlər

- EventStorming (hadisə fırtınası)
- Chaotic Exploration (xaotik kəşf)
- Pivotal Event (mərkəzi hadisə)
- Swim Lane (üzüş zolağı)
- Storytelling (hekayə danma)
- Bounded Context (məhdud kontekst)
- Ubiquitous Language (ümumi dil)
- Temporal Event (zaman hadisəsi)
- Gherkin / godog
- ADR (Architecture Decision Record)
- Modular Monolith (modul monolit)

## Praktik nəticə

- EventStorming: texnologiyasız, nə-fokuslu; hekayə danışma eksik
  event-ləri üzə çıxarır — UI interaction-ları unutma
- Scenario-lar domain eksperti tərəfindən oxuna bilən Gherkin-də yazılır,
  godog ilə executable testə çevrilir
- Hər əhəmiyyətli qərar ADR formatında loglanır — "niyə beli" sualına
  cavab tarixçədə qalır
- Modular monolith — mikroservisə keçid üçün hazır strukturlu başlanğıc

## Mənbə

Pages: 62-93 (Chapter 3, Event-Driven Architecture in Golang)
