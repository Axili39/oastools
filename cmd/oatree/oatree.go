package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Axili39/oastools/asciitree"
	"github.com/Axili39/oastools/oasmodel"
)

func main() {
	var root = flag.String("r", "root", "root node")
	var file = flag.String("f", "test.yaml", "file")
	flag.Parse()
	//Test Function (to be removed)
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 60, 80, 0, '\t', 0)

	defer w.Flush()

	oa := oasmodel.OpenAPI{}
	err := oa.Load(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s : %v", *file, err)
		os.Exit(1)
	}
	asciitree.Components2AscTree(&oa, w, *root)
}
