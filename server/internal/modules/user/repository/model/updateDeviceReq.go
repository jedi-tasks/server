package model

import "github.com/google/uuid"

type UpdateDeviceReq struct {
	UserID                 uuid.UUID
	DeviceID               string
	RefreshToken           *string
	NotificationToken      *string
	ApplicationInformation UpdateApplicationInformationReq
	DeviceInformation      UpdateDeviceInformationReq
	LastSyncedAuditLogID   *uint32
}
