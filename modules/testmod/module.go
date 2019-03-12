package main

import (
	"log"
	"vgof/core"
)

type testModule struct {
}

// ModuleEntry 模块入口点
var ModuleEntry = testModule{}

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

func main() {
	log.Fatal("This is a vgof module, please build this package with \"-buildmode=plugin\".")
}
