BIN_DIR = ./bin

VERSION = 1.0.0
LDFLAGS = -ldflags=-X=main.version=$(VERSION)

$(BIN_DIR):
	mkdir -p $@

OGI = $(BIN_DIR)/ogi
$(OGI): main.go | $(BIN_DIR)
	go build -o $@ $(LDFLAGS) $?

.PHONY: all
all: $(OGI)

.PHONY: clean
clean:
	rm -f $(OGI)

.DEFAULT_GOAL := all
