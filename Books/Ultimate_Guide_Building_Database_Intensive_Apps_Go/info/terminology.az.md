# database/sql (Ultimate Guide) — Terminoloji Lüğət (Azərbaycanca)

## B

**Bağlantı pool-u (connection pool)** — sql.DB daxilində idarə olunan
bağlantı dəsti; lazy yaradılır, 10x ölü-bağlantı retry.

## C

**Client-side prepared statement** — saxta prepare (string birləşdirmə);
Perl DBI-nin default-u; Go HƏQİQƏTƏN server tərəfində prepare edir.

**Connection leak (bağlantı sızması)** — pool-a qayıtmayan bağlantı;
ignor edilmiş rows.Close() unutması → "too many connections".

**Context** — ləğv/timeout/trace daşıyıcısı; 1.8+: QueryContext; köhnə
metodlar thin wrapper-dir.

**ConnMaxLifetime** — bağlantının maksimal ömrü; köhnəlmiş bağlantıların
avtomatik qapatılması.

## D

**driver.Valuer** — Value() (driver.Value, error) — DB-yə YAZILARKƏN
transformasiya interfeysi.

**DSN (Data Source Name)** — sürücü-specific bağlantı sətri; sql.Open-ın
2-ci parametri.

## E

**Exec vs Query** — Exec: nəticə seti YOX, bağlantı DƏRHAL qaytarır;
Query: rows.Close-a qədər TUTULUR.

## I

**Idle connections** — pool-da saxlanan boş bağlantılar; SetMaxIdleConns
default 0 = dərhal qapanır → thrashing.

## L

**Lazy connection** — ilk real əməliyyatda açılan bağlantı; Open yalnız
obyekt yaradır.

## N

**NextResultSet** — çoxlu nəticə dəstləri (stored procedure-lər); 1.8+.

**NullXXX (NullString/NullFloat64/...)** — NULL-scan üçün {T; Valid bool}
strukturları.

## P

**Prepared statement** — server-də əvvəlcədən təhlil edilmiş, parametrli
soruğ; sql.Stmt; parametrli db.Exec GİZLİ prepare edir.

**Pitfall (tələ)** — 13 siyahı: döngüdə defer, çoxlu db, ignor Query,
uint64 MSB, NULL sürprizi...

## R

**RawBytes** — kopyasız scan; DB-yə məxsus yaddaş, MƏHDUD ömür;
QueryRow.Scan ilə İŞLƏMİR.

**rows.Err()** — döngüdən SONRA mütləq yoxlanmalı xəta (io.EOF özü həll
olunur).

## S

**Side-effect import (`_`)** — yalnız init() (driver qeydiyyatı) üçün
import; ad sahəsi bağlı deyil.

**sql.Conn** — tək-bağlantı zəmanəti (1.9+); tranzaksiyasız; yalnız
Context-vari metodlar.

**sql.Scanner** — Scan(src interface{}) error — DB-dən OXUNARKƏN
transformasiya interfeysi.

**sql.Tx** — tranzaksiya; tək bağlantı; db ilə ƏLAQƏSİZ; daxilində retry
QEYRİ-AKTİV; paralel sorğu QADAĞA.

## T

**Thrashing** — idle qadağası səbəbindən bağlantıların davamlı açılıb-
qapanması (latency + yüklənmə).

## W

**WaitCount / WaitDuration** — db.Stats() sahələri: pool-un dolu
olmasından gözləmələrin sayı/cəmi — "sorğu yavaşdımı, pool gözlədimi?"
sualının cavabı.
