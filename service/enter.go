package service

import "sephiroth-go/service/sys"

type ServiceGroup struct {
	SystemServiceGroup sys.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
