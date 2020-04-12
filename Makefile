all: mdlexplore 
clean:
	rm -f bundle.zip
	rm -f mdlexplore
mdlexplore: cmd/mdlexplore.go pkg/oasjstree/oasjstree.go
	go build cmd/mdlexplore.go
bundle: mdlexplore
	mkdir -p /tmp/mdlbuild
	cp -r dist /tmp/mdlbuild
	cp -r templates /tmp/mdlbuild
	cp mdlexplore* /tmp/mdlbuild
	cd /tmp/mdlbuild && zip -r ${PWD}/bundle.zip *  && cd -
	rm -rf /tmp/mdlbuild

