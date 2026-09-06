# Go Web Programming — Mündəricat (TOC)

**Müəllif:** Sau Sheong Chang
**Nəşriyyat:** Manning Publications Co., 2016 · **ISBN:** 9781617292569
**Səhifə:** 314 (PDF) · **Book pages:** 280 + indeks · **10 chapter, 3 hissə**

> PDF offset: book page N → PDF page N+21

## Hissələr

- **Part 1 — Go and web applications** (Ch 1-2): giriş + tam ChitChat forum tətbiqi
- **Part 2 — Basic web applications** (Ch 3-6): request emalı, formalar, şablonlar, data saxlama
- **Part 3 — Being real** (Ch 7-10): web services, testing, concurrency, deployment

## Chapter-lər

| Chapter | Başlıq | Book səh. | PDF səh. |
|---------|--------|-----------|----------|
| 1 | Go and web applications | 3-21 | 24-42 |
| 2 | Go ChitChat | 22-44 | 43-65 |
| 3 | Handling requests | 47-68 | 68-89 |
| 4 | Processing requests | 69-95 | 90-116 |
| 5 | Displaying content | 96-124 | 117-145 |
| 6 | Storing data | 125-152 | 146-173 |
| 7 | Go web services | 155-189 | 175-210 |
| 8 | Testing your application | 190-221 | 211-242 |
| 9 | Leveraging Go concurrency | 223-254 | 244-276 |
| 10 | Deploying Go | 256-280 | 277-301 |

## Section-lər

### Ch 1 — Go and web applications (pdf 24-42)
1.1 Using Go for web applications (4) · 1.2 How web applications work (6) · 1.3 A quick introduction to HTTP (8) · 1.4 The coming of web applications (8) · 1.5 HTTP request (9) · 1.6 HTTP response (12) · 1.7 URI (14) · 1.8 Introducing HTTP/2 (16) · 1.9 Parts of a web app (16) · 1.10 Hello Go (18) · 1.11 Summary (21)

### Ch 2 — Go ChitChat (pdf 43-65)
2.1 Let's ChitChat (23) · 2.2 Application design (24) · 2.3 Data model (26) · 2.4 Receiving and processing requests (27) · 2.5 Generating HTML responses with templates (32) · 2.6 Installing PostgreSQL (37) · 2.7 Interfacing with the database (38) · 2.8 Starting the server (43) · 2.9 Wrapping up (43) · 2.10 Summary (44)

### Ch 3 — Handling requests (pdf 68-89)
3.1 The Go net/http library (48) · 3.2 Serving Go (50) · 3.3 Handlers and handler functions (55) · 3.4 Using HTTP/2 (66) · 3.5 Summary (68)

### Ch 4 — Processing requests (pdf 90-116)
4.1 Requests and responses (69) · 4.2 HTML forms and Go (74) · 4.3 ResponseWriter (82) · 4.4 Cookies (87) · 4.5 Summary (95)

### Ch 5 — Displaying content (pdf 117-145)
5.1 Templates and template engines (97) · 5.2 The Go template engine (98) · 5.3 Actions (102) · 5.4 Arguments, variables, and pipelines (110) · 5.5 Functions (111) · 5.6 Context awareness (113) · 5.7 Nesting templates (119) · 5.8 Using the block action (123) · 5.9 Summary (124)

### Ch 6 — Storing data (pdf 146-173)
6.1 In-memory storage (126) · 6.2 File storage (128) · 6.3 Go and SQL (134) · 6.4 Go and SQL relationships (143) · 6.5 Go relational mappers (147) · 6.6 Summary (152)

### Ch 7 — Go web services (pdf 175-210)
7.1 Introducing web services (155) · 7.2 Introducing SOAP-based web services (157) · 7.3 Introducing REST-based web services (160) · 7.4 Parsing and creating XML with Go (163) · 7.5 Parsing and creating JSON with Go (174) · 7.6 Creating Go web services (181) · 7.7 Summary (188)

### Ch 8 — Testing your application (pdf 211-242)
8.1 Go and testing (191) · 8.2 Unit testing with Go (191) · 8.3 HTTP testing with Go (200) · 8.4 Test doubles and dependency injection (204) · 8.5 Third-party Go testing libraries (210) · 8.6 Summary (221)

### Ch 9 — Leveraging Go concurrency (pdf 244-276)
9.1 Concurrency isn't parallelism (223) · 9.2 Goroutines (225) · 9.3 Channels (232) · 9.4 Concurrency for web applications (240) · 9.5 Summary (254)

### Ch 10 — Deploying Go (pdf 277-301)
10.1 Deploying to servers (257) · 10.2 Deploying to Heroku (263) · 10.3 Deploying to Google App Engine (266) · 10.4 Deploying to Docker (271) · 10.5 Comparison of deployment methods (279) · 10.6 Summary (280)

## Layihə davamlılığı (project_continuity)

```json
{
  "project_continuity": true,
  "project_name": "chitchat + ws (web service)",
  "chapters_involved": [2, 7],
  "note": "Chapter 2 tam ChitChat forumu qurur. Chapter 7 ayrıca ws (web service) layihəsi qurur (data.go + server.go + posts.go). Chapter 10-da ws layihəsi deploy olunur. project-state/ qovluğu lazım deyil — kitabda fərqli addımlar ayrıca göstərilir, kumulyativ fayl dəsti yoxdur."
}
```
