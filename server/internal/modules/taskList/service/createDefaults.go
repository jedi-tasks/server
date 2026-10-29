package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	"server/internal/enum/sortMode"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/taskList/model"
	repoModel "server/internal/modules/taskList/repository/model"
	"server/internal/utils/necessary"
)

// CreateDefaults создает системный Inbox нового пользователя внутри уже открытой транзакции
func (s *TaskListService) CreateDefaults(ctx context.Context, userID uuid.UUID, deviceID string, now time.Time) error {
	ctx, span := tracer.Start(ctx, "CreateDefaults")
	defer span.End()

	// Создаем Inbox
	inboxID := uuid.New()
	if err := s.taskListRepository.CreateInbox(ctx, repoModel.CreateTaskListReq{
		ID:        inboxID,
		UserID:    userID,
		Name:      "Inbox",
		ColorHex:  "#3E9BF0",
		Icon:      "tray.fill",
		SortOrder: 0,
		SortMode:  sortMode.Manual,
		CreatedAt: now,
	}); err != nil {
		return err
	}

	// Получаем созданный Inbox для слепка "после" в аудит-логе
	after, err := slices.FirstWithError(s.taskListRepository.GetTaskLists(ctx, model.GetTaskListsReq{ //nolint:exhaustruct
		Necessary: necessary.NecessaryUserInformation{UserID: userID, DeviceID: deviceID},
		IDs:       []uuid.UUID{inboxID},
	}))
	if err != nil {
		return err
	}

	// Фиксируем создание в аудит-логе
	_, err = s.auditLogService.TrackMutation(ctx, auditLogModel.TrackMutationReq{
		Entity:   auditLogEntity.TaskList,
		Method:   auditLogMethod.Create,
		EntityID: inboxID.String(),
		Before:   nil,
		After:    after,
		UserID:   userID,
		DeviceID: deviceID,
	})
	return err
}
