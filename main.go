package main

import (
	"github.com/jigang-duan/gopractices/bootstrap"
	"github.com/jigang-duan/gopractices/conf"
	"github.com/jigang-duan/gopractices/datasource/entitys"
	"github.com/jigang-duan/gopractices/web"
	"github.com/jigang-duan/gopractices/web/middleware/identity"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New()
	app.Bootstrap()
	app.Configure(
		conf.Configure,
		identity.Configure,
		entitys.Configure,
		web.Routes,
	)
	return app
}

func main() {
	app := newApp()
	app.Listen(":8080")
}
