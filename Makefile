export GOPATH=$(PWD)
all:
	go build -gcflags "-N" rbcalc
clean:
	rm -f rbcalc
