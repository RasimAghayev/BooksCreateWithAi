# Event-Driven Architecture in Golang — Müəllim Qeydləri (AZ)

📖 **Kitab deyir:**

Kitab MallBots üzərindən EDA-nın tam dövrünü göstərir: EventStorming ilə
kəşf → DDD/hexagonal modullar → domain hadisələri → event sourcing →
JetStream inteqrasiyası → state transfer keşləri → saga → inbox/outbox →
test piramidi → AWS deploy → observability.

👨‍🏫 **Müəllim qeydi (əlavə tövsiyələr):**

1. **Productionda NATS seçimi.** Kitab JetStream-i göstərir, amma real sistemdə
   Kafka (çox böyük ekosistem, tam dəstək) və ya Redpanda (Kafka-uyğun,
   daha yüngül) seçimlərini də müqayisə edin. Kiçik/orta sistemlər üçün NATS
   sadəliyi real üstünlükdür.

2. **Event store üçün xüsusi DB.** Kitabdakı Postgres events cədvəli öyrədici
   məqsədlidir. Miqyaslanan sistemlərdə EventStoreDB, Marten (Postgres üzərində)
   və ya Axon kimi hazır həllərə baxın — snapshot/dedup/versioning artıq
   daxildir.

3. **"Just enough" prinsipi.** Kitabın ən vacib mesajlarından biri də budur —
   CQRS-i, event sourcing-i, saga-nı HƏR YERDƏ etmə. Sadə CRUD modulu üçün
   bunlar yükdür. Core domain-ə saxla.

4. **Saga vs workflow engine.** SEC öz implementasiyası maarifləndiricidir, amma
   Temporal (Go SDK) və ya Cadence kimi hazır workflow engine-lər productionda
   davamlılıq, replay, vizuallaşdırma verir. Kiçik sistemlərdə sadə saga kifayətdir.

5. **Idempotensiya testləri.** Kitab dedupe mexanizmlərini göstərir, amma
   testlərin duplicate mesaj ssenarilərini (eyni mesajı 2 dəfə göndər) avtomatik
   yoxlamasını təkidlə tövsiyə edirəm — migration/rebalance zamanı bu qaçınılmazdır.

6. **AsyncAPI + schema registry.** AsyncAPI sənədləşməsini schema registry
   (Confluent Schema Registry və ya sadə git-də protobuf versiyalanması) ilə
   birləşdirin — hadisə kontraktının drift-i ən gizli EDA xətasıdır.

## Ən vacib 5 fikir

1. Hadisə immutable fakt; cari vəziyyət = hadisələrin toplamı (və ya
   proyeksiyası)
2. Dual write EDA-nın anadangəlmə günahıdır — outbox mütləqdir
3. Eventually consistent sistemi idarə etmək üçün UI (task-based) və monitoring
   dəyişməlidir
4. Saga = kiçik tranzaksiyalar + kompensasiyalar; 2PC paylanmış sistemə yaraşmır
5. Contract testlər (Pact) modullararası əlaqələri inteqrasiya testlərindən
   ucuz qoruyur

## Kitabın ən dəyərli hissəsi

Chapter 9 (Transactional Messaging) — inbox/outbox pattern-in tam
implementasiyası. Bu pattern istənilən EDA sistemində (dilindən asılı olmayaraq)
istifadə olunur və ən az sənədlənən mövzulardan biridir. Kodu real layihəyə
kopyalayaraq başlamaq olar.

## Sonrakı oxuma

- *Versioning in an Event-Sourced System* (Derek Comartin) — hadisə versiyası
- *Microservices Patterns* (Chris Richardson) — saga/2PC daha dərin
- *Learning Domain-Driven Design* (Vlad Khononov) — DDD strukturasiyası
- NATS rəsmi docs — JetStream konfiqurasiyası (retention, limitlər)
