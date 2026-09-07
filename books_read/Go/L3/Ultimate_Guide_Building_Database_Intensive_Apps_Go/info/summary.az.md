# The Ultimate Guide To Building Database-Intensive Apps with Go — Xülasə (Azərbaycanca)

**Müəllif:** Baron Schwartz | **Nəşriyyat:** VividCortex | **İl:** 2019 (Go 1.12-yə qədər) | **Səviyyə:** L3

## Kitabın ümumi məqsədi

Bu 40 səhifəlik pulsuz e-kitab Go-nun standart `database/sql` paketini
production keyfiyyətində istifadə etməyi öyrədir. Müəllif (High Performance
MySQL müəllifi, VividCortex qurucusu) illər boyu toplanmış təcrübəni
sıxılmış formada verir: bağlantı pool-unun DƏQİQ işləmə mexanizmi, nəticə
setlərinin həyat dövrü, prepared statement-lərin gizli davranışları,
tranzaksiya tələləri, xəta idarəetməsi, şəffaf data transformasiyaları və
13 real dünya pitfall-u. Kitabın əsas tezisi: database/sql "maqiyasızdır"
— içini bilsən, connection pool, retry, çevrilmələr hamısı SƏNİN üçün işləyir;
bilməsən — leak və fəlakətlə qarşılaşırsan.

## Əsas mövzuların xülasəsi

1. **sql.DB = pool:** Open bağlantı AÇMIR (lazy); Ping ilə doğrula; bir dəfə
   yarat, daimi saxla; hər sorğu üçün yeni db = TCP TIME_WAIT fəlakəti.
2. **Bağlantı sahibliyi metoddan asılıdır:** Exec dərhal qaytarır; Query —
   rows.Close-a qədər; QueryRow — Scan-a qədər; Begin — Commit/Rollback-a
   qədər. İgnor edilmiş rows = leak = "too many connections".
3. **rows üçlüyü:** defer/erkən Close + döngüdən sonra mütləq rows.Err().
4. **Scan ağıllıdır:** VARCHAR→float avtomatik; NULL üçün NullXXX ({T,
   Valid}) tipləri; sıra üzrə işlənir — ilk xətada qalanlar dolu qalır.
5. **Prepared statements:** parametrli çağırışlar GİZLİ prepare/execute/close
   edir; tək-istifadə = 3x raund-trip; yüksək konkurrentsdə pool səbəbi ilə
   təkrar-təkrar prepare = server-də statement çoxalması.
6. **sql.Tx:** db.Exec("BEGIN") QATİ YANLIŞ (hər Exec başqa bağlantıda!);
   tx tək bağlantı tutur — tx içində paralel sorğu YOX; retry avtomatik
   DEYİL; db-yə müraciət tranzaksiyadan kənardadır.
7. **sql.Conn (1.9+):** tranzaksiyasız tək-bağlantı zəmanəti (temp
   cədvəllər, USE, lock-lar).
8. **Xəta idarəetməsi:** string-match (dil asılı!) YOX → driver struct
   assertion + adlı konstantlar (mysqlerr.ER_LOCK_DEADLOCK).
9. **Valuer/Scanner:** yazılış/oxunuş transformasiyası — lowercase, gzip,
   şəffaf ENKRİPSİYA; biznes kodu təmiz qalır.
10. **db.Stats() + context:** pool metriyaları (WaitCount/Duration — pool
    gözləmələri SORĞU yavaşlığından ayrılır); ContextContext 1.8+ (ləğv/
    timeout; driver dəstəyi fərqlidir).
11. **13 pitfall:** döngüdə defer, çoxlu db, Close unutmaları, uint64 MSB,
    NULL sürprizi, eyni-bağlantı gümanı...

## Kitabın əsas mesajları

1. database/sql SƏNİN kodunu təmiz saxlamaq üçün ÇOX işi özü edir — pool,
   retry (10x), çevrilmələr; amma resurs sahibliyi SƏNƏ qalıb.
2. Leak-in 2 mənbəyi: ignor edilmiş Query nəticəsi və unudulmuş Close.
3. Tranzaksiya = tək bağlantı bağlanması; db və tx ayrı dünyalardır.
4. Driver spesifikaları (parametr sintaksisi, xəta kodları, context dəstəyi)
   — "write once, run anywhere" YOX, driver-i TANU.

## Kim üçündür

Go-da DB ilə işləyən HƏR KƏS üçün mütləq oxunuş — hətta ORM (GORM və s.)
istifadə etsən belə, altındaki database/sql davranışını bilmək leak və
performance problemlərini izah edir. L3: Go əsasları + SQL bilik tələb edir.
