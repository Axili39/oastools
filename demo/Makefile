demoTool:
	PATH=${PATH}:~/go/bin ../oatoolgen -f demo.yaml
	go build cmd/demoTool.go
test: demoTool
	./test.sh
clean:
	rm -f demoTool
	rm -f output.* 
	rm -rf cmd
	rm -rf demo


