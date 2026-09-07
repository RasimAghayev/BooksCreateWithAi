# Chapter 9 — Metaprogramming with Go (səh. 317-357)

## Bu chapter nədən bəhs edir?

Proqramçı öz alətlərini düzəldən sənətkardır: API inteqrasiyaları ilə
avtomatlaşdırma, proqramları koordinə edən Go alətləri (`exec`), və kod
generasiyası (AST + `text/template`). Hər birinin meyarı: alətə sərf olunan vaxt
< qənaət olunan vaxt.

## Əsas fikirlər

### 1. API inteqrasiyaları (PaaS/SaaS avtomatlaşdırması)
**Nədir:** Gündəlik istifadə etdiyimiz servislərin (GitLab, PagerDuty, Google
Calendar, AWS...) API-ləri ilə təkrarlanan işlərin avtomatlaşdırılması.

**Nümunələr:**
- Merge request yaradanda 30 saniyəlik reviewer əlavə etmə → kiçik app API ilə
  əlavə edir (team member-lər hardcode belə ola bilər)
- On-call cədvəli + tətil tarixi (Calendar API) → cədvəl özü yenilənir
- PagerDuty istifadəçi contact-method yoxlaması (kitabın əsas nümunəsi)

**PagerDuty nümunəsinin quruluşu:**
```go
func (u *UsersAPI) buildRequest(ctx context.Context) (*http.Request, error) {
    params := &url.Values{}
    params.Set("include[]", "contact_methods") // contact method-lar istənilir
    uri := u.apiBaseURL + "/users?" + params.Encode()
    req, err := http.NewRequestWithContext(ctx, "GET", uri, http.NoBody)
    req.Header.Set("Authorization", "Token token="+u.apiToken)
    return req, err
}
```
**Vacib transformasiya qaydası:** API response struct-ı (`apiResponse`) birbaşa
paketi sızdırməq əvəzinə, **öz domain modelinə** çevir:
```go
type User struct {
    Name, Email string
    EmailIsSet, PhoneIsSet, SMSIsSet, PushIsSet bool // Current settings
}
```
Encapsulation + davranış zamanı API formatı dəyişsə, yalnız konversiya dəyişir.

**CLI konfiqurasiya (flag paketi):**
```go
flag.BoolVar(&onlyErrors, "errors", true, "print only those users with invalid settings")
flag.BoolVar(&requireSMS, "sms", false, "require SMS setting")
flag.Parse()
```
**Çıxış:** `fmt.Fprintf(w io.Writer, ...)` — yazmaq haraya main()-də qərar
verilir (test üçün `bytes.Buffer`).

**Dəyər meyarı:** 5 işçi üçün alət dəyməz; 100+ işçi üçün real fayda. Aləti
genişləndirmək (auto-email göndərmək) mümkündür — amma dəyər/qiymət soruşulmalı.

### 2. Go ilə proqramların koordinasiyası (exec)
**Nədir:** Shell skriptlərinin yerinə Go — yeni bacarıq öyrənmədən, tipli və
test oluna bilən koordinasiya.

**Nümunə: dəyişən paketlərin testi** (yalnız diff olunan paketlər):
```go
// 1. Dəyişən faylların siyahısı
cmd := exec.CommandContext(ctx, "git", "diff", "--name-only", "-M100%", "master")
cmd.Dir = dir
cmd.Env = os.Environ()
output, err := cmd.CombinedOutput()

// 2. Fayl yollarını paketlərə çevir + dedupe (map[string]struct{})

// 3. Repo kökü (yollar ona nisbidir)
cmd = exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")

// 4. Hər dəyişən paketdə test
exec.CommandContext(ctx, "go", "test", "-race", pkg)
```
**Sub-kod izahı:**
- `-M100%` → rename-ləri filterlə (100% oxşar = dəyişiklik deyil)
- `map[string]struct{}` → dedupe üçün empty struct set-i (Chapter 8 pattern-i)
- Ardıcıllıq: format → test → lint bir Go binary-də → mental yük azalır

### 3. Code Generation — kod generasiyası
**Nədir:** AST-ni oxuyub, `text/template` ilə kod yaratmaq (Mockery-nin öz
versiyasını yazırıq: interfeysdən stub generator).

**Addım-addım:**

**(a) Kodu parse et:**
```go
fs := token.NewFileSet()
parsedFile, err := parser.ParseFile(fs, filename, nil, 0)
// → *ast.File — kodun ağacı
```

**(b) Declaration-ları süz:**
```go
for _, decl := range parsedFile.Decls {
    generalDecl, ok := decl.(*ast.GenDecl) // import/const/type/var
    if !ok {
        continue // funksiya declaration-ları skip
    }
    // TypeSpec → interface → *ast.InterfaceType → Methods
}
```

**(c) Metodları çıxar:**
```go
funcType, ok := method.Type.(*ast.FuncType) // imza
outputMethod.Inputs = parseFieldList(funcType.Params)
outputMethod.Outputs = parseFieldList(funcType.Results)
```
**Fənd:** `func add(x, y int)` — bir tip iki ada aiddir → `param.Names`
 üzərində ikinci loop LABİLDİR (tip bir dəfə, adlar çox).

**(d) Template:**
```go
tmpl := template.New("generator")
tmpl.Funcs(template.FuncMap{ // custom funksiyalar
    "isNotLast": func(len, index int, s string) string {
        if index != len-1 { return s } // vergül üçün
        return ""
    },
    "stubValue": stubValue,
})
```
```go
var stubTemplate = `
package {{ .Interface.PackageName }}
type Stub{{ .Interface.Name }} struct {}
{{ range .Interface.Methods -}}
func (s *Stub{{ $.Interface.Name }}) {{ .Name }}(
{{- range $i, $p := .Inputs }}{{ $p.Name }} {{ $p.Type }}{{ isNotLast (len .Inputs) $i ", " }}{{ end -}}
) ({{/* çıxışlar */}}) {
    return {{ stubValue .Type }}
}
{{ end -}}
`
```

**(e) Stub dəyərlər:**
```go
func stubValue(typeName string) string {
    switch {
    case strings.HasPrefix(typeName, "int"):  return "0"
    case typeName == "bool":                  return "false"
    case typeName == "string":                return `""`
    case typeName == "error":                 return "nil"
    }
    return "nil"
}
```

**(f) Generator testi (TDD ilə inkişaf):** gözlənilən nəticəni string kimi testə
qoy → template-i iterativ doldur.

**Generasiya qaydaları:**
- Generasiya olunmuş kod **repo-ya commit edilir** (build repeatable olsun)
- Faylın başında `// @generated` şərhi (və ya `z_[name].go` adı) — insanlar
  redaktə etməsin, gofmt-ə daxil edilsin
- `go generate` ilə bağlı; dəyişiklikdən sonra yenidən generasiya

### 4. Nə vaxt metaproqramlaşdırma?
Sıralama (müəllifin tövsiyə etdiyi try order): sadə kod → design pattern-lər →
generics → kod generasiyası. Hər dəfə problemi həll edən ən sadə vasitə seçilir.

**ROI hesabı:** avtomatlaşdırma 30 saniyə qənaət edirsə, gündə 10 dəfə = ilə
~21 saat; alətə 2 saat sərf etmək sərfəlidir. Amma vaqt itirməmək üçün alətin
özünü də məhdud saxla — "tool building intoxicatingdır" (alət düzəltmək
vərdirədir).

## Əsas terminlər

- Metaprogramming (metaproqramlaşdırma)
- AST (Abstract Syntax Tree / abstrakt sintaksis ağacı)
- Code Generation (kod generasiyası)
- `go:generate` direktivi
- Domain Model (domain modelinə transformasiya)
- ROI (Return on Investment / investisiya gəliri)

## Praktik nəticə

- Təkrarlanan SaaS əməliyyatlarını kiçik Go tool-lara çevir (hardcode ilə başla)
- API response-u dərhal istifadə etmə — öz modelinə çevir
- Shell skriptləri əvəzinə `exec.CommandContext` + dedupe set-lər
- Kod generasiyası: parse → süz → template + custom Funcs → commit + @generated
- Ən sadə həlldən başla: sadə kod → pattern → generics → generasiya

## Mənbə

Pages: 317-357 (Chapter 9, Beyond Effective Go Part 2)
