package controllers

import (
	"github.com/astaxie/beego"
)

//FeedbackController 反馈节点
type FeedbackController struct {
	beego.Controller
}

//Suggestion .
// @Title 反馈建议
// @Description 反馈建议
// @Param   content     	formData    string  	true       "反馈内容"
// @Param   img		    	formData    string  	false       "图片"
// @Success 200 {object} utils.ResultDTO
// @router /suggestion [post]
func (c *FeedbackController) Suggestion() {

}
