package core

import (
	"errors"
	"io/ioutil"
	"log"
	"path"
	"plugin"
	"runtime/debug"
	"strings"
	"sync"
)

// 错误值
var (
	ErrUnknowError = errors.New("Unknow error")
	ErrBadParame   = errors.New("Bad parames")
	ErrNoMoudles   = errors.New("No modules will load")
)

// Module 模块接口
type Module interface {
	CheckDepend(Service) bool // 检查Service依赖,如果为true则表示依赖满足,可以执行模块初始化
	Start(Service) bool       // 启动模块,true代表安装了新的service
}

// Option 配置回调函数
type Option func(*object)

type object struct {
	// 服务列表
	serviceTable sync.Map
	modulePath   string
}

// Loader 加载器接口
type Loader interface {
	LoadAll() error           // 加载模块目录下全部的模块
	LoadList(...string) error // 加载指定的模块
	Start() error             // 执行最终应用入口
}

var _ Loader = (*object)(nil)
var _ Service = (*object)(nil)

// NewLoader 创建新的加载器
func NewLoader(opts ...Option) Loader {
	obj := &object{
		//
	}
	// 执行配置
	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

// Load 加载目录下所有模块
func (o *object) LoadAll() error {
	modules := make([]string, 0)
	// 获取所有文件
	modulePath, _ := ioutil.ReadDir(o.modulePath)
	for _, file := range modulePath {
		if file.IsDir() {
			continue
		} else {
			// TODO 过滤文件后缀
			ext := path.Ext(file.Name())
			if ext != ".so" {
				continue
			}
			modules = append(modules, strings.TrimSuffix(file.Name(), ext))
		}
	}
	if len(modules) == 0 {
		return ErrNoMoudles
	}
	// 加载module
	return o.loadModules(modules...)
}

// LoadList 加载指定的模块
func (o *object) LoadList(modules ...string) error {
	// 加载module
	return o.loadModules(modules...)
}

func (o *object) loadModules(modules ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Print(r)
			switch r.(type) {
			case error:
				err = r.(error)
			default:
				err = ErrUnknowError
			}
		}
	}()
	if modules == nil || len(modules) == 0 {
		return ErrBadParame
	}
	// 已初始化的模块列表
	moduleInited := make(map[string]bool)
REINIT:
	reinit := false
	for _, m := range modules {
		if moduleInited[m] == true { // 跳过已初始化的模块
			continue
		}
		where := o.modulePath + m + ".so"
		mObj, err := plugin.Open(where)
		if err != nil {
			panic(err)
		}
		var sym plugin.Symbol
		// 找到模块对象接口
		if sym, err = mObj.Lookup("ModuleEntry"); err != nil {
			panic(err)
		} else {
			module := sym.(Module)
			if module.CheckDepend(o) { // 检查模块的依赖是否满足
				log.Printf("Start module %s.\n", m)
				reinit = module.Start(o) // 执行模块初始化,并且更新 是否需要重新初始化 标志
				moduleInited[m] = true   // 记录已经初始化
			} else {
				log.Printf("Module %s dependence fail.\n", m)
			}
		}
	}
	if reinit {
		goto REINIT
	}
	// 打印出未初始化的模块
	for _, m := range modules {
		if moduleInited[m] == true { // 跳过已初始化的模块
			continue
		}
		log.Printf("Module %s has not been install.\n", m)
	}
	return nil
}

// Start 开始应用
func (o *object) Start() (err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Print(r)
			switch r.(type) {
			case error:
				err = r.(error)
			default:
				err = ErrUnknowError
			}
		}
	}()
	// 找到SrvApplicationUUID这个服务
	var srv Service = o // 获得服务接口
	var app interface{} // 应用接口
	app, err = srv.LocateService(SrvApplicationUUID)
	(app.(Application)).Main(srv)
	log.Print("Application returned.")
	return nil
}

// OptModulePath 设置模块目录
func OptModulePath(path string) Option {
	return func(o *object) {
		path = strings.TrimRight(path, "/") + "/"
		o.modulePath = path
	}
}
