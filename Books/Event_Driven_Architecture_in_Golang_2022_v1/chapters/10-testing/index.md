# Chapter 10 — Testing (səh. 278-315)

## Bu chapter nədən bəhs edir?

Test strategiyası: unit (mocks/fakes), integration (Docker konteynerləri,
testify suite), contract testlər (CDCT — Pact REST + mesajlaşma) və end-to-end
(BDD/Gherkin, godog).

## Əsas fikirlər

### 1. Test piramidi (strategiya)
| Səviyyə | Nə yoxlanır | Qiymət |
|---|---|---|
| Unit | Domain + application funksiyaları, asılılıqlar mock | ucuz/sürətli |
| Integration | Real asılılıqlar (Postgres, NATS) konteynerdə | orta |
| Contract | Modullararası müqavilələr | orta |
| E2E | Bütün sistem, istifadəçi axını | bahalı |

### 2. Unit testlər
- Yalnız funksiya kodu; asılılıqlar **mock/fake** ilə
- Table-driven + subtest (`t.Run`) — hər ssenari ayrı başlıqlı
- Mock nümunəsi: test cədvəlində `on func(f mocks)` — hər ssenari öz mock
  konfiqurasiyasını daşıyır

### 3. Integration testlər
Real DB lazımdır → **Docker konteyneri** test üçün qalxır:
- **testcontainers** (və ya testify) ilə Postgres konteyneri: port təyinatı
  (`port.Port()`), healthcheck/retry, `Timeout(5*time.Second)`
- **Testify suite:** SetupTest/TearDownTest — hər test eyni "slate"-dən başlayır
  (şərh silinir, state reset)
- Test faylları build tag-lərlə idarə olunur: `//go:build integration` →
  `go test ./... -tags integration` (normal run-daə skip)
- Xüsusi fayl/direktoriya/run filtri də mümkündür

### 4. Contract testlər (CDCT)
Modullararası əlaqə nöqtələri çoxdur; hər cütlük üçün ayrı integration test
baha başa gəlir. Əvəzinə **consumer-driven** müqavilə:
- Consumer öz gözləntisini "pact" kimi yazır (axios çağırışı + MatchersV3)
- Provider öz tərəfindən pact-fə qarşı yoxlanılır
- **Pact Broker:** pact-ların anbarı; CI/CD inteqrasiyası — müqavilə dəyişəndə
  provider dərhal xəbərdar olur
- Provider state (`given()`): consumer-in gözlədiyi başlanğıc vəziyyətlər
  provider-də qurulur

**REST nümunəsi (consumer):**
```js
await provider.executeTest((mockServer) => {
  const client = new BasketClient(mockServer.url)
  return client.addItem(...).then(response => {
    expect(response.status).to.eq(400)
  })
})
```

**Messaging nümunəsi:** Consumer gözlədiyi mesajı `pact.AddMessage()` ilə
müqaviləyə yazır; provider tərəfində real module qaldırılıb mesaj buraxılır →
müqavilədəki mesajla müqayisə. Mock provider/consumer YOX — real mesajın
özünü doğrulama.

### 5. End-to-end testlər + BDD
- **İkili loop:** BDD (dış, Gherkin) + TDD (daxili, unit)
- Gherkin senariləri → godog addımları (regex parametr tutumu: `a store called
  "([^"]*)"`)
- REST client-lər hər modul üçün; UI yoxdursa API səviyyəsində E2E
- MallBots-un öz addım kitabçası: "creating an order" axını

## Əsas terminlər

- Testing Strategy / Test Pyramid
- Mock vs Fake (yalançı impl)
- testcontainers / Testify Suite
- Build Tags (integration testlərin ayrılması)
- CDCT (Consumer-Driven Contract Testing)
- Pact / Pact Broker / Provider State
- BDD + TDD double loop
- Gherkin / godog

## Praktik nəticə

- Strategiya: unit (mock) → integration (container) → contract (pact) → e2e (bdd)
- DB testləri üçün konteyner + suite teardown — hər test təmiz slate-dən
- Build tag ilə ağır testləri standart run-dan ayır
- Modullararası müqavilələri Pact ilə dondur; Broker-u CI-ya bağla
- Event mesajlarını da contract testlə — real publisher-dən doğrula

## Mənbə

Pages: 278-315 (Chapter 10, Event-Driven Architecture in Golang)
