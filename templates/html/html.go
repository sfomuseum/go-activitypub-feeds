package html

import (
	"context"
	"embed"
	"io/fs"
	
	sfomuseum_html "github.com/sfomuseum/go-template/html"
	"html/template"
)

//go:embed *.html
var FS embed.FS

func LoadTemplates(ctx context.Context, filesystems ...fs.FS) (*template.Template, error) {

	return sfomuseum_html.LoadTemplates(ctx, filesystems...)
}
