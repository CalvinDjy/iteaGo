package itea

import (
	"github.com/CalvinDjy/iteaGo/ilog"
)

func InitLog() {
	logfile, rotate := "", false
	if logConf, ok := conf.Config(LOG_CONFIG).(map[string]interface{}); ok {
		if v, ok := logConf["logfile"]; ok {
			logfile = v.(string)
		}
		if v, ok := logConf["rotate"]; ok {
			rotate = v.(bool)
		}
	}
	ilog.Init(logfile, rotate)
}
