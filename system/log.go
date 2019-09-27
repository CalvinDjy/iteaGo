package system

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"strings"
	"fmt"
)

const LOG_KEY = "log"

type Log struct {
	Type 		string
	Logfile 	string
	Rotate 		bool
}

func InitLog() {
	logtype, logfile, rotate := "", "", false
	if c := Conf.GetStruct(fmt.Sprintf("%s.%s", Conf.FileName, LOG_KEY), Log{}); c != nil {
		logConf := c.(*Log)
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
