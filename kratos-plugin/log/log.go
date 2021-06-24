package log

import "github.com/google/wire"

// wire data set
var ProviderSet = wire.NewSet(NewCoreLogger, NewKratosLog)
