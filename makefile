APPLICATION_NAME=alert-system
GO_BIN=go
VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')

##
## Testing
##---------------------------------------------------------------------------

test:                    ## Run the tests
	$(GO_BIN) test ./...

##
## Project binary build
##---------------------------------------------------------------------------

get-deps:                ## Update the project's dependencies
	$(GO_BIN) get -u

build:                   ## Compile the binary
	@echo "building ${APPLICATION_NAME} ${VERSION}"
	$(GO_BIN) build -o ${APPLICATION_NAME}