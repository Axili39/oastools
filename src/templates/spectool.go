package main

import(
	"fmt"
	"os"
	"flag"
	"github.com/Axili39/oastools/oatool"
	"google.golang.org/protobuf/proto"
	"github.com/Axili39/oastools/{{.Package}}/{{.Package}}"
)

{{ $package := .Package }}

/* Package */
func GetObjByName(node string) proto.Message {
	switch node {
	{{ range $val := .Components }}
		case "{{ $val }}":
			var obj {{ $package }}.{{ $val }}
			return &obj
	{{end}}
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
  -r string : {{ range $val := .Components }}
		{{ $val }}{{ end }}
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,"Usage of %s:\n%s\n", os.Args[0], strUsage)
}
	
	oatool.MainOAFileTool(GetObjByName)
}



