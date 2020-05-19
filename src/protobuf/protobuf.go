package protobuf

import (
	"fmt"
	"os"
	"sort"
	"github.com/Axili39/oastools/oasmodel"
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
	Declare(w *os.File, indent string)
	Name() string
	Repeated() bool
}

//TypeName simple type or reference (by-name)
type TypeName struct {
	name    string
	refName string //useless
}

//Declare : ProtoType interface realization
func (t *TypeName) Declare(w *os.File, indent string) {
	// does't exist in protobuf
}

//Name :  ProtoType interface realization
func (t *TypeName) Name() string {
	return t.name
}

//Repeated ProtoType interface realization
func (t *TypeName) Repeated() bool {
	return false
}

//Enum simple type or reference (by-name)
type Enum struct {
	name   string
	values []string
}

//Declare : ProtoType interface realization
/* Example :
enum Corpus {
    UNIVERSAL = 0;
    WEB = 1;
    IMAGES = 2;
    LOCAL = 3;
    NEWS = 4;
    PRODUCTS = 5;
    VIDEO = 6;
  }
*/
func (t *Enum) Declare(w *os.File, indent string) {
	fmt.Fprintf(w, "%senum %s {\n", indent, t.name)
	values := 0
	for i := range t.values {
		fmt.Fprintf(w, "%s\t%s = %d;\n", indent, t.values[i], values)
		values++
	}
	fmt.Fprintf(w, "%s}\n",indent)
}

//Name :  ProtoType interface realization
func (t *Enum) Name() string {
	return t.name
}

//Repeated :ProtoType interface realization
func (t *Enum) Repeated() bool {
	return false
}

//Map object, used to represents AdditionalProperties
type Map struct {
	name  string
	key   string
	value ProtoType
}

//Declare : ProtoType interface realization
func (t *Map) Declare(w *os.File , indent string) {
	fmt.Fprintf(w, "map<%s, value>", t.key, t.value.Name())
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
	name string
	members []MessageMembers
}

//Declare : ProtoType interface realization
func (t *Oneof) Declare(w *os.File, indent string) {
	fmt.Fprintf(w, "%smessage %s {\n", indent, t.name)
	fmt.Fprintf(w, "%s\toneof select {\n", indent)
	// body
	for m := range t.members {
		t.members[m].Declare(w, indent+"\t\t")
	}
	fmt.Fprintf(w, "\t%s}\n%s}\n", indent,indent)
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
func (t *Array) Declare(w *os.File, indent string) {
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
func (t *MessageMembers) Declare(w *os.File, indent string) {
	fmt.Fprintf(w,"%s", indent)
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
func (t *Message) Declare(w *os.File, indent string) {
	fmt.Fprintf(w, "%smessage %s {\n", indent, t.name)
	// nested
	for n := range t.nested {
		t.nested[n].Declare(w, indent + "\t")
	}
	// body
	for m := range t.body {
		t.body[m].Declare(w, indent+"\t")
	}
	fmt.Fprintf(w, "%s}\n", indent)
}

//Name :  ProtoType interface realization
func (t *Message) Name() string {
	return t.name
}

//Repeated : ProtoType interface realization
func (t *Message) Repeated() bool {
	return false
}

//CreateType : convert OAS Schema to internal ProtoType
func CreateType(name string, schema *oasmodel.SchemaOrRef, Parent *Message) ProtoType {
	if schema.Ref != nil {
		node := TypeName{schema.Ref.RefName, ""}
		return &node
	}
	if schema.Val.OneOf != nil {
		node := Oneof{name, nil}
		num := 0
		for i := range schema.Val.OneOf {
			num++
			prop := schema.Val.OneOf[i]
			t := CreateType("YYY", prop, Parent)
			f := MessageMembers{t, t.Name() + "Value", num, nil}			
			node.members = append(node.members, f)
		}
		return &node
	}
	if schema.Val.AllOf != nil {
		node := Message{name, nil, nil}
		num := 0
		// parse all allOf members
		for i := range schema.Val.AllOf {
			current := schema.Val.AllOf[i].Schema()
			var keys[]string
			if len(current.XPropertiesOrder) > 0 {
				keys = current.XPropertiesOrder
			} else {
				keys = keysorder(current.Properties)
			}
			for i := range keys {
				m := keys[i]
				num++
				f := MessageMembers{nil, m, num, nil}
				prop := current.Properties[m]
				t := CreateType(name+"_"+m, prop, &node)
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
		objType := CreateType(name+"Elem", schema.Val.AdditionalProperties.Schema, Parent)
		node := Map{name, "string", objType}
		return &node
	}

	if schema.Val.Type == "object" {

		// otherwise
		node := Message{name, nil, nil}
		// parse elements
		num := 0
		var keys []string
		if len(schema.Val.XPropertiesOrder) > 0 {
			keys = schema.Val.XPropertiesOrder
		} else {
			keys = keysorder(schema.Val.Properties)
		}
		for i := range keys {
			m := keys[i]
			num++
			f := MessageMembers{nil, m, num, nil}
			prop := schema.Val.Properties[m]
			t := CreateType(name+"_"+m, prop, &node)
			f.typedecl = t
			node.body = append(node.body, f)
		}
		if Parent != nil {
			Parent.nested = append(Parent.nested, &node)
		}
		return &node
	}

	if schema.Val.Type == "array" {
		t := CreateType(name+"Elem", schema.Val.Items, Parent)
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

	// Enums
	if schema.Val.Type == "string" && len(schema.Val.Enum) > 0 {
		node := Enum{name, nil}
		for i := range schema.Val.Enum {
			node.values = append(node.values, schema.Val.Enum[i])
		}
		if Parent != nil {
			Parent.nested = append(Parent.nested, &node)
		}
		return &node
	}

	// bytes
	if schema.Val.Type == "string" && schema.Val.Format == "binary" {
		node := TypeName{"bytes", ""}
		return &node
	}

	node := TypeName{schema.Val.Type, ""}
	return &node
}

func keysorder(m map[string]*oasmodel.SchemaOrRef) ([]string) {
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
func Components2Proto(oa *oasmodel.OpenAPI, f *os.File, packageName string) {
	oa.ResolveRefs()
	nodeList := make([]ProtoType, 0, 10)
	// create first level Nodes
	for _,k := range keysorder(oa.Components.Schemas) {
			v := oa.Components.Schemas[k]
			node := CreateType(k, v, nil)
			nodeList = append(nodeList, node)
	}

	fmt.Fprintf(f, "syntax = \"proto3\";\n")
	if packageName != "" {
		fmt.Fprintf(f, "option go_package = \"%s\";\n", packageName)
	}
	for n := range nodeList {
		nodeList[n].Declare(f, "")
	}
}
