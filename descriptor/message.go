package descriptor

import (
	"fmt"
	"strings"

	"github.com/kk-no/protoc-gen-gqls/types"
	"google.golang.org/protobuf/types/descriptorpb"
)

type messages []*message

func (m messages) GQL() string {
	s := make([]string, len(m))
	for i, message := range m {
		s[i] = message.GQL()
	}
	return strings.Join(s, "\n")
}

func (m messages) String() string {
	s := make([]string, len(m))
	for i, message := range m {
		s[i] = message.String()
	}
	return strings.Join(s, "")
}

type message struct {
	messageType string
	name        string
	fields      []string
}

func (m message) GQL() string {
	return strings.Join([]string{
		m.messageType + m.name + Open,
		"\n" + strings.Join(m.fields, "\n"),
		"\n" + Close + "\n",
	}, "")
}

func (m message) String() string {
	return fmt.Sprintf("%s %s: %s", m.messageType, m.name, strings.Join(m.fields, ","))
}

func getMessages(messageTypes []*descriptorpb.DescriptorProto) messages {
	messages := make(messages, len(messageTypes))

	for i, messageType := range messageTypes {
		message := &message{
			name: messageType.GetName(),
		}

		if strings.Contains(message.name, "Request") {
			message.messageType = Input
		} else {
			message.messageType = Type
		}

		if fields := messageType.GetField(); len(fields) != 0 {
			message.fields = make([]string, len(fields))
			for i, field := range fields {
				var typeName string
				switch field.GetType() {
				case descriptorpb.FieldDescriptorProto_TYPE_ENUM, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
					// TODO: resolve dependency.
					typeName = pop(strings.Split(field.GetTypeName(), "."))
				default:
					typeName = types.Type[field.GetType()]
				}
				message.fields[i] = Indent + field.GetName() + ": " + types.Label(field.GetLabel()).GQL(typeName)
			}
		} else {
			message.fields = []string{Indent + "_: Boolean # noop field"}
		}
		messages[i] = message
	}

	return messages
}
