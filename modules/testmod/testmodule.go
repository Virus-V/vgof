package main

import (
	"errors"
	"log"
	"vgof/core"

	"go.uber.org/zap"
)

// Application 服务接口
type app struct {
	logger *zap.Logger
}

func (a *app) Main(s core.Service) {
	sTmp := core.GlobalServices()
	if sTmp != s {
		log.Print("Not match")
	}
	log.Print("Hello Module! adsadsad ")
	a.logger.Info("Hello world zlog!")
	panic(errors.New("test error"))
}
