package itea

type Application struct {
	Process []Process				`yaml:"process"`
	Database string					`yaml:"database"`
	Import []string					`yaml:"import"`
	Log Log							`yaml:"log,omitempty"`
}

type Process struct {
	Name string 					`yaml:"name"`
	Class string 					`yaml:"class"`
	ExecuteMethod string			`yaml:"execute"`
	Params map[string]interface{}	`yaml:"construct-params"`
}

type Log struct {
	Type string						`yaml:"type"`
	Logfile string					`yaml:"logfile"`
	Rotate bool						`yaml:"rotate"`
}

type StorageConf struct {
	Connections map[string]DatabaseConf	`yaml:"connections"`
	Redis RedisConf						`yaml:"redis"`
}

type DatabaseConf struct {
	Driver string					`yaml:"driver"`
	Ip string						`yaml:"ip"`
	Port string						`yaml:"port"`
	Database string					`yaml:"database"`
	Username string					`yaml:"username"`
	Password string					`yaml:"password"`
	Charset string					`yaml:"charset"`
	MaxConn int						`yaml:"maxConn"`
	MaxIdle int						`yaml:"maxIdle"`
	ConnMaxLift int					`yaml:"connMaxLift"`
}

type RedisConf struct {
	Host string						`yaml:"host"`
	Port string 					`yaml:"port"`
	Database int					`yaml:"database"`
	Password string 				`yaml:"password"`
	MaxIdle int 					`yaml:"max_idle"`
	MaxActive int					`yaml:"max_active"`
	IdleTimeout int 				`yaml:"idle_timeout"`
	MaxConnLifetime int 			`yaml:"max_conn_lifetime"`
	IdleCheck int					`yaml:"idle_check"`
}