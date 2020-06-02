//go:install go generate
//go:generate res2go -package main -prefix Rsrc -o resources.go resources/*.template
package main

// TODO : directly generate binary
// TODO : supports command name choice
import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/Axili39/oastools/oasmodel"
	"github.com/Axili39/oastools/protobuf"
)

type genCtx struct {
	Package    string
	Component  string
}

func (g *genCtx) generate(wr *os.File) error {
	fileTemplate := template.Must(template.New("").Parse(string(RsrcFiles["resources/objtool.go.template"])))
	
	err := fileTemplate.Execute(wr,g)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing template", err)
	}
	// TODO format generated code with go/format
	return err
}


func main() {
	var file = flag.String("f", "", "oas file")
	var component = flag.String("c", "", "component name in spec file")
	var outputfile = flag.String("o", "", "output filename")

	flag.Parse()
	RsrcInit()

	// Some Checks
	if *file == "" {
		fmt.Fprintln(os.Stderr, "missing spec file")
		os.Exit(1)
	}

	if *component == "" { 
		fmt.Fprintln(os.Stderr, "missing component name")
		os.Exit(1)
	}

	// If package name isn't specified, we take spec file name
	var output string
	if *outputfile == "" {
		sl := strings.Split(*file, ".yaml")
		output = sl[0]	
	} else {
		output = *outputfile
	}

	// create directory
	os.MkdirAll(output, 0750)
	
	protofilename := output + "/" + output + ".proto"
	// Step 1: Generate .proto with oa2proto
	w, err := os.Create(protofilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}
	defer w.Close()

	oa := oasmodel.OpenAPI{}
	oa.Load(*file)

	if oa.Components.Schemas[*component] == nil {
		fmt.Fprintf(os.Stderr, "component %s doesn't exists, candidate are :\n", *component)
		for k := range oa.Components.Schemas {
			fmt.Fprintln(os.Stderr,"\t",k)			
		}		
		os.Exit(1)
	}

	protobuf.Components2Proto(&oa, w, "main")

	//Step 2: Generate package with protoc
	cmd := exec.Command("protoc", "--go_out=.", protofilename)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running protoc command %v\n", err)
		os.Exit(1)
	}

	// Step 3: Generate filetoolcmd for package
	wr, err := os.Create(output +  "/main.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
	defer wr.Close()
	// Verify that component exists
	g := genCtx{"main", *component}

	err = g.generate(wr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
}
