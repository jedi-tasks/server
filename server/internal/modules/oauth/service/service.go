// Package service — Authorization Server OAuth 2.1 для MCP-клиентов: регистрация клиентов,
// выдача кодов после входа пользователя, обмен кодов и refresh-токенов на JWT сервера
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	authModel "server/internal/modules/auth/model"
	authService "server/internal/modules/auth/service"
	"server/internal/modules/oauth/model"
	oauthRepository "server/internal/modules/oauth/repository"
)

var tracer = otel.Tracer("/server/internal/modules/oauth/service")

var _ OAuthRepository = new(oauthRepository.OAuthRepository)

// OAuthRepository — хранение клиентов и кодов
type OAuthRepository interface {
	CreateClient(context.Context, model.CreateClientReq) error
	GetClient(ctx context.Context, id uuid.UUID) (model.Client, bool, error)
	TouchClient(ctx context.Context, id uuid.UUID, now time.Time) error
	CreateCode(context.Context, model.CreateCodeReq) error
	ConsumeCode(ctx context.Context, code string, now time.Time) (model.Code, bool, error)
}

var _ AuthService = new(authService.AuthService)

// AuthService — проверка пароля и выдача токенов сервера
type AuthService interface {
	VerifyPassword(ctx context.Context, email, password string) (uuid.UUID, error)
	IssueTokens(context.Context, authModel.IssueTokensReq) (authModel.Tokens, error)
	RefreshTokens(context.Context, authModel.RefreshTokensReq) (authModel.RefreshTokensRes, error)
}

// Config — параметры Authorization Server
type Config struct {
	PublicURL      string        // Канонический URL MCP, он же issuer и resource
	AccessTokenTTL time.Duration // Срок жизни access-токена, для expires_in
	Version        string        // Версия сервера, для информации об устройстве
}

// OAuthService — Authorization Server
type OAuthService struct {
	oauthRepository OAuthRepository
	authService     AuthService
	config          Config
}

// NewOAuthService создает Authorization Server
func NewOAuthService(oauthRepository OAuthRepository, authService AuthService, config Config) *OAuthService {
	return &OAuthService{oauthRepository: oauthRepository, authService: authService, config: config}
}

// PublicURL возвращает канонический URL MCP
func (s *OAuthService) PublicURL() string {
	return s.config.PublicURL
}
