package service

import (
	"github.com/icosmos-space/iadmin/server/service/example"
	"github.com/icosmos-space/iadmin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
}
