package main

import (
	"log"
	"vgof/core"
)

// Application 服务接口
type app struct {
}

func (a *app) Main(s core.Service) {
	log.Print("Hello Module! adsadsad ")
}
