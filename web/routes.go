package web

import (
	"github.com/jigang-duan/gopractices/bootstrap"
	"github.com/jigang-duan/gopractices/conver"
	"github.com/jigang-duan/gopractices/datasource/entitys"
	"github.com/jigang-duan/gopractices/services"
	"github.com/jigang-duan/gopractices/web/controllers"
	"github.com/jigang-duan/gopractices/web/middleware"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/mvc"
)

// Routes Configure
func Routes(b *bootstrap.Bootstrapper) {
	userService := services.NewUserService(entitys.ORM, conver.UserConverter{})

	users := mvc.New(b.Party("/api/users"))
	users.Router.Use(middleware.BasicAuth)
	users.Register(userService)
	users.Handle(new(controllers.UsersController))

	user := mvc.New(b.Party("/user"))
	user.Register(userService, b.Sessions.Start)
	user.Handle(new(controllers.UserController))

	b.Any("/debug/pprof/{action:path}", pprof.New())
}
