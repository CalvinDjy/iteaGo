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

type Response struct {
	Data interface{}
	Header map[string]string
}

func (r *Response) SetHeader(key string, value string) {
	r.Header[key] = value
}

type HttpServer struct {
	Name			string
	Ip 				string
	Port 			string
	ReadTimeout 	int
	WriteTimeout 	int
	Route			string
	Ctx             context.Context
	Ioc 			*Ioc
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
	hs.Router.InitRoute(hs.Route, conf.Env)

	//Create route manager
	mux := http.NewServeMux()

	for _, a := range hs.Router.Actions {
		hs.wg.Add(1)
		go func(action *action) {
			defer hs.wg.Done()
			method := hs.extractMethod(action)

			//Get action interceptor list
			interceptor := ActionInterceptor(action.Middleware, hs.Ioc)
			
			mux.HandleFunc(action.Uri, hs.handler(action, method, interceptor))
		}(a)
	}
	hs.wg.Wait()

	hs.ser.Handler = mux
	//Start http server
	hs.start()
}

//Http handler
func (hs *HttpServer) handler(action *action, method reflect.Value, interceptor []IInterceptor) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){

		hs.wg.Add(1)

		response := &Response{
			Header: make(map[string]string),
		}
		rr, rw := reflect.ValueOf(r), reflect.ValueOf(w)

		defer hs.output(w, response)

		if !strings.EqualFold(r.Method, action.Method) {
			response.Data = "Method not allowed"
			return
		}

		n := method.Type().NumIn()
		if n > 2 {
			panic("Action params must be (*http.Request) or (*http.Request, http.ResponseWriter)")
		}

		p := []reflect.Value{rr}
		if n == 2 {
			p[1] = rw
		}

		f := func(*http.Request) (*Response, error){
			res := method.Call(p)
			switch len(res) {
			case 0:
				return response, nil
			case 1:
				if err, ok := res[0].Interface().(error); ok {
					return response, err
				}
				response.Data = res[0].Interface()
				return response, nil
			default:
				err := res[1].Interface()
				response.Data = res[0].Interface()
				if err != nil {
					return response, err.(error)
				}
				return response, nil
			}
		}

		for _, i := range interceptor {
			f = i.Handle(f)
		}

		_, err := f(r)
		if err != nil {
			response.Data = err.Error()
		}
	}
}

func (hs *HttpServer) extractMethod(a *action) reflect.Value{
	c := reflect.ValueOf(hs.Ioc.InsByName(a.Controller))
	if !c.IsValid() {
		panic(fmt.Sprint("Controller [", a.Controller, "] need register"))
	}
	m := c.MethodByName(a.Action)

	if !m.IsValid() {
		panic(fmt.Sprint("Can not find method [", a.Action, "] in [", a.Controller, "]"))
	}

	return m
}

//Http server start
func (hs *HttpServer) start() {
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
func (hs *HttpServer) stop() {
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