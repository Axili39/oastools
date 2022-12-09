package protobuf

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Axili39/oastools/oasmodel"
)

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
	repeated bool
	//	Options  []Option
	comment string
}

// Declare : Message Member declaration
func (t *MessageMembers) Declare(w io.Writer, indent string) {
	fmt.Fprintf(w, "%s", indent)
	// repeated
	if t.repeated {
		fmt.Fprintf(w, "repeated ")
	}
	// field decl
	fmt.Fprintf(w, "%s %s = %d;", normalizeName(t.typedecl.Name()), normalizeName(t.name), t.number)
	// TODO : options
	fmt.Fprintf(w, " /* %s */", t.comment)
	fmt.Fprintf(w, "\n")
}

// Name : Member Name
func (t *MessageMembers) Name() string {
	return t.name
}

// Message structure
type Message struct {
	name    string
	nested  []ProtoType      // Nested definitions
	body    []MessageMembers // Message Fields
	comment string
}

// Declare : ProtoType interface realization
func (t *Message) Declare(w io.Writer, indent string) {
	fmt.Fprintf(w, "%s/* Type : %s */\n", indent, t.comment)
	fmt.Fprintf(w, "%smessage %s {\n", indent, normalizeName(t.name))
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

// Name :  ProtoType interface realization
func (t *Message) Name() string {
	return t.name
}

func isRepeated(schema *oasmodel.SchemaOrRef) bool {
	if schema.Schema() == nil {
		return false
	}
	return schema.Schema().Type == "array"
}

func createMessage(name string, schema *oasmodel.Schema, parent *Message, genOpts GenerationOptions) (ProtoType, error) {
	var err error

	node := Message{name, nil, nil, schema.Description}
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
		prop := schema.Properties[m]
		if prop == nil {
			fmt.Fprintln(os.Stderr, "bad property name : ", m)
			os.Exit(1)
		}
		f := MessageMembers{nil, m, num, isRepeated(prop), prop.Description()}

		f.typedecl, err = CreateType(name+"_"+m, prop, &node, genOpts)

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

func createMessageArray(name string, schema *oasmodel.Schema, genOpts GenerationOptions) (ProtoType, error) {
	var err error

	node := Message{name + "Array", nil, nil, schema.Description}

	f := MessageMembers{nil, "Items", 1, true, schema.Items.Description()}
	f.typedecl, err = CreateType(name, schema.Items, &node, genOpts)
	if err != nil {
		return nil, err
	}
	node.body = append(node.body, f)

	return &node, nil
}

// Array : array of Prototype
type Oneof struct {
	name    string
	members []MessageMembers
}

// Declare : ProtoType interface realization
func (t *Oneof) Declare(w io.Writer, indent string) {
	fmt.Fprintf(w, "%smessage %s {\n", indent, normalizeName(t.name))
	fmt.Fprintf(w, "%s\toneof select {\n", indent)
	// body
	for m := range t.members {
		t.members[m].Declare(w, indent+"\t\t")
	}
	fmt.Fprintf(w, "\t%s}\n%s}\n", indent, indent)
}

// Name :  ProtoType interface realization
func (t *Oneof) Name() string {
	return t.name
}

func createOneOf(name string, oneof []*oasmodel.SchemaOrRef, parent *Message, genOpts GenerationOptions) (ProtoType, error) {
	node := Oneof{name, nil}
	num := 0
	for _, prop := range oneof {
		num++
		t, err := CreateType("YYY", prop, parent, genOpts)
		if err != nil {
			return nil, err
		}
		// items in array don't have name, we use typename as member name, we must clean this name frome namespace prefix
		fieldname := t.Name()
		index := strings.LastIndex(fieldname, ".")
		if index >= 0 {
			fieldname = fieldname[index+1:]
		}
		f := MessageMembers{t, fieldname + "Value", num, isRepeated(prop), prop.Description()}
		node.members = append(node.members, f)
	}
	return &node, nil
}

func createAllOf(name string, allOf []*oasmodel.SchemaOrRef, parent *Message, genOpts GenerationOptions) (ProtoType, error) {
	node := Message{name, nil, nil, ""}
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
			prop := current.Properties[m]
			f := MessageMembers{nil, m, num, isRepeated(prop), prop.Description()}
			t, err := CreateType(name+"_"+m, prop, &node, genOpts)
			if err != nil {
				return nil, err
			}
			f.typedecl = t
			node.body = append(node.body, f)
		}
	}
	return &node, nil
}
