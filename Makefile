GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

export GO111MODULE=on

build:
	$(GOGET)
	$(GOBUILD) homekit-geiger-counter.go

run:
	$(GORUN) homekit-geiger-counter.go
