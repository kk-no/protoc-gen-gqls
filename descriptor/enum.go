package descriptor

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type enums []*enum

func (e enums) GQL() string {
	s := make([]string, len(e))
	for i, enum := range e {
		s[i] = enum.GQL()
	}
	return strings.Join(s, "\n")
}

func (e enums) String() string {
	s := make([]string, len(e))
	for i, enum := range e {
		s[i] = enum.String()
	}
	return strings.Join(s, "")
}

type enum struct {
	name   string
	values []string
}

func (e enum) GQL() string {
	return strings.Join([]string{
		Enum + e.name + Open,
		strings.Join(e.values, "\n"),
		Close + "\n",
	}, "\n")
}

func (e enum) String() string {
	return fmt.Sprintf("%s: %s", e.name, strings.Join(e.values, ","))
}

func getEnums(enumTypes []*descriptorpb.EnumDescriptorProto) enums {
	enums := make(enums, len(enumTypes))

	for i, enumType := range enumTypes {
		enum := &enum{
			name: enumType.GetName(),
		}
		values := make([]string, len(enumType.GetValue()))
		for i, value := range enumType.GetValue() {
			values[i] = Indent + value.GetName()
		}
		enum.values = values
		enums[i] = enum
	}

	return enums
}
