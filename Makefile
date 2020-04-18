all: mdlexplore oa2proto
clean:
	rm -f bundle.zip
	rm -f mdlexplore
oa2proto: cmd/oa2proto.go pkg/protobuf/protobuf.go
	go build cmd/oa2proto.go

mdlexplore: cmd/mdlexplore.go pkg/oasjstree/oasjstree.go
	go build cmd/mdlexplore.go
	GOOS=windows GOARCH=386 go build -o mdlexplore.exe cmd/mdlexplore.go
	
bundle: mdlexplore
	mkdir -p /tmp/mdlbuild
	cp -r dist /tmp/mdlbuild
	cp -r templates /tmp/mdlbuild
	cp mdlexplore* /tmp/mdlbuild
	cd /tmp/mdlbuild && zip -r ${PWD}/bundle.zip *  && cd -
	rm -rf /tmp/mdlbuild

