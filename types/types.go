package types

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var GQL = map[descriptor.FieldDescriptorProto_Type]string{
	descriptor.FieldDescriptorProto_TYPE_INT32:  "Int",
	descriptor.FieldDescriptorProto_TYPE_INT64:  "Int",
	descriptor.FieldDescriptorProto_TYPE_FLOAT:  "Float",
	descriptor.FieldDescriptorProto_TYPE_STRING: "String",
	descriptor.FieldDescriptorProto_TYPE_BOOL:   "Boolean",
}
