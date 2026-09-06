# Chapter 0 — Acknowledgments / About Go Optimizations 101

## Bu chapter nədən bəhs edir?
Müəllifin minnətdarlıqlarına və kitabın məqsəd/üsuluna: performans optimizasiyaları trade-off dünyasında yaşayır.

## Əsas fikirlər

### Kitabın məqsədi
Kod performansı optimizasiya fəndləri, məsləhətləri — əsasən rəsmi Go compiler/runtime implementasiyasına əsaslanır (v1.19 istinad nöqtəsi).

### Trade-off fəlsəfəsi (Chapter 1)
- Proqramlaşdırma trade-offlarla doludur: oxunaqlılıq, saxlanıla bilmə, inkişaf sürəti, proqram effektiviliyi
- Hətta effektivliyin İÇİNDƏ belə: yaddaş qənaəti vs icra sürəti vs implementasiya çətinliyi
- **Praktik prinsip:** Layihə kodunun əksər hissəsi yüksək performans tələb ETMİR — saxlanıla bilmə və oxunaqlıq daha vacibdir; hot path-ləri tap və yalnız oranı optimallaşdır

### Kitabın əhatə qaydaları (müəllifin xəbərdarlıqları)
- Bəzi məsləhətlər hər platformada, bəziləri yalnız müəyyən CPU/arxitekturada işləyir → **production mühitində benchmark et**
- Compiler/runtime detalları versiyadan-versiyaya dəyişir → bəzi fəndlər gələcəkdə işləməyə bilər
- Kitabın formatı: fəsl-fəsil open-source olacaq (2023/iyun hədəfi)

### Müəllif
Tapir Liu — Go 101 kitabının da müəllifi (go101.org); indie oyun developer keçmişi (tapirgames.com).

## Praktik nəticə
Optimizasiyaya başlamazdan əvvəl: (1) profile et, (2) hot path-i tap, (3) yalnız oranı optimallaşdır, (4) hər fəndlə istehad mühitində benchmark ilə təsdiqlə.

## Mənbə
Pages: 6-8 (PDF səh. 6-8)
