include $(GOROOT)/src/Make.inc

TARG=ggigl
GOFILES=$(wildcard *.go)

include $(GOROOT)/src/Make.cmd

.PHONY: fmt

fmt:
	gofmt -w $(GOFILES)
