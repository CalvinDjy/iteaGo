package itea

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	projectPath string
	conf map[string]*json.RawMessage
}

func InitConf(appConfig string) *Config {
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dat, err := ioutil.ReadFile(projectPath + appConfig)
	if err != nil {
		panic("Application config not find")
	}
	var conf map[string]*json.RawMessage
	err = json.Unmarshal(dat, &conf)
	if err != nil {
		panic("Application config extract error")
	}
	return &Config{
		projectPath:projectPath,
		conf:conf,
	}
}

func (c *Config) Path(name string) string{
	if v, ok := c.conf[name]; ok {
		var s string
		json.Unmarshal(*v, &s)
		return c.projectPath + s
	}
	return ""
}

func (c *Config) PathList(name string) []string{
	if v, ok := c.conf[name]; ok {
		var sl []string
		json.Unmarshal(*v, &sl)
		if sl != nil {
			for i, _ := range sl {
				sl[i] = c.projectPath + sl[i]
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
