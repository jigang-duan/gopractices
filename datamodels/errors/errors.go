package errors

import (
	"fmt"
	"github.com/kataras/iris"
)

const (
	CodeSuccess     = 0
	CodeCannotFound = 10000 + iota
	CodeDatabaseFound
	CodeCannotUpdate = 10010 + iota
	CodeVerifyForm
	CodePasswordGenerate
	CodeValidatePassword
	CodeDatabaseInsert
	CodeDatabaseUpdate
	CodeCannotDelete = 10020 + iota
)

var codeMessage = map[int]string{
	CodeSuccess:          "成功",
	CodeCannotFound:      "不能找到资源",
	CodeDatabaseFound:    "数据库找不到记录",
	CodeCannotUpdate:     "更新资源不成功",
	CodeVerifyForm:       "验证表单错误",
	CodePasswordGenerate: "密码生成出错",
	CodeValidatePassword: "密码验证错误",
	CodeDatabaseInsert:   "数据库插入记录失败",
	CodeDatabaseUpdate:   "数据库更新记录失败",
	CodeCannotDelete:     "不能删除",
}

var statusOfCode = map[int]int{
	CodeSuccess:          iris.StatusOK,
	CodeCannotFound:      iris.StatusFound,
	CodeDatabaseFound:    iris.StatusBadRequest,
	CodeCannotUpdate:     iris.StatusInternalServerError,
	CodeVerifyForm:       iris.StatusUnprocessableEntity,
	CodePasswordGenerate: iris.StatusInternalServerError,
	CodeValidatePassword: iris.StatusUnauthorized,
	CodeDatabaseInsert:   iris.StatusInternalServerError,
	CodeDatabaseUpdate:   iris.StatusInternalServerError,
	CodeCannotDelete:     iris.StatusInternalServerError,
}

var statusMessage = map[int]string{
	200: "服务器成功返回请求的数据",
	201: "新建或修改数据成功",
	202: "一个请求已经进入后台排队（异步任务）",
	204: "删除数据成功",
	400: "发出的请求有错误，服务器没有进行新建或修改数据的操作",
	401: "用户没有权限（令牌、用户名、密码错误）",
	403: "用户得到授权，但是访问是被禁止的",
	404: "发出的请求针对的是不存在的记录，服务器没有进行操作",
	406: "请求的格式不可得",
	410: "请求的资源被永久删除，且不会再得到的",
	422: "当创建一个对象时，发生一个验证错误",
	500: "服务器发生错误，请检查服务器",
	502: "网关错误",
	503: "服务不可用，服务器暂时过载或维护",
	504: "网关超时",
}

type HttpError struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HttpError) Error() string {
	if e != nil {
		return ""
	}
	return fmt.Sprintf("HTTP ERROR -> %s | <%d>: %s", statusMessage[e.Status], e.Code, e.Message)
}

func build(status, code int, message string) *HttpError {
	return &HttpError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func New(code int, message string) *HttpError {
	return build(statusOfCode[code], code, message)
}

func NewErr(code int, err error) *HttpError {
	if err == nil {
		return New(code, codeMessage[code])
	}
	return New(code, err.Error())
}

func NewSuccess() *HttpError {
	return New(CodeSuccess, "")
}

func (e *HttpError) Store(ctx iris.Context) {
	ctx.StatusCode(e.Status)
	ctx.Values().SetImmutable("error", *e)
	//
	//ctx.Values().Set("status", e.Status)
	//ctx.Values().Set("message", e.Message)
	//ctx.Values().Set("code", e.Code)
}

func (e *HttpError) StatusMessage() string {
	return statusMessage[e.Status]
}
