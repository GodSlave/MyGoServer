--- Makefile ----

# This how we want to name the binary output
BINARY=server

# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS_f1=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

# Builds the project
build:
    go build ${LDFLAGS_f1} -o ${BINARY}_f1


# Installs our project: copies binaries
install:
    go install ${LDFLAGS_f1}

# Cleans our project: deletes binaries
clean:
    if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install