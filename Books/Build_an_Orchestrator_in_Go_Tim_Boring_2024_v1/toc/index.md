# Build an Orchestrator in Go (From Scratch)

**Müəllif:** Tim Boring (Manning, 2024) · **Səviyyə:** L4 Advanced · 410 səhifə

"Cube" adlı minimal konteyner orkestratoru sıfırdan: task-lar, worker-lər,
manager, scheduler, REST API-lər, metrics, persistent storage (BoltDB), CLI
(Cobra). Kubernetes-in daxili mexanikasını öyrənmək üçün.

## Fəsillər

| # | Mövzu | Səh. | Fayl |
|---|---|---|---|
| 1 | Orkestrator nədir? (komponentlər) | 24-58 | [01](../chapters/01-what-is-orchestrator/index.md) |
| 2 | Skelet kod | 59-77 | [02](../chapters/02-skeleton-code/index.md) |
| 3 | Task skeleti (Docker API) | 78-103 | [03](../chapters/03-task-skeleton-flesh/index.md) |
| 4 | Worker-lər | 104-128 | [04](../chapters/04-workers/index.md) |
| 5 | Worker API | 129-154 | [05](../chapters/05-worker-api/index.md) |
| 6 | Metrics (/proc, gopsutil) | 155-177 | [06](../chapters/06-metrics/index.md) |
| 7 | Manager | 178-202 | [07](../chapters/07-manager/index.md) |
| 8 | Manager API | 203-226 | [08](../chapters/08-manager-api/index.md) |
| 9 | Xətalar və sağlamlıq | 227-259 | [09](../chapters/09-failures-resiliency/index.md) |
| 10 | Mürəkkəb scheduler (Elephant algorithm) | 260-296 | [10](../chapters/10-sophisticated-scheduler/index.md) |
| 11 | BoltDB persistent storage | 297-341 | [11](../chapters/11-persistent-storage/index.md) |
| 12 | Cobra CLI | 342-382 | [12](../chapters/12-cli/index.md) |
| 13 | Növbəti addımlar | 383-404 | [13](../chapters/13-now-what/index.md) |
