package descriptor

import (
	"log"
	"strings"

	"github.com/kk-no/protoc-gen-gqls/types"
	"google.golang.org/protobuf/types/descriptorpb"
)

type message []string

func (m message) String() string {
	return strings.Join(m, "\n")
}

func getMessage(packageName string, messageTypes []*descriptorpb.DescriptorProto) message {
	messages := make(message, 0, len(messageTypes))

	for _, messageType := range messageTypes {
		fields := messageType.GetField()
		if strings.Contains(messageType.GetName(), "Request") {
			messages = append(messages, Input+messageType.GetName()+Open)
		} else {
			messages = append(messages, Type+messageType.GetName()+Open)
		}

		if len(fields) != 0 {
			for _, field := range fields {
				var sb strings.Builder
				sb.WriteString(Indent + field.GetName() + ": ")

				var typeName string
				switch field.GetType() {
				case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
					typeName = pop(strings.Split(field.GetTypeName(), "."))
				case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
					if strings.HasPrefix(field.GetTypeName(), "."+packageName) {
						typeName = pop(strings.Split(field.GetTypeName(), "."))
					} else {
						typeName = "String" // TODO: fix dependency
					}
				default:
					typeName = types.Type[field.GetType()]
				}

				sb.WriteString(types.Label(field.GetLabel()).GQLStr(typeName))
				messages = append(messages, sb.String())
			}
		} else {
			messages = append(messages, Indent+"_: Boolean # noop field")
		}
		messages = append(messages, Close+"\n")
		log.Printf("%s: %+v", messageType.GetName(), messageType)
	}

	return messages
}
