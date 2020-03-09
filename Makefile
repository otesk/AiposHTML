BIN_DIR:=./bin
HTML:=$(BIN_DIR)/html

build: $(HTML)

$(HTML):
	go build -ldflags "-linkmode external -extldflags -static" -o $(HTML) cmd/service/main.go

.PHONY: dep
dep:
	go get ./...
