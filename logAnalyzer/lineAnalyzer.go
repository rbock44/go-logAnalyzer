package logAnalyzer

// IsLineOK checks whether the line contains a regex
func IsLineOK(regularExpressions []NamedRegEx, regexToIdentifyIgnoredParts []IgnoreRegEx, stringToAnalyze string) (result bool, hitString string, hit *NamedRegEx) {
	// should the current line be ignored ???
	for _, regEx := range regexToIdentifyIgnoredParts {
		if regEx.Level == REGEX_LEVEL_SOURCEFILE {
			ignored := regEx.RegEx.MatchString(stringToAnalyze)
			if ignored {
				return true, "", nil
			}
		}
	}

	for _, regEx := range regularExpressions {
		if regEx.Level == REGEX_LEVEL_LINE {
			match := regEx.RegEx.MatchString(stringToAnalyze)
			if match {
				return false, stringToAnalyze, &regEx
			}
		}
	}

	return true, "", nil
}
