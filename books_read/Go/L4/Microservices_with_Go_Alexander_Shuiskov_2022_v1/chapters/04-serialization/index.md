# Chapter 4 — Serialization (səh. 77-92)

## Bu fəsil nədən bəhs edir?

Serializasiyanın əsasları, populyar formatlar (XML, YAML, Thrift, Avro),
Protocol Buffers ilə şema təyini, kod generasiyası və JSON/XML/Proto
ölçü-sürət müqayisəsi (benchmark).

## Əsas fikirlər

### 1. Serializasiya nədir?
- **Serialization (Serializasiya):** datanın ötürülə/bilə saxlana bilən formata
  çevrilməsi; əks proses — **Deserialization**
- 2 əsas istifadə: servislərarası ötürmə (ümumi dil) + mürəkkəb strukturların
  byte array kimi saxlanması (DB, konfiq, loglar)

### 2. Format müqayisəsi
| Format | İl | Xüsusiyyət | Problem |
|---|---|---|---|
| XML | 1998 | element ağacı, geniş qəbul | ən BÖYÜK ölçü, yavaş |
| JSON | — | oxunaqlı, browser dəstəyi | ölçü/sürət optimal deyil |
| YAML | 2001 | insan üçün oxunaqlı, şərhlər | konfiq üçün; servis-gözə YOX |
| Thrift | Facebook | serializasiya + RPC protokolu; struct + **service** təyini | populyarlıq azalıb, sənədləmə zəif |
| Avro | — | JSON/IDL şema; **versiyalı şema** — köhnə+yeni birgə | — |
| **Protobuf** | Google, 2008 | sadə dil, kiçik output, sürətli, service təyini, çoxdilli kodgenerasiya | binary — oxunaqlı deyil |

### 3. Protocol Buffers praktikası
```proto
syntax = "proto3";
option go_package = "/gen";

message Metadata {
    string id = 1;
    string title = 2;
    string description = 3;
    string director = 4;
}
message MovieDetails {
    float rating = 1;
    Metadata metadata = 2;
}
```
- Quraşdırma: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- Generasiya: `protoc -I=api --go_out=. movie.proto` → `gen/movie.pb.go`
- Sahə NÖMRƏLƏRİ (1,2,3...) wire formatda açar rolunu oynayur — adlar yox

**Şemanın faydaları:**
- Explicit schema — data tipləri koddan ayrıqı, API aydındır
- Kod generasiyası (Go + digər dillər) — model dəyişsə, tək komanda
- Cross-language dəstək

### 4. Benchmark nəticələri (Metadata strukturu)
| Format | Ölçü | Sürət (ns/op) |
|---|---|---|
| XML | 148B | 2519 (13x yavaş) |
| JSON | 106B | 342 (2x yavaş) |
| **Protobuf** | **63B** | **185.7 (ən sürətli)** |

- Protobuf: JSON-dan 40% kiçik, 2x sürətli; XML-dən 2x+ kiçik, 13x sürətli
- Benchmark: `go test -bench=.` — b.N iterasiya, ns/op orta sürət

### 5. Best practices
- **Backward compatibility:** sahə adını silmə/adlandırma/tip dəyişmə = break
- Client-server şema versiyaları sinxron saxla
- İmplicit detalları şemada şərh et (məs. boş sahə qadağan)
- Vaxt üçün `google.protobuf.Timestamp` (int timestamp YOX)
- Ardıcıl adlandırma + rəsmi stil qidları

## Termindirmə (AZ)
- Serialization — Serializasiya (datanın ötürülə bilən formata çevrilməsi)
- Schema — Şema (datanın strukturlu tərifi)
- Wire Format — Ötürülmə Formatı (şəbəkədəki binary təsvir)
- Backward Compatibility — Geri Uyumluluq
- Benchmark — Performans Ölçməsi
- Code Generation — Kod Generasiyası

## Kviz sualları
1. Protobuf JSON-dan nə qədər kiçik/sürətlidir? (40% kiçik, ~2x sürətli)
2. Avro-nun unikal üstünlüyü nədir? (Versiyalı şemalar — köhnə/yeni formatlar
   paralel saxlanır, incompatible dəyişikliklər də idarə olunur)
3. Protobuf-da sahə nömrələri niyə vacibdir? (Wire formatda sahələr nömrə ilə
   tanınır — ad dəyişməsi break etmir, nömrə dəyişməsi break edir)
4. Timestamp üçün protobuf-də nə istifadə olunmalıdır? (google.protobuf.Timestamp)
