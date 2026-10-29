package config

import (
	"github.com/jedi-tasks/go-server-grpc/proto"
)

func GetOpenForEveryoneMethods() []string {
	return []string{
		proto.AuthEndpoint_SignIn_FullMethodName,
		proto.AuthEndpoint_SignUp_FullMethodName,
		proto.AuthEndpoint_RefreshTokens_FullMethodName,
		proto.SettingsEndpoint_GetVersion_FullMethodName,
	}
}
