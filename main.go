package main

import (
	"idea/core"
	"log"
)

func main() {
	loader := core.NewLoader(core.OptModulePath("./modules"))
	err := loader.LoadAll()
	if err != nil {
		log.Fatal(err)
	}
	// 开始任务
	err = loader.Start()
	if err != nil {
		log.Fatal(err)
	}
}
