//go:generate sh -c "cd resources ; res2go *.template"
package main

// TODO : supports -r root node preselection
// TODO : supports -output package
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

	"github.com/Axili39/oastools/cmd/oatoolgen/resources"
	"github.com/Axili39/oastools/oasmodel"
	"github.com/Axili39/oastools/protobuf"
)

type genCtx struct {
	Package    string
	Components []string
}

func (g *genCtx) generate(wr *os.File) error {
	fileTemplate := template.Must(template.New("").Parse(string(resources.Files["spectool.go.template"])))
	
	err := fileTemplate.Execute(wr,g)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing template", err)
	}
	// TODO format generated code with go/format
	return err
}


func main() {
	var file = flag.String("f", "", "oas file")
	flag.Parse()
	resources.Init()

	sl := strings.Split(*file, ".yaml")
	packageName := sl[0]
	// create directory
	os.MkdirAll(packageName, 0750)
	os.MkdirAll("cmd", 0750)
	protofilename := packageName + "/" + packageName + ".proto"
	// Step 1: Generate .proto with oa2proto
	w, err := os.Create(protofilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
	defer w.Close()

	oa := oasmodel.OpenAPI{}
	oa.Load(*file)
	protobuf.Components2Proto(&oa, w, packageName)

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
	wr, err := os.Create("cmd/" + packageName + "Tool.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
	defer wr.Close()

	g := genCtx{packageName, nil}
	for k := range oa.Components.Schemas {
		fmt.Println(k)
		g.Components = append(g.Components, k)
	}

	err = g.generate(wr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
}
