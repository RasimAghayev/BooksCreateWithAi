# Microservices with Go — Terminologiya (AZ)

## Arxitektura
| English | Azərbaycanca | İzah |
|---|---|---|
| Microservice Architecture | Mikroservis Arxitekturası | Aplikasiyanın biznes qabiliyyətlərinə görə bölünmüş müstəqil servislər toplusu |
| Monolith | Monolit | Tək proqram kimi qurulmuş aplikasiya |
| Business Capability | Biznes Qabiliyyəti | Servisin cavab verdiyi iş sahəsi |
| Scaffolding | Skelet Qurulumu | Proyektin ilkin strukturunun yaradılması |
| Service Discovery | Servis Kəşfiyyəti | Servislərin dinamik tapılması |
| Registry | Reyestr | Aktiv instansların unvan xəritəsi |
| Load Balancer | Yük Bölücü | Request-ləri instanslar arasında paylayan komponent |
| Gateway (pattern) | Keçid qatı | Başqa servislərə çağırış kodu |
| Repository (pattern) | Repozitoriya qatı | Data saxlamanın kod qatı |
| Aggregation | Aqreqasiya | Dəyərlərin toplanması/birləşdirilməsi |
| Graceful Degradation | Zərif Enmə | Xəta zamanı məhdud amma işlək rejim |
| Fallback | Əvəzedici məntiq | Əsas əməliyyat uğursuz olanda alternativ yol |
| Ownership | Sahiblik | Servisə kimin məsul olması |

## Kommunikasiya
| English | Azərbaycanca | İzah |
|---|---|---|
| Synchronous Communication | Sinxron Rabitə | Request-response, dərhal cavab gözlənilir |
| Asynchronous Communication | Asinxron Rabitə | Göndərən dərhal cavab gözləmir |
| RPC (Remote Procedure Call) | Uzaq Prosedur Çağırışı | Şəbəkə üzərindən funksiya çağırışı |
| Serialization | Serializasiya | Datanın ötürülə bilən formata çevrilməsi |
| Wire Format | Ötürülmə Formatı | Şəbəkədəki binary təsvir |
| Protocol Buffers | Protobuf | Google-un binary serializasiya formatı |
| Publisher-Subscriber | Nəşr edən-Abunə olan | Broadcast ünsiyyət modeli |
| Message Broker | Mesaj Vasitəçisi | Mesajları çatdıran ara komponent |
| Topic | Mövzu | Kafka-da mesaj axınının kanalı |
| Partition | Bölmə | Mövzunun paralel emal hissəsi |
| Offset | Ofset | Mesajın mövzu daxilində mövqeyi |
| Delivery Guarantee | Çatdırılmə Zəmanəti | at-most/least/exactly-once |
| Backoff | Gecikməli Retry | Retry-lər arası artan gecikmə |
| Jittering | Təsadüfi Dəyişkənlik | Retry vaxtına random əlavə — burst qarşısı |
| Rate Limiting | Sorğu Limitasiyası | Paralel request sayının sınırı |
| Token Bucket | Jeton Vedrəsi | Limit alqoritmi (b həcm, r dolğun sürəti) |
| Timeout / Deadline | Vaxt limiti / Son an | Request-in max gözləməsi |
| Authentication | Kimliyin Təsdiqi | Kim olduğunu yoxlama (login) |
| Authorization | Səlahiyyətləndirmə | Nəyə icazəsi olduğunu yoxlama (role) |
| JWT (JSON Web Token) | JWT tokeni | Header.Payload.Signature strukturlu təhlükəsizlik tokeni |
| Claim | İddia | JWT payload-dakı identifikasiya sahəsi |
| Bearer Token | Daşıyıcı tokeni | Sahibinə giriş verən token |

## Reliability / İnkishaf
| English | Azərbaycanca | İzah |
|---|---|---|
| Reliability | Etibarlılıq | Gözlənilən iş + açıq limitlər |
| Blast Radius | Təsir Dairəsi | Xətanın yayılma sahəsi |
| Fault Isolation | Qəza İzolyasiyası | Bir hissənin ölümü sistemi öldürmür |
| Design for Failure | Uğursuzluq üçün Dizayn | Xətaları gözləyən dizayn |
| Graceful Shutdown | Zəif Dayandırma | Resursları bağlayaraq dayanma |
| Connection Leak | Bağlantı Sızıntısı | Bağlanılmayan DB bağlantısı |
| Canary Deployment | Kanari Yayımı | Kiçik fraqmentdə sınaq yayımı |
| Continuous Deployment (CD) | Fasiləsiz Yayım | Hər dəyişiklikdə avtomatik deploy |
| Rollback | Geri Qaytarma | Uğursuz deployment-in geri alınması |
| Environment | Mühit | dev/staging/prod |
| On-call | Növbədə | İnsident bildirişlərini qəbul edən mühəndis |
| Runbook | Əməliyyat Kitabçası | İnsident mitigation addımları |
| Postmortem | Hadisə Hesabatı | İnsidentdən sonra öyrənilən dərslər |
| Five Whys | Beş Niyə | Kök səbəb analizi üsulu |
| TTR (Time to Repair) | Təmir Vaxtı | İnsidentin həlli müddəti |
| Reliability Drill | Etibarlılıq Məşqi | Planlaşdırılmış failure sınağı |
| Escalation Policy | Eskalasiya Siyasəti | Cavabsız bildirişin yuxarıya ötürülməsi |

## Observability
| English | Azərbaycanca | İzah |
|---|---|---|
| Observability | Müşahidə Qabiliyyəti | Sistemin daxili halının ölçülə biləməsi |
| Telemetry | Telemetriya | Logs + Metrics + Traces |
| Structured Logging | Strukturlaşdırılmış Loglama | JSON + level + sahələr |
| Log Level | Log Səviyyəsi | info/warn/error/fatal/debug |
| Metric | Metrik | Zaman seriyalı kəmiyyət ölçməsi |
| Counter / Gauge / Histogram | Saya / Göstərici / Histoqram | Metric növləri |
| Cardinality | Kardinalite | Unikal tag dəyərlərinin sayı |
| Distributed Tracing | Paylanmış İzləmə | Request-in servislər arası səyahəti |
| Span | Aralıq | Bir əməliyyatın trace elementi |
| Context Propagation | Kontekst Yayımı | Metadata-nın çağırış zənciri ilə ötürülməsi |
| Profiling | Profilləşdirmə | CPU/heap performans analizi |
| Percentile (p90/p95/p99) | Persentil | Paylanmanın üst hədləri |
| Saturation | Doluluq | Resurs istifadə səviyyəsi |
| Retention | Saxlama Müddəti | Datanın nə qədər saxlanması |
| Inversion of Control (IoC) | İdarənin İnversiyası | Framework-ün icranı özünə götürməsi |
| Black-box Monitoring | Qara qutu monitorinqi | Yalnız xarici göstəricilərlə izləmə |
| White-box Monitoring | Ağ qutu monitorinqi | Daxili data ilə izləmə |
