.PHONY: all clean bundle build
all: teamspeak3-viewer

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o spolyr cmd/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -a  -o spolyr.exe cmd/main.go

clean:
	rm -f spolyr 
	rm -f spolyr.exe
	rm -f spolyr-linux-amd64.tar.gz
	rm -f spolyr-windows-amd64.tar.gz
	rm -rf dist

bundle: build
	mkdir -p ./dist
	tar -czvf dist/spolyr-linux-amd64.tar.gz public spolyr template docker-compose.yml
	tar -czvf dist/spolyr-windows-amd64.tar.gz public spolyr.exe template docker-compose.yml