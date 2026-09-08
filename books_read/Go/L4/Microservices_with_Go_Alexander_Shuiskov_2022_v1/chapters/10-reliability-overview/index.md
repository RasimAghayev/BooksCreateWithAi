# Chapter 10 — Reliability Overview (səh. 193-214)

## Bu fəsil nədən bəhs edir?

Reliability anlayışı, avtomatlaşdırma texnikaları (retry, backoff, timeout,
fallback, rate limiting, graceful shutdown) və proses/mədəniyyət (on-call,
incident management, drills).

## Əsas fikirlər

### 1. Reliability (Etibarlılıq) nədir?
- Gözlənilən şəkildə işləmə + **limitlərin AÇIQ təyini**
- Kod izolyasiyada işləyir, amma sistemdə: **DoS** (həddindən artıq request →
  CPU/memory limiti çatır), backward-incompatible API dəyişikliyi zəng edənləri sındırır
- Eksplisit ol: throughput (max RPS), **congestion policy** (100 paralel
  requestdən çoxu → HTTP 429 Too Many Requests)
- 3 kateqoriya: **Prevention / Detection / Mitigation** (qarşısını alma /
  aşkarlama / yumşaltma) — avtomatlaşdırma və proseslərlə

### 2. Retriable vs Non-retriable xətalar
| Non-retriable | Retriable |
|---|---|
| InvalidArgument (validasiya) | **DeadlineExceeded** (timeout) |
| NotFound | **ResourceExhausted** (quota/resurs) |
| | **Unavailable** (müvəqqəti offline) |

```go
func shouldRetry(err error) bool {
    e, ok := status.FromError(err)
    if !ok { return false }
    return e.Code() == codes.DeadlineExceeded ||
        e.Code() == codes.ResourceExhausted || e.Code() == codes.Unavailable
}
// Get: maxRetries=5 dövrü, shouldRetry → continue, yoxsa return err
```

### 3. Backoff və jittering
- Dərhal retry = request BURST-ləri (server conjestionda boğulur)
- **Backoff:** retry-lər arası gecikmə; constant vs **exponential** (100ms →
  400ms → 900ms — serverə bərpa şansı verir; cenkalti/backoff kitabxanası)
- **Jittering:** gecikmə ±10% random — bütün clientlər eyni anda retry
  etməsin, yük bərabər paylansın

### 4. Timeouts və deadlines
- Request gözləməsi sonsuz olmasın: timeout = müddət (10s), deadline = an
  (2074-01-01T00:00:00Z) — texniki olaraq eyni məqsəd
```go
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)   // və ya
ctx, cancel := context.WithDeadline(ctx, deadline)
defer cancel()
```

### 5. Fallback (əvəzedici məntiq)
Retry-lər də uğursuz → 3 seçim: error qaytar / panic / **fallback**
- Rating DB əlçatmaz → keşdəki (map) rating-ləri qaytar
- **Graceful degradation** (Zərif Enmə) — məhdud amma işlək funksionallıq
- Fallback qəsdən olduğunu koddə göstər + log/metrics ötür

### 6. Rate limiting (token bucket)
Səviyyələr: client / server / network (load balancer). Client/server instans
başınadır — 1000 instans × 100 RPS = 100,000! Global limit üçün balancer lazımdır.

**Token bucket:** bucket ölçüsü b, dolum sürəti r/s; hər request 1 token alır:
```go
limiter := rate.NewLimiter(rate.Limit(limit), burst)  // golang.org/x/time/rate
// limiter.Allow() → icazə/dürüstlük
```
gRPC-ə qoşulması — **UnaryInterceptor** (movie server):
```go
srv := grpc.NewServer(grpc.UnaryInterceptor(
    ratelimit.UnaryServerInterceptor(newLimiter(100, 100))))
```
Limit aşılırsa `codes.ResourceExhausted`. Limit çox aşağı = istifadəçi
senzurasiyası — benchmark ilə ədalətli sərhəd tap.

### 7. Graceful shutdown (zəif dayandırma)
Səbəblər: SIGINT (Ctrl+C), SIGTERM/SIGKILL, panic. Risklər: düşmüş
requestlər, **connection leak** (DB bağlantısı boşalmır).

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
ctx, cancel := context.WithCancel(context.Background())
go func() { s := <-sigChan; cancel(); srv.GracefulStop() }()  // Stop YOX
wg.Wait()
```
- `GracefulStop` vs `Stop`: işlənməkdə olan requestləri gözləyir
- Panic üçün `defer recover()`; komponentlərə xəbər = **context cancellation**

### 8. On-call prosesi
- Rotation (növbə, ~1 həftə shift), **escalation policy** (primary →
  secondary → menecer → ... CTO), **shadow** rol (təcrübəsizlər üçün), ACK
  məcburi
- PagerDuty — populyar platforma (SMS/email/call, Slack/Zoom inteqrasiyası)
- Çətinliklər: rotation ownership xəritəsi (servis → növbə), cross-service
  hadisələr → mərkəzi komanda (Uber-də **Ring0**)

### 9. Incident management
- **Mitigation** (həll) + **Prevention** (təkrar olmasın)
- **Runbook** — qısa, actionable addımlar; `rating_service_fd_limit_reached:
  mitigation: Restart the service`; log/dashboard linkləri
- Hədəf: minimal **TTR (time to repair)**
- **Postmortem** sənədi: başlıq, müəllif, detection/mitigation, kontekst,
  root cause, təsir, timeline, dərslər, action items (Google SRE nümunəsi)
- **Five whys** — "niyə?"-ni kök səbəb tapana qədər soruş (DB offline ← yüksək
  yük ← movie servicedəki bug)

### 10. Reliability drills (məşqlər)
- Backup var ≠ restore edə bilirsən — PERİODİK sınaq lazımdır
- Növlər: DB backup/restore, network failure simulyasiyası
- Faydalar: gözlənilməz xətalar/panic-lər kontrollü şəraitdə aşkarlanır;
  **transitive/circular asılılıqlar** üzə çıxır; gələcək mitigation sürətlənir
- Drill = planlaşdırılmış incident → postmortem + runbook yenilənməsi

## Termindirmə (AZ)
- Reliability — Etibarlılıq
- DoS (Denial of Service) — Xidmətdən Məhrumetmə
- Backoff — Gecikməli Yenidən Cəhd
- Jittering — Təsadüfi Dəyişkənlik
- Fallback — Əvəzedici Məntiq
- Graceful Degradation — Zərif Enmə
- Rate Limiting — Sorğu Limitasiyası
- Token Bucket — Jeton Vedrəsi alqoritmi
- Graceful Shutdown — Zəif Dayandırma
- Runbook — Əməliyyat Kitabçası
- Postmortem — Hadisədən Sonra Hesabat
- TTR (Time to Repair) — Təmir Vaxtı

## Kviz sualları
1. Retry üçün hansı gRPC kodları uyğundur? (DeadlineExceeded,
   ResourceExhausted, Unavailable)
2. Jittering nə üçün lazımdır? (Bütün clientlər sinxron retry edib burst
   yaratmasınlar)
3. GracefulStop Stop-dan fərqi? (Cari requestlərin işlənməsini gözləyir)
4. Five whys nə edir? (Simptomdan kök səbəbə zəncirvari "niyə?" sorğusu)
