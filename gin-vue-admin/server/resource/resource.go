package resource

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:page
var pageFS embed.FS

func GetPageFS() http.FileSystem {
	subFS, err := fs.Sub(pageFS, "page")
	if err != nil {
		return http.FS(pageFS)
	}
	return http.FS(subFS)
}
