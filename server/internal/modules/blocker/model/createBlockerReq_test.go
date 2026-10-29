package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// TestCreateBlockerReqValidate проверяет правило «ровно одно из blockingTaskID и reason» и формат даты
func TestCreateBlockerReqValidate(t *testing.T) {
	taskID, blockingTaskID := uuid.New(), uuid.New()
	reason := "ждем ответа"

	testCases := []struct {
		name           string     // Название сценария
		blockingTaskID *uuid.UUID // Блокирующая задача
		reason         *string    // Текстовая причина
		sinceDay       string     // Дата начала
		wantError      bool       // Ожидается отказ
	}{
		{name: "только задача", blockingTaskID: &blockingTaskID, reason: nil, sinceDay: "2026-10-29", wantError: false},
		{name: "только причина", blockingTaskID: nil, reason: &reason, sinceDay: "2026-10-29", wantError: false},
		{name: "ни того ни другого", blockingTaskID: nil, reason: nil, sinceDay: "2026-10-29", wantError: true},
		{name: "и то и другое", blockingTaskID: &blockingTaskID, reason: &reason, sinceDay: "2026-10-29", wantError: true},
		{name: "блокирует сама себя", blockingTaskID: &taskID, reason: nil, sinceDay: "2026-10-29", wantError: true},
		{name: "неверная дата", blockingTaskID: nil, reason: &reason, sinceDay: "29.10.2026", wantError: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := CreateBlockerReq{TaskID: taskID, BlockingTaskID: testCase.blockingTaskID, Reason: testCase.reason, SinceDay: testCase.sinceDay} //nolint:exhaustruct
			if err := request.Validate(context.Background()); (err != nil) != testCase.wantError {
				t.Fatalf("ожидался отказ: %v, получено: %v", testCase.wantError, err)
			}
		})
	}
}
