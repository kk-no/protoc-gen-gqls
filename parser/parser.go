package parser

import (
	"io"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"google.golang.org/protobuf/proto"
)

// Parse will convert the input to proto request format.
func Parse(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return parse(b)
}

func parse(b []byte) (*plugin.CodeGeneratorRequest, error) {
	var req plugin.CodeGeneratorRequest // do not pointer
	if err := proto.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	return &req, nil
}
