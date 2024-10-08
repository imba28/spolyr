.PHONY: all clean bundle build test frontend

open_api_tools_version=v6.2.0

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -tags netgo -o spolyr ./main.go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -a -tags netgo -o spolyr.exe ./main.go

build: build-linux build-windows

clean:
	rm -f spolyr 
	rm -f spolyr.exe
	rm -f spolyr-linux-amd64.tar.gz
	rm -f spolyr-windows-amd64.tar.gz
	rm -rf dist
	rm -rf pkg/openapi
	rm -rf assets/openapi

bundle: build frontend
	mkdir -p ./dist
	tar -czvf dist/spolyr-linux-amd64.tar.gz public spolyr
	tar -czvf dist/spolyr-windows-amd64.tar.gz public spolyr.exe

test:
	DATABASE_USER=root DATABASE_PASSWORD=example DATABASE_HOST=127.0.0.1 go test -coverprofile cover.out ./pkg/...
	go tool cover -html=cover.out -o cover.html

test-units:
	go test -short ./pkg/...

test-e2e: build-linux frontend
	DATABASE_USER=root DATABASE_PASSWORD=example DATABASE_HOST=127.0.0.1 ./spolyr fixtures
	DATABASE_USER=root DATABASE_PASSWORD=example DATABASE_HOST=127.0.0.1 ./spolyr web > /tmp/backend.log 2>&1 &
	npm run test:e2e:ci:chromium
	killall spolyr
	cat /tmp/backend.log

coverage: test
	go tool cover -func cover.out | tail -n 1 | awk '{print $3}'

node_modules:
	npm ci

frontend: node_modules
	npm run build

lint-frontend: node_modules
	npm run lint

test-frontend: node_modules
	npm run test:unit

openapi-spec:
	docker run --rm -v "${PWD}/oapi-spec.yaml:/local/oapi-spec.yaml" -v "${PWD}/pkg/openapi/:/local/pkg/openapi/openapi" openapitools/openapi-generator-cli:$(open_api_tools_version) generate -g go-server -i /local/oapi-spec.yaml -o /local/pkg/openapi --additional-properties=outputAsLibrary=true,onlyInterfaces=true,sourceFolder=openapi,addResponseHeaders=true
	sudo chown -R $(USER): pkg/openapi/
	sed -i -e 's/"github.com\/gorilla\/mux"//g' pkg/openapi/api_auth.go
	sed -i -e 's/"encoding\/json"//g' pkg/openapi/api_import.go
	sed -i -e 's/"encoding\/json"//g' -e 's/"github.com\/gorilla\/mux"//g' pkg/openapi/api_playlists.go
	docker run --rm -v "${PWD}/oapi-spec.yaml:/local/oapi-spec.yaml" -v "${PWD}/assets/openapi/:/local/assets/openapi/src" openapitools/openapi-generator-cli:$(open_api_tools_version) generate -g javascript -i /local/oapi-spec.yaml -o /local/assets/openapi --additional-properties=usePromises=true,moduleName=@/openapi --global-property models,modelTests=false --global-property apis,apiTests=false --global-property supportingFiles

