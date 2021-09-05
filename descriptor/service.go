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
		s.gqlStr(&sb, QueryType, s.query)
	}
	if len(s.mutation) != 0 {
		s.gqlStr(&sb, MutationType, s.mutation)
	}
	return sb.String()
}

func (s service) gqlStr(sb *strings.Builder, opType string, opContents []string) {
	sb.WriteString("\n" + opType)
	sb.WriteString("\n" + strings.Join(opContents, "\n"))
	sb.WriteString("\n" + Close + "\n")
}

func getService(packageName string, serviceTypes []*descriptorpb.ServiceDescriptorProto) service {
	service := service{}

	for _, serviceType := range serviceTypes {
		for _, method := range serviceType.GetMethod() {
			in := pop(strings.Split(method.GetInputType(), "."))
			out := pop(strings.Split(method.GetOutputType(), "."))
			httpOpt := proto.GetExtension(method.GetOptions(), options.E_Http).(*options.HttpRule)
			methodStr := fmt.Sprintf("%s%s(req: %s): %s", Indent, method.GetName(), in, out)

			switch {
			case httpOpt.GetGet() != "":
				service.query = append(service.query, methodStr)
			case httpOpt.GetPost() != "", httpOpt.GetPut() != "", httpOpt.GetPatch() != "", httpOpt.GetDelete() != "":
				service.mutation = append(service.mutation, methodStr)
			default:
				continue // No HTTP option
			}
		}
	}

	return service
}
