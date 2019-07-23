package itea

import (
	"flag"
	"fmt"
	"github.com/goinggo/mapstructure"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
)

const (
	SEARCH_ENV  	= "{env}"
	DEFAULT_ENV		= "dev"
	IMPORT_KEY		= "import"
	DATABASE_KEY	= "database"
)

var (
	Help		bool
	Start 		bool
	Stop 		bool
	Env 		string	//Environment
	projpath 	string	//Application proj base path
)

func init ()  {
	flag.BoolVar(&Help, "h", false, "Get help")
	flag.BoolVar(&Start, "start", true, "Start application")
	flag.BoolVar(&Stop, "stop", false, "Stop application")
	flag.StringVar(&Env, "e", DEFAULT_ENV, "Set application environment")
	flag.Parse()
	if Help {
		fmt.Fprintf(os.Stderr, `iteaGo version: iteaGo/%s
Usage: main [-start|-stop] [-e env]
Options:
`, ITEAGO_VERSION)
		flag.PrintDefaults()
	}
}

//Get file path
func filePath(f string) string {
	return projpath + strings.Replace(f, SEARCH_ENV, Env, -1)
}

//Get file name
func fileName(p string) string {
	filenameWithSuffix := path.Base(p)
	fileSuffix := path.Ext(filenameWithSuffix)
	return strings.TrimSuffix(filenameWithSuffix, fileSuffix)
}

//Find config
func find(k []string, l int, conf map[interface{}]interface{}) interface{} {
	if l == 1 {
		return conf[k[0]]
	}
	if c, ok := conf[k[0]];ok {
		l--
		return find(k[1:], l, c.(map[interface{}]interface{}))
	} else {
		return nil
	}
}

func decode(v interface{}, t reflect.Type) (interface{}, error){
	ins := reflect.New(t).Interface()
	if err := mapstructure.Decode(v, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

type Config struct {
	FileName 		string
	config			map[interface{}]interface{}
}

func InitConf(file string) *Config {
	var err error
	projpath, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	dat, err := ioutil.ReadFile(filePath(file))
	if err != nil {
		panic("Application config not find")
	}
	var application map[interface{}]interface{}
	err = yaml.Unmarshal(dat, &application)
	if err != nil {
		panic("Application config extract error")
	}

	FileName := fileName(file)
	c := &Config{
		FileName: FileName,
		config: make(map[interface{}]interface{}),
	}
	c.config[FileName] = application
	
	var wg sync.WaitGroup
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.importConfig()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.dbConfig()
	}()
	
	wg.Wait()

	return c
}

//Extract database config
func (c *Config) dbConfig() {
	if f := c.GetString(fmt.Sprintf("%s.%s", c.FileName, DATABASE_KEY));!strings.EqualFold(f, "") {
		dat, err := ioutil.ReadFile(filePath(f))
		if err != nil {
			panic("database config not find")
		}
		var databases map[interface{}]interface{}
		err = yaml.Unmarshal(dat, &databases)
		if err != nil {
			panic(err)
		}
		c.config[DATABASE_KEY] = databases
	}
}

//Extract import config
func (c *Config) importConfig() {
	imp := c.GetArray(fmt.Sprintf("%s.%s", c.FileName, IMPORT_KEY))
	if len(imp) <= 0 {
		return
	}
	
	l := len(imp)
	
	ch := make(chan []interface{}, l)
	
	for _, f := range imp {
		go func(f string) {
			dat, err := ioutil.ReadFile(filePath(f))
			if err != nil {
				ch <- nil
			}
			var conf map[interface{}]interface{}
			yaml.Unmarshal(dat, &conf)
			ch <- []interface{}{
				fileName(f), conf,
			}
		}(f.(string))
	}
	
	for i := 0; i < l; i++ {
		v := <-ch
		c.config[v[0].(string)] = v[1]
	}
}

func (c *Config) value(key string) interface{} {
	arr := strings.Split(key, ".")
	l := len(arr)
	return find(arr, l, c.config)
}

//Get string value
func (c *Config) GetString(key string) string {
	v := c.value(key)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

//Get config array
func (c *Config) GetArray(key string) []interface{} {
	v := c.value(key)
	if v == nil {
		return nil
	}
	if array, ok := v.([]interface{}); ok {
		return array
	}
	return nil
}

func (c *Config) GetStruct(key string, s interface{}) interface{} {
	v := c.value(key)
	if v == nil {
		return nil
	}
	ins, err := decode(v, reflect.TypeOf(s))
	if err != nil {
		fmt.Println("GetStruct error : ", err)
	}
	return ins
}

func (c *Config) GetStructArray(key string, s interface{}) []interface{} {
	v := c.value(key)
	if v == nil {
		return nil
	}
	if av, ok := v.([]interface{}); ok {
		var list []interface{}
		t := reflect.TypeOf(s)
		for _, item := range av {
			ins, err := decode(item, t)
			if err != nil {
				fmt.Println("GetStructArray error : ", err)
				continue
			}
			list = append(list, ins)
		}
		return list
	}
	return nil
}

func (c *Config) GetStructMap(key string, s interface{}) map[string]interface{} {
	v := c.value(key)
	if v == nil {
		return nil
	}
	if mv, ok := v.(map[interface{}]interface{}); ok {
		m := make(map[string]interface{})
		t := reflect.TypeOf(s)
		for k, item := range mv {
			ins, err := decode(item, t)
			if err != nil {
				fmt.Println("GetStructMap error : ", err)
				continue
			}
			m[k.(string)] = ins
		}
		return m
	}
	return nil
}

func String(key string) string {
	return config.GetString(key)
}

func Array(key string) []interface{} {
	return config.GetArray(key)
}

func Struct(key string, s interface{}) interface{} {
	return config.GetStruct(key, s)
}

func StructArray(key string, s interface{}) []interface{} {
	return config.GetStructArray(key, s)
}

func StructMap(key string, s interface{}) map[string]interface{} {
	return config.GetStructMap(key, s)
}