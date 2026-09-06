# Chapter 8 — Слабая связанность

## Bu chapter nədən bəhs edir?

Paylanmış sistemlərdə komponentlərin bir-birindən ne qədər asılı olduğunu — zəif əlaqəlilik (loose coupling) — izah edir. Protokol çeşidləri (REST, gRPC, SOAP, pub/sub), vaxt əlaqəliliyi (time coupling), plugin sistemləri və Heksagonal Arxitektura (ports and adapters) müzakirə edilir.

## Əsas fikirlər

### 1. Zəif əlaqəlilik (Loose Coupling) və yağılı bağlılıq (Tight Coupling)
**Nədir:** Komponentlərin bir-birindən asılılıq dərəcəsi.

**Necə işləyir:**
- **Tight coupling (sıx bağlılıq):** Bir komponent dəyişsə, bir neçə komponent dəyişmək məcburiyyətində qalır — distribüed monolit yaradır
- **Loose coupling (zəif əlaqəlilik):** Komponentlər standart protokollar və müəyyən müqavilələr (contracts) üzərində əlaqələnir — bir tərəf dəyişsə, digəri təsir olmadan işləməyə davam edir

### 2. Protokol çeşidləri
**Nədir:** Komponentlər arasında data və əmr mübadiləsi üçün istifadə edilən standartlar.

**Necə işləyir:**
- **REST:** Uniform interfeys, HTTP metodları (GET, POST, PUT, DELETE), JSON/XML. Sadə, lakin kontrakt yoxdur — tərəflər arasında ifrat sövdələşmə (over-negotiation) yarana bilər.
- **gRPC:** Kontrakt ilk (contract-first) dizayn. `.proto` faylından avtomatik olaraq Go interfeysləri (stub) yaradılır. HTTP/2 üzərində, binary serialization (Protocol Buffers), yüksək performans.
- **SOAP:** Çox ağır, XML əsaslı, tight coupling yaradır — "anti-pattern" hesab edilir.
- **Pub/Sub (Yayın/Abunə):** Zaman əlaqəliliyini azaldır — istehsalçı (publisher) və istehlakçı (consumer) eyni anda olmasına ehtiyac yoxdur.

### 3. Vaxt əlaqəliliyi (Time Coupling)
**Nədir:** Komponentlərin eyni anda işləməsi və ya müəyyən bir sıra ilə sorğu göndərməsi ehtiyacı.

**Necə işləyir:**
- **Sinxron (synchronous):** Sorğu göndərir, cavabı gözləyir — zaman əlaqəliliyi var
- **Asinxron (asynchronous):** Sorğu göndərir, cavabı gələndə alır — zaman əlaqəliliyi azalır
- Pub/Sub modelləri vaxt əlaqəliliyini tamamilə aradan qaldırır

### 4. Plugin sistemləri
**Nədir:** Tətbiyin funksionallığını yenidən kompilyasiya etmədən dinamik olaraq genişləndirmək imkanı.

**Necə işləyir:**
- Go standart `plugin` paketi: eyni Go versiyası və build teqləri (build tags) tələb edir
- `hashicorp/go-plugin`: RPC əsaslı, fərqli binary-lər arasında işləyir — ən güclü plugin sistemi
- Plugin-lər interface vasitəsilə qeyd edilir və `pluginMap` ilə yüklənir

**Kitabdan kod nümunəsi (Hashicorp plugin konsepti):**
```go
type Sayer interface {
    Say() error
}

type SayerPlugin struct {
    Impl Sayer
}

var SayerPluginMap = map[string]plugin.Plugin{
    "sayer": &plugin.SayerPlugin{},
}
```

### 5. Heksagonal Arxitektura (Ports and Adapters)
**Nədir:** Alıcı (core) domain məntiqini xarici interfeyslərdən (REST, gRPC, CLI, DB) izolyasiya edən arxitektura modeli.

**Necə işləyir:**
- **Ports:** Alıcı tərəfindən təyin edilən interfeyslər (müqavilələr) — "necə işləməli"
- **Adapters (Adaptorlar):** Portları həyata keçirən xarici kod — REST, gRPC, CLI, verilənlər bazası
- **Core:** Biznes məntiqi, heç bir xarici asılılıq olmadan
- Xarici interfeys dəyişərsə, yalnız uyğun adapter dəyişir, core qalır

**Üstünlükləri:**
- Testləmə asanlığı — adapter-lər mock edilə bilir
- Yeni interfeyslər əlavə etmək asan — yeni adapter yaradılır
- Core təmiz və test edilə bilir

## Əsas terminlər
- Loose coupling (zəif əlaqəlilik) — minimal asılılıq
- Tight coupling (sıx bağlılıq) — yüksək asılılıq, distribüed monolit
- REST — uniform interfeys üzərində HTTP
- gRPC — contract-first, HTTP/2, binary serialization
- SOAP — ağır, XML əsaslı, tight coupling
- Pub/Sub — asinxron mesajlaşma
- Time coupling (vaxt əlaqəliliyi) — sinxron/asinxron asılılığı
- Plugin — dinamik funksionallıq genişlənməsi
- Hexagonal Architecture (Heksagonal Arxitektura) — core/ports/adapters
- Port — interfeys müqaviləsi
- Adapter — portun həyata keçirilməsi

## Praktik nəticə
Zəif əlaqəlilik paylanmış sistemlərin müstəqil inkişafı, deploy və miqyaslanması üçün əsas şərtdir. REST və gRPC ən çox istifadə edilen protokollərdir, pub/sub isə vaxt əlaqəliliyini tamamilə azaldır. Heksagonal arxitektura biznes məntiqini xarətən asılılıqlardan izolyasiya edərək təkcə testləməni deyil, həm də sistemin davamlılığını artırır.

## Mənbə
Pages: 234-277 (PDF səh. 234-277)
