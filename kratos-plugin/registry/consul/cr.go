package consul

import (
	"github.com/go-kratos/consul/registry"
	"github.com/hashicorp/consul/api"
	"time"
)

var _address = "127.0.0.1:8500"

func RegistrarDiscovery(address string) *registry.Registry {
	cli, err := api.NewClient(
		&api.Config{
			Address:  address,
			WaitTime: time.Second * 6,
		})
	if err != nil {
		panic(err)
	}
	return registry.New(cli)
}
