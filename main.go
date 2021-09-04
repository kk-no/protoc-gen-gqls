package main

import (
	"io"
	"log"
	"os"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/kk-no/protoc-gen-gqls/descriptor"
	"google.golang.org/protobuf/proto"
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
	return descriptor.Load(req)
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
