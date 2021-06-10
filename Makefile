#!make

include .env
export $(shell sed 's/=.*//' .env)

GIT_BUILD = $(shell git describe --tags)
SCOPE = github.com/diauweb/xmcl/config

CONSTS = -X ${SCOPE}.GIT_BUILD=${GIT_BUILD} \
		 -X ${SCOPE}.CONFIG_ENDPOINT=${CONFIG_ENDPOINT} \
		 -X ${SCOPE}.PRODUCT_NAME=XMCL

.PHONY: clean

version:
	echo "${GIT_BUILD}"

run:
	go run -ldflags "${CONSTS} -X ${SCOPE}.MODE=DEBUG" .

build:
	GOOS=windows go build -ldflags "${CONSTS}" .

build2:
	go build -ldflags "${CONSTS}" .

clean:
	rm xmcl*.exe
