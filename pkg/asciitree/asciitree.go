package asciitree

import (
	"fmt"
	"os"
	"io"
	"sort"

	"../../pkg/oasmodel"
)


//ProtoType Field Type protocol buffer interface
type ProtoType interface {
	Tree(w io.Writer, name string, desc string, indent string, flag int)
}
const (
	S_ARRAY      	= "[ ] "
	S_OBJECT		= "{ } "
	S_MAP			= "M<s>"
	S_INT			= "int "
	S_BOOL			= "bool"
	S_STRING		= "str "
	S_ENUM			= "enum"
)
const (
	FLAG_FIRST int = iota
	FLAG_MIDDLE
	FLAG_LAST
)
func DrawLine(w io.Writer, symbol string, name string, desc string, indent string, flag int) string {
	if flag == FLAG_FIRST {
		fmt.Fprintf(w, "%-40s\t\t# %-40s\n", fmt.Sprintf("%s%s ── %s", indent, symbol, name), desc);
		return indent
	}
	if flag == FLAG_LAST {
		fmt.Fprintf(w, "%-40s\t\t# %-40s\n", fmt.Sprintf("%s└── %s ── %s", indent, symbol, name), desc);
		return indent + "    "
	} 

	fmt.Fprintf(w, "%-40s\t\t# %-40s\n", fmt.Sprintf("%s├── %s ── %s", indent, symbol, name), desc);
	return indent + "│   "
	
}


//TypeRef simple type or reference (by-name)
type TypeName struct {
	name    string
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
	DrawLine(w, S_ENUM, name, desc, indent, flag)
}

//Map object, used to represents AdditionalProperties
type Map struct {
	key   string
	value ProtoType
}

//Tree : ProtoType interface realization
func (t *Map) Tree(w io.Writer, name string, desc string, indent string, flag int) {
	newindent := DrawLine(w, S_MAP, name, desc,  indent, flag)
	t.value.Tree(w, "", desc, newindent, FLAG_LAST)
}
// Array : array of Prototype
type Array struct {
	typedecl ProtoType
}

//Tree : ProtoType interface realization
func (t *Array) Tree(w io.Writer, name string, desc string, indent string, flag int) {
	newindent := DrawLine(w, S_ARRAY, name, desc, indent, flag)
	t.typedecl.Tree(w, "[i]", "", newindent,  FLAG_LAST)
}

// MESSAGE

//MessageMembers Message Field definition
type ObjectMembers struct {
	//repeated bool
	typedecl ProtoType
	name     string
	desc	 string
}

//Tree : Message Member declaration
func (t *ObjectMembers) Tree(w io.Writer, indent string, flag int) {
	t.typedecl.Tree(w, t.name, t.desc, indent, flag)
}

//Object structure
type Object struct {
	body   []ObjectMembers // Message Fields
}

//Tree : ProtoType interface realization
func (t *Object) Tree(w io.Writer, name string, desc string, indent string, flag int) {
	newindent := DrawLine(w, S_OBJECT, name, desc, indent, flag)
	sort.Slice(t.body, func(i, j int) bool {
		return t.body[i].name < t.body[j].name
	})
	for m := range t.body {
		flag = FLAG_MIDDLE
		if m == len(t.body)-1 {
			flag = FLAG_LAST
		}
		t.body[m].Tree(w, newindent,  flag)
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
		node := Array{t}
		return &node
	}

	if schema.Type == "boolean" {
		node := TypeName{S_BOOL}
		return &node
	}

	if schema.Type == "integer" {
		node := TypeName{S_INT}
		return &node
	}

	// Enums
	if schema.Type == "string" && len(schema.Enum) > 0 {
		node := Enum{nil}
		for i := range schema.Enum {
			node.values = append(node.values, schema.Enum[i])
		}
		return &node
	}
	node := TypeName{S_STRING}
	return &node
}

//Components2Proto : generate proto file from Parsed OpenAPI definition
func Components2AscTree(oa *oasmodel.OpenAPI, f io.Writer, root string) {
	oa.ResolveRefs()
	// create first level Nodes
	for k,v := range oa.Components.Schemas {
		if k == root {
			node := CreateType(v.Schema())
			node.Tree(f, k, v.Schema().Description, "", FLAG_FIRST)
		}
	}
}
