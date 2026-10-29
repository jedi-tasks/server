package mcp

import (
	"net/http"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	pkgErrors "pkg/errors"
	"pkg/log"
)

// NewHandler оборачивает MCP-сервер в HTTP: проверка Bearer-токена, затем Streamable HTTP без
// сессий — каждый запрос независим, как и остальной сервер
func NewHandler(server *sdk.Server, authenticator *Authenticator) http.Handler {
	transport := sdk.NewStreamableHTTPHandler(
		func(*http.Request) *sdk.Server { return server },
		&sdk.StreamableHTTPOptions{ //nolint:exhaustruct
			Stateless:    true,
			JSONResponse: true,
		},
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := authenticator.Authorize(r.Context(), r.Header.Get("Authorization"))
		if err != nil {
			log.Warning(pkgErrors.CastError(err))
			w.Header().Set("WWW-Authenticate", `Bearer realm="tasks"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		transport.ServeHTTP(w, r.WithContext(ctx))
	})
}
