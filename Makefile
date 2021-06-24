# kratos-plugin/log/config
.PHONY: init
init:
	go get github.com/favadi/protoc-go-inject-tag

.PHONY: log_config
#
log_config:
	 protoc --proto_path=/usr/local/include/ \
           --proto_path=. \
           --go_out=. \
           --go-grpc_out=. \
           ./zap-log/init/config/log_config.proto

	protoc-go-inject-tag -input=./zap-log/init/config/log_config.pb.go
