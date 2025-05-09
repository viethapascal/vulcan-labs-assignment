# Tools
GO := go
all: deps build
deps:
	$(GO) mod tidy
	${GO} mod download
build:
	$(GO) build -o cine
protob:
	@echo "---> Building protocol buffers"
	protoc --go_out=. --go_opt=paths=source_relative \
 	--go-grpc_opt=paths=source_relative --go-grpc_out=. proto/seat.proto

.PHONY: protob build