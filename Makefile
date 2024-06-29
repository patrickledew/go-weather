SRC_DIR := ./src
OUT_DIR := ./bin
CC := go build
CFLAGS := -o $(OUT_DIR)/weather.exe

SOURCES := $(wildcard $(SRC_DIR)/*.go)

all: $(SOURCES) clean
	$(CC) $(CFLAGS) $(SOURCES) 

clean:
	rm -r ./bin

run: all
	./bin/weather.exe