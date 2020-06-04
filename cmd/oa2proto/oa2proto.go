package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Axili39/oastools/oasmodel"
	"github.com/Axili39/oastools/protobuf"
)

func main() {
	file := flag.String("f", "test.yaml", "yaml file to parse")
	out := flag.String("o", "", "output file")
	packageName := flag.String("p", "", "package name")
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
	err := oa.Load(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s : %v", *file, err)
		os.Exit(1)
	}
	err = protobuf.Components2Proto(&oa, output, *packageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %s : %v", *file, err)
		os.Exit(1)
	}
}
