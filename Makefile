.PHONY: all clean bundle build test frontend

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -tags netgo -o spolyr cmd/main.go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -a -tags netgo -o spolyr.exe cmd/main.go

build: build-linux build-windows

clean:
	rm -f spolyr 
	rm -f spolyr.exe
	rm -f spolyr-linux-amd64.tar.gz
	rm -f spolyr-windows-amd64.tar.gz
	rm -rf dist

bundle: build
	mkdir -p ./dist
	tar -czvf dist/spolyr-linux-amd64.tar.gz public spolyr
	tar -czvf dist/spolyr-windows-amd64.tar.gz public spolyr.exe

test:
	DATABASE_USER=root DATABASE_PASSWORD=example DATABASE_HOST=127.0.0.1 go test -coverprofile cover.out ./internal/...
	go tool cover -html=cover.out -o cover.html

test-units:
	go test -short ./internal/...

coverage: test
	go tool cover -func cover.out | tail -n 1 | awk '{print $3}'

node_modules:
	npm ci

frontend: node_modules
	npm run build

lint-frontend: node_modules
	npm run lint