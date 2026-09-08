# Effective Go Recipes — Müəllim qeydləri (AZ)

## Ümumi qiymətləndirmə

Miki Tebeka (353 Solutions, Ardan Labs dərsçisi) Pragmatic seriyasında
Go-nun gündəlik problemlərini 84 reseptdə topluyub. Hər resept real iş
ssenarisi (task email/issue formasında) ilə açılır — bu, kitabı adi
" cookbook"-dan fərqləndirir: oxucu HANSI problem olduqda HANSI reseptə
getdiyini bilir.

## Güclü tərəflər

1. **Realistik problem çərçivəsi** — hər resept komanda iqtişafı ilə
   başlayır ("Jessie belə yazır...") — tələbə özünü real layihədə görür.
2. **Benchmark/profiling mədəniyyəti** — cumSum (ch6), Mutex vs Atomic
   (ch10), time.Now (ch10): hər optimizasiya iddiası rəqəmlə sübut
   olunur.
3. **Trade-off şəffaflığı** — mock="aldatma" (ch13), cgo="Go deyil"
   (ch12), statik build=DNS riski (ch14): kitab həllərin qiymətini gizlətmir.
4. **Müasirlik** — generics (Go 1.18+), embed (1.16+), fuzzing (1.18+),
   errors.Join (1.20) — köhnəlməmiş material.

## Zəif tərəflər və müəllim tövsiyələri

1. **Structural bütövlük** — fəsillər müstəqildir; tələbə üçün təkrar
   oxunuş planı verin: ch1-2 → ch4 → ch6-8 → ch9-10 → ch13 → ch14-15.
2. **Zerolog/Go Recipes bənzərliyi** — Build Systems with Go (Tirado)
   ilə üst-üstə düşən mövzular var (interfeyslər, konkurensiya, test);
   bu kitabı nüans-kitabı kimi təqdim edin.
3. **Recipe 30 (go:linkname)** — vacib incəlik: müəllif özü "dirty
   trick" deyir; tələbələrə bunun "son çarə + issue linki" olduğunu
   ikiqat vurğulayın.
4. **Recipe 60 (NTP)** — şərh: real sistemdə ntpd/chrony istifadə edin;
   resept UDP binary protokol dərsi kimi dəyərlidir.
5. **Viper/Cobra** — kitab ardanlabs/conf seçib; daha geniş ekosistem
   üçün Viper-in rolunu əlavə edin.

## Sərt xatırlatmalar (tələbələr üçün)

- Recipe 31: map-dən oxuyarkən həmişə comma-ok — 0.0 "endirim" bug-ı
  klassik müsahibə sualıdır.
- Recipe 33: append-ə güvənmə — böyük slice-lərdə make(0, n) ilə başlayın.
- Recipe 49: "Never start a goroutine without knowing how it will stop"
  — Dave Cheney sitatı əzbərlənməlidir.
- Recipe 55: `-race` yalnız development-də yox, CI-də də standartdır.
- Recipe 82: bahalı parametrlər Check arxasında — production
  performansının sessiz öldürücüsü.

## Kitabın ən dəyərli hissəsi

Chapter 13 (Testing) — YAML table-dan linter yazmağa qədər tam arsenal;
 Chapter 15 (Shipping) — konfiqurasiyadan Delve-ə qədər production
 hazırlığı. Bu iki fəsil "işə hazır" developer ilə "kod yazan"
 developer arasındakı fərqi yaradır.

## Nə oxumaq davam etmək üçün

- **100 Go Mistakes** (Harsanyi) — bu reseptlərin "anti-pattern"
  tərəfindən eyni material
- **Learning Go** (Bodner) — generics və interfeys nəzəriyyəsi dərinliyi
- **The Art of Unix Programming** (Raymond) — kitabın多次 istinad
  etdiyi mənbə (Rule of Generation, Rule of Separation)
