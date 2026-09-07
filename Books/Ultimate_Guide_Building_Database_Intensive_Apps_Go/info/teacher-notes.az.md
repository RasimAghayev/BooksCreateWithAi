# The Ultimate Guide (database/sql) — Müəllim Qeydləri (Azərbaycanca)

📖 **Kitab deyir:**

database/sql kiçik API-yə malikdir amma dərin semantikası var: sql.DB bir
bağlantı deyil, pool-dur; hər metodun bağlantı sahibliyi qaydası fərqlidir
(Exec dərhal, Query Close-a qədər); prepared statement-lər gizli şəkildə
parametrli çağırışlarda işə düşür və pool səbəbi ilə təkrar prepare oluna
bilir; tranzaksiya tək bağlantı deməkdir və tx daxilində konkurrensiya
yoxdur; Scan avtomatik çevrirmələr və NullXXX-lərlə kömək edir; driver
səviyyəsində xətalar kodlarla tutulmalıdır; Valuer/Scanner şəffaf
transformasiya qapısıdır; db.Stats() pool sağlamlığının göstəricisidir.

👨‍🏫 **Müəllim qeydi:**

Bu kitabı hər Go backend tələbəsinin İLK həftəsində oxutmaq istəyirdim —
çünki real layihələrdə görülən "misterioz" DB problemlərinin 90%-i buradakı
pitfall-lardır. Praktik şərhlərim:

1. **"Exec vs Query leak" dərsləri hələ də aktualdır** — müasir GORM/sqlx
   istifadəçiləri belə arxada eyni mexanizmlərlə üzləşirlər. GORM-un
   `Rows()` +RowsAffected lazımlandıqda unudulan `rows.Close()` eyni leak-i
   yaradır.
2. **Kitabın 2019 tarixi haqqında:** context bölməsi bugün daha da vacibdir —
   lib/pq artıq pgx ilə əvəzlənir (pgx native context dəstəyi verir);
   sql.NullXXX yerinə çoxları `sql.Null[T]` (Go 1.22 generics) istifadə
   edir; amma KONSEPTUAL dərslər dəyişməzdir.
3. **Maraqlı detail:** uint64 MSB pitfall — driver dəstəyi yoxdursa
   `fmt.Sprint()` həlli hələ də forumlarda tapılır.
4. **Təcrübə tövsiyəm:** kitabı oxuyub bitirəndən sonra öz kod bazanızda
   "db.Query(" axtarışı edib nəticəsinin ignor edilib-edilmədiyini yoxlayın —
   production leak auditinin ən ucuz forması budur.

## Ən vacib 5 fikir

1. **sql.DB pool-dur:** global, uzunömürlü, Ping ilə doğrulanan — sorğu
   başına Open YOX (səh. 7-10).
2. **Bağlantı sahibliyi metoddan asılıdır:** ignor edilmiş Query nəticəsi =
   leak = "too many connections" (səh. 17-18) — kitabın ən vacib səhifəsi.
3. **Tranzaksiya = tək bağlantı:** tx.Exec (db.Exec YOX); tx içində
   paralellik QADAĞA; retry öz üzərində (səh. 23-26).
4. **Scan ağıllı, amma NULL-a qarşı kor:** NullXXX + schema audit — "testdə
   işləyir, prod-da partlayır" klassikası (səh. 14-15, 36).
5. **Xətalar KOD ilə tutulur:** driver struct assertion + adlı konstantlar;
   string match = dil əsaslı partlama (səh. 27-28).

## Kitabın ən dəyərli hissəsi

"Common Pitfalls" (səh. 35-36) + "Modifying Data With db.Exec()" (səh. 17-18)
— birlikdə 2 səhifə, amma illər DB debugging təcrübəsini əvəz edir.

## Kim üçündür

Go backend başlayanlar üçün MÜTLƏQ; ORM istifadəçiləri üçün də faydalı
(arxada nə baş verdiyini bilmək); DBA-lar üçün Go tərəfini anlamaq üçün ideal
qısa mətn. Səviyyə: L3 (SQL + Go əsasları tələb olunur).
