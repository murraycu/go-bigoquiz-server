.PHONY: all build clean deploy

all: build

build:
	go build

# This runs only the "short" tests.
# (not the tests that require more setup, such as the datastore emulator.)
test:
	go test ./... -short -cover

# This runs all tests,
# including the ones that require more setup, such as the datastore emulator.)
# TODO: Stop the datastore emulator after the tests have run.
# This also outputs a coverage file and processes it to produce a coverage.html report.
full-test:
	gcloud config set project bigoquiz ; \
	(gcloud beta emulators datastore start --no-store-on-disk & ) ; \
	export DATASTORE_EMULATOR_HOST="localhost:8081" ; \
	go test ./... -coverprofile=coverage.out ; \
	go tool cover -html=coverage.out -o coverage.html

clean:
	go clean

deploy:
	gcloud app deploy .

format:
	go fmt ./...

local_run: build
	(./start_datastore_emulator.sh & ) ; \
	export DATASTORE_EMULATOR_HOST="localhost:8025" ; \
        go run .

