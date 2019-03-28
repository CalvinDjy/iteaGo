package itea

import (
	"encoding/json"
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
	conf 			map[string]*json.RawMessage
	confMap			map[string]interface{}
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
	var conf map[string]*json.RawMessage
	err = json.Unmarshal(dat, &conf)
	if err != nil {
		panic("Application config extract error")
	}
	config := &Config{
		projectPath: projectPath,
		Env: env,
		conf: conf,
		confMap: make(map[string]interface{}),
	}
	config.wg.Add(3)
	go config.dbConfig()
	go config.importConfig()
	go config.logConfig()
	config.wg.Wait()
	return config
}

//Env
func env() string {
	num := len(os.Args)
	if num > 1 {
		return os.Args[1]
	}
	return DEFAULT_ENV
}

//Extract database config
func (c *Config) dbConfig() {
	defer c.wg.Done()
	if v, ok := c.conf[DB_CONFIG]; ok {
		var s string
		json.Unmarshal(*v, &s)
		path := c.projectPath + strings.Replace(s, SEARCH_ENV, c.Env, -1)
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			panic("[" + path + "] config not find")
		}
		var databases map[string]*json.RawMessage
		err = json.Unmarshal(dat, &databases)
		if err != nil {
			panic(err)
		}
		for k, v := range databases {
			mutex.Lock()
			c.confMap[k] = v
			mutex.Unlock()
		}
	}
}

//Extract import config
func (c *Config) importConfig() {
	defer c.wg.Done()
	if v, ok := c.conf[IMPORT_CONFIG]; ok {
		var pathList []string
		json.Unmarshal(*v, &pathList)
		if pathList != nil {
			for i, _ := range pathList {
				pathList[i] = c.projectPath + strings.Replace(pathList[i], SEARCH_ENV, c.Env, -1)
			}
			mutex.Lock()
			c.confMap[IMPORT_CONFIG] = pathList
			mutex.Unlock()
		}
	}
}

//Extract log config
func (c *Config) logConfig() {
	defer c.wg.Done()
	if v, ok := c.conf[LOG_CONFIG]; ok {
		var log map[string]interface{}
		json.Unmarshal(*v, &log)
		mutex.Lock()
		c.confMap[LOG_CONFIG] = log
		mutex.Unlock()
	}
}

func (c *Config) Config(name string) interface{} {
	if v, ok := c.confMap[name]; ok {
		return v
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
