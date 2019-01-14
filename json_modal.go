package itea

type Bean struct {
	Name string 					`json:"name"`
	Class string 					`json:"class"`
	ExecuteMethod string			`json:"execute"`
	Params map[string]interface{}	`json:"construct-params"`
}

type DatabaseConf struct {
	Driver string					`json:"driver"`
	Ip string						`json:"ip"`
	Port string						`json:"port"`
	Database string					`json:"database"`
	Username string					`json:"username"`
	Password string					`json:"password"`
	Charset string					`json:"charset"`
	MaxConn int						`json:"maxConn"`
	MaxIdle int						`json:"maxIdle"`
	ConnMaxLift int					`json:"connMaxLift"`
}

type RedisConf struct {
	Host string						`json:"host"`
	Port string 					`json:"port"`
	Database int					`json:"database"`
	Password string 				`json:"password"`
	MaxIdle int 					`json:"max_idle"`
	MaxActive int					`json:"max_active"`
	IdleTimeout int 				`json:"idle_timeout"`
	MaxConnLifetime int 			`json:"max_conn_lifetime"`
	IdleCheck int					`json:"idle_check"`
}