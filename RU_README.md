# Сервис отслеживания отправлений (Shipment Tracking Service)

---

## Как запустить сервис

### 1. Подготовка окружения

Убедитесь, что у вас установлены **Docker** и **Docker Compose**.
Создайте файл `.env` в корне проекта (или проверьте существующий):

```env
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret_password
POSTGRES_DB=shipment_db
GRPC_PORT=50051
LOG_LEVEL=DEBUG
```

### 2. Запуск через Docker Compose (рекомендуется)

Выполните команду для сборки и запуска всех компонентов:

```bash
docker compose up -d --build
```

Это автоматически:

1. Поднимет базу данных **PostgreSQL**.
2. Запустит сервис **migrator**, который применит все SQL-миграции.
3. Соберет и запустит основной **gRPC сервер**.

### 3. Ручной запуск (для разработки)

Если вы хотите запустить сервер локально:

1. Поднимите БД: `docker compose up -d postgres`.
2. Загрузите зависимости: `go mod download`.
3. Запустите сервер: `go run cmd/main.go`.

---

## Как запустить тесты

### 1. Юнит-тесты и тесты бизнес-логики

Проект покрыт тестами на уровне домена и прикладного слоя. Для запуска выполните:

```bash
go test ./...
```

## Обзор архитектуры

Сервис построен по принципам **чистой архитектуры** (Clean Architecture) и с применением **DDD** (Domain-Driven Design). Всего реализовано 4 слоя: Presentation, Application, Domain и Infrastructure.
- **Presentation Layer**: gRPC сервер, который обрабатывает входящие запросы и возвращает ответы.
- **Application Layer**: Это промежуточный слой, который связывает бизнес-логику с внешним миром.
- **Domain Layer**: Содержит бизнес-логику, модели и спецификации. Здесь реализованы все правила и ограничения для управления отправлениями.
- **Infrastructure Layer**: Реализация доступа к данным (репозиторий)

### Слои взаимодействуют друг с другом через интерфейсы

Интерфейс Репозитория:
```
type ShipmentRepository interface {
	Create(ctx context.Context, shipment domain.Shipment) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Shipment, error)
	AddEvent(ctx context.Context, shipmentID uuid.UUID, event kernel.DomainEvent) error
	GetHistory(ctx context.Context, shipmentID uuid.UUID) ([]EventDTO, error)
	UpdateShipmentStatus(ctx context.Context, shipmentID uuid.UUID, newStatus domain.Status) error
}
```

Интерфейс Сервиса:
```
type ShipmentService interface {
	UpdateShipmentStatus(context.Context, string, string) (domain.Shipment, error)
	CreateShipment(context.Context, string, string, domain.Details, domain.DriverDetails) (domain.Shipment, error)
	History(context.Context, string) ([]EventDTO, error)
	GetShipment(context.Context, string) (domain.Shipment, error)
}
```


## Паттерны проектирования:
Для валидации бизнес-правил я использую паттерн "Спецификация" (Specification), который позволяет инкапсулировать сложные правила валидации в отдельные объекты. A также позволяет легко добавлять новые правила без изменения существующего кода.

```type ShipmentStatusSpec interface {
    Check(shipment domain.Shipment, newStatus domain.Status) (bool, error)
}
```

также для решения проблемы валидации состояния я использовал матрицу смежности, которая определяет допустимые переходы между статусами отправления. поскольку валидация статусов может быть сложной, я решил инкапсулировать эту логику в отдельной структуре, которая использует эту матрицу для проверки допустимости переходов.
(matrix of allowed status transitions){./image.png}
поскольку колличество правил переходов зависит от колличества статусов.
```go
func DefaultTransitionSpec() (*transitionValidation, error) {
	return NewTransitionSpec(
		WithRule(domain.StatusPending, domain.StatusInTransit, domain.AlwaysAllowRule),
		WithRule(domain.StatusPending, domain.StatusCancelled, domain.AlwaysAllowRule),
		WithRule(domain.StatusPending, domain.StatusDelivered, domain.AlwaysDenyRule),

		...

		WithRule(domain.StatusCancelled, domain.StatusPending, domain.AlwaysDenyRule),
		WithRule(domain.StatusCancelled, domain.StatusInTransit, domain.AlwaysDenyRule),
		WithRule(domain.StatusCancelled, domain.StatusDelivered, domain.AlwaysDenyRule),
	)
}
```
Я использую Фабрику для создания экземпляра спецификации, которая инициализирует эту матрицу на основе текущих статусов. и следит за количеством правил, чтобы не допустить undefined behavior при переходе в несуществующий статус. или при добавлении нового статуса, который не был учтен в матрице.

```go
type transitionValidation struct {
	registry map[domain.Status]map[domain.Status]domain.Rule
}

type TransitionValidationOption func(*transitionValidation)

func NewTransitionSpec(opts ...TransitionValidationOption) (*transitionValidation,error){

	ts := &transitionValidation{
		registry: make(map[domain.Status]map[domain.Status]domain.Rule),
	}

	for _, opt := range opts {
		opt(ts)
	}

...

	return ts, nil
}

```

Также я использую DDD паттерны как Value-Object Агрегат и Enitiy, для разделния бизнес сущностей. Например, Shipment является агрегатом, который содержит в себе сущности Event и Value-Object Details и DriverDetails. Это позволяет мне четко разделить ответственность между различными частями модели и обеспечить целостность данных внутри агрегата.

```go
type AggregateRoot struct {
	domainEvents []DomainEvent
}

func (ar *AggregateRoot) ApplyDomain(e DomainEvent) {
	ar.domainEvents = append(ar.domainEvents, e)
}

func (ar *AggregateRoot) DomainEvents() []DomainEvent {
	return ar.domainEvents
}

```
Агрегат в моей модели также сущность которая производит события, которые описывают изменения в состоянии отправления. Это позволяет объявлять события в бизнес-логике и обрабатывать их в инфраструктуре, например, для сохранения в базу данных или отправки уведомлений.
Например так я создаю события:
```go
	shipment.ApplyDomain(domain.ShipmentStatusUpdatedEvent{
		ShipmentID: shipment.ID.String(),
		NewStatus:  newStatus,
	})
```
А так я обрабатываю события в application слое:
```go
    for _, event := range shipment.DomainEvents() {
        err = s.EventBus.Publish(ctx, EventBusKey, event)
		if err != nil {
			return domain.Shipment{}, err
		}
    }
```
Чтобы продемонстировать использование событий, я реализовал простой EventBus, который позволяет публиковать события и подписываться на них.
```go
type EventHandler interface {
	Handle(ctx context.Context, event kernel.DomainEvent) error
}

type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}
```
как пример подписчика на события я реализовал простой логгер, который выводит информацию о каждом событии в консоль.
```go
func (l *LogSubscriber) Handle(ctx context.Context, event kernel.DomainEvent) error {
	fmt.Printf("Event received: %s, Payload: %s\n", event.Name(), event.Payload())
	return nil
}
```
Для сборки всего приложения я сделал простой DI контейнер, который позволяет регистрировать зависимости и легко получать их в нужных местах. Также для удобства сборки я использую Фабрику, для управления зависимостями
```go
type Deps struct {
	Repository app.ShipmentRepository
	AppService app.ShipmentService

	EventBus *app.EventBus
	Pool     *pgxpool.Pool
	Server   *grpc.Server
}
```

Для SQL-запросов я использую sqlc — это генератор кода, который позволяет писать типобезопасный SQL. Такой подход гораздо надёжнее обычных запросов и работает быстрее, чем ORM.

А для генерации gRPC кода я использую buf — это инструмент для управления Protocol Buffers.

