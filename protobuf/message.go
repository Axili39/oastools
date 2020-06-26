package protobuf

import (
	"fmt"
	"io"

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

func createMessage(name string, schema *oasmodel.Schema, parent *Message) (ProtoType, error) {
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
