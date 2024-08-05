package main

import (
	"context"
	"os"

	"github.com/ohayao/log"
	// gormlogger "gorm.io/gorm/logger"
)

func main() {
	sh := log.NewStreamHandler(os.Stdout)
	logger := log.NewLogger(sh)
	logger.SetLevels(log.LevelAll)
	glogger := log.WrapLoggerForGorm(context.TODO(), logger)
	glogger.Logger.SetFlags(log.FlagTime, log.FlagColor, log.FlagLevel, log.FlagLongFile)
	// db, err = gorm.Open(mysql.New(mysql.Config{
	// 	Conn: sqldb,
	// }), &gorm.Config{
	// 	Logger: glogger.LogMode(gormlogger.Info),
	// })
}
