package app

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
)

type DashboardApp struct {
	Tmpl *template.Template
}

func (d *DashboardApp) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /dashboard", d.renderDashboard)
	mux.HandleFunc("GET /prices", d.renderPrices)
}

func (d *DashboardApp) renderDashboard(w http.ResponseWriter, r *http.Request) {
	d.Tmpl.ExecuteTemplate(w, "dashboard", nil)
}

type PricesRender struct {
	BTC string
	ETH string
	ZEC string
}

func (d *DashboardApp) renderPrices(w http.ResponseWriter, r *http.Request) {
	data := PricesRender{
		BTC: getBTC(),
		ETH: getETH(),
		ZEC: getZEC(),
	}
	d.Tmpl.ExecuteTemplate(w, "prices", data)
}

func getBTC() string {
	return randPrice(99950, 100050)
}
func getETH() string {
	return randPrice(3750, 3850)
}
func getZEC() string {
	return randPrice(57, 59)
}

func randPrice(min, max float64) string {
	return fmt.Sprintf("$%f", min+(rand.Float64()*(max-min)))
}
