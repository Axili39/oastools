package oatool

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"

	"gopkg.in/yaml.v3"
)

func y2jConvert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = y2jConvert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = y2jConvert(v)
		}

	default:
		fmt.Printf("unmatch type %v\n", reflect.TypeOf(i))
	}
	return i
}

// yaml2Json : credits : stackoverflow ;)
func yaml2Json(buf []byte) []byte {
	var body interface{}
	if err := yaml.Unmarshal(buf, &body); err != nil {
		panic(err)
	}

	body = y2jConvert(body)

	if b, err := json.Marshal(body); err != nil {
		panic(err)
	} else {
		fmt.Printf("Output: %s\n", b)
		return b
	}
}
func j2yConvert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[string]interface{}:
		m2 := map[interface{}]interface{}{}
		for k, v := range x {
			m2[j2yConvert(k)] = j2yConvert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = j2yConvert(v)
		}
	case string:
		return x
	default:
		fmt.Printf("unmatch type %v\n", reflect.TypeOf(i))
	}
	return i
}
func json2Yaml(buf []byte) []byte {
	var body interface{}
	if err := json.Unmarshal(buf, &body); err != nil {
		panic(err)
	}

	body = j2yConvert(body)

	if b, err := yaml.Marshal(body); err != nil {
		panic(err)
	} else {
		return b
	}
}
func data2Object(obj interface{}, data []byte, intype string) error {
	var err error
	switch intype {
	case "json":
		err = json.Unmarshal(data, obj)
	case "yaml":
		jsonFile := yaml2Json(data)
		err = json.Unmarshal(jsonFile, obj)
	case "bin":
		err = proto.Unmarshal(data, obj.(proto.Message))
	}
	return err
}
func objet2Data(obj interface{}, outtype string) ([]byte, error) {
	var out []byte
	var err error
	switch outtype {
	case "json":
		out, err = json.Marshal(obj)
	case "yaml":
		var jbuf []byte
		jbuf, err = json.Marshal(obj)
		out = json2Yaml(jbuf)
	case "bin":
		out, err = proto.Marshal(obj.(proto.Message))

	}
	return out, err
}

func MainOAFileTool(getObj func(string) interface{}) {
	var file = flag.String("f", "", "input file .json/.yaml/.bin")
	var format = flag.String("o", "bin", "json|yaml|bin")
	var root = flag.String("r", "", "")
	var generate = flag.Bool("g", false, "generate empty file")
	var out []byte
	var err error
	flag.Parse()
	obj := getObj(*root)
	if obj == nil {
		fmt.Fprintf(os.Stderr, "error %v")
		return
	}
	if !*generate {
		var data []byte
		sl := strings.Split(*file, ".")
		intype := sl[len(sl)-1]

		data, err = ioutil.ReadFile(*file)
		if err != nil {
			log.Printf("error opening %s  #%v ", *file, err)
			return
		}

		err = data2Object(obj, data, intype)
		if err != nil {
			log.Printf("error Unmarshalling data %s  #%v ", *file, err)
			return
		}
	}
	out, err = objet2Data(obj, *format)

	if err != nil {
		log.Fatalf("error Marshalling: %v", err)
	}

	os.Stdout.Write(out)
	ioutil.WriteFile("output."+*format, out, 0666)
}
