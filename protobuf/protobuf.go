package protobuf

import (
	"fmt"
	"io"
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
	Declare(w io.Writer, indent string)
	Name() string
	Repeated() bool
}

//TypeName simple type or reference (by-name)
type TypeName struct {
	name string
}

//Declare : ProtoType interface realization
func (t *TypeName) Declare(w io.Writer, indent string) {
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
func (t *Enum) Declare(w io.Writer, indent string) {
	fmt.Fprintf(w, "%senum %s {\n", indent, t.name)
	values := 0
	for i := range t.values {
		fmt.Fprintf(w, "%s\t%s = %d;\n", indent, t.values[i], values)
		values++
	}
	fmt.Fprintf(w, "%s}\n", indent)
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
func (t *Map) Declare(w io.Writer, indent string) {
	fmt.Fprintf(w, "map<%s, %s>", t.key, t.value.Name())
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

// MESSAGE
/*
//Option for message fields
type Option struct {
	name  string
	value string
}
*/
//MessageMembers Message Field definition
type MessageMembers struct {
	//repeated bool
	typedecl ProtoType
	name     string
	number   int
	//	Options  []Option
}

//Declare : Message Member declaration
func (t *MessageMembers) Declare(w io.Writer, indent string) {
	fmt.Fprintf(w, "%s", indent)
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
func (t *Message) Declare(w io.Writer, indent string) {
	fmt.Fprintf(w, "%smessage %s {\n", indent, t.name)
	// nested
	for n := range t.nested {
		t.nested[n].Declare(w, indent+"\t")
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

// Basic Types
// OAS specify format for type :
// Type		Format	Description
// number	–		Any numbers.
// number	float	Floating-point numbers.
// number	double	Floating-point numbers with double precision.
// integer	–		Integer numbers.
// integer	int32	Signed 32-bit integers (commonly used integer type).
// integer	int64	Signed 64-bit integers (long type).

// Protobuf has much more specs:
// type = "double" | "float" | "int32" | "int64" | "uint32" | "uint64"
//      | "sint32" | "sint64" | "fixed32" | "fixed64" | "sfixed32" | "sfixed64"
//      | "bool" | "string" | "bytes" | messageType | enumType
func createTypename(typename, format string) string {
	switch typename {
	case "number":
		switch format {
		case "float":
			return "float"
		case "double":
			return "double"
		default:
			return "double"
		}
	case "integer":
		switch format {
		case "int32":
			return "int32"
		case "int64":
			return "int32"
		case "uint32":
			return "uint32"
		case "uint64":
			return "uint64"
		default:
			return "int32"
		}
	case "boolean":
		return "bool"
	}

	return typename
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

func createObject(name string, schema *oasmodel.Schema, parent *Message) (ProtoType, error) {
	var err error

	node := Message{name, nil, nil}
	num := 0
	// sorting Properties Name
	var keys []string
	if len(schema.XPropertiesOrder) > 0 {
		keys = schema.XPropertiesOrder
	} else {
		keys = keysorder(schema.Properties)
	}

	// Add each Properties as message Member
	for _, m := range keys {
		num++
		f := MessageMembers{nil, m, num}
		prop := schema.Properties[m]
		f.typedecl, err = CreateType(name+"_"+m, prop, &node)
		if err != nil {
			return nil, err
		}
		node.body = append(node.body, f)
	}
	// if has parent insert as nested message
	if parent != nil {
		parent.nested = append(parent.nested, &node)
	}
	return &node, nil
}

func createAdditionalProperties(name string, schema *oasmodel.Schema, parent *Message) (ProtoType, error) {
	if parent == nil {
		return nil, fmt.Errorf("Warning: Additional Properties can only be a message member (ignore)")
	}
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
			return &TypeName{schemaOrRef.Ref.RefName}, nil
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
		return createObject(name, schema, parent)
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
		node := Enum{name, nil}
		for i := range schema.Enum {
			node.values = append(node.values, schema.Enum[i])
		}
		if parent != nil {
			parent.nested = append(parent.nested, &node)
		}
		return &node, nil
	}

	// bytes
	if schema.Type == "string" && schema.Format == "binary" {
		return &TypeName{"bytes"}, nil
	}

	return &TypeName{createTypename(schema.Type, schema.Format)}, nil
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
