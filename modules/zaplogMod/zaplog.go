package main

import (
	"vgof/core"
	"vgof/types"

	"go.uber.org/zap"
)

type zaplog struct {
	init   bool // 对象是否初始化
	logger *zap.Logger
}

var _ types.ZapLogService = (*zaplog)(nil)

func (z *zaplog) GetZapLogger(s core.Service) interface{} {
	var err error
	if !z.init {
		z.logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}
	return z.logger
}
