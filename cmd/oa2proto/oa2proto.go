package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/Axili39/oastools/oasmodel"
	"github.com/Axili39/oastools/protobuf"
)

// Multiples file in command lines
type stringList []string

func (i *stringList) String() string {
	return ""
}

func (i *stringList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func compileProto(protofilename string, directory string) {
	protoPath := filepath.Dir(protofilename)
	cmd := exec.Command("protoc", "--go_out="+directory, "--proto_path="+protoPath, protofilename)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %v running protoc command: %v\n", err, cmd)
		os.Exit(1)
	}
}

func main() {
	build := flag.String("build", "", "build with protoc")
	file := flag.String("f", "", "yaml file to parse")
	out := flag.String("o", "", "output file")
	AddEnumPrefix := flag.Bool("add-enum-prefix", false, "Auto add prefix on Enums")
	packageName := flag.String("p", "", "package name eg: foo.bar")
	showversion := flag.Bool("v", false, "show version")
	var options stringList
	flag.Var(&options, "option", "add directive option in .proto file (multi)")
	var filteredNodes stringList
	flag.Var(&filteredNodes, "node", "select component (multi)")
	var packageNameMap stringList
	flag.Var(&packageNameMap, "rename-package", "rename package imports")
	flag.Parse()

	genOpts := protobuf.GenerationOptions{AddEnumPrefix: *AddEnumPrefix, Imports: make(map[string]bool), PackageNames: map[string]string{}}

	if *showversion {
		if info, available := debug.ReadBuildInfo(); available {
			fmt.Println(info.Main.Version)
		} else {
			fmt.Println("unknown")
		}
		os.Exit(0)
	}

	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
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

	// Chargement du fichier
	var err error
	oa := oasmodel.OpenAPI{}
	if *file == "" {
		err = oa.Read(os.Stdin)
	} else {
		err = oa.Load(*file)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s : %v", *file, err)
		os.Exit(1)
	}

	// load package Map
	for _, p := range packageNameMap {
		pair := strings.Split(p, ":")
		if len(pair) != 2 {
			fmt.Fprintf(os.Stderr, "error bad package rename :", p)
		}
		genOpts.PackageNames[pair[0]] = pair[1]
	}

	err = protobuf.Components2Proto(&oa, output, *packageName, genOpts, filteredNodes, options...)
	if err != nil {
		fmt.Println("error")
		fmt.Fprintf(os.Stderr, "error parsing %s : %v", *file, err)
		os.Exit(1)
	}

	if *build != "" && *out != "" {
		compileProto(*out, *build)
	}
}
