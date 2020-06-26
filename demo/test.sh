cycletest()
{
	format=$1
	./demoTool -if data.yaml -of output.$format > /dev/null
	./demoTool -if output.$format -of output.yaml > /dev/null
	diff data.yaml output.yaml
	[ $? -ne 0 ] && exit "$format test error"
}
# consider reference : yaml
cycletest json
cycletest bin
cycletest yaml
exit 0
