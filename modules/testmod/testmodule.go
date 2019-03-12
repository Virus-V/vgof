package main

import (
	"log"
	"vgof/core"

	"go.uber.org/zap"
)

// Application 服务接口
type app struct {
	logger *zap.Logger
}

func (a *app) Main(s core.Service) {
	log.Print("Hello Module! adsadsad ")
	a.logger.Info("Hello world zlog!")
}
