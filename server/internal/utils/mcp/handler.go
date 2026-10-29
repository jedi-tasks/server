package mcp

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	sdkAuth "github.com/modelcontextprotocol/go-sdk/auth"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/oauthex"

	"server/internal/utils/contextKeys"
	"server/internal/utils/necessary"
)

// ResourceMetadataPath — путь метаданных защищенного ресурса (RFC 9728)
const ResourceMetadataPath = "/.well-known/oauth-protected-resource"

// NewHandler оборачивает MCP-сервер в HTTP: Bearer-токен проверяется как JWT сервера, без токена
// клиент получает 401 с адресом метаданных ресурса и начинает OAuth. Затем Streamable HTTP без
// сессий — каждый запрос независим, как и остальной сервер
func NewHandler(server *sdk.Server, publicURL string) http.Handler {
	transport := sdk.NewStreamableHTTPHandler(
		func(*http.Request) *sdk.Server { return server },
		&sdk.StreamableHTTPOptions{ //nolint:exhaustruct
			Stateless:    true,
			JSONResponse: true,
		},
	)

	// Пользователь и устройство из токена — в контекст, как после gRPC-интерсептора авторизации
	withNecessary := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := sdkAuth.TokenInfoFromContext(r.Context())
		userID, err := uuid.Parse(info.UserID)
		deviceID, _ := info.Extra[DeviceIDKey].(string)
		if err != nil || deviceID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := contextKeys.SetUserID(r.Context(), userID)
		ctx = contextKeys.SetDeviceID(ctx, deviceID)
		transport.ServeHTTP(w, r.WithContext(ctx))
	})

	requireToken := sdkAuth.RequireBearerToken(NewTokenVerifier(), &sdkAuth.RequireBearerTokenOptions{ //nolint:exhaustruct
		ResourceMetadataURL:    publicURL + ResourceMetadataPath,
		AllowMissingExpiration: true,
	})
	return requireToken(withNecessary)
}

// NewResourceMetadataHandler отдает метаданные защищенного ресурса: клиент узнает из них, где
// Authorization Server. У нас это тот же хост
func NewResourceMetadataHandler(publicURL string) http.Handler {
	return sdkAuth.ProtectedResourceMetadataHandler(&oauthex.ProtectedResourceMetadata{ //nolint:exhaustruct
		Resource:               publicURL,
		AuthorizationServers:   []string{publicURL},
		BearerMethodsSupported: []string{"header"},
		ResourceName:           "Tasks",
	})
}

// Necessary возвращает пользователя и устройство из контекста запроса для сервисных моделей
func Necessary(ctx context.Context) (necessary.NecessaryUserInformation, error) {
	var dest struct {
		Necessary necessary.NecessaryUserInformation
	}
	if err := necessary.ParseNecessary(ctx, &dest); err != nil {
		return dest.Necessary, err
	}
	return dest.Necessary, nil
}
