package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/i18n"
	"github.com/kataras/iris/v12/mvc"
	_ "fabric-client/inits"
	"fabric-client/sdkInit"
	"fabric-client/service"
	"fabric-client/web/controllers"
	"log"
)

var clientMap map[string]*sdkInit.Client

func main() {
	var err error
	clientMap, err = sdkInit.InitClientMap()
	service.ClientMap=clientMap
	defer sdkInit.CloseClientMap(clientMap)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	app := iris.New()

	app.Logger().SetLevel("debug")

	app.Use(i18n.New(i18n.Config{
		Default:      "en",
		URLParameter: "lang",
		Languages: map[string]string{
			"en": "./locale/locale_en-US.ini",
			"zh": "./locale/locale_zh-CN.ini"}}))

	mvcApp := mvc.New(app.Party("/api"))
	mvcApp.Register(clientMap)
	mvcApp.Handle(new(controllers.FabricSDKController))

	// 启动服务
	err = app.Run(
		iris.Addr(":8080"),                            // 地址
		iris.WithCharset("UTF-8"),                     // 国际化
		iris.WithOptimizations,                        // 自动优化
		iris.WithoutServerError(iris.ErrServerClosed), // 忽略框架错误
	)

	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
}
