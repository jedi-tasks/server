package systemInfo

type SystemInfo struct {
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
	Build    string `json:"build"`
	Commit   string `json:"commit"`
	Env      string `json:"env"`
}

type UserInfo struct {
	UserID    *int64  `json:"userID"`
	DeviceID  *string `json:"deviceID"`
	RequestID *string `json:"requestID"`
}
