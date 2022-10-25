# Copyright 2015 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO     := go
GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
PROMU  := $(GOPATH)/bin/promu
pkgs    = $(shell $(GO) list ./... | grep -v /vendor/)

PREFIX              ?= $(shell pwd)
OUTPREFIX           ?= $(shell pwd)/cmd/vertica-prometheus-exporter/
BIN_DIR             ?= $(shell pwd)/cmd/vertica-prometheus-exporter/
DOCKER_IMAGE_NAME   ?= vertica-prometheus-exporter
DOCKER_IMAGE_TAG    ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

VERSION := v1.0.2

all: format test build


test:
	@echo ">> running tests"
	@$(GO) test -short -race $(PREFIX)/test

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)


build: promu
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(OUTPREFIX)
#	@mv  vertica-prometheus-exporter cmd/vertica-prometheus-exporter/

tarball: promu
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(OUTPREFIX) $(BIN_DIR)

promu:
	@echo $(VERSION) > VERSION
	@GOOS=$(shell uname -s | tr A-Z a-z) \
		GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
		$(GO) install github.com/prometheus/promu@v0.13.0

.PHONY: all format build test promu tarball
