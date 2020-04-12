all: mdlexplore 
mdlexplore: cmd/mdlexplore.go pkg/oasjstree/oasjstree.go
	go build cmd/mdlexplore.go
