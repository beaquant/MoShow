package controllers

import (
	"MoShow/models"
	"MoShow/utils"
	"boys/common"
	"errors"
	"runtime/debug"
	"strconv"
	"time"

	netease "github.com/MrSong0607/netease-im"
	"github.com/MrSong0607/wechat/oauth"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/garyburd/redigo/redis"
)

var (
	timeFormat = "2006-01-02T15:04:05.000Z"
	adminPhone = beego.AppConfig.String("adminPhoneNum")
	adminCode  = beego.AppConfig.String("adminCode")
	faker      = &fakerUserInfo{}
	actives    = &CacheActiveInfo{}
)

//AuthController 短信登陆，微信登陆，发送验证码，退出登陆等
type AuthController struct {
	beego.Controller
}

type codeInfo struct {
	Code string
	Time int64
}

type fakerUserInfo struct {
	Users map[string]models.User
	Time  time.Time
}

//ActiveInfo .
type ActiveInfo struct {
	models.Active
	Detail models.ActiveDetail
}

//CacheActiveInfo .
type CacheActiveInfo struct {
	Actives []ActiveInfo
	Time    time.Time
}

//SendCode .
// @Title 发送验证码
// @Description 发送验证码
// @Param   phone     path    string  true        "接收验证码的手机号"
// @Success 200 {object} utils.ResultDTO
// @router /:phone/sendcode [post]
func (c *AuthController) SendCode() {
	dto := &utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)
	con := utils.RedisPool.Get()
	defer con.Close()

	num := c.Ctx.Input.Param(":phone")
	codeEx, err := redis.String(con.Do("HGET", SmsCodeRedisKey, num))
	if err != nil && err != redis.ErrNil {
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	ci := &codeInfo{}
	utils.JSONUnMarshal(codeEx, ci)

	if ci != nil && ci.Time > time.Now().Add(time.Minute*27).Unix() {
		dto.Message = "验证码获取太频繁"
		return
	}

	code := strconv.Itoa(utils.RandNumber(1000, 9999))

	if res, err := utils.SendMsgByAPIKey(num, code); err != nil {
		beego.Error("发送验证码失败:", num, res, err)
		dto.Message = err.Error()
	} else {
		cs, _ := utils.JSONMarshalToString(&codeInfo{Code: code, Time: time.Now().Add(time.Minute * 30).Unix()})

		con.Do("HSET", SmsCodeRedisKey, num, cs)
		dto.Sucess = true
		dto.Message = "验证码发送成功"
	}
}

//Login .
// @Title 登陆或者注册
// @Description 登陆或者注册
// @Param   phone     path    	  string  true	 "手机号"
// @Param   code	  formData     string  true	 "验证码"
// @Success 200 {object} utils.ResultDTO
// @router /:phone/login [post]
func (c *AuthController) Login() {
	dto := utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)
	con := utils.RedisPool.Get()
	defer con.Close()

	phoneNum := c.Ctx.Input.Param(":phone")
	code := c.GetString("code")

	if phoneNum != adminPhone && !isInternalAcct(phoneNum) {
		codeEx, err := redis.String(con.Do("HGET", SmsCodeRedisKey, phoneNum))
		if err != nil {
			if err == redis.ErrNil {
				dto.Message = "请先获取验证码"
				return
			}

			beego.Error("读取验证码失败", err)
			dto.Message = "读取验证码失败" + err.Error()
			return
		}

		ci := &codeInfo{}
		utils.JSONUnMarshal(codeEx, ci)

		if ci.Time < time.Now().Unix() {
			dto.Message = "验证码已过期,请重新获取"
			return
		}

		if ci.Code != code {
			dto.Message = "验证码错误"
			return
		}
	} else if code != adminCode {
		dto.Message = "验证码错误"
		return
	}

	u := &models.User{PhoneNumber: phoneNum}
	if err := u.ReadFromPhoneNumber(); err != nil {
		beego.Error(err)
		dto.Message = err.Error()
		return
	}

	tk := &Token{}
	if u.ID == 0 { //该手机号未注册，执行注册逻辑
		up, err := c.InitUser(u, models.AcctTypeTelephone)
		if err != nil {
			beego.Error("注册用户失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "注册用户失败\t" + err.Error()
			return
		}

		tk.ID = up.ID
		tk.AcctStatus = u.AcctStatus
		tk.UserType = up.UserType
		if err := SetToken(c.Ctx, tk); err != nil {
			beego.Error("设置token失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "设置token失败\t" + err.Error()
			return
		}

		dto.Message = "注册成功"
		dto.Data = &UserProfileInfo{UserProfile: *up, ImTk: up.ImToken}
		dto.Sucess = true
	} else {
		if u.AcctStatus != models.AcctStatusShield {
			if IsCheckMode4Context(c.Ctx) {
				(&models.UserProfile{ID: u.ID}).SetFaker()
			}

			if u.AcctType == models.AcctTypeInternal {
				go addAcctLoginInfo(u.ID, c.Ctx) //添加登陆信息
			}

			if err := (&models.UserProfile{ID: u.ID}).UpdateOnlineStatus(models.OnlineStatusOnline); err != nil {
				beego.Error("更新在线状态失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "更新在线状态失败\t" + err.Error()
				return
			}

			up := &models.UserProfile{ID: u.ID}
			if err := up.Read(); err != nil {
				beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "获取用户信息失败\t" + err.Error()
				return
			}

			var err error
			if dto.Data, err = genSelfUserPorfileInfo(up, nil); err != nil {
				beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "获取用户信息失败\t" + err.Error()
				return
			}

			tk.ID = u.ID
			tk.AcctStatus = u.AcctStatus
			tk.UserType = up.UserType
			if err := SetToken(c.Ctx, tk); err != nil {
				beego.Error("设置token失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "设置token失败\t" + err.Error()
				return
			}

			dto.Sucess = true
			dto.Message = "登陆成功"
		} else {
			dto.Message = "您因涉及违规，已被封号。可联系平台运营"
			dto.Code = utils.DtoStatusAccountDisable
		}
	}
}

//WechatLogin .
// @Title 登陆或者注册
// @Description 登陆或者注册
// @Param   AccessToken     formData     string  true        "The email for login"
// @Param   OpenID          formData     string  true        "The email for login"
// @Success 200 {object} utils.ResultDTO
// @router /wechatlogin [post]
func (c *AuthController) WechatLogin() {
	dto := utils.ResultDTO{Sucess: false}
	defer dto.JSONResult(&c.Controller)

	AccessToken := c.GetString("AccessToken")
	OpenID := c.GetString("OpenID")
	Ivt, _ := strconv.ParseUint(c.GetString("Ivt"), 10, 64) //邀请人信息

	o := oauth.NewOauth(nil)
	info, err := o.GetUserInfo(AccessToken, OpenID)
	if err != nil {
		dto.Message = err.Error()
		beego.Error(err)
		return
	}

	u := &models.User{WeChatID: info.Unionid}
	err = u.ReadFromWechatID()
	if err != nil {
		dto.Message = err.Error()
		beego.Error(err)
		return
	}

	tk := &Token{}
	if u.ID == 0 { //执行微信注册
		u.InvitedBy = Ivt
		up, err := c.InitUser(u, models.AcctTypeWechat)
		if err != nil {
			beego.Error("注册用户失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "注册用户失败\t" + err.Error()
			return
		}

		tk.ID = up.ID
		tk.AcctStatus = u.AcctStatus
		tk.UserType = up.UserType
		if err := SetToken(c.Ctx, tk); err != nil {
			beego.Error("设置token失败", err, c.Ctx.Request.UserAgent())
			dto.Message = "设置token失败\t" + err.Error()
			return
		}

		dto.Message = "注册成功"
		dto.Data = &UserProfileInfo{UserProfile: *up, ImTk: up.ImToken}
		dto.Sucess = true
	} else {
		if u.AcctStatus != models.AcctStatusShield {
			if IsCheckMode4Context(c.Ctx) {
				(&models.UserProfile{ID: u.ID}).SetFaker()
			}

			if err := (&models.UserProfile{ID: u.ID}).UpdateOnlineStatus(models.OnlineStatusOnline); err != nil {
				beego.Error("更新在线状态失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "更新在线状态失败\t" + err.Error()
				return
			}

			up := &models.UserProfile{ID: u.ID}
			if err := up.Read(); err != nil {
				beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "获取用户信息失败\t" + err.Error()
				return
			}

			if dto.Data, err = genSelfUserPorfileInfo(up, nil); err != nil {
				beego.Error("获取用户信息失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "获取用户信息失败\t" + err.Error()
				return
			}

			tk.ID = u.ID
			tk.AcctStatus = u.AcctStatus
			tk.UserType = up.UserType
			if err := SetToken(c.Ctx, tk); err != nil {
				beego.Error("设置token失败", err, c.Ctx.Request.UserAgent())
				dto.Message = "设置token失败\t" + err.Error()
				return
			}

			dto.Sucess = true
			dto.Message = "登陆成功"
		} else {
			dto.Message = "您因涉及违规，已被封号。可联系平台运营"
			dto.Code = utils.DtoStatusAccountDisable
		}
	}
}

//InitUser .
func (c *AuthController) InitUser(u *models.User, acctType int) (*models.UserProfile, error) {
	trans := models.TransactionGen()

	if u.InvitedBy != 0 {
		if err := (&models.User{ID: u.InvitedBy}).Read(); err != nil { //检测邀请人是否存在
			u.InvitedBy = 0
		} else {
			(&models.UserExtra{ID: u.InvitedBy}).AddInviteCount(trans)
		}
	}

	u.AcctType = acctType
	u.AcctStatus = models.AcctStatusNormal
	u.CreatedAt = time.Now().Unix()
	uli := &models.UserLoginInfo{UserAgent: c.Ctx.Request.UserAgent(), IPAddress: c.Ctx.Input.IP(), Time: time.Now().Unix()}
	u.LastLoginInfo, _ = utils.JSONMarshalToString(uli)

	if err := u.Add(trans); err != nil {
		models.TransactionRollback(trans)
		return nil, errors.New("添加用户失败\t" + err.Error())
	}

	imUser := &netease.ImUser{ID: strconv.FormatUint(u.ID, 10)}
	imtk, err := utils.ImCreateUser(imUser)
	if err != nil {
		models.TransactionRollback(trans)
		return nil, errors.New("创建IMUser失败\t" + err.Error())
	}

	up := &models.UserProfile{ID: u.ID}
	index := common.RandNumber(0, len(randomName))
	up.Alias = randomName[index] + strconv.FormatUint(u.ID, 10) //随机生成花名
	up.ImToken = imtk.Token
	up.Birthday = 820425600
	up.CoverPic = `{"cover_pic_info": {"image_url": "` + defaultBoysAvatar + `", "cloud_porn_check": true}}`
	up.OnlineStatus = models.OnlineStatusOnline
	up.Description = "你不主动我们怎么会有故事"
	up.Location = "北京市"
	up.Price = 200
	if IsCheckMode4Context(c.Ctx) {
		up.UserType = models.UserTypeFaker
		up.Price = 0
	}

	if err := up.Add(trans); err != nil {
		models.TransactionRollback(trans)
		return nil, errors.New("创建用户详情失败\t" + err.Error())
	}

	if err := (&models.Subscribe{ID: u.ID}).Add(trans); err != nil {
		models.TransactionRollback(trans)
		return nil, errors.New("创建关注信息失败\t" + err.Error())
	}

	if err := (&models.UserExtra{ID: u.ID}).Add(trans); err != nil {
		models.TransactionRollback(trans)
		return nil, errors.New("创建用户附加信息失败\t" + err.Error())
	}

	if err := (&models.ProfileChg{ID: u.ID}).Add(trans); err != nil {
		models.TransactionRollback(trans)
		return nil, errors.New("创建用户资料变动信息失败\t" + err.Error())
	}

	models.TransactionCommit(trans)
	return up, nil
}

//Logout .
// @Title 注销登录
// @Description 注销登录
// @Success 200 {object} utils.ResultDTO
// @router /logout [get]
func (c *AuthController) Logout() {
	dto, tk := utils.ResultDTO{Message: "退出登陆成功"}, GetToken(c.Ctx)
	defer dto.JSONResult(&c.Controller)

	if err := (&models.UserProfile{ID: tk.ID}).UpdateOnlineStatus(models.OnlineStatusOffline); err != nil {
		beego.Error("更新在线状态失败", err, c.Ctx.Request.UserAgent())
		dto.Message = "更新在线状态失败\t" + err.Error()
		return
	}

	ClearToken(&c.Controller)
	dto.Sucess = true
}

func genSelfUserPorfileInfo(up *models.UserProfile, pc *models.ProfileChg) (*UserProfileInfo, error) { //获取用户自己信息时,给出审核状态，已经在审核状态的图片等信息
	upi := &UserProfileInfo{UserProfile: *up, ImTk: up.ImToken}

	upi.Alipay = &models.AlipayAcctInfo{}
	utils.JSONUnMarshal(upi.AlipayAcct, upi.Alipay) //忽略json解析错误

	cv := up.GetCover()
	genUserPorfileInfoCommon(upi, cv)

	if pc == nil {
		pc = &models.ProfileChg{ID: up.ID}
		if err := pc.Read(nil); err != nil {
			beego.Error("获取个人信息变动失败", err)
			return nil, err
		}
	}

	//当前用户获取自身信息时，给出正在审核的图片和视频
	if len(pc.CoverPic) > 0 && pc.CoverPicCheckStatus == models.CheckStatusUncheck {
		upi.Avatar = pc.CoverPic
	}

	if len(pc.Video) > 0 && pc.VideoCheckStatus == models.CheckStatusUncheck {
		upi.Video = pc.Video
	}
	upi.CheckStatus = pc

	if up.Gender == models.GenderMan {
		upi.IsFill = true
	} else if up.Gender == models.GenderWoman && len(up.Location) > 0 && up.Birthday > 0 && len(upi.Avatar) > 0 && len(up.Alias) > 0 {
		upi.IsFill = true
	} else {
		upi.IsFill = false
	}
	return upi, nil
}

func genUserPorfileInfoCommon(upi *UserProfileInfo, cv *models.UserCoverInfo) {
	if cv != nil {
		if cv.CoverPicture != nil {
			upi.Avatar = utils.TransCosToCDN(cv.CoverPicture.ImageURL)
		}

		if cv.DesVideo != nil {
			upi.Video = utils.TransCosToCDN(cv.DesVideo.VideoURL)
		}

		if cv.Gallery != nil && len(cv.Gallery) > 0 {
			var g []string
			for index := range cv.Gallery {
				g = append(g, utils.TransCosToCDN(cv.Gallery[index].ImageURL))
			}
			upi.Gallery = g
		}
	}

	if upi.DialAccept+upi.DialDeny > 0 {
		upi.AnswerRate = upi.DialAccept * 100 / (upi.DialAccept + upi.DialDeny) //计算接通率
	} else {
		upi.AnswerRate = 100
	}

	if upi.UserType == models.UserTypeFaker { //马甲号隐藏掉金额相关字段
		upi.Balance = 0
		upi.Income = 0
	}

	upi.Wallet = upi.Balance + upi.Income
}

func isInternalAcct(num string) bool {
	if faker.Users == nil || faker.Time.Before(time.Now()) {
		uarr, _ := (&models.User{}).GetInternalAcct()
		mp := make(map[string]models.User)

		for index := range uarr {
			u := uarr[index]
			if len(u.PhoneNumber) == 0 {
				continue
			}
			mp[u.PhoneNumber] = u
		}

		faker.Users = mp
		faker.Time = time.Now().Add(5 * time.Minute)
	}

	if _, ok := faker.Users[num]; ok {
		return true
	}
	return false
}

func addAcctLoginInfo(uid uint64, ctx *context.Context) {
	defer func() {
		if err := recover(); err != nil {
			beego.Error(err)
			debug.PrintStack()
		}
	}()

	ali := &models.AcctLoginInfo{UserID: uid, IPAddress: ctx.Input.IP(), Agent: ctx.Request.UserAgent(), Time: time.Now().Unix()}

	ii, err := utils.GetIPInfo(ali.IPAddress)
	if err != nil {
		beego.Error("获取IP信息失败:", ali.IPAddress)
	}

	ali.IPInfo, _ = utils.JSONMarshalToString(ii)

	if err := ali.Add(nil); err != nil {
		beego.Error("添加登陆信息失败:", err, ali)
	}
}

//SendActivity 发送促活消息
func SendActivity(uid uint64) error {
	if actives.Time.Before(time.Now()) {
		newActives := &CacheActiveInfo{Time: time.Now().Add(5 * time.Minute)}

		acts, err := (models.Active{}).GetActive()
		if err != nil {
			return err
		}

		for index := range acts {
			acti := ActiveInfo{Active: acts[index]}
			utils.JSONUnMarshal(acts[index].Content, &acti.Detail)
			newActives.Actives = append(newActives.Actives, acti)
		}

		actives = newActives
	}

	tmp, startTime := make([]ActiveInfo, len(actives.Actives)), time.Now()
	copy(tmp, actives.Actives)

	for len(tmp) > 0 {
		i := 0
		for _, x := range tmp {
			if time.Now().After(startTime.Add(time.Duration(x.DelayTime) * time.Second)) {
				var err error
				fromid, toid := strconv.FormatUint(x.UserID, 10), strconv.FormatUint(uid, 10)
				switch x.Type {
				case models.ActiveTypeMessage:
					err = utils.SendP2PMessage(fromid, toid, x.Detail.Message)
				case models.ActiveTypeImage:
					err = utils.SendP2PImageMessage(x.Detail.FileURL, fromid, []string{toid})
				case models.ActiveTypeVoice:
					err = utils.SendP2PVoiceMessage(x.Detail.FileURL, x.Detail.Duration*1000, fromid, []string{toid})
				case models.ActiveTypeVideo:
					err = utils.SendP2PVideoMessage(x.Detail.FileURL, x.Detail.Duration*1000, fromid, []string{toid})
				}
				if err != nil {
					beego.Error("发送促活消息失败 active_id:", x.ID, "to:", toid, "from:", fromid, err)
				}
			} else {
				// copy and increment index
				tmp[i] = x
				i++
			}
		}
		time.Sleep(time.Second)
		tmp = tmp[:i]
	}

	return nil
}
