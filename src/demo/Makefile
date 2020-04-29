demoTool:
	PATH=${PATH}:~/go/bin ../../bin/oatoolgen -f demo.yaml
	cd cmd && go build -o ../$@
test: demoTool
	./test.sh
clean:
	rm -f demoTool
	rm -f output.* 
	rm -rf cmd
	rm -rf demo


