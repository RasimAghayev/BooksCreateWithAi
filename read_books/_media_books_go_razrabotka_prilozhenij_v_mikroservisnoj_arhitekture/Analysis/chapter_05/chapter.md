# Chapter 5 — Chapter Summary

## Book: Разработка приложений в микросервисной архитектуре с нуля
## Author: Юлия Апапова

### Chapter Title
Развертывание микросервисов (Deployment of Microservices)

### Chapter Goal
Cover deployment strategies and CI/CD automation for the microservice system. Introduces GitHub Actions for continuous integration testing and deployment, along with infrastructure setup (nginx, Docker, Makefile).

---

## Deployment Strategies (5 Approaches)

| Strategy | Description | Tradeoffs |
|----------|-------------|-----------|
| 1. One server, several processes | Local dev style — multiple processes on different ports | Single point of failure, not production-ready |
| 2. Several servers, several processes | Multiple VMs behind load balancer (nginx/Apache) | More resilient but manual setup |
| 3. Containers | Docker-packaged microservices | Portable, consistent environments |
| 4. Orchestrator (Kubernetes/Nomad) | Full container orchestration platform | Complex but handles scaling/service discovery |
| 5. Serverless | Cloud functions — no server management | Pay-per-execution, reduced ops burden |

---

## CI/CD with GitHub Actions

### Purpose
- Automated test execution
- Code linting/formatting checks
- Build verification
- Automated deployment to staging/production

### Workflow Structure (`.github/workflows/master.yaml`)
```yaml
name: Transaction
on: workflow_dispatch    # or push, pull_request, branch filters
jobs:
  stage_dev:
    runs-on: ubuntu-latest
    steps:
      - name: get code
        uses: actions/checkout@v4.5.3
      - name: install deps
        run: go mod download
      # Additional steps: lint, test, build, deploy
```

### Key Parameters
- **name**: Workflow identifier (service name or stage)
- **on**: Trigger event (`workflow_dispatch`, `push`, `pull_request`, branch-specific)
- **jobs**: Set of work to execute
- **runs-on**: Execution environment (ubuntu-latest, macos-latest, windows-latest)
- **steps**: Ordered pipeline steps
- **uses**: Pre-built action reference (e.g., `actions/checkout@v4.5.3`)
- **run**: Shell command execution

---

## Makefile (Complete Build Automation)

| Target | Description |
|--------|-------------|
| `build` | Compile binary with version/time ldflags |
| `run` | Clean + build + execute |
| `deps` | Download dependencies |
| `deps-update` | Update dependencies |
| `fmt` | Format code (gofmt, gofumpt, goimports) |
| `lint` | Run golangci-lint |
| `test` | Run all tests |
| `test-coverage` | Generate HTML coverage report |
| `bench` | Run benchmarks |
| `generate-mocks` | Auto-generate test mocks via mockgen |
| `ci` | Full pipeline: deps → fmt → lint → test → build |

---

## Unit Testing Approach

### Tools
- **gomock**: Mock interface implementation generator
- **testify/assert**: Assertion library for test verification
- **zerolog.Nop()**: No-op logger for tests

### Test Pattern: Table-Driven Tests

1. Define test struct with: name, inputs, setupMocks, expectedError, expectedResult
2. Use `gomock.NewController(t)` + `defer ctrl.Finish()` for mock management
3. Create mock implementations for Repository, KafkaPublisher, AccountService
4. Set expectations: `mockRepo.EXPECT().Deposit(...).Return(...)` 
5. Verify results with `assert.Equal()`, `assert.NoError()`, `assert.Error()`

### Tests Implemented
- **TestTransactionService_Deposit**: Success, repository error, Kafka publish failure
- **TestTransactionService_Withdraw**: Success, repository error, Kafka publish failure
- **TestTransactionService_Transfer**: Success, repository error
- **TestTransactionService_HandleAccountResponse**: Valid/invalid responses, completed/failed status
- **TestTransactionService_GetTransactionsWithDetails**: Success, repository error, partial failures

### Mock Generation
```bash
go install go.uber.org/mock/mockgen@latest
mockgen -source=internal/service/service.go \
        -destination=internal/service/mocks/mocks.go \
        -package=mocks
```

---

## nginx Configuration
```nginx
location / {
    proxy_pass http://localhost:50053;  # Proxy to Gateway service
}
```
- File edited: `/etc/nginx/sites-available/default`
- Installed via: `sudo apt install nginx` on Ubuntu 24.04

---

## What's Next
Chapter 6 likely covers final project wrap-up, Docker containerization, and production deployment considerations.
