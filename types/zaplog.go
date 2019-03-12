package types

import (
	"vgof/core"

	"github.com/google/uuid"
)

/**
zap Log 模块接口
*/

// SrvZapLogUUID Zap 日志模块uuid
var SrvZapLogUUID = uuid.UUID{0x3a, 0x80, 0x5e, 0x12, 0x5f, 0xe1, 0x4c, 0x92, 0x8b, 0xa, 0x86, 0x1a, 0xde, 0x55, 0x8e, 0x66}

// ZapLogService Zap日志模块接口
type ZapLogService interface {
	// 获得Zap对象
	GetZapLogger(core.Service) interface{}
}
