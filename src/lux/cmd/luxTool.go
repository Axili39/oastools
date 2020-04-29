package main

import(
	"fmt"
	"os"
	"flag"
	"../../pkg/oatool"
	"google.golang.org/protobuf/proto"
	"../lux"
)



/* Package */
func GetObjByName(node string) proto.Message {
	switch node {
	
		case "Topology":
			var obj lux.Topology
			return &obj
	
		case "TopologyStatus":
			var obj lux.TopologyStatus
			return &obj
	
		case "Error":
			var obj lux.Error
			return &obj
	
		case "TopologyDef":
			var obj lux.TopologyDef
			return &obj
	
		case "Node":
			var obj lux.Node
			return &obj
	
		case "Network":
			var obj lux.Network
			return &obj
	
	}
	return nil
}


var strUsage = 
`
  -f string
        input file .json/.yaml/.bin
  -g    generate empty file
  -o string
        json|yaml|bin (default "bin")
  -r string : 
		Topology
		TopologyStatus
		Error
		TopologyDef
		Node
		Network
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,"Usage of %s:\n%s\n", os.Args[0], strUsage)
}
	
	oatool.MainOAFileTool(GetObjByName)
}



