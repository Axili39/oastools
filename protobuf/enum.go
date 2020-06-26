package protobuf

import (
	"fmt"
	"io"

	"github.com/Axili39/oastools/oasmodel"
)

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

func createEnum(name string, schema *oasmodel.Schema, parent *Message) (ProtoType, error) {
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
	return nil, fmt.Errorf("Enum must be string and have non empte Enum Array")
}
