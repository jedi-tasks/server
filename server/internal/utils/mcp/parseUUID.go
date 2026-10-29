package mcp

import (
	"fmt"

	"github.com/google/uuid"
)

// ParseUUID разбирает обязательный идентификатор из аргумента инструмента
func ParseUUID(field, value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: ожидается UUID, получено %q", field, value)
	}
	return id, nil
}

// ParseOptionalUUID разбирает необязательный идентификатор: nil-аргумент — nil, пустая строка —
// uuid.Nil (сброс значения там, где сервис так его трактует)
func ParseOptionalUUID(field string, value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return &uuid.Nil, nil
	}
	id, err := ParseUUID(field, *value)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// ParseUUIDs разбирает список идентификаторов
func ParseUUIDs(field string, values []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		id, err := ParseUUID(field, value)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// NewOrParseUUID возвращает идентификатор из аргумента или новый, если агент его не передал
func NewOrParseUUID(field string, value *string) (uuid.UUID, error) {
	if value == nil || *value == "" {
		return uuid.New(), nil
	}
	return ParseUUID(field, *value)
}

// UUIDString — указатель на UUID строкой для ответов, nil остается nil
func UUIDString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}
