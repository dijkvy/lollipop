package etcd

import (
	zap_log "github.com/laxiaohong/lollipop/zap-log"
	"github.com/laxiaohong/lollipop/zap-log/init/config"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	defer logger.Sync()
}

var logger = zap_log.NewZapLogger(&config.ZapLoggerConfig{Path: "logs", Level: zap_log.LevelDebug})

func TestRegistrarDiscovery(t *testing.T) {
	rd := RegistrarDiscovery(Config{Endpoints: []string{"127.0.0.1:2379"}}, logger)
	_ = rd
}
