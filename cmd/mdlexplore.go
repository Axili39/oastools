package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"

	"../pkg/oasjstree"
)

var help = flag.Bool("h", false, "show help")
var file = flag.String("f", "file", "model file")
var root = flag.String("root", "", "root object to explore")
var output = flag.String("output", "", "output to file instead of stdout")
var outJSON = flag.Bool("json", false, "output to json")
var outHTML = flag.Bool("html", false, "output to html")
var unfold = flag.Bool("unfold", false, "")

//Usage
func Usage() {
	/* options :
	-h		Usage
	-f 		file to open
	-o		output filename
	-root	root node in components oas section
	-json 	output json data
	-html	output sample html explorer
	-unfold	default unfold the entire tree
	*/

	/* mdlexplorer-server :
	-server launch web server to generate interactive
	-bind	server bind info, default : 0.0.0.0:8096
	*/
	flag.Usage()
}

type genCtx struct {
	Title string
	Data  *oasjstree.JstNode
}

func genHTML(node *oasjstree.JstNode, out *os.File) {
	t, err := template.ParseFiles("templates/mexplorer.html")

	if err != nil {
		fmt.Println("unable to open template")
		return
	}

	g := genCtx{node.Text, node}
	err = t.ExecuteTemplate(out, "mexplorer.html", g)
}

func main() {
	flag.Parse()

	// Usage
	if *help {
		Usage()
		return
	}

	// Not server
	// Generate
	node := oasjstree.GetJstree("test.yaml", "TopologyDef", *unfold)

	b, err := json.Marshal(node)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Output
	f := os.Stdout
	if *output != "" {
		f, err = os.Create(*output)
		if err != nil {
			fmt.Println("unable to open output")
			return
		}
	}
	defer f.Close()

	if *outJSON {
		f.Write(b)
	}

	if *outHTML {
		genHTML(node, f)
	}

}
