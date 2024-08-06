package main

import (
	"os"

	"github.com/ohayao/log"
	// gormlogger "gorm.io/gorm/logger"
)

func main() {
	// init handler
	sh := log.NewStreamHandler(os.Stdout)
	// init gormloghandler with handler
	glh := log.NewGormLoggerHandler(sh)
	// init logger
	logger := log.NewLogger(glh)
	// set basic logger
	logger.SetLevels(log.LevelAll)
	logger.SetFlags(log.FlagLevel, log.FlagColor)
	// you can also set gormlogger
	glh.SetColorful(true)
	glh.SetIgnoreRecordNotFoundError(true)

	// ....
	// mysql config set

	// db, err = gorm.Open(mysql.New(mysql.Config{
	// 	Conn: sqldb,
	// }), &gorm.Config{
	// 	Logger: glh.LogMode(gormlogger.Info).(*log.GormLoggerHandler).SetSlowThreshold(time.Millisecond * 3),
	// })
}
