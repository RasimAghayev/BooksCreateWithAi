# Chapter 10 — Управляемость

## Bu chapter nədən bəhs edir?

Sistemin davranışını kod dəyişdirmədən dəyişdirmə imkanı — idarəetmə (manageability) — izah edir. Konfiqurasiya idarəetməsi (environment variables, CLI flags, config faylları, dinamik yeniləmə), xüsusiyyət flaqları (feature flags), CLI dizaynı (Cobra) və birlikdə konfiqurasiya (Viper) müzakirə edilir.

## Əsas fikirlər

### 1. İdarəetmə (Manageability) anlayışı
**Nədir:** Sistemin davranışını yenidən kod yazmaq və ya yenidən deploy etmək tələb etmədən dəyişdirmə imkanı.

**Necə işləyir:**
- URL-lər, portlar, timeout-lar kimi dəyərlər koddakı sabitlər (hardcoded) olmamalıdır
- Konfiqurasiya xarici mənbələrdə (environment variable, config faylı, CLI flag) saxlanmalıdır
- Dəyişikliklər yalnız konfiqurasiyanı dəyişdirməklə tətbiq edilə bilməlidir

### 2. Konfiqurasiya iyerarxiyası (Hierarchy)
**Nədir:** Konfiqurasiya dəyərlərinin üstünlük sırası — ən yüksək prioritetli dəyər ən aşağıdakıni ləğv edir.

**Necə işləyir (yüksəkdən aşağıya):**
1. CLI flags — ən yüksək prioritet
2. Environment variables
3. Config faylları (JSON, YAML, TOML)
4. Default dəyərlər — ən aşağı prioritet

**Kitabdan kod nümunəsi (Viper ilə konfiqurasiya):**
```go
import "github.com/spf13/viper"

viper.SetDefault("id", "13")
viper.SetConfigName("config")
viper.AddConfigPath("/etc/myapp/")
viper.BindEnv("id")
```

### 3. Dinamik konfiqurasiya yeniləməsi (Dynamic Reloading)
**Nədir:** Konfiqurasiya faylı dəyişdikdə xidməti yenidən başlatmaq tələbi olmadan dəyişikliklərin tətbiq edilməsi.

**Necə işləyir:**
- `viper.WatchConfig()` — konfiqurasiya faylını izləyir
- Dəyişiklik aşkar edildikdə `viper.OnConfigChange()` callback-i çağırılır
- Bəzi hallarda yenidən başlatmaq (restart) tələb olunur — məsələn, log səviyyəsi dəyişəndə

### 4. Xüsusiyyət flaqları (Feature Flags / Feature Toggles)
**Nədir:** Funksionallığı deploy etmədən açıb-bağlamaq imkanı verən mexanizm.

**Necə işləyir:**
- **Statik flaqlar:** Konfiqurasiya faylında saxlanır, hər sorğu üçün eyni dəyər qaytarır
- **Dinamik flaqlar:** Hər sorğu üçün fərqli dəyər qaytara bilər — məsələn, müəyyən IP diapazonundan gələn client-lər üçün yeni funksiya aktiv, digərləri üçün deyil

**Kitabdan kod nümunəsi (dinamik flaq):**
```go
type Enabled func(flag string, r *http.Request) (bool, error)

func fromPrivateIP(flag string, r *http.Request) (bool, error) {
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    userIP := net.ParseIP(ip)
    _, cidr, _ := net.ParseCIDR("10.0.0.0/8")
    return cidr.Contains(userIP), nil
}
```

### 5. Cobra və Viper birlikdə istifadəsi
**Nədir:** CLI və konfiqurasiya idarəetməsi üçün ən çox istifadə olunan Go kitabxanaları.

**Necə işləyir:**
- **Cobra:** CLI framework — komandalar (commands), alt-komandalar (subcommands), avtomatik kömək (help) yaradır
- **Viper:** Konfiqurasiya idarəetməsi — JSON, YAML, TOML, env vars, CLI flags, uzaq konfiqurasiya (remote config) dəstəyi
- Birlikdə: Cobra CLI flag-ləri Viper-ə bağlayır, Viper konfiqurasiya iyerarxiyasını idarə edir

## Əsas terminlər
- Manageability (idarəetmə) — kod dəyişdirmədən davranışı dəyişdirmə
- Configuration hierarchy (konfiqurasiya iyerarxiyası) — defaults → config → env → flags
- Feature flags (xüsusiyyət flaqları) — deploy etmədən funksionallığı aç/bağla
- Dynamic reloading (dinamik yeniləmə) — yenidən başlatmaq olmadan konfiqurasiya yeniləmə
- Cobra — CLI framework (spf13/cobra)
- Viper — konfiqurasiya kitabxanası (spf13/viper)
- YAML struct tags — `yaml:"name"`, `yaml:",inline"`, `yaml:"flow"`
- RFC1918 — şəxsi IP diapazonları (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)

## Praktik nəticə
İdarəetmə (manageability) bulud sistemlərində ən vacib xüsusiyyətlərdən biridir — "operate" mərhələsini asanlaşdırır. Konfiqurasiya iyerarxiyası sayəsində fərqli mühitlərdə (dev, staging, prod) eyni kodu istifadə edə bilərik. Xüsusiyyət flaqları (feature flags) isə riski minimuma endirərək yeni funksionallığı test və iskəslənməyə buraxa bilməyimizə imkan verir.

## Mənbə
Pages: 315-353 (PDF səh. 315-353)
