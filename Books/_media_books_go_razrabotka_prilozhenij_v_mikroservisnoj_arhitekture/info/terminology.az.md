# Go mikroservisləri (Popova 2026) — Terminoloji Lüğət (Azərbaycanca)

## A

**ACID** — tranzaksiya tələbləri: Atomicity (hamısı və ya heç biri),
Consistency (konsistent haldan hala), Isolation (paralellərdən qoruma),
Durability (commit sonrası dəyişməzlik).

**AMQP** — Advanced Message Queuing Protocol; RabbitMQ-nun əsası.

**API Gateway** — xarici dünya üçün tək giriş nöqtəsi; avtorizasiya, rate
limit, keş, aqreqasiya.

**At Most/Least/Exactly Once** — mesaj çatdırılma zəmanət səviyyələri.

## B

**BalanceOperationType** — Deposit/Credit (balans artır/azalt) əməliyyat tipi.

**B-tree indeks** — balanslaşmış ağac; PostgreSQL default indeks növü.

**bcrypt** — parol hash funksiyası; salt + cost factor (DefaultCost=10,
bilərəkdən yavaş → bruteforse qarşı).

**Bearer** — `Authorization: Bearer [token]` sxemi.

## C

**Claims** — JWT payload sahələri (iss, aud, exp, sub).

**Compact topic** — Kafka-da hər açarın yalnız SON dəyərinin saxlandığı topic.

**Cyclic import** — A→B→A; protoc üçün FATAL — modulluq pozulub.

## D

**DEBIT/CREDIT** — buxgalteriya istiqamətləri: balans artır / azaldır.

**Dependency Injection (DI)** — asılılıqlar xaricdən; Go-da əl ilə
(internal/app kompozitoru).

**Determinant** — funksional asılılığın sol tərəfi; НФБК-də potensial açar
olmalıdır.

**Dirty read** — commit OLMAMIŞ datanın oxunması.

**Double-entry** — hər əməliyyatın ≥2 qeyd ilə (göndərən CREDIT + alıcı
DEBIT).

**ДКНФ** — hər məhdudiyyət domen+açar məhdudiyyətlərinin məntiqi nəticəsi.

## E

**Exchange** — AMQP-da mesajların daxil olduğu yönləndirici; fanout/direct/topic.

**Exactly Once** — dəqiq bir dəfə çatdırılma; idempotentlik + unikal ID tələb edir.

## F

**Fanout** — bütün bağlanmış queue-lara yayımlama.

**Fingerprint (kitab kontekstində YOX — bu layihədə)** — icaze verilmir, bu
lüğət kitaba aiddir.

**Фантомное чтение (Phantom read)** — sorğunun YENİ sətirlər görməsi (say dəyişir).

## G

**goose** — Go miqrasiya aləti; up/down cütləri, tranzaksiya daxilində icra.

**GORM** — Go ORM; WithContext, clause.OnConflict (upsert), Updates.

**gRPC** — HTTP/2 + Protobuf RPC protokolu; unary/stream/bidirectional.

## H

**HS256** — HMAC-SHA256 JWT imza üsulu; interceptor-da alg YOXLANMALI.

## I

**Interceptor** — gRPC middleware; unary/stream; JWT yoxlama + context-ə userID.

**Isolation levels** — Read Uncommitted → Read Committed → Repeatable Read →
Serializable; PostgreSQL: RU-dan dirty read MÜMKÜN DEYİL.

## J

**JWT** — JSON Web Token; header.payload.signature; RegisteredClaims.

## K

**kafka-go** — segmentio-nun Go Kafka kitabxanası; Writer/Reader.

**KRaft** — Kafka-nın ZooKeeper-sız idarə rejimi (broker+controller rolları).

**kubectl** — K8s CLI; apply/get/describe/logs.

## L

**Lost update** — son yazan digərinin dəyişikliyini İTİRİR.

## M

**Manifest** — K8s deklarativ YAML (Deployment və s.).

**Mapper** — qatlararası model çevirici (PbToUser, UserToRepoUser).

**Multi-stage build** — build və runtime ayrı image-lər; alpine ~5MB runtime.

## N

**НФБК (Boyce-Codd)** — determinant = potensial açar; tək açarli 3НФ = НФБК.

**Non-repeatable read** — eyni tranzaksiyada təkrar oxu FƏRLİ nəticə verir.

## O

**Offset** — Kafka partition daxili mövqe; təkrar oxunmanın qarşısı.

**OSI model** — 7 qatlı şəbəkə modeli (fiziki→tətbiqi).

## P

**Partition** — Kafka topic bölgüsü; daxilində SIRALI mesajlar.

**Phantom read** — phantom.

**Pod** — K8s-də 1+ konteyner; klaster-daxili unikal IP.

**publicMethods** — interceptor-da autentifikasiyasız metodlar siyahısı
(Register/Login/Refresh).

## Q

**Qəpik int64** — pul məbləği *100 kimi tam ədəd; float YOX.

## R

**Race condition** — paralel yazışlar; -race detector, izolyasiya.

**Refresh token** — uzunömürlü DB-da saxlanılan SHA-256 token; revoked_at.

**Repository pattern** — biznes məntiqini DB detallarından ayıran qat.

**Rollback (miqrasiya)** — TƏRS sıra: indekslər → PK → cədvəllər.

## S

**Saga** — event əsaslı distributed tranzaksiya həlli; kompensasiya əməliyyatları.

**Salt** — hash-ə qarışan təsadüfi baytlar; bcrypt avtomatik.

**Serialization anomaly** — heç bir ardıcıl ekvivalent yoxdur.

**Stencilling** — (Know Go-dan) generics compile-time instansiasiya.

**Supply Chain Attack** — paketi dəyişib versiyanı saxlama; go.sum hash müdafiəsi.

## T

**Table-driven test** — ssenari cədvəli ilə test patterni.

**2PC (two-phase commit)** — koordinator + prepare/commit fazaları; yavaş,
bloklayıcı.

**Topic (Kafka)** — mesaj kanalı; adi (müddət limitli) və compact.

**TransactionEntry** — hesab hərəkəti; transaction_id + account_id +
direction + amount.

## U

**Upsert** — clause.OnConflict{UpdateAll: true} — mövcuddursa yenilə.

**U_ operatorları** — 6НФ temporal DB: U_JOIN və s.

## V

**Volume (K8s)** — pod daxili paylaşılan saxlanc.

**Volume (Docker)** — host-da data saxlancı (postgres_data və s.).

## W

**workflow_dispatch** — GitHub Actions manual trigger.

**Writer/Reader** — kafka-go producer/consumer obyektləri.

## Y

**YAGNI** — "You ain't gonna need it"; lazım olmayan funksiya YAZMA
(Transaction-da DELETE yoxdur).
