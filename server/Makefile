.PHONY: build doc fmt lint run test vendor_clean vendor_get vendor_update vet

# Prepend our _vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/_vendor:${GOPATH}
PROJECT_NAME := TOSET
export GOPATH

default: build

build:
	go build -v -o ./bin/$(PROJECT_NAME) ./src/*.go

doc:
	godoc -http=:6060 -index -goroot .

fmt:
	go fmt ./

lint:
	golint .

SUBDOMAIN=TOSET
launch:
	HOST=$$(hostname); \
	PORT='8080'; \
	./bin/$(PROJECT_NAME) -port=$$PORT

run: build launch

run-dev: build launch
	export MONGO_URL=localhost:27017

test:
	cd src; go test ./routes/

test-verbose:
	cd src; go test -v ./routes/

test-target:
	go test -v -run '$(TARGET)' ./src/routes

cover:
	go test -cover ./src/routes

ROUTES=src/routes/
cover-html:
	go test -coverprofile=cover.out ./src/routes
	sed -i -e "s#.*\($(ROUTES).*\.go\)#\./\\1#" cover.out #Format cover.out
	go tool cover -html=cover.out -o coverage.html
	rm cover.out

cover-func:
	@go test -coverprofile=cover.out ./src/routes
	@sed -i -e "s#.*\($(ROUTES).*\.go\)#\./\\1#" cover.out #Format cover.out
	@go tool cover -func=cover.out
	@rm cover.out

vendor_clean:
	rm -rf ./_vendor/src

# We have to set GOPATH to just the _vendor
# directory to ensure that `go get` doesn't
# update packages in our primary GOPATH instead.
# This will happen if you already have the package
# installed in GOPATH since `go get` will use
# that existing location as the destination.

vendor_get: vendor_clean
	GOPATH=${PWD}/_vendor go get -d -u -v \
		github.com/getsentry/raven-go \
		github.com/gorilla/mux \
		gopkg.in/mgo.v2 \
		github.com/rs/cors \
		github.com/gorilla/handlers \
		github.com/gorilla/securecookie \
		golang.org/x/crypto/bcrypt \
		github.com/k0kubun/pp \
		github.com/kylelemons/godebug/pretty

vendor_update: vendor_get

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ./...
