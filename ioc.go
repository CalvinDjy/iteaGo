package itea

import (
	"reflect"
	"sync"
	"strings"
	"encoding/json"
	"context"
	"io/ioutil"
)

const (
	PROCESS_CONFIG 	= "process"
	DB_CONFIG 		= "database"
	IMPORT_CONFIG 	= "@import"

	NAME_KEY 		= "Name"
	IOC_KEY 		= "Ioc"
	CTX_KEY 		= "Ctx"
	CONSTRUCT_FUNC 	= "Construct"
	INIT_FUNC 		= "Init"
	EXEC_FUNC 		= "Execute"
)

type Ioc struct {
	ctx 				context.Context
	register 			*Register
	typeN 				map[string]reflect.Type
	insN 				map[string]interface{}
	insT 				map[reflect.Type]interface{}
	idT					map[reflect.Type]int
	imports				map[string]map[string]string
	mutex 				*sync.Mutex
	wg 					sync.WaitGroup
	tid					int
}

//Create ioc
func NewIoc(ctx context.Context) (*Ioc) {
	register := NewRegister()
	ioc := &Ioc{
		ctx:ctx,
		register:register,
		typeN:make(map[string]reflect.Type),
		insN:make(map[string]interface{}),
		insT:make(map[reflect.Type]interface{}),
		idT:make(map[reflect.Type]int),
		mutex:new(sync.Mutex),
		imports:make(map[string]map[string]string),
		tid:0,
	}
	
	ioc.wg.Add(3)
	
	go func() {
		if !strings.EqualFold(ctx.Value(DB_CONFIG).(string), "") {
			ioc.registerDatabase(ctx.Value(DB_CONFIG).(string))
		}
		ioc.wg.Done()
	}()
	go func() {

		if len(ctx.Value(IMPORT_CONFIG).([]string)) != 0 {
			ioc.importConfig(ctx.Value(IMPORT_CONFIG).([]string))
		}
		ioc.wg.Done()
	}()

	go func() {
		tl := register.Init()
		if len(tl) > 0 {
			for _, t := range tl {
				ioc.typeN[t.Name()] = t
				ioc.tid++
				ioc.idT[t] = ioc.tid
			}
		}
		ioc.wg.Done()
	}()

	ioc.wg.Wait()

	return ioc
}

//Register database config
func (ioc *Ioc) registerDatabase(dbConfig string) {
	dat, err := ioutil.ReadFile(dbConfig)
	if err != nil {
		panic(err)
	}
	var databases map[string]*json.RawMessage
	err = json.Unmarshal(dat, &databases)
	if err != nil {
		panic(err)
	}
	for k, v := range databases {
		ioc.ctx = context.WithValue(ioc.ctx, k, v)
	}
}

//Import config
func (ioc *Ioc) importConfig(imports []string) {
	for _, filePath := range imports {
		dat, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic("Import [" + filePath + "] config not find")
		}
		var imconfig map[string]map[string]string
		err = json.Unmarshal(dat, &imconfig)
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
			ioc.tid++
			ioc.idT[t] = ioc.tid
		}
	}
}

//Init process of application
func (ioc *Ioc) InitProcess(ctx context.Context, bean Bean) {
	p := reflect.New(ioc.getType(bean.Class))

	ioc.wg.Add(1)
	go func() {
		for k, v := range bean.Params {
			setField(p, k, v)
		}
		ioc.wg.Done()
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
func (ioc *Ioc) GetInstanceByName(name string) (interface{}) {
	defer ioc.mutex.Unlock()
	ioc.mutex.Lock()
	return ioc.getInstanceByName(name)
}

//Get instance by Type
func (ioc *Ioc) GetInstanceByType(t reflect.Type) (interface{}) {
	defer ioc.mutex.Unlock()
	ioc.mutex.Lock()
	return ioc.getInstanceByType(t)
}

func (ioc *Ioc)getInstanceByName(name string) (interface{}) {
	var(
		instance interface{}
		exist bool
	)
	if instance, exist = ioc.insN[name];!exist {
		instance = ioc.buildInstance(ioc.getType(name))
	}
	return instance
}

func (ioc *Ioc)getInstanceByType(t reflect.Type) (interface{}) {
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

	setField(ins, CTX_KEY, ioc.ctx)

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
			if i := ioc.getInstanceByType(f.Type()); i != nil {
				f.Set(reflect.ValueOf(i).Elem())
			}
			break
		case reflect.Ptr:
			if i := ioc.getInstanceByType(f.Type().Elem()); i != nil {
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