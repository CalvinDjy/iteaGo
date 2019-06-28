package itea

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"strings"
)

func InitLog() {
	logtype, logfile, rotate := "", "", false
	if conf.Config(LOG_CONFIG) != nil {
		logConf := conf.Config(LOG_CONFIG).(Log)
		if !strings.EqualFold(logConf.Type, "") {
			logtype = logConf.Type
		}
		if !strings.EqualFold(logConf.Logfile, "") {
			logfile = logConf.Logfile
		}
		rotate = logConf.Rotate
	}
	ilog.Init(logtype, logfile, rotate)
}
