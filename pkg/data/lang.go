package data

type Language string

const (
	English  Language = "en-us"
	Japanese          = "ja"
)

type Lang struct {
	currentLang Language
	langStrings map[string]string
}

var langObj *Lang

// Initializes language object, defaults to English
func InitLang() error {
	langStrings, err := ReadLanguageFile(English)

	if err != nil {
		return err
	}
	langObj = &Lang{
		currentLang: English,
		langStrings: langStrings,
	}
	return err
}

// Reads in the new language file and replaces all matching keys in the map
// Allows partial replacement of language in case there are missing translations
func ChangeLang(lang Language) {
	println("changin lang to ", lang)
	if langObj == nil {
		InitLang()
	}

	// water u doiin
	if lang == langObj.currentLang {
		return
	}

	// get the lang
	langStrings, err := ReadLanguageFile(lang)

	// if there's an error, let's just handle it gracefully
	if err != nil {
		println("oumch, had issies getin dat: ", err)
	}

	for key, value := range langStrings {
		langObj.langStrings[key] = value
	}

	langObj.currentLang = lang
}

// Retrieve a given string from map
func GiveMeString(code string) string {
	if langObj == nil {
		InitLang()
	}
	str, ok := langObj.langStrings[code]
	if !ok {
		return code
	}
	return str
}
