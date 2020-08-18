package db

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/kataras/golog"
	_ "github.com/go-sql-driver/mysql"
	"fabric-client/inits/parse"
	"sync"
	"time"
)

var (
	masterEngine *xorm.Engine
	slaveEngine  *xorm.Engine
	lock         sync.Mutex
)

func MasterEngine() *xorm.Engine {
	if masterEngine != nil {
		return masterEngine
	}

	lock.Lock()
	defer lock.Unlock()

	if masterEngine != nil {
		return masterEngine
	}

	masterDB := parse.DB.MasterDB

	engine, err := xorm.NewEngine(masterDB.Dialect, GetDBConnURL(&masterDB))
	if err != nil {
		golog.Fatalf("连接主数据库失败:%s", err)
		return nil
	}
	settings(engine, &masterDB)
	masterEngine = engine
	return masterEngine
}

func SlaveEngine() *xorm.Engine {
	if slaveEngine != nil {
		return slaveEngine
	}

	lock.Lock()
	defer lock.Unlock()

	masterDB := parse.DB.MasterDB

	engine, err := xorm.NewEngine(masterDB.Dialect, GetDBConnURL(&masterDB))
	if err != nil {
		golog.Fatalf("连接从数据库失败:%s", err)
		return nil
	}
	settings(engine, &masterDB)
	slaveEngine = engine
	return slaveEngine
}


func settings(engine *xorm.Engine, info *parse.DBYamlConfig) {
	engine.ShowSQL(info.ShowSql)
	localTime, _ := time.LoadLocation("Asia/Shanghai")
	engine.SetTZLocation(localTime)
	if info.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(info.MaxIdleConns)
	}
	if info.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(info.MaxOpenConns)
	}
}

func GetDBConnURL(info *parse.DBYamlConfig) (url string) {
	url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		info.User,
		info.Password,
		info.Host,
		info.Port,
		info.Database,
		info.Charset)
	return url
}
