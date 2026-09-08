# Learn Go with Pocket-Sized Projects — Terminologiya (AZ)

| Termin | Azərbaycanca | İzah |
|---|---|---|
| Example test | nümunə testi | stdout çıxışını `// Output:` ilə yoxlayır; dokumentasiya |
| Table-driven test | cədvəl testi | test data map/slice + t.Run subtestlər |
| Custom type | xüsusi tip | `type language string` — nominal typing |
| iota | iota | enum üçün ardıcıl sabit sayğacı |
| defer | defer | funksiya sonunda icra (fayl bağlama, unlock) |
| Comma-ok idiom | v, ok forma | `v, ok := map[k]` — açar varlıq yoxlaması |
| Stringer | Stringer | `String() string` metodu — çap formatı |
| Sentinel error | sentinel xəta | `const ErrX = tipX("...")` + errors.Is |
| Error wrapping | xəta sarılması | `fmt.Errorf("%w", err)` — zəncir |
| Variadic | variadik | `args ...any` — istənilən sayda arqument |
| Functional options | funksional seçənəklər | `New(x, WithOutput(w))` |
| io.Writer / io.Reader | yazı/oxu interfeysi | test/mock üçün injektə olunan asılılıq |
| json.NewDecoder | JSON dekoderi | stream-əsaslı dekodlaşdırma |
| Fixed-point | sabit nöqtə | int subunits + precision (float əvəzi) |
| ISO-4217 | valyuta kodu | USD, EUR + kəsr dəqiqliyi |
| struct tag | struktur teqi | `json:"name"` / `xml:"Cube"` |
| httptest | httptest | mock HTTP server (ts.URL) |
| Dependency injection | asılılıq inyeksiyası | asılılıq parametr/ininterfeys kimi ötürülür |
| Generics | generiklər | `[T any]` / `[K comparable, V any]` |
| Type constraint | tip məhdudiyyəti | any, comparable, öz interfeysi |
| Goroutine | gorutin | `go f()` — yüngül paralel vahid |
| Channel | kanal | goroutine-lərarası kommunikasiya; close/istiqamət |
| select | select | çoxkanal gözləmə; quit pattern |
| Data race | məlumat yarışı | paralel map yazışı; `go test -race` |
| Mutex / RWMutex | mutex | Lock/Unlock — serializasiya; defer ilə |
| WaitGroup | gözləmə qrupu | Add/Done/Wait — goroutine bitməsi |
| errgroup | errgroup | xəta toplayan WaitGroup (x/sync) |
| TTL | yaşam müddəti | entry expires time.Now() + ttl |
| Loop variable shadowing | dəyişən kölgəsi | `p := p` — closure təhlükəsizliyi (Go ≤1.22) |
| http.ServeMux | router | "POST /game" formatında marşrut |
| Repository pattern | depo nümunəsi | saxlanc abstraksiyası (in-memory → DB) |
| Linked list | zəncir siyahı | previousStep pointer — path strukturu |
| Protobuf | protobuf | binar serializasiya; versioning |
| gRPC | gRPC | RPC framework; status kodları (3, 13...) |
| go:generate | generasiya direktivi | protoc/minimock çağırışları |
| Typed error | tipli xəta | struct{field,reason} + errors.As |
| context.Context | kontekst | deadline/Done/Value; hər çağırışa öz dəsti |
| Integration test | inteqrasiya testi | real server + client; -short ilə skip |
| go:embed | daxiletmə | binary-yə fayl məzmunu |
| text/template | şablon | {{ . }} {{ range }} {{ if }} FuncMap define |
| Golden file | qızıl fayl | referens çıxış müqayisəsi |
| WebAssembly (Wasm) | VebAssembly | GOOS=js GOARCH=wasm; brauzerdə Go |
| js.FuncOf | JS körpüsü | Go funksiyasını JS-ə açır |
| TinyGo | TinyGo | mikrokontroler Go-ları; machine paketi |
| Build tag | build teqi | //go:build arduino_nano33 — platforma seçimi |
| Flashing | flesh etmə | mikrokontrolerə proqram yazma |
| Zero value | sıfır dəyər | nil / 0 / "" / false — init edilməyənin dəyəri |
