package integrationTest

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"time"

	"pkg/jwtManager"
	authModel "server/internal/modules/auth/model"
	authService "server/internal/modules/auth/service"
	oauthEndpointHTTP "server/internal/modules/oauth/endpoint/http"
	oauthRepository "server/internal/modules/oauth/repository"
	oauthService "server/internal/modules/oauth/service"
	userRepository "server/internal/modules/user/repository"
	"server/internal/utils/auth"
)

const oauthPublicURL = "https://mcp.test"

// oauthServer поднимает Authorization Server поверх сервисов фикстуры и заводит пользователя с паролем
func (f *fixture) oauthServer(email, password string) *httptest.Server {
	f.t.Helper()
	jwtManager.Init([]byte("test-signing-key"), 30*time.Minute, 720*time.Hour)

	authSvc := authService.NewAuthService(userRepository.NewUserRepository(database), f.transactor, []byte("salt"), false, f.auditLog,
		[]authService.DefaultsCreator{f.taskList, f.priority})
	if _, err := authSvc.EnsureUser(f.ctx, authModel.EnsureUserReq{Name: "OAuth", Email: email, Password: password}); err != nil {
		f.t.Fatal(err)
	}

	service := oauthService.NewOAuthService(oauthRepository.NewOAuthRepository(database), authSvc,
		oauthService.Config{PublicURL: oauthPublicURL, AccessTokenTTL: 30 * time.Minute, Version: "test"})
	mux := http.NewServeMux()
	oauthEndpointHTTP.NewOAuthEndpoint(service).Register(mux)
	server := httptest.NewServer(mux)
	f.t.Cleanup(server.Close)
	return server
}

// noRedirect — клиент, который не ходит по редиректам: нам нужен сам Location
var noRedirect = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

func TestOAuth_AuthorizationCodeFlow(t *testing.T) {
	f := newFixture(t)
	email := f.necessary.UserID.String() + "-oauth@test.local"
	server := f.oauthServer(email, "secret-password")

	// Метаданные
	res, err := http.Get(server.URL + "/.well-known/oauth-authorization-server")
	f.ok(err)
	var metadata map[string]any
	f.ok(json.NewDecoder(res.Body).Decode(&metadata))
	if metadata["issuer"] != oauthPublicURL || metadata["token_endpoint"] != oauthPublicURL+"/oauth/token" {
		t.Fatalf("метаданные: %v", metadata)
	}

	// Регистрация клиента
	res, err = http.Post(server.URL+"/oauth/register", "application/json",
		strings.NewReader(`{"client_name":"Claude Test","redirect_uris":["http://localhost:9/cb"]}`))
	f.ok(err)
	var client struct {
		ClientID string `json:"client_id"`
	}
	f.ok(json.NewDecoder(res.Body).Decode(&client))
	if res.StatusCode != http.StatusCreated || client.ClientID == "" {
		t.Fatalf("регистрация: %d %+v", res.StatusCode, client)
	}

	// PKCE
	verifier := "verifier-verifier-verifier-verifier-1234567890"
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	params := url.Values{
		"client_id": {client.ClientID}, "redirect_uri": {"http://localhost:9/cb"}, "response_type": {"code"},
		"code_challenge": {challenge}, "code_challenge_method": {"S256"}, "state": {"xyz"}, "resource": {oauthPublicURL},
	}

	// Форма входа показывает имя клиента
	res, err = http.Get(server.URL + "/oauth/authorize?" + params.Encode())
	f.ok(err)
	page, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK || !strings.Contains(string(page), "Claude Test") {
		t.Fatalf("форма: %d %.200s", res.StatusCode, page)
	}

	// Неверный пароль — форма снова, без редиректа
	form := url.Values{}
	for k, v := range params {
		form[k] = v
	}
	form.Set("email", email)
	form.Set("password", "wrong")
	res, err = noRedirect.PostForm(server.URL+"/oauth/authorize", form)
	f.ok(err)
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("неверный пароль: %d", res.StatusCode)
	}

	// Верный пароль — редирект с кодом и state
	form.Set("password", "secret-password")
	res, err = noRedirect.PostForm(server.URL+"/oauth/authorize", form)
	f.ok(err)
	location, err := url.Parse(res.Header.Get("Location"))
	f.ok(err)
	code := location.Query().Get("code")
	if res.StatusCode != http.StatusFound || !strings.HasPrefix(location.String(), "http://localhost:9/cb?") || code == "" || location.Query().Get("state") != "xyz" {
		t.Fatalf("редирект: %d %s", res.StatusCode, res.Header.Get("Location"))
	}

	// Неверный verifier — invalid_grant, код при этом сгорает, поэтому берем новый
	res, err = http.PostForm(server.URL+"/oauth/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {code}, "client_id": {client.ClientID},
		"redirect_uri": {"http://localhost:9/cb"}, "code_verifier": {"wrong"},
	})
	f.ok(err)
	var oauthErr struct {
		Error string `json:"error"`
	}
	f.ok(json.NewDecoder(res.Body).Decode(&oauthErr))
	if res.StatusCode != http.StatusBadRequest || oauthErr.Error != "invalid_grant" {
		t.Fatalf("неверный verifier: %d %+v", res.StatusCode, oauthErr)
	}
	res, err = noRedirect.PostForm(server.URL+"/oauth/authorize", form)
	f.ok(err)
	location, _ = url.Parse(res.Header.Get("Location"))
	code = location.Query().Get("code")

	// Обмен кода: JWT сервера с пользователем и stateless-устройством oauth-*
	res, err = http.PostForm(server.URL+"/oauth/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {code}, "client_id": {client.ClientID},
		"redirect_uri": {"http://localhost:9/cb"}, "code_verifier": {verifier}, "resource": {oauthPublicURL},
	})
	f.ok(err)
	var tokens struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	f.ok(json.NewDecoder(res.Body).Decode(&tokens))
	if res.StatusCode != http.StatusOK || tokens.TokenType != "Bearer" || tokens.ExpiresIn != 1800 || tokens.RefreshToken == "" {
		t.Fatalf("токены: %d %+v", res.StatusCode, tokens)
	}
	claims, err := jwtManager.ParseToken[auth.Claims](f.ctx, tokens.AccessToken, jwtManager.AccessToken)
	f.ok(err)
	if !strings.HasPrefix(claims.DeviceID, "oauth-") {
		t.Fatalf("устройство: %+v", claims)
	}
	if n := f.count(squirrelSelectDevices(claims.DeviceID)); n != 1 {
		t.Fatalf("устройство не создано stateless: %d", n)
	}

	// Повторный обмен того же кода — invalid_grant
	res, err = http.PostForm(server.URL+"/oauth/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {code}, "client_id": {client.ClientID},
		"redirect_uri": {"http://localhost:9/cb"}, "code_verifier": {verifier},
	})
	f.ok(err)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("повторный обмен прошел: %d", res.StatusCode)
	}

	// Обновление по refresh-токену
	res, err = http.PostForm(server.URL+"/oauth/token", url.Values{"grant_type": {"refresh_token"}, "refresh_token": {tokens.RefreshToken}})
	f.ok(err)
	var refreshed struct {
		AccessToken string `json:"access_token"`
	}
	f.ok(json.NewDecoder(res.Body).Decode(&refreshed))
	if res.StatusCode != http.StatusOK || refreshed.AccessToken == "" {
		t.Fatalf("refresh: %d", res.StatusCode)
	}
}

// squirrelSelectDevices — счетчик stateless-устройств с данным device_id
func squirrelSelectDevices(deviceID string) sq.SelectBuilder {
	return sq.Select("count(*)").From("tasks.devices").Where(sq.Eq{"device_id": deviceID, "is_stateless": true})
}
