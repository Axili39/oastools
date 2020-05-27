package main
import (
	"github.com/Axili39/oastools/oasmodel"
	"github.com/Axili39/oastools/asciitree"
	"os"
	"flag"
	"text/tabwriter"
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
	oa.Load(*file)
	asciitree.Components2AscTree(&oa, w, *root)
}