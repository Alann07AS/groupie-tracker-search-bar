package config

import (
	"html/template"
)

type Config struct {
	TemplateCache map[string]*template.Template
	Port          string
	Api           string
}

var appConfig Config

func ConfigLoad() *Config {
	return &appConfig
}
