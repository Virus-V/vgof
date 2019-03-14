package main

import (
	"log"
	"vgof/core"
)

func main() {
	loader := core.NewLoader(
		core.OptModulePath("./modules"),
		core.OptGlobalServices())
	modules, err := loader.LoadAll()
	if err != nil {
		log.Fatal(err)
	}
	// 调度模块列表
	err = loader.Dispatch(modules)
	if err != nil {
		log.Fatal(err)
	}
	// 开始任务
	err = loader.Start()
	if err != nil {
		log.Fatalf("Application exit with error: %s\n", err)
	}
}
