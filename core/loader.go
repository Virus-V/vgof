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
	String() string           // 获得模块名称
	CheckDepend(Service) bool // 检查Service依赖,如果为true则表示依赖满足,可以执行模块初始化
	Start(Service) bool       // 启动模块,true代表安装了新的service
	Stop(Service)             // 关闭插件
}

// Option 配置回调函数
type Option func(*object)

type object struct {
	// 服务列表
	serviceTable sync.Map
	modulePath   string
	moduleList   []Module
}

// Loader 加载器接口
type Loader interface {
	LoadAll() ([]Module, error)          // 加载模块目录下全部的模块
	LoadList([]string) ([]Module, error) // 加载指定的模块
	Dispatch(modules []Module) error     // 调度列表里的所有模块
	Start() error                        // 执行最终应用入口
}

var _ Loader = (*object)(nil)
var _ Service = (*object)(nil)

// NewLoader 创建新的加载器
func NewLoader(opts ...Option) Loader {
	obj := &object{
		moduleList: make([]Module, 0),
	}
	// 执行配置
	for _, opt := range opts {
		opt(obj)
	}
	return obj
}

// Load 加载目录下所有模块文件
func (o *object) LoadAll() ([]Module, error) {
	modules := make([]string, 0)
	// 获取所有文件
	modulePath, err := ioutil.ReadDir(o.modulePath)
	if err != nil {
		return nil, err
	}
	for _, file := range modulePath {
		if file.IsDir() {
			continue
		} else {
			// 过滤文件后缀
			ext := path.Ext(file.Name())
			if ext != ".so" {
				continue
			}
			modules = append(modules, strings.TrimSuffix(file.Name(), ext))
		}
	}
	if len(modules) == 0 {
		return nil, ErrNoMoudles
	}
	// 加载module
	return o.loadModules(modules)
}

// LoadList 加载指定的模块文件
func (o *object) LoadList(modules []string) ([]Module, error) {
	// 加载module
	return o.loadModules(modules)
}

func (o *object) loadModules(modules []string) ([]Module, error) {
	if modules == nil || len(modules) == 0 {
		return nil, ErrBadParame
	}
	// 已初始化的模块列表
	moduleList := make([]Module, 0)
	for _, m := range modules {
		where := o.modulePath + m + ".so"
		mObj, err := plugin.Open(where)
		if err != nil {
			return nil, err
		}
		var sym plugin.Symbol
		// 找到模块对象接口
		if sym, err = mObj.Lookup("ModuleEntry"); err != nil {
			return nil, err
		} else {
			// 放入模块列表
			moduleList = append(moduleList, sym.(Module))
		}
	}
	// 调度模块列表
	return moduleList, nil
}

// Dispatch 调度模块列表
func (o *object) Dispatch(modules []Module) (err error) {
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
	// 已初始化的模块列表
	moduleInited := make(map[Module]bool)
REINIT:
	reinit := false
	for _, m := range modules {
		if moduleInited[m] == true { // 跳过已初始化的模块
			continue
		}
		// 检查模块的依赖是否满足
		if m.CheckDepend(o) {
			log.Printf("Start module: \"%s\".\n", m)
			reinit = m.Start(o)                    // 执行模块初始化,并且更新 是否需要重新初始化 标志
			moduleInited[m] = true                 // 记录已经初始化
			o.moduleList = append(o.moduleList, m) // 记录插件列表
		} else {
			log.Printf("Module \"%s\" dependence fail.\n", m)
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
		log.Printf("Module \"%s\" has not been install.\n", m)
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
		// 停止模块
		for i := len(o.moduleList) - 1; i >= 0; i-- {
			o.moduleList[i].Stop(o)
		}
	}()
	// 找到SrvApplicationUUID这个服务
	var srv Service = o // 获得服务接口
	var app interface{} // 应用接口
	app, err = srv.LocateService(SrvApplicationUUID)
	if err != nil {
		panic(err)
	}
	(app.(ApplicationService)).Main(srv)
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

// OptGlobalServices 将该对象设置为全局服务表
func OptGlobalServices() Option {
	return func(o *object) {
		globalServices = o
	}
}
