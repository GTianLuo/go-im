package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-im/common/dao"
	"go-im/common/model"
	"go-im/common/result"
	"go-im/common/util"
	"go-im/ipConfig/serviceManage"
)

type UserLogic struct {
	NickName string `form:"nickName"`
	Account  string `form:"account"`
	Password string `form:"password"`
}

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	user := &UserLogic{}
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(200, result.Fail(result.InvalidParam, err))
		return
	}

	if user.NickName == "" || user.Account == "" || user.Password == "" {
		ctx.JSON(200, result.Fail(result.InvalidParam, nil))
		return
	}

	userDao := dao.NewUserDao()

	//判断账号的唯一性
	if exist, err := userDao.UserIsExist(user.Account); err != nil {
		ctx.JSON(result.Ok, result.Fail(result.Error, err))
		return
	} else if exist {
		ctx.JSON(result.Ok, result.Fail(result.AccountHasExist, nil))
		return
	}
	if err := userDao.SaveUser(&model.User{NickName: user.NickName,
		Account: user.Account, Password: util.Encryption(user.Password)}); err != nil {
		ctx.JSON(200, result.Fail(result.Error, err))
		return
	}
	ctx.JSON(result.Ok, result.Success(result.Ok, true))
	return
}

// UserLogin 用户登陆
func UserLogin(ctx *gin.Context) {
	user := &UserLogic{}
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(200, result.Fail(result.InvalidParam, err))
		return
	}
	userDao := dao.NewUserDao()
	//校验账号密码
	exist, err := userDao.UserIsExist(user.Account)
	if err != nil {
		ctx.JSON(result.Ok, result.Fail(result.Error, err))
		return
	} else if !exist {
		ctx.JSON(result.Ok, result.Fail(result.AccountNotExist, nil))
		return
	}
	userModel, err := userDao.GetUserByAccount(user.Account)
	if err != nil {
		ctx.JSON(200, result.Fail(result.Error, err))
		return
	}
	user.NickName = userModel.NickName
	if util.Encryption(user.Password) != userModel.Password {
		ctx.JSON(result.Ok, result.Fail(result.WrongPassword, nil))
		return
	}
	user.NickName = userModel.NickName

	//获取登陆列表
	paths := serviceManage.DisPatch()
	//生成token
	token := uuid.NewString()
	//登陆状态保存至redis
	if err := userDao.SaveLoginStatus(user.Account, token, user.NickName); err != nil {
		ctx.JSON(result.Ok, result.Fail(result.Error, err))
		return
	}

	ctx.JSON(result.Ok, result.Success(result.Ok, result.UserResult{Token: token, NickName: user.NickName, IpList: paths}))
	return
}
