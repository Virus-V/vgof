package core

import (
	"errors"

	"github.com/google/uuid"
)

// 错误码
var (
	ErrNotInstall = errors.New("The service has not installed")
	ErrInstalled  = errors.New("The service has already installed")
)

// Service 服务接口
type Service interface {
	LocateService(name uuid.UUID) (interface{}, error)
	InstallService(name uuid.UUID, service interface{}) error
	ReplaceService(name uuid.UUID, service interface{}) interface{}
	UninstallServices(name ...uuid.UUID)
	CheckServices(name ...uuid.UUID) bool
}

// LocateService 获得服务
func (o *object) LocateService(name uuid.UUID) (interface{}, error) {
	if service, exists := o.serviceTable.Load(name); exists {
		return service, nil
	}
	return nil, ErrNotInstall
}

// InstallService 安装服务
func (o *object) InstallService(name uuid.UUID, service interface{}) error {
	if _, loaded := o.serviceTable.LoadOrStore(name, service); loaded {
		return ErrInstalled
	}
	return nil
}

// ReplaceService 替换服务
func (o *object) ReplaceService(name uuid.UUID, service interface{}) interface{} {
	if orign, loaded := o.serviceTable.LoadOrStore(name, service); loaded {
		o.serviceTable.Store(name, service)
		return orign
	}
	return nil
}

// UninstallServices 卸载服务
func (o *object) UninstallServices(name ...uuid.UUID) {
	if name == nil || len(name) == 0 {
		return
	}
	for _, v := range name {
		o.serviceTable.Delete(v)
	}
	return
}

// CheckServices 检查服务是否安装
func (o *object) CheckServices(name ...uuid.UUID) bool {
	if name == nil || len(name) == 0 {
		return true
	}
	cnt := len(name)
	var existed int
	o.serviceTable.Range(func(k, v interface{}) bool {
		for _, v := range name {
			if v == k.(uuid.UUID) {
				existed++
			}
		}
		return true
	})
	if cnt != existed {
		return false
	}
	return true
}
