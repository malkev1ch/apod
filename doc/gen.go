package doc

//go:generate oapi-codegen -generate types -o ../gen/v1/openapi_types.gen.go -package gen v1/picture.yaml
//go:generate oapi-codegen -generate server -o ../gen/v1/openapi_server.gen.go -package gen v1/picture.yaml
