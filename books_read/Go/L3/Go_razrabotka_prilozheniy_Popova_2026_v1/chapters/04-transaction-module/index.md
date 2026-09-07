# Глава 4 — Transaction modulunun hazırlanması (səh. 157-266)

## Bu fəsil nədən bəhs edir?

Normal formaların davamı (НФБК — Boyce-Codd, 4НФ — çoxdəyərli asılılıqlar,
5НФ — asılı birləşmələr, ДКНФ — domen-aşkar məhdudiyyətlər, 6НФ — temporal
U_ operatorları), miqrasiyalar (sütun dəyişmə alqoritmi, enum dəyişmə
strategiyası), indekslər (kitab predmet-göstəricisi metaforası, B-tree default,
sintez indekslər, rollback sırası: indekslər → PK → cədvəllər), tranzaksiyalar
+ ACID, paralel tranzaksiya anomalayaları (lost update, dirty/non-repeatable/
phantom read, serialization anomaly) və izolyasiya səviyyələri (Read
Uncommitted → Serializable; PostgreSQL fərqləri), Transaction servisi:
buxgalteriya cüt yazılışı (Transaction + Entry DEBIT/CREDIT), docker-compose
3 DB, gRPC müqavilələri (Deposit/Withdraw/Transfer/GetTransactions),
distribusiya problemi və həlləri: 2PC (koordinator, 2 faza) vs Saga ( Kafka
implementasiyası — kafka-go Producer/Consumer wrapper-lər, transaction_data /
transaction_response topic-lər, pending → completed/failed status axını,
HandleAccountResponse/HandleTransaction).

## Əsas fikirlər

### 1. Normal Formaların Davamı (4-6НФ)
- **НФБК (Boyce-Codd):** 3НФ-in xüsusi formu — hər qeyri-trivial
  funksional asılılığın DETERMİNANTI potensial açar olmalıdır. Nümunə:
  dayanacaq rezervasyonu {Dayanacaq, Başlanğıc, Bitmə, Tarif} — Tarif →
  Dayanacaq asılılığı (determinant Tarif açar DEYİL) → НФБК POZULUR →
  Tariflər və Rezervasiya cədvəllərinə parçalanır. Bir potensial açar
  olduqda НФБК = 3НФ.
- **4НФ:** НФБК + çoxdəyərli asılılıqlar (multi-valued) funksional olmalıdır.
  {Restoran, Pizza növü, Rayon} — Restoran → Pizzalar, Restoran → Rayonlar
  (bir-birindən asılı deyil) → yeni pizza əlavə edəndə HƏR rayona yeni yazı →
  parçala: Restoran-Pizza və Restoran-Rayon. Amma qiymət funksional
  asılılığı əlavə edilərsə — parçalama itki verir.
- **5НФ:** dövrəvi asılı birləşmə (A→B→C→A) YOX; təchizatçı-ticarət-alıcı
  nümunəsi — 3 ayrı cədvəl; birləşdirmə 3 cədvəlin JOIN-u ilə (2-lük JOIN
  yanlış nəticə verir). Reallıqda 4+ atributlarla praktiki rast gəlinmir.
- **ДКНФ:** bütün məhdudiyyətlər domen+açar məhdudiyyətlərinin məntiqi
  nəticəsidir → 5НФ-i avtomatik verir; gərgin tələb, həmişə mümkün deyil.
- **6НФ:** yalnız trivia join asılılıqları — parçalana bilməz; temporal DB-lər
  üçün (U_JOIN, U_ operatorları) — işçi/ünvan/vəzifə tarixçəsi ayrı saxlanılır.

### 2. Miqrasiya Mexanikası
- **Niyə:** insan amili (kirill "c" ilə "city" — proqramı sındıran səhv);
  köçürmə/fəlakət bərpasında avtomatik təkrar; start-da avtomatik icra.
- **Sütun tipi dəyişməsi (status):** status_old-a rename → yeni status yarat
  → data çevir → status_old sil. Milyon sətirlərdə vacib faza.
- **Dəyər dəyişməsi:** manual miqrasiya (sətir-sətir keçid); **enum halı:**
  1) enum-a yeni dəyər əlavə et → 2) dəyərləri köçür → 3) köhnə enum dəyərini
  sil — 3 AYRI miqrasiya.
- **Rollback sırası (vacib):** indekslər → PK-lar → cədvəllər (asılılıqları
  qırmaq üçün).

### 3. İndekslər
- **Niyə:** 10 000 istifadəçi, axtarılan 9999-cu — tam cədvəl skan;
  çoxsaylı belə sorğu = ölü yükdür. Predmet göstəricisi metaforası: termin →
  səhifə (sətir → link); göstəricinin ÖZÜ bütün kitabı əhatə etsə — axtarış
  yavaşlar → yalnız sıx sorğu sahələri indekslə.
- **Nümunə:** `CREATE INDEX ix_account_login ON users (login);` — auth-da
  login axtarışı üçün (çəkilmiş: `DROP INDEX ... ON users`).
- **Tiplər:** B-tree (default, çox hallarda optimal), hash, GiST, SP-GiST,
  GIN, BRIN, bloom.
- **Trade-off:** yazı (insert/update/delete) üçün əlavə yükdür — avtomatik
  yenidən qurulur.
- **Sintez (composite):** login + phone eyni sorğuda → birgə indeks.

### 4. Tranzaksiya və ACID
- **Tərif:** məntiqi vahid — hamısı və ya heç biri.
- **ACID:**
  - **Atomicity** — xəta olarsa hamısı rollback (əlaqə kəsildi, DB çökdü,
    tip xətası)
  - **Consistency** — bir konsistent haldan digərinə; qismən icra YOX
  - **Isolation** — paralel tranzaksiyalar nəticəyə təsir etməməli
  - **Durability** — commit təsdiqi gəldikdən sonra heç nə dəyişdirə bilməz
- **Pul köçürməsi** — kanonik nümunə: mənfi balans və ya itməmiş vəsait
  QAĐAĐANDIR.

### 5. Paralel Tranzaksiya Anomaliyaları
| Anomaliya | Ssenari (Sаşa/Pyotr ortaq hesab) | Nə baş verir |
|---|---|---|
| **Lost update** | hər ikisi 5000-dən əlavə edir | son yazan digərinin əlavəsini İTİRİR (5500 ≠ 7500) |
| **Dirty read** | ödənişdən sonra balans yoxlanılır | ödənilməmiş aralıq balans görünür (3000, real 500) |
| **Non-repeatable read** | eyni sorğu 2 dəfə | dəyər DƏYİŞƏNDƏ fərqli cavab |
| **Phantom read** | aylıq hesabat yazılır | sətir sayı artır/azalır (əlavə/sil) — non-repeatable-dan fərq: dəyəşmə YOX, SAY dəyişir |
| **Serialization anomaly** | paralel qrup | heç bir ardıcıl varianta uyğun gəlmir |

### 6. İzolyasiya Səviyyələri (SQL)
| Səviyyə | Nə həll edir | Qeyd |
|---|---|---|
| Read Uncommitted | — (ən aşağı) | dirty read buraxılır |
| Read Committed | dirty read | yalnız fiksədilmiş oxu; digəri gözləyir |
| Repeatable Read | lost update + non-repeatable | oxunmuş datanı başqası dəyişə bilməz |
| Serializable | phantom + hamısı | tam ardıcıl; şərtə uyğun sətirlər BLOKLANIR |

- **PostgreSQL spesifikası:** Read Uncommitted-də belə dirty read YOXDUR
  (= Read Committed kimi); Repeatable Read phantom-a İCAZƏ VERMİR.

### 7. Transaction Servis Dizaynı
- **Buxgalteriya cüt yazılışı (double-entry):**
  - `Transaction` — əməliyyat faktı: id, user_id, amount (QƏPİK-lərdə!),
    type (TRANSFER/DEPOSIT/WITHDRAW), status (INITIATED→pending→completed/
    failed/cancelled)
  - `TransactionEntry` — hesab üzrə hərəkət: transaction_id, account_id,
    direction (DEBIT=artır / CREDIT=azalt), amount DECIMAL(20,2)
  - Transfer = 2 entry: göndərən CREDIT + alıcı DEBIT — bir əməliyyatda
- **DB sxemi:** CHECK məhdudiyyətləri (type/status/direction IN ...);
  FK ON DELETE CASCADE; 6 indeks (user_id, status, created_at hər ikisində)
- **Docker-compose:** 3 PostgreSQL (account 5432 / auth 5433 / transaction
  5434) — portlar FƏRLİ olmalıdır
- **GORM tranzaksiya idiomu:**
```go
err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // hamısı tx ilə; xəta → avtomatik rollback
    return nil
})
```

### 8. Distributed Tranzaksiya Problemi
- **Ssenari:** Transaction DB tranzaksiya daxilində Account-a gRPC çağırışı →
  xətada yerli rollback, amma Account dəyişikliyi QALIR; rollback cəhdi
  xətanın ÖZÜNDƏ gəlsə — əks köçürmə (Pyotr→Saşa) YALNİŞ istiqamət;
  Account çatılmazsa (elektrik kəsildi) — tranzaksiya ASILI qalır.
- **Tərif:** bir neçə DB-yə yayılan tranzaksiya = **distributed**; əsas problem
  ACID (konsistensiya) pozuntusu. **Məsləhət 1: mümkünsə ümumiyyətlə YAYINDIR.**

### 9. 2PC — İki Fazalı Commit
- **Koordinator** + 2 faza:
  1. **Prepare:** iştirakçılara hazırlıq sorğusu → hər ikisi yazıları
     BLOKLAYIR və hazır deyir
  2. **Commit/Rollback:** hamısı hazırsa — commit əmri; biri belə
     çatışmırsa — hamı üçün rollback
- **Artıları:** atomarlıq + izolyasiya; sinxron — dərhal istifadəçiyə nəticə.
- **Eksiləri:** yavaş (çoxsaylı raund-triplər); koordinatora güclü asılılıq
  (yüksək yükdə boğaz); bloklanmış yazılar dar nöqtə + qarşılıqlı bloklanma.

### 10. Saga Pattern — Event-Based Həll
- **Tərif:** hər servis öz datasını dəyişəndə EVENT nəşr edir; digərləri
  abunə olub öz datasını yeniləyir — distribusiya tranzaksiyaları UMUMİYYƏTLƏ
  İSTİFADƏ OLUNMUR.
- **Kitabın implementasiyası (Kafka):**
  - Transaction: pending tranzaksiya yaradır → `transaction_data` topic-ə
    request (request_type/user_id/amount/operation_id/timestamp) göndərir
  - Account: `transaction_data`-ya abunə → HandleTransaction →
    deposit/withdraw/transfer icra → `transaction_response` topic-ə
    nəticə (result: true/false) nəşr edir
  - Transaction: `transaction_response`-a abunə → HandleAccountResponse →
    status: completed/failed
  - Statuslar istifadəçiyə "in processed" göstərməklə intermediate hal
    görünür → çaşqınlıq azalır
- **Artıları:** zəif bağlılıq (servislər bir-birini TANIMIR), bloklama YOX,
  yüksək yüklər + miqyaslana bilənlik.
- **Eksiləri:** mürəkkəblik (uğursuz + KOMPENSASIYA ssenariləri yazmaq
  lazımdır); izolyasiya pozulur (oxunuşda) — status/qarşı-tədbirlərlə
  yumşaldılır.
- **Kafka seçimi:** ödənişlər kimi həssas data — emaldan SONRA da saxlanmaq
  (replay) bacarığı.

### 11. kafka-go Wrapper (segmentio/kafka-go)
- **Producer:** Writer (Balancer: LeastBytes, BatchSize, RequiredAcks:
  RequireOne, Async: false); SendMessage (raw) + SendJSONMessage (auto
  marshal + content-type header)
- **Consumer:** Reader (GroupID, MinBytes 10KB / MaxBytes 10MB, MaxWait 1s,
  CommitInterval); Start(ctx, handler) — sonsuz döngü, ctx.Done → dayanma,
  xətada 1s gözlə və davam
- **Client interfeysi:** Publish / Subscribe / Close; Subscribe hər topic üçün
  unique groupID (base + "_" + topic) ilə AYRI consumer+goroutine yaradır
- **Env:** KAFKA_BROKER_HOST (envSeparator ","), KAFKA_CONSUMER_GROUP,
  KAFKA_TRANSACTION_TOPIC; docker-compose: confluentinc/cp-kafka:7.5.0
  (KRaft: broker+controller rolları, PLAINTEXT 9092 daxili / 29092 host) +
  kafka-ui (provectuslabs, 29093)
- **Docker tam silmə:** `docker-compose down -v` (-v: volume-larla)

### 12. YAGNI və Yekun Dizayn
- **YAGNI** ("You ain't gonna need it"): Transaction-da DELETE metodu YOX —
  maliyyə hərəkətləri HƏMİŞƏ saxlanılır; yalnız indi lazım olan funksiya yaz.
- Parametrlər üçün ayrı strukturlar (DepositParams/WithdrawParams/
  TransferParams/GetTransactionsParams) — 1 sahədən çox olanda.
- Gateway-ə transaction metodlarının əlavəsi — arxivdəki koddadır (eyni
  pattern).

## Əsas terminlələr
- НФБК (Boyce-Codd NF) — determinant = potensial açar tələbi
- Multi-valued dependency (çoxdəyərli asılılıq) — 4НF pozuntusunun kökü
- Lost update / dirty / non-repeatable / phantom read — paralellik anomaliyaları
- Serialization anomaly — ardıcıllaşdırılmış ekvivalentin olmaması
- B-tree indeks — balanslaşmış ağac (default)
- Composite (sintez) indeks — çoxsahəli axtarış üçün
- Double-entry (cüt yazılış) — DEBIT/CREDIT hərəkətləri
- Koordinator — 2PC-nin mərkəzi idarəçisi
- Saga — event əsaslı, kompensasiyalı distribusiya həlli
- Compensating transaction (kompensasiya tranzaksiyası) — əvvəlki addımın
  tərsi
- KRaft — Kafka-nın ZooKeeper-sız rejimi (broker+controller)

## Praktik Nəticə

1. **Pulu QƏPİK-lərdə saxla:** `int64(amount * 100)` — float yuvarlaqlaşdırma
   xətaları maliyyədə QƏBULOLUNMAZDIR.
2. **DB tranzaksiyası repozitoriyada, status axını servis qatında:** pending →
   (uzaq çağırış) → completed/failed; uzaq əməliyyat statusu YERLİ
   tranzaksiyadan AYRI vaxtda gəlir.
3. **Distribusiya tranzaksiyasından QAÇIN:** əvvəl birləşdirmə düşün; qaçılmazdırsa
   Saga (kompensasiya + statuslar) və ya 2PC (sıx tutarlılıq, sinxrondur).
4. **Saga-da status = izolyasiya əvəzi:** "in processed" görünən ara hallar
   anomal oxunuşların təsirini yumşaldır.
5. **Rollback miqrasiyaları TƏRS SIRADA:** indeks → PK → cədvəl.
6. **İndeks yalnız sorğu-ləvazimatlı sahələrə:** yazı xərci; B-tree default.
7. **Kafka go wrapper-ləri:** Producer/Consumer/Client strukturları —
   hər abunəlik üçün unique groupID + goroutine.

## Mənbə
Pages: 157-266 (PDF 158-267)
