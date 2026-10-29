package integrationTest

import (
	"testing"
	"time"

	"github.com/google/uuid"

	oauthModel "server/internal/modules/oauth/model"
	oauthRepository "server/internal/modules/oauth/repository"
)

func TestOAuth_CodeIsSingleUse(t *testing.T) {
	f := newFixture(t)
	repo := oauthRepository.NewOAuthRepository(database)
	now := time.Now()

	// Клиент регистрируется, redirect_uri хранятся массивом
	clientID := uuid.New()
	f.ok(repo.CreateClient(f.ctx, oauthModel.CreateClientReq{
		ID: clientID, Name: "Claude", RedirectURIs: []string{"http://localhost:1/cb", "https://claude.ai/cb"}, CreatedAt: now,
	}))
	client, found, err := repo.GetClient(f.ctx, clientID)
	f.ok(err)
	if !found || len(client.RedirectURIs) != 2 || client.Name != "Claude" {
		t.Fatalf("клиент не прочитан: %+v %v", client, found)
	}
	if _, found, err := repo.GetClient(f.ctx, uuid.New()); err != nil || found {
		t.Fatalf("неизвестный клиент найден: %v %v", found, err)
	}

	// Код обменивается ровно один раз
	f.ok(repo.CreateCode(f.ctx, oauthModel.CreateCodeReq{
		Code: "code-1", ClientID: clientID, UserID: f.necessary.UserID, RedirectURI: "http://localhost:1/cb",
		CodeChallenge: "challenge", Resource: "https://mcp.test/", ExpiresAt: now.Add(10 * time.Minute),
	}))
	code, found, err := repo.ConsumeCode(f.ctx, "code-1", now)
	f.ok(err)
	if !found || code.UserID != f.necessary.UserID || code.CodeChallenge != "challenge" {
		t.Fatalf("код не обменян: %+v %v", code, found)
	}
	if _, found, err := repo.ConsumeCode(f.ctx, "code-1", now); err != nil || found {
		t.Fatalf("код обменян повторно: %v %v", found, err)
	}

	// Просроченный код не обменивается
	f.ok(repo.CreateCode(f.ctx, oauthModel.CreateCodeReq{
		Code: "code-2", ClientID: clientID, UserID: f.necessary.UserID, RedirectURI: "http://localhost:1/cb",
		CodeChallenge: "challenge", Resource: "https://mcp.test/", ExpiresAt: now.Add(-time.Minute),
	}))
	if _, found, err := repo.ConsumeCode(f.ctx, "code-2", now); err != nil || found {
		t.Fatalf("просроченный код обменян: %v %v", found, err)
	}
}
