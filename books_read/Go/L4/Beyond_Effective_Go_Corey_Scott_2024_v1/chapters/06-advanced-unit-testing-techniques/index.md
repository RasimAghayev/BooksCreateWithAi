# Chapter 6 — Advanced Unit Testing Techniques (səh. 148-236)

## Bu chapter nədən bəhs edir?

Unit testlərin nə üçün, nə vaxt, nə qədər və necə yazılması: test motivasiyası,
test piramidi (70/28/2), behavior vs implementation, table-driven tests (TDT),
senari quruluşu, mocks/stubs/test recorders, test UX, resilience, az tanınan
test fndləri (context timeout, testdata, test latch), concurrent kodun testi və
test-induced damage.

## Əsas fikirlər

### 1. Niyə test yazırıq?
Testlər əmr yox, **programmer-in iş həyatını asanlaşdıran alətdir**:
daha sürətli və cəsarətli dəyişikliklər (regression qorxusu olmadan),
intentin sənədləşməsi, yeni üzvlərin tez onboarding-i.

### 2. Test piramidi (səy bölgüsü)
| Səviyyə | Pay | Məqsəd | Zəif tərəfi |
|---|---|---|---|
| Unit | ~70% | Davranış səviyyəsində tez, ucuz testlər | Kiçik scope — sistemi yoxlamır |
| UAT (user acceptance) | ~28% | İstifadəçinin gözləntilərinin sistemi | Mock-lar üzərindən — xarici asılılıqlar yox |
| E2E (end-to-end) | ~2% | Xarici sistemlərin fərz edildiyi kimi işləməsi | Bahalı, yavaş, zəif signal/noise |

- Coverage hədəfi: **~70%** — xətt yox, **davranış coverage**-i sayılır
- Konfiqurasiya səhvləri E2E ilə yox, unit/UAT səviyyəsində yoxlanılır

### 3. Unit = obyekt + entry point
**Nədir:** Test uniti ayrı-ayrı private metodlar deyil, **整个 struct + public API**
(OrderManager.Process). Private metodlar implementation detail-dır — refactor
onları dəyişəndə testlər break olmasın deyə.

**Sad path/sad path senariləri:**
- Happy path: bütün daxilolmalar düzgün + asılılıqlar işləyir → gözlənilən nəticə
- Input xətası: doğrulama pozuntusu (validation sad path)
- Dependency xətası: mock-u failed qaytarır (məs. `bank.On("Charge", ...).Return("", errors.New("failed"))`)

### 4. Behavior yox implementation
**QADAĞANDIR:** exact error yoxlama, mock-un neçə dəfə çağırıldığını yoxlama,
çağırış ardıcıllığını yoxlama — bunlar implementation-a couple edir və refactor-i
çətinləşdirir. Testdə yalnız: "xəta gəldi/gəlmədi" (`expectAnErr bool`) və
çıxış dəyəri. **Müstəsna:** "charge fail olsa receipt GÖNDƏRİLMƏMƏLİDIR" — bu
behaviorın özüdür, mock-un `Run()`-unda `assert.FailNow` ilə səsli çatdırılır.

### 5. Table-Driven Tests (TDT)
```go
scenarios := []struct {
    scenarioDesc        string
    inputOrder          Order
    configureMockBank   func(bank *MockBank)
    configureMockSender func(sender *MockReceiptSender)
    expectedReceiptNo   string
    expectAnErr         bool
}{ /* ... */ }

for _, s := range scenarios {
    scenario := s
    t.Run(scenario.scenarioDesc, func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
        defer cancel()
        // mocks setup → call → validate
    })
}
```
**Sub-kod izahı:**
- `scenarioDesc` → subtest adı (test output grouping)
- `scenario := s` → closure dəyişən tutumu (t.Parallel təhlükəsizliyi)
- **Validasiya sırası:** əvvəl error, sonra result — gözlənilməz xəta vaxtı noise az olur
- `require` (FailNow) assertion-larda, `assert` (Fail) davam edən yoxlamalarda

### 6. Mocks, Stubs, Test Recorders
| Alət | Nə edir | Nə vaxt |
|---|---|---|
| Mock (Mockery) | Çağırışları gözləyir + qeyd edir | Asılılığın dəqiq davranışı lazım olanda |
| Stub | Sadə sabit dəyär qaytarır | Sadə keçid/canlı-davranış lazım deyilsə |
| Test Recorder | Daxilolmaları qeyd edir | Logger/instrumentation (çıxış yoxdur) |

**Mockery ilə generasiya:**
```go
//go:generate mockery -name=UserRepository -case underscore -testonly -inpkg
```
- `-testonly` → mock yalnız test faylında yaşayır (production kodu təmiz qalır)
- `mock.Anything` → arqument konkret əhəmiyyət daşımırsa
- Ardıcıllıq fərqli cavablar: `.Return(...).Once(...)` ilə "birinci xəta, ikinci uğur"
  (transient error simulyasiyası)
- `Run(func(args mock.Arguments))` → mutasiya edən metodlar, dinamik nəticə

**Quick stub (embedding ilə):** Tək metod lazımdırsa:
```go
type MockSave struct {
    Crud // interface-i embed et — qalan metodlar "hazır"
}
func (m *MockSave) Save(u *User) error { /* ... */ }
```

### 7. Test UX — 5 addımlı test strukturu
1. İnputları qur → 2. Asılılıqları (mock) qur → 3. Test edilən obyekt yarat →
4. Gözləntiləri təyin → 5. Çıxışları validasiya et.

Təkrarların azaldılması:
- Mock konfiq closure-larını **adlandırılmış funksiyalara** çıxar
  (`happyPathBankCharge`, `happyPathReceiptSend`)
- Test inputlarını dəyişənlərə (`validTestOrder`, `testReceiptNo`)
- Yerləşdirmə: senarilərdən əvvəl / global var / t.Run closure başında

### 8. Lesser-known test alətləri
```bash
go test -run=TestOrderManager_Process ./...   # ad filtri (regex)
go test -timeout=30s ./...                    # panik ilə timeout
go test -short ./...                           # t.Short() ilə uzun testləri ötür
```
- **Context veriləndə:** `context.Background()` yox,
  `context.WithTimeout(context.Background(), 1*time.Second)` — testlər həmişə
  bitər (flaky yox)
- **testdata/ direktiv:** Go tooling onu skip edir — fixture faylları orada
- **Fixture generator:** env var switch ilə test eyni zamanda fixture yaradır

### 9. Test-only quruluş (private constructor-lar)
Test-only asılılıqları (timeout, `sendMail` funksiyası, config) **private
constructor** ilə daxil et — public API-nin UX-i pozulmur:
```go
type sendMail func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

func NewEmailer(...) *Emailer { return &Emailer{sendFunc: smtp.SendMail} }
```
Monkey patching-ə ehtiyac qalmır (data race + restore problemi yoxdur).

### 10. Fluent API testi
ElasticSearch kimi fluent API-lər mock-lanmır. İstifadə yeri custom funksiya
tipinə çıxarılır (`historySearchFunc`), member var stub-lanır. Biz fluent API-nin
özünü yox, kodumuzun **cavablarına reaksiyasını** test edirik.

### 11. Concurrent kodun testi
- **Avoidance:** `time.Sleep` QADAĞANDIR — `select { case <-resultCh: ...; case <-time.After(1*time.Second): assert.Fail(...) }`
- **Test latch:** NOOP closure member var (`orderDropped func(*Order)`) —
  production-da boş, testdə nəticə verir (test-induced damage, amma yalnız
  obyektin özünə məhdud)
- **Benchmark:** `b.RunParallel()` — resurs contention ölçür
- Interface implementasiya yoxlaması: `var _ Iface = (*T)(nil)`

### 12. Skip strategiyaları
- init()-də ping → `mysqlAvailable bool` → `t.Skip`
- **Env var (tövsiyə):** `if os.Getenv("MYSQL_HOST") == "" { t.Skip(...) }` —
  clone-dan dərhal sonra bütün testlər keçməli
- `exec.LookPath("mysql")` helper ilə yoxlama

### 13. Test-Induced Damage (testin yaratdığı zərər)
DHH-nin termini — yalnız test üçün kodun/architecture-nin pisləşməsi.
**Təhlükə siqnalları və həlləri:**
- Yalnız test üçün parametrlər → private constructor + member var
- Yalnız test üçün çıxışlar → test latch (NOOP closure)
- Testdə URL/DSN hardcode → asılılıq obyektə çıxarılmalı, DI
- Mock-lar production faylında → `-testonly` generasiya
- **100% coverage hədəfi QADAĞANDIR:** hər testin yazma/işlətmə/baxım qiyməti
  var; `json.Marshal`-ı mock etmək kimi "mümkün" şeylər bəzən heç bir fayda
  gətirmir

### 14. Make it work, make it clean, then maybe make it fast
Yeni davranış əlavə edərkən yalnız o bir işə fokuslan; design/UX/performance
sonra refactor mərhələsində. Eyni vaxtda hər üçünü etmək mümkün deyil.

## Əsas terminlər

- Behavior Coverage (davranış örtüyü)
- Table-Driven Test (cədvəl əsaslı test)
- Test Pyramid (test piramidi)
- Mock / Stub / Test Recorder
- Test Latch (test kilidi)
- Test-Induced Damage (testin yaratdığı zərər)
- Fixture Generator (test datası generatoru)
- Transient Error (keçici xəta) simulyasiyası

## Praktik nəticə

- Unit = struct + public API; private metodlar dolayı yoxlanılır
- TDT + adlandırılmış mock funksiyaları + require/assert sırası
- Mock yoxlaması: yalnız davranış — çağırış sayı/sırası YOX
- `context.WithTimeout` hər testdə; `time.Sleep` əvəzinə `select`
- ~70% davranış coverage; test-in qiyməti vs faydası

## Mənbə

Pages: 148-236 (Chapter 6, Beyond Effective Go Part 2)
