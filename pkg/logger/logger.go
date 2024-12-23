/*
@author: panfengguo
@since: 2024/11/9
@desc: desc
*/
package logger

import (
	"github.com/AgentGuo/spike/cmd/server/config"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var (
	log          *logrus.Logger
	logOnce      sync.Once
	logStringMap = map[string]logrus.Level{
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}
)

func GetLogger() *logrus.Logger {
	logOnce.Do(func() {
		log = logrus.New()
		if lvl, ok := logStringMap[config.GetConfig().LogLevel]; ok {
			log.SetLevel(lvl)
		} else {
			log.SetLevel(logrus.InfoLevel)
		}
		log.SetReportCaller(true)

		// 获取当前时间并格式化为字符串
		currentTime := time.Now().Format("2006-01-02_15-04-05")
		logFileName := currentTime + ".log"

		// 创建日志文件
		file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Warn("Failed to log to file, using default stderr")
		}
	})
	return log
}
