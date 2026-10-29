package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"server/internal/enum/sortMode"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// DefaultSettings — настройки пользователя, который их еще не менял
func DefaultSettings(userID uuid.UUID) Settings {
	return Settings{
		UserID:               userID,
		DailyCapacityMinutes: 360,
		TodaySortMode:        sortMode.Priority,
		UpdatedAt:            time.Time{},
	}
}

// Settings — настройки пользователя, общие для всех его устройств
type Settings struct {
	UserID               uuid.UUID         `json:"userID" db:"user_id"`                              // Идентификатор пользователя
	DailyCapacityMinutes uint32            `json:"dailyCapacityMinutes" db:"daily_capacity_minutes"` // Дневной лимит загрузки в минутах
	TodaySortMode        sortMode.SortMode `json:"todaySortMode" db:"today_sort_mode"`               // Режим сортировки экрана «Сегодня»
	UpdatedAt            time.Time         `json:"updatedAt" db:"updated_at"`                        // Дата и время последнего изменения
}

// ConvertToProto преобразует Settings в proto-формат
func (m Settings) ConvertToProto() (*proto.Settings, error) {
	todaySortMode, err := m.TodaySortMode.ConvertToProto()
	if err != nil {
		return nil, err
	}

	return &proto.Settings{
		DailyCapacityMinutes: m.DailyCapacityMinutes,
		TodaySortMode:        todaySortMode,
		UpdatedAt:            timestamppb.New(m.UpdatedAt),
	}, nil
}
