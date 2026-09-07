# The Ultimate Guide To Building Database-Intensive Apps with Go — Mündəricat

**Müəllif:** Baron Schwartz | **Nəşriyyat:** VividCortex | **İl:** 2019 | **Səviyyə:** L3 (Intermediate)

Qısa (40 səh.) amma SİXILI texniki kitab: Go `database/sql` paketinin bütün
dərinliyi — pool-dan leak-ə, prepared statement gizli davranışından
tranzaksiya tələlərinə qədər.

## Mövzular (tək böyük chapter kimi emal)

| # | Mövzu | Səhifə |
|---|-------|--------|
| 1 | [database/sql — tam dərs](../chapters/01-database-sql/index.md) | 1-39 |

Chapter daxilində bölmələr: sql.DB və pool mexanizmi → nəticə setləri
(rows.Next/Scan/Close) → tək sətir (QueryRow) → NULL tipləri → naməlum
sütunlar → çoxlu nəticə dəstləri → Exec və Result → prepared statement-lər →
tranzaksiyalar (sql.Tx) → sql.Conn → xəta idarəetməsi → Valuer/Scanner →
db.Stats monitoring → context → driverlər → 13 pitfall.

## Kitabın xarakteri

"Yeganə yararlı cheat-sheet" tipli kitabdır: kod nümunələri minimal amma
hər biri production dərsi ilə; müəllif High Performance MySQL (O'Reilly)
müəllifi — DB performans perspektivi güclüdür. Oxunuş müddəti: ~1 saat;
istinad kimi daimi dəyər.
