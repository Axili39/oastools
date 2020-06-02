//go:generate sh -c "cd resources ; res2go gowapi*.template"

package main

import (
	"github.com/Axili39/oastools/protobuf"
	"github.com/Axili39/oastools/oasmodel"
	"github.com/Axili39/oastools/cmd/oa2gowapi/resources"
	"flag"
	"os"
	"fmt"
	"text/template"
	"strings"
	"os/exec"
	"bytes"
	"go/format"
)

type operation struct {
	Name	string
	PbMessageName string 
}

type entityInterface struct {
	Ops []operation
	ContainerMessage string
}

type wsHandler struct {
	Name string
	Path string
	ServerInterface entityInterface
	ClientInterface entityInterface
}

type httpHandler struct {
	Operation 	string
	Path		string
	Method		string
}
type genContext struct {
	PackageName string
	HTTPHandlers []httpHandler
	WSHandlers []wsHandler
}

func (ctx *genContext)addHTTP(path string, oaOp *oasmodel.Operation, method string) {
	if oaOp != nil {
		h := httpHandler{strings.Title(oaOp.OperationID),path,method}
		ctx.HTTPHandlers = append(ctx.HTTPHandlers, h)
	}
}
func makeGenCtx(packageName string, oa *oasmodel.OpenAPI) genContext {
	var ctx genContext
	
	ctx.PackageName = packageName
	// Path Analyzing
	for k,v := range oa.Paths {
		ctx.addHTTP(k, v.Get, "GET")
		ctx.addHTTP(k, v.Put, "PUT")
		ctx.addHTTP(k, v.Post, "POST")
		ctx.addHTTP(k, v.Delete, "DELETE")
		ctx.addHTTP(k, v.Options, "OPTIONS")
		ctx.addHTTP(k, v.Head, "HEAD")
		ctx.addHTTP(k, v.Patch, "PATCH")
		ctx.addHTTP(k, v.Trace, "TRACE")
	}

	// x-ws-rpc Analyzing
	for k,v := range oa.XWsRPC {
		var wsh wsHandler
		wsh.Name = k
		wsh.Path = v.Server.OpenPath

		// fill ServerInterface
		var schServer oasmodel.Schema
		wsh.ServerInterface.ContainerMessage = v.Server.Name + "Itf"

		for i := range v.Server.Interface {
			name :=  v.Server.Interface[i].Name
			pbname := name // TODO read reference in schema and retrieve message name
			wsh.ServerInterface.Ops = append(wsh.ServerInterface.Ops, operation{name,pbname})
			schServer.OneOf = append(schServer.OneOf, v.Server.Interface[i].Schema) 
		}
		// Add new Component to Model wich is a protobuf message
		oa.Components.Schemas[wsh.ServerInterface.ContainerMessage] = &oasmodel.SchemaOrRef{nil, &schServer}

		// fill Client interface
		var schClient oasmodel.Schema
		wsh.ClientInterface.ContainerMessage = v.Client.Name + "Itf"
		for i := range v.Client.Interface {
			name :=  v.Client.Interface[i].Name
			pbname := name // TODO read reference in schema and retrieve message name
			wsh.ClientInterface.Ops = append(wsh.ClientInterface.Ops, operation{name,pbname})
			schClient.OneOf = append(schClient.OneOf, v.Client.Interface[i].Schema) 
		}
		oa.Components.Schemas[wsh.ClientInterface.ContainerMessage] = &oasmodel.SchemaOrRef{nil, &schClient}

		// Add WebSocket Handler
		ctx.WSHandlers = append(ctx.WSHandlers, wsh)
		fmt.Printf("%s\n",k)
	}

	return ctx
}

// genFile generate a file based on template & genCTX.
func genFile(outputFile string, tmpl string, ctx genContext, packageName string) error {
	fmt.Printf("Generating %s code...\n", outputFile)
	tmplData := resources.Files[tmpl]
	if tmplData == nil {
		return fmt.Errorf("Can't find template %s", tmpl)
	}
	generator := template.Must(template.New("").Parse(string(tmplData)))

	goOut, err := os.Create(packageName + "/" + outputFile)
	if err != nil {
		fmt.Printf("error opening %s : %v", outputFile, err)
		return err
	}
	defer goOut.Close()

	stream := &bytes.Buffer{}
	err = generator.Execute(stream, ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing template", err)
		return err
	}

	data, _ := format.Source(stream.Bytes())
	goOut.Write(data)

	err = generator.Execute(goOut, ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", err)
		return err
	}
	return nil
}

func main() {
	resources.Init()
	file := flag.String("f", "test.yaml", "yaml file to parse")
	output := flag.String("o", "output", "output file")
	packageName := flag.String("p", "gendefault", "package name")
	genSqueletons := flag.Bool("squel", false, "generate squeletons")
	genClient := flag.Bool("client", false, "generate client")
	flag.Parse()
	protofilename := *packageName + "/" + *output + ".proto"
	
	// Load Spec file
	oa := oasmodel.OpenAPI{}
	err := oa.Load(*file)
	if err != nil {
		flag.Usage()
		return
	} 

	// Create Generation Context
	ctx := makeGenCtx(*packageName, &oa)
	fmt.Printf("ctx %v\n", ctx)

	// Make package directory
	os.MkdirAll(*packageName, 0750)

	// Protobuf generation
	protoOut, err := os.Create(protofilename)
	if err != nil {
		fmt.Printf("error opening %s.proto : %v", *output, err)
		return
	}
	defer protoOut.Close()
	protobuf.Components2Proto(&oa, protoOut, *packageName)

	// protobuf compilation (generate go-code)
	cmd := exec.Command("protoc", "--go_out=.", protofilename)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running protoc command %v\n", err)
		os.Exit(1)
	}

	// Server go file
	err = genFile(*packageName + "-server.go", "gowapi-server.go.template", ctx, *packageName)
	if err != nil  {
		fmt.Fprintf(os.Stderr, "error : can't generate server code (%v)", err)
	}

	if *genSqueletons {
		err = genFile(*packageName + "-impl.go", "gowapi-server-squel.go.template", ctx, *packageName)
		if err != nil  {
			fmt.Fprintf(os.Stderr, "error : can't generate server sequeleton code (%v)", err)
		}
	}

	if *genClient {
		err = genFile(*packageName + "-client.go", "gowapi-client.go.template", ctx, *packageName)
		if err != nil  {
			fmt.Fprintf(os.Stderr, "error : can't generate client code (%v)", err)
		}
	}
}
