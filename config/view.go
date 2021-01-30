package config

import (
	"github.com/gofiber/template/html"
	"github.com/markbates/pkger"
	"path/filepath"
)

type TemplateConfig struct {
	TemplateEngine *html.Engine
	Path           string `yaml:"path" env-default:"resources/views"`
	Extension      string `yaml:"extension" env-default:".html"`
}
type ViewConfig struct {
	Template TemplateConfig `yaml:"template"`
}

func (v *ViewConfig) Load(path string) {
	path = MakeDir(filepath.Join(path, v.Template.Path))
	v.Template.TemplateEngine = html.NewFileSystem(pkger.Dir(path), v.Template.Extension)
}
