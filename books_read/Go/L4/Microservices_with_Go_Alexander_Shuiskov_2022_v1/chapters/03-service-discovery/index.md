# Chapter 3 — Service Discovery (səh. 55-76)

## Bu fəsil nədən bəhs edir?

Statik unvanların problemini, service registry anlayışını, client-side vs
server-side discovery modellərini, health monitoring-i və Consul əsaslı
implementasiyanı əhatə edir. Movie aplikasiyasına dinamik discovery əlavə olunur.

## Əsas fikirlər

### 1. Problem: statik unvanlar
- Əvvəlki fəsildə unvanlar (localhost:8081/8082/8083) kodda sərt yazılmışdı
- Çoxlu instans olduqda: hansı unvanı çağırmalıq? İnstans söndükdə nə etməli?
- Siyahı konfiqurasiyası ilə idarə = hər əlavyə/silmə üçün bütün zəng edən
  servislərin konfiqini dəyişmək lazımdır — çevik deyil

### 2. Service Registry (Servis Reyestri)
| Əməliyyat | Təsvir |
|---|---|
| Register | instans qeydiyyata alınır (başlanğıcda) |
| Deregister | instans silinir (dayanma vaxtı) |
| ServiceAddresses | aktiv instansların unvan siyahısı |

- Registry sadəcə `service name → address list` xəritəsidir
- **Health monitoring:** Pull modeli (registry özü yoxlayır) vs Push modeli
  (servis öz vəziyyətini yeniləyir)

### 3. İki model
| | Client-side | Server-side |
|---|---|---|
| Kim registry-yə müraciət edir | tətbiqin özü | **load balancer** |
| Load balancing | tətbiqdə (kod mürəkkəbləşir) | balancerdə |
| Plus | sadə, birbaşa | tətbiq registry haqqında bilimir |
| Minus | LB məntiqi tətbiqdə | balancer qurmaq lazımdır |

### 4. Həllər
- **HashiCorp Consul:** Go-da yazılıb; PUT /catalog/register,
  PUT /catalog/deregister, GET /catalog/services; API və ya DNS
- **Kubernetes:** öz service API + endpoints; load balancer qoşmaqla
  server-side da mümkün (Ch8-də ətraflı)

### 5. Texnologiya-neytral interfeys
```go
type Registry interface {
    Register(ctx, instanceID, serviceName, hostPort) error
    Deregister(ctx, instanceID, serviceName) error
    ServiceAddresses(ctx, serviceID) ([]string, error)
    ReportHealthyState(instanceID, serviceName) error  // push health
}
```
- Consul tiplərini (consul.ServiceEntry) kod boyu yaymaq YOX — sonradan dəyişmək
  çətinləşir; `[]string` qaytar (technology-agnostic)
- `GenerateInstanceID(serviceName)` → `name-random` formatında ID

### 6. İki implementasiya
**In-memory (test üçün):**
```go
map[serviceName]map[instanceID]*serviceInstance  // hostPort + lastActive
```
- sync.RWMutex ilə qorunur; ServiceAddresses yalnız son 5 san.-də sağlam
  (lastActive) instansları qaytarır

**Consul əsaslı (production):**
- Register → `Agent().ServiceRegister` + TTL check `"5s"`
- ServiceAddresses → `Health().Service(name, "", true, nil)` — yalnız
  passing-health instanslar
- ReportHealthyState → `Agent().PassTTL` (hər 1 san.-də goroutine-da)
- Deregister → `Agent().ServiceDeregister`

### 7. Movie servisinə inteqrasiya
- Gateway-lər artıq `addr string` deyil, `registry discovery.Registry` alır
- Hər çağırışda: registry-dən unvanlar → **rand.Intn ilə təsadüfi seçim**
  (sadə load balancing)
- main(): flag ilə port → consul registry-yə qeydiyyat → sağlamlıq goroutine
  (1 san. interval) → `defer Deregister` → handler başladır
- Consul-u işə salma: `docker run -d -p 8500:8500 consul agent -server ...`
- UI: http://localhost:8500 — Passing/Critical vəziyyətləri görünür
- Yeni instans: `go run *.go --port 8084` — kod dəyişmədən scale

## Termindirmə (AZ)
- Service Discovery — Servis Kəşfiyyəti (servisin dinamik tapılması)
- Registry — Reyestr (servis instanslarının qeydiyyatı)
- Load Balancer — Yük Bölücü
- Health Check — Sağlamlıq Yoxlaması
- Deregister — Qeydiyyatdan Çıxarma
- TTL (Time To Live) — Yaşama Müddəti

## Kviz sualları
1. Client-side discovery-nin əsas çatışmazlığı nədir? (LB məntiqi tətbiq koduna
   düşür — mürəkkəblik)
2. In-memory registry "aktiv" instansı necə müəyyənləşdirir? (lastActive son 5
   san.-dən gec deyilsə)
3. Consul-da ReportHealthyState nə edir? (Agent().PassTTL — TTL check-i yeniləyir)
4. Niyə Registry interfeysi `[]string` qaytarır, consul tipləri yox? (Texnologiya
   dəyişəndə bütün çağırış yerlərini yeniləmək məcburiyyəti olmasın)
