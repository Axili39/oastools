demoTool:
	PATH=${PATH}:~/go/bin ../../bin/objtoolgen -f demo.yaml -c A
	cd demo && go build -o ../$@
test: demoTool
	./test.sh
clean:
	rm -f demoTool
	rm -f output.* 
	rm -rf cmd
	rm -rf demo


