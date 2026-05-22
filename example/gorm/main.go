package main

import (
	"fmt"
	"time"

	"github.com/ohayao/log/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var LG *log.Logger

type Test01 struct {
	ID   int64
	Name string
}

func init() {
	LG, _ = log.NewFileLogger("./logx.log", 1000*1000*1,
		log.WithColor(true),
		log.WithShortName(true),
		log.WithMinLevel(log.LV_DEBUG))
	loger := log.NewGormLogger(LG, true, time.Millisecond*30)
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: loger.LogMode(logger.Info),
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		"root",
		"123456",
		"localhost",
		3306,
		"test",
		"charset=utf8&parseTime=True&loc=Local",
	)
	if db, err := gorm.Open(mysql.Open(dsn), &gormConfig); err != nil {
		panic(err)
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(140)
		sqlDB.SetConnMaxLifetime(time.Second * 30)
		DB = db
	}
}

func main() {
	LG.Errorf("Test: %d", time.Now().Unix())
	var list []Test01
	if err := DB.Model(&Test01{}).Find(&list).Error; err != nil {
		LG.Error(err)
	} else {
		LG.Infof("list lenght: %d", len(list))
	}
	_ = DB.Model(&Test01{}).Where("names = ''").Find(&list).Error
}
