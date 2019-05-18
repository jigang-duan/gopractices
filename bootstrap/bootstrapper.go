package bootstrap

import (
	"github.com/gorilla/securecookie"
	"github.com/jigang-duan/gopractices/datamodels/errors"
	"github.com/jigang-duan/gopractices/web/middleware"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/sessions"
	"strings"
	"time"
)

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time

	ConfigFile string

	Sessions *sessions.Sessions
}

type Configurator func(*Bootstrapper)

func New(cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      "Awesome",
		AppOwner:     "Jigang Duan",
		ConfigFile:   "conf/config.yaml",
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

func (b *Bootstrapper) SetupViews(viewsDir string) {
	b.RegisterView(iris.HTML(viewsDir, ".html").Layout("shared/layout.html"))
}

func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + strings.ReplaceAll(b.AppName, " ", "_"),
		Expires:  expires,
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

func matchRest(ctx iris.Context) bool {
	match := false
	if route := ctx.GetCurrentRoute(); route != nil {
		if path := route.Path(); path[0:4] == "/api" {
			match = true
		}
	}
	if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
		match = true
	}
	return match
}

func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(middleware.Logger, func(ctx iris.Context) {
		if err, ok := ctx.Values().Get("error").(errors.HttpError); ok {
			ctx.StatusCode(err.Status)
			if matchRest(ctx) {
				ctx.JSON(err)
				return
			}
			ctx.ViewData("Err", err)
		}
		ctx.ViewData("Title", b.AppName)
		ctx.View("shared/error.html")
	})
}

const (
	// StaticAssets is the root directory for public assets like images, css, js.
	StaticAssets = "./web/public/"
	// Favicon is the relative 9to the "StaticAssets") favicon path for our app.
	Favicon = "favicon.ico"
)

func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews("./web/views")
	b.SetupSessions(24*time.Hour,
		[]byte("the-big-and-secret-fash-key-here"),
		[]byte("lot-secret-of-characters-big-too"),
	)
	b.SetupErrorHandlers()

	// static files
	b.Favicon(StaticAssets + Favicon)
	b.StaticWeb("/public", StaticAssets)

	b.Logger().SetLevel("debug")

	// middleware, after static files
	b.Use(recover.New())
	b.Use(middleware.Logger)
	b.Use(middleware.NewYaag())

	return b
}

func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	_ = b.Run(iris.Addr(addr), cfgs...)
}
