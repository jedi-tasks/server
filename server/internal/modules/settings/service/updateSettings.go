package service

import (
	"context"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/settings/model"
	"server/internal/utils/errors"
)

// UpdateSettings меняет переданные поля настроек. Изменение старее сохраненного игнорируется
func (s *SettingsService) UpdateSettings(ctx context.Context, req model.UpdateSettingsReq) error {
	ctx, span := tracer.Start(ctx, "UpdateSettings")
	defer span.End()

	// Проверяем, что есть поля для обновления
	if req.DailyCapacityMinutes == nil && req.TodaySortMode == nil {
		return errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Получаем настройки для слепка "до" в аудит-логе
		before, err := s.GetSettings(ctxTx, req.Necessary.UserID)
		if err != nil {
			return 0, err
		}

		// Применяем изменения
		after := before
		if req.DailyCapacityMinutes != nil {
			after.DailyCapacityMinutes = *req.DailyCapacityMinutes
		}
		if req.TodaySortMode != nil {
			after.TodaySortMode = *req.TodaySortMode
		}
		after.UpdatedAt = req.UpdatedAt

		// Сохраняем настройки
		saved, err := s.settingsRepository.SaveSettings(ctxTx, after)
		if err != nil || !saved {
			return 0, err
		}

		// Фиксируем изменение в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Settings,
			Method:   auditLogMethod.Update,
			EntityID: req.Necessary.UserID.String(),
			Before:   before,
			After:    after,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}
