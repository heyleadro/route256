CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=route256/notifications/cmd/app

include ../protomake/Makefile

PROTOC = PATH="$$PATH:$(BINDIR)" protoc

all: format build test lint

build: bindir generate tidy
	go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(BINDIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(BINDIR) go install github.com/pressly/goose/v3/cmd/goose@latest

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

test-notifications-package:
	cd internal/service/notifications && \
	go test -coverprofile=cov.out && go tool cover -html=cov.out -o cov.html


generate: vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/validate
	mkdir -p pkg/notifications_v1
	$(PROTOC) -I api/notifications/v1 -I vendor-proto \
	--go_out pkg/notifications_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/notifications_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/notifications_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/notifications_v1" \
	api/notifications/v1/service.proto

tidy:
	go mod tidy