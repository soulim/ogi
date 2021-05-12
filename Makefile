VERSION = 1.0.0
LDFLAGS = -ldflags=-X=main.version=$(VERSION)

ogi: main.go
	go build -o $@ $(LDFLAGS) $?

.PHONY: all
all: ogi

.PHONY: clean
clean:
	rm -f ogi

.DEFAULT_GOAL := all
