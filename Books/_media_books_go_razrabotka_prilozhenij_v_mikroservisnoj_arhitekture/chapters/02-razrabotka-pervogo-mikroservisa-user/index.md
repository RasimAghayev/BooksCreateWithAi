# Chapter 1 — Разработка первого микросервиса (User)

Pages: 13-86

## Bu chapter nədən bəhs edir?
Bu chapter praktiki olaraq ilk mikroservisi (User) necə quracağımızı
addım-addım izah edir. Lokal mühitin qurulmasından, Go quraşdırılmasından,
proyekt strukturu yaradılmasından, PostgreSQL bazasının layihələndirilməsindən,
biznes məntiqə hazırlanmasından və testləmədən bəhs edir.

## Əsas fikirlər

### 1. Lokal mühit qurulması
- **Code Editor (kod redaktoru):** Visual Studio Code tövsiyə olunur. Go
  Team Google tərəfindən IntelliSense plaqini tövsiyə olunur.
- **Go quraşdırma:** https://go.dev/dl/ saytından endirilir. `go version`
  komandası ilə yoxlanılır.
- **Package Manager (paket meneceri):** Go Modules (Go 1.11+), go.mod və
  go.sum faylları ilə asılılıqların idarə edilməsi.
- **GVM (Go Version Manager):** Bir neçə Go versiyası arasında keçid üçün
  istifadə olunur. Versiyanı yüksəltmək poetik (addım-addım) olmalıdır.
- **Protobuf:** Interface təyin etmək üçün istifadə olunur.
- **Git:** Versiya idarəsi sistemi.
- **Docker & Docker Compose:** Konteynerləşdirmə və çoxkonteynerli tətbiqlər
  üçün alətlər.

### 2. Go Modules və asılılıqların idarə edilməsi
- **$GOPATH** əvəz edilib Go Modules gəlmişdir.
- `go.mod` faylı: module adı, Go versiyası, asılılıqlar.
- `go.sum` faylı: asılılıqların hash cəmləri, təhlükəsizlik üçün.
- **Package Manager tarixi:**
  - godep → glide → govendor → Go Modules (standart)
  - Go 1.16+ ilə Go Modules standart oldu.

### 3. Proyekt strukturu
- Təmiz arxitektura (Clean Architecture) prinsipləri
- Dependency Injection (asılılıqların inyeksiyası)
- Environment variable-lərdən istifadə
- Logging (logrus və s.)

### 4. Verilənlər bazası
- ORM (GORM) istifadəsi
- PostgreSQL layihələndirilməsi
- Normal forms (1NF, 2NF, 3NF)
- Swagger dokumentasiyası

### 5. Biznes məntiqə və routing
- REST API dizaynı
- HTTP handler-lər
- Validasiya
- Testləmə (unit, integration)

## Əsas terminlər
- Microservice (mikroservis)
- REST API
- ORM (Object-Relational Mapping / obyekt-əlaqəli xəritələmə)
- Docker (konteynerləşdirmə platforması)
- Docker Compose (çoxkonteynerli tətbiqlər üçün alətlər)
- Swagger (API dokumentasiyası)
- PostgreSQL ( obyekt-relational verilənlər bazası)
- GORM (Go ORM kitabxanası)
- Protobuf (Protocol Buffers, serializasiya formatı)
- Git (versiya idarəsi sistemi)
- Go Modules (Go asılılıqların idarə edilməsi sistemi)
- go.mod (modul tərifi faylı)
- go.sum (hash cəmləri faylı)
- $GOPATH (köhnə Go iş qovluğu dəyişəni)
- GVM (Go Version Manager)
- IntelliSense (kod avtomatik tamamlama)
- Editor (kod redaktoru)

## Praktik nəticə
Bu chapter-dən sonra siz tam funksionallığa malik bir User mikroservisinə
sahibsınız. Bu servis:
- Qeydiyyat və giriş funksiyalarına malikdir
- PostgreSQL bazası ilə işləyir
- Docker konteynerində işləyir
- REST API vasitəsilə interfeys təqdim edir

## Mənbə
Pages: 13-86
