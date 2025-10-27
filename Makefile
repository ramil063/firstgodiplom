ifneq (,$(wildcard ./.env))
include .env
export
endif

build:
	/usr/local/go/bin/go build -o gophermart cmd/gophermart/main.go

clean:
	rm gophermart

run-accrual:
	./cmd/accrual/accrual_linux_amd64 -a :8081

run-accrual-win:
	./cmd/accrual/accrual_windows_amd64 -a :8081

run:
	./gophermart -d $(DATABASE_URI) -r http://localhost:8081

test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm -f coverage.out

test-cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@open coverage.html || xdg-open coverage.html

.PHONY: build clean run