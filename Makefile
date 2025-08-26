# ===== Settings =====
PROTO_DIR := proto
SERVER_OUT := .
CLIENT_DIR := client
CLIENT_GEN := $(CLIENT_DIR)/src/gen
PROTOC_VER := 29.3
GRPC_WEB_VER := 1.5.0

# protoc plugins (Go 쪽). 로컬에 설치되어 있어야 함:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
PROTOC_GEN_GO := $(shell which protoc-gen-go 2>/dev/null)
PROTOC_GEN_GO_GRPC := $(shell which protoc-gen-go-grpc 2>/dev/null)

# ===== Public targets =====
.PHONY: gen-go gen-web gen-all watch clean

## 서버용 Go 스텁 생성 (.pb.go / _grpc.pb.go)
gen-go:
ifndef PROTOC_GEN_GO
	@echo "error: protoc-gen-go not found. Run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" && exit 1
endif
ifndef PROTOC_GEN_GO_GRPC
	@echo "error: protoc-gen-go-grpc not found. Run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest" && exit 1
endif
	@mkdir -p $(SERVER_OUT)
	protoc -I $(PROTO_DIR) \
	  $(PROTO_DIR)/*.proto \
	  --go_out=$(SERVER_OUT) \
	  --go-grpc_out=$(SERVER_OUT)

## 브라우저용 gRPC-Web 스텁 생성 (컨테이너 안에서 실행: 맥/의존성 이슈 회피)
gen-web:
	@mkdir -p $(CLIENT_GEN)
	docker run --rm --platform=linux/arm64/v8 \
	  -v "$(PWD)/$(PROTO_DIR):/proto" \
	  -v "$(PWD)/$(CLIENT_DIR):/client" \
	  alpine:3.20 sh -lc '\
	    set -e; apk add --no-cache curl unzip; \
	    PROTOC_VER=$(PROTOC_VER); \
	    mkdir -p /tmp/protoc && cd /tmp/protoc; \
	    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v$$PROTOC_VER/protoc-$$PROTOC_VER-linux-aarch_64.zip; \
	    unzip -q protoc-$$PROTOC_VER-linux-aarch_64.zip -d /usr/local; \
	    curl -L https://github.com/grpc/grpc-web/releases/download/$(GRPC_WEB_VER)/protoc-gen-grpc-web-$(GRPC_WEB_VER)-linux-aarch64 -o /usr/local/bin/protoc-gen-grpc-web; \
	    chmod +x /usr/local/bin/protoc-gen-grpc-web; \
	    /usr/local/bin/protoc -I /proto /proto/*.proto \
	      --grpc-web_out=import_style=commonjs,mode=grpcwebtext:/client/src/gen \
	  '

## 둘 다
gen-all: gen-go gen-web

## 파일 변경 감지 자동 재생성 (entr 필요: brew install entr)
watch:
	@ls $(PROTO_DIR)/*.proto | entr -r make gen-all

## 생성물 정리(원하면 사용)
clean:
	@rm -rf $(CLIENT_GEN)
	@find . -name "*_pb.go" -o -name "*_grpc.pb.go" | xargs rm -f