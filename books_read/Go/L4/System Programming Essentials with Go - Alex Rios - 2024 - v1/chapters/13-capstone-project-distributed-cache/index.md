# Chapter 13 — Capstone Project: Distributed Cache (Paylanmış Keş)

## Bu chapter nədən bəhs edir?
Kitabın final layihəsi: sıfırdan paylanmış (distributed) keş sistemi — thread safety (sync.RWMutex), HTTP interface (TCP/gRPC/WebSocket müqayisəsi), eviction (TTL + LRU — container/list ilə), P2P replikasiya (X-Replication-Request header) və consistent hashing (hash ring, SHA-1) ilə sharding. Hər addımda trade-off analizi.

## Əsas fikirlər

### 1. Paylanmış keş — nədir, nəyə lazım?
Tətbiq ilə əsas data store arasındakı yaddaş qatı: tez-tez oxunan datanı yaddaşda saxlayır → latency aşağı, read throughput yuxarı. Çətinliklər: çoxlu node üzrə tutarlılıq (CAP theorem), node xətalarına davamlılıq, eviction idarəsi, datanın köhnəlməməsi.

**Tələblər (DNA):**
- **Performance:** ms cavab, sürətli strukturlar
- **Scalability:** horizontal — sharding ilə yeni node-lar
- **Fault tolerance:** replikasiya, node xətasına zərif davranış
- **Eviction:** TTL + LRU ilə yaddaş idarəsi
- **Monitoring:** hit/miss nisbəti, alert-lər
- **Security:** auth, şifrələmə, təhlükəsiz kanal

### 2. Baza cache + thread safety
```go
type CacheItem struct {
    Value      string
    ExpiryTime time.Time
}
type Cache struct {
    mu       sync.RWMutex
    items    map[string]CacheItem
}

func (c *Cache) Set(key, value string, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = CacheItem{Value: value, ExpiryTime: time.Now().Add(ttl)}
}
func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, found := c.items[key]
    if !found || time.Now().After(item.ExpiryTime) { return "", false }
    return item.Value, true
}
```
**Sub-kod izahı:** yazma üçün tam Lock, oxuma üçün RLock — paralel oxumaqa icazə, yazma eksklüziv.

**Concurrency seçimləri:**

| Variant | Nə vaxt |
|---|---|
| sync.Mutex | sadəlik; amma ağır yükdə contention |
| sync.RWMutex | oxuma > yazma (bu layihə seçimi) |
| sync.Map / ristretto / go-cache / golang-lru | tez oxu, nadir yazı; hazır feature-lar |
| Atomic əməliyyatlar | lock-free, dizaynı çətin |
| Channel + seriya goroutine | bütün əməliyyatlar kanaldan |
| Sharded cache (hissə-hissə kilid) | maksimal performans — yük kilidlər arası bölünür |

### 3. Interface seçimi — HTTP
| | TCP | HTTP | gRPC | WebSocket |
|---|---|---|---|---|
| Sürət | ən yüksək (aşağı overhead) | header/tekst overhead | yüksək (HTTP/2 + Protobuf) | yüksək |
| Sadəlik | çətin (manual connection, framing) | asan (request/response) | orta (HTTP/2+proto tələbi) | çətin (long-lived idarəsi) |
| Standartlıq | app-level protokol yox | universal | güclü kontrakt | real-time ikiistiqamət |
| State | connection-oriented | stateless → scaling asan | — | persistent |

**HTTP seçimi:** sadəlik + standart + zəngin ekosistem + stateless scaling; tədris layihəsi üçün konseptlərə fokus. Sonradan başqa protokol əlavə oluna bilər.

**CacheServer:**
```go
type CacheServer struct {
    cache *Cache
    peers []string
    mu    sync.Mutex
}
func (cs *CacheServer) SetHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Key   string `json:"key"`
        Value string `json:"value"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    cs.cache.Set(req.Key, req.Value, 1*time.Hour)
    w.WriteHeader(http.StatusOK)
}
func (cs *CacheServer) GetHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    value, found := cs.cache.Get(key)
    if !found { http.NotFound(w, r); return }
    json.NewEncoder(w).Encode(map[string]string{"value": value})
}
```

### 4. Eviction siyasətləri
| Siyasət | Nə edir | Güclü | Zəif |
|---|---|---|---|
| **LRU** | ən az istifadə olunanı atır | access pattern-a uyğun, hit-rate yaxşı | sıra izləmə overhead-i |
| **TTL** | vaxtı bitəni atır | data təzə qalır | istifadəyə baxmır; periodik yoxlama resursu |
| **FIFO** | ən köhnə əlavəni atır | ən sadə | tez istifadə olunanı da ata bilər |

**TTL implementasiyası — 2 yol:**
1. **Goroutine + Ticker:** proaktiv təmizlik (Get daha təmiz), amma əlavə goroutine overhead-i.
2. **Get-də yoxlama:** sadə, on-demand; amma expiry item-lar access olmayana qədər qalır.

Layihə hər ikisini birləşdirir: ticker-lü background təmizlik + Get-də expiry yoxlaması.

```go
func (c *Cache) startEvictionTicker(d time.Duration) {
    ticker := time.NewTicker(d)
    go func() {
        for range ticker.C {
            c.evictExpiredItems()
        }
    }()
}
```

### 5. LRU — container/list (doubly-linked list)
```go
type Cache struct {
    mu       sync.RWMutex
    items    map[string]*list.Element // key → list element
    eviction *list.List               // istifadə sırası (front = yeni)
    capacity int
}
type entry struct {
    key   string
    value CacheItem
}

func (c *Cache) Set(key, value string, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if elem, found := c.items[key]; found {          // köhnə dəyəri təmizlə
        c.eviction.Remove(elem)
        delete(c.items, key)
    }
    if c.eviction.Len() >= c.capacity {              // doludursa LRU at
        c.evictLRU()
    }
    elem := c.eviction.PushFront(&entry{key, CacheItem{value, time.Now().Add(ttl)}})
    c.items[key] = elem
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.Lock()                                      // MoveToFront yazmadır → tam Lock
    defer c.mu.Unlock()
    elem, found := c.items[key]
    if !found || time.Now().After(elem.Value.(*entry).value.ExpiryTime) {
        if found {
            c.eviction.Remove(elem)
            delete(c.items, key)
        }
        return "", false
    }
    c.eviction.MoveToFront(elem)                     // istifadə olundu → öndə
    return elem.Value.(*entry).value.Value, true
}

func (c *Cache) evictLRU() {
    elem := c.eviction.Back()                        // ən son istifadə olunan
    if elem != nil {
        c.eviction.Remove(elem)
        kv := elem.Value.(*entry)
        delete(c.items, kv.key)
    }
}
```
**Sub-kod izahı:** map O(1) axtarış, list istifadə sırası; hər Get access-ı PushFront-ə daşıyır; Back() = LRU qurbanı. PushFront/MoveToFront yazma əməliyyatı olduğundan Get-də RLock deyil, tam Lock istifadə olunur.

### 6. Replikasiya — P2P
4 strategiya: Primary-Replica (master bottleneck + SPOF), Pub/Sub (broker SPOF), Raft/Paxos consensus (güclü tutarlılıq, kompleks), **P2P** (bütün node-lar bərabər — SPOF yoxdur, horizontal scaling təbii; eventual consistency + konflikt həlli öz üzərinə).

**Replikasiya funksiyası:**
```go
const replicationHeader = "X-Replication-Request"

func (cs *CacheServer) replicateSet(key, value string) {
    cs.mu.Lock()
    defer cs.mu.Unlock()
    req := struct {
        Key   string `json:"key"`
        Value string `json:"value"`
    }{Key: key, Value: value}
    data, _ := json.Marshal(req)
    for _, peer := range cs.peers {
        if peer != cs.selfID {
            go func(peer string) {                    // hər peer paralel
                client := &http.Client{}
                req, err := http.NewRequest("POST", peer+"/set", bytes.NewReader(data))
                if err != nil { log.Printf("Failed to create replication request: %v", err); return }
                req.Header.Set("Content-Type", "application/json")
                req.Header.Set(replicationHeader, "true")   // replikasiya işarəsi
                _, err = client.Do(req)
                if err != nil { log.Printf("Failed to replicate to peer %s: %v", peer, err) }
            }(peer)
        }
    }
}
```
**Sonsuz dövrün qarşısı:** SetHandler-də `if r.Header.Get(replicationHeader) == "" { go cs.replicateSet(...) }` — replikasiya request-i gələndə header doludur → yeni replikasiya başlamır (infinite loop block).

**Test:**
```bash
go run main.go -port=:8080 -peers=http://localhost:8081
go run main.go -port=:8081 -peers=http://localhost:8080
curl -X POST -d '{"key":"foo","value":"bar"}' -H "Content-Type: application/json" http://localhost:8080/set
curl http://localhost:8081/get?key=foo     # → {"value":"bar"} replikasiya işləyir
```

### 7. Sharding — 3 yanaşma
| | Range | Hash | Consistent Hashing |
|---|---|---|---|
| Bölgü | ardıcıl aralıqlar | hash(key) % n | halqa üzərində node/key hash-i |
| Bərabərlik | skew hot-spot riski | bərabər | yaxşı balans |
| Re-sharding | — | çox data hərəkəti | MINIMAL hərəkət (yalnız qonşu aralıq) |
| Range query | effektiv | zəif | zəif |
| Komplekslik | sadə | sadə | daha yüksək |

**Consistent hashing (hash ring):** node-lar və key-lər eyni halqaya (hash məkanı) düşür; key sahibi = saat əqrəbi istiqamətində **ilk node**. Node çıxanda yalnız onun aralığındakı key-lər qonşuya keçir.

**HashRing strukturu:**
```go
type Node struct {
    ID   string
    Addr string
}
type HashRing struct {
    nodes  []Node
    hashes []uint32       // sorted — binary search üçün
    lock   sync.RWMutex
}
// hash()          — SHA-1 ilə node ID / key hash
// AddNode()       — hash hesabla, slices-ə daxil et, sırala
// GetNode(key)    — binary search: key hash-indən ≥ ilk node = sahib
```

**Sharding + forwardinq SetHandler:**
```go
func (cs *CacheServer) SetHandler(w http.ResponseWriter, r *http.Request) {
    // ... JSON decode ...
    targetNode := cs.hashRing.GetNode(req.Key)
    if targetNode.Addr == "self" {
        cs.cache.Set(req.Key, req.Value, 1*time.Hour)
        if r.Header.Get(replicationHeader) == "" {
            go cs.replicateSet(req.Key, req.Value)
        }
        w.WriteHeader(http.StatusOK)
    } else {
        cs.forwardRequest(w, targetNode, r)      // düyün sahibinə yönləndir
    }
}
```
**GetHandler-də loop qorunması:** `X-Forwarded-For` header-i öz selfID-dirsə → "Loop detected" (400) — forward zənciri özünə qayıdır.

**forwardRequest:** GET → query string ilə yeni URL; POST → body oxu, təkrar qur, header-ləri köçür, `client.Do`, `resp` statusu + body `io.Copy` ilə cavaba köçür.

```bash
# Test — key-lər hash-ə görə node-lara paylanır:
curl -X POST -H "Content-Type: application/json" -d '{"key":"key1","value":"value1"}' localhost:8080/set
curl -X POST -H "Content-Type: application/json" -d '{"key":"key2","value":"value2"}' localhost:8080/set
curl localhost:8080/get?key=key1     # 8080 VƏ YA 8083 — hash-dən asılı
curl localhost:8083/get?key=key2     # istənilən node cavab verir (forward edir)
```

### 8. İrəliləyüş yolları
- **Alqoritmlər:** LIRS, ARC (LRU-dan adaptiv), TTL/LRU parametr tuning, sıxışdırma (compression), connection pooling.
- **Metrikalar:** hit/miss nisbəti, eviction rate, latency, throughput, memory — Grafana dashboard + threshold alert-lər.
- **Profiling:** CPU (isti funksiyalar) + memory (leak/alloc) — Chapter 9 alətləri.

## Əsas terminlər
- Distributed cache (paylanmış keş) — app ↔ data store arası yaddaş qatı
- CAP theorem — consistency / availability / partition tolerance balansı
- sync.RWMutex — paralel oxu + eksklüziv yaz
- sync.Map / ristretto / golang-lru / go-cache — hazır həllər
- TTL (Time To Live) — vaxt_bitmə əsaslı eviction
- LRU (Least Recently Used) + FIFO — eviction siyasətləri
- container/list — doubly-linked list (PushFront/MoveToFront/Back)
- time.NewTicker — periodik background təmizlik
- P2P replikasiya vs Primary-Replica / Pub-Sub / Raft-Paxos consensus
- X-Replication-Request — replikasiya dövrünü kəsən header
- Sharding: range-based / hash-based / consistent hashing
- Hash ring — dairəvi hash məkanı; saat əqrəbi ilk node = sahib
- HashRing: SHA-1, sorted hashes, binary search
- X-Forwarded-For — forward loop aşkarlanması
- Hit/miss ratio, eviction rate — keş metrikaları
- LIRS / ARC — inkişaf etmiş replacement alqoritmləri

## Praktik nəticə
1. Oxuma-agır keş üçün RWMutex; daha irəli: sharded locks və ya hazır ristretto/golang-lru.
2. Replikasiya request-lərini xüsusi header ilə işarələ — əks halda node-lar sonsuz qarşılıqlı replikasiya dövrünə düşür.
3. LRU üçün map + doubly-linked list cütlüyü: map O(1) tapır, list sıranı saxlayır; Get yazma sayıldığından (MoveToFront) tam Lock.
4. Consistent hashing re-sharding-da yalnız minimal datanı hərəkət etdirir — node əlavə/çıxarmada sistem davam edir.
5. Forward zəncirində loop qoruması (X-Forwarded-For yoxlaması) mütləqdir — əks halda node-lar bir-birinə sonsuz request göndərir.
6. Keşi ölç: hit ratio, latency, eviction rate — observability olmadan optimizasiya kor təsadüfdür.

## Mənbə
Pages: 267-303 (PDF səh. 288-325)
