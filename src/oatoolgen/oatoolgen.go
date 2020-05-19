package main

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
	Components []string
}

func main() {
	var file = flag.String("f", "", "oas file")
	flag.Parse()

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
	generator, err := template.ParseFiles("../templates/spectool.go")
	g := genCtx{packageName, nil}
	for k := range oa.Components.Schemas {
		fmt.Println(k)
		g.Components = append(g.Components, k)
	}
	err = generator.ExecuteTemplate(wr, "spectool.go", g)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
	}
}
