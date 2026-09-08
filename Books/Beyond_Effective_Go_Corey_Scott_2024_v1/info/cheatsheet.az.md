# Beyond Effective Go — Part 2 — Cheat Sheet (Azərbaycanca)

## 1. Design prinsipləri (Ch4)

| Prinsip | Go-dakı tətbiqi |
|---|---|
| Unix Philosophy | kiçik, tək-məsuliyyətli paketlər; birləşdirilə bilən alətlər |
| KISS > DRY | DRY axiomatik DEYİL — təsadüfi oxşarlıq abstraksiya tələb etmir |
| Composition over inheritance | struct embed + interface; miras YOX |
| Accept interfaces, return structs | istifadəçi genişləndirə bilsin; genişlənmə kontraktı interface-də |
| DIP | yüksək qat aşağı qatı DƏQİQ ŞEYL DEYİL, interfeyslə tanır |
| ISP | interfeysi istifadəçisinin EHTİYACINA görə böl |

## 2. Pattern implementasiyaları (Ch4)

```go
// Singleton (lazy, thread-safe):
var once sync.Once; var instance *DB
once.Do(func() { instance = connect() })
// Observer: channel mübadiləsi (sub chan<- Event)
// Factory: func NewRed(cache Cache) *Client {...}  — closure ilə
// Adapter: struct { http.RoundTripper }  — embed edib genişləndir
```

## 3. Code UX üçlüyü (Ch5)

1. **Clarity:** ad = niyyət; funksiya < 1 ekran; nesting ≤ 3-4;
   erkən return
2. **Consistency:** layihə VƏ Go ekosistemi ilə vahid üslub
3. **Predictability:** principle of least astonishment — zero value-ya
   hörmət, error-ları wrap et, API surface-i kiçildən saxla

## 4. Table-driven test kanonu (Ch6)

```go
func TestPay(t *testing.T) {
    t.Parallel()
    tt := []struct{ name string; amount int; want error }{
        {"valid", 100, nil},
        {"zero", 0, ErrAmount},
    }
    for _, tc := range tt {
        tc := tc                       // capture!
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            if err := Pay(tc.amount); !errors.Is(err, tc.want) {
                t.Errorf("Pay(%d): want %v, got %v", tc.amount, tc.want, err)
            }
        })
    }
}
```
- Piramid: ~70% unit, ~28% integration, ~2% E2E
- **Behavior yox implementation:** mock-un DAXİLİNİ (necə çağırıldığını)
  yox, NƏTİCƏNİ yoxla; recorder = strukturu YOX, çağırışları qeyd edir

## 5. Mock vs Stub

- **Stub:** hazır dəyər qaytarır (state tələb etmir)
- **Mock:** gözləntiləri YOXLanIR (call-ların sayı/sırası)
- mockery: `//go:generate mockery` — interfeysdən avtomatik

## 6. Test helper imkanları

```go
t.Cleanup(func() { ... })          // defer-in test versiyası
t.Parallel(); tc := tc             // cədvəl + parallel
ctx, cancel := context.WithTimeout(...)
t.Cleanup(cancel)
// testdata/ qovluğu — fixture faylları; package-level dəyişməz!
```

## 7. Race/concurrency testi

```bash
go test -race -count=10 ./...
```
- Test latch: `wg := sync.WaitGroup{}; wg.Add(1)` — goroutine startını
  gözlə; blocking kanal üçün: ayrı-goroutine-da başlat, timeoutlu select

## 8. Productivlik (Ch7)

```bash
staticcheck ./...          # SA* class-ları; go vet-dən geniş
gopls                      # IDE dili serveri — rename/findref
dlv debug                  # breakpoint, qədim print-dən yaxşı
go test -gcflags="all=-N -l" -coverprofile=c.out && go tool cover -html=c.out
```
- Kiçik commit + interactive rebase (bir məntiq = bir commit)
- Conventional Commits: feat:/fix:/docs: ...

## 9. FP və function pattern-lər (Ch8)

```go
// Middleware:
func Logging(next http.Handler) http.Handler { ... }
// Functional options (az option üçün sadə variantlar üstündür):
func NewServer(opts ...Option) *Server
// Decorator: func(http.Handler) http.Handler zənciri
// Futures:
func Fetch() <-chan Result { ch := make(chan Result, 1); go func(){ ch <- work() }(); return ch }
// Empty struct: struct{} — 0 bayt; signal-only kanal: chan struct{}
// noCopy:
type T struct{}; func (*T) Lock() {}   // go vet copylock tapır
```

## 10. Metaproqramlaşdırma (Ch9)

```go
// (1) API:
req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
// (2) exec:
out, err := exec.CommandContext(ctx, "git", "log", "--oneline").CombinedOutput()
cmd.Stdin, cmd.Stdout = &buf, os.Stdout       // pipe-lar
// (3) codegen:
fset := token.NewFileSet(); f, _ := parser.ParseFile(fset, "in.go", nil, 0)
ast.Inspect(f, func(n ast.Node) bool { ... })
tmpl.Execute(outFile, data)                    // text/template
```
- Meyar: **alət vaxtı < qənaət vaxtı** — yoxsa YAZMA

## 11. Tələsiklər

| Tələsik | Düzgünü |
|---|---|
| Erkən DRY | oxşarlıq 2 dəfə görünəndə saxla, 3-cüdə abstrakt et |
| Implementation testi | davranış + sərhədlər |
| Singleton hər yerdə | explicit dependency daha aydın |
| Currying Go-da | adi parametr funksiyası kifayətdir |
| test-də os.Exit | t.Fatal/Fatalf |
| deep nesting | early return, extract funksiya |
| API design: konkret type qəbul | interface qəbul et, konkret qaytar |
