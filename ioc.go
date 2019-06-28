package itea

import (
	"gopkg.in/yaml.v2"
	"context"
	"reflect"
	"sync"
	"strings"
	"io/ioutil"
)

const (
	PROCESS_CONFIG 	= "process"
	DB_CONFIG 		= "database"
	IMPORT_CONFIG 	= "import"
	LOG_CONFIG 		= "log"
	DEBUG			= "debug"
	CONNECTION_CONFIG		= "connections"
	REDIS_CONFIG			= "redis"

	NAME_KEY 		= "Name"
	IOC_KEY 		= "Ioc"
	CTX_KEY 		= "Ctx"
	CONSTRUCT_FUNC 	= "Construct"
	INIT_FUNC 		= "Init"
	EXEC_FUNC 		= "Execute"
)

type Ioc struct {
	register 			*Register
	typeN 				map[string]reflect.Type
	insN 				map[string]interface{}
	insT 				map[reflect.Type]interface{}
	imports				map[string]map[string]string
	mutex 				*sync.Mutex
	wg 					sync.WaitGroup
}

//Create ioc
func NewIoc() (*Ioc) {
	register := NewRegister()
	ioc := &Ioc{
		register:register,
		typeN:make(map[string]reflect.Type),
		insN:make(map[string]interface{}),
		insT:make(map[reflect.Type]interface{}),
		mutex:new(sync.Mutex),
		imports:make(map[string]map[string]string),
	}
	
	ioc.wg.Add(2)
	
	go func() {
		defer ioc.wg.Done()
		imp := conf.Config(IMPORT_CONFIG)
		if imp != nil && len(imp.([]string)) != 0 {
			ioc.importConfig(imp.([]string))
		}
	}()

	go func() {
		defer ioc.wg.Done()
		tl := register.Init()
		if len(tl) > 0 {
			for _, t := range tl {
				ioc.typeN[t.Name()] = t
			}
		}
	}()

	ioc.wg.Wait()

	return ioc
}

//Import config
func (ioc *Ioc) importConfig(imports []string) {
	for _, filePath := range imports {
		dat, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic("Import [" + filePath + "] config not find")
		}
		var imconfig map[string]map[string]string
		err = yaml.Unmarshal(dat, &imconfig)
		if err != nil {
			panic("Import [" + filePath + "] config extract error")
		}
		for k, v := range imconfig {
			ioc.imports[k] = v
		}
	}
}

//Register beans
func (ioc *Ioc) RegisterBeans(beans [] interface{}) {
	tl := ioc.register.Register(beans)
	if len(tl) > 0 {
		for _, t := range tl {
			ioc.typeN[t.Name()] = t
		}
	}
}

//Init process of application
func (ioc *Ioc) InitProcess(ctx context.Context, bean Bean) {
	if strings.EqualFold(bean.Class, "") {
		return
	}

	t := ioc.getType(bean.Class)
	if t == nil {
		return
	}

	p := reflect.New(t)

	ioc.wg.Add(1)
	go func() {
		defer ioc.wg.Done()
		for k, v := range bean.Params {
			setField(p, k, v)
		}
	}()

	setField(p, NAME_KEY, bean.Name)
	setField(p, CTX_KEY, ctx)
	setField(p, IOC_KEY, ioc)

	ioc.wg.Wait()

	// Do execute
	var exec string
	if !strings.EqualFold(bean.ExecuteMethod, "") {
		exec = bean.ExecuteMethod
	} else if p.MethodByName(EXEC_FUNC) != reflect.ValueOf(nil) {
		exec = EXEC_FUNC
	}

	if !strings.EqualFold(exec, "") {
		p.MethodByName(exec).Call([]reflect.Value{})
	}
}

//Get instance by name
func (ioc *Ioc) InsByName(name string) (interface{}) {
	defer ioc.mutex.Unlock()
	ioc.mutex.Lock()
	return ioc.instanceByName(name)
}

//Get instance by Type
func (ioc *Ioc) InsByType(t reflect.Type) (interface{}) {
	defer ioc.mutex.Unlock()
	ioc.mutex.Lock()
	return ioc.instanceByType(t)
}

func (ioc *Ioc) instanceByName(name string) (interface{}) {
	var(
		instance interface{}
		exist bool
	)
	if instance, exist = ioc.insN[name];!exist {
		instance = ioc.buildInstance(ioc.getType(name))
	}
	return instance
}

func (ioc *Ioc) instanceByType(t reflect.Type) (interface{}) {
	var(
		instance interface{}
		exist bool
	)
	if instance, exist = ioc.insT[t];!exist {
		instance = ioc.buildInstance(t)
	}
	return instance
}

//Create new instance
func (ioc *Ioc) buildInstance(t reflect.Type) (interface{}) {
	if t == nil {
		return nil
	}
	ins := reflect.New(t)

	setField(ins, CTX_KEY, ctx)

	//Execute construct method of instance
	cm := ins.MethodByName(CONSTRUCT_FUNC)
	if cm.IsValid() {
		cm.Call(nil)
	}

	//Inject construct params
	for index := 0; index < t.NumField(); index++ {
		f := ins.Elem().FieldByIndex([]int{index})
		if !f.CanSet() {
			continue
		}
		switch f.Kind() {
		case reflect.Struct:
			if i := ioc.instanceByType(f.Type()); i != nil {
				f.Set(reflect.ValueOf(i).Elem())
			}
			break
		case reflect.Ptr:
			if i := ioc.instanceByType(f.Type().Elem()); i != nil {
				f.Set(reflect.ValueOf(i))
			}
			break
		case reflect.String:
			if c, ok := ioc.imports[t.Name()]; ok {
				if v, ok := c[t.Field(index).Name]; ok {
					f.Set(reflect.ValueOf(v))
				}
			}
			break
		default:
			break
		}
	}

	//Execute init method of instance
	im := ins.MethodByName(INIT_FUNC)
	if im.IsValid() {
		if im.Type().NumIn() > 0 {
			im.Call([]reflect.Value{ins})
		} else {
			im.Call(nil)
		}
	}

	if ins.Interface() != nil {
		ioc.insN[t.Name()] = ins.Interface()
		ioc.insT[t] = ins.Interface()
	}

	return ins.Interface()
}

//Get type of bean
func (ioc *Ioc) getType(name string) reflect.Type{
	if t, ok := ioc.typeN[name]; ok {
		return t
	}
	return nil
}

//Set field of instance
func setField(i reflect.Value, n string, v interface{}) {
	field := i.Elem().FieldByName(n)
	if field.CanSet() {
		field.Set(reflect.ValueOf(v))
	}
}