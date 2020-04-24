all:oatree  mdlexplore oa2proto oatoolgen 
clean:
	rm -f bundle.zip
	rm -f mdlexplore
	rm -f oa2proto
	rm -f oatoolgen
	rm -f oatree
	rm -f *.exe

oatree: cmd/oatree.go pkg/oasmodel/oasmodel.go pkg/asciitree/asciitree.go
	go build cmd/oatree.go
	GOOS=windows GOARCH=386 go build -o oatree.exe cmd/oatree.go

oatoolgen: cmd/oatoolgen.go pkg/protobuf/protobuf.go pkg/oasmodel/oasmodel.go
	go build cmd/oatoolgen.go
	GOOS=windows GOARCH=386 go build -o oatoolgen.exe cmd/oatoolgen.go

oa2proto: cmd/oa2proto.go pkg/protobuf/protobuf.go
	go build cmd/oa2proto.go
	GOOS=windows GOARCH=386 go build -o oa2proto.exe cmd/oa2proto.go

mdlexplore: cmd/mdlexplore.go pkg/oasjstree/oasjstree.go
	go build cmd/mdlexplore.go
	GOOS=windows GOARCH=386 go build -o mdlexplore.exe cmd/mdlexplore.go
	
bundle: all
	mkdir -p /tmp/mdlbuild
	cp -r dist /tmp/mdlbuild
	cp -r templates /tmp/mdlbuild
	cp mdlexplore* /tmp/mdlbuild
	cp oa2proto* /tmp/mdlbuild
	cp oatoolgen* /tmp/mdlbuild
	cp oatree* /tmp/mdlbuild
	cd /tmp/mdlbuild && zip -r ${PWD}/bundle.zip *  && cd -
	rm -rf /tmp/mdlbuild

