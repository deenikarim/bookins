package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
)

//PART-3 ;;;; SETTING UP APPLICATION WIDE CONFIGURATION AND GO MAIN.GO FOR FINAL SETUP

//AppConfig setting application wide configuration
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
}
