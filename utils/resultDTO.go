package utils

import (
	"github.com/astaxie/beego"
)

const (
	//DtoStatusNormal 正常状态
	DtoStatusNormal int = iota
	//DtoStatusUnkownError 未知错误
	DtoStatusUnkownError
	//DtoStatusAuthError 登录验证失败
	DtoStatusAuthError
	//DtoStatusFrequencyError 访问频率错误
	DtoStatusFrequencyError
	//DtoStatusDatabaseError 数据库操作错误
	DtoStatusDatabaseError
	//DtoStatusParamError 请求参数错误
	DtoStatusParamError
)

//ResultDTO .
type ResultDTO struct {
	Sucess  bool        `description:"API调用成功与否"`
	Data    interface{} `description:"API返回的结果数据"`
	Message string      `description:"API返回的消息"`
	Code    int         `description:"API调用状态码"`
}

//JSONResult .
func (dto *ResultDTO) JSONResult(c *beego.Controller) {
	// if err := recover(); err != nil {
	// 	dto.Sucess = false

	// 	if e, ok := err.(error); ok {
	// 		dto.Message = e.Error()
	// 	}
	// 	beego.Error(err)
	// }

	if !dto.Sucess && len(dto.Message) == 0 && dto.Code == 0 {
		dto.Message = "未知错误"
		dto.Code = DtoStatusUnkownError
	}

	c.Data["json"] = dto
	c.ServeJSON()
}
