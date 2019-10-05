.PHONY: all build clean deploy

all: build

build:
	go build

test:
	go test ./...

clean:
	go clean

deploy:
	gcloud app deploy .

format:
	go fmt ./...

local_run: build
	(gcloud beta emulators datastore start & ) ; \
	export DATASTORE_EMULATOR_HOST="localhost:8081" ; \
        go run .

