#!/bin/sh

tinygo \
	build \
	-o ./tgls.wasm \
	-target=wasip1 \
	-opt=z \
	-no-debug \
	./ls.go
