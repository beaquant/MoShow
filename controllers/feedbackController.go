package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"strconv"

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
	utils.JSONUnMarshal(c.GetString("img"), &r.Img)
	r.Content = c.GetString("content")
	f.LogFile = c.GetString("log")

	if err := f.AddSuggestion(r); err != nil {
		beego.Error(err)
		dto.Message = "添加反馈记录失败\t" + err.Error()
		return
	}

	utils.SendP2PSysMessage("已收到您的反馈，请耐心等待运营人员处理。", strconv.FormatUint(tk.ID, 10))
	dto.Message = "反馈成功!"
	dto.Sucess = true
}
