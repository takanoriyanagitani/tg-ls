#!/bin/sh

wazero \
	run \
	-env ENV_DIR_NAME=/guest.d \
	-mount "${PWD}:/guest.d:ro" \
	./tgls.wasm
