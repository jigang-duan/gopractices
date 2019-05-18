package conver

import (
	"github.com/jigang-duan/gopractices/datamodels"
	"github.com/jigang-duan/gopractices/datasource/entitys"
)

type UserConverter struct {
}

func (c UserConverter) Reverse(user entitys.User) datamodels.User {
	return datamodels.User{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

func (c UserConverter) Convert(user datamodels.User) entitys.User {
	entity := entitys.User{
		Username: user.Username,
		Nickname: user.Nickname,
	}
	return entity
}

func (c UserConverter) ConvertWithPassword(user datamodels.User, userPassword string) entitys.User {
	entity := c.Convert(user)
	entity.GeneratePassword(userPassword)
	return entity
}
