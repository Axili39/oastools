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
}

//Map object, used to represents AdditionalProperties
type Map struct {
	name  string
	key   string
	value ProtoType
}

//Declare : ProtoType interface realization
func (t *Map) Declare(w io.Writer, indent string) {
}

//Name :  ProtoType interface realization
func (t *Map) Name() string {
	return "map<" + t.key + ", " + t.value.Name() + ">"
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
		return CreateType(name, schema.Items, parent)
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
	}
	for _, opt := range options {
		fmt.Fprintln(f, "option ", opt, ";")
	}
	for n := range nodeList {
		nodeList[n].Declare(f, "")
	}
	return nil
}
