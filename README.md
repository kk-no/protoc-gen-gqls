# protoc-gen-gqls
protoc-gen-gqls is a protoc plugin that generates graphql schema (.glaphqls) from Protocol Buffers.

## Use
```
$ protoc --plugin=protoc-gen-gqls --gqls_out=DIR example/proto/*.proto
```