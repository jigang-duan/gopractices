package controllers

import (
	"github.com/jigang-duan/gopractices/datamodels"
	"github.com/jigang-duan/gopractices/datamodels/errors"
	"github.com/jigang-duan/gopractices/services"
	"github.com/kataras/iris"
)

type UsersController struct {
	Ctx iris.Context
	Service services.UserService
}

func (c *UsersController) Get() (results []datamodels.User) {
	return c.Service.GetAll()
}

func (c *UsersController) GetBy(id int64) (datamodels.User, error) {
	u, err := c.Service.GetByID(id)
	if err == nil {
		return u, nil
	}
	err.Store(c.Ctx)
	return u , err
}

func (c *UsersController) PutBy(id int64) (datamodels.User, error) {
	u := datamodels.User{}
	if err := c.Ctx.ReadJSON(&u); err != nil {
		return u, err
	}

	user, err := c.Service.Update(id, u)
	if err == nil {
		return user, nil
	}
	err.Store(c.Ctx)
	return user , err
}

func (c *UsersController) DeleteBy(id int64) interface{} {
	err := c.Service.DeleteByID(id)
	if err == nil {
		return map[string]interface{}{"deleted": id}
	}
	if httpError, match := err.(*errors.HttpError); match {
		httpError.Store(c.Ctx)
		return httpError.Status
	}
	return iris.StatusNotExtended
}