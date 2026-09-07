# Go mikroservisləri (Popova 2026) — Müəllim Qeydləri (Azərbaycanca)

📖 **Kitab deyir:**

Mikroservis arxitekturası Go ilə sıfırdan: 5 servis (Account, Auth, Gateway,
Transaction + DB-lər), Protobuf kontraktlar ayrı repo-da, 3-qatlı təmiz
arxitektura hər qatın öz modeli ilə, JWT interceptor, double-entry
maliyyə, Saga pattern Kafka ilə, sonda Docker multi-stage + Docker Hub +
Kubernetes kompose miqrasiyası. Kitab "vahid layihə" prinsipi ilə gedir —
hər fəsil əvvəlkinin üzərinə qurulur.

👨‍🏫 **Müəllim qeydi:**

Bu kitab rusdilli ekosistemdə Go mikroservisləri üzrə nadir və dəyərli
praktik materialdır. Güclü tərəfi: kod arxivi real işləyən sistemdir və
kitab boyu ARXİTEKTURA PROQRESİYASI var — tək servisdən paylanmış sistemə.
Tələbəyə tövsiyələrim:

1. **3-qatlı model ayrılığını BAŞ ADIM kimi qəbul et:** kitabın ən dəyərli
   dərsi budur — gRPC/DB tiplərinin biznesə sızması ANTİPATTERN. Bu
   prinsipi hər yeni layihədə birinci gündən tətbiq et.
2. **Saga vs 2PC müqayisəsini mənimsə:** fəsil 4 bu mövzunu rus dilində ən
   yaxşı izahlardan birində verir. Kompensasiya əməliyyatlarının mürəkkəbliyi
   real sistemlərdə haqlıdır — status maşını (pending → completed/failed)
   bu mürəkkəbliyi görünür edir.
3. **Kitabın zəif nöqtələrinə diqqət:** (a) Kubernetes səthidir — prod üçün
   readiness/liveness probe-lar, resource limitləri, HPA, ingress, secret
   idarəetməsi YOXDUR; (b) observability demək olar ki, yoxdur (Prometheus/
   Loki yalnız adlarıyla); (c) grpc.Dial deprecated-yə meyllidir —
   grpc.NewClient (fəsil 3-də düzgündür, fəsil 4-də köhnə); (d) interfeys
   tipləri bəzi yerlərdə float32 balans ilə — real pulla işləməzdən əvvəl
   int64 qəpiklərə keçmək lazımdır (kitabın ÖZÜ bunu Transaction-da edir,
   Account-da yaddan çıxır).
4. **Öyrənmə yolu:** arxivi klonla → hər fəsilin kodunu ÖZÜN yaz → testləri
   başqa ssenarilərlə genişləndir (kitabın 84.4% coverage-i yaxşı başlanğıcdır) →
   compose-u minikube-də qaldır → bir komponenti sındır (Kafka-nı dayandır)
   və Saga status maşınının davranışını İZLƏ. Bu sonuncu təcrübə kitabın ən
   dərin dərsini — paylanmış sistemin xəta davranışını — canlı göstərir.

## Ən vacib 5 fikir

1. **Hər qatın öz modeli:** dəyişikliklər mapper-də lokallaşır; gRPC-dən
   asılı biznes kodu böyük refaktor tələb edir (F1, F3).
2. **Pul = int64 qəpik + double-entry:** hər əməliyyat DEBIT+CREDIT
   qeydləri ilə; float maliyyədə QADAĞANDIR (F4).
3. **Saga = status + kompensasiya:** distributed tranzaksiyadan qaç;
   qaçılmazdırsa pending/in-processed statusları intermediate halları
   görünür edir (F4).
4. **Kontraktlar ayrı repo + semver:** cyclic import = dizayn xətası;
   breaking change = major (F1, F3).
5. **Avtomatlaşdırma pilləli:** Makefile → CI/CD → Docker → K8s; hər addım
   manual xərci azaldır, amma 5-ci üsulu (serverless) istisnadır — öz
   infrastruktur qərarını texniki AMAL deyil, biznes ehtiyacı müəyyən edir.

## Kitabın ən dəyərli hissəsi

Fəsil 4 (səh. 157-266) — 110 səhifə: normal formalardan Saga-nın tam
Kafka implementasiyasınadək. Bu fəsil kitabı digər Go giriş kitablarından
ayıran yerdir. İkinci əhəmiyyətli: Fəsil 1-in DI + 3-qatlı arxitektura
bölməsi — burada qoyulan struktur bütün kitab boyu işləyir.

## Kim üçündür

Go sintaksisini bilən (kitab dili öyrətmir!), amma mikroservis
arxitekturasında təcrübəsi az olanlar üçün. L3 tövsiyə olunur: gRPC,
Kafka, PostgreSQL əsaslarından xəbərdarlıq fərqli olaraq L2 oxucu fəsil
3-4-də çətinlik çəkəcək.
