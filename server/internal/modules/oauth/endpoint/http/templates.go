package http

import (
	"html/template"
	"net/http"
	"net/url"
)

// formData — данные формы входа: имя клиента, параметры авторизации для скрытых полей, ошибка
type formData struct {
	Client string     // Название клиента, запрашивающего доступ
	Params url.Values // Параметры запроса авторизации, повторяются скрытыми полями
	Error  string     // Текст ошибки входа, пусто — нет
}

// hiddenParams — параметры авторизации, которые форма возвращает обратно
var hiddenParams = []string{"client_id", "redirect_uri", "response_type", "code_challenge", "code_challenge_method", "resource", "state", "scope"}

var formTemplate = template.Must(template.New("form").Parse(`<!doctype html>
<html lang="ru"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Tasks — вход</title>
<style>
body{font:16px -apple-system,system-ui,sans-serif;background:#f2f2f7;margin:0;display:flex;min-height:100vh;align-items:center;justify-content:center}
form{background:#fff;padding:28px;border-radius:14px;width:320px;box-shadow:0 2px 12px rgba(0,0,0,.08)}
h1{font-size:20px;margin:0 0 6px}p{margin:0 0 18px;color:#555}
input{width:100%;box-sizing:border-box;padding:10px;margin:0 0 12px;border:1px solid #ccc;border-radius:8px;font-size:16px}
button{width:100%;padding:11px;border:0;border-radius:8px;background:#3e9bf0;color:#fff;font-size:16px}
.err{color:#c62828;margin:0 0 12px}
</style></head><body>
<form method="post" action="/oauth/authorize">
<h1>Tasks</h1>
<p><b>{{.Client}}</b> запрашивает доступ к вашим задачам.</p>
{{if .Error}}<p class="err">{{.Error}}</p>{{end}}
<input type="email" name="email" placeholder="Почта" required autofocus autocomplete="username">
<input type="password" name="password" placeholder="Пароль" required autocomplete="current-password">
{{range $name, $value := .Hidden}}<input type="hidden" name="{{$name}}" value="{{$value}}">{{end}}
<button type="submit">Разрешить</button>
</form></body></html>`))

var errorTemplate = template.Must(template.New("error").Parse(`<!doctype html>
<html lang="ru"><head><meta charset="utf-8"><title>Tasks — ошибка</title>
<style>body{font:16px -apple-system,system-ui,sans-serif;padding:40px;color:#333}</style></head>
<body><h1>Не удалось авторизовать клиента</h1><p>{{.}}</p></body></html>`))

// renderForm показывает форму входа
func renderForm(w http.ResponseWriter, status int, data formData) {
	hidden := make(map[string]string, len(hiddenParams))
	for _, name := range hiddenParams {
		if value := data.Params.Get(name); value != "" {
			hidden[name] = value
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = formTemplate.Execute(w, struct {
		Client, Error string
		Hidden        map[string]string
	}{data.Client, data.Error, hidden})
}

// renderError показывает ошибку, которую нельзя отдать клиенту редиректом
func renderError(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = errorTemplate.Execute(w, text)
}
