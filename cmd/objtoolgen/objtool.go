//go:generate res2go -package main -prefix Rsrc -o resources.go resources/*.template
package main

// TODO : embed schema
// TODO : add option in generated tool to dump schema

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Axili39/oastools/oasmodel"
	"github.com/Axili39/oastools/protobuf"
)

type genCtx struct {
	Package   string
	Component string
}

func (g *genCtx) generate(wr io.Writer) error {
	fileTemplate := template.Must(template.New("").Parse(string(RsrcFiles["resources/objtool.go.template"])))

	err := fileTemplate.Execute(wr, g)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing template", err)
	}
	// TODO format generated code with go/format
	return err
}

func genProto(file string, protofilename string, component string) {

	// Step 1: Generate .proto with oa2proto
	w, err := os.Create(protofilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		os.Exit(1)
	}
	defer w.Close()

	oa := oasmodel.OpenAPI{}
	err = oa.Load(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s : %v", file, err)
		os.Exit(1)
	}
	if oa.Components.Schemas[component] == nil {
		fmt.Fprintf(os.Stderr, "component %s doesn't exists, candidate are :\n", component)
		for k := range oa.Components.Schemas {
			fmt.Fprintln(os.Stderr, "\t", k)
		}
		os.Exit(1)
	}

	err = protobuf.Components2Proto(&oa, w, "foo.bar", nil, "go_package=\".;main\"")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %s : %v", file, err)
		os.Exit(1)
	}
}

func compileProto(protofilename string, directory string) {
	cmd := exec.Command("protoc", "--go_out="+directory, protofilename)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running protoc command %v\n", err)
		os.Exit(1)
	}
}

func genCfgTool(directory string, component string) {
	wr, err := os.Create(directory + "/main.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
	defer wr.Close()
	// Verify that component exists
	g := genCtx{"main", component}

	err = g.generate(wr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
}

func buildCfgTool(directory string) {
	cmd := exec.Command("go", "build")
	cmd.Dir = "./" + directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running build command %v\n", err)
		os.Exit(1)
	}
}

func main() {
	var file = flag.String("f", "", "oas file")
	var component = flag.String("c", "", "component name in spec file")
	var outputfile = flag.String("o", "", "output filename")
	var build = flag.Bool("build", false, "build tool with go build")

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
	protofilename := output + "/" + output + ".proto"

	// create directory
	err := os.MkdirAll(output, 0750)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s : %v", *file, err)
		os.Exit(1)
	}

	// Step 1: gen .proto file
	genProto(*file, protofilename, *component)

	// Step 2: Generate package with protoc
	compileProto(protofilename, output)

	// Step 3: Generate filetoolcmd for package
	genCfgTool(output, *component)

	// Step 4: Build if requested
	if *build {
		buildCfgTool(output)
	}
}
