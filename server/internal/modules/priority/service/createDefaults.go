package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/priority/model"
	repoModel "server/internal/modules/priority/repository/model"
	"server/internal/utils/necessary"
)

// CreateDefaults создает приоритеты системы «Путь джедая» для нового пользователя внутри уже открытой транзакции
func (s *PriorityService) CreateDefaults(ctx context.Context, userID uuid.UUID, deviceID string, now time.Time) error {
	ctx, span := tracer.Start(ctx, "CreateDefaults")
	defer span.End()

	// Уровни по умолчанию: название, цвет, знаки
	defaults := []struct {
		name     string // Название уровня
		colorHex string // Цвет уровня
		marks    string // Знаки в строке задачи
	}{
		{name: "1 — ближайшие дни", colorHex: "#E5484D", marks: "!!!"},
		{name: "5 — неделя–две", colorHex: "#F5A524", marks: "!!"},
		{name: "9 — когда-нибудь", colorHex: "#3E9BF0", marks: "!"},
	}

	necessaryInfo := necessary.NecessaryUserInformation{UserID: userID, DeviceID: deviceID}
	for index, priority := range defaults {

		// Создаем уровень
		id := uuid.New()
		if _, err := s.priorityRepository.CreatePriority(ctx, repoModel.CreatePriorityReq{
			ID:        id,
			UserID:    userID,
			Name:      priority.name,
			ColorHex:  priority.colorHex,
			Marks:     priority.marks,
			SortOrder: float64(index),
			UpdatedAt: now,
		}); err != nil {
			return err
		}

		// Получаем созданный уровень для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.priorityRepository.GetPriorities(ctx, model.GetPrioritiesReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       []uuid.UUID{id},
		}))
		if err != nil {
			return err
		}

		// Фиксируем создание в аудит-логе
		if _, err = s.auditLogService.TrackMutation(ctx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Priority,
			Method:   auditLogMethod.Create,
			EntityID: id.String(),
			Before:   nil,
			After:    after,
			UserID:   userID,
			DeviceID: deviceID,
		}); err != nil {
			return err
		}
	}

	return nil
}
