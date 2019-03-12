package main

import (
	"idea/core"
	"log"
)

type testModule struct {
}

// Application 服务接口
type app struct {
}

// Module 模块入口点
var Module = &testModule{}

var _ core.Module = (*testModule)(nil)

func (t *testModule) CheckDepend(s core.Service) bool {
	return true
}

// Start 实现module的初始化，如果返回true，则表示安装了新的service
func (t *testModule) Start(s core.Service) bool {
	appObj := &app{}
	err := s.InstallService(core.SrvApplicationUUID, appObj)
	if err != nil { // 如果出错
		panic(err)
	}
	return true
}

func (a *app) Main(s core.Service) {
	log.Print("Hello Module!")
}

func main() {
	log.Fatal("this is a module!!")
}
