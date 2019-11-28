package ihttp

import (
	"context"
	"fmt"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/ioc/iface"
	"github.com/CalvinDjy/iteaGo/system"
	"github.com/json-iterator/go"
	"io"
	"net/http"
	"reflect"
	"sync"
	"strings"
	"time"
)

const (
	DEFAULT_READ_TIMEOUT 	= 1
	DEFAULT_WRITE_TIMEOUT 	= 30
)

type Response struct {
	Data 	interface{}
	Header 	map[string]string
}

func (r *Response) SetHeader(key string, value string) {
	r.Header[key] = value
}

type routeAction struct {
	exec reflect.Value
	method string
	interceptor []IInterceptor
}

type HttpServer struct {
	Ctx             context.Context
	Ioc 			iface.IIoc
	Name			string
	Ip 				string
	Port 			int
	ReadTimeout 	int
	WriteTimeout 	int
	Route			string
	Router			Route
	ser 			*http.Server
	wg 				sync.WaitGroup
}

//Http server init
func (hs *HttpServer) Execute() {

	//Create http server
	hs.ser = &http.Server{
		ReadTimeout : DEFAULT_READ_TIMEOUT * time.Second,
		WriteTimeout : DEFAULT_WRITE_TIMEOUT * time.Second,
	}

	//Init route
	hs.Router.InitRoute(hs.Route, system.Env)

	//Create route manager
	mux := http.NewServeMux()

	for p, as := range hs.Router.Actions {
		hs.wg.Add(1)
		go func(path string, actions []*action) {
			defer hs.wg.Done()
			
			var routeActions []routeAction
			for _, a := range actions {
				exec := hs.extractExec(a)

				if exec == reflect.ValueOf(nil) {
					continue
				}
				
				//Get action interceptor list
				interceptor := ActionInterceptor(a.Middleware, hs.Ioc)
				
				routeActions = append(routeActions, routeAction{
					exec: exec,
					method: a.Method,
					interceptor: interceptor,
				})
			}
			
			mux.HandleFunc(path, hs.handler(routeActions))
		}(p, as)
	}

	hs.wg.Wait()

	hs.ser.Handler = mux
	//Start http server
	hs.start()
}

//Http handler
func (hs *HttpServer) handler(routeActions []routeAction) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		r.ParseForm()
		
		hs.wg.Add(1)

		response := &Response{
			Header: make(map[string]string),
		}
		rr, rw := reflect.ValueOf(r), reflect.ValueOf(w)

		defer hs.output(w, response)
		
		var exec reflect.Value
		var interceptor []IInterceptor
		
		if len(routeActions) == 0 {
			response.Data = "Page not found"
			return
		}
		
		for _, ra := range routeActions {
			if strings.EqualFold(ra.method, r.Method) {
				exec, interceptor = ra.exec, ra.interceptor
				break
			}
		}
		
		if !exec.IsValid() {
			response.Data = "Method not allowed"
			return
		}

		n := exec.Type().NumIn()
		if n > 2 {
			panic("Action params must be (*http.Request) or (*http.Request, http.ResponseWriter)")
		}

		p := []reflect.Value{}
		if n == 1 {
			p = append(p, rr)
		} else if n == 2 {
			p = append(p, rr, rw)
		}

		f := func(request *http.Request, response *Response) error {
			res := exec.Call(p)
			switch len(res) {
			case 0:
				return nil
			case 1:
				if err, ok := res[0].Interface().(error); ok {
					return err
				}
				response.Data = res[0].Interface()
				return nil
			case 2:
				response.Data = res[0].Interface()
				err := res[1].Interface()
				if err != nil {
					if e, ok := err.(error); ok {
						return e
					} else {
						ilog.Error("invalid type of the second return param, error accepted")
						return nil
					}
				}
				return nil
			default:
				ilog.Error("invalid num of return, 2 or less is accepted")
				return nil
			}
		}

		for _, i := range interceptor {
			f = i.Handle(f)
		}

		err := f(r, response)
		if err != nil {
			response.Data = err.Error()
		}
	}
}

func (hs *HttpServer) extractExec(a *action) reflect.Value {
	c := reflect.ValueOf(hs.Ioc.InsByName(a.Controller))
	if !c.IsValid() {
		ilog.Error(fmt.Sprintf("controller [%s] is nil, please check out if [%s] is registed", a.Controller, a.Controller))
		//panic(fmt.Sprintf("Controller [%s] need register", a.Controller))
		return reflect.ValueOf(nil)
	}
	m := c.MethodByName(a.Action)

	if !m.IsValid() {
		ilog.Error(fmt.Sprintf("can not find method [%s] in [%s]", a.Action, a.Controller))
		return reflect.ValueOf(nil)
		//panic(fmt.Sprintf("Can not find method [%s] in [%s]", a.Action, a.Controller))
	}

	return m
}

//Http server start
func (hs *HttpServer) start() {
	hs.ser.Addr = fmt.Sprintf("%s:%d", hs.Ip, hs.Port)
	if hs.ReadTimeout != 0 {
		hs.ser.ReadTimeout = time.Duration(hs.ReadTimeout) * time.Second
	}
	if hs.WriteTimeout != 0 {
		hs.ser.WriteTimeout = time.Duration(hs.WriteTimeout) * time.Second
	}

	go hs.stop()

	ilog.Info(fmt.Sprintf("=== 【Http】Server [%s] start [%s] ===", hs.Name,  hs.ser.Addr))
	if err := hs.ser.ListenAndServe(); err != nil {
		ilog.Info(fmt.Sprintf("http server [%s] stop [%s]", hs.Name, err))
	}
}

//Http server stop
func (hs *HttpServer) stop() {
	for {
		select {
		case <-	hs.Ctx.Done():
			ilog.Info("http server stop ...")
			ilog.Info("wait for all http requests return ...")
			hs.wg.Wait()
			hs.ser.Shutdown(hs.Ctx)
			ilog.Info("http server stop success")
			return
		}
	}
}

//Http server output
func (hs *HttpServer) output(w http.ResponseWriter, response *Response) {
	defer hs.wg.Done()
	var json = jsoniter.Config{
		EscapeHTML:             false,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
	if response.Header != nil {
		for k, v := range response.Header {
			w.Header().Set(k, v)
		}
	}
	if _, ok := (*response).Data.(string); !ok {
		r, err := json.Marshal((*response).Data)
		if err != nil {
			ilog.Error(err)
		}
		io.WriteString(w, string(r))
	} else {
		io.WriteString(w, (*response).Data.(string))
	}
}