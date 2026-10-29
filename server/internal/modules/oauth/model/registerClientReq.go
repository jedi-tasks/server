package model

// RegisterClientReq — динамическая регистрация клиента (RFC 7591), поля из тела запроса
type RegisterClientReq struct {
	ClientName   string   `json:"client_name"`   // Название клиента
	RedirectURIs []string `json:"redirect_uris"` // Разрешенные redirect_uri
}
