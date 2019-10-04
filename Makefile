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

