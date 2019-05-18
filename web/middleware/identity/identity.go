package identity

import (
	"fmt"
	"github.com/jigang-duan/gopractices/bootstrap"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"time"
)

func New(b *bootstrap.Bootstrapper) iris.Handler {
	return func(ctx iris.Context) {
		ctx.Header("App-Name", b.AppName)
		ctx.Header("App-Owner", b.AppOwner)
		ctx.Header("App-Since", time.Since(b.AppSpawnDate).String())

		ctx.Header("Server", fmt.Sprintf("Web [%s]", viper.Get("app.sitedomain")))

		ctx.ViewData("AppName", b.AppName)
		ctx.ViewData("AppOwner", b.AppOwner)
		ctx.ViewData("Title", b.AppName)
		ctx.Next()
	}
}

func Configure(b *bootstrap.Bootstrapper) {
	h := New(b)
	b.UseGlobal(h)
}
