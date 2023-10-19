package shared

type Language string

const (
	LanguageEnglish Language = "en"
	LanguageGerman  Language = "de"
)

func DecodeLanguage(language string) Language {
	switch language {
	case "English":
		return LanguageEnglish
	case "German":
		return LanguageGerman
	default:
		return LanguageGerman
	}
}
