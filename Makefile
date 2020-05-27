ifeq ($(GOOS),windows)
	EXT=.exe
	ifeq ($(GOARCH),)
		GOARCH=amd64
	endif
endif
BIN=$(shell pwd)/bin
all: oa2proto oatoolgen oatree oa2gowapi
clean:
	rm -f bin/*
install: all
	cp bin/* ~/go/bin

oa2proto: cmd/oa2proto/oa2proto.go
	cd cmd/$@ && go build -o ${BIN}/$@${EXT}

oatoolgen: cmd/oatoolgen/oatoolgen.go
	cd cmd/$@ && go generate
	cd cmd/$@ && go build -o ${BIN}/$@${EXT}

oatree: cmd/oatree/oatree.go
	cd cmd/$@ && go build -o ${BIN}/$@${EXT}

oa2gowapi: cmd/oa2gowapi/oa2gowapi.go
	cd cmd/$@ && go generate
	cd cmd/$@ && go build -o ${BIN}/$@${EXT}
