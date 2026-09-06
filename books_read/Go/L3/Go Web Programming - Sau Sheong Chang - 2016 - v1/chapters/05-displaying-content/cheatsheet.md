# Chapter 5 — Displaying content Cheatsheet

## `template.ParseFiles("tmpl.html")`

**Nə edir:** Template faylını parse edib `*Template` qaytarır. Sonradan `Execute(w, data)` ilə render olunur.

**Parametrlər:**
- fayl adı — parse ediləcək template faylı (ixtiyari sayda, variadic)

```go
t, _ := template.ParseFiles("tmpl.html")
t.Execute(w, "Hello World!")
```

**Mənbə:** Chapter 5, page 99

---

## `template.Must(template.ParseFiles(...))`

**Nə edir:** `ParseFiles`-ə sarğıdır; parse xətası olarsa `panic` edir (development üçün faydalı).

```go
t := template.Must(template.ParseFiles("tmpl.html"))
```

**Mənbə:** Chapter 5, page 101

---

## `t.Execute(w, data)` / `t.ExecuteTemplate(w, "name", data)`

**Nə edir:** Template-i `ResponseWriter`-ə render edir. `Execute` birinci template-i, `ExecuteTemplate` isə adla çağırır.

```go
t, _ := template.ParseFiles("t1.html", "t2.html")
t.ExecuteTemplate(w, "t2.html", "Hello World!")
```

**Mənbə:** Chapter 5, page 102

---

## `{{ if arg }} ... {{ else }} ... {{ end }}`

**Nə edir:** Şərti render — `arg` truthy olarsa birinci blok, əks halda `else` bloku göstərilir.

```html
{{ if . }}
  Number is greater than 5!
{{ else }}
  Number is 5 or less!
{{ end }}
```

**Mənbə:** Chapter 5, page 103

---

## `{{ range . }} ... {{ else }} ... {{ end }}`

**Nə edir:** Slice, array, map və ya channel iterasiyası. `.` iteration zamanı elementi göstərir. Boş olduqda `else` bloku göstərilir.

```html
<ul>
{{ range . }}
  <li>{{ . }}</li>
{{ else }}
  <li>Nothing to show</li>
{{ end }}
</ul>
```

**Mənbə:** Chapter 5, page 104

---

## `{{ with arg }} ... {{ else }} ... {{ end }}`

**Nə edir:** `arg` truthy olduqda scope açır və `.`-u `arg`-a dəyişir. `end`-dən sonra köhnə `.`-a qayıdır.

```html
<div>The dot is {{ . }}</div>
{{ with "world" }}
  Now the dot is {{ . }}
{{ end }}
```

**Mənbə:** Chapter 5, page 105

---

## `{{ template "name" arg }}`

**Nə edir:** Başqa template-i daxil edir. `arg` included template-ə data ötürür (boş buraxılarsa, `.` boş olur).

```html
{{ template "t2.html" . }}
```

**Mənbə:** Chapter 5, page 109

---

## `{{ define "name" }} ... {{ end }}`

**Nə edir:** Template faylında konkret template adı ilə blok təyin edir. Bir faylda bir neçə template müəyyən etmək olur.

```html
{{ define "layout" }}
<html><body>{{ template "content" }}</body></html>
{{ end }}
{{ define "content" }}<h1>Hello</h1>{{ end }}
```

**Mənbə:** Chapter 5, page 120

---

## `{{ block "name" arg }} ... {{ end }}`

**Nə edir:** `define` + `template` birləşməsi. Default content təyin edir; başqa faylda eyni adlı template parse olunarsa override olunur (Go 1.6+).

```html
{{ define "layout" }}
<body>{{ block "content" . }}<h1>Default</h1>{{ end }}</body>
{{ end }}
```

**Mənbə:** Chapter 5, page 123

---

## `$variable := value` (template dəyişəni)

**Nə edir:** `$` ilə başlayan dəyişən. `range` zamanı açar/dəyər saxlamaq üçün istifadə olunur.

```html
{{ range $key, $value := . }}
  Key: {{ $key }}, Value: {{ $value }}
{{ end }}
```

**Mənbə:** Chapter 5, page 110

---

## `{{ p1 | p2 | p3 }}` (Pipeline)

**Nə edir:** Arqumentləri/funksiyaları `|` ilə zəncir edir; `p1` nəticəsi `p2`-yə, sonra `p3`-ə ötürülür.

```html
{{ 12.3456 | printf "%.2f" }}
```

**Mənbə:** Chapter 5, page 111

---

## `template.FuncMap{"name": fn}` + `t.Funcs(funcMap)`

**Nə edir:** Template-ə xüsusi funksiya əlavə edir. `Funcs` **parse-dən əvvəl** çağırılmalıdır.

```go
funcMap := template.FuncMap{"fdate": formatDate}
t := template.New("tmpl.html").Funcs(funcMap)
t, _ = t.ParseFiles("tmpl.html")
```

**Mənbə:** Chapter 5, page 111

---

## `template.HTML(s)` (escape-i söndürmə)

**Nə edir:** `html/template`-ın avtomatik escape-ini söndürür; daxil olan HTML kodu olduğu kimi render olunur. Yalnız etibarlı məzmun üçün!

```go
t.Execute(w, template.HTML(r.FormValue("comment")))
```

**Mənbə:** Chapter 5, page 119
