package itea

import (
	"github.com/CalvinDjy/iteaGo/ilog"
)

func InitLog() {
	logtype, logfile, rotate := "", "", false
	if logConf, ok := conf.Config(LOG_CONFIG).(map[string]interface{}); ok {
		if v, ok := logConf["type"]; ok {
			logtype = v.(string)
		}
		if v, ok := logConf["logfile"]; ok {
			logfile = v.(string)
		}
		if v, ok := logConf["rotate"]; ok {
			rotate = v.(bool)
		}
	}
	ilog.Init(logtype, logfile, rotate)
}
