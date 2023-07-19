CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=route256/checkout/cmd/app
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


generate: vendor-proto/validate vendor-proto/google/api vendor-proto/google/protobuf generate-service generate-loms generate-product

generate-service: 
	mkdir -p pkg/checkout_v1 
	$(PROTOC) -I api/checkout/v1 -I vendor-proto \
	--go_out pkg/checkout_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/checkout_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/checkout_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/checkout_v1" \
	api/checkout/v1/checkout_service.proto

generate-loms: 
	mkdir -p pkg/loms_v1
	$(PROTOC) -I api/loms/v1 -I vendor-proto \
	--go_out pkg/loms_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/loms_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/loms_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/loms_v1" \
	api/loms/v1/loms_client.proto

generate-product:
	mkdir -p pkg/product_v1
	$(PROTOC) -I api/product/v1 -I vendor-proto \
	--go_out pkg/product_v1 --go_opt paths=source_relative \
	--go-grpc_out pkg/product_v1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out pkg/product_v1 --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:pkg/product_v1" \
	api/product/v1/product_service.proto

tidy:
	go mod tidy

integration-test: down-docker up-docker wait-for-docker
	./migration.sh && \
	cd internal/repository/postgres/integration_test && \
	go test && docker-compose down -v


up-docker:
	cd internal/repository/postgres/integration_test && \
	docker-compose up -d --force-recreate --build

down-docker:
	cd internal/repository/postgres/integration_test && \
	docker-compose down -v

wait-for-docker:
	cd internal/repository/postgres/integration_test && \
	until docker-compose exec -T checkout-db pg_isready -U user -d checkout; do \
		sleep 1; \
	done

test-checkout-package:
	cd internal/service/checkout && \
	go test -coverprofile=cov.out && go tool cover -html=cov.out -o cov.html