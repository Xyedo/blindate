package locale

import (
	"embed"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed active.*.toml
var LocaleFS embed.FS
var Bundle *i18n.Bundle

func New() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := Bundle.LoadMessageFileFS(LocaleFS, "active.en.toml")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = Bundle.LoadMessageFileFS(LocaleFS, "active.id.toml")
	if err != nil {
		log.Fatalln(err)
	}
}
