package router

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// TODO: depracted. using consul function instead.
type ConsulHandler interface {
	Service(consulAddr, serviceName string) (srv *api.CatalogService, err error)
}

type consulHandler struct {
}

func NewConsulHandler() ConsulHandler {
	return &consulHandler{}
}

func (c consulHandler) Service(consulAddr, serviceName string) (srv *api.CatalogService, err error) {
	consulClient, _ := api.NewClient(&api.Config{
		Address: consulAddr,
	})

	services, _, err := consulClient.Catalog().Service(serviceName, "", nil)
	if err != nil {
		return
	}

	if len(services) < 1 {
		err = fmt.Errorf("No consul service found for '%s'", serviceName)
		return
	}

	srv = services[0]
	return
}
