package main

import (
	"fmt"
	"io"
	"log"
	"os"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"google.golang.org/protobuf/proto"
)

func main() {
	req, err := parse(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	res := process(req)

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
func process(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	var res plugin.CodeGeneratorResponse
	for _, name := range req.GetFileToGenerate() {
		out := name + ".graphqls"
		res.File = append(res.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(out),
			Content: proto.String("type PingRequest {}\n\ntype PingResponse {}"),
		})
	}
	return &res
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
