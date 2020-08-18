package controllers

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/i18n"
)

type Result struct {
	Code    int
	Message string
	Data    interface{}
}

const (
	OK                     = 0  //成功
	ParseParamsError       = 1  //解析参数错误
	SignExpiredError       = 2  //签名过期
	SignInvalidError       = 3  //签名错误
	GetAndCheckClientError = 4  //获取并检查客户端错误
	CreateChannelError     = 5  //创建通道失败
	JoinChannelError       = 6  //加入通道失败
	InstallCCError         = 7  //安装链码失败
	InstantiateCCError     = 8  //初始化失败
	UpgradeCCError         = 9  //更新失败
	ArgsError              = 10 //参数或者参数长度错误
	NewChannelClientError  = 11 //新建通道客户端错误
	ExecCCError            = 12 //执行失败
	QueryCCError           = 13 //查询失败
	NewLedgerClientError   = 14 //新建账本客户端错误
	QueryBlockError        = 15 //查询block失败
	QueryBlockByIdError    = 16 //根据txid查询block失败
)

func parseJson(ctx iris.Context, jsonObjectPtr interface{}) Result {
	err := ctx.ReadJSON(jsonObjectPtr)
	if err != nil {
		return getBadRequestResult(ctx, ParseParamsError, i18n.Translate(ctx, "parse_params_fail"), err.Error())
	}

	return Result{Code: OK}
}

func checkSign(ctx iris.Context, timestamp int64, sign string, src string) Result {
	currentTimestamp := time.Now().Unix()
	ctx.Application().Logger().Println(currentTimestamp)
	if currentTimestamp-timestamp > 120 {
		return getInternalServerError(ctx, SignExpiredError, i18n.Translate(ctx, "sign_expired"), nil)
	}

	currentSign := getSign(src)
	ctx.Application().Logger().Println(currentSign + ":" + sign)
	if currentSign != sign {
		return getInternalServerError(ctx, SignInvalidError, i18n.Translate(ctx, "sign_invalid"), src+":"+currentSign+":"+sign)
	}

	return Result{Code: OK}
}

func getSign(src string) string {
	sha512Bytes := sha512.Sum512([]byte(src))
	sha512 := hex.EncodeToString(sha512Bytes[:])
	md5Bytes := md5.Sum([]byte(sha512))
	sign := hex.EncodeToString(md5Bytes[:])
	return sign
}

func getBadRequestResult(ctx iris.Context, code int, message string, data interface{}) Result {
	ctx.StatusCode(iris.StatusBadRequest)
	return Result{Code: code, Message: message, Data: data}
}

func getInternalServerError(ctx iris.Context, code int, message string, data interface{}) Result {
	ctx.StatusCode(iris.StatusInternalServerError)
	return Result{Code: code, Message: message, Data: data}
}
