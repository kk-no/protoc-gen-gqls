package descriptor

import (
	"fmt"
	"strings"

	options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type services []*service

func (s services) GQL() string {
	gql := make([]string, 0, len(s))
	for _, service := range s {
		gql = append(gql, service.GQL())
	}
	return strings.Join(gql, "\n")
}

func (s services) String() string {
	return "" // TODO: implements
}

type service struct {
	name     string
	query    *query
	mutation *mutation
}

func (s service) GQL() string {
	gql := make([]string, 0, 2)
	if len(s.query.methods) != 0 {
		gql = append(gql, s.query.GQL())
	}
	if len(s.mutation.methods) != 0 {
		gql = append(gql, s.mutation.GQL())
	}
	return strings.Join(gql, "\n")
}

func (s service) String() string {
	return fmt.Sprintf("%s: %s, %s", s.name, s.query.String(), s.mutation.String())
}

type query struct {
	methods methods
}

func (q query) GQL() string {
	if len(q.methods) == 0 {
		return ""
	}
	return strings.Join([]string{
		QueryType,
		q.methods.GQL(),
		Close + "\n",
	}, "\n")
}

func (q query) String() string {
	return fmt.Sprintf("Query: %s", q.methods.String())
}

type mutation struct {
	methods methods
}

func (m mutation) GQL() string {
	if len(m.methods) == 0 {
		return ""
	}
	return strings.Join([]string{
		MutationType,
		m.methods.GQL(),
		Close + "\n",
	}, "\n")
}

func (m mutation) String() string {
	return fmt.Sprintf("Mutation: %s", m.methods.String())
}

type methods []*method

func (m methods) GQL() string {
	s := make([]string, len(m))
	for i, method := range m {
		s[i] = method.GQL()
	}
	return strings.Join(s, "\n")
}

func (m methods) String() string {
	s := make([]string, len(m))
	for i, method := range m {
		s[i] = method.String()
	}
	return strings.Join(s, "")
}

type method struct {
	name     string
	request  string
	response string
}

func (m method) GQL() string {
	return fmt.Sprintf("%s%s(req: %s): %s", Indent, m.name, m.request, m.response)
}

func (m method) String() string {
	return fmt.Sprintf("%s(req: %s): %s", m.name, m.request, m.response)
}

func getServices(serviceTypes []*descriptorpb.ServiceDescriptorProto) services {
	services := make(services, len(serviceTypes))

	for i, serviceType := range serviceTypes {
		serviceMethods := serviceType.GetMethod()

		service := &service{
			name: serviceType.GetName(),
			query: &query{
				methods: make(methods, 0, len(serviceMethods)),
			},
			mutation: &mutation{
				methods: make(methods, 0, len(serviceMethods)),
			},
		}
		for _, serviceMethod := range serviceMethods {
			method := &method{
				name:     serviceMethod.GetName(),
				request:  pop(strings.Split(serviceMethod.GetInputType(), ".")),
				response: pop(strings.Split(serviceMethod.GetOutputType(), ".")),
			}

			httpOpt := proto.GetExtension(serviceMethod.GetOptions(), options.E_Http).(*options.HttpRule)
			switch {
			case httpOpt.GetGet() != "":
				service.query.methods = append(service.query.methods, method)
			case httpOpt.GetPost() != "", httpOpt.GetPut() != "", httpOpt.GetPatch() != "", httpOpt.GetDelete() != "":
				service.mutation.methods = append(service.mutation.methods, method)
			default:
				continue // No HTTP option
			}
		}
		services[i] = service
	}

	return services
}
