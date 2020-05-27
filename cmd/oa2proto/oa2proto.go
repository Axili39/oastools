package main

import (
	"github.com/Axili39/oastools/protobuf"
	"github.com/Axili39/oastools/oasmodel"
	"flag"
	"os"
	"fmt"
)

func main() {
	file := flag.String("f", "test.yaml", "yaml file to parse")
	out := flag.String("o", "", "output file")
	packageName := flag.String("p","", "package name")
	flag.Parse()

	var output *os.File
	if *out != "" {
		var err error
		output, err = os.Create(*out)
		if err != nil {
			fmt.Printf("error opening %s : %v", *out, err)
			return
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}
    
	oa := oasmodel.OpenAPI{}
	oa.Load(*file)
	protobuf.Components2Proto(&oa, output, *packageName)
}
