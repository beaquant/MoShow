package utils

import (
	"github.com/astaxie/beego"
)

const (
	//DtoStatusNormal 正常状态
	DtoStatusNormal int = iota
	//DtoStatusAuthError 登录验证失败
	DtoStatusAuthError
	//DtoStatusFrequencyError 访问频率错误
	DtoStatusFrequencyError
	//DtoStatusDatabaseError 数据库操作错误
	DtoStatusDatabaseError
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
	if err := recover(); err != nil {
		dto.Sucess = false

		if e, ok := err.(error); ok {
			dto.Message = e.Error()
		}
		beego.Error(err)
	}

	c.Data["json"] = dto
	c.ServeJSON()
}
