# oastools
Open Api Specs tools

First clone project:
```
git clone https://github.com/Axili39/oastools
```

Build Project:
```
make
```

Build zip bundle to export portable App (if wanted)
```
make bundle
```

Data Model Explorer
-------------------
usage
```
./mdlexplore -h

Usage of ./mdlexplore:
  -bind string
    	HTTP Server address (default "0.0.0.0:8096")
  -f string
    	model file (default "file")
  -h	show help
  -html
    	output to html
  -json
    	output to json
  -output string
    	output to file instead of stdout
  -root string
    	root object to explore
  -server
    	start HTTP Server
  -unfold
```

generate JSON output from command line 
```
./mdlexplore -f test/test.yaml -root TopologyDef -json
```

generate HTML output from command line 
```
./mdlexplore -f test/test.yaml -root TopologyDef -html
```

start interactive mode withembedded web-server
```
./mdlexplore -server
```


