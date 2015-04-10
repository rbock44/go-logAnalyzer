package logAnalyzer

const (
	REGEX_LEVEL_DUPLICATE = "DUPLICATE"
	REGEX_LEVEL_LINE      = "LINE"
)

type NamedRegEx struct {
	Level string
	Name  string
	RegEx string
}
