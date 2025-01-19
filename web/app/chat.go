package app

import (
	"html/template"
	"net/http"
)

type ChatApp struct {
	Tmpl *template.Template
}

func (c *ChatApp) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /chat", c.renderChat)
	mux.HandleFunc("POST /chat", c.handleChat)
}

func (c *ChatApp) renderChat(w http.ResponseWriter, r *http.Request) {
	c.Tmpl.ExecuteTemplate(w, "chat", nil)
}

type ResponseRender struct {
	UserMsg      string
	AssistantMsg string
}

func (c *ChatApp) handleChat(w http.ResponseWriter, r *http.Request) {
	userMsg := r.FormValue("user-message")
	assistantMsg := `
		This is a placeholder response for an AI chat app.
		Connect your AI API to create your own app.
	`
	data := ResponseRender{
		UserMsg:      userMsg,
		AssistantMsg: assistantMsg,
	}
	c.Tmpl.ExecuteTemplate(w, "chatresponse", data)
}
