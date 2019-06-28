package itea

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"os"
	"sync"
)

type routeConf struct {
	Groups []groupConf					`yaml:"groups"`
	ActionConf map[string]actionConf	`yaml:"action"`
}

type groupConf struct {
	Name string						`yaml:"name"`
	Prefix string					`yaml:"prefix"`
	Middleware string				`yaml:"middleware"`
}

type actionConf struct {
	Method string					`yaml:"method"`
	Uses string						`yaml:"uses"`
	Middleware string				`yaml:"middleware"`
	Group string					`yaml:"group"`
}

type Route struct {
	Groups		map[string]groupConf
	Actions 	map[string]*action
	mutex 		*sync.Mutex
	wg			sync.WaitGroup
}

type action struct {
	Method 		string
	Controller 	string
	Action 		string
	Middleware  []string
}

func (r *Route) Init(routeConfig string, env string) (route *Route){
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(projectPath + strings.Replace(routeConfig, SEARCH_ENV, env, -1))
	if err != nil {
		panic("Route config not find")
	}
	var routeConf routeConf
	err = yaml.Unmarshal(data, &routeConf)
	if err != nil {
		panic("Route config extract fail")
	}
	r.Groups = make(map[string]groupConf)
	r.Actions = make(map[string]*action)
	r.mutex = new(sync.Mutex)
	for _, gConf := range routeConf.Groups {
		r.Groups[gConf.Name] = gConf
	}
	r.extract(routeConf.ActionConf)
	return r
}

func (r *Route) extract(actionConf map[string]actionConf) {
	for uri, conf := range actionConf {
		u, c := uri, conf
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			method, controller, deal, middleware := "get", "", "", []string{}
			uArray := strings.Split(u, " ")
			if len(uArray) == 2 {
				method = uArray[0]
				u = uArray[1]
			}
			if !strings.EqualFold(c.Method, "") {
				method = c.Method
			}
			if strings.EqualFold(c.Uses, "") {
				return
			}
			pathArray := strings.Split(c.Uses, "@")
			if len(pathArray) != 2 {
				return
			}
			controller, deal = pathArray[0], pathArray[1]
			if !strings.EqualFold(c.Group, "") {
				groupNames := strings.Split(c.Group, "|")
				for _, groupName := range groupNames {
					if group, ok := r.Groups[groupName]; ok {
						if !strings.EqualFold(group.Prefix, "") {
							u = group.Prefix + u
						}
						if !strings.EqualFold(group.Middleware, "") {
							middleware = append(middleware, strings.Split(group.Middleware, "|")...)
						}
					}
				}
			}
			if !strings.EqualFold(c.Middleware, "") {
				middleware = append(middleware, strings.Split(c.Middleware, "|")...)
			}
			r.mutex.Lock()
			r.Actions[u] = &action{
				Method:method,
				Controller:controller,
				Action:deal,
				Middleware:middleware,
			}
			r.mutex.Unlock()
		}()
	}
	r.wg.Wait()
}
