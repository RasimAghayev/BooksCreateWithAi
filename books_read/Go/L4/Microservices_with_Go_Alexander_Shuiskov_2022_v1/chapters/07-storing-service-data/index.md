# Chapter 7 — Storing Service Data (səh. 133-146)

## Bu fəsil nədən bəhs edir?

In-memory repozitoriyaların çatışmazlıqları, database növləri, ACID,
replication və Movie servislərinin MySQL-ə köçürülməsi.

## Əsas fikirlər

### 1. Niyə in-memory kifayət etmir?
- **Davamlılıq yoxdur:** restart → bütün data itir
- **Çox instansalı inkonsistensiya:** movie servisi yazmağa instans-0 seçir,
  oxumağa instans-1 düşür → `ErrNotFound`! Hər instans müstəqil dataset saxlayır
- Həll: bütün instansların paylaşdığı tək DB

### 2. Database-lərin təklif etdiyi zəmanətlər
- **Durability (Davamlılıq):** yazılan data restartlardan sağ çıxır
- **ACID transaction-lar:**
  - **A**tomicity — ya hamısı, ya heç biri
  - **C**onsistency — bir valid vəziyyətdən digərinə
  - **I**solation — paralel dəyişikliklər ardıcıl kimi görünür
  - **D**urability — dəyişikliklər persist olunur
- **Replication:** replikalara köçürmə — resilience + oxu latensiyasının azalması
- Query dəstəyi (SQL, JOIN-lər)

### 3. Database növləri
| Növ | Model | Nə vaxt |
|---|---|---|
| Key-value | açar→dəyər, yalnız get/put | ən sürətli; sadə identifikatorlu data |
| Relational | cədvəllər (sətir/sütun), SQL, JOIN | universal; strukturlu, şemalı data |
| Document | JSON/XML sənədlər, ŞEMA YOX | fərqli strukturlu sənəd topluları |
| Graph | vertex + edge; traversal/connectivity | əlaqə data (sosial şəbəkə) |
| Blob | binary, immutable, append-only | fayl/medianın saxlanması |

- "One-size-fits-all" YOXDUR; amma relational model (1970, E.F. Codd) çoxlu
  məsələni qarşılayır → MySQL seçilir

### 4. MySQL praktikası
- Docker: `docker run --name movieexample_db -e MYSQL_ROOT_PASSWORD=... -p 3306:3306 -d mysql:latest`
- Şema (schema/schema.sql):
```sql
CREATE TABLE IF NOT EXISTS movies (id VARCHAR(255), title VARCHAR(255),
    description TEXT, director VARCHAR(255));
CREATE TABLE IF NOT EXISTS ratings (record_id VARCHAR(255),
    record_type VARCHAR(255), user_id VARCHAR(255), value INT);
```
- VARCHAR(255) vs TEXT (uzun mətn — description)
- Başlatma: `docker exec -i movieexample_db mysql ... < schema/schema.sql`

### 5. Go repo implementasiyası
```go
import _ "github.com/go-sql-driver/mysql"   // driver (yalnız side-effect)

db, _ := sql.Open("mysql", "root:password@/movieexample")  // connection string
```
- **Metadata Get:** `QueryRowContext` + `row.Scan`; `sql.ErrNoRows` →
  `repository.ErrNotFound` mapping
- **Put:** `ExecContext` — INSERT
- **Rating Get:** `QueryContext` → `rows.Next()` dövrü + `rows.Scan`;
  boş nəticə → ErrNotFound
- Driver importu `_` prefiksi ilə — database/sql üçün qeydiyyat
- Credential-ların kodda saxlanması PİS təcrübə — ayrı konfiq faylı (Ch8-də)

### 6. Test (grpcurl ilə)
```
grpcurl -plaintext -d '{"record_id":"1","record_type":"movie"}' \
  localhost:8082 RatingService/GetAggregatedRating   → NotFound
grpcurl ... PutRating (user alex, 5) → yazılır
grpcurl ... GetAggregatedRating → {"ratingValue": 5}
```
- Servisi RESTART edib yenidən oxu → data QALIR = persistence təsdiqləndi

## Termindirmə (AZ)
- Durability — Davamlılıq (data restartlardan sonra da mövcuddur)
- ACID — Atomluq/Konsistensiya/İzolyasiya/Davamlılıq
- Replication — Replikasiya (data nüsxələnməsi)
- Replica — Replika (nüsxə saxlayan qovşaq)
- Connection String — Qoşulma Sətri
- Blob (Binary Large Object) — İkili Böyük Obyekt

## Kviz sualları
1. İki metadata instansı ilə in-memory repo hansı bug verir? (Yazma instans-0-də,
   oxuma instans-1-də → data yoxdur — inkonsistensiya)
2. TEXT ilə VARCHAR(255) fərqi? (TEXT uzun mətnlər üçün, VARCHAR maks 255 simvol)
3. `_ "github.com/go-sql-driver/mysql"` importu nə üçün lazımdır? (Driveri
   database/sql-ə qeydiyyata alır — side-effect import)
4. ACID-nin I hərfi nə deyir? (Isolation — paralel dəyişikliklər ardıcıl
   icra kimi görünür)
