package descriptor

import (
	"fmt"
	"strings"

	options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type service struct {
	query    []string
	mutation []string
}

func (s service) String() string {
	var sb strings.Builder
	if len(s.query) != 0 {
		sb.WriteString("\n" + QueryType)
		sb.WriteString("\n" + strings.Join(s.query, "\n"))
		sb.WriteString("\n" + Close + "\n")
	}
	if len(s.mutation) != 0 {
		sb.WriteString("\n" + MutationType)
		sb.WriteString("\n" + strings.Join(s.mutation, "\n"))
		sb.WriteString("\n" + Close + "\n")
	}
	return sb.String()
}

func getService(serviceTypes []*descriptorpb.ServiceDescriptorProto) service {
	service := service{}

	for _, serviceType := range serviceTypes {
		for _, method := range serviceType.GetMethod() {
			in := pop(strings.Split(method.GetInputType(), "."))
			out := pop(strings.Split(method.GetOutputType(), "."))
			httpOpts := proto.GetExtension(method.GetOptions(), options.E_Http).(*options.HttpRule)
			switch {
			case httpOpts.GetGet() != "":
				service.query = append(service.query, fmt.Sprintf("%s%s(req: %s): %s", Indent, method.GetName(), in, out))
			case httpOpts.GetPost() != "":
				service.mutation = append(service.mutation, fmt.Sprintf("%s%s(req: %s): %s", Indent, method.GetName(), in, out))
			}
		}
	}

	return service
}
