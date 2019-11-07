package logAnalyzer

import "regexp"

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
	RegEx *regexp.Regexp
}

type IgnoreRegEx struct {
	Level string
	RegEx *regexp.Regexp
}
