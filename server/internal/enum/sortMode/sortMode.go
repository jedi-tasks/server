package sortMode

import (
	"context"

	"pkg/maps"

	"server/internal/utils/errors"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

// SortMode - режим сортировки задач в списке
type SortMode string

// enum:"manual,priority,plannedDay,deadline,duration"
const (
	Manual     = SortMode("manual")
	Priority   = SortMode("priority")
	PlannedDay = SortMode("plannedDay")
	Deadline   = SortMode("deadline")
	Duration   = SortMode("duration")
)

// Validate проверяет, что режим сортировки принадлежит известному набору значений
func (m SortMode) Validate(ctx context.Context) error {
	switch m {
	case Manual, Priority, PlannedDay, Deadline, Duration:
	default:
		return errors.BadRequest.New("Unknown sort mode").
			WithContextParams(ctx).
			WithParams("sortMode", m)
	}
	return nil
}

// mappingProtoToModel содержит соответствие между значениями proto.SortMode и SortMode
var mappingProtoToModel = map[proto.SortMode]SortMode{
	proto.SortMode_Manual:     Manual,
	proto.SortMode_Priority:   Priority,
	proto.SortMode_PlannedDay: PlannedDay,
	proto.SortMode_Deadline:   Deadline,
	proto.SortMode_Duration:   Duration,
}

// ConvertToProto преобразует SortMode в proto.SortMode
func (m SortMode) ConvertToProto() (proto.SortMode, error) {

	// Разворачиваем мапу
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	// Получаем значение
	protoSortMode, ok := mappingModelToProto[m]
	if !ok {
		return protoSortMode, errors.BadRequest.New("SortMode undefined")
	}

	return protoSortMode, nil
}

type ProtoSortMode struct {
	proto.SortMode
}

// ConvertToModel преобразует ProtoSortMode в SortMode
func (p ProtoSortMode) ConvertToModel() (SortMode, error) {

	// Проверяем наличие значения
	sortMode, ok := mappingProtoToModel[p.SortMode]
	if !ok {
		return sortMode, errors.BadRequest.New("SortMode undefined")
	}

	return sortMode, nil
}
