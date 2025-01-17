package main

import (
	"html/template"
	"net/http"
)

type MenuItem struct {
	Label  string
	Href   string
	Icon   string // This can be an SVG or font-awesome class
	Active bool
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Define the sidebar menu items
		menuItems := []MenuItem{
			{Label: "Dashboard", Href: "/dashboard", Icon: `<svg width="20" height="20" fill="white" xmlns="http://www.w3.org/2000/svg"><circle cx="10" cy="10" r="9" stroke="white"/></svg>`, Active: r.URL.Path == "/dashboard"},
			{Label: "TX Monitor", Href: "/tx-monitor", Icon: `<svg width="20" height="20" fill="white" xmlns="http://www.w3.org/2000/svg"><rect width="20" height="20" fill="white"/></svg>`, Active: r.URL.Path == "/tx-monitor"},
			{Label: "RX Monitor", Href: "/rx-monitor", Icon: `<svg width="20" height="20" fill="white" xmlns="http://www.w3.org/2000/svg"><path d="M10 2L2 18h16L10 2z" fill="white"/></svg>`, Active: r.URL.Path == "/rx-monitor"},
		}

		t, err := template.ParseFiles("templates/base.html", "templates/sidebar.html")
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
			return
		}

		// Render the template with the menu items
		t.Execute(w, map[string]interface{}{
			"MenuItems": menuItems,
		})
	})

	// Static file server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Dashboard</h1>")) // Replace this with a template if needed
	})

	http.HandleFunc("/tx-monitor", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>TX Monitor</h1>"))
	})

	// Run the server
	http.ListenAndServe(":8080", nil)
}
