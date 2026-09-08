# Build an Orchestrator in Go (From Scratch) — Xülasə (AZ)

## Kitab kimin üçündür?

"Kubernetes necə işləyir" sualına cavab axtaran Go developer-ləri üçün —
fərdi, kiçik "Cube" orkestratoru sıfırdan quraraq prod sistemlərinin (K8s,
Nomad, Borg) daxili mexanikasını öyrənmək. Goroutine xaricində qabaqcıl
dil xüsusiyyətləri qəsdən istifadə olunmur — fokus SİSTEM DİZAYNındadır.

## Layihə xətti (13 fəsil, 410 səh.)

| Hissə | Mövzu | Nə tikilir |
|---|---|---|
| 1 (Ch1-3) | Konsept + Task | komponent modeli, state machine, Docker SDK |
| 2 (Ch4-6) | Worker | queue+DB+runtime, chi REST API, /proc metrics |
| 3 (Ch7-9) | Manager | naive scheduler, SendWork, API, health checklər |
| 4 (Ch10-11) | Refactor | Scheduler interface + E-PVM, Store + BoltDB |
| 5 (Ch12-13) | CLI + bitiş | Cobra komandaları, hara getməli |

## Ən vacib 5 fikir

1. **Bütün orkestratorlar eyni komponentlərdir:** task/worker/manager/scheduler —
   adlar fərqli (Borg: Borglet/BorgMaster; K8s: kubelet/control plane; Nomad:
   client/server), arxitektura eyni
2. **State machine + keçid cədvəli:** task vəziyyətləri (Pending→Scheduled→
   Running→Completed/Failed) QANUNLA təsbit olunur; hər keçid = TaskEvent
   (audit izi)
3. **İnterfeyslər dəyişməni təcrid edir:** Scheduler (3 metod: feasibility/
   scoring/picking) və Store (Put/Get/List/Count) — manager kodu DƏYİŞMƏDƏN
   round-robin→E-PVM, memory→BoltDB keçidləri
4. **Məsuliyyət ayrımı:** restart məntiqi Docker-da deyil, manager-də (kim
   kimə cavab verir bulanıqlaşmasın); API handler icra etmir — queue + cavab
5. **Metrics → scheduler qida mənbəyi:** /proc (meminfo/stat/loadavg) +
   gopsutil/goprocinfo → /stats endpoint → E-PVM ballandırması

## Kitabın ən dəyərli hissəsi

Chapter 10 (E-PVM): Google Borg-un atası olan alqoritm — marginal cost
hesabı, CPU delta ölçmə (2 snapshot), ballandırma formulunu real kodla göstərir.
"Production tarixi" olan sadə kod parçası.

## Qeydlər

- Yalnız Linux (cgroups/proc); Docker SDK mütləqdir
- Kod: github.com/timboring/orchestrator-in-go
- Müəllif Tim Boring — Google-da Ganeti (Borg-dan əvvəl) təcrübəsi ilə
  "2007-də istərdim bu kitab olsaydı" motivasiyası
