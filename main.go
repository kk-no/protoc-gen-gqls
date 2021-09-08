package main

import (
	"log"
	"os"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/kk-no/protoc-gen-gqls/descriptor"
	"github.com/kk-no/protoc-gen-gqls/parser"
	"google.golang.org/protobuf/proto"
)

func main() {
	req, err := parser.Parse(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := descriptor.Load(req)
	if err != nil {
		log.Fatalln(err)
	}

	if err := emit(res); err != nil {
		log.Fatalln(err)
	}
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
