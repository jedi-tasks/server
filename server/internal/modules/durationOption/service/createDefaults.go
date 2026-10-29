package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/durationOption/model"
	repoModel "server/internal/modules/durationOption/repository/model"
	"server/internal/utils/necessary"
)

// CreateDefaults создает справочник длительностей для нового пользователя внутри уже открытой транзакции
func (s *DurationOptionService) CreateDefaults(ctx context.Context, userID uuid.UUID, deviceID string, now time.Time) error {
	ctx, span := tracer.Start(ctx, "CreateDefaults")
	defer span.End()

	// Значения по умолчанию, как в приложении
	defaults := []struct {
		title   string // Подпись
		minutes int32  // Минуты, меньше нуля — непрогнозируемо долго
	}{
		{title: "5 мин", minutes: 5},
		{title: "15 мин", minutes: 15},
		{title: "30 мин", minutes: 30},
		{title: "1 ч", minutes: 60},
		{title: "2 ч", minutes: 120},
		{title: "5 ч", minutes: 300},
		{title: "Непрогнозируемо долго", minutes: -1},
	}

	necessaryInfo := necessary.NecessaryUserInformation{UserID: userID, DeviceID: deviceID}
	for _, option := range defaults {

		// Создаем значение
		id := uuid.New()
		if _, err := s.durationOptionRepository.CreateDurationOption(ctx, repoModel.CreateDurationOptionReq{
			ID:        id,
			UserID:    userID,
			Title:     option.title,
			Minutes:   option.minutes,
			UpdatedAt: now,
		}); err != nil {
			return err
		}

		// Получаем созданное значение для слепка "после" в аудит-логе
		after, err := slices.FirstWithError(s.durationOptionRepository.GetDurationOptions(ctx, model.GetDurationOptionsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       []uuid.UUID{id},
		}))
		if err != nil {
			return err
		}

		// Фиксируем создание в аудит-логе
		if _, err = s.auditLogService.TrackMutation(ctx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.DurationOption,
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
