package controllers

import (
	"MoShow/models"
	"MoShow/utils"

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
	tk, dto := GetToken(c.Ctx), &utils.ResultDTO{}
	defer dto.JSONResult(&c.Controller)

	f := &models.FeedBack{UserID: tk.ID}
	r := &models.FeedBackSuggestion{}
	r.Img = c.GetString("img")
	r.Content = c.GetString("content")

	if err := f.AddSuggestion(r); err != nil {
		beego.Error(err)
		dto.Message = "添加反馈记录失败\t" + err.Error()
		return
	}

	dto.Message = "反馈成功!"
	dto.Sucess = true
}
