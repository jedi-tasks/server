package model

type UpdateDeviceInformationReq struct {
	VersionOS *string
	IPAddress *string
	UserAgent *string
	Stateless *bool // Устройство без локальной копии данных
}
