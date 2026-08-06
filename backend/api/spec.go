package api

import (
	"embed"
	"io/fs"
)

//go:embed openapi-v1.yaml
var SpecFS embed.FS

var SpecBytesV1 []byte

func init() {
	var err error
	SpecBytesV1, err = fs.ReadFile(SpecFS, "openapi-v1.yaml")
	if err != nil {
		panic("openapi-v1.yaml not found in embed: " + err.Error())
	}
}
