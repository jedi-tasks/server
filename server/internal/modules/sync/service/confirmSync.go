package service

import (
	"context"

	"server/internal/modules/sync/model"
)

// ConfirmSync принимает от устройства подтверждение, что оно корректно применило изменения,
// полученные последним вызовом Sync, и фиксирует новый чекпоинт устройства
func (s *SyncService) ConfirmSync(ctx context.Context, req model.ConfirmSyncReq) error {
	ctx, span := tracer.Start(ctx, "ConfirmSync")
	defer span.End()

	return s.userService.ConfirmDeviceSync(ctx, req.Necessary.UserID, req.Necessary.DeviceID, req.PendingSyncToken)
}
