# Методика выполнения задачи
## Книга: «Разработка приложений в микросервисной архитектуре с нуля» (Julia Popova, БХВ-Петербург, 2026)

## Информация о книге
- **Автор**: Юлия Попова
- **Издательство**: БХВ-Петербург
- **Год издания**: 2026
- **ISBN**: 978-5-9775-2120-8
- **Страниц**: 322
- **Язык**: Русский
- **Ссылка на архив**: https://zip.bhv.ru/9785977521208.zip

---

## Глава 1: Введение и настройка окружения (страницы 1-50)

### Введение (стр. 1-15)
- **Монолитная архитектура**: все компоненты приложения работают как единое целое
  - **Плюсы**: простота разработки, отладки и тестирования, отсутствие сетевых задержек
  - **Минусы**: высокая связность, сложность масштабирования, долгое время сборки
- **Микросервисная архитектура**: приложение разбивается на маленькие независимые сервисы
  - **Плюсы**: независимая масштабируемость, технологическая разнородность, устойчивость
  - **Минусы**: распределенная сложность, сетевые задержки, сложность отладки
- **Ключевые проблемы микросервисов**:
  - Согласованность данных (распределенные транзакции)
  - Управление сервисом и обнаружение (Service Discovery)
  - Сетевые задержки и сбои
  - Мониторинг и трассировка

### Настройка окружения (стр. 16-50)
- **Git workflow**:
  - Использование git flow (feature branches, merge requests)
  - .gitignore для Go-проектов (.env, vendor, bin, .idea)
- **Go установка**:
  - GVM (Go Version Manager) для управления версиями
  - Go Modules (go.mod/go.sum) вместо старых менеджеров (glide, godep)
  - `$GOPATH` — путь к рабочей директории Go
- **Линтеры и форматеры**:
  - `gofmt` — стандартный форматер
  - `goimports` — управление импортами
  - `golangci-lint` — комплексная проверка (lint, vet, staticcheck)
- **Makefile**: стандартные команды (build, test, run, migrate, lint)
- **Переменные окружения**: вынесение конфигурации в .env файлы
- **ZeroLog**: структурированное логирование в JSON формате
- **Docker Compose**: PostgreSQL для локальной разработки

---

## Глава 2: Работа с базами данных и ORM (страницы 51-90)

### ORM и репозитории (стр. 51-70)
- **GORM**:
  - ActiveRecord подход: модели содержат методы CRUD
  - AutoMigrate для автоматического создания таблиц
  - Associations: HasOne, HasMany, BelongsTo, ManyToMany
  - Hooks: BeforeCreate, AfterCreate, BeforeUpdate, AfterUpdate
- **Repository Pattern**:
  - Интерфейсы для репозиториев определяются в доменном слое
  - GORM-реализация репозитория в инфраструктурном слое
  - Моки для тестирования сервисного слоя
- **Миграции**:
  - Использование goose для управления миграциями
  - SQL-файлы с версионированием

### Protobuf и gRPC (стр. 71-90)
- **Protocol Buffers (Proto3)**:
  - Определение сервисов, сообщений, методов
  - `protoc` компилятор с плагинами `protoc-gen-go` и `protoc-gen-go-grpc`
  - Типы: string, uint64, double, bool, repeated, oneof, enum
- **gRPC**:
  - HTTP/2 транспорт
  - Unary вызовы (один запрос → один ответ)
  - Server-streaming, Client-streaming, Bidirectional streaming
  - Interceptor для middleware (логирование, аутентификация)
- **Архитектурные слои**:
  - Domain layer: бизнес-логика, модели
  - Repository layer: интерфейсы для доступа к данным
  - Transport/Delivery layer: gRPC/REST обработчики
  - Service layer: оркестрация бизнес-операций
  - Mapper: преобразование между domain и transport моделями

---

## Глава 3: Сервисы (страницы 91-150)

### HTTP и протоколы (стр. 91-110)
- **HTTP/1.1 vs HTTP/2**:
  - HTTP/2: мультиплексирование, бинарный протокол, сжатие заголовков
- **gRPC vs REST**:
  - gRPC: бинарный протокол, схема в proto, типобезопасность
  - REST: текстовый (JSON), человекочитаемый, широкая поддержка
- **RabbitMQ (AMQP)**:
  - Обмен сообщениями: Producer → Exchange → Queue → Consumer
  - Типы exchanges: direct, fanout, topic, headers
  - Виртуальные хосты и учетные записи
- **Apache Kafka**:
  - Топики, партиции, реплики
  - Producers publish сообщения в топики
  - Consumers читают из топиков (consumer group)
  - Использование библиотеки `segmentio/kafka-go`

### Gateway и аутентификация (стр. 111-130)
- **API Gateway**:
  - Единая точка входа для всех клиентских запросов
  - Маршрутизация запросов к соответствующим микросервим
  - Interceptor для аутентификации и логирования
  - Реализация: Register, Login, Logout, Refresh, GetCurrentUser
- **Auth service**:
  - bcrypt для хеширования паролей
  - JWT (JSON Web Token): Access Token и Refresh Token
  - Claims: user_id, exp (время истечения), iat (время создания)
  - Access Token: короткий TTL (60 минут)
  - Refresh Token: длинный TTL (30 дней)
  - Переменные: JWT_SECRET, ACCESS_TOKEN_TTL, REFRESH_TOKEN_TTL

---

## Глава 4: Транзакции и распределенные системы (страницы 131-190)

### Transaction module (стр. 131-170)
- **Proto контракты**:
  - service TransactionService с методами Deposit, Withdraw, Transfer, GetTransactions
  - Сообщения: DepositRequest, WithdrawRequest, TransferRequest, TransactionDetails
- **Модели**:
  - Transaction: ID, UserID, AccountID, Amount, Status, Type, CreatedAt, UpdatedAt
  - TransactionEntry: запись в журнале транзакций
  - TransactionStatus: pending, completed, failed
  - TransactionType: deposit, withdraw, transfer
- **Repository**:
  - Deposit: создание транзакции с типом deposit
  - Withdraw: проверка баланса, создание транзакции
  - Transfer: списание с одного счета, зачисление на другой
  - UpdateTransactionStatus
- **Service layer**:
  - Публикация события в Kafka после создания транзакции (status=pending)
  - Обработка ответа от Account service (HandleAccountResponse)
  - Обновление статуса транзакции на completed/failed

### Нормализация БД и транзакции (стр. 171-190)
- **Нормальные формы**:
  - 1NF: атомарность значений
  - 2NF: устранение частичной зависимости
  - 3NF: устранение транзитивных зависимостей
  - 4NF: устранение многозначительных зависимостей
  - 5NF: устранение проекционных зависимостей
  - DKNF: зависимость только от доменных констант
  - 6NF: требования к временным типам
- **ACID свойства**:
  - Атомарность, Консистентность, Изолированность, Прочность
- **Уровни изоляции**:
  - Read Uncommitted → Read Committed → Repeatable Read → Serializable
  - Race condition: потерянное обновление (lost update)
- **Distributed Transactions**:
  - 2PC (Двухфазная фиксация): prepare → commit/rollback
  - **Saga pattern**: последовательность локальных транзакций с компенсирующими действиями

### Account service интеграция (стр. 191-210)
- **Account repository**:
  - updateBalance: изменение баланса
  - transferBalance: перевод между счетами
  - GetBalance: получение текущего баланса
- **Account service gRPC**:
  - Deposit, Withdraw, Transfer, GetBalance методы
  - Публикация BalanceChanged события в Kafka
- **Server handlers**: обработка gRPC запросов

---

## Глава 5: Интеграция Kafka и тестирование (страницы 211-250)

### Kafka Producer и Consumer (стр. 211-230)
- **Producer**:
  - `kafka.NewWriter()` с конфигом: брокеры, топик
  - `kafka.Message{Topic, Key, Value}` — сообщение
  - Публикация событий: TransactionSaved, BalanceChanged
- **Consumer**:
  - `kafka.NewReader()` с конфигом: GroupID, Brokers, Topic
  - Потребление сообщений в цикле (FetchMessage)
  - Обработка типов запросов: deposit, withdraw, transfer
- **Event-Driven Architecture**:
  - PUBLISH TransactionSaved → Account service обрабатывает
  - LISTEN BalanceChanged → Transaction service обновляет статус

### Transaction service с Kafka (стр. 231-250)
- **Service layer**:
  - Deposit/Withdraw/Transfer: создание транзакции → публикация в Kafka
  - Если публикация в Kafka не удалась → UpdateTransactionStatus(failed)
- **Kafka message handler**:
  - Маршалинг/анмаршалинг JSON для AccountResponse
  - Маршрутизация по RequestType (deposit, withdraw, transfer)
  - Публикация TransactionResponse в topic transaction_response
- **Dependency Injection**:
  - app.go: инициализация миграций, Kafka, gRPC подключений
  - Конфигурация через переменные окружения

### Тестирование (стр. 251-280)
- **Unit тесты с GoMock**:
  - Создание моков для Repository, AccountService, KafkaPublisher
  - gomock.NewController(t) + defer ctrl.Finish()
  - Table-driven tests с setupMocks
  - Test cases: успешные операции + ошибки (repository error, kafka publish error)
- **Coverage**: 84.4% для service слоя
- **CI/CD**:
  - GitHub Actions workflow с триггерами: `[workflow_dispatch, push]`
  - Шаги: checkout → go mod download → go test -v -race -coverprofile=coverage.out ./...

---

## Глава 6: Развертывание и деплой (страницы 281-322)

### Docker (стр. 281-300)
- **Dockerfile (мульти-стейдж сборка)**:
  ```dockerfile
  FROM golang:1.24.5-alpine AS build
  WORKDIR /app
  COPY go.mod go.sum ./
  RUN go mod download
  COPY ..
  RUN go build -o app ./cmd
  FROM alpine:3.18
  WORKDIR /app
  COPY --from=build /app/app
  COPY --from=build /app/internal/migrations ./internal/migrations
  CMD ["./app"]
  ```
- **Оптимизация Docker layers**:
  - Сначала копируются go.mod/go.sum → go mod download (кэшируется)
  - Затем копируется исходный код → go build
  - Если зависимости не менялись, слой кэшируется
- **.dockerignore**: .idea, vendor, .git
- **docker-compose.yaml**:
  - 3 PostgreSQL базы: account_db, auth_db, transaction_db
  - Healthcheck: `pg_isready -U postgres`
  - Kafka: confluentinc/cp-kafka:7.5.0 с KRaft режимом
  - Kafka UI: provectuslabs/kafka-ui:latest
  - 4 микросервиса: transaction, account, auth, gateway
  - app-network (bridge driver)

### Docker Hub (стр. 301-310)
- **Регистрация**: https://hub.docker.com/
- **Авторизация**: `docker login`
- **Тегирование**: `docker tag book_all-auth yuliapopova/book_all-auth:latest`
- **Публикация**: `docker push yuliapopova/book_all-auth:latest`
- **Обновление docker-compose**: замена `build:` на `image: yuliapopova/service:latest`
- **Структура проекта**:
  ```
  book-all/
  ├── docker-compose.yaml
  ├── account/
  ├── auth/
  ├── transaction/
  ├── gateway/
  └── contracts/
  ```

### Kubernetes (стр. 311-322)
- **Kubernetes**: оркестратор контейнеров с открытым исходным кодом
- **Ключевые понятия**:
  - **Node**: физическая/виртуальная машина в кластере
  - **Pod**: базовая единица — один или несколько контейнеров на одном узле
  - **Volume**: общий ресурс хранилища для контейнеров в поде
  - **Replication Controller**: масштабирование количества подов
  - **Service**: абстракция над группой подов + политика доступа
- **kubectl команды**:
  - `kubectl run nginx --image=nginx:latest --port=80` (создание пода)
  - `kubectl describe pods nginx` (информация о поде)
  - `kubectl apply -f pod.yaml` (применение манифеста)
  - `kubectl get pods` (список подов)
  - `kubectl logs <pod-name>` (просмотр логов)
- **Манифест Deployment**:
  ```yaml
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: my-nginx-deployment
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: nginx
    template:
      metadata:
        labels:
          app: nginx
      spec:
        containers:
        - name: nginx
          image: cr.yandex/crpv7tlcpgb30qpgkiij/ubuntu-nginx:latest
  ```
- **Kompose**: утилита для миграции docker-compose.yaml → Kubernetes манифесты
  - `kompose convert -f docker-compose.yaml`
  - Генерирует deployment и service для каждого сервиса
- **Заключение**:
  - Миграция Docker Compose → Kubernetes через kompose
  - Финальный цикл: Commit → CI/CD (GitHub Actions) → Docker image → Docker Hub → Pull на VM → Run (docker-compose или Kubernetes)

---

## Предметный указатель (ключевые понятия)
- **2PC (Двухфазная фиксация)**: 228
- **ACID**: 166
- **AMQP**: 120
- **API Gateway**: 126
- **AWS/GCP Cloud**: 294 (Yandex Cloud пример)
- **bcrypt**: 110
- **CI/CD**: 270
- **Docker**: 22, 294
- **Docker Hub**: 300
- **Docker-compose**: 22
- **Dockerfile**: 294
- **GORM**: 40
- **Git**: 19
- **git flow**: 20
- **GitHub Actions**: 272
- **GoMock**: 292
- **Golang**: 14
- **Go Modules**: 14
- **GORM**: 40
- **HTTP**: 111
- **HTTP/2**: 111
- **gRPC**: 116
- **JWT**: 92
- **Kafka**: 123
- **Keycloak**: 93
- **Kubernetes**: 267, 306
- **kompose**: 309
- **Ключ доступа (Access Key)**: 91
- **Микросервисы (микросервисная архитектура)**: 14-15
- **Миграции (Migrations)**: 161
- **Nginx**: 115, 268
- **Node Kubernetes**: 306
- **ORM**: 40
- **PAGINATION (Пагинация)**: 66
- **POSTGRES**: 47
- **PostgreSQL**: 22, 47
- **Protobuf**: 18
- **RACe condition**: 167
- **Read Committed**: 168
- **Read Uncommitted**: 168
- **refresh-token**: 109
- **Repeatable Read**: 168
- **RESПул API**: 43
- **Saga**: 229
- **Serializable**: 168
- **Service Discovery**: 14-15
- **СУБД (Базы данных)**: 47
- **Тестирование**: 25, 274-280
- **Транзакция**: 166
- **Workflow**: 272
- **ZeroLog**: 37

---

## Архитектурные слои проекта
1. **Domain layer** (`/internal/domain`): модели, интерфейсы репозиториев, бизнес-логика
2. **Repository layer** (`/internal/repository`): GORM-реализация интерфейсов доступа к данным
3. **Service layer** (`/internal/service`): оркестрация бизнес-операций, работа с событиями
4. **Transport layer** (`/internal/transport`): gRPC сервер, HTTP обработчики, interceptors
5. **Config layer** (`/internal/config`): конфигурация через переменные окружения
6. **Migrations** (`/internal/migrations`): SQL-файлы миграций (goose)
7. **Contracts** (`/contracts`): Proto-файлы (.proto) для gRPC

---

## Сервисы микросервисов
1. **Account service**: управление счетами и балансами
   - gRPC методы: Deposit, Withdraw, Transfer, GetBalance
   - Kafka: публикация BalanceChanged событий

2. **Auth service**: аутентификация и авторизация
   - gRPC методы: Register, Login, Logout, Refresh, GetCurrentUser
   - bcrypt + JWT (access/refresh tokens)

3. **Transaction service**: управление финансовыми транзакциями
   - gRPC методы: Deposit, Withdraw, Transfer, GetTransactions
   - Kafka: публикация TransactionSaved, обработка BalanceChanged

4. **Gateway service**: API Gateway
   - Единая точка входа
   - Маршрутизация к микросервисам
   - Аутентификация через JWT interceptor
