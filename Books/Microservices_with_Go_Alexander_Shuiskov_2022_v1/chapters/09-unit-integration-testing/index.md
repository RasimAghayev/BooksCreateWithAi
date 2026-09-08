# Chapter 9 — Unit and Integration Testing (səh. 163-192)

## Bu fəsil nədən bəhs edir?

Go testing əsasları (table-driven, subtests, skipping), mocking (manual vs
gomock), metadata controller unit testləri və 3 servisin tam integration testi.

## Əsas fikirlər

### 1. Go testing əsasları
- `example.go` ↔ `example_test.go` konvensiyası; `func TestXxx(t *testing.T)`
- `go test` — qovluqdakı bütün testlər

**Table-driven tests:** inputlar cədvəl kimi, yoxlama məntiqi tək dövr:
```go
tests := []struct{ a, b, want int }{{1,2,3}, {-1,-2,-3}, ...}
for _, tt := range tests {
    assert.Equal(t, tt.want, Add(tt.a, tt.b), ...)
}
```
- Təkrarsızlıq, oxunaqlıq; testify (`assert.Equal`) və ya standart
  `t.Errorf` ilə müqayisə

**Subtests:** `t.Run(name, func(t){...})`
- `-v` ilə hər case ayrı görünür; `-run` ilə tək case işə düşür
- `t.Parallel()` — yavaş funksiyalar paralel icra (dövr dəyişənini lokal kopyala!)

**Skipping:** `t.Skip("səbəb")` + şərt
- `testing.Short()` + `go test -test.short` — yavaş testləri gündəlik işlətdən
  çıxarır

### 2. Mocking (imitasiya)
Mock — komponentin saxta versiyası; test üçün müəyyən cavablar qaytarır.

**Manual mock:**
```go
type mockMetadataRepository struct {
    returnRes *model.Metadata; returnErr error
}
func (m *...) setReturnValues(res, err)  // cavabı proqramlaşdır
```
- Plus: kitabxanasız; Minus: vaxt aparır, interfeys dəyişsə yenilə, az funksiya

**gomock (mockgen):**
```go
mockgen -package=repository -source=metadata/internal/controller/metadata/controller.go
// → gen/mock/metadata/repository/repository.go
```
```go
ctrl := gomock.NewController(t); defer ctrl.Finish()
m := gen.NewMockmetadataRepository(ctrl)
m.EXPECT().Get(ctx, id).Return(nil, repository.ErrNotFound).Times(1)
```
- `EXPECT()` ilə gözlənilən çağırışlar; `Times(n)` — çağırış sayı yoxlanılır
- Generasiya = konsistensiya + boilerplate yox + geniş funksiyalar → gomock seçilir

### 3. Controller unit testi (metadata)
3 test case (cədvəl): not found (ErrNotFound mapping), unexpected error,
success:
```go
t.Run(tt.name, func(t *testing.T) {
    ctrl := gomock.NewController(t); defer ctrl.Finish()
    repoMock := gen.NewMockmetadataRepository(ctrl)
    c := New(repoMock)
    repoMock.EXPECT().Get(ctx, id).Return(tt.expRepoRes, tt.expRepoErr)
    res, err := c.Get(ctx, id)
    assert.Equal(t, tt.wantRes, res, tt.name)
    assert.Equal(t, tt.wantErr, err, tt.name)
})
```

### 4. Integration testlər
- Unit = hissələr; **Integration = hissələrin BİRLİKŞƏ işi**
- **Black box testing** — yalnız publik API ilə (daxili funksiyalar YOX)
- Struktur: **Setup** (server+client instansları) → **Operations+verify** →
  **Teardown** (GracefulStop, Close)

**8-addımlı test planı:** PutMetadata → GetMetadata (müqayisə) →
GetMovieDetails (yalnız metadata) → PutRating(5) → GetAggregatedRating(5) →
PutRating(1) → GetAggregatedRating((5+1)/2=3) → GetMovieDetails (rating=3)

**Kod detalları:**
- `srv.Serve(l)` goroutine içində — bloklamasın; `defer srv.GracefulStop()`
- 6 komponent: 3 server + 3 client; in-memory registry/repo istifadə
- Testutil paketləri: `NewTestMetadataGRPCServer()` və s.

**Persistent DB ilə integration test:** staging-də işlə! Random ID-lər
(`uuid.New()`) + test sonunda cleanup — istifadəçi datası ilə qarışmasın

### 5. Best practices
**Köməkçi mesajlar:**
```go
t.Errorf("got %v, want %v", got, want)               // əvvəl ACTUAL, sonra EXPECTED (Go konvensiyası!)
t.Errorf("YourFunc(%v) = %v, want %v", in, got, want) // funksiya + input da göstər
```
**Fatal-dan qaçın:** `t.Fatalf` testi DAYANDIRIR → az test icra olunur →
az məlumat. `t.Errorf` + loop-da `continue` davam etdirir.

**cmp kitabxanası:** struct pointer-ları `reflect.DeepEqual` ilə yox, `cmp.Diff`
ilə müqayisə — insan-oxunaqlı diff output:
```
-  Title: "Title"
+  IPAddress: "The Title"
```
- `cmpopts.IgnoreUnexported(gen.Metadata{})` — generasiya olunmuş private
  sahələri (protoimpl) nəzərə almasın

## Termindirmə (AZ)
- Unit Test — Vahid Testi (tək komponentin testi)
- Integration Test — İnteqrasiya Testi (komponentlərin birgə testi)
- Mock — İmitasiya (saxta komponent)
- Table-driven Test — Cədvəlli Test
- Black Box Testing — Qara Qutu Testi
- Teardown — Sökmə (test sonrası təmizlik)

## Kviz sualları
1. `t.Parallel()` nə vaxt faydalıdır? (Yavaş test caseləri paralel işə salmaqla
   ümumi vaxtı azaldır)
2. gomock-un manual mock-dan üstünlüyü? (Kod generasiyası, Times() çağırış
   sayı yoxlaması, konsistensiya)
3. Test loglarında dəyərlərin konvensional sırası? (Əvvəl actual (got), sonra
   expected (want))
4. Integration test niyə black box olmalıdır? (Yalnız publik API test olunur —
   daxili detallar dəyişsə test xarab olmasın)
