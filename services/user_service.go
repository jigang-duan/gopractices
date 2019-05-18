package services

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/jigang-duan/gopractices/conver"
	"github.com/jigang-duan/gopractices/datamodels"
	"github.com/jigang-duan/gopractices/datamodels/errors"
	"github.com/jigang-duan/gopractices/datasource/entitys"
)

type UserService interface {
	GetAll() []datamodels.User
	GetByID(id int64) (datamodels.User, *errors.HttpError)
	GetByUsername(username string) (datamodels.User, *errors.HttpError)
	GetByUsernameAndPassword(username, userPassword string) (datamodels.User, *errors.HttpError)
	DeleteByID(id int64) error

	Update(id int64, user datamodels.User) (datamodels.User, *errors.HttpError)
	UpdatePassword(id int64, newPassword string) (datamodels.User, *errors.HttpError)
	UpdateUsername(id int64, newUsername string) (datamodels.User, *errors.HttpError)

	Create(userPassword string, user datamodels.User) (datamodels.User, *errors.HttpError)
}

func NewUserService(eg *xorm.EngineGroup, converter conver.UserConverter) UserService {
	return &userService{
		eg:        eg,
		converter: converter,
	}
}

type userService struct {
	eg        *xorm.EngineGroup
	converter conver.UserConverter
}

func (s *userService) GetAll() []datamodels.User {
	var (
		results []datamodels.User
		users   []entitys.User
	)
	if err := s.eg.Asc("username").Find(&users); err == nil {
		for _, user := range users {
			results = append(results, s.converter.Reverse(user))
		}
	}
	return results
}

func generateDatabaseFoundError(has bool, err error) (datamodels.User, httpError) {
	if err != nil {
		herr := errors.NewErr(errors.CodeDatabaseFound, err)
		return datamodels.User{}, herr
	}
	herr := errors.New(errors.CodeDatabaseFound, "不存在的记录")
	return datamodels.User{}, herr
}

func generateFoundError(has bool, err error) (datamodels.User, error) {
	if err != nil {
		herr := errors.NewErr(errors.CodeDatabaseFound, err)
		return datamodels.User{}, herr
	}
	herr := errors.New(errors.CodeDatabaseFound, "不存在的记录")
	return datamodels.User{}, herr
}

func (s *userService) GetByID(id int64) (datamodels.User, httpError) {
	user := new(entitys.User)
	if has, err := s.eg.Id(id).Get(user); err != nil || !has {
		return generateDatabaseFoundError(has, err)
	}
	result := s.converter.Reverse(*user)
	return result, nil
}

func (s *userService) GetByUsername(username string) (datamodels.User, httpError) {
	user := &entitys.User{Username: username}
	if has, err := s.eg.Get(user); err != nil || !has {
		return generateDatabaseFoundError(has, err)
	}
	result := s.converter.Reverse(*user)
	return result, nil
}

func (s *userService) GetByUsernameAndPassword(username, userPassword string) (datamodels.User, httpError) {
	if username == "" || userPassword == "" {
		herr := errors.New(errors.CodeVerifyForm, "用户名或密码不能为空")
		return datamodels.User{}, herr
	}

	user := &entitys.User{Username: username}
	if has, err := s.eg.Get(user); err != nil || !has {
		return generateDatabaseFoundError(has, err)
	}

	if pass := user.ValidatePassword(userPassword, user.Password); pass {
		result := s.converter.Reverse(*user)
		return result, nil
	}
	herr := errors.New(errors.CodeValidatePassword, "用户名或密码不能为空")
	return datamodels.User{}, herr
}

func (s *userService) Update(id int64, user datamodels.User) (datamodels.User, httpError) {
	user.ID = id
	entity := s.converter.Convert(user)

	if affected, err := s.eg.Update(entity); affected == 0 {
		return user, errors.NewErr(errors.CodeDatabaseInsert, err)
	}

	return s.GetByID(id)
}

func (s *userService) UpdatePassword(id int64, newPassword string) (datamodels.User, httpError) {
	user := new(entitys.User)
	if has, err := s.eg.Id(id).Get(user); err != nil || !has {
		return generateDatabaseFoundError(has, err)
	}

	user.UpdatePassword(newPassword)
	if affected, err := s.eg.Update(user); affected == 0 {
		herr := errors.NewErr(errors.CodeDatabaseUpdate, err)
		return datamodels.User{}, herr
	}
	return s.GetByID(id)
}

func (s *userService) UpdateUsername(id int64, newUsername string) (datamodels.User, httpError) {
	user := new(entitys.User)
	if has, err := s.eg.Id(id).Get(user); err != nil || !has {
		return generateDatabaseFoundError(has, err)
	}
	user.Username = newUsername
	if affected, err := s.eg.Update(user); affected == 0 {
		herr := errors.NewErr(errors.CodeDatabaseUpdate, err)
		return datamodels.User{}, herr
	}
	return s.GetByID(id)
}

func (s *userService) Create(userPassword string, user datamodels.User) (datamodels.User, httpError) {
	if user.ID > 0 || userPassword == "" || user.Nickname == "" || user.Username == "" {
		herr := errors.New(errors.CodeVerifyForm, "无法创建此用户")
		return datamodels.User{}, herr
	}

	entity := s.converter.Convert(user)
	hashed := entity.GeneratePassword(userPassword)
	if !hashed {
		herr := errors.New(errors.CodePasswordGenerate, "密码没有生存成功")
		return datamodels.User{}, herr
	}
	if affected, err := s.eg.InsertOne(entity); affected == 0 {
		return user, errors.NewErr(errors.CodeDatabaseInsert, err)
	}
	return s.GetByUsername(user.Username)
}

func (s *userService) DeleteByID(id int64) error {
	user := new(entitys.User)
	if has, err := s.eg.ID(id).Exist(user); !has {
		return errors.NewErr(errors.CodeDatabaseFound, err)
	}

	affected, err := s.eg.ID(id).Delete(user)
	fmt.Print(affected)
	httpErr := errors.NewErr(errors.CodeCannotDelete, err)
	return httpErr
}
