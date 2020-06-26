package protobuf

import (
	"io"
)

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
func createTypename(typename, format string) (ProtoType, error) {
	if typename == "number" {
		if format == "" {
			return &TypeName{"double"}, nil
		}
		return &TypeName{format}, nil
	}

	if typename == "integer" {
		if format == "" {
			return &TypeName{"int32"}, nil
		}
		return &TypeName{format}, nil
	}
	if typename == "boolean" {
		return &TypeName{"bool"}, nil
	}

	if typename == "string" && format == "binary" {
		return &TypeName{"bytes"}, nil
	}
	return &TypeName{typename}, nil
}
