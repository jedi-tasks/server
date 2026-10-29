package integrationTest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"time"

	"pkg/jwtManager"
	taskEndpointMCP "server/internal/modules/task/endpoint/mcp"
	"server/internal/utils/auth"
	mcpTransport "server/internal/utils/mcp"
)

// mcpHTTPServer поднимает MCP по HTTP с проверкой Bearer-JWT, как в main.go
func (f *fixture) mcpHTTPServer() *httptest.Server {
	f.t.Helper()
	jwtManager.Init([]byte("test-signing-key"), 30*time.Minute, 720*time.Hour)
	server := mcpTransport.NewServer("test")
	taskEndpointMCP.Register(server, f.task)

	mux := http.NewServeMux()
	mux.Handle("GET "+mcpTransport.ResourceMetadataPath, mcpTransport.NewResourceMetadataHandler(oauthPublicURL))
	mux.Handle("/", mcpTransport.NewHandler(server, oauthPublicURL))
	httpServer := httptest.NewServer(mux)
	f.t.Cleanup(httpServer.Close)
	return httpServer
}

// mcpPost шлет JSON-RPC запрос на MCP с заголовком авторизации
func mcpPost(url, authorization, body string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	return http.DefaultClient.Do(req)
}

func TestMCP_HTTPRequiresOAuthToken(t *testing.T) {
	f := newFixture(t)
	server := f.mcpHTTPServer()
	toolsList := `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`

	// Без токена — 401 с адресом метаданных ресурса, с него клиент начинает OAuth
	res, err := mcpPost(server.URL+"/", "", toolsList)
	f.ok(err)
	challenge := res.Header.Get("WWW-Authenticate")
	if res.StatusCode != http.StatusUnauthorized || !strings.Contains(challenge, `resource_metadata="`+oauthPublicURL+mcpTransport.ResourceMetadataPath+`"`) {
		t.Fatalf("без токена: %d %q", res.StatusCode, challenge)
	}

	// Метаданные ресурса указывают на наш Authorization Server
	res, err = http.Get(server.URL + mcpTransport.ResourceMetadataPath)
	f.ok(err)
	var metadata struct {
		Resource             string   `json:"resource"`
		AuthorizationServers []string `json:"authorization_servers"`
	}
	f.ok(json.NewDecoder(res.Body).Decode(&metadata))
	if metadata.Resource != oauthPublicURL || len(metadata.AuthorizationServers) != 1 || metadata.AuthorizationServers[0] != oauthPublicURL {
		t.Fatalf("метаданные: %+v", metadata)
	}

	// Чужой токен — 401
	res, err = mcpPost(server.URL+"/", "Bearer not-a-jwt", toolsList)
	f.ok(err)
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("мусорный токен: %d", res.StatusCode)
	}

	// Refresh-токен вместо access — 401
	refresh, err := jwtManager.GenerateToken(f.ctx, jwtManager.RefreshToken, auth.Claims{UserID: f.necessary.UserID, DeviceID: f.necessary.DeviceID})
	f.ok(err)
	res, err = mcpPost(server.URL+"/", "Bearer "+refresh, toolsList)
	f.ok(err)
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("refresh вместо access: %d", res.StatusCode)
	}

	// Access-JWT сервера — инструменты доступны, правка идет от пользователя из токена
	access, err := jwtManager.GenerateToken(f.ctx, jwtManager.AccessToken, auth.Claims{UserID: f.necessary.UserID, DeviceID: f.necessary.DeviceID})
	f.ok(err)
	res, err = mcpPost(server.URL+"/", "Bearer "+access, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"create_task","arguments":{"title":"Через HTTP"}}}`)
	f.ok(err)
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK || !strings.Contains(string(body), `"id"`) || strings.Contains(string(body), `"isError":true`) {
		t.Fatalf("create_task по HTTP: %d %s", res.StatusCode, body)
	}
	if n := f.count(squirrelCountUserTasks(f, "Через HTTP")); n != 1 {
		t.Fatalf("задача не создана у пользователя из токена: %d", n)
	}
}

// squirrelCountUserTasks — счетчик живых задач пользователя фикстуры с данным названием
func squirrelCountUserTasks(f *fixture, title string) sq.SelectBuilder {
	return sq.Select("count(*)").From("tasks.tasks").Where(sq.Eq{"user_id": f.necessary.UserID, "title": title, "deleted_at": nil})
}
