# Chapter 7 — Go web services Cheatsheet

## `xml.Unmarshal(data, &struct)`

**Nə edir:** XML mətni (və ya `[]byte`) struct-a çevirir. Struct tag-ları ilə xəritələşmə: `xml:"name"`, `xml:"name,attr"`, `xml:",chardata"`, `xml:",innerxml"`, `xml:"a>b>c"`.

```go
type Post struct {
    XMLName xml.Name `xml:"post"`
    Id      string   `xml:"id,attr"`
    Content string   `xml:"content"`
}
var post Post
xml.Unmarshal(xmlData, &post)
```

**Mənbə:** Chapter 7, page 165

---

## `xml.NewDecoder(r).Token()` / `DecodeElement(&v, &se)`

**Nə edir:** Streaming XML-i token-token oxuyur. `StartElement` tutulur, `DecodeElement` ilə struct-a parse olunur.

```go
decoder := xml.NewDecoder(xmlFile)
for {
    t, err := decoder.Token()
    if err == io.EOF { break }
    switch se := t.(type) {
    case xml.StartElement:
        if se.Name.Local == "comment" {
            var comment Comment
            decoder.DecodeElement(&comment, &se)
        }
    }
}
```

**Mənbə:** Chapter 7, page 170

---

## `xml.Marshal(&v)` / `xml.MarshalIndent(&v, prefix, indent)`

**Nə edir:** Struct-ı XML-ə çevirir. `MarshalIndent` oxunaqlı format verir. `xml.Header` XML deklarasiyası üçün.

```go
output, _ := xml.MarshalIndent(&post, "", "\t")
ioutil.WriteFile("post.xml", []byte(xml.Header+string(output)), 0644)
```

**Mənbə:** Chapter 7, page 172

---

## `xml.NewEncoder(w).Encode(&v)`

**Nə edir:** Struct-ı `io.Writer`-ə birbaşa XML kimi yazır. `encoder.Indent(prefix, indent)` ilə formatlama.

```go
xmlFile, _ := os.Create("post.xml")
encoder := xml.NewEncoder(xmlFile)
encoder.Indent("", "\t")
encoder.Encode(&post)
```

**Mənbə:** Chapter 7, page 174

---

## `json.Unmarshal(data, &struct)`

**Nə edir:** JSON mətni struct-a çevirir. Struct tag: `json:"key"`. Nested obyektlər üçün başqa struct, list-lər üçün slice.

```go
type Post struct {
    Id      int       `json:"id"`
    Content string    `json:"content"`
    Comments []Comment `json:"comments"`
}
var post Post
json.Unmarshal(jsonData, &post)
```

**Mənbə:** Chapter 7, page 177

---

## `json.NewDecoder(r).Decode(&v)`

**Nə edir:** `io.Reader`-dən JSON-i streaming oxuyur. `http.Request.Body` üçün ideal. EOF ilə dayanır.

```go
decoder := json.NewDecoder(jsonFile)
for {
    var post Post
    if err := decoder.Decode(&post); err == io.EOF { break }
}
```

**Mənbə:** Chapter 7, page 178

---

## `json.MarshalIndent(&v, prefix, indent)` / `json.NewEncoder(w).Encode(&v)`

**Nə edir:** Struct-ı JSON-a çevirir. `MarshalIndent` formatlanmış, `Encoder.Encode` isə writer-ə birbaşa yazır (avtomatik newline).

```go
output, _ := json.MarshalIndent(&post, "", "\t\t")
ioutil.WriteFile("post.json", output, 0644)

// və ya:
jsonFile, _ := os.Create("post.json")
json.NewEncoder(jsonFile).Encode(&post)
```

**Mənbə:** Chapter 7, page 179

---

## REST handler multiplexer: `switch r.Method`

**Nə edir:** Tək handler URL prefix-inə bütün HTTP metodlarını (GET/POST/PUT/DELETE) yönləndirir.

```go
http.HandleFunc("/post/", handleRequest)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    var err error
    switch r.Method {
    case "GET":    err = handleGet(w, r)
    case "POST":   err = handlePost(w, r)
    case "PUT":    err = handlePut(w, r)
    case "DELETE": err = handleDelete(w, r)
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

**Mənbə:** Chapter 7, page 183

---

## URL-dən resurs ID-si: `strconv.Atoi(path.Base(r.URL.Path))`

**Nə edir:** `/post/123` URL-in son seqmentini (123) `int`-ə çevirir.

```go
id, err := strconv.Atoi(path.Base(r.URL.Path))
```

**Mənbə:** Chapter 7, page 183

---

## Sorğu body-ni oxuma: `r.Body.Read(body)` + `json.Unmarshal`

**Nə edir:** HTTP sorğu body-sindəki JSON-i Post struct-una parse edir.

```go
len := r.ContentLength
body := make([]byte, len)
r.Body.Read(body)
var post Post
json.Unmarshal(body, &post)
```

**Mənbə:** Chapter 7, page 185

---

## JSON cavab yazma: `w.Header().Set("Content-Type", "application/json")` + `w.Write(output)`

**Nə edir:** JSON cavabı client-ə göndərir. Content-Type header mütləqdir.

```go
output, _ := json.MarshalIndent(&post, "", "\t\t")
w.Header().Set("Content-Type", "application/json")
w.Write(output)
```

**Mənbə:** Chapter 7, page 186
