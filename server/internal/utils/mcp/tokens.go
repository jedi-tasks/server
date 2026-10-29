// Package mcp — транспорт MCP (Model Context Protocol) для агентов: HTTP-эндпоинт с Bearer-токенами,
// за которым стоят те же сервисы, что и за gRPC. Инструменты регистрируют модули в своих
// endpoint/mcp, здесь — только общая часть: авторизация, устройство агента, логирование вызовов.
package mcp

import (
	"crypto/subtle"
	"strings"

	"server/internal/utils/errors"
)

// Tokens — соответствие токенов доступа и почты их владельцев
type Tokens struct {
	byToken map[string]string // Токен → email
}

// ParseTokens разбирает MCP_TOKENS вида "<токен>:<email>,<токен>:<email>"
func ParseTokens(raw string) (Tokens, error) {
	tokens := Tokens{byToken: make(map[string]string)}

	// Разбираем пары через запятую, пробелы вокруг игнорируем
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		token, email, found := strings.Cut(pair, ":")
		if !found || token == "" || email == "" {
			return tokens, errors.InternalServer.New("MCP_TOKENS: ожидается <токен>:<email>")
		}
		if len(token) < 32 {
			return tokens, errors.InternalServer.New("MCP_TOKENS: токен короче 32 символов").WithParams("email", email)
		}
		tokens.byToken[token] = strings.TrimSpace(email)
	}

	return tokens, nil
}

// IsEmpty — токенов нет, MCP запускать незачем
func (t Tokens) IsEmpty() bool {
	return len(t.byToken) == 0
}

// Lookup возвращает владельца токена. Сравнение за постоянное время по всем токенам, чтобы по
// скорости ответа нельзя было подобрать токен посимвольно
func (t Tokens) Lookup(token string) (email string, ok bool) {
	for candidate, owner := range t.byToken {
		if subtle.ConstantTimeCompare([]byte(candidate), []byte(token)) == 1 {
			email, ok = owner, true
		}
	}
	return email, ok
}
