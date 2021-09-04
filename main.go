package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/kk-no/protoc-gen-gqls/types"
	"google.golang.org/protobuf/proto"
)

var (
	defaultIndent = 4
	indent        = strings.Repeat(" ", defaultIndent)
)

func main() {
	req, err := parse(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := process(req)
	if err != nil {
		log.Fatalln(err)
	}

	if err := emit(res); err != nil {
		log.Fatalln(err)
	}
}

// parse will convert the input to proto request format.
func parse(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var req plugin.CodeGeneratorRequest // do not pointer
	if err := proto.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// process will parse request and format the file into the output format.
func process(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, f := range req.GetProtoFile() {
		files[f.GetName()] = f
	}

	var res plugin.CodeGeneratorResponse
	for _, name := range req.GetFileToGenerate() {
		var content strings.Builder

		f := files[name]

		// Service
		serviceTypes := f.GetService()
		services := make([]string, 0, len(serviceTypes))

		services = append(services, "type Query {") // TODO: Implement other operations such as Mutation.
		for _, serviceType := range serviceTypes {
			services = append(services, indent+"# "+serviceType.GetName())
			for _, method := range serviceType.GetMethod() {
				in := pop(strings.Split(method.GetInputType(), "."))
				out := pop(strings.Split(method.GetOutputType(), "."))
				services = append(services, fmt.Sprintf("%s%s(req: %s): %s", indent, method.GetName(), in, out))
			}
			log.Printf("%s: %+v", serviceType.GetName(), serviceType)
		}
		services = append(services, "}\n\n")

		if _, err := content.WriteString(strings.Join(services, "\n")); err != nil {
			return nil, err
		}

		// Message
		messageTypes := f.GetMessageType()
		messages := make([]string, 0, len(messageTypes))

		for _, messageType := range messageTypes {
			fields := messageType.GetField()
			if strings.Contains(messageType.GetName(), "Request") {
				messages = append(messages, "input "+messageType.GetName()+" {")
			} else {
				messages = append(messages, "type "+messageType.GetName()+" {")
			}
			if len(fields) != 0 {
				for _, field := range fields {
					messages = append(messages, indent+field.GetName()+": "+types.GQL[field.GetType()])
				}
			} else {
				messages = append(messages, indent+"_: Boolean # noop field")
			}
			messages = append(messages, "}\n")
			log.Printf("%s: %+v", messageType.GetName(), messageType)
		}

		if _, err := content.WriteString(strings.Join(messages, "\n")); err != nil {
			return nil, err
		}

		out := name + ".graphqls"
		res.File = append(res.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(out),
			Content: proto.String(content.String()),
		})
	}
	return &res, nil
}

// emit will output the file in the name and content passed to it.
func emit(res *plugin.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		return err
	}
	return nil
}

// pop extracts the end value of the slice.
func pop(s []string) string {
	return s[len(s)-1]
}
