package descriptor

import (
	"strings"
	"sync"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"google.golang.org/protobuf/proto"
)

type contents struct {
	service service
	message message
	enum    enum
}

func (c *contents) String() string {
	var sb strings.Builder // TODO: use sb.Glow()

	sb.WriteString(c.service.String())
	sb.WriteString("\n" + c.message.String())
	sb.WriteString("\n" + c.enum.String())

	return sb.String()
}

func Load(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, f := range req.GetProtoFile() {
		files[f.GetName()] = f
	}

	var wg sync.WaitGroup
	var res plugin.CodeGeneratorResponse
	for _, name := range req.GetFileToGenerate() {
		var content contents

		f := files[name]

		wg.Add(1)
		go func() {
			defer wg.Done()
			content.service = getService(f.GetService())
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			content.message = getMessage(f.GetMessageType())
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			content.enum = getEnum(f.GetEnumType())
		}()

		wg.Wait()

		res.File = append(res.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(name + ".graphqls"),
			Content: proto.String(content.String()),
		})
	}

	return &res, nil
}
