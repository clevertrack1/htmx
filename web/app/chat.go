package app

import (
	"github.com/clevertrack1/mach"
	"html/template"
)

type ChatApp struct {
	Tmpl *template.Template
}

func (c *ChatApp) RegisterRoutes(app *mach.App) {
	app.GET("/chat", c.renderChat)
	app.POST("/chat", c.handleChat)
}

func (c *ChatApp) renderChat(ctx *mach.Context) {
	c.Tmpl.ExecuteTemplate(ctx.Response, "chat", nil)
}

type ResponseRender struct {
	UserMsg      string
	AssistantMsg string
}

func (c *ChatApp) handleChat(ctx *mach.Context) {
	userMsg := ctx.Request.FormValue("user-message")
	assistantMsg := `
		This is a placeholder response for an AI chat app.
		Connect your AI API to create your own app.
	`
	data := ResponseRender{
		UserMsg:      userMsg,
		AssistantMsg: assistantMsg,
	}
	c.Tmpl.ExecuteTemplate(ctx.Response, "chatresponse", data)
}
