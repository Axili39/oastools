ifeq ($(GOOS),windows)
	EXT=.exe
	ifeq ($(GOARCH),)
		GOARCH=amd64
	endif
endif
BIN=$(shell pwd)/bin
all: oa2proto oatoolgen oatree mdlexplore
clean:
	rm -f bin/*

oa2proto: src/oa2proto/oa2proto.go
	cd src/$@ && go build -o ${BIN}/$@${EXT}

oatoolgen: src/oatoolgen/oatoolgen.go
	cd src/$@ && go build -o ${BIN}/$@${EXT}

oatree: src/oatree/oatree.go
	cd src/$@ && go build -o ${BIN}/$@${EXT}

mdlexplore: src/mdlexplore/mdlexplore.go
	cd src/$@ && go build -o ${BIN}/$@${EXT}
