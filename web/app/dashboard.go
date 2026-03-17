package app

import (
	"fmt"
	"github.com/clevertrack1/mach"
	"html/template"
	"math/rand"
)

type DashboardApp struct {
	Tmpl *template.Template
}

func (d *DashboardApp) RegisterRoutes(app *mach.App) {
	app.GET("/dashboard", d.renderDashboard)
	app.GET("/prices", d.renderPrices)
}

func (d *DashboardApp) renderDashboard(c *mach.Context) {
	d.Tmpl.ExecuteTemplate(c.Response, "dashboard", nil)
}

type PricesRender struct {
	BTC string
	ETH string
	ZEC string
}

func (d *DashboardApp) renderPrices(c *mach.Context) {
	data := PricesRender{
		BTC: getBTC(),
		ETH: getETH(),
		ZEC: getZEC(),
	}
	d.Tmpl.ExecuteTemplate(c.Response, "prices", data)
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
