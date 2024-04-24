GREP_NOT_TESTS= { grep -v 'no test files'; true; }

unit:
	go test -failfast -count 1 ./... -tags unit_tests | ${GREP_NOT_TESTS}

smoke:
	go test -failfast -count 1 ./... -tags smoke_test | ${GREP_NOT_TESTS}

integration:
	go test -failfast -count 1 ./... -tags integration_tests | ${GREP_NOT_TESTS}

test:
	clear
	$(MAKE) unit
	$(MAKE) smoke
	$(MAKE) integration

build:
	go build -v -o server

dev:
	APP_ENV=development
	$(MAKE) build
	./server

prod:
	APP_ENV=production
	$(MAKE) build
	./server