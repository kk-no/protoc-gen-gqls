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
		return nil, fmt.Errorf("parse/ReadAll: %w", err)
	}

	var req plugin.CodeGeneratorRequest // do not pointer
	if err := proto.Unmarshal(b, &req); err != nil {
		return nil, fmt.Errorf("parse/Unmarshal: %w", err)
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

		messageTypes := files[name].GetMessageType()
		messages := make([]string, 0, len(messageTypes))

		for _, messageType := range messageTypes {
			if fields := messageType.GetField(); len(fields) != 0 {
				messages = append(messages, "type "+messageType.GetName()+" {")
				for _, field := range fields {
					messages = append(messages, indent+field.GetName()+": "+types.GQL[field.GetType()])
				}
				messages = append(messages, "}\n")
			}
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
		return fmt.Errorf("emit/Marsgal: %w", err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		return fmt.Errorf("emit/Write: %w", err)
	}
	return nil
}
