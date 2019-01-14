package itea

import (
	"io/ioutil"
	"encoding/json"
	"strings"
	"os"
)

type Route struct {
	Actions 	map[string]*action
}

type action struct {
	Method 		string
	Controller 	string
	Action 		string
}

func NewRoute() (r *Route) {
	return &Route{}
}

func (r *Route)Init(routeConfig string) (route *Route){
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(projectPath + routeConfig)
	if err != nil {
		panic("Route config not find")
	}
	var mapping map[string][]string
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		panic("Route config extract fail")
	}
	r.Actions = make(map[string]*action)
	r.extract(mapping)
	return r
}

func (r *Route)extract(mapping map[string][]string) {
	for uri, path := range mapping {
		pathArray := strings.Split(path[1], "/")
		r.Actions[uri] = &action{
			Method:path[0],
			Controller:pathArray[0],
			Action:pathArray[1],
		}
	}
}
