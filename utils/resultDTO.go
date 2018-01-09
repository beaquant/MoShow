package utils

const (
	//StatusNormal 正常状态
	StatusNormal int = iota
	//StatusAuthError 登录验证失败
	StatusAuthError
)

//ResultDTO .
type ResultDTO struct {
	Sucess  bool        `description:"API调用成功与否"`
	Data    interface{} `description:"API返回的结果数据"`
	Message string      `description:"API返回的消息"`
	Code    int         `description:"API调用状态码"`
}
