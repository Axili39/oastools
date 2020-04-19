cycletest()
{
	format=$1
	./demoTool -r A -f data.yaml -o $format > /dev/null
	./demoTool -r A -f output.$format -o yaml > /dev/null
	diff data.yaml output.yaml
	[ $? -ne 0 ] && exit "$format test error"
}
# consider reference : yaml
cycletest json
cycletest bin
cycletest yaml
exit 0
