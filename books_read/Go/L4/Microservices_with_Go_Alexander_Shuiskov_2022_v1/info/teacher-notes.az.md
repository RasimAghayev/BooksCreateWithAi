# Microservices with Go — Müəllim Qeydləri (AZ)

## Tədris planı (tövsiyə olunan sıra)
1. **Ch1-2 (konsept + scaffolding):** əvvəl monolit/mikro müzakirəsi, sonra
   Movie aplikasiyasının dizaynını (metadata/rating/movie bölünməsini)
   whiteboard-da tələbələrlə birgə çıxarın — kod yazmaqdan əvvəl.
2. **Ch3 (discovery):** Registry interfeysini verin, tələbələrdən in-memory
   implementasiyasını yazmalarını istəyin; Consul sonra.
3. **Ch4-5 (protobuf + gRPC):** benchmark nəticələrini (63B vs 106B vs 148B)
   canlı göstərin — motivasiya güclü olur.
4. **Ch9 (test):** gomock EXPECT()/Times() oyunu — "çağırılmadı" xətasını
   qəsdən törədin.
5. **Ch10-12 (reliability üçlüyüsü):** bir service-i qəsdən "sındırın"
   (Ctrl+C, limit aşımı), alert firing-i canlı izləyin.

## Çətin anlaşılan yerlər
- **Ch2:** Niyə hər komponentdə AYRI ErrNotFound? (rating vs metadata xətası
  qarışmasın — errors.Is ilə müqayisə düzgün qalsın)
- **Ch5:** internal vs generated model — "kod dublikatı deyil, qat ayrımıdır";
  nil-check yükü və format coupling arqumentlərini göstərin
- **Ch6:** ReadMessage(-1) — "həmişə başdan oxu"məntiqi; ReadMessage ctx-dən
  asılı DEYİL — goroutine + ctx.Done() şərti
- **Ch10:** retriable seçim klientin məsuliyyətidir — server kodları yalnız
  İMAKAN verir (DeadlineExceeded ≠ həmişə retry-lıq)
- **Ch11:** context propagation — reqId header nümunəsi ilə çətin
  istifadəçilər üçün HTTP-metodu ilə izah edin

## Sıx edilən səhvlər
1. `t.Fatalf` test dövrünü kəsir — `t.Errorf` + continue
2. Tag-lərdə UUID (cardinality partlayışı)
3. DB credential-ların koda yazılması (kitab özü Ch8-də etiraf edir)
4. `Stop` yerinə `GracefulStop` (connection leak)
5. JWT parse-da alqoritm yoxlamasının unudulması (alg confusion attack)

## Müzakirə sualları
- Rating servisinə `RecordType` əlavə etmək yaxşı abstraksiyadırmı, yoxsa
  erkən over-engineering? (Kitabın "6-12 ay" testi)
- Minik əvəzinə Consul nə vaxt haqq qazanır? (çox instans + xarici LB ehtiyacı)
- Alert nə vaxt noise olur? (actionable deyilsə — `for:` müddəti)

## Ev tapşırığı ideyaları
- Rating servisinə `RatingEventTypeDelete`-i dəstəkləyən kod əlavə et
- Consul əvəzinə K8s service discovery ilə gateway yenidən yaz
- Metadata üçün mysql repo-nun unit testini gomock ilə yaz
- pprof ilə öz servisindəki ən yavaş funksiyanı tap və optimallaşdır

## Qiymətləndirmə rubrikası
| Bacarıq | Əla | Kafi |
|---|---|---|
| Qat ayrımı | handler/controller/repo təmiz | logiki handler-də qalır |
| Xəta idarəsi | errors.Is/As + wrap | if err != nil boş |
| Observability | 3 sütun inteqre | yalnız log.Printf |
| Reliability | retry+backoff+limit | heç biri |
