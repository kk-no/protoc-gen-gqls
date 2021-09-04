package descriptor

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type service []string

func (s service) String() string {
	return strings.Join(s, "\n")
}

func getService(serviceTypes []*descriptorpb.ServiceDescriptorProto) service {
	services := make(service, 0, len(serviceTypes))

	services = append(services, QueryType) // TODO: Implement other operations such as Mutation.
	for _, serviceType := range serviceTypes {
		services = append(services, Indent+"# "+serviceType.GetName())
		for _, method := range serviceType.GetMethod() {
			in := pop(strings.Split(method.GetInputType(), "."))
			out := pop(strings.Split(method.GetOutputType(), "."))
			services = append(services, fmt.Sprintf("%s%s(req: %s): %s", Indent, method.GetName(), in, out))
		}
		log.Printf("%s: %+v", serviceType.GetName(), serviceType)
	}
	services = append(services, Close+"\n")

	return services
}
