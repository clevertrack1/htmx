package main

import (
	"bitbucket.org/paulcrfi/htmx/utils"
	"bitbucket.org/paulcrfi/htmx/web/app"
	"flag"
	"fmt"
	"github.com/clevertrack1/mach"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func main() {
	address := flag.String("address", "localhost:8085", "Port to run server on")
	flag.Parse()

	// Parse all templates once
	funcMap := template.FuncMap{"dict": utils.Dict}
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("web/templates/*/*")
	if err != nil {
		panic(err)
	}

	appInstance := mach.Default()

	appInstance.Static("/static/", "./web/static")

	appInstance.GET("/", func(c *mach.Context) {
		// Extract the filename name from the URL
		templateName := strings.TrimPrefix(c.Path(), "/")
		if templateName == "" {
			templateName = "index"
		}

		templatePath := fmt.Sprintf("./web/templates/views/%s.tmpl", templateName)

		t, err := template.ParseFiles(templatePath)
		if err != nil {
			c.NoContent(http.StatusNotFound)
			return
		}

		t.Execute(c.Response, nil)
	})

	// Register App routes
	emailApp := app.EmailApp{Tmpl: tmpl}
	emailApp.RegisterRoutes(appInstance)

	dashApp := app.DashboardApp{Tmpl: tmpl}
	dashApp.RegisterRoutes(appInstance)

	chatApp := app.ChatApp{Tmpl: tmpl}
	chatApp.RegisterRoutes(appInstance)

	// Search route
	search := utils.InitSearchIndex("web/templates", tmpl)
	appInstance.POST("/search", func(c *mach.Context) {
		searchInput := c.Request.FormValue("searchInput")
		results := search.QueryIndex(searchInput)
		c.Response.Header().Set("Content-Type", "text/html")
		c.Response.Write([]byte(results))
	})

	// Run the server with graceful shutdown
	if err := appInstance.Run(*address, mach.WithGracefulShutdown(5*time.Second)); err != nil {
		panic(fmt.Sprintf("Server failed: %s", err))
	}
}
