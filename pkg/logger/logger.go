/*
@author: panfengguo
@since: 2024/11/9
@desc: desc
*/
package logger

import (
	"github.com/AgentGuo/spike/cmd/server/config"
	"github.com/sirupsen/logrus"
	"io"
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

		if config.GetConfig().LogToFile {
			currentTime := time.Now().Format("2006-01-02_15-04-05")
			logFileName := currentTime + ".log"
			// 创建日志文件
			file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Warn("Failed to log to file, using default stderr")
			} else {
				mw := io.MultiWriter(os.Stdout, file)
				log.SetOutput(mw)
			}
		}
	})
	return log
}
