package util

import (
	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
	"os"
)

func ReadYamlConfig(path string,config interface{}) error {
	f,err:=os.Open(path)
	defer f.Close()
	if err!=nil{
		golog.Fatalf("打开文件失败:%s",err)
	}

	yaml.NewDecoder(f).Decode(config)
	return nil
}