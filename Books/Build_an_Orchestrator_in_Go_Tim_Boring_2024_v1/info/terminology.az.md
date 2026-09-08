# Build an Orchestrator in Go — Terminologiya (AZ)

| Termin | Azərbaycanca | İzah |
|---|---|---|
| Orchestrator | orkestrator | konteynerlərin avtomatik deploy/scale/idarəsi |
| Container | konteyner | namespace+cgroups izolyasiyası; kernel YOXDUR (konsept!) |
| VM vs Container | VM vs konteyner | VM: öz kernel; konteyner: host kernel paylaşır |
| Task | tapşırıq | ən kiçik iş vahidi = 1 konteyner (1:1) |
| Job | iş | yüksək səviyyə: n instans, datacenter, update siyasəti |
| Worker | işçi | task icra edən node (K8s: kubelet; Borg: Borglet) |
| Manager | menecer | beyin/control plane (K8s: control plane; Borg: BorgMaster) |
| Scheduler | planlayıcı | task→node seçimi |
| Cluster | klaster | bütün komponentlərin məntiqi qrupu |
| Control/Data plane | idarə/məlumat müstəvisi | manager qərar / worker icra |
| State machine | vəziyyət maşını | Pending→Scheduled→Running→Completed/Failed |
| State transition map | keçid xəritəsi | valid keçidlərin cədvəli (mənbə həqiqəti) |
| TaskEvent | tapşırıq hadisəsi | hər keçidin jurnal qeydi |
| Feasibility | uyğunluq | resurs filtri (kənarlaşdırma) |
| Scoring | ballandırma | hər namizədə rəqəmsal dəyər |
| Picking | seçim | ən yaxşı ballı node |
| E-PVM | E-PVM | Enhanced PVM — Borg-un əcdadı; marginal cost əsaslı |
| Marginal cost | kənar yük | node-a əlavə taskın gətirəcəyi yük |
| Round-robin | dairəvi növbə | sadə ardıcıl seçim (naive baseline) |
| Docker Engine API | Docker API | SDK: ImagePull/ContainerCreate/Start/Stop/Inspect |
| chi router | chi | router + URLParam middleware |
| DisallowUnknownFields | sərt dekodlaşdırma | naməlum JSON sahə = 400 |
| /proc filesystem | /proc | meminfo, stat, loadavg, diskstats — kernel pəncərəsi |
| Load average | yük ortalaması | 1/5/15 dəq queue ölçüsü |
| CpuPercentUsed | CPU istifadəsi | (total-idle)/total — 2 snapshot deltası |
| Health check | sağlamlıq yoxlaması | GET /health + status 200 |
| RestartCount | yenidənbaşatma sayı | 3 limit — sonsuz loop qarşısı |
| Exponential backoff | eksponensial gecikmə | retry aralarının böyüməsi |
| Store interface | anbar interfeysi | Put/Get/List/Count — Delete YOX (tarix!) |
| BoltDB | BoltDB | embedded key-value; bucket + tx; fayl = DB |
| Bucket | buckıt | BoltDB namespace-i |
| View / Update | View / Update | oxu / yazma transaction-u |
| Cobra | Cobra | CLI framework (kubectl, docker, hugo...) |
| PersistentFlags | persistent bayraqlar | alt-komandalara miras |
| cobra-cli | cobra-cli | scaffold aləti (init/add) |
| Service discovery | xidmət kəşfi | komponentlər bir-birini necə tapır |
| Consensus | konsensus | paylanmış dəyər razılığı (kitabda kənar) |
| Load balancing | yük tarazlığı | request-lərin node-lar arası paylanması |

## Komponent ↔ real sistemlər

| Konsept | Cube | Borg | K8s | Nomad |
|---|---|---|---|---|
| Manager | Manager | BorgMaster | control plane | server |
| Worker | Worker | Borglet | kubelet | client |
| Task | Task | alloc | pod | task |
| Scheduler | Scheduler/ | (içində) | kube-scheduler | (serverdə) |
