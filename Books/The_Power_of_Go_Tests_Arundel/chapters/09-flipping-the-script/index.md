# Chapter 9 — Flipping the script (Ssenarini Çevirmək)

## Bu fəsil nədən bəhs edir?

testscript paketi (rogpeppe/go-internal): CLI alətlərinin shell-vari skriptlə
testi. exec (proqram işə düşür + exit status 0 assert), stdout/stderr (regex
match), ! neqasiyası, qoşma pravilaları (single quote, double quote literal),
hello.Main delegate + testscript.RunMain (custom binary $PATH-ə), TestMain
tələsi (os.Exit unudulsa hamısı "uğur" olur), coverage (total coverage: 100% —
subprocess), cmp/cmpenv (golden files), exists/grep/-count, txtar formatı
(-- file -- markerləri, qovluq ağacları), stdin (fayldan / stdin stdout
pipeline), cp/mv/mkdir/cd/rm/symlink, shell-dən fərqləri (glob/pipe YOX —
exec sh -c), şərhlər = fazalar (yalnız fail olan fazanın logu), conditions
([exec:sh], [go1.18], [darwin], [!arm64], [unix], skip), env dəyişənləri
(HOME=/no-home, $WORK, $TMPDIR, $exe), Setup+env.Setenv (randomLocalAddr
SERVER_ADDR), & fon rejimi + wait, standalone testscript runner (-v, -e, -work),
shebang #!/usr/bin/env testscript, txtar-c (kataloqu arxivləşdirib -script ilə
birləşdir), issue repro-lar, ! exec go test (fail gözlənilən modullar).

## Əsas fikirlər

### 1. Problem və Həll — "Flip the Script"
Binary testi əllə: go build → binary işə düşür → args ver → exit status yoxla →
output parse et. Çox iş! Əvəzinə: sadələşdirilmiş shell skript dili ilə Go testi.

```go
// skript faylı (testdata/script/hello.txtar):
exec echo 'hello world'
stdout 'hello world\n'
```
```go
// Go testi:
func TestScript(t *testing.T) {
    testscript.Run(t, testscript.Params{
        Dir: "testdata/script",
    })
}
```
Hər .txtar skript = TestScript-in PARALEL SUBTEST-i (ad = fayl adı Extensionsuz:
`TestScript/hello`). Bir skripti tək işə salma: `go test -run TestScript/hello`.
Skriptlər kiçik və 1-2 davranışa fokuslu saxlanılır (çox işləyən skript = çox
işləyən test kimi oxunmaz).

### 2. exec — İşə Salma + Assert
`exec echo 'hello world'` — proqramı arqumentlə işə salır. **Eyni zamanda
ASSERT-dür:** exit status 0 OLMALIDIR. Fail → skript dərhal dayanır (t.Fatal
kimi) — qalan sətirlər ötürülür.

Hər skript ÖZ unikal müvəqqəti work qovluğunda işləyir, sonra avtomatik
silinir. Workdir BOŞ başlayır → reproduce: proqram lazımi fayl YOXDUSA panic
edirsə, skriptdən ilk çağırışda TUTULUR. "Works on my laptop" dəliyi bağlanır.
("What's true of every bug found in the field? It passed all the tests!" — Rich
Hickey)

### 3. stdout/stderr — Regex Assertlər
```go
stdout 'hello world\n'      # output ən azı bu regex-i özündə saxlamalı
stderr 'cat: doesntexist: No such file or directory'
```
Fail output nümunəsi — çap > ilə sətir-sətir, sonra FAIL + fayl:sətir:
```
> exec echo 'hello world'
[stdout]
hello world
> stdout 'goodbye world\n'
FAIL: testdata/script/hello.txtar:2: no match for `goodbye world\n` found in stdout
```

### 4. ! Neqasiyası — "Olmamalı"
```go
! exec cat doesntexist     # cat FAIL ETMƏLİDİR (non-zero exit)
! stdout .                 # stdout BOŞ olmalı (regex . = hər hansı mətn)
! stderr .                 # stderr boş olmalı
```
Klassik kombinasiya — invalid input davranışı:
```go
# With no arguments, fail and print a usage message
! exec hello
! stdout .
stderr 'usage: hello NAME'
# With an argument, print a greeting using that value
exec hello Joumana
stdout 'Hello to you, Joumana'
! stderr .
```
! exec = "proqram iflas etməlidir" — invalid flag/arg istifadəçi tərəfindən
verildikdə konvensional davranışın (error message + non-zero exit) asserti.

### 5. Arqument Keçirmə Qaydaları
- Boşluqla ayrılmış hər söz = AYRI arqument: `exec cat data file.txt` → İKİ fayl!
- Single quote GRUPLAYIR: `exec cat 'data file.txt'` → bir arqument
- Single quote escape: `''` → `'`: `exec echo 'Here''s how...'` → Here's how...
- **Double quote XÜSUSİ DEYİL** — literal çap olunur! `exec echo "foo"` →
  "foo" (dırnaqlarla). Shell refleksi təhlükəli — diqqət!

### 6. Custom Binary — hello.Main Delegate
main() test olunmur → bütün iş `hello.Main() int`-ə keçir:
```go
func main() {
    status := hello.Main()
    os.Exit(status)
}

func Main() int {
    fmt.Println("hello world")
    return 0
}
```
**TestMain + RunMain ilə qeydiyyat:**
```go
func TestMain(m *testing.M) {
    os.Exit(testscript.RunMain(m, map[string]func() int{
        "hello": hello.Main,
    }))
}
```
RunMain: testlərdən ƏVVƏL "hello" binary-si yaradılır (main→hello.Main çağırır),
müvəqqəti qovluğa qoyulur və $PATH-ə əlavə olunur. Skriptdə: `exec hello` —
sanki əllə compile etmiş kimi. Test bitəndə binary avtomatik silinir.

**Doğrulama:** hello.Main "goodbye world" çap etsə → test FAIL → həqiqətən BİZİM
binary işləyir. return 1 etsə → `[exit status 1] ... unexpected command failure`.

**Tələ:** `os.Exit(...)` ƏTRAFINA alın! Unutsan RunMain-in statusu ignore
olunur → hello həmişə "uğurlu" görünər — gizli bug.

**testscript mənşəyi:** Go tool-un öz testlərindən törəyib — "as complex a
command-line tool as any". github.com/rogpeppe/go-internal.

### 7. Coverage — total coverage
```bash
go test -coverprofile=cover.out
# coverage: 0.0% of statements      ← Go testləri birbaşa 0 işlədir
# total coverage: 100.0% of statements  ← SUBPROCESS-lər daxil! BAXILACAQ rəqəm
```
hello.Main subprocess-də icra olunur — adi coverage onu görmür, total coverage
görür. Happy-path-only skript → `total coverage: 60.0%`. cover.out + go tool
cover / IDE ilə qırmızı yolları gör, skriptləri genişləndir. Coverage öz-özlüyündə
keyfiyyət sübutu deyil (icra ≠ düzgünlük) — amma statistika İTİRMİR.

### 8. cmp / cmpenv — Golden Fayl Muqayisəsi
```go
exec hello
cmp stdout golden.txt
-- golden.txt --
hello world
```
- `-- ad --` marker sətri: ardınca gələn hər şey fayla YAZILIR, workdir-ə qoyulur
- cmp = BÜTÜN fayl bərabərliyi (fər varsa unified diff göstərilir)
- stdout/stderr xüsusi adları: əvvəlki exec-in outputu ilə muqayisə
- `! cmp` = fərq olmalı
- **cmpenv:** golden fayldakı $ENV dəyişənləri EXPAND olunur:
```go
exec echo Running with home directory $HOME
cmpenv stdout golden.txt
-- golden.txt --
Running with home directory $HOME
```
$HOME kənar dəyişəndir — cmpenv flake-i önləyir.

### 9. exists / grep / -count
```go
exec myprog -o results.txt
exists results.txt                    # fayl YARANMALI (məzmun əhəmiyyətsizdirsə)
cmp results.txt golden.txt           # tam məzmun yoxlanmalısa
grep '^hello' results.txt            # hissəvi regex (ən azı 1 match)
grep -count=1 'beep' result.txt       # DƏQİQ say: "have 2 matches for `beep`, want 1" FAIL
```
Workdir avtomatik silinir — troubleshoot üçün `-testwork` flag-i saxlayır və
WORK= yolunu çap edir.

### 10. txtar Formatı — Fayl Ağacları
```
exec cat a.txt b.txt c.txt
-- a.txt --
...
-- b.txt --
...
-- misc/subfolder/b.txt --    # SLAŞ-li ad = qovluq ağacı yaradır!
...
```
Marker: `-- ` ilə başlayır, ` --` ilə bitir, arada fayl adı. Növbəti markere
qədər = məzmun. İstənilən sayda fayl/qovluq. txtar = "text archive" —
bir faylda çoxlu faylın təmsili. **Müstəqil istifadə:** txtar package import
edilir (istifadəçi input-u kimi, alət oxuyucusu kimi). VS Code: vscode-txtar
extension (skript + daxili .go faylları highlight).

### 11. stdin — İstifadəçi Input Simulyasiyası
```go
stdin input.txt
exec greet
stdout 'Hello, John!'
-- input.txt --
John
```
stdin göstərişi: növbəti exec-in standard inputu fayldan gəlir. Çoxsətirli
"söhbət": hər sətir = bir Scan.
```go
-- input.txt --
Kim
barbecue
```
**Pipeline — stdin stdout:**
```go
exec echo hello
stdin stdout          # əvvəlki exec-in OUTPUTU = növbəti input
exec cat
stdout 'hello'
```
Shell pipe-dan fərqi: konkurensi YOXDUR — stdin bütün inputu oxuyur, sonra
davam edir (əvvəlki exec bitib stream bağlanana qədər skript gözləyir).

### 12. Fayl Əməliyyatları
```go
cp a.txt b.txt            # copy (birinci arg stdout/stderr ola bilər)
cp stdout tmp.txt         # exec outputunu fayla yaz
mv a.txt b.txt            # rename
mkdir data
cp a.txt b.txt c.txt data # çoxlu faylı qovluğa
cd data                   # sonrakı exec-lər üçün cwd
rm data                   # REKURSIV (rm -rf kimi)
symlink source -> target  # -> istiqaməti MÜTLƏQ
```

### 13. Shell-dən Fərqlər
- Control flow YOX: loop, funksiya, if yoxdur (fail = bail out)
- exec SHELLSİZ işə salınır → glob (`*`) işləmir, pipe (`|`) işləmir
- Lazımdırsa shell çağır: `exec sh -c 'ls .*'`, `exec sh -c 'echo hello | wc -l'`
- `#` şərh — həmçinin FAZA ayrıcı: fail olanda yalnız CARİ fazanın logu
  çap olunur, əvvəlki fazalardan yalnız şərh + vaxt qalır:
```
# run an existing command: this will succeed
(0.003s)
# try to run a command that doesn't exist
(0.000s)
> exec bogus
[exec: "bogus": executable file not found in $PATH]
```

### 14. Conditions — [kvadrat] Prefiksləri
```go
[exec:sh] exec echo yay, we have a shell     # sh $PATH-də VARSA
[exec:/bin/sh] ...                            # konkret yol VARSA
[go1.16] ...                                  # Go 1.16+ DİLSƏ
[darwin] ...                                  # GOOS şərti
[!arm64] ...                                  # GOARCH deyilsə
[unix] ...                                    # Unix-like (linux, darwin, freebsd...)
[!go1.18] skip                                # şərt + skip kombinasiyası
[linux] skip
```
Şərt true → sətir icra; false → sətir ignore. ! = neqasiya. skip = testi ötür.

### 15. env — Environment Dəyişənləri
```go
env MYVAR=hello                # təyin
env AUTH_TOKEN=                # BOŞ təyin
! exec myprog
stderr 'AUTH_TOKEN must be set'
```
- Skript BOŞ environment ilə başlayır ($PATH + bir neçə predefined xaric)
- HOME=/no-home (çox proqram $HOME gözləyir)
- $WORK = workdir absolute yolu (çapda həmişə literal "$WORK" — deterministik
  output; real dəyir hər run-da fərqlidir)
- $TMPDIR = $WORK/.tmp
- **$exe** — Windows-da `.exe`, başqa platformada boş string: `exec prog$exe`
  → cross-platform skriptlər üçün
- `$` istinadları shell kimi expand olunur: `exec echo $PATH`

### 16. Setup — Dinamik Dəyərlərin Ötürülməsi
Test zamanı bilinən dəyəri (random port!) skriptə ötürmək:
```go
func TestScriptWithExtraEnvVars(t *testing.T) {
    t.Parallel()
    addr := randomLocalAddr(t)                 // port 0 fəndi (keçən fəsil)
    testscript.Run(t, testscript.Params{
        Dir: "testdata/script",
        Setup: func(env *testscript.Env) error {
            env.Setenv("SERVER_ADDR", addr)     // skriptdə $SERVER_ADDR
            return nil
        },
    })
}
```
Setup fayllar EXTRACT-olunduqdan sonra, skript BAŞLAMAZDAN ƏVVƏL işə düşür.
Skriptdə: `exec curl -s --retry-connrefused --retry 1 $SERVER_ADDR`.

### 17. & Fon Rejimi + wait
```go
exec listen $SERVER_ADDR &       # fon: davam et, output buffer olunur
exec curl -s --retry-connrefused --retry 1 $SERVER_ADDR
stdout 'Hello from the Go web server'
wait                             # bütün fon proqramları bitənə qədər gözlə
```
- curl --retry-connrefused: server hələ açılmayıbsa retry (keçən fəslin
  wait-for-success analoqu — klient tərəfdə)
- Script sonunda işləyən fon proqramlar os.Interrupt (Ctrl-C) / os.Kill ilə
  dayandırılır → dayandırma xətası istəmirsənsə wait + server özünü bağlasın
  (demo üçün 1 request-dən sonra shutdown; real serverlər belə etməz)
- stdout/stderr/cmp assertləri wait-dən SONRA buffer-lənmiş outputa tətbiq
  oluna bilər

### 18. Standalone Runner — Test Xaricində
```bash
go install github.com/rogpeppe/go-internal/cmd/testscript@latest

testscript testdata/script/*    # bütün skriptlər; exit 1 = fail (CI üçün)
testscript -v echo.txtar        # verbose həmişə
testscript -e VAR1=hello -e VAR2=goodbye script.txtar   # env ötür
testscript -work script.txtar   # workdir-i saxla (yolu çap edir)
```
**Shebang istifadəsi (Unix):**
```
#!/usr/bin/env testscript
exec echo hello
stdout 'hello'
```
chmod +x → `./hello.txtar` → PASS. (#! = interpreter göstərişi.)

### 19. Bug Report = txtar Repro
```
# I was promised 'Go 2', are we there yet?
exec go version
stdout 'go version go2.\d+'
```
Maintainer: `testscript repro.txtar` → FAIL (go1.18) → fix → PASS → skript
layihənin öz testlərinə əlavə olunur. Kent Beck: izahları test case şəklində
istə — testscript automatik testi GƏNİŞ yerlərə yumşaq girişidir.

### 20. Test İçində Test — ! exec go test
Müəllifin kitab nümunələri: bəzi nümunələr BY DESIGN fail olmalıdır (TDD —
test əvvəl yazılıb, implementation sonra). Skriptlə idarə:
```go
! exec go test                     # BU modul FAIL ETMƏLİDİR
stdout 'want 4, got 5'            # bu mesajla

exec go test                       # başqa modul isə PASS etməli
```
Kod skriptə köçürülmür — REPO-dan gəlir. **txtar-c aləti:**
```bash
go install github.com/bitfield/txtar-c@latest
txtar-c .                          # kataloqu txtar-a çevir
txtar-c -script test.txtar .       # + test skriptini başlıq edir
txtar-c -script test.txtar . | testscript   # boru ilə birbaşa işə sal
# PASS
```
Nəticə: minik kodla (build, exec, output compare) məşğul olmadan davranış
notasiyası — "By relieving the brain of all unnecessary work, a good notation
sets it free to concentrate on more advanced problems." (Whitehead)

## Əsas terminlələr
- testscript — go-internal-dən DSL test skript aləti
- .txtar — text archive: skript + daxili fayllar bir faylda
- exec — proqram işə sal + exit 0 assert
- stdout/stderr — regex match assertləri
- ! — neqasiya (fail GÖZLƏ)
- testscript.Run(t, Params{Dir}) — skriptləri paralel subtest kimi işə salır
- TestMain + RunMain — custom binary-lərin $PATH qeydiyyatı
- hello.Main() int — delegate main pattern (exit status qaytarır)
- os.Exit(RunMain(...)) — unutma — yoxsa həmişə uğur
- total coverage — subprocess daxil coverage
- cmp/cmpenv — golden fayl / env-expand muqayisəsi
- exists/grep/-count — fayl mövcudluğu / regex / dəqiq say
- stdin (fayl/stdout) — input simulyasiyası, pipeline
- cp/mv/mkdir/cd/rm/symlink — fayl əməliyyatları
- Fazalar (# şərhlər) — fail outputunu qısaldır
- [condition] — exec:/go1.x/darwin/!arm64/unix; skip
- env / $WORK / $TMPDIR / $HOME=/no-home / $exe
- Setup + env.Setenv — dinamik dəyər ötürülməsi
- & + wait — fon proqramları + gözləmə
- Standalone testscript — CI/automation, -v/-e/-work, shebang
- txtar-c — kataloq + skript birləşdirici
- Repro — bug report-da txtar test case

## Praktik nəticə
(1) CLI alətini test edirsənsə — əllə build/run/parse YOX: testscript.Run +
Dir: testdata/script. (2) exec = çağırış + exit-0 assert; ! = əks; ! stdout .
= boşluq assert. (3) Custom binary: main → Delegate Main() int, TestMain-də
RunMain mapinə qeyd et, MUTLƏQ os.Exit-ə sarı. (4) coverage-i total coverage
rəqəmindən oxu — subprocess-lər daxil; cover.out + go tool cover işləyir.
(5) Golden fayl: txtar daxilinde `-- golden.txt --`; cmp tam, grep hissəvi,
-count dəqiq say; env-asılı output üçün cmpenv. (6) İstifadəçi inputu: stdin
fayl + sətir-sətir "söhbət"; `stdin stdout` = pipeline ( konkurensi yoxdur).
(7) Double quote escape ETMİR — literal keçir; single quote '' ilə escape.
(8) Şərhlər fazaları ayırır — fail olanda yalnız cari fazanın logu görünür.
(9) Platforma asılılıqları: [darwin]/[unix]/[$exe] ilə cross-platform skript.
(10) Dinamik dəyər lazımdırsa (port!) — Setup + env.Setenv + $VAR. (11) Server
testi: exec listen & + curl --retry-connrefused + wait. (12) Bug report =
testscript repro.txtar — reproduce + fix yoxlaması + regression testi birdən.
(13) txtar-c . | testscript — mövcud kataloqları minik kodla yoxla; "! exec go
test" — BY DESIGN fail olan TDD nümunələrini idarə et.

## Mənbə
Pages: 259-308 (PDF 271-320)
