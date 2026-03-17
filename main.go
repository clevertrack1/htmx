package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/clevertrack1/htmx/utils"
	"github.com/clevertrack1/htmx/web/app"
	"github.com/clevertrack1/mach"
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
		// Extract the template name from the URL
		templateName := strings.TrimPrefix(c.Path(), "/")
		if templateName == "" {
			templateName = "index.tmpl"
		} else {
			// Append .tmpl if not present
			if !strings.HasSuffix(templateName, ".tmpl") && !strings.HasSuffix(templateName, ".html") {
				templateName += ".tmpl"
			}
		}

		err := tmpl.ExecuteTemplate(c.Response, templateName, nil)
		if err != nil {
			c.NoContent(http.StatusNotFound)
			return
		}
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
