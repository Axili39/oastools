package asciitree

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Axili39/oastools/oasmodel"
)

//ProtoType Field Type protocol buffer interface
type ProtoType interface {
	Tree(w io.Writer, name string, desc string, indent string, flag int)
}

const (
	AscTypeArray  = "[ ] "
	AscTypeObject = "{ } "
	AscTypeMap    = "M<s>"
	AscTypeInt    = "int "
	AscTypeBool   = "bool"
	AscTypeString = "str "
	AscTypeEnum   = "enum"
)
const (
	FlagFirst int = iota
	FlagMiddle
	FlagLast
)

func DrawLine(w io.Writer, symbol string, name string, desc string, indent string, flag int) string {
	if flag == FlagFirst {
		fmt.Fprintf(w, "%-40s\t\t# %-40s\n", fmt.Sprintf("%s%s ── %s", indent, symbol, name), desc)
		return indent
	}
	if flag == FlagLast {
		fmt.Fprintf(w, "%-40s\t\t# %-40s\n", fmt.Sprintf("%s└── %s ── %s", indent, symbol, name), desc)
		return indent + "    "
	}

	fmt.Fprintf(w, "%-40s\t\t# %-40s\n", fmt.Sprintf("%s├── %s ── %s", indent, symbol, name), desc)
	return indent + "│   "

}

//TypeRef simple type or reference (by-name)
type TypeName struct {
	name string
}

//Tree : ProtoType interface realization
func (t *TypeName) Tree(w io.Writer, name string, desc string, indent string, flag int) {
	DrawLine(w, t.name, name, desc, indent, flag)
}

//Enum simple type or reference (by-name)
type Enum struct {
	values []string
}

func (t *Enum) Tree(w io.Writer, name string, desc string, indent string, flag int) {
	DrawLine(w, AscTypeEnum, name, desc, indent, flag)
}

//Map object, used to represents AdditionalProperties
type Map struct {
	key   string
	value ProtoType
}

//Tree : ProtoType interface realization
func (t *Map) Tree(w io.Writer, name string, desc string, indent string, flag int) {
	newindent := DrawLine(w, AscTypeMap, name, desc, indent, flag)
	t.value.Tree(w, "", desc, newindent, FlagLast)
}

// Array : array of Prototype
type Array struct {
	typedecl ProtoType
}

//Tree : ProtoType interface realization
func (t *Array) Tree(w io.Writer, name string, desc string, indent string, flag int) {
	newindent := DrawLine(w, AscTypeArray, name, desc, indent, flag)
	t.typedecl.Tree(w, "[i]", "", newindent, FlagLast)
}

// MESSAGE

//MessageMembers Message Field definition
type ObjectMembers struct {
	//repeated bool
	typedecl ProtoType
	name     string
	desc     string
}

//Tree : Message Member declaration
func (t *ObjectMembers) Tree(w io.Writer, indent string, flag int) {
	t.typedecl.Tree(w, t.name, t.desc, indent, flag)
}

//Object structure
type Object struct {
	body []ObjectMembers // Message Fields
}

//Tree : ProtoType interface realization
func (t *Object) Tree(w io.Writer, name string, desc string, indent string, flag int) {
	newindent := DrawLine(w, AscTypeObject, name, desc, indent, flag)
	sort.Slice(t.body, func(i, j int) bool {
		return t.body[i].name < t.body[j].name
	})
	for m := range t.body {
		flag = FlagMiddle
		if m == len(t.body)-1 {
			flag = FlagLast
		}
		t.body[m].Tree(w, newindent, flag)
	}
}

//CreateType : convert OAS Schema to internal ProtoType
func CreateType(schema *oasmodel.Schema) ProtoType {
	if schema.AllOf != nil {
		node := Object{nil}
		// parse all allOf members
		for i := range schema.AllOf {
			current := schema.AllOf[i].Schema()
			for m := range current.Properties {
				prop := current.Properties[m].Schema()
				f := ObjectMembers{CreateType(prop), m, prop.Description}
				node.body = append(node.body, f)
			}
		}
		return &node
	}

	// Case AdditionnalProperties
	if schema.AdditionalProperties != nil {
		// MUST be type object
		if schema.Type != "object" {
			fmt.Fprintf(os.Stderr, "Schema with Additional Properties MUST be an object\n")
		}
		objType := CreateType(schema.AdditionalProperties.Schema.Schema())
		node := Map{"string", objType}
		return &node
	}

	if schema.Type == "object" {
		// otherwise
		node := Object{nil}
		for m := range schema.Properties {
			prop := schema.Properties[m].Schema()
			f := ObjectMembers{CreateType(prop), m, prop.Description}
			node.body = append(node.body, f)
		}
		return &node
	}

	if schema.Type == "array" {
		t := CreateType(schema.Items.Schema())
		return &Array{t}
	}

	if schema.Type == "boolean" {
		return &TypeName{AscTypeBool}
	}

	if schema.Type == "integer" {
		return &TypeName{AscTypeInt}
	}

	// Enums
	if schema.Type == "string" && len(schema.Enum) > 0 {
		node := Enum{nil}
		for i := range schema.Enum {
			node.values = append(node.values, schema.Enum[i])
		}
		return &node
	}
	return &TypeName{AscTypeString}
}

//Components2Proto : generate proto file from Parsed OpenAPI definition
func Components2AscTree(oa *oasmodel.OpenAPI, f io.Writer, root string) {
	oa.ResolveRefs()
	// create first level Nodes
	for k, v := range oa.Components.Schemas {
		if k == root {
			node := CreateType(v.Schema())
			node.Tree(f, k, v.Schema().Description, "", FlagFirst)
		}
	}
}
