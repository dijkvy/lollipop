# kratos-plugin/log/config
.PHONY: log_config
#
log_config:
	 protoc --proto_path=/usr/local/include/ \
           --proto_path=. \
           --go_out=. \
           --go-grpc_out=. \
           ./kratos-plugin/log/config/log_config.proto


