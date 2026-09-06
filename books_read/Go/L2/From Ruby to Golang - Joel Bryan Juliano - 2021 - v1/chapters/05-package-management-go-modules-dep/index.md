# Chapter 5 — Package Management: Go Modules və dep (book səh. 82-92)

## Bu chapter nədən bəhs edir?

RubyGems/Bundler/Gemfile dünyasının Go ekvivalentləri: paket paylaşımı (GOPATH, import), Ruby module+include → Go package+import xəritələnməsi, Go Modules (go.mod — manual, source-dən, dep-dən generasiya), require/indirect asılılıqlar və dep paket meneceri (Gopkg.toml/lock, TOML Array of Tables).

---

## Əsas fikirlər

### 1. Ekosistem xəritəsi

| Ruby | Go |
|------|-----|
| RubyGems (install) | go get |
| Gemfile + Bundler | go.mod + Go Modules (və ya legacy: dep) |
| module + include | package + import |
| Gemfile.lock | go.sum (müasir) / Gopkg.lock (dep) |

### 2. Paket paylaşımı — GOPATH + import

Go funksiyaları paketlərə yığılıb paylanır. **import** = namespace-i yüklə (Ruby `include` analoqu). GOPATH — paketlərin/kitabxanaların/binarilərin saxlandığı mühit dəyişəni; uzaq paketlər `$GOPATH/src`-də yaşayır.

```go
import (
    "time"
    "github.com/foo/greeter"
    weather "github.com/bar/climate"   // CUSTOM AD — alias
)

greeter.SayHello("John")
weather.SayWeather(time.Now())
```

### 3. Ruby module+include → Go package+import

```ruby
# Ruby:
module Coffee
  def types; %w[espresso cappuccino macchiato]; end
end
class Drink
  include Coffee
  def menu; puts types; end
end
```

```go
// $GOHOME/src/local/coffee.go:
package coffee

func Types() []string {    // BÖYÜK hərf — public!
    return []string{"espresso", "cappuccino", "macchiato"}
}

// main proqramı:
import "local/coffee"

for _, val := range coffee.Types() { fmt.Println(val) }
```

**Dönüşüm:** module → package faylı; include → import; metodlar → BÖYÜK hərflə export olunmuş funksiyalar.

### 4. Go Modules (Go 1.11+)

Mühit: `go version` ≥ 1.11; `GO111MODULE=on|off|auto` (default auto: layihə $GOPATH/src-də DEYİL + go.mod mövcuddursa ON).

Modul sistemi GOPATH-dən asılılıq aradan qaldırır — paket UZAQ resurs olmalıdır (github.com/foo/coffee, tag v0.0.1).

**Manual go.mod:**

```
module menu

require github.com/foo/coffee v0.0.1
```

Versiyalar: semver prefiks (v1 = həmin prefikslə son tag), müqayisə operatorları (`<`, `<=`, `>`, `>=`), `latest`. Sonra: `go run .` — avtomatik download; və ya `go get -u`. Modulu public etmək üçün `module` sətri uzaq yola dəyişilir: `module github.com/foo/menu`.

**Source-dən generasiya (`go mod -sync`):** import-ları oxuyur. VCS tələb edir (git, bzr...). **Subpackage tələsi:** `github.com/quux/foo/baz` importu, repo `github.com/quux/foo` tələb etmirsə — həllər:
1. `baz`-ı ayrı repoya çıxar (github.com/quux/baz)
2. go.mod-da `// indirect` kimi əlavə et:

```
require (
    github.com/quux/foo v0.0.1 // indirect
)
```

Versiya tapmaq: `curl -s https://api.github.com/repos/quux/foo/releases/latest | jq -r .tag_name`. Çox indirect → paketin öz reposunu düşün.

**dep-dən migrasiya:** `dep ensure` → `go get -u` — Gopkg.lock-dan go.mod-a köçürür.

**Refresh:** modullar `$GOPATH/src/mod`-da READ-ONLY — silmək üçün yazıla bilən et + rekursiv sil.

### 5. dep (legacy paket meneceri)

Production-da sabit, Go Modules-ə çevrilməmiş köhnə layihələr üçün. `dep init`:
1. Gopkg.toml (konfiq)
2. Gopkg.lock (versiya kilidi)
3. vendor/ qovluğuna download

**Gopkg.toml — TOML Array of Tables:**

```toml
[[constraint]]
name = "github.com/foo/bar"
version = "0.0.1"        # və ya: branch = "baz"
source = "github.com/quux/baz"   # FORK-ə yönləndirmə
```

JSON/Ruby ekvivalenti: `{"constraint": [{"name": ..., "version": ...}]}`. Dəyişiklikdən sonra `dep ensure` — vendor yenilənir. **TOML:** YAML kimi, amma aydın semantika ön planda.

**Status:** Go Modules standartlaşdı — dep köhnəldi (obsolet).

---

## əsas terminlər

| Termin | İzah |
|--------|------|
| GOPATH | Paket/kitabxana/binari/derleyici mühit yolu; $GOPATH/src |
| go get | RubyGems-ə bənzər paket çəkmə (`-u` update) |
| GO111MODULE | on/off/auto — modul rejimi keçidi |
| go.mod | Gemfile analoqu: module adı + require-lər |
| semver prefiks/müqayisə | v1, latest, <, <=, >, >= versiya qaydaları |
| `// indirect` | Sub-paket asılılığının modul faylında işarəsi |
| go mod -sync | Import-lardan go.mod generasiyası |
| Gopkg.toml/.lock | dep konfiq/kilid; TOML Array of Tables |
| [[constraint]] | dep-də paket qeydi (name/version/branch/source) |
| vendor/ | Asılılıqların layihə daxilində nüsxəsi |

---

## Praktik nəticə

1. **Yeni layihə = Go Modules:** `go mod init github.com/foo/menu` + require-lər; GOPATH-dən asılılıq YOX.
2. **Modul public-etmə qaydası:** `module` sətri uzaq yola + release tag-lər (v0.0.1...).
3. **Subpackage xətası:** ayrıca repo YOXSA `// indirect` + semver tap (curl+jq).
4. **Legacy dep layihəsi:** `dep ensure` saxla, amma `go get -u` ilə go.mod-a miqrasiya mümkündür.
5. **Custom import adı:** `weather "github.com/bar/climate"` — ad toqquşması həlli.
6. **Export = böyük hərf:** paket funksiyaları `Types()` kimi — kiçik hərf yalnız paket daxili.

---

## Mənbə

- Kitab: *From Ruby to Golang* — Joel Bryan Juliano, Leanpub, 2021
- Chapter 5: Package Management, book səh. 82-92
- PDF səhifələri: 88-98
- İstinadlar: golang.org/doc/code.html (ImportPaths), TOML spec
