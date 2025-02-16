# Shell plugin

## Generate Go code from protoc file

```
cd plugins
protoc --go_out=. --go-grpc_out=. proto/shell.proto
```
