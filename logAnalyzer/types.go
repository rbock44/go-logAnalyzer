package logAnalyzer

const (
	REGEX_LEVEL_DUPLICATE        = "DUPLICATE"
	REGEX_LEVEL_SOURCEFILE       = "SOURCEFILE"
	REGEX_LEVEL_DUPLICATE_TRIM   = "DUPLICATE-TRIM"
	REGEX_LEVEL_DUPLICATE_DIFFER = "DUPLICATE-DIFFER"
	REGEX_LEVEL_LINE             = "LINE"
)

type NamedRegEx struct {
	Level string
	Name  string
	RegEx string
}

type IgnoreRegEx struct {
	Level string
	RegEx string
}
