package main

import (
	"bitbucket.org/paulcrfi/htmx/app"
	"bitbucket.org/paulcrfi/htmx/utils"
	"context"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MenuItem struct {
	Label  string
	Href   string
	Icon   string // This can be an SVG or font-awesome class
	Active bool
}

func main() {
	// Create a channel to catch OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	address := flag.String("address", "localhost:8085", "Port to run server on")
	flag.Parse()

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	// Define the sidebar menu items
	//	menuItems := []MenuItem{
	//		{Label: "Dashboard", Href: "/dashboard", Icon: `<svg width="20" height="20" fill="white" xmlns="http://www.w3.org/2000/svg"><circle cx="10" cy="10" r="9" stroke="white"/></svg>`, Active: r.URL.Path == "/dashboard"},
	//		{Label: "TX Monitor", Href: "/tx-monitor", Icon: `<svg width="20" height="20" fill="white" xmlns="http://www.w3.org/2000/svg"><rect width="20" height="20" fill="white"/></svg>`, Active: r.URL.Path == "/tx-monitor"},
	//		{Label: "RX Monitor", Href: "/rx-monitor", Icon: `<svg width="20" height="20" fill="white" xmlns="http://www.w3.org/2000/svg"><path d="M10 2L2 18h16L10 2z" fill="white"/></svg>`, Active: r.URL.Path == "/rx-monitor"},
	//	}
	//
	//	t, err := template.ParseFiles("templates/base.html", "templates/sidebar.html")
	//	if err != nil {
	//		http.Error(w, "Could not load template", http.StatusInternalServerError)
	//		return
	//	}
	//
	//	// Render the template with the menu items
	//	t.Execute(w, map[string]interface{}{
	//		"MenuItems": menuItems,
	//	})
	//})
	//
	//// Static file server
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//
	//http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("<h1>Dashboard</h1>")) // Replace this with a template if needed
	//})
	//
	//http.HandleFunc("/tx-monitor", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("<h1>TX Monitor</h1>"))
	//})

	var err error
	funcMap := template.FuncMap{"dict": utils.Dict}
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("web/templates/*/*")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", utils.RenderTemplate(tmpl, "index.html", nil))
	mux.HandleFunc("GET /cheatsheet", utils.RenderTemplate(tmpl, "cheatsheet.html", nil))
	mux.HandleFunc("GET /test", utils.RenderTemplate(tmpl, "test.html", nil))
	mux.HandleFunc("GET /docs/{doc}", utils.RenderDoc(tmpl, nil))
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("web/public"))))

	email := app.EmailApp{Tmpl: tmpl}
	email.RegisterRoutes(mux)

	dash := app.DashboardApp{Tmpl: tmpl}
	dash.RegisterRoutes(mux)

	chat := app.ChatApp{Tmpl: tmpl}
	chat.RegisterRoutes(mux)

	search := utils.InitSearchIndex("web/templates", tmpl)
	mux.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		searchInput := r.FormValue("searchInput")
		results := search.QueryIndex(searchInput)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(results))
	})

	// Create an HTTP server instance
	server := &http.Server{
		Addr:    *address,
		Handler: mux,
	}

	// Run the server
	// Start the server in a goroutine
	go func() {
		fmt.Printf("Listening at %s\n", *address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Server failed: %s", err))
		}
	}()
	// Block until we receive a signal to shut down
	<-signalChan
	fmt.Println("Shutting down server...")

	// Create a context with timeout to allow graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %s\n", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}
