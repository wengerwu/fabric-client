package inits

import (
	"github.com/kataras/golog"
	"fabric-client/inits/parse"
	"fabric-client/util"
)

//init 初始化配置文件
func init() {
	//初始化数据库配置
	err := util.ReadYamlConfig("config/db.yaml",&parse.DB)
	if err!=nil {
		golog.Errorf("初始化配置文件错误：%s", err)
	}
	golog.Infof("db.yaml配置文件解析:%v",parse.DB)

}

