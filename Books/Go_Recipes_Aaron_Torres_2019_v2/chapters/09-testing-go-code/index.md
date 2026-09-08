# Chapter 9 — Тестирование Go-кода (Go kodunun testi)

## Bu chapter nədən bəhs edir?

Standart kitabxana ilə mocking (Patch/Restore), mockgen ilə interfeys
mock-ları, table-driven testlər və coverage, üçüncü tərəf alətlər (gocov,
goconvey) və davranış testi (BDD — godog/Gherkin).

## Əsas fikirlər

### 1. Mocking — standart kitabxana ilə
**Nədir:** Interfeysin test versiyasını yazmaq + Patch/Restore texnikası ilə
qlobal funksiya dəyişənlərinin runtime əvəzi.

**Kitabdan kod nümunəsi:**
```go
// İnterfeş:
type DoStuffer interface {
    DoStuff(input string) error
}

// MOCK — closure sahəli struct:
type MockDoStuffer struct {
    MockDoStuff func(input string) error   // runtime-da doldurulur
}
func (m *MockDoStuffer) DoStuff(input string) error {
    if m.MockDoStuff != nil {
        return m.MockDoStuff(input)        // mock ediləibsə çağır
    }
    return nil                             // default hal
}

// PATCH/RESTORE — reflection ilə dəyişən əvəzi:
func Patch(dest, value interface{}) Restorer {
    destv := reflect.ValueOf(dest).Elem()     // pointer-in dəyəri
    oldv := reflect.New(destv.Type()).Elem()
    oldv.Set(destv)                            // köhnə dəyəri saxla
    valuev := reflect.ValueOf(value)
    destv.Set(valuev)                           // yeni dəyəri qoy
    return func() {                             // restorer qaytar
        destv.Set(oldv)
    }
}

// Patch oluna bilən funksiya (VAR kimi!):
var ThrowError = func() error {
    return errors.New("always fails")
}

// TEST — table-driven + hər iki texnika:
func TestDoSomeStuff(t *testing.T) {
    tests := []struct {
        name      string
        DoStuff   error
        ThrowErr  error
        wantErr   bool
    }{
        {"base-case", nil, nil, false},
        {"DoStuff error", errors.New("failed"), nil, true},
        {"ThrowError error", nil, errors.New("failed"), true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            d := MockDoStuffer{}
            d.MockDoStuff = func(string) error { return tt.DoStuff }
            defer Patch(&ThrowError, func() error { return tt.ThrowErr }).Restore()
            if err := DoSomeStuff(&d); (err != nil) != tt.wantErr {
                t.Errorf("DoSomeStuff() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

**Sub-kod izahı:**
- Mock struct closure sahə ilə → davranış test-dən idarə olunur
- `Patch(&var, val)` → var-ın qöhnə dəyərini yadda saxlayır, yenisi qoyur;
  `defer .Restore()` → test sonunda geri qaytarır
- Patch yalnız `var A = func()` üçün işləyir — `func A()` deyil!
- Xarici paket üçün: `var packageDoSomething = package.DoSomething` wrapper

### 2. mockgen (gomock) ilə mock generasiyası
**Nədir:** İnterfeyslərdən avtomatik mock kodu generasiya edən alət —
 gözləntilər (EXPECT), çağırış sayları, Return dəyərləri.

**Kitabdan kod nümunəsi:**
```bash
mockgen -destination internal/mocks.go -package internal \
    github.com/.../mockgen GetSetter
```
```go
type GetSetter interface {
    Set(key, val string) error
    Get(key string) (string, error)
}

// TEST — gomock controller:
func TestController_Set(t *testing.T) {
    tests := []struct {
        name          string
        getReturnVal  string
        getReturnErr  error
        setReturnErr  error
        wantErr       bool
    }{
        {"get error", "value", errors.New("failed"), nil, true},
        {"value match", "value", nil, nil, false},
        {"no errors", "not set", nil, nil, false},
        {"set error", "not set", nil, errors.New("failed"), true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()          // gözləntilər doğrulanır
            mockGetSetter := internal.NewMockGetSetter(ctrl)
            mockGetSetter.EXPECT().Get("key").AnyTimes().
                Return(tt.getReturnVal, tt.getReturnErr)
            mockGetSetter.EXPECT().Set("key", gomock.Any()).AnyTimes().
                Return(tt.setReturnErr)
            c := &Controller{GetSetter: mockGetSetter}  // mock inject
            if err := c.GetThenSet("key", "value"); (err != nil) != tt.wantErr {
                t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

// Do + Return — side effect ilə:
mockGetSetter.EXPECT().Get("key").Do(func(key string) {
    k = key                       // arqumenti yaxala
}).Return("", nil)
// gomock.Any() → istənilən arqument; AnyTimes() → 0+ dəfə çağırıla bilər
```

**Sub-kod izahı:**
- `EXPECT()` → hansı metod, hansı arqumentlərlə, nə qaytarmalı
- `.Times(n)` / `.AnyTimes()` → çağırış sayı gözləntisi (default: dəqiq 1)
- `ctrl.Finish()` → bütün gözləntilər yerinə yetirilməyibsə test FALL edir
- `internal/` paketinə generasiya → mock-lar xarici istifadəçiyə görünmur

### 3. Table-driven testlər və coverage
**Nədir:** Test halları struct slice kimi — az kodla çox hal; `gotests`
skelet generasiyası; `-cover` metrikası.

**Kitabdan kod nümunəsi:**
```go
// gotests -all -w → avtomatik skelet:
func TestCoverage(t *testing.T) {
    type args struct { condition bool }
    tests := []struct {
        name    string
        args    args
        wantErr bool
    }{
        // TODO: Add test cases.
        {"no condition", args{true}, true},
        {"condition", args{false}, false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := Coverage(tt.args.condition); (err != nil) != tt.wantErr {
                t.Errorf("Coverage() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```
```bash
go test -cover                          # coverage: 100.0% of statements
go test -coverprofile=cover.out
go tool cover -html=cover.out -o coverage.html   # qrafik hesabat
```

**Sub-kod izahı:**
- `t.Run(tt.name, ...)` → hər hal ayrıca subtest kimi görünür
- 100% coverage hədəfi məcburi deyil — xarici səthin hərtərəfli testi
  daha dəyərlidir (məs. gob.Encode xətasını süni yaratmaq çətindir)

### 4. Üçüncü tərəf alətlər (gocov, goconvey)
**Nədir:** gocov → funksiya-səviyyəli coverage; goconvey → assertion
kİTABXANASI + test runner + veb UI.

**Kitabdan kod nümunəsi:**
```go
import (. "github.com/smartystreets/goconvey/convey")

func Test_example(t *testing.T) {
    tests := []struct{ name string }{{"base-case"}}
    for _, tt := range tests {
        Convey(tt.name, t, func() {
            res := example()
            So(res, ShouldBeNil)                    // assertion!
            So(example2(), ShouldBeGreaterThanOrEqualTo, 1)
        })
    }
}
```
```bash
gocov test | gocov report    # funksiya üzrə coverage cədvəli
goconvey                     # veb UI + avtomatik yenidən işə salma
```

**Sub-kod izahı:**
- `So(actual, Should..., expected)` → sətir-xəbərdarlıqlı assertion
- `Convey` blokları t.Run-dan FƏRQLİ: nested bloklar SIRAYLA yenidən icra
  olunur (outer setup üçün) — t.Run isə ardıcıl, təkrarsızdır
- goconvey-də notification: kod yadda saxlananda testlər avtomatik işə düşür

### 5. Davranış testi (BDD — godog)
**Nədir:** Gherkin dili ilə (ingiliscə addımlar) sınaq ssenariləri —
qeyri-tezniki insanlar üçün oxunaqlı; godog addımları Go-da realizə edir.

**Kitabdan kod nümunəsi:**
```gherkin
Feature: Bad Method
Scenario: Good request
  Given we create a HandlerRequest payload with:
    | reader |
    | coder  |
    | other  |
  And we POST the HandlerRequest to /hello
  Then the response code should be 200
  And the response body should be:
    | BDD testing reader |
    | BDD testing coder  |
    | BDD testing other  |
```
```go
// Addım realizasiyaları (feature faylına regex ilə bağlanır):
func weCreateAHandlerRequestPayloadWith(arg1 *gherkin.DataTable) error {
    for _, row := range arg1.Rows {
        h := HandlerRequest{Name: row.Cells[0].Value}
        payloads = append(payloads, h)
    }
    return nil
}

func wePOSTTheHandlerRequestToHello() error {
    for _, p := range payloads {
        v, _ := json.Marshal(p)
        w := httptest.NewRecorder()                    // ResponseWriter mock!
        r := httptest.NewRequest("POST", "/hello", bytes.NewBuffer(v))
        Handler(w, r)                                  // birbaşa handler çağır
        resps = append(resps, w)
    }
    return nil
}

func theResponseCodeShouldBe(arg1 int) error {
    for _, r := range resps {
        if got, want := r.Code, arg1; got != want {
            return fmt.Errorf("got: %d; want %d", got, want)
        }
    }
    return nil
}

// Qeydiyyat — regex → funksiya:
func FeatureContext(s *godog.Suite) {
    s.Step(`^we create a HandlerRequest payload with:$`, weCreateAHandlerRequestPayloadWith)
    s.Step(`^we POST the HandlerRequest to /hello$`, wePOSTTheHandlerRequestToHello)
    s.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
    s.Step(`^the response body should be:$`, theResponseBodyShouldBe)
}
```

**Sub-kod izahı:**
- `httptest.NewRecorder()` → ResponseWriter-in yaddaş mock-u (body+code)
- `httptest.NewRequest(...)` → real HTTP sorğusu YARADAN server
- DataTable (`| col |`) → ssenari data cədvəli; addım funksiyaları onu
  arbitrary Go koduna çevirir
- Godog məhdudiyyətləri: examples/context ötürülməsi yoxdur; öz test
  runner-i var → `go test -cover` ilə birləşmir

## Əsas terminlər

- Mock / Stub / Test Double (sınaq ikiqatı)
- Closure-based Mock
- Patch / Restore (patch/geri qaytarma)
- gomock / mockgen / EXPECT
- Table-driven Tests (cədvəlli testlər)
- t.Run (subtestlər)
- Test Coverage (test örtüklüyü) / coverprofile
- gotests (test skeleti generatoru)
- Assertion (iddia) / Convey / So
- BDD (Behavior-Driven Development) / Gherkin / Cucumber / godog
- httptest.NewRecorder / httptest.NewRequest

## Praktik nəticə

- İnterfeyslərə görə kod yaz → mock asan; branch-ləri azalt → test sadə
- Mock: closure sahə strukturu; qlobal funksiya üçün `var A = func()` + Patch
- Çox interfeys → mockgen; internal paketə generasiya et
- Table-driven + gotests → coverage sürətlə artır
- gocov funksiya-səviyyə boşluqları göstərir; goconvey assertion + UI verir
- BDD: ssenarilər İNSAN dilində; httptest ilə real handler-ları network-süz
  sına

## Mənbə

Pages: 297-323 (Chapter 9, Go Programming Cookbook 2nd ed)
