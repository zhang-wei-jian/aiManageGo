package models

import (
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var err error

var Rdb *redis.Client

func init() {

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	//dsn := "root:1379@tcp(127.0.0.1:3306)/gin?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "airBlog:137930@tcp(60.204.170.225)/airBlog?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		println(err)
	}

	//连接redis
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "60.204.170.225:6379",        //连接地址
		Password: "jhkdjhkjdhsIUTYURTU_352CcZ", //连接密码
		DB:       0,                            //默认连接库
		PoolSize: 100,                          //连接池大小
	})

	println("数据库链接成功", Rdb)

}
