proto-auth:
	@protoc \
	--proto_path=protobuf "protobuf/auth.proto" \
	--go_out=services/common/genproto/auth \
	--go_opt=paths=source_relative \
	--go-grpc_out=services/common/genproto/auth \
	--go-grpc_opt=paths=source_relative \

proto-user:
	@protoc \
	--proto_path=protobuf "protobuf/user.proto" \
	--go_out=services/common/genproto/user \
	--go_opt=paths=source_relative \
	--go-grpc_out=services/common/genproto/user \
	--go-grpc_opt=paths=source_relative \

proto-notification:
	@protoc \
	--proto_path=protobuf "protobuf/notification.proto" \
	--go_out=services/common/genproto/notification \
	--go_opt=paths=source_relative \
	--go-grpc_out=services/common/genproto/notification \
	--go-grpc_opt=paths=source_relative \

proto-chat:
	@protoc \
	--proto_path=protobuf "protobuf/chat.proto" \
	--go_out=services/common/genproto/chat \
	--go_opt=paths=source_relative \
	--go-grpc_out=services/common/genproto/chat \
	--go-grpc_opt=paths=source_relative \

run-auth:
	@go run ./services/3-auth/*.go
