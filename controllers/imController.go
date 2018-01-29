package controllers

import "github.com/astaxie/beego"

//ImController .
type ImController struct {
	beego.Controller
}

//RefreshToken .
// @Title 获取新token
// @Description 获取新token
// @Success 200 {object} utils.ResultDTO
// @router /refreshtoken [get]
func (c *ImController) RefreshToken() {

}
