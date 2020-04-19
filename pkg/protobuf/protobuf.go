package protobuf

import (
	"fmt"
	"os"

	"../../pkg/oasmodel"
)

/*protobuf node
 Specs:
 	message = "message" messageName messageBody
	messageBody = "{" { field | enum | message | option | oneof | mapField |
	reserved | emptyStatement } "}"

	type = "double" | "float" | "int32" | "int64" | "uint32" | "uint64"
      | "sint32" | "sint64" | "fixed32" | "fixed64" | "sfixed32" | "sfixed64"
      | "bool" | "string" | "bytes" | messageType | enumType
	fieldNumber = intLit;

	field = [ "repeated" ] type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
	fieldOptions = fieldOption { ","  fieldOption }
	fieldOption = optionName "=" constant


	Example :

syntax = "proto3";
import public "other.proto";
option java_package = "com.example.foo";
enum EnumAllowingAlias {
  option allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2 [(custom_option) = "hello world"];
}
message outer {
  option (my_option).a = true;
  message inner {   // Level 2
    int64 ival = 1;
  }
  repeated inner inner_message = 2;
  EnumAllowingAlias enum_field =3;
  map<int32, string> my_map = 4;
}
*/

//ProtoType Field Type protocol buffer interface
type ProtoType interface {
	Declare(w *os.File)
	Name() string
	Repeated() bool
}

//TypeName simple type or reference (by-name)
type TypeName struct {
	name    string
	refName string //useless
}

//Declare : ProtoType interface realization
func (t *TypeName) Declare(w *os.File) {
	// does't exist in protobuf
}

//Name :  ProtoType interface realization
func (t *TypeName) Name() string {
	return t.name
}

//Repeated
func (t *TypeName) Repeated() bool {
	return false
}

//TypeName simple type or reference (by-name)
type Map struct {
	name  string
	key   string
	value ProtoType
}

//Declare : ProtoType interface realization
func (t *Map) Declare(w *os.File) {
	fmt.Fprintf(w, "map<%s, value>", t.key, t.value.Name())
}

//Name :  ProtoType interface realization
func (t *Map) Name() string {
	return "map<" + t.key + ", " + t.value.Name() + ">"
}

//Repeated
func (t *Map) Repeated() bool {
	return false
}

// ARRAY
type Array struct {
	typedecl ProtoType
}

//Declare : ProtoType interface realization
func (t *Array) Declare(w *os.File) {
	// does't exist in protobuf
	t.typedecl.Declare(w)
}

//Name :  ProtoType interface realization
func (t *Array) Name() string {
	return t.typedecl.Name()
}

//Repeated
func (t *Array) Repeated() bool {
	return true
}

// MESSAGE

//Option for message fields
type Option struct {
	name  string
	value string
}

//MessageMembers Message Field definition
type MessageMembers struct {
	//repeated bool
	typedecl ProtoType
	name     string
	number   int
	Options  []Option
}

//Declare : Message Member declaration
func (t *MessageMembers) Declare(w *os.File) {
	// repeated
	if t.typedecl.Repeated() {
		fmt.Fprintf(w, "repeated ")
	}
	// field decl
	fmt.Fprintf(w, "%s %s = %d;", t.typedecl.Name(), t.name, t.number)
	// TODO : options
	fmt.Fprintf(w, "\n")
}

//Name : Member Name
func (t *MessageMembers) Name() string {
	return t.name
}

//Message structure
type Message struct {
	name   string
	nested []ProtoType      // Nested definitions
	body   []MessageMembers // Message Fields
}

//Declare : ProtoType interface realization
func (t *Message) Declare(w *os.File) {
	fmt.Fprintf(w, "message %s {\n", t.name)
	// nested
	for n := range t.nested {
		t.nested[n].Declare(w)
	}
	// body
	for m := range t.body {
		t.body[m].Declare(w)
	}
	fmt.Fprintf(w, "}\n")
}

//Name :  ProtoType interface realization
func (t *Message) Name() string {
	return t.name
}

//Repeated
func (t *Message) Repeated() bool {
	return false
}

//refIndexElement part of map used to resolve réfences (shall be built-in in OAS package ...)
type refIndexElement struct {
	name   string
	schema *oasmodel.SchemaOrRef
}

func makeRefindex(oa *oasmodel.OpenAPI) map[string]refIndexElement {
	refindex := make(map[string]refIndexElement)
	// create index of all Schemas (only level 1)
	for k := range oa.Components.Schemas {
		refindex["#/components/schemas/"+k] = refIndexElement{k, oa.Components.Schemas[k]}
		// TODO: go down items / properties to create new Path
	}
	return refindex
}

//CreateType : convert OAS Schema to internal ProtoType
func CreateType(name string, schema *oasmodel.SchemaOrRef, refIndex map[string]refIndexElement, Parent *Message) ProtoType {
	if schema.Ref != nil {
		// search referenced object
		// TODO : make recursive method to support refs of refs
		if elem, ok := refIndex[schema.Ref.Ref]; ok {
			node := TypeName{elem.name, ""}
			return &node
		}
		return nil
	}

	if schema.Val.AllOf != nil {
		node := Message{name, nil, nil}
		num := 0
		// parse all allOf members
		for i := range schema.Val.AllOf {
			current := schema.Val.AllOf[i]
			var defVal *oasmodel.Schema
			defVal = nil
			// Verify that composer is a object
			if current.Ref != nil {
				if elem, ok := refIndex[current.Ref.Ref]; ok {
					defVal = elem.schema.Val // assume there is no refs of refs
				} else {
					fmt.Fprintf(os.Stderr, "can't find ref %s\n", current.Ref.Ref)
					return nil
				}
			} else {
				defVal = current.Val
			}
			if defVal.Type != "object" {
				fmt.Fprintf(os.Stderr, "can't support allOf without ref %s\n", current.Ref.Ref)
				return nil
			}
			// now we have an object => copy properties

			for m := range defVal.Properties {
				num++
				f := MessageMembers{nil, m, num, nil}
				prop := defVal.Properties[m]
				t := CreateType(name+"_"+m, prop, refIndex, &node)
				f.typedecl = t
				node.body = append(node.body, f)
			}
		}

		return &node
	}

	// Case AdditionnalProperties
	if schema.Val.AdditionalProperties != nil {
		// MUST be type object
		if schema.Val.Type != "object" {
			fmt.Fprintf(os.Stderr, "Schema %s with Additional Properties MUST be an object\n", name)
		}
		objType := CreateType(name+"Elem", schema.Val.AdditionalProperties.Schema, refIndex, Parent)
		node := Map{name, "string", objType}
		return &node
	}

	if schema.Val.Type == "object" {

		// otherwise
		node := Message{name, nil, nil}
		// parse elements
		num := 0

		for m := range schema.Val.Properties {
			num++
			f := MessageMembers{nil, m, num, nil}
			prop := schema.Val.Properties[m]
			t := CreateType(name+"_"+m, prop, refIndex, &node)
			f.typedecl = t
			node.body = append(node.body, f)
		}
		if Parent != nil {
			Parent.nested = append(Parent.nested, &node)
		}
		return &node
	}

	if schema.Val.Type == "array" {
		t := CreateType(name+"Elem", schema.Val.Items, refIndex, Parent)
		node := Array{t}
		return &node
	}

	if schema.Val.Type == "boolean" {
		node := TypeName{"bool", ""}
		return &node
	}

	if schema.Val.Type == "integer" {
		node := TypeName{"int32", ""}
		return &node
	}
	node := TypeName{schema.Val.Type, ""}
	return &node
}

//Components2Proto : generate proto file from Parsed OpenAPI definition
func Components2Proto(oa *oasmodel.OpenAPI, f *os.File) {
	refindex := makeRefindex(oa)
	var nodeList []ProtoType
	// create first level Nodes
	for k := range oa.Components.Schemas {
		elem := oa.Components.Schemas[k]
		node := CreateType(k, elem, refindex, nil)
		nodeList = append(nodeList, node)
	}

	fmt.Fprintf(f, "syntax = \"proto3\";\n")
	for n := range nodeList {
		nodeList[n].Declare(f)
	}
}

//Test Function (to be removed)
func Test() {
	oa := oasmodel.OpenAPI{}
	oa.Load("test.yaml")

	Components2Proto(&oa, os.Stdout)
}
