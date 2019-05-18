package controllers

import (
	"github.com/jigang-duan/gopractices/datamodels"
	"github.com/jigang-duan/gopractices/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

type UserController struct {
	Ctx iris.Context

	Service services.UserService

	Session *sessions.Session
}

const userIDKey = "UserID"

func (c *UserController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey, 0)
	return userID
}

func (c *UserController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *UserController) logout()  {
	c.Session.Destroy()
}

var registerStaticView = mvc.View{
	Name: "user/register.html",
	Data: iris.Map{"Title": "用户注册"},
}

func (c *UserController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}

	return registerStaticView
}

func (c *UserController) PostRegister() mvc.Result {
	var (
		nickname = c.Ctx.FormValue("nickname")
		username = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
	)

	u, err := c.Service.Create(password, datamodels.User{
		Username: username,
		Nickname: nickname,
	})

	if err != nil {
		c.Session.Set(username, u.ID)
		return mvc.Response{ Err: err }
	}

	return mvc.Response{
		Path: "/user/me",
	}
}

var loginStaticView = mvc.View{
	Name: "user/login.html",
	Data: iris.Map{"Title": "User Login"},
}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserController) GetLogin() mvc.Result {
	if c.isLoggedIn() {
		// if it's already logged in then destroy the previous session.
		c.logout()
	}

	return loginStaticView
}

// PostLogin handles POST: http://localhost:8080/user/lotin.
func (c *UserController) PostLogin() mvc.Result {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)

	u, err := c.Service.GetByUsernameAndPassword(username, password)
	if err != nil {
		return mvc.Response{
			Path: "/user/register",
		}
	}

	c.Session.Set(userIDKey, u.ID)

	return mvc.Response{
		Path: "/user/me",
	}
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserController) GetMe() mvc.Result {
	if !c.isLoggedIn() {
		// if it's not logged in then redirect user to the login page.
		return mvc.Response{Path: "/user/login"}
	}

	u, err := c.Service.GetByID(c.getCurrentUserID())
	if err != nil {
		// if the  session exists but for some reason the user doesn't exist in the "database"
		// then logout and re-execute the function, it will redirect the client to the
		// /user/login page.
		c.logout()
		return c.GetMe()
	}

	return mvc.View{
		Name: "user/me.html",
		Data: iris.Map{
			"Title": "Profile of " + u.Username,
			"User":  u,
		},
	}
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *UserController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}

	c.Ctx.Redirect("/user/login")
}
