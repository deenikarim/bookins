package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/deenikarim/bookings/internal/models"
	"html/template"
	"log"
)

//PART-3 ;;;; SETTING UP APPLICATION WIDE CONFIGURATION AND GO MAIN.GO FOR FINAL SETUP

//AppConfig setting application wide configuration
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	MailChan      chan models.MailData //channels for mainData type
}
