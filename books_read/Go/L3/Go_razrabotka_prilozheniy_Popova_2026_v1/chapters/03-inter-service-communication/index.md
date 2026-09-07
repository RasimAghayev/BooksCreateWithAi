# Глава 3 — Mikroservislər arası əlaqə üsulları (səh. 111-156)

## Bu fəsil nədən bəhs edir?

Sinxron (HTTP, gRPC) və asinxron (RabbitMQ, Kafka) əlaqə üsulları. OSI
modelinin 7 qatı (fiziki → tətbiqi), HTTP strukturu (metodlar, sorğu/cavab),
gRPC (HTTP/2 + Protobuf, 4 ünsiyyət nümunəsi: unary/server-stream/client-stream/
bidirectional; validasiya limitləri — enum 0 problem), RabbitMQ/AMQP (exchange
növləri: fanout/direct/topic; delivery garantiyaları: At Most/Least/Exactly
Once), Kafka (pull model, partition/offset, topic növləri, Connect API),
Redis (in-memory key-value, cache/queue), API Gateway pattern — JWT
interceptor (publicMethods, Bearer, claims, context-ə userID), Gateway
servisinin auth+account inteqrasiyası (grpc.NewClient, insecure credentials,
interfeys izolyasiyası), /users/me üçün GetUserIDFromContext.

## Əsas fikirlər

### 1. Sinxron vs Asinxron
- **Sinxron:** cavab DƏRHAL lazımdır → HTTP, gRPC
- **Asinxron:** cavab "sonra da ola bilər" → Kafka, RabbitMQ (nümunə: pul
  çıxarışı əməliyyatı + sonra bildiriş)

### 2. OSI Modeli — 7 Qat
| # | Qat | Nə edir | Nümunə |
|---|-----|---------|--------|
| 1 | Fiziki | bitlərin ötürülməsi (optik/kabel/radio); konnektorlar, gərgilik | şəbəkə adapteri |
| 2 | Kanal (Data Link) | kadr (frame) qruplaşdırma, checksum, xəta düzəltmə (yenidən göndərmə) | MAC |
| 3 | Şəbəkə (Network) | IP: PAKETLƏR ünvanlanır; ZƏMANƏT YOXDUR (sıra pozula, itə, dublikat) | router |
| 4 | Nəqliyyət (Transport) | TCP: bağlantı + tam zəmanət (itən = yenidən sorğu, dublikat təmizlənir); UDP: datagram — bağlantısız, gecikmə kritik (multimedia) | TCP/UDP |
| 5 | Sessiya | seansın qurulması/saxlanması; sinxron audio/video (qrup zəngi) | RPC, RTCP, SMPP |
| 6 | Təqdimat (Presentation) | format çevirmə: sıxma, kodlama, ŞİFRƏLƏMƏ | SSL/TLS, ISO 8583 |
| 7 | Tətbiqi (Application) | istifadəçiyə çatan format | HTTP, SMTP, POP3 |

- **Qruplaşma:** 1-3 şəbəkə-asılı; 5-7 tətbiq-yönlü; 4 (transport) — aradakı,
  aşağı detalları yuxarıdan gizlədən qat
- **Sərhəd qeydi:** HTTPS vəziyyətən görə tətbiqi VƏ YA təqdimat qatı ola
  bilər (ISO 8583 şifrələmə daşıyıcısı kimi)

### 3. HTTP Strukturu
- **Sorğu:** metod (GET — oxu; HEAD — GET bədənsiz; POST — yarat; PUT — TAM
  əvəz et — ötürülməyənləri SİL; PATCH — HİSSƏVİ yenilə; DELETE — sil;
  CONNECT — tunel; OPTIONS — parametrlər; TRACE — debug) + URL + versiya +
  başlıqlar + body (POST/PUT)
- **Cavab:** versiya + status kodu (1xx info, 2xx uğur, 3xx redirect, 4xx
  klient, 5xx server) + mesaj + başlıqlar + body
- **Proxy/nginx:** hər sorğu üçün axın parçalanır — işçi connection-lar
  (hərəsi ayrıca proses, 1024 sorğuya qədər); nəticələr blokla birləşir

### 4. gRPC — HTTP/2 + Protobuf
```protobuf
service UserService {
  rpc GetUsers(UsersRequest) returns (UsersResponse);
}
message UsersRequest {
  repeated string client_ids = 1;   // massive = repeated
  uint32 take = 3;  uint32 skip = 4;  // float32/float64 də var
}
```
- **Protobuf:** Google; XML-in binar alternativi — minimal ölçü, sürət;
  sahələr İNDEKSLənir (ad fərqli, indeks+tip eyni olsa da qəbul edilir →
  qarışıqlıq təhlükəsi → ümumi contracts repo)
- **HTTP metodları YOXDUR:** bütün sorğular eyni görünür; ayrım ADLA
  (CreateUser, GetUserById) — developer özü qərar verir
- **İkili protokol dəstəyi:** `import "google/api/annotations.proto"` +
  `option (google.api.http) = { get: "/v1/users" }` — eyni servis həm gRPC,
  həm REST
- **4 kommunikasiya nümunəsi:**
  1. Unary — 1 sorğu → 1 cavab (funksiya çağırışı kimi)
  2. Server streaming — 1 sorğu → mesaj AXINI
  3. Client streaming — mesaj axını → 1 cavab (yazıb bitirəndə gözlə)
  4. Bidirectional — hər iki tərəf müstəqil oxu/yazı axınları
- **Validasiya limiti:** binar serializasiyada "array-mi, tək dəyər-mi"
  ayırd etmək olmaz; **enum-0 problem:** Protobuf 0-ı avtomatik qoyur →
  UNDEFINED ilk dəyər konvensiyası ilə mübarizə
- **Test:** proto faylları lazımdır (Postman, BloomRPC, Enevs — konsol);
  health-check mürəkkəbləşir; gücü: sürət + kontrakt "box-dan"

### 5. RabbitMQ — AMQP Brokeri
- **AMQP konseptləri:**
  - **Message** — məzmun broker tərəfindən interpretasiya OLUNMUR; başlıqlar
    strukturqlu
  - **Exchange** — mesajlar buraya gəlir; növləri:
    - *Fanout* — BÜTÜN bağlanmış queue-lara
    - *Direct* — routing key = queue adı ilə dəqiq uyğunluq
    - *Topic* — maska uyğunluğu: `app.notification.sms.*` → sms prefiksli
      bütün açarlar
  - **Queue** — mesajlar götürülənə qədər saxlanılır
- **İştirakçılar:** broker (vasitəçi), publisher (göndərən), consumer (abunə)
- **Axın:** qoşul → queue yarat → publish → route → saxlanış → çatdırılma →
  **acknowledgment** (uğurlu emaldan sonra broker mesajı SİLİR)
- **Delivery garantiyaları:**
  - *At Most Once* — zəmanət YOX; itə bilər, retry YOX
  - *At Least Once* — mütləq çatır; retry var → DUBLİKAT mümkün
  - *Exactly Once* — dəqiq bir dəfə; ən çətin — idempotentlik + unikal ID
    tələb edir (uygulama səviyyəsində həll)
- **Kanallar:** bir bağlantıda çox kanal (paralel komandalar); kanal = yaddaş
  xərci; hər əməliyyət üçün YENİ bağlantı AÇMA (resurs israfı)
- **Fon emalı:** böyük sənəd generasiyası HTTP timeout-a düşməsin —
  nəticə RabbitMQ-ya → növbədən çatdırılma/bildiriş

### 6. Apache Kafka
- **Pull model (RabbitMQ push-dan fərqli):** consumer özü soruşur
  (bir neçə saniyəyə bir) → batch qruplaşdırma + yüksək throughput; gecikmə
  mümkün, balans pozula bilər; push isə balans saxlayır, amma QoS limitləri
  tələb edir
- **Saxlanış fərqi:** mesajlar emaldan SONRA DA saxlanılır (tərsil mümkün);
  təkrar oxunma problem yaratmır — **offset** (partition daxili mövqe) +
  yaradılma vaxtı indeksi
- **Struktur:** producers → key-value mesajlar → **topics** → **partitions**
  (içərisində SIRALI); distributed cluster — broker node-ları; partition
  REPLİKASIYASI (etibarlılıq); 0.11+ tranzaksiya modeli → Streams API ilə
  "bir dəfə emal"
- **Topic növləri:**
  - *Adi* — müddət/yer limiti; köhnələnlər avtomatik silinir
  - *Compact* — hər açar üçün YALNIZ son dəyər saxlanılır; boş dəyər göndər
    → açarın silinməsi
- **Connect API:** xarici sistemlərdən import/export konnektorları (0.9+);
  dil reallaşdırmaları cəmiyyət tərəfindən

### 7. Redis — In-Memory Key-Value
- **Yerləşmə:** PostgreSQL kimi disk DB-lərin ƏLAVƏSİ (əvəzedici deyil) —
  RAM-də → sürətli; RAM diskdən BAHA → tamamlayıcı
- **Strukturlar:** string, hash, bitmassiv, list, set, bitfield, sorted set,
  geo, HyperLogLog, stream
- **İstifadə halı:**
  - Az dəyişən, tez-tez çatılan (header-lər, api-key-lər)
  - Tez dəyişən, kritik olmayan (caches hesabatlar)
- Keçmiş: Memcached kimi keş → indi queue, stream emal da

### 8. API Gateway Pattern
- **Problemlər (gateway-siz):** klient çox sayda ünvanı bilməli; hər çağırış
  ayrıca avtorizasiya; topologiya dəyişikliyi klient kodu dəyişir; təkrar
  infrastruktur logiki; təhlükəsizlik/monitorinq çətin
- **Həll:** tək giriş nöqtəsi — "gömrük və dirijor": avtorizasiya yoxlayır,
  başlıq əlavə edir, format çevirir, yönləndirir; cavabı birləşdirir
  (mobil "profil" = 3 servisdən aqreqasiya)
- **Mərkəzləşdirilən cross-cutting:** token yoxlama, rate limit, keş, API
  versiya routransı, data aqreqasiyası

### 9. Gateway Kontraktı
```protobuf
service Gateway {
  // auth: /api/v1/auth/register|login|refresh|logout|validate
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  ...
  // users (autentikasiya tələb edir): /api/v1/users + {user_id} + /me
  rpc GetCurrentUser(google.protobuf.Empty) returns (GetCurrentUserResponse);
  rpc UpdateCurrentUser(UpdateCurrentUserRequest) returns (...);
  rpc DeleteCurrentUser(google.protobuf.Empty) returns (...);
}
message RegisterRequest { account.User user = 1; string password = 2; }
message RegisterResponse { account.User user = 1; auth.TokenPair tokens = 2; }
```
- **Yeni /me metodu:** user_id token-dən çıxarılır (path-də YOX)
- Kontrakt digər servis modellərini (account.User, auth.TokenPair) İMPORT
  edir — fasad hər şeyi təqdim edir

### 10. JWT Interceptor — gRPC Middleware
```go
type JWTClaims struct {
    UserID uint64 `json:"user_id"`
    jwt.RegisteredClaims
}

var publicMethods = map[string]bool{
    "/gateway.Gateway/Register": true,
    "/gateway.Gateway/Login":   true,
    "/gateway.Gateway/Refresh":  true,
}
```
- **Növlər:** Unary Interceptor (adı çağırışlar) + Stream Interceptor (axın)
- **Axın (unary):**
  1. publicMethods yoxlaması → atla
  2. metadata-dan Authorization başlığı (`metadata.FromIncomingContext`)
  3. Bearer prefiksi sil → JWT parse (`ParseWithClaims` + secret)
     - **Alg yoxlaması:** `token.Method.(*jwt.SigningMethodHMAC)` — gözlənilən
       imza ÜSULU təsdiqlənməlidir
  4. claims-dən UserID çıxar
  5. `ctx = context.WithValue(ctx, UserIDKey, userID)`
  6. handler çağır
- **Stream wrapper:** ServerStream embed + Context() override — yeni
  konteksti ötürmək üçün
- **Faydası:** auth-a şəbəkə çağırışı YOX (lokal JWT validasiyası) → sürətli;
  bütün qorunmuş metodlar avtomatik əhatə olunur — dublikatsız
- **Handler-də:** `userID, err := interceptor.GetUserIDFromContext(ctx)`

### 11. Servis İnteqrasiyası (auth + account)
```go
// client yaratma (konfiqurasiyadan host):
conn, _ := grpc.NewClient(cfg.AuthGrpcHost,
    grpc.WithTransportCredentials(insecure.NewCredentials()))
client := authpb.NewAuthClient(conn)
a.authService = auth.NewService(client)
```
- `.env`: `ACCOUNT_GRPC_HOST=localhost:50051`, `AUTH_GRPC_HOST=localhost:50052`
- **İzolyasiya qatı:** internal/auth, internal/account paketləri — gRPC
  client-ləri model+mapper ilə bürüyür; daxili servis dəyişsə — YALNIZ bir
  fayl düzəlir
- **GatewayService interfeyslərlə işləyir** (AccountService, AuthService) —
  biznes məntiqi gRPC haqqında bilir YOX

### 12. Kompozit Register Əməliyyatı
```go
func (s *GatewayService) Register(ctx, newUser, password) (User, TokenPair, error) {
    user, err := s.accountService.CreateUser(ctx, newUser)  // 1. yarat (PAROLSUZ!)
    s.authService.Register(ctx, user.ID, login, email, password)  // 2. auth-da qeyd
    token, err := s.authService.Login(ctx, login, password)  // 3. dərhal token
}
```
- **Parol Account-a ÖTÜRÜLMÜR** — yalnız şifrələnmiş halda auth-da
- DeleteUser: hər ikisində silinir (auth + account)

## Əsas terminlər
- OSI modeli — 7 qatlı şəbəkə standartı
- Frame/Packet/Datagram — kanal/şəbəkə/UDP vahidləri
- Exchange (dəyişim nöqtəsi) — AMQP router (fanout/direct/topic)
- Routing key — mesajın queue seçim açarı
- Acknowledgment (təsdiq) — uğurlu emaldan sonra mesaj silinir
- At Most/Least/Exactly Once — çatdırılma zəmanət səviyyələri
- Pull/Push model — consumer özü soruşur / server özü göndərir
- Partition/Offset — Kafka topic bölgüsü / mövqe indeksi
- Compact topic — hər açarın son dəyəri saxlanılır
- Interceptor (müdaxiləçi) — gRPC middleware: avtorizasiya, metrik
- Cross-cutting concern (kəsən maraq) — bütün metodlara aid ümümi məntiq
- API Gateway (fasad) — xarici dünya üçün tək giriş

## Praktik nəticə

1. **Əlaqə üsulu seçimi:** dərhal cavab → HTTP/gRPC; bildiriş/fon emalı →
   Kafka (stream, replay) / RabbitMQ (queue, routing).
2. **gRPC kontraktı REST-lə paralel:** google.api.http annotations — eyni
   servis 2 protokolla; amma validasiya limitlərini (enum-0, repeated)
   nəzərə al.
3. **JWT yoxlaması Gateway-də interceptor ilə:** auth servisinə hər dəfə
   çağırış YOX; publicMethods siyahısı açıq uçot üçün.
4. **Alg yoxlaması ZƏRURİ:** ParseWithClaims callbackində imza üsulunu
   təsdiqlə — token.Method.(*jwt.SigningMethodHMAC).
5. **Parol yalnız auth-da:** Account parolu bilir YOX; Register kompozisiyası
   = account.CreateUser + auth.Register + auth.Login.
6. **Xarici servis client-ləri izolə:** xüsusi paket + mapper; interfeys
   arxasında — gRPC detalları biznes məntiqinə sızmasın.
7. **Kafka vs RabbitMQ:** replay/son dəyər lazımdırsa — Kafka; routing
   qaydaları + acknowledgment — RabbitMQ.

## Mənbə
Pages: 111-156 (PDF 112-157)
