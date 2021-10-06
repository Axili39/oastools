demoTool:
	PATH=${PATH}:~/go/bin ../bin/objtoolgen -f demo.yaml -c A --build -o demo
	cp demo/demo ./demoTool
test: demoTool
	./test.sh
clean:
	rm -f demoTool
	rm -f output.* 
	rm -rf cmd
	rm -rf demo


