# Chapter 3 — Getting up and running with gRPC and Golang (gRPC və Go ilə işə başlama)

## Bu chapter nədən bəhs edir?
Protobuf dərinliyinə: mesaj təyinatı (field rules/types/numbers), binary encoding mexanikası (varint, MSB), protoc ilə stub generasiyası, proto-ların ayrı repo-da saxlanması + GitHub Actions avtomatizasiyası və backward/forward compatibility qaydaları.

## Əsas fikirlər

### 1. Protobuf mesaj təyini
```protobuf
syntax = "proto3";

message CreateOrderRequest {
    int64 user_id = 1;        // sahib
    repeated Item items = 2;  // sifariş elementləri
    float amount = 3;         // ödəniləcək məbləğ
}
```
**Field komponentləri:**
- **Rules:** Singular (max 1, default) / Repeated (0+, sıra qorunur — Go-da slice)
- **Types:** scalar (string/int64/float...), enum, embedded mesaj (Item)
- **Names:** lowercase + underscore (`user_id`); compiler çoxdilli generasiya üçün qaydalara ehtiyac duyar
- **Numbers:** sahələrin binary-formada unikal identifikatoru — DƏYİŞDİRİLMƏZ!

### 2. Reserved sahələr — compatibility qorunması
```protobuf
message CreateOrderRequest {
    reserved 1, 2, 3 to 7;         // nömrə ilə (və range!)
    reserved "customer_id";         // AD ilə
    int64 user_id = 7;
    repeated Item items = 8;
    float amount = 9;
}
```
**Niyə:** köhnə nömrə/ad yeni TIPLİ sahəyə verilsə — köhnə client + yeni server = data uyğunsuzluğu. Reserved = gələcəkdə təkrar istifadəyə qadağa.

**Performans qaydası:** 1-15 arası nömrələr metadata-da 1 bayt; 16-2047 — 2 bayt → TEZ-İSTİFADƏ olunan sahələri (correlation_id kimi) 1-15-də saxla!

### 3. Protobuf encoding — binary wire format
**Nümunə:** `UserId: 65`, field №1, int → marshal → `[]byte`:

**Metadata baytı:** wire type (ilk 3 bit: 000 = Varint int üçün) + field number.
**Data baytı:** MSB (devam biti) + 7 bit data.

**65 = 1 bayta sığır** → metadata 1 bayt + data 1 bayt.

**Böyük dəyər nümunəsi (21567):**
1. 21567 → binary: `101010000111111`
2. 7-bit bloklara böl: `0000001-0101000-0111111`
3. Sıranı tərs çevir (little-endian): `0111111-0101000-0000001`
4. MSB: ilk iki blokda 1 (davam var), üçüncüdə 0 (son)

**Nəticə:** 127+ dəyərlər çoxbaytlı; field № 15+ metadata çoxbaytlı → performans hər iki tərəfdən asılıdır.

### 4. protoc — stub generasiyası
**Qurulum:**
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Generasiya:**
```bash
protoc -I ./proto \
   --go_out ./golang \
   --go_opt paths=source_relative \
   --go-grpc_out ./golang \
   --go-grpc_opt paths=source_relative \
   ./proto/order.proto
```
**Flag izahı:**
- `-I` — import axtarış yolu
- `--go_out` — mesaj kodlarının yeri (order.pb.go)
- `--go-grpc_out` — servis funksiyaların yeri (order_grpc.pb.go)
- `paths=source_relative` — eyni qovluq strukturu qorunur

**Konvensiya:** `New<ServiceName>Client(...)` → `<ServiceName>Client` interface; istifadə:
```go
client := order.NewOrderClient(...)
client.Create(ctx, &CreateOrderRequest{UserId: 123})
```

### 5. Ayrı proto repo strukturu
**Niyə ayrı:** Go mikro servis layihəsində Java kodu saxlamaq yanlış; hər dil consumer-ı öz dependency-sini alsın.

```
microservices-proto/
├── golang/
│   ├── order/     (go.mod, go.sum, order.pb.go, order_grpc.pb.go)
│   ├── payment/   (...)
│   └── shipping/  (...)
├── order/order.proto
├── payment/payment.proto
└── shipping/shipping.proto
```

**Git tagging (Go module subfolder convention):**
```bash
git tag -a golang/order/v1.2.3 -m "golang/order/v1.2.3"
git push --tags
# Consumer tərəfi:
go get -u github.com/huseyinbabal/microservices-proto/golang/order@v1.2.3
```

### 6. GitHub Actions avtomatizasiyası
**run.sh (bash script):**
```bash
#!/bin/bash
SERVICE_NAME=$1
RELEASE_VERSION=$2
sudo apt-get install -y protobuf-compiler golang-goprotobuf-dev
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=./golang --go_opt=paths=source_relative \
  --go-grpc_out=./golang --go-grpc_opt=paths=source_relative \
  ./${SERVICE_NAME}/*.proto
cd golang/${SERVICE_NAME}
go mod init github.com/huseyinbabal/microservices-proto/golang/${SERVICE_NAME} || true
go mod tidy
cd ../../
git add . && git commit -am "proto update" || true
git tag -fa golang/${SERVICE_NAME}/${RELEASE_VERSION} \
  -m "golang/${SERVICE_NAME}/${RELEASE_VERSION}"
git push origin refs/tags/golang/${SERVICE_NAME}/${RELEASE_VERSION}
```

**Workflow (matrix strategy ilə):**
```yaml
name: "Protocol Buffer Go Stubs Generation"
on:
  push:
    tags: [ "v**" ]
jobs:
  protoc:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service: ["order", "payment", "shipping"]   # hər servis üçün job
    steps:
      - uses: actions/setup-go@v2
        with: { go-version: 1.17 }
      - uses: actions/checkout@v2
      - run: echo "RELEASE_VERSION=${GITHUB_REF#refs/*/}" >> $GITHUB_ENV
      - shell: bash
        run: |
          chmod +x "${GITHUB_WORKSPACE}/run.sh"
          ./run.sh ${{ matrix.service }} ${{ env.RELEASE_VERSION }}
```
**Axın:** tag push → workflow trigger → hər servis üçün stub generate → go mod init → commit → yeni tag.

### 7. Backward/Forward compatibility
**Backward (geriyə uyğunluq):** yeni server + köhnə client işləsin (istifadəçi güclə yenilənməsin).
**Forward (irəli uyğunluq):** köhnə server yeni input-u bağışlasın (brauzer/yeni HTML kimi).

| Ssenari | Nəticə |
|---|---|
| **Yeni sahə əlavə** | ƏSASƏN təhlükəsiz — reserved yoxlamanı unutma |
| **Server yeni, client köhnə** | Server default dəyərlə işləsin: |
| **Client yeni, server köhnə** | Server yeni sahəni İGNORR edir — xətasız amma itir |

**Default dəyər pattern-i:**
```go
vat := VAT
if req.Vat > 0 { vat = req.Vat }   // köhnə client-siz də düzgün işləsin
return &CreatePaymentResponse{TotalPrice: vat + req.Price}, nil
```

**oneof tələləri:**
```protobuf
message CreatePaymentRequest {
    oneof payment_method {
        CreditCard credit_card = 1;
        PromoCode promo_code = 2;
    }
}
```
- oneof-dan sahə SİLMƏK = **backward-incompatible** (köhnə client promo_code göndərir → server-də İTİR)
- oneof-a sahə ƏLAVƏ = forward-incompatible
- oneof-aİRİ/ÇIXARMaq = data loss (ikisi birlikdə göndərilə bilməz — biri düşür)
- **Breaking change → semver MAJOR** bump + release notes

## Əsas terminlər
- proto3 — sadələşdirilmiş syntax (proto1 deprecated, proto2 köhnə)
- Field rule — singular/repeated
- Field number — binary identifikator; 1-15 = 1 bayt
- Varint — int wire type (type 0)
- MSB (Most Significant Bit) — davam biti
- Little-endian — blokların tərs sırası
- reserved — nömrə/adın gələcək istifadəyə qorunması
- protoc / protoc-gen-go / protoc-gen-go-grpc
- go_package — generasiya olunan modulun yolu
- Matrix strategy — çox-job workflow dəyişəni
- Backward/Forward compatibility
- oneof — ya-o-ya sahə qrupu
- semver — breaking change major

## Praktik nəticə
1. Tez-istifadə sahələrinə 1-15 nömrə ver; silinən sahələri mütləq `reserved` et.
2. Proto-ları ayrı repo-da saxla, dil-özəl qovluqlar + git tag convention (golang/order/v1.2.3).
3. Stub generasiyasını GitHub Actions matrix ilə avtomatlaşdır — tag push → generate → tag.
4. Yeni sahələrdə server tərəfdə default fallback yaz — köhnə client qırılmasın.
5. oneof dəyişiklikləri həmişə breaking-dir — semver major + release notes.
6. Client yeni/server köhnə ssenarisində data İTİRİLİR (xəta YOX) — səssiz bug, diqqət!

## Mənbə
Pages: 31-46 (PDF səh. 48-63)
