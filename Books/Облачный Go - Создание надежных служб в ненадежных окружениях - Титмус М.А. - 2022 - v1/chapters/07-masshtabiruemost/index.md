# Chapter 7 — Масштабируемость

## Bu chapter nədən bəhs edir?

Sistemin artan yüke (load) uyğunlaşaraq performansını saxlayabilməsi — miqyaslana bilənliyi (scalability) izah edir. Stateful (vəziyyətli) və stateless (vəziyyətsiz) dizaynlar arasındakı fərq, yaddaş sızıntıları (memory leaks), gorutin sızıntıları, polling nümunələri və xidmət arxitekturası (monolith vs microservices) müzakirə edilir.

## Əsas fikirlər

### 1. Vəziyyətli (Stateful) və Vəziyyətsiz (Stateless) dizayn
**Nədir:** Tətbiq vəziyyəti (application state) və resurs vəziyyəti (resource state) arasındakı fərq.

**Necə işləyir:**
- **Tətbiq vəziyyəti (Application state):** Müəyyən bir server instance-na bağlı olan sorğu konteksti, sessiya məlumatları, yaddaşda saxlanan müvəqqəti data.
- **Resurs vəziyyəti (Resource state):** Verilənlər bazası, cache (Redis, Memcached), mesaj növü — bu resurslar paylaşıla bilər və təhlükəsizdir.

**Statelesslik üstünlükləri:**
- Hər hansı bir instance hər hansı bir sorğunu emal edə bilər
- Instance-lər yox edilərkən və ya əlavə edilərkən data itkisi yoxdur
- API daha sadə və proqnozlaşdırıla bilən olur
- Yalnızca resurs vəziyyətini (DB, cache) idarə etmək lazımdır

### 2. Yaddaş sızıntıları (Memory Leaks)
**Nədir:** İstifadə edilməyən və ya artıq olmayan məlumatların yaddaşda qalması, bu da vaxtla bərabər yaddaşın bitməsinə səbəb olur.

**Necə işləyir:**
- **Sonsuz böyüyən məlumat strukturları:** `map[string]int` kimi struktur-a heç vaxt aradan qaldırılmayan elementlər əlavə edilir
- **Gorutin sızıntıları:** Başlanmış gorutinlərdən bəziləri heç vaxt bitməyə bilər — məsələn, kanal (channel) gözləməkdə qalan gorutin
- **Unudulmuş `time.Ticker`-lar:** `ticker.Stop()` çağırılmadan gorutin bitirilərsə, ticker axını (ticker goroutine) işləməyə davam edir
- **Həll yolları:** LRU (Least Recently Used) cache, kanalların bağlanması (close), `ticker.Stop()` unutmayın

### 3. Polling nümunələri
**Nədir:** Resursların vəziyyətini müəyyən aralıqlarla yoxlamaq üçün istifadə edilən texnika.

**Necə işləyir:**
- **Nüvəvi (Threaded) yanaşma:** Hər resurs üçün ayrı OS thread, mütləq məxaric (mutex) ilə ən yeni yoxlanmış resursu tapmaq
- **Gorutin əsaslı yanaşma:** Hər resurs üçün gorutin, kanal (channel) vasitəsilə vəziyyətləri göndərmək
- Go gorutinləri çox yüngüldür (~2 KB stack-dən başlayaraq dinamik böyüyür), ona görə də onlardan etibarlı şəkildə istifadə edilə bilər

**Kitabdan kod nümunəsi (gorutin əsaslı polling):**
```go
func Poller(stop <-chan bool, in <-chan *Resource, out chan<- *Resource) {
    for {
        select {
        case r := <-in:
            go func(res *Resource) {
                state := getState(res)
                select {
                case out <- res:
                case <-stop:
                    return
                }
            }(r)
        case <-stop:
            close(out)
            return
        }
    }
}
```

### 4. Monolith vəya Mikroservislər
**Nədir:** Xidmət arxitekturasının iki əsas modeli.

**Necə işləyir:**
- **Monolith:** Bütün funksionallıq tək prosesdə, tək kod bazasında. Sadə başlanğıc, lakin miqyaslandıqda çətinlik yaradır.
- **Mikroservislər:** Kiçik, müstəqil deploy edilə bilən xidmətlər, yüngül mexanizmlər (HTTP, gRPC) ilə əlaqələnir. Hər xidmət öz bazasını idarə edir, müstəqil miqyaslana bilir.

## Əsas terminlər
- Scalability (miqyaslana bilənlik) — artan yükə uyğunlaşma
- Stateful (vəziyyətli) — serverə bağlı vəziyyət saxlayan
- Stateless (vəziyyətsiz) — serverə bağlı vəziyyət saxlamayan
- Memory leak (yaddaş sızıntısı) — istifadə edilməyən yaddaşın qalması
- Goroutine leak (gorutin sızıntısı) — bitməyən gorutin
- time.Ticker — periyodik taymer, `Stop()` lazımdır
- LRU cache (ən yeni istifadə edilən cache) — yaddaş məhdudiyyəti ilə
- Monolith (monolit) — tək prosesdəki bütün funksionallıq
- Microservices (mikroservislər) — müstəqil kiçik xidmətlər

## Praktik nəticə
Miqyaslana bilənlik əsasən vəziyyətsiz (stateless) dizayn və resurs vəziyyətini (DB, cache) xarici saxlamaqla təmin edilir. Yaddaş sızıntıları (gorutin, ticker, map-lər) diqqətsizliklə yaranır və vaxtla kritik problem halına gəlir. Go gorutinləri və kanalları polling və konkurensiya üçün güclü, amma `Stop()` və kanal bağlama qaydaları unudulmamalıdır.

## Mənbə
Pages: 206-225 (PDF səh. 206-225)
