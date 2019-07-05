package itea

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

const (
	SEARCH_ENV  = "{env}"
	DEFAULT_ENV	= "dev"
)

type Config struct {
	Env 			string
	projectPath 	string
	appConf 		Application
	config			map[string]interface{}
	wg				sync.WaitGroup
}

func InitConf(appConfig string) *Config {
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	env := env()
	dat, err := ioutil.ReadFile(projectPath + strings.Replace(appConfig, SEARCH_ENV, env, -1))
	if err != nil {
		panic("Application config not find")
	}
	var appConf Application
	err = yaml.Unmarshal(dat, &appConf)
	if err != nil {
		panic("Application config extract error")
	}
	config := &Config{
		projectPath: projectPath,
		Env: env,
		appConf: appConf,
		config: make(map[string]interface{}),
	}
	config.wg.Add(3)
	go config.dbConfig()
	go config.logConfig()
	go config.importConfig()
	config.wg.Wait()
	return config
}

//Env
func env() string {
	num := len(os.Args)
	if num <= 1 {
		return DEFAULT_ENV
	}
	for i, arg := range os.Args {
		if !strings.EqualFold(arg, "-e") {
			continue
		}
		if i + 1 <= num - 1 {
			return os.Args[i + 1]
		}
	}
	return DEFAULT_ENV
}

//Extract database config
func (c *Config) dbConfig() {
	defer c.wg.Done()
	if !strings.EqualFold(c.appConf.Database, "") {
		s := c.appConf.Database
		path := c.projectPath + strings.Replace(s, SEARCH_ENV, c.Env, -1)
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			panic("[" + path + "] config not find")
		}
		var databases StorageConf
		err = yaml.Unmarshal(dat, &databases)
		if err != nil {
			panic(err)
		}
		c.config[CONNECTION_CONFIG] = databases.Connections
		c.config[REDIS_CONFIG] = databases.Redis
	}
}

//Extract log config
func (c *Config) logConfig() {
	defer c.wg.Done()
	if !strings.EqualFold(c.appConf.Log.Type, "") {
		c.config[LOG_CONFIG] = c.appConf.Log
	}
}

//Extract import config
func (c *Config) importConfig() {
	defer c.wg.Done()
	if len(c.appConf.Import) > 0 {
		pathList := c.appConf.Import
		for i, _ := range pathList {
			pathList[i] = c.projectPath + strings.Replace(pathList[i], SEARCH_ENV, c.Env, -1)
		}
		c.config[IMPORT_CONFIG] = pathList
	}
}

func (c *Config) Config(name string) interface{} {
	if v, ok := c.config[name]; ok {
		return v
	}
	return nil
}

func (c *Config) Beans(name string) []Process{
	return c.appConf.Process
}
