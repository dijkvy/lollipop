package pprof

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/laxiaohong/lollipop/kratos-plugin/log"
	zap_log "github.com/laxiaohong/lollipop/zap-log"
	"github.com/laxiaohong/lollipop/zap-log/init/config"
	"testing"
	"time"
)

func TestRegisterPprof(t *testing.T) {
	logger := zap_log.NewZapLogger(&config.ZapLoggerConfig{Level: zap_log.LevelDebug, Path: "logs"})
	defer logger.Sync()
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(log.NewKratosLog(logger)),
			metrics.Server(),
			validate.Validator(),
		),
	}
	opts = append(opts, http.Network("tcp"))
	opts = append(opts, http.Address("0.0.0.0:1234"))
	opts = append(opts, http.Timeout(time.Second*8))

	log.NewCoreLogger(logger).WithContext(context.TODO()).Info("demo")
	server := http.NewServer(opts...)
	RegisterPprof(server)
	if err := kratos.New(
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			server,
		),
	).Run(); err != nil {

	}
}
