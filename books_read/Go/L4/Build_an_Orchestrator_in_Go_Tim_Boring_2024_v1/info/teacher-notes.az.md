# Build an Orchestrator in Go — Müəllim Qeydləri (AZ)

📖 Kitab deyir: minimal "Cube" orkestratoru ilə Kubernetes-in daxili mexanikasını
öyrən — scheduler-dən CLI-yə bütün hissələri öz əlinlə yaz.

👨‍🏫 Müəllim qeydi: kurs kimi əla qurulub (hissə-hissə, hər fəsil işləyən kod);
amma oxucu 3 praqmatik qeydi bilsin:

1. **Sadəlikdən güzəştələr AÇIQ deyilən yerlər:** tək manager (HA yoxdur),
   map-lər mutex-siz (paralel yazış yoxdur deyə "işləyir"), kanal istifadəsi
   demək olar minimal. Bu, tədris seçimidir — real sistemdə manager
   replikasiyası (consensus: Raft) və kilidləmə MÜTLƏQ olardı. Kitab bunu
   düzgün qeyd edir, amma oxucu "prod patterns" siqnallarını ayırmalıdır.

2. **E-PVM ballandırması sadələşdirilib:** real E-PVM paper-da daha çox amil
   var; kitabda 0.6/0.4 tipli ağırlıqlar empirik göstərilir. Dərs: scoring
   FUNKSİYASI dəyişəndir — biznes tələbindən asılı; universal formula YOXDUR.
   K8s-də eyni: plugin-lərlə genişlənən ballandırma.

3. **Worker API-də avtorizasiya yoxdur:** hər hansı POST /tasks worker-ə
  Komanda verə bilər. Kitab Chapter 1-də security bölməsində bunu "kənar
   mövzu" kimi qeyd edir — amma DEMO belə olsa, real deployda manager-worker
   arası mTLS/token ilk gün əlavə olunmalıdır. Nomad/K8s bunu daxili olaraq
   həll edir.

## Ən vacib 5 fikir

1. Komponent modeli universal: task/worker/manager/scheduler
2. State machine + keçid cədvəli + TaskEvent = auditable sistem
3. İnterfeyslər (Scheduler, Store) — dəyişmə bir yerdə cəmlənir
4. Restart məsuliyyəti manager-də; limit (RestartCount<3) sonsuz loop qarşısı
5. Metrics axını: /proc → worker /stats → scheduler qərarı

## Kitabın ən dəyərli hissəsi

Chapter 2-3 keçidi: soyuq konseptdən (state cədvəli) işləyən Docker API
çağırışına — "skeletdən ətə" transformasiya prosesi bütün sistem layihələrində
təkrarlanan qəlibdir.

⚠️ Uyğunsuzluq YOXDUR; kod və mətn tutarlıdır. Yalnız versiya qeydi: kitab
Go 1.18-ə yazılıb (generics istifadəsindən qəsdən imtina açıq izah olunur);
müasir oxucu interface{}-ni any kimi oxuya bilər.
