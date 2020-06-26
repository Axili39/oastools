package protobuf

import (
	"fmt"
	"io"
	"sort"

	"github.com/Axili39/oastools/oasmodel"
)

//ProtoType Field Type protocol buffer interface
type ProtoType interface {
	Declare(w io.Writer, indent string)
	Name() string
	Repeated() bool
}

//Map object, used to represents AdditionalProperties
type Map struct {
	name  string
	key   string
	value ProtoType
}

//Declare : ProtoType interface realization
func (t *Map) Declare(w io.Writer, indent string) {
	//log.Println("called MAP.Declare()")
	//fmt.Fprintf(w, "map<%s, %s>", t.key, t.value.Name())
}

//Name :  ProtoType interface realization
func (t *Map) Name() string {
	return "map<" + t.key + ", " + t.value.Name() + ">"
}

//Repeated :  ProtoType interface realization
func (t *Map) Repeated() bool {
	return false
}

// Array : array of Prototype
type Oneof struct {
	name    string
	members []MessageMembers
}

//Declare : ProtoType interface realization
func (t *Oneof) Declare(w io.Writer, indent string) {
	fmt.Fprintf(w, "%smessage %s {\n", indent, t.name)
	fmt.Fprintf(w, "%s\toneof select {\n", indent)
	// body
	for m := range t.members {
		t.members[m].Declare(w, indent+"\t\t")
	}
	fmt.Fprintf(w, "\t%s}\n%s}\n", indent, indent)
}

//Name :  ProtoType interface realization
func (t *Oneof) Name() string {
	return t.name
}

//Repeated :  ProtoType interface realization
func (t *Oneof) Repeated() bool {
	return false
}

// Array : array of Prototype
type Array struct {
	typedecl ProtoType
}

//Declare : ProtoType interface realization
func (t *Array) Declare(w io.Writer, indent string) {
	// does't exist in protobuf
	t.typedecl.Declare(w, indent)
}

//Name :  ProtoType interface realization
func (t *Array) Name() string {
	return t.typedecl.Name()
}

//Repeated :  ProtoType interface realization
func (t *Array) Repeated() bool {
	return true
}

func createOneOf(name string, oneof []*oasmodel.SchemaOrRef, parent *Message) (ProtoType, error) {
	node := Oneof{name, nil}
	num := 0
	for _, prop := range oneof {
		num++
		t, err := CreateType("YYY", prop, parent)
		if err != nil {
			return nil, err
		}
		f := MessageMembers{t, t.Name() + "Value", num}
		node.members = append(node.members, f)
	}
	return &node, nil
}

func createAllOf(name string, allOf []*oasmodel.SchemaOrRef, parent *Message) (ProtoType, error) {
	node := Message{name, nil, nil}
	num := 0

	// parse all allOf members
	for _, val := range allOf {
		current := val.Schema()
		var keys []string
		if len(current.XPropertiesOrder) > 0 {
			keys = current.XPropertiesOrder
		} else {
			keys = keysorder(current.Properties)
		}
		for _, m := range keys {
			num++
			f := MessageMembers{nil, m, num}
			prop := current.Properties[m]
			t, err := CreateType(name+"_"+m, prop, &node)
			if err != nil {
				return nil, err
			}
			f.typedecl = t
			node.body = append(node.body, f)
		}
	}
	return &node, nil
}

func createAdditionalProperties(name string, schema *oasmodel.Schema, parent *Message) (ProtoType, error) {
	if schema.Type != "object" {
		return nil, fmt.Errorf("Schema %s with Additional Properties must be an object", name)
	}

	objType, err := CreateType(name+"Elem", schema.AdditionalProperties.Schema, parent)
	if err != nil {
		return nil, err
	}
	return &Map{name, "string", objType}, nil

}

//CreateType : convert OAS Schema to internal ProtoType
func CreateType(name string, schemaOrRef *oasmodel.SchemaOrRef, parent *Message) (ProtoType, error) {
	schema := schemaOrRef.Schema()

	// In case of Ref, we need to get the corresponding type name
	if schemaOrRef.Ref != nil {
		if schema.AllOf != nil || schema.Type == "object" && schema.AdditionalProperties == nil || (schema.Type == "string" && len(schema.Enum) > 0) {
			// in case of Ref, reference type name only for messages :
			return createTypename(schemaOrRef.Ref.RefName, "")
		}
	}

	// case Oneof
	if schema.OneOf != nil {
		return createOneOf(name, schema.OneOf, parent)
	}

	// case AllOf
	if schema.AllOf != nil {
		return createAllOf(name, schema.AllOf, parent)
	}

	// Case AdditionalProperties
	if schema.AdditionalProperties != nil {
		return createAdditionalProperties(name, schema, parent)
	}

	// case Object
	if schema.Type == "object" {
		return createMessage(name, schema, parent)
	}

	// case Array
	if schema.Type == "array" {
		t, err := CreateType(name+"Elem", schema.Items, parent)
		if err != nil {
			return nil, err
		}
		return &Array{t}, nil
	}

	// Enums
	if schema.Type == "string" && len(schema.Enum) > 0 {
		return createEnum(name, schema, parent)
	}

	return createTypename(schema.Type, schema.Format)
}

func keysorder(m map[string]*oasmodel.SchemaOrRef) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

//Components2Proto : generate proto file from Parsed OpenAPI definition
func Components2Proto(oa *oasmodel.OpenAPI, f io.Writer, packageName string, options ...string) error {
	oa.ResolveRefs()
	nodeList := make([]ProtoType, 0, 10)
	// create first level Nodes
	for _, k := range keysorder(oa.Components.Schemas) {
		v := oa.Components.Schemas[k]
		node, err := CreateType(k, v, nil)
		if err != nil {
			// silentely ignore it
			continue
		}
		nodeList = append(nodeList, node)
	}

	fmt.Fprintf(f, "syntax = \"proto3\";\n")
	if packageName != "" {
		fmt.Fprintln(f, "package ", packageName, ";")
		//fmt.Fprintf(f, "option go_package = \".;%s\";\n", packageName)
	}
	for _, opt := range options {
		fmt.Fprintln(f, "option ", opt, ";")
	}
	for n := range nodeList {
		nodeList[n].Declare(f, "")
	}
	return nil
}
