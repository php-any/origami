.PHONY: build test run-chat run-hello clean

# 必须在 extensions/wails 目录下执行（不要在仓库根目录 go build -o wails ./cmd/）
build:
	go build -o wailsrunner ./cmd/

test:
	go test ./...

run-hello: build
	./wailsrunner examples/hello.php

run-chat: build
	./wailsrunner examples/chat_demo.php

clean:
	rm -f wailsrunner wails
