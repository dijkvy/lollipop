package etcd

import (
	kratosEtcd "github.com/go-kratos/etcd/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func RegistrarDiscovery(cfg Config, logger *zap.Logger) *kratosEtcd.Registry {
	client, err := clientv3.New(
		clientv3.Config{
			Endpoints: cfg.Endpoints,
			Logger:    logger,
		},
	)
	if err != nil {
		panic(err)
	}
	return kratosEtcd.New(client)
}
