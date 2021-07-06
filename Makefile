all: build

.PHONY: install
install:
	go install ogl/cmd/ogl

.PHONY: build
build: plugin
	go build ogl/cmd/ogl

.PHONY: plugin
plugin:
	go build -buildmode=plugin -o ogl.so ./plugin/plugin.go

.PHONY: test
test:
	go test -modcacherw -v ogl/...

.PHONY: clean
clean:
	rm -f ogl.so ogl
