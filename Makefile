# SPDX-FileCopyrightText: 2026 <Your Name>
# SPDX-License-Identifier: GPL-3.0-or-later

.PHONY: all build run clean

BINARY=cmatrix

all: build

build:
	go build -o $(BINARY) cmatrix.go

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
