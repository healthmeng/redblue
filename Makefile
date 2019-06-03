export GOPATH=$(PWD)
all:
	go build -gcflags "-N" rbcalc
