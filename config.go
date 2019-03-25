package itea

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"context"
	"strings"
)

const (
	ENV 		= "env"
	SEARCH_ENV  = "{env}"
	DEFAULT_ENV	= "dev"
)

type Config struct {
	projectPath 	string
	env 			string
	conf 			map[string]*json.RawMessage
}

func InitConf(ctx context.Context, appConfig string) *Config {
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	env := ctx.Value(ENV).(string)
	dat, err := ioutil.ReadFile(projectPath + strings.Replace(appConfig, SEARCH_ENV, env, -1))
	if err != nil {
		panic("Application config not find")
	}
	var conf map[string]*json.RawMessage
	err = json.Unmarshal(dat, &conf)
	if err != nil {
		panic("Application config extract error")
	}
	return &Config{
		projectPath: projectPath,
		env: env,
		conf: conf,
	}
}

func (c *Config) Path(name string) string{
	if v, ok := c.conf[name]; ok {
		var s string
		json.Unmarshal(*v, &s)
		return c.projectPath + strings.Replace(s, SEARCH_ENV, c.env, -1)
	}
	return ""
}

func (c *Config) PathList(name string) []string{
	if v, ok := c.conf[name]; ok {
		var sl []string
		json.Unmarshal(*v, &sl)
		if sl != nil {
			for i, _ := range sl {
				sl[i] = c.projectPath + strings.Replace(sl[i], SEARCH_ENV, c.env, -1)
			}
		}
		return sl
	}
	return nil
}

func (c *Config) Beans(name string) []Bean{
	if v, ok := c.conf[name]; ok {
		var beans []Bean
		json.Unmarshal(*v, &beans)
		return beans
	}
	return nil
}
