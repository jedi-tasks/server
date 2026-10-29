# Architecture & Code Standards

## Layered Architecture

Each domain module lives under `internal/modules/{module}/` and is split into exactly these layers:

```
internal/modules/{module}/
├── endpoint/grpc/     # gRPC handlers
├── service/           # Business logic
├── repository/        # Database access
├── network/           # External HTTP/gRPC/TCP calls
└── model/             # Data structures and DTOs
```

| Layer | Responsibility |
|---|---|
| **Endpoint** | Receive proto, extract context, convert to model, validate, call service |
| **Service** | Business logic, orchestrate repo + service calls, manage transactions |
| **Repository** | All DB/cache access, convert DB models → service models |
| **Network** | All outbound HTTP/gRPC/TCP calls |
| **Model** | Data structures, request/response DTOs |

### One File Per Function and Struct

Every function and every struct lives in its **own separate file** named after its contents:

```
service/
├── service.go       # Struct + constructor + all dependency interfaces
├── createSSP.go     # CreateSSP method only
└── findSSPs.go      # FindSSPs method only

model/
├── ssp.go           # SSP struct only
└── createSSPReq.go  # CreateSSPReq struct only
```

### Service Interaction Rules

A service may only communicate with its **own** repository, its **own** network layer, and **other services via interfaces**. A service must never bypass another service to call its network or repository directly.

```go
// Верно
ssps, err := s.sspService.GetSSPsFromCache(ctx, req)

// Неверно — нельзя лезть в чужой репозиторий напрямую
ssps, err := s.sspRepository.FindSSPs(ctx, req)
```

### Server Makes No Hidden Changes

The server is only the central relay: a mutation changes exactly the rows the client sent, nothing else. No cascades (no deleting children, detaching references or resolving blockers on the server). Instead, a mutation is rejected with `errors.BadRequest` and a Russian `WithCustomHumanText` if afterwards a live row (`deleted_at IS NULL`) would reference a soft-deleted one. Checks run inside the mutation's transaction: a delete checks for live children after its UPDATE (the parent row is locked), a create/update locks referenced rows with `LockLive` (`FOR SHARE`). Rules and texts are listed in `README.md` («Правила целостности») — keep them in sync.

---

## Endpoint Layer

Endpoints use compile-time interface checks, convert proto → model, validate, then call the service.

```go
// Проверка реализации контракта на уровне компилятора
var _ SSPService = new(sspService.SSPService)

// SSPService — интерфейс бизнес-логики SSP
type SSPService interface {
	CreateSSP(context.Context, model.CreateSSPReq) error
}

// SSPEndpoint — gRPC-обработчик для SSP
type SSPEndpoint struct {
	proto.UnsafeSSPEndpointServer
	sspService SSPService
}

// CreateSSP обрабатывает gRPC-запрос на создание SSP
func (e *SSPEndpoint) CreateSSP(ctx context.Context, r *proto.CreateSSPRequest) (*proto.CreateSSPResponse, error) {
	res := new(proto.CreateSSPResponse)

	// Маппим proto-структуру во внутреннюю модель
	req, err := model.ProtoCreateSSPReq{CreateSSPRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Валидируем входные данные
	if err = validator.Validate(req); err != nil {
		return res, err
	}

	// Создаём SSP через сервис
	return res, e.sspService.CreateSSP(ctx, req)
}
```

---

## Service Layer

All dependencies are declared as interfaces with compile-time checks. Never accept concrete types.

```go
// SSPService — сервис управления SSP
type SSPService struct {
	sspRepository SSPRepository
	transactor    Transactor
	auditLog      AuditLog
}

// Проверка реализации контракта репозитория
var _ SSPRepository = new(sspRepository.SSPRepository)

// SSPRepository — интерфейс репозитория SSP
type SSPRepository interface {
	CreateSSP(context.Context, model.CreateSSPReq) error
	FindSSPs(context.Context, model.FindSSPsReq) ([]model.SSP, error)
}

// CreateSSP создаёт SSP и фиксирует событие в аудит-логе
func (s *SSPService) CreateSSP(ctx context.Context, req model.CreateSSPReq) error {
	return s.transactor.WithinTransaction(ctx, func(tx context.Context) error {

		// Создаём запись в базе данных
		if err := s.sspRepository.CreateSSP(tx, req); err != nil {
			return err
		}

		// Фиксируем создание в аудит-логе через другой сервис
		return s.auditLog.TrackMutation(tx, auditLogModel.ActionCreate, "CreateSSP", "SSP", req.Slug, nil, nil)
	})
}
```

---

## Repository Layer

All queries use **Squirrel** (`github.com/Masterminds/squirrel`). Table/column names come from DDL constants in `internal/ddl/pgsqlDDL/`. Repository models are separate from service models and converted via `ConvertToModel()`.

```go
// FindSSPs возвращает список SSP по фильтрам, сортировке и пагинации
func (r *SSPRepository) FindSSPs(ctx context.Context, req sspModel.FindSSPsReq) ([]sspModel.SSP, error) {

	// Формируем базовый SELECT-запрос
	q := sq.Select(ddlHelper.SelectAll).From(sspDDL.Table)

	// Применяем фильтры
	q = GetSSPsFilters(req.Filters, q)

	// Применяем пагинацию
	if req.Pagination.Size != 0 {
		q = q.Offset(uint64((req.Pagination.Page-1)*req.Pagination.Size)).
			Limit(uint64(req.Pagination.Size))
	}

	// Выполняем запрос
	var rows []sspRepoModel.SSP
	if err := r.pgsql.Select(ctx, &rows, q); err != nil {
		return nil, err
	}

	// Конвертируем в сервисные модели
	return slices.Map(rows, func(s sspRepoModel.SSP) sspModel.SSP { return s.ConvertToModel() }), nil
}

// GetSSPsFilters применяет фильтры к запросу выборки SSP
func GetSSPsFilters(filters sspFilters.SSPFilters, q squirrel.SelectBuilder) squirrel.SelectBuilder {

	// Фильтрация по slug
	if len(filters.Slugs) > 0 {
		q = q.Where(squirrel.Eq{sspDDL.ColumnSlug: filters.Slugs})
	}

	// Исключаем мягко удалённые записи
	return q.Where(squirrel.Eq{sspDDL.ColumnIsDeleted: false})
}
```

DDL constants in `internal/ddl/pgsqlDDL/{entity}DDL/`:

```go
const (
	Table           = "ssps"
	ColumnSlug      = "slug"       // Системное имя
	ColumnIsEnable  = "is_enable"  // Признак активности
	ColumnIsDeleted = "is_deleted" // Признак мягкого удаления
)
```

---

## Model Layer

Service model and repository model (with `db` tags) are separate structs. The repo model implements `ConvertToModel()`.

```go
// SSP — внутренняя модель SSP-платформы
type SSP struct {
	Slug        string // Системное имя (уникальный идентификатор)
	Name        string // Отображаемое название
	IsEnable    bool   // Признак активности для входящего трафика
	Timeout     *int32 // Таймаут ответа в миллисекундах (nil = без ограничений)
	Currency    string // Валюта расчётов
	SupportGzip bool   // Поддержка gzip-сжатия в ответах
}
```

```go
// SSP — репозиторная модель, соответствует строке в таблице ssps
type SSP struct {
	Slug        string `db:"slug"`          // Системное имя
	Name        string `db:"name"`          // Отображаемое название
	IsEnable    bool   `db:"is_enable"`     // Признак активности
	Timeout     *int32 `db:"timeout"`       // Таймаут в миллисекундах
	Currency    string `db:"currency"`      // Валюта
	SupportGzip bool   `db:"supports_gzip"` // Поддержка gzip
}

// ConvertToModel конвертирует репозиторную модель в сервисную
func (s SSP) ConvertToModel() model.SSP {
	return model.SSP{Slug: s.Slug, Name: s.Name, IsEnable: s.IsEnable,
		Timeout: s.Timeout, Currency: s.Currency, SupportGzip: s.SupportGzip}
}
```

---

## Error Handling

Never use `fmt.Errorf` or standard `errors.New`. Always use the custom error type from `internal/utils/errors`.

```go
// Создание новой ошибки
return errors.InternalServer.New("не удалось подключиться к базе данных")

// Оборачивание существующей ошибки
return errors.BadRequest.Wrap(err)

// С параметрами для отладки
return errors.InternalServer.New("ошибка запроса").
    WithParams("sspSlug", slug, "requestID", id).
    WithContextParams(ctx)

// С кастомным текстом для пользователя
return errors.NotFound.New("ssp not found").
    WithCustomHumanText("SSP с таким slug не существует")
```

Error types are defined in `internal/utils/errors/errorTypes.go`:

```go
var InternalServer = errors.ErrorType{Name: "InternalServer", HTTPCode: http.StatusInternalServerError,
    LogAs: errors.LogAsError, HumanText: "Произошла непредвиденная ошибка"}

var BadRequest = errors.ErrorType{Name: "BadRequest", HTTPCode: http.StatusBadRequest,
    LogAs: errors.LogAsWarning, HumanText: "Введены неверные данные"}

var NotFound = errors.ErrorType{Name: "NotFound", HTTPCode: http.StatusNotFound,
    LogAs: errors.LogAsWarning, HumanText: "Данные не найдены"}

var Unauthorized = errors.ErrorType{Name: "Unauthorized", HTTPCode: http.StatusUnauthorized,
    LogAs: errors.LogAsWarning, HumanText: "Пользователь не авторизован"}
```

---

## Naming Rules

Always use full, unabbreviated names for variables, fields, parameters, and identifiers. Never shorten names.

```go
// Верно
userService UserService
sspRepository SSPRepository

// Неверно — сокращения запрещены
userSvc UserService
sspRepo SSPRepository
```

---

## Comment Rules

1. **Every function** — one Russian comment on the line directly above `func`.
2. **Every struct** — one Russian comment on the line directly above `type`.
3. **Every struct field** — inline Russian comment after the type (or db tag).
4. **Every meaningful code block** inside a function — one Russian comment above the block.
