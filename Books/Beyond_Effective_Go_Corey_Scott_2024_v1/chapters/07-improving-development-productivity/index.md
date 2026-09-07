# Chapter 7 — Improving your Development Productivity (səh. 237-276)

## Bu chapter nədən bəhs edir?

Proqramlaşdırmanın "periferik" tərəfləri: alətlərin mükəmməl işlədilməsi,
təkrarlanan işlərin avtomatlaşdırılması, IDE bacarığı, git axını və kiçik
dəyişikliklərlə işləmə. Fəlsəfə: **"Be lazy"** — daha az səylə daha çox dəyər.

## Əsas fikirlər

### 1. Be Lazy — az səylə çox dəyər
Tənbəllik = işi görməmək deyil, **minimum səylə görmək**. Nümunələr:
- **Less is often better:** az istifadə olunan feature-ləri silmək təklif et —
  baxım qiyməti (kod anlaşılması, testlər, asılılıq update-ləri) davamlı ödənilir
- **Dokumentasiya yazmaq vaxt qənaətidir:** tez-tez soruşulan sual → cavab
  dokumentasiyada olsun (API doc, RFC). "Bu sualı kimisə yenidən soruşacaqmı?"
  testi ilə qərar verilir

### 2. Be Observant — təkrara qızğın ol
Təkrarlanan/yavaş/qəliz işə görülür → hirsini həllə çevir:
- `git push -u origin HEAD` → `pr` alias/script-inə
- Race detection → git pre-commit hook
- Adi səhvlər → linter qayda / template / avtomatlaşdırma

### 3. Clean as you go
McDonald's-dan gələn qayda: kod bazası mükəmməl olmayacaq — göründüyü yerdə
kiçik təmizliklər et (boy scout rule).

### 4. Be Introspective + Adventurous
- Aylıq özünü yoxla: nə yavaşladır? nə yaxşı işləyir?
- **Meticulous caution ən yavaş yoldur.** Riski idarə etməyin yolu ehtiyatlı
  olmaq deyil, **təhlükəsizlik şəbəkəsi** qurmaqdır: automated tests + plan
  (RFC) — planın dərinliyi problemin ölçüsünə mütənasib olmalıdır

### 5. Master Your Tooling — skriptlər
Müəllifin bash skript dəsti (makefile-ə çevrilə bilər):

**Coverage hesablama:**
```bash
package-coverage -a -i $COVERAGE_EXCLUDE -m 70 -prefix $BASE_PKG ${@:2} $PKG_DIR
# -m 70 → minimum 70% coverage limiti
```

**HTML coverage:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Race testləri:**
```bash
go test -race $(go list $1 | grep -v /vendor)
```

**Format + importlar (goclean):**
```bash
gofmt -w -s -l $DIRS
goimports -w -l $DIRS
```

**Generate + format (gogen):**
```bash
go generate $1
goclean $1
```

**Dependency graph:**
```bash
godepgraph -s -p "$EXCLUSIONS" -o "$BASE_PKG" $BASE_PKG/${PKG} |
  sed "s/$BASE_PKG_DELIMITED//g" | dot -Gsplines=true -Tpng -o "$DEST_FILE"
```

**Ship it (ən qiymətli):** bütün paketlər üzrə clean → race testlər → coverage
yoxlaması → `read -p "Ship it? (y/n)"` → `git add . && git commit && git push`.
Dəyişməmiş kod skip olunur — vaxt itmir.

### 6. Master your IDE (Goland)
**Context actions:** cursor mövqeyinə görə təkliflər (struct field-lərini doldur,
issue quick-fix). **F2** — Next Highlighted Error.
**Test runner + External tools:** klaviatura qısayolu ilə testlər + format skripti
birlikdə; kod generatorlar IDE-dən çağırılır.
**Live templates:** məs. `cctx` →
```go
ctx, cancel := context.WithTimeout(context.Background(), $TIMEOUT$ * time.Second)
defer cancel()
```
TDT skeleton template-i — bütün 5-addımlı test strukturu bir abzasla.
**Vacib shortcuts:** ⇧⌘O (fayl aç + :line), ⌘B (declaration/usages), ⌘N
(generate test/interface), ⇧⌘N (scratch file).
**Plugins:** IdeaVim (klaviatura mədəniyyəti), Key Promoter (mouse-dan qaçmağı
öyrədir), GitHub Copilot (tekrar construct-lar: error formatları, test
assertion-ları).

### 7. Git axını skriptləri
**newfeature** (təmiz başlanğıc):
```bash
IS_CLEAN=$(git status --short | wc -l)
if [ ${IS_CLEAN} -ne 0 ]; then git stash; fi
git checkout master
git pull origin master
git checkout -b $1
if [ ${IS_CLEAN} -ne 0 ]; then git stash pop; fi
```

**switch** (branch-lar arası + rebase):
```bash
git checkout $1
git rebase master
```
— bitmiş sayıb PR açanda rebase stressi yoxdur.

**Aliaslar:**
```bash
alias ga='git add ';  alias gc='git commit ';  alias gd='git diff '
alias gco='git checkout ';  alias gp='git pull ';  alias grm='git rebase master '
alias swm='sw master'
alias got='git ';  alias gut='git '   # typos tutulur
```

### 8. Make Small Changes
CD ideologiyası: **"aqrıdan bir şey varsa, onu daha tez-tez et"** — böyük PR-lər
böyük risk, böyük review yükü. Kiçik PR-lər üçün:
- Feature-i public API-yə BAĞLAMADAN merge et
- Implementation dəyişikliyi feature flag (off) ilə qoru
- Kod merge olunub, amma istifadə olunmur → risk yoxdur

## Əsas terminlər

- Continuous Delivery (fasiləsiz çatdırılma)
- Feature Flag (xüsusiyyət bayrağı)
- Pre-commit Hook (commit-öncəsi qarma)
- Live Template (canlı şablon)
- Dependency Graph (asılılıq qrafı)
- Boy Scout Rule (təmiz-təmizlə qaydası)

## Praktik nəticə

- Gündəlik əməliyyatları skript/alias-a çevir; "ship it" prosesini vahid et
- IDE-ni klaviatura ilə işlət: context actions, live templates, test runner
- Hər ay introspeksiya: nə yavaşladır?
- Təkrarlanan ağrılı işi (PR, rebase, review) KİÇİK və TEZ-TEZ et

## Mənbə

Pages: 237-276 (Chapter 7, Beyond Effective Go Part 2)
