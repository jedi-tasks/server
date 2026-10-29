package model

import "github.com/google/uuid"

type GetDevicesReq struct {
	IDs       []uuid.UUID
	DeviceIDs []string
	UserIDs   []uuid.UUID
}
