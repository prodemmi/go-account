GREP_NOT_TESTS= { grep -v 'no test files'; true; }
IMAGE_TAG=$(shell git describe --tags --abbrev=0 | cut -c2-)
IMAGE_NAME=ghcr.io/prodemmi/go-account
IMAGE_FULL_NAME=ghcr.io/prodemmi/go-account:${IMAGE_TAG}

build-image:
	@echo ":::building image"
	HTTPS_PROXY=http://fodev.org:8118/
	echo ${CR_PAT} | docker login ghcr.io -u prodemmi --password-stdin
	docker build . -t ${IMAGE_FULL_NAME}
	slim build --http-probe-exit-on-failure=false --tag ${IMAGE_NAME}.slim:${IMAGE_TAG} ${IMAGE_FULL_NAME}
	docker push ${IMAGE_NAME}.slim:${IMAGE_TAG}
	docker rmi ${IMAGE_FULL_NAME}

unit:
	@echo ":::unit testing"
	go test -failfast -count 1 ./... -tags unit_tests | ${GREP_NOT_TESTS}

smoke:
	@echo ":::smoke testing"
	go test -failfast -count 1 ./... -tags smoke_test | ${GREP_NOT_TESTS}

integration:
	@echo ":::integration testing"
	go test -failfast -count 1 ./... -tags integration_tests | ${GREP_NOT_TESTS}

test:
	@echo ":::full testing"information_schema
	clear
	$(MAKE) unit
	$(MAKE) smoke
	$(MAKE) integration

build:
	@echo ":::building"
	go build -v -o server

dev:
	@echo ":::running development"
	APP_ENV=development
	$(MAKE) build
	./server

prod:
	@echo ":::running production"
	APP_ENV=production
	$(MAKE) build
	./server