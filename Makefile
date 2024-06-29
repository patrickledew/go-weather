SRC_DIR := ./src
OUT_DIR := ./bin
CC := go build
CFLAGS := -o $(OUT_DIR)/weather.exe

SOURCES := $(wildcard $(SRC_DIR)/*.go)

all: $(SOURCES) clean
	$(CC) $(CFLAGS) $(SOURCES) 

clean:
	if [ -d "$(OUT_DIR)" ]; then \
		rm -r $(OUT_DIR); \
	fi

run: all
	./bin/weather.exe