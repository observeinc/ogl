all: build

.PHONY: build
build: plugin
	go build ./cmd/ogl/main.go

.PHONY: plugin
plugin:
	go build -buildmode=plugin -o ogl.so ./plugin/plugin.go

.PHONY: test
test:
	go test -modcacherw -v ./...

.PHONY: clean
clean:
	rm -f ogl.so ogl
