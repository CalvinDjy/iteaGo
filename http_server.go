package itea

import (
	"net/http"
	"time"
	"log"
	"reflect"
	"io"
	"github.com/json-iterator/go"
	"fmt"
	"context"
	"sync"
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
	ser 			*http.Server
	wg 				sync.WaitGroup
}

//Http server init
func (hs *HttpServer)Execute() {

	//Create http server
	hs.ser = &http.Server{
		ReadTimeout : DEFAULT_READ_TIMEOUT * time.Second,
		WriteTimeout : DEFAULT_WRITE_TIMEOUT * time.Second,
	}

	//Init route
	route := NewRoute().Init(hs.Route)

	//Get interceptor list
	interceptor := NewInterceptorManager(hs.Ioc).GetInterceptor()

	//Create route manager
	mux := http.NewServeMux()

	for u, a := range route.Actions {
		uri, action := u, a
		method := reflect.ValueOf(hs.Ioc.GetInstanceByName(action.Controller)).MethodByName(action.Action)

		if !method.IsValid() {
			log.Println("Can not find method [", action.Action, "] in [", action.Controller, "]")
		}

		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request){

			hs.wg.Add(1)
			//log.Println(r.Method)

			var result interface{}

			defer hs.output(w, &result)

			for _, ins := range interceptor {
				err := ins[0].Call([]reflect.Value{reflect.ValueOf(r)})[0].Interface()
				if err != nil {
					result = reflect.ValueOf(err)
					break
				}
				defer ins[1].Call([]reflect.Value{reflect.ValueOf(r), reflect.ValueOf(&result)})[0].
					Call([]reflect.Value{})
			}

			if reflect.ValueOf(result).Kind() == reflect.Invalid {
				res := method.Call([]reflect.Value{reflect.ValueOf(r)});
				if len(res) > 1 {
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

	log.Println("=== Http server [", hs.Name, "] start [", hs.ser.Addr, "] ===")
	err := hs.ser.ListenAndServe()
	if err != nil {
		log.Println("=== Http server [", hs.Name, "] stop [", err, "] ===")
	}
}

//Http server stop
func (hs *HttpServer)stop() {
	for {
		select {
		case <-	hs.Ctx.Done():
			log.Println("Http server stop ...")
			log.Println("Wait for all requests return ...")
			hs.wg.Wait()
			hs.ser.Shutdown(hs.Ctx)
			log.Println("Http server stop success")
			return
		default:
			break
		}
	}
}

//Http server output
func (hs *HttpServer) output(w http.ResponseWriter, result *interface{}) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if _, ok := (*result).(string); !ok {
		r, err := json.Marshal(*result)
		if err != nil {
			log.Println(err)
		}
		io.WriteString(w, string(r))
	} else {
		io.WriteString(w, (*result).(string))
	}
	hs.wg.Done()
}