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

.PHONY: build clean run