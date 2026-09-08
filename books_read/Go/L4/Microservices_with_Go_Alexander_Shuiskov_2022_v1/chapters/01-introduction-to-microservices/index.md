# Chapter 1 — Introduction to Microservices (səh. 1-12)

## Bu fəsil nədən bəhs edir?

Mikroservis arxitekturasının tərifi, monolit yanaşmanın çatışmazlıqları,
mikroservislərin faydaları/riskləri, nə vaxt istifadə edilməli və Go-nun
mikroservis inkişafındakı rolu.

## Əsas fikirlər

### 1. Mikroservis nədir?
- **Mikroservis arxitekturası (Microservice architecture):** aplikasiyanın
  hər biri müəyyən **biznes qabiliyyətinə (business capability)** cavab
  verən müstəqil servislər toplusu kimi təşkili
- Nümunə: onlayn marketplace — search, cart, payments, order history hər biri
  ayrı servis ola bilər; search ilə payments-in texniki ortaq cəhəti yoxdur

### 2. Monolitin çatışmazlıqları (motivasiya)
| Problem | Təsviri |
|---|---|
| Böyük ölçü, yavaş deploy | build/start/deploy dəqiqələr çəkə bilər |
| Hissəni müstəqil deploy edilməsi | mümkün deyil — bottleneck |
| Blast radius (təsir dairəsi) | bir bug bütün sistemə təsir edir |
| Vertical scaling həddi | CPU/RAM limitinə çatır |
| Interference (qarşılıqlı maneə) | bir hissə CPU/IO/RAM-ı yükləyir, hamısı əziyyət çəkir |
| Artıq asılılıqlar | refaktor bir anda payments-i sındıra bilər |
| Təhlükəsizlik | bir loophole = bütün komponentlərə giriş |

- Monolit nə vaxt YAXŞIDIR: kiçik kod bazası, early-stage tez iterasiya,
  dar scope (məs. parol generator servisi)

### 3. Mikroservislərin faydaları
- Sürətli build/deploy, kiçik deployable ölçüsü
- Xüsusi deploy cadence + monitoring + müstəqil avtomatik testlər
- Cross-language dəstək, texnoloji azadlıq
- Sadə (fine-grained) API-lər
- Horizontal scaling ( Horizontal Scaling (Üfüqi Miqyaslandırma)) — hər servis
  ayrıca miqyaslanır, ucuz instanslar
- Fault isolation (Qəza İzolyasiyası) — qismən sıradan çıxma sistemi öldürmür
- Başa düşülə bilmə, asan refaktor, paylanmış komanda işi, xərj optimizasiyası

### 4. Mikroservislərin riskləri
- **Resurs overhead:** şəbəkə trafiki, latency, IO; hər komponent ayrı proses
- **Debug çətinliyi:** bir request N servisin loglarını yoxlamaq tələb edir
- **Integration testing** dəsti böyüyür
- **Konsistensiya və transaction-lar:** data sistemə səpələnir, atomik dəyişiklik çətin
- Divergence (Fərqlənmə): fərqli kitabxana versiyaları, upgrade çətinliyi
- Tech debt idarəsi çətinləşir; Observability (Müşahidə Qabiliyyəti) — log/trace/metrics hamı üçün
- Funksiya dublikasiyası; Ownership (Sahiblik) müqavilələri vacibdir

### 5. Nə vaxt və necə istifadə etməli
- **Çox tez mikroservisə keçmə** — məhsul hələ dəyişkəndə monolitlə başla,
  sərhədlər aydınlaşanda böl
- "No size fits all" — komanda ölçüsü, coğrafiya, paylanma nəzərə alın
- Best practices: **design for failure** (uğursuzluq üçün dizayn),
  **embrace automation** (avtomatlaşdırmayı qucaqla),
  **don't ship hierarchy** (təşkilati struktur üzrə deyil, domain üzrə böl —
  Conway qanununa qarşı çıx), integration testlərə investisiya,
  backward compatibility (Geri Uyumluluq) + versioning

### 6. Go-nun rolu
- Sadəlik, şəbəkə aplikasiyaları yazmağın asanlığı, native binary-yə sürətli compile
- Standart kitabxana production-keyfiyyətli (net/http, encoding/json...)
- Explicit error handling (Açıq Xəta İdarəsi) — "design for failure" prinsipinə uyğun
- Concurrency (Eyni Zamanlılık): goroutine + channel + sync paketi — bir servis
  bir neçə başqasını paralel çağırıb nəticələri birləşdirir
- Böyük community, zəngin tooling ekosistemi

## Termindirmə (AZ)
- Microservice Architecture — Mikroservis Arxitekturası
- Monolith — Monolit (tək proqram kimi təşkil olunan aplikasiya)
- Business Capability — Biznes Qabiliyyəti
- Blast Radius — Təsir Dairəsi (xətanın yayılma sahəsi)
- Fault Isolation — Qəza İzolyasiyası
- Design for Failure — Uğursuzluq üçün Dizayn
- Observability — Müşahidə Qabiliyyəti (log/metrics/trace)

## Kviz sualları
1. Monolit "blast radius" probleminin səbəbi nədir? (Geniş istifadə olunan
   funksiya/kitabxanadakı bug bütün hissələrə təsir edir)
2. "Don't ship hierarchy" nə deməkdir? (Servisləri komanda strukturuna görə
   deyil, biznes domain-lərinə görə böl)
3. Go-nun hansı 2 xüsusiyyəti mikroservis üçün idealdır? (Explicit error
   handling; asan concurrency — goroutine/channel)
