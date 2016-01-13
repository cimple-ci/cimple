package frontend

import "github.com/lukesmith/cimple/server/web_application"

func RegisterAssets(app *web_application.Application) {
	app.Static("/assets/js/application.js")
}
