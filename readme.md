# log
> a golang log library
> 

## 示例
``` go
func testDefault() {
	log.UseOption(log.DEFAULT,
		log.WithColor(true),
		log.WithShortName(true),
		log.WithTimeStyle(log.FLAG_TIME_DATETIME),
		log.WithMinLevel(log.LV_DEBUG),
	)

	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Errorf("recover %v", err))
		}
	}()

	log.Info("hello world")
	log.Warnf("hello world, it's %d", time.Now().Unix())
	log.Debugln("hello world")
	log.Error("hello world")

	if time.Now().Unix()%2 == 0 {
		log.Panic("hello world")
	} else {
		log.Fatal("fatal")
	}
}
```