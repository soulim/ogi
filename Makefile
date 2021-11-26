BIN_DIR = ./bin
EXAMPLES_DIR = ./docs/examples
SOURCES_DIR = $(EXAMPLES_DIR)/sources
SOURCE_FILES = $(shell find $(SOURCES_DIR) -type f -name "*.txt")
EXAMPLE_FILES = $(patsubst $(SOURCES_DIR)/%.txt, $(EXAMPLES_DIR)/%.png, $(SOURCE_FILES))

VERSION = 1.0.0
LDFLAGS = -ldflags=-X=main.version=$(VERSION)

$(BIN_DIR):
	mkdir -p $@

OGI = $(BIN_DIR)/ogi
$(OGI): main.go | $(BIN_DIR)
	go build -o $@ $(LDFLAGS) $?

$(EXAMPLE_FILES): $(EXAMPLES_DIR)/%.png : $(SOURCES_DIR)/%.txt | $(OGI)
	$(OGI) --text="$(file < $?)" \
	       --note="$(shell basename $? .txt)" \
	       --width=1200 \
	       --height=628 \
	       --pattern="nested-squares" \
	> $@

.PHONY: all
all: $(OGI)

.PHONY: examples
examples: $(EXAMPLE_FILES)

.PHONY: clean
clean:
	rm -f $(OGI)
	rm -f $(EXAMPLE_FILES)

.DEFAULT_GOAL := all
