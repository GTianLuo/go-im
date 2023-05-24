package dbConf

import "C"
import (
	"go-im/common/conf"
	"strings"
)

func InitDbService() {

	//配置mysql
	readPath := strings.Join([]string{conf.V.GetString("mysql.dbUser"), ":", conf.V.GetString("mysql.dbPassword"), "@tcp(", conf.V.GetString("mysql.dbHost"), ":", conf.V.GetString("mysql.dbPort"), ")/", conf.V.GetString("mysql.dbName"), "?charset=utf8&parseTime=true"}, "")
	writePath := ""
	database(readPath, writePath)

	//配置redis
	cache(conf.V.GetString("redis.addr"), conf.V.GetString("redis.password"), conf.V.GetInt("redis.db"))

}
