package itea

import (
	"net/http"
	"time"
	"github.com/CalvinDjy/iteaGo/ilog"
	"reflect"
	"io"
	"github.com/json-iterator/go"
	"fmt"
	"context"
	"sync"
	"strings"
)

const (
	DEFAULT_READ_TIMEOUT 	= 1
	DEFAULT_WRITE_TIMEOUT 	= 30
)

type HttpServer struct {
	Name			string
	Ip 				string
	Port 			string
	ReadTimeout 	float64
	WriteTimeout 	float64
	Route			string
	Ctx             context.Context
	Ioc 			*Ioc
	router			Route
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
	route := hs.router.Init(hs.Route, conf.Env)

	//Get interceptor list
	interceptor := GetInterceptor(hs.Ioc)

	//Create route manager
	mux := http.NewServeMux()

	for u, a := range route.Actions {
		uri, action := u, a
		controller := reflect.ValueOf(hs.Ioc.GetInstanceByName(action.Controller))
		if !controller.IsValid() {
			panic(fmt.Sprint("Controller [", action.Controller, "] need register"))
		}
		method := controller.MethodByName(action.Action)
		
		if !method.IsValid() {
			ilog.Error("Can not find method [", action.Action, "] in [", action.Controller, "]")
		}

		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request){

			hs.wg.Add(1)

			var result interface{}
			rr, rw := reflect.ValueOf(r), reflect.ValueOf(w)
			
			defer hs.output(w, &result)

			for _, ins := range interceptor {
				err := ins[0].Call([]reflect.Value{rr})[0].Interface()
				if err != nil {
					result = reflect.ValueOf(err).Interface()
					break
				}
				defer ins[1].Call([]reflect.Value{rr, reflect.ValueOf(&result)})
			}
			
			if !strings.EqualFold(r.Method, action.Method) {
				result = NewServerError("Method not allowed")
			}
			
			if reflect.ValueOf(result).Kind() == reflect.Invalid {
				n := method.Type().NumIn()
				if n > 2 {
					panic("Action params must be empty or (*http.Request) or (*http.Request, http.ResponseWriter)")
				}

				var p, res []reflect.Value
				if n == 2 {
					p = []reflect.Value{rr, rw};
				} else if n == 1 {
					p = []reflect.Value{rr};
				}
				res = method.Call(p);

				rl := len(res)

				if rl == 0 {
					result = nil
					return
				}

				if rl > 1 {
					if _, ok := res[1].Interface().(error); ok {
						result = res[1].Interface()
						return
					}
				}
				result = res[0].Interface()
			}
		})
	}

	hs.ser.Handler = mux
	//Start http server
	hs.start()
}

//Http server start
func (hs *HttpServer)start() {
	hs.ser.Addr = fmt.Sprintf("%s:%s", hs.Ip, hs.Port)
	if hs.ReadTimeout != 0 {
		hs.ser.ReadTimeout = time.Duration(hs.ReadTimeout) * time.Second
	}
	if hs.WriteTimeout != 0 {
		hs.ser.WriteTimeout = time.Duration(hs.WriteTimeout) * time.Second
	}

	go hs.stop()

	ilog.Info("=== Http server [", hs.Name, "] start [", hs.ser.Addr, "] ===")
	err := hs.ser.ListenAndServe()
	if err != nil {
		ilog.Info("=== Http server [", hs.Name, "] stop [", err, "] ===")
	}
}

//Http server stop
func (hs *HttpServer)stop() {
	for {
		select {
		case <-	hs.Ctx.Done():
			ilog.Info("Http server stop ...")
			ilog.Info("Wait for all http requests return ...")
			hs.wg.Wait()
			hs.ser.Shutdown(hs.Ctx)
			ilog.Info("Http server stop success")
			return
		default:
			break
		}
	}
}

//Http server output
func (hs *HttpServer) output(w http.ResponseWriter, result *interface{}) {
	defer hs.wg.Done()
	var json = jsoniter.Config{
		EscapeHTML:             false,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
	if _, ok := (*result).(string); !ok {
		r, err := json.Marshal(*result)
		if err != nil {
			ilog.Error(err)
		}
		io.WriteString(w, string(r))
	} else {
		io.WriteString(w, (*result).(string))
	}
}