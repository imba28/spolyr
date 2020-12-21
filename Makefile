.PHONY: all clean bundle
all: static/css/main.css teamspeak3-viewer

build:
	go build -ldflags "-linkmode external -extldflags -static" -a -o spolyr cmd/main.go

clean:
	rm teamspeak3-viewer

bundle: build
	tar -czvf spolyr-linux-amd64.tar.gz public spolyr template