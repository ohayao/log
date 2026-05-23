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

type Test01 struct {
	ID   int64
	Name string
}

func init() {
	opts := []log.Option{
		log.WithColor(true),
		log.WithShortName(true),
		log.WithMinLevel(log.LV_DEBUG),
	}
	logger, err := log.NewFileLogger("./output/log.log", 3000, opts...)
	if err != nil {
		panic(err)
	}
	log.DEFAULT = logger
}

func initDb() {
	gormLogger := log.NewGormLogger(log.DEFAULT, true, time.Millisecond*50)
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormLogger.LogMode(logger.Info),
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
	initDb()
	log.Errorf("Test: %d", time.Now().Unix())
	var list []Test01
	if err := DB.Model(&Test01{}).Find(&list).Error; err != nil {
		log.Error(err)
	} else {
		log.Infof("list lenght: %d", len(list))
	}
	_ = DB.Model(&Test01{}).Where("names = ''").Find(&list).Error
}
