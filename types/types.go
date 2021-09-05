package types

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var Type = map[descriptor.FieldDescriptorProto_Type]string{
	descriptor.FieldDescriptorProto_TYPE_INT32:  "Int",
	descriptor.FieldDescriptorProto_TYPE_INT64:  "Int",
	descriptor.FieldDescriptorProto_TYPE_FLOAT:  "Float",
	descriptor.FieldDescriptorProto_TYPE_STRING: "String",
	descriptor.FieldDescriptorProto_TYPE_BOOL:   "Boolean",
}

type Label descriptor.FieldDescriptorProto_Label

const (
	Optional = Label(descriptor.FieldDescriptorProto_LABEL_OPTIONAL)
	Required = Label(descriptor.FieldDescriptorProto_LABEL_REQUIRED)
	Repeated = Label(descriptor.FieldDescriptorProto_LABEL_REPEATED)
)

func (l Label) GQL(s string) string {
	switch l {
	case Optional:
		return s
	case Required:
		return s + "!"
	case Repeated:
		return "[" + s + "]"
	default:
		return s
	}
}
