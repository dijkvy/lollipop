package etcd

import (
	kratosEtcd "github.com/go-kratos/etcd/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func RegistrarDiscovery(cfg Config, logger *zap.Logger) *kratosEtcd.Registry {
	client, err := clientv3.New(
		clientv3.Config{
			Endpoints:   cfg.Endpoints,
			Logger:      logger,
			DialTimeout: time.Second,
			DialOptions: []grpc.DialOption{},
		},
	)
	if err != nil {
		panic(err)
	}
	return kratosEtcd.New(client)
}
