package main

import (
	"log"
	"vgof/core"
	"vgof/types"

	"go.uber.org/zap"
)

type testModule struct {
}

// ModuleEntry 模块入口点
var ModuleEntry = testModule{}

var _ core.Module = (*testModule)(nil)

func (t *testModule) CheckDepend(s core.Service) bool {
	// 该模块需要zaplog模块
	return s.CheckServices(types.SrvZapLogUUID)
}

// Start 实现module的初始化，如果返回true，则表示安装了新的service
func (t *testModule) Start(s core.Service) bool {
	appObj := &app{}
	zaplogSrv, err := s.LocateService(types.SrvZapLogUUID)
	if err != nil {
		panic(err)
	}
	// 安装日志插件
	appObj.logger = ((zaplogSrv.(types.ZapLogService)).GetZapLogger(s)).(*zap.Logger)
	// 安装Application服务
	err = s.InstallService(core.SrvApplicationUUID, appObj)
	if err != nil { // 如果出错
		panic(err)
	}
	return true
}

// Stop 关闭插件
func (t *testModule) Stop(s core.Service) {
	log.Print("Test module stoped!")
}

func main() {
	log.Fatal("This is a vgof module, please build this package with \"-buildmode=plugin\".")
}
