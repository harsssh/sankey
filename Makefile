ROOT_DIR = $(dir $(abspath $(firstword $(MAKEFILE_LIST))))
GO_DIR=$(ROOT_DIR)src/go
PY_DIR=$(ROOT_DIR)src/python
GO_SRC=$(GO_DIR)/aggregate.go $(GO_DIR)/model/log.go $(GO_DIR)/parser/parser.go
PY_SRC=$(PY_DIR)/plot.py
BIN_NAME=$(ROOT_DIR)bin/sankey

.PHONY: all
all: $(BIN_NAME) plot

$(BIN_NAME): $(GO_SRC)
	@cd $(GO_DIR) && go build -o $@

.PHONY: plot
plot: $(BIN_NAME)
	@$(BIN_NAME) | python $(PY_SRC)

.PHONY: clean
	$(RM) $(BIN_NAME)
