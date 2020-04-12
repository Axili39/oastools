package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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
var server = flag.Bool("server", false, "start HTTP Server")
var bind = flag.String("bind", "0.0.0.0:8096", "HTTP Server address")

//Usage print mdlexplorer help
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

var generator *template.Template

func genHTML(node *oasjstree.JstNode, out *os.File) {
	g := genCtx{node.Text, node}
	err := generator.ExecuteTemplate(out, "mexplorer.html", g)

	if err != nil {
		log.Printf("error generating HTML content %v", err)
	}
}
func httpServeFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
func httpMain(w http.ResponseWriter, r *http.Request) {
	generator.ExecuteTemplate(w, "index.html", nil)
}
func httpUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	file, hdl, err := r.FormFile("oasFile")
	if err != nil {
		log.Printf("Error Retrieving the File : %v\n", err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Successfully Uploaded File : %s\n", hdl.Filename)

	// Retrieve Root Object to explore
	object := r.FormValue("Root")
	log.Printf("Looking for %s node\n", object)
	node := oasjstree.GetJstreeFromData(fileBytes, object, *unfold)
	if node != nil {
		g := genCtx{node.Text, node}
		generator.ExecuteTemplate(w, "mexplorer.html", g)
	} else {
		fmt.Fprintf(w, "Error can't find object : %s", object)
	}
}

//StartHTTP start HTTP server
func StartHTTP() {
	http.HandleFunc("/", httpMain)
	http.HandleFunc("/dist/", httpServeFiles)
	http.HandleFunc("/upload", httpUpload)

	log.Printf("Start OAS model explorer HTTP server : %s\n", *bind)

	err := http.ListenAndServe(*bind, nil)
	if err != nil {
		log.Printf("error starting http server %v", err)
	}
}
func main() {
	var err error
	generator, err = template.ParseGlob("./templates/*")

	if err != nil {
		fmt.Println("unable to open template")
		return
	}

	flag.Parse()

	// Usage
	if *help {
		Usage()
		return
	}

	// HTTP Server Case
	if *server {
		StartHTTP()
		return
	}

	// Not server
	// Generate
	node := oasjstree.GetJstree(*file, *root, *unfold)

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
