# source : bin
./demoTool -r A -f data.bin -o json > /dev/null
diff data.json output.json
[ $? -ne 0 ] && exit "bin test error"
./demoTool -r A -f data.bin -o yaml > /dev/null
diff data.yaml output.yaml
[ $? -ne 0 ] && exit "bin test error"
./demoTool -r A -f data.bin -o bin > /dev/null
diff data.bin output.bin
[ $? -ne 0 ] && exit "bin test error"

# source json
./demoTool -r A -f data.json -o json > /dev/null
diff data.json output.json
[ $? -ne 0 ] && exit "json test error"
./demoTool -r A -f data.json -o yaml > /dev/null
diff data.yaml output.yaml
[ $? -ne 0 ] && exit "json test error"
./demoTool -r A -f data.json -o bin > /dev/null
diff data.bin output.bin
[ $? -ne 0 ] && exit "json test error"

#source yaml
./demoTool -r A -f data.yaml -o json > /dev/null
diff data.json output.json
[ $? -ne 0 ] && exit "yaml test error"
./demoTool -r A -f data.yaml -o yaml > /dev/null
diff data.yaml output.yaml
[ $? -ne 0 ] && exit "yaml test error"
./demoTool -r A -f data.yaml -o bin > /dev/null
diff data.bin output.bin
[ $? -ne 0 ] && exit "yaml test error"
