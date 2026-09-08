# Go Programming B2P — Müəllim qeydləri (AZ)

## Ümumi qiymətləndirmə

Samantha Coyle-un Packt üçün yazdığı 2nd Edition — 21 fəsil, 6 hissə
strukturu ilə Go-nu "scripts → components → modules → applications →
web → professional" təbii təkamül xəttində təqdim edir. Kitabın əsas
gücü: hər anlayış Exercise formatında (adım-adım + gözlənilən çıxış)
verilir; hər Activity əvvəlki 2-3 fəslin biliklərini birləşdirir.

## Güclü tərəflər

1. **Bug-driven tədris** — shadowing (1.04), slice link (4.x), race
   condition (18.03): hər incəlik əvvəl SƏHV nəticə göstərilir, sonra
   düzəldilir — tələbə "niyə" sualına cavab alır.
2. **Qida zənciri strukturu** — Activity-lər kumulativdir (bank app
   6.01-6.04 progression-da eyni domain üçün error→panic→recover).
3. **Modern** — Go 1.18+ generics, go.work workspaces, OTel, sync.Cond,
   sync.Map, fuzzing, coverage integration (1.20) hamısı var.
4. **Sənədə bağlılıq** — hər fəslin kodu GitHub-da; şəkillər gözlənilən
   terminal çıxışlarını göstərir.

## Zəif tərəflər və müəllim tövsiyələri

1. **Bəzən səthi** — fəsillər geniş, amma dərinlik məhdud (sync.Cond
   bir exercise; embedding promotion 1 səhifə). Tam kurs üçün:
   ch18-ı 2-3 dərsə bölün (WaitGroup/atomic → kanallar → context).
2. **Exercise addımları bəzən çox qranulyar** — 15 addım "hello world"-ə
   gəlir; kodu kopyalamaq əvəzinə tələbədən addımları BİRLƏŞDİRMƏSİNİ
   istəyin (daha yüksək kognitiv yük, daha yaxşı sabitləşmə).
3. **Typos** — kitabda az-çox var (misal: Exercise nömrələri 19↔20
   qarışır kod kommentlərində); tələbələri xəbərdar edin.
4. **http.DefaultClient istifadəsi** (ch17-nin bir hissəsində) — kitab
   özü "production-da yox" deyir; bunu vurğulayın.
5. **Zənginləşdirilə biləcək yerlər:** slog (Go 1.21 structured logging)
   kitabda yoxdur — zap əvəzinə/z əlavəsində göstərmək olar; errors.Is/As
   ch8-də tam açılmır.

## Sərt xatırlatmalar (tələbələr üçün)

- 0644 OKTALDIR — baş sıfır şərt (411-dəki 511 fərqini göstərin!)
- `len(string)` BAYT sayır — Ü simvolu testində sübut edin
- Slice kopyası: `s2 := s1` LINKDIR — 4.x-in "cap-link" ssenarisini
  board-da çəkin
- loop dəyişəni closure-da: `test := test` (19.01)
- time.Sleep = konkurensiya əksidir — WaitGroup yazana qədər qəbul
  oluna bilər, sonra QADAĞANDIR (18.02-nin öz sözləri)

## Kitabın ən dəyərli hissəsi

Chapter 18 + 19 kombinasiyası: konkurensiyanın təkamül zənciri
(WaitGroup → atomic → Mutex → kanal → context → Cond/Map) və test
piramidası — bu ikilik "professional" başlığının əsl dolğunluğudur.
Həmçinin Chapter 21 — Prometheus/OTel/Docker/K8s bir fəsildə: yalnız
giriş səviyyəsində, amma deployment-ə hazırlıq checkpoint-list kimi
istifadə oluna bilər.

## Nə oxumaq davam etmək üçün

- **Learning Go** (Jon Bodner) — reflection, interface daxili mexanika
  dərinliyi
- **100 Go Mistakes** (Harsanyi) — bu kitabın bütün "bug" mövzularının
  sistemləşdirilmiş versiyası
- **Go in Action** (Kennedy) — kanal semantikası və scheduler dərinliyi
