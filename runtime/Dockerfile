# Setup Go configuration variables. (Despite being ARG instructions, these should not be set at runtime via `docker build --build-arg`, the
# default values should be used or changed here).
ARG GO_VERSION=1.15.6
ARG GOPATH=/go
ARG GOROOT_BOOTSTRAP=/usr/local/gotemp/go
ARG GOROOT=/usr/local/go

# Setup PATH.
ARG PATH=${GOPATH}/bin:${GOROOT_BOOTSTRAP}/bin:${GOROOT}/bin:${PATH}

### BUILD STAGE
FROM buildpack-deps:jessie-scm as build

# Download dependencies. gcc is for cgo.
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		patch

# Create the parent directory of GOROOT_BOOTSTRAP.
RUN mkdir -p /usr/local/gotemp

# Download and install the Go binary to be used for bootstrapping.
ENV GO_BOOTSTRAP_DOWNLOAD_URL=https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
ENV GO_BOOTSTRAP_SHA256=3918e6cc85e7eaaa6f859f1bdbaac772e7a825b0eb423c63d3ae68b21f84b844
RUN curl -fsSL ${GO_BOOTSTRAP_DOWNLOAD_URL} -o bootstrap.tar.gz \
	&& echo ${GO_BOOTSTRAP_SHA256} bootstrap.tar.gz | sha256sum -c - \
	&& tar -C /usr/local/gotemp -xzf bootstrap.tar.gz

# Download the Go source code to be patched.
ENV GO_DOWNLOAD_URL=https://golang.org/dl/go${GO_VERSION}.src.tar.gz
ENV GO_SHA256=890bba73c5e2b19ffb1180e385ea225059eb008eb91b694875dd86ea48675817
RUN curl -fsSL ${GO_DOWNLOAD_URL} -o golang.tar.gz \
	&& echo ${GO_SHA256} golang.tar.gz | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz

# Copy the patch file in.
COPY go${GO_VERSION}-gotrace.patch /runtime/gotrace.patch

# Patch the Go source code.
RUN cd ${GOROOT}/src \
	&& patch -p2 < /runtime/gotrace.patch \
	&& ./make.bash

### RELEASE STAGE
FROM alpine:latest

# Setup GOPATH requisite folders.
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin \
	&& chmod -R 777 ${GOPATH}

# Copy in patched Go from build stage.
COPY --from=build ${GOROOT} ${GOROOT}

WORKDIR ${GOPATH}
