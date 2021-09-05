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

func getMessage(messageTypes []*descriptorpb.DescriptorProto) message {
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
				if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
					messages = append(messages, Indent+field.GetName()+": "+pop(strings.Split(field.GetTypeName(), ".")))
				} else {
					messages = append(messages, Indent+field.GetName()+": "+types.GQL[field.GetType()])
				}
			}
		} else {
			messages = append(messages, Indent+"_: Boolean # noop field")
		}
		messages = append(messages, Close+"\n")
		log.Printf("%s: %+v", messageType.GetName(), messageType)
	}

	return messages
}
