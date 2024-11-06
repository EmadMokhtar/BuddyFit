# Makefile for building Go binaries

# Define the commands
CMDS = fetcher importer buddyfit

# Default target
all: $(CMDS)

# Build each command
$(CMDS):
	go build -o ./bin/$@ ./cmd/$@

# Clean up binaries
clean:
	rm -f bin/*

.PHONY: all clean $(CMDS)