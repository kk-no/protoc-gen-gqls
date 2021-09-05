package descriptor

import (
	"log"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type enum []string

func (e enum) String() string {
	return strings.Join(e, "\n")
}

func getEnum(packageName string, enumTypes []*descriptorpb.EnumDescriptorProto) enum {
	enums := make([]string, 0, len(enumTypes))

	for _, enumType := range enumTypes {
		enums = append(enums, Enum+enumType.GetName()+Open)
		for _, value := range enumType.GetValue() {
			enums = append(enums, Indent+value.GetName())
		}
		enums = append(enums, Close+"\n")
		log.Printf("%s: %+v", enumType.GetName(), enumType)
	}

	return enums
}
