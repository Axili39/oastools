package protobuf

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Axili39/oastools/oasmodel"
)

// ProtoType Field Type protocol buffer interface
type ProtoType interface {
	Declare(w io.Writer, indent string)
	Name() string
}

// Map object, used to represents AdditionalProperties
type Map struct {
	name  string
	key   string
	value ProtoType
}

// Generation Options
type GenerationOptions struct {
	AddEnumPrefix bool
	PackageNames  map[string]string
	Imports       map[string]bool
}

// Declare : ProtoType interface realization
func (t *Map) Declare(w io.Writer, indent string) {
}

// Name :  ProtoType interface realization
func (t *Map) Name() string {
	return "map<" + t.key + ", " + t.value.Name() + ">"
}

func normalizeName(name string) string {
	return strings.Replace(name, "-", "_", -1)
}

func createAdditionalProperties(name string, schema *oasmodel.Schema, parent *Message, genOpts GenerationOptions) (ProtoType, error) {
	if schema.Type != "object" {
		return nil, fmt.Errorf("Schema %s with Additional Properties must be an object", name)
	}

	if schema.AdditionalProperties.Schema == nil {
		return nil, fmt.Errorf("Schema %s with Additional Properties : Unsupported AdditionalProperties with boolean value for protobuf generation, Schema MUST be provided", name)
	}

	objType, err := CreateType(name+"Elem", schema.AdditionalProperties.Schema, parent, genOpts)
	if err != nil {
		return nil, err
	}
	return &Map{name, "string", objType}, nil

}

// CreateType : convert OAS Schema to internal ProtoType
func CreateType(name string, schemaOrRef *oasmodel.SchemaOrRef, parent *Message, genOpts GenerationOptions) (ProtoType, error) {
	fmt.Println(name, schemaOrRef)
	schema := schemaOrRef.Schema()
	// In case of Ref, we need to get the corresponding type name
	if schemaOrRef.Ref != nil {
		if schemaOrRef.Ref.External != "" {
			fmt.Println("external", schemaOrRef.Ref.External)
			packageName := schemaOrRef.Ref.External

			// rename package
			if v, ok := genOpts.PackageNames[packageName]; ok {
				// change package name
				packageName = v
			}

			// store package in imports
			genOpts.Imports[schemaOrRef.Ref.External] = true

			return createTypename(packageName+"."+schemaOrRef.Ref.RefName, "")
		}
		if schema == nil {
			return nil, fmt.Errorf("bad ref")
		}
		if schema.AllOf != nil || schema.Type == "object" && schema.AdditionalProperties == nil || (schema.Type == "string" && len(schema.Enum) > 0) {
			// in case of Ref, reference type name only for messages :
			return createTypename(schemaOrRef.Ref.RefName, "")
		}
	}
	// case Oneof
	if schema.OneOf != nil {
		return createOneOf(name, schema.OneOf, parent, genOpts)
	}
	// case AllOf
	if schema.AllOf != nil {
		return createAllOf(name, schema.AllOf, parent, genOpts)
	}
	// Case AdditionalProperties
	if schema.AdditionalProperties != nil {
		return createAdditionalProperties(name, schema, parent, genOpts)
	}
	// case Object
	if schema.Type == "object" {
		return createMessage(name, schema, parent, genOpts)
	}
	// case Array
	if schema.Type == "array" {
		if parent == nil {
			return createMessageArray(name, schema, genOpts)
		}
		return CreateType(name, schema.Items, parent, genOpts)
	}
	// Enums
	if schema.Type == "string" && len(schema.Enum) > 0 {
		return createEnum(name, schema, parent, genOpts)
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

// Components2Proto : generate proto file from Parsed OpenAPI definition
func Components2Proto(oa *oasmodel.OpenAPI, f io.Writer, packageName string, genOpts GenerationOptions, filternodes []string, options ...string) error {
	var items []string
	if filternodes == nil {
		oa.ResolveRefs()
		items = keysorder(oa.Components.Schemas)
	} else {
		items = keysorder(oa.ResolveRefsWithFilter(filternodes))
	}
	nodeList := make([]ProtoType, 0, 10)
	// create first level Nodes
	for _, k := range items {
		v := oa.Components.Schemas[k]
		node, err := CreateType(k, v, nil, genOpts)
		if err != nil {
			fmt.Println("error : ", err)
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
	for packageFile, _ := range genOpts.Imports {
		fmt.Fprintf(f, "import \"%s.proto\";\n", packageFile)
	}
	for n := range nodeList {
		nodeList[n].Declare(f, "")
	}
	return nil
}
