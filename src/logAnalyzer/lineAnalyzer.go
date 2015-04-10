package logAnalyzer

import (
	"regexp"
)

func IsLineOK(regularExpressions []NamedRegEx, ignoredSrcFiles []string, stringToAnalyze string) (result bool, hitString string, hit NamedRegEx) {
	var nilNamedRegEx NamedRegEx

	for _, ignSrcFile := range ignoredSrcFiles {
		srcFileIgnored, _ := regexp.MatchString(ignSrcFile, stringToAnalyze)
		if srcFileIgnored {
			return true, "", nilNamedRegEx
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
