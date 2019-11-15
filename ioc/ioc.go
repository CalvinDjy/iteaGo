package ioc

import (
	"context"
	"github.com/CalvinDjy/iteaGo/ioc/bean"
	"github.com/CalvinDjy/iteaGo/process"
	"github.com/CalvinDjy/iteaGo/system"
	"reflect"
	"strings"
	"sync"
	"fmt"
)

const (
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
	beansN				map[string]*bean.Bean
	beansT				map[reflect.Type]*bean.Bean
	insN 				map[string]interface{}
	insT 				map[reflect.Type]interface{}
	mutex 				*sync.Mutex
}

//Create ioc
func NewIoc(ctx context.Context) *Ioc {
	register := NewRegister()
	ioc := &Ioc{
		ctx:ctx,
		register:register,
		beansN:make(map[string]*bean.Bean),
		beansT:make(map[reflect.Type]*bean.Bean),
		insN:make(map[string]interface{}),
		insT:make(map[reflect.Type]interface{}),
		mutex:new(sync.Mutex),
	}
	
	ioc.appendBeans(register.Init())
	return ioc
}

//Register beans
func (ioc *Ioc) Register(beans [] interface{}) {
	ioc.appendBeans(ioc.register.Register(beans))
}

//Register beans
func (ioc *Ioc) RegisterBeans(beans []*bean.Bean) {
	ioc.appendBeans(ioc.register.RegisterBeans(beans))
}

func (ioc *Ioc) appendBeans(beans []*bean.Bean) {
	if len(beans) > 0 {
		for _, bean := range beans {
			ioc.beansN[bean.Name] = bean
			ioc.beansT[bean.GetAbstractType()] = bean
		}
	}
}

//Exec process of application
func (ioc *Ioc) ExecProcess(ctx context.Context, process *process.Process) {
	if strings.EqualFold(process.Class, "") {
		return
	}

	t := ioc.getType(process.Class)
	if t == nil {
		return
	}

	p := reflect.New(t)

	var wg sync.WaitGroup
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		for k, v := range process.Params {
			setField(p, k, v)
		}
	}()

	setField(p, NAME_KEY, process.Name)
	setField(p, CTX_KEY, ctx)
	setField(p, IOC_KEY, ioc)

	wg.Wait()

	// Do execute
	var exec string
	if !strings.EqualFold(process.ExecuteMethod, "") {
		exec = process.ExecuteMethod
	} else if p.MethodByName(EXEC_FUNC) != reflect.ValueOf(nil) {
		exec = EXEC_FUNC
	}

	if !strings.EqualFold(exec, "") {
		p.MethodByName(exec).Call([]reflect.Value{})
	}
}

func (ioc *Ioc) BeansByName(name string) *bean.Bean {
	if b, ok := ioc.beansN[name]; ok {
		return b
	}
	return nil
}

//Get instance by name
func (ioc *Ioc) InsByClass(i interface{}) interface{} {
	ioc.mutex.Lock()
	defer ioc.mutex.Unlock()
	return ioc.instanceByType(reflect.TypeOf(i))
}

//Get instance by name
func (ioc *Ioc) InsByName(name string) interface{} {
	ioc.mutex.Lock()
	defer ioc.mutex.Unlock()
	return ioc.instanceByName(name)
}

//Get instance by Type
func (ioc *Ioc) InsByType(t reflect.Type) interface{} {
	ioc.mutex.Lock()
	defer ioc.mutex.Unlock()
	return ioc.instanceByType(t)
}

func (ioc *Ioc) instanceByName(name string) interface{} {
	var(
		instance interface{}
		exist bool
	)
	if instance, exist = ioc.insN[name];!exist {
		instance = ioc.buildInstance(ioc.getType(name))
	}
	return instance
}

func (ioc *Ioc) instanceByType(t reflect.Type) interface{} {
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
func (ioc *Ioc) buildInstance(t reflect.Type) interface{} {
	if t == nil {
		return nil
	}
	
	scope := SINGLETON
	
	if bean, ok := ioc.beansT[t]; ok {
		t = bean.GetConcreteType()
		scope = bean.Scope
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
		f := ins.Elem().Field(index)
		if !f.CanSet() {
			continue
		}
		switch f.Kind() {
		case reflect.Struct:
			if tag := t.Field(index).Tag.Get("wired"); !strings.EqualFold(tag, "") {
				if i := ioc.instanceByType(f.Type()); i != nil {
					f.Set(reflect.ValueOf(i).Elem())
				}
			}
			break
		case reflect.Ptr:
			if tag := t.Field(index).Tag.Get("wired"); !strings.EqualFold(tag, "") {
				if i := ioc.instanceByType(f.Type().Elem()); i != nil {
					f.Set(reflect.ValueOf(i))
				}
			}
			break
		case reflect.String:
			if tag := t.Field(index).Tag.Get("value"); !strings.EqualFold(tag, "") {
				if v := system.Conf.GetString(tag); !strings.EqualFold(v, "") {
					f.Set(reflect.ValueOf(v))
				}
			}
			break
		case reflect.Int:
			if tag := t.Field(index).Tag.Get("value"); !strings.EqualFold(tag, "") {
				f.Set(reflect.ValueOf(system.Conf.GetInt(tag)))
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

	if ins.Interface() != nil && strings.EqualFold(scope, SINGLETON) {
		ioc.insN[t.Name()] = ins.Interface()
		ioc.insT[t] = ins.Interface()
	}

	return ins.Interface()
}

//Get type of bean
func (ioc *Ioc) getType(name string) reflect.Type{
	if t, ok := ioc.beansN[name]; ok {
		return t.GetConcreteType()
	}
	return nil
}

//Set field of instance
func setField(i reflect.Value, n string, v interface{}) {
	field := i.Elem().FieldByName(n)
	if field.CanSet() && v != nil {
		ft, vt := field.Type(), reflect.TypeOf(v)
		if ft == vt || (field.Kind().String() == "interface" && vt.Implements(ft)) {
			field.Set(reflect.ValueOf(v))
		} else {
			panic(fmt.Sprintf("can not inject %s(%s) with %s", n, ft.String(), vt.String()))
		}
	}
}