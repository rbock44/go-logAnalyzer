package logAnalyzer

import (
	"regexp"
)

// IsLineOK checks whether the line contains a regex
func IsLineOK(regularExpressions []NamedRegEx, regexToIdentifyIgnoredParts []IgnoreRegEx, stringToAnalyze string) (result bool, hitString string, hit NamedRegEx) {
	var nilNamedRegEx NamedRegEx

	// should the current line be ignored ???
	for _, regEx := range regexToIdentifyIgnoredParts {
		if regEx.Level == REGEX_LEVEL_SOURCEFILE {
			ignored, _ := regexp.MatchString(regEx.RegEx, stringToAnalyze)
			if ignored {
				return true, "", nilNamedRegEx
			}
		}
	}

	for _, regEx := range regularExpressions {
		if regEx.Level == REGEX_LEVEL_LINE {
			match, _ := regexp.MatchString(regEx.RegEx, stringToAnalyze)
			if match {
				return false, stringToAnalyze, regEx
			}
		}
	}

	return true, "", nilNamedRegEx
}
