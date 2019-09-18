package itea

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"os"
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
	Actions 	[]*action
}

type action struct {
	Uri 		string
	Method 		string
	Controller 	string
	Action 		string
	Middleware  []string
}

func (r *Route) InitRoute(routeConfig string, env string) {
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
	for _, gConf := range routeConf.Groups {
		r.Groups[gConf.Name] = gConf
	}
	r.Actions = extract(routeConf.ActionConf, r.Groups)
}

func extract(actionConf map[string]actionConf, groups map[string]groupConf) []*action{
	l := len(actionConf)
	ch := make(chan *action, l)
	defer close(ch)
	
	var actions []*action
	
	for uri, conf := range actionConf {
		u, c := uri, conf
		go func() {
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
				ch <- &action{Uri:""}
				return
			}
			pathArray := strings.Split(c.Uses, "@")
			if len(pathArray) != 2 {
				ch <- &action{Uri:""}
				return
			}
			controller, deal = pathArray[0], pathArray[1]
			if !strings.EqualFold(c.Group, "") {
				groupNames := strings.Split(c.Group, "|")
				for _, groupName := range groupNames {
					if group, ok := groups[groupName]; ok {
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
			ch <- &action{
				Uri:u,
				Method:method,
				Controller:controller,
				Action:deal,
				Middleware:middleware,
			}
		}()
	}
	
	for i := 0; i < l; i++ {
		a := <-ch
		if strings.EqualFold(a.Uri, "") {
			continue
		}
		actions = append(actions, a)
	}
	return actions
}
