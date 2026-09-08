# Chapter 13 — Protocol Buffers (səh. 260-281)

## Bu fəsil nədən bəhs edir?

Google-un Protocol Buffers (PB) serializasiya mexanizmi: `.proto` faylının
strukturu, `protoc` kod generasiyası, kompleks mesajlar (enum, repeated),
import, nested tiplər, Any, oneof, map və PB↔JSON uyğunluğu. Chapter 14-dəki
gRPC-nin təməli.

## Əsas fikirlər

### 1. Proto faylı və protoc
**Nədir:** PB — dil/platforma-müneyyən serializasiya. Mesaj tərifləri `.proto`
faylında; `protoc` aləti buradan Go (və digər dillər üçün) kod yaradır.
XML/JSON-dan daha kiçik və sürətlidir.

**Workflow (Figure 13.1):** .proto tərifi → protoc → .pb.go → kodda
istifadə → tərif dəyişərsə yenidən protoc (kod dəyişmir).

**Kitabdan kod nümunəsi (user.proto):**
```protobuf
syntax = "proto3";
package user;

option go_package="github.com/juanmanuel-tirado/savetheworldwithgo/12_protocolbuffers/pb/example_01/user";

message User {
    string user_id = 1;
    string email = 2;
}
```

**Kompilyasiya:**
```bash
protoc --go_out=$GOPATH/src user.proto
# Nəticə: user.pb.go — User tipi + marshal/unmarshal kodu
```

**Sub-izahlar:**
- `syntax = "proto3"` → PB versiyası
- `package user` → PB paketi; digər proto-lardan import oluna bilər
- `option go_package` → yaradılan Go kodunun yolu
- Sahə sintaksisi: `tip ad = tag;` → tag = marshaldəki mövqe (user_id=1
  birinci, email=2 ikinci)
- Tiplər Go-ya bənzəyir: int32, uint32, string, bool; fərqlər: float,
  double, sint64

**Go-da istifadə (marshal/unmarshal):**
```go
import (
    "github.com/juanmanuel-tirado/savetheworldwithgo/.../user"
    "google.golang.org/protobuf/proto"
)

u := user.User{Email: "john@gmail.com", UserId: "John"}

encoded, err := proto.Marshal(&u)      // Go → bytes (şəbəkə ilə göndərilə bilər)
if err != nil {
    panic(err)
}

v := user.User{}
err = proto.Unmarshal(encoded, &v)      // bytes → Go (istənilən dil/platf.)
fmt.Println("Recovered:", v.String())
```
- `google.golang.org/protobuf/proto` — aktual paket; köhnə
  `github.com/golang/protobuf` superseded (tövsiyə olunmur)
- Asılılıq üçün `go mod init` + `go build` (Chapter 12)

### 2. Kompleks mesajlar (enum + repeated)
**Kitabdan kod nümunəsi:**
```protobuf
enum Category {
    DEVELOPER = 0;
    OPERATOR = 1;
}

message Group {
    int32 id = 1;
    Category category = 2;
    float score = 3;
    repeated User users = 4;
}
```

**Go tərcüməsi:**
```go
userA := user.User{UserId: "John", Email: "john@gmail.com"}
userB := user.User{UserId: "Mary", Email: "mary@gmail.com"}

g := user.Group{Id: 1,
    Score: 42.0,
    Category: user.Category_DEVELOPER,   // enum → Paket_AD konstantu
    Users: []*user.User{&userA, &userB}, // repeated → pointer slice
}

encoded, err := proto.Marshal(&g)
recovered := user.Group{}
proto.Unmarshal(encoded, &recovered)
```

**Sub-izahlar:**
- `enum Ad { ADE = 0; }` → Go-da `Ad_ADE` konstantu; enum dəyəri ≠ field tag
- `repeated` → slice; hədd limiti yoxdur, boş ola bilər; Go-da `[]*T`

### 3. Proto importu
```protobuf
// group.proto
syntax = "proto3";
package group;
option go_package="group";

import "user.proto";       // user.proto-dakı mesajlar əlçatandır

message Group {
    // ...
    repeated user.User users = 4;   // Paket.Mesaj forması ilə
}
```
Go tərəfdə hər iki paket import edilir: `group.Group{... Users:
[]*user.User{...}}`.

### 4. Nested tiplər
**Nədir:** Yalnız başqa mesajın daxilində mənalı olan tiplər.

```protobuf
message Group {
    int32 id = 1;
    Category category = 2;
    float score = 3;
    message User {                    // daxili mesaj
        string user_id = 1;
        string email = 2;
    }
    repeated User users = 4;
}

message Winner {
    Group.User user = 1;              // Group.User istinadı
    Category category = 2;
}
```
- Go-da nested mesaj `Group_User` adlanır:
```go
userA := group.Group_User{UserId: "John", Email: "john@gmail.com"}
g := group.Group{..., Users: []*group.Group_User{&userA, &userB}}
```

### 5. Any tipi
**Nədir:** Tipi əvvəlcədən bilinməyən sahələr — ixtiyari ölçülü bayt
serializasiyası + unikal identifikator olan type URL.

```protobuf
import "google/protobuf/any.proto";

message User {
    string user_id = 1;
    string email = 2;
    repeated google.protobuf.Any info = 3;
}
```

```go
import "google.golang.org/protobuf/types/known/anypb"

info := anypb.Any{Value: []byte(`John rules`), TypeUrl: "urltype"}
userA := user.User{UserId: "John", Email: "john@gmail.com",
    Info: []*anypb.Any{&info}}

encoded, _ := proto.Marshal(&userA)
// info:{type_url:"urltype" value:"John rules"}
```

### 6. oneof tipi
**Nədir:** Bir neçə sahədən yalnız BİRİNİN göndərilməsi.

```protobuf
message Developer {
    string language = 1;
}
message Operator {
    string platform = 1;
}

message User {
    string user_id = 1;
    string email = 2;
    oneof type {
        Developer developer = 3;
        Operator operator = 4;
    }
}
```

**Go tərcüməsi (wrapper tipləri):**
```go
goDeveloper := user.Developer{Language: "go"}
userA := user.User{UserId: "John", Email: "john@gmail.com",
    Type: &user.User_Developer{&goDeveloper}}   // User_Developer wrapper

aksOperator := user.Operator{Platform: "aks"}
userB := user.User{UserId: "Mary", Email: "mary@gmail.com",
    Type: &user.User_Operator{&aksOperator}}

// Marshal ediləndə YALNIZ seçilmiş variant gedir:
// RecoveredA: ... developer:{language:"go"}
// RecoveredB: ... operator:{platform:"aks"}
```

### 7. Map-lər
```protobuf
message User {
    string user_id = 1;
    string email = 2;
}
message UserList {
    repeated User users = 1;
}
message Teams {
    map<string, UserList> teams = 1;
}
```

```go
userA := user.User{UserId: "John", Email: "john@gmail.com"}
userB := user.User{UserId: "Mary", Email: "mary@gmail.com"}

teams := map[string]*user.UserList {
    "teamA": &user.UserList{Users: []*user.User{&userA, &userB}},
    "teamB": nil,
}

teamsPB := user.Teams{Teams: teams}   // adi Go map kimi
encoded, err := proto.Marshal(&teamsPB)
recovered := user.Teams{}
proto.Unmarshal(encoded, &recovered)
// teams:{key:"teamA" value:{users:{...}}} teams:{key:"teamB" value:{}}
```

### 8. PB ↔ JSON uyğunluğu
**Nədir:** Yaradılan `.pb.go` fayllarındakı struct-lar JSON tag-ləri daşıyır
— adi `encoding/json` ilə işləyir.

```go
encoded, err := json.Marshal(&teamsPB)   // PB mesajı → JSON
recovered := user.Teams{}
json.Unmarshal(encoded, &recovered)      // JSON → PB mesajı
```
- Əlavə kod tələb olunmur; sistemlər arasında hibrid mübadilə imkanı

## Əsas terminlər
- Protocol Buffers (PB) — dil-müneyyən binary serializasiya
- proto file — mesajların IDL tərifi
- protoc — .proto → .pb.go kod generatoru
- Field tag — sahənin marshaldəki identifikatoru ( mövqeyi)
- repeated — təkrarlanan sahə (Go slice)
- oneof — "yalnız biri" sahə qrupu
- Any — tipi bilinməyən sahə (baytlar + type URL)
- .pb.go — protoc tərəfindən yaradılan Go faylı

## Praktik nəticə
PB seçimi: kiçik mesaj ölçüsü + sürətli serializasiya lazımdırsa (mikroservis
arası, gRPC). Axın: .proto yaz → protoc işə sal → .pb.go import et →
proto.Marshal/Unmarshal. Enum `Paket_A` konstantlarına, nested `Group_User`-ə,
repeated `[]*T`-yə çevrilir. oneof wrapper struct-ları ilə, Any isə anypb ilə
işlənir. Eyni mesaj JSON-a da çevrilə bilər (pb.go tag-ləri hazır).

## Mənbə
Pages: 260-281 (PDF 260-281)
