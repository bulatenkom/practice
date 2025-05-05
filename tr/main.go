package main

import (
	"fmt"
	"os"
)

var app struct {
	lang string // "ru", "en"
	dict map[string]map[string]string
}

func init() {
	lang, _ := os.LookupEnv("LANG")
	lcAll, _ := os.LookupEnv("LC_ALL")

	if lang == "" && lcAll == "" {
		app.lang = "en"
		return
	}
	if lcAll != "" {
		lang = lcAll
	}

	switch lang {
	case "ru_RU.UTF-8":
		app.lang = "ru"
	case "C.UTF-8", "en_US.UTF-8", "en_GB.UTF-8":
	default:
		app.lang = "en"
	}

	app.dict = dict
}

func tr(text string) string {
	translations, ok := app.dict[text]
	if !ok {
		return text
	}

	if translation, ok := translations[app.lang]; ok {
		return translation
	}
	return text
}

func main() {
	fmt.Println(tr("Apple"))
	fmt.Println(tr("Something went wrong"))
	fmt.Println(tr("wtf?"))
	fmt.Println(tr("bla-bla"))
}
