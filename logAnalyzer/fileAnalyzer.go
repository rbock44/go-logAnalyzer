package logAnalyzer

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func IsFileOK(regularExpressions []NamedRegEx, regexToIdentifyIgnoredParts []IgnoreRegEx, duplicateBufferSize int, showProgress bool, fileUrlToAnalyze string) (result bool, hits []string) {

	result = true

	fileToAnalyze, err := os.Open(fileUrlToAnalyze)
	if err != nil {
		fmt.Println("[  ERROR  ] Could not open logfile.")
		fmt.Println()
		os.Exit(1)
	}
	defer fileToAnalyze.Close()

	fileStat, err := fileToAnalyze.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSizeInBytes := fileStat.Size()
	bar := pb.New64(fileSizeInBytes).SetUnits(pb.U_BYTES)

	if showProgress {
		bar.Start()
	}

	lineBuffer := make([]string, duplicateBufferSize)

	scanner := bufio.NewScanner(fileToAnalyze)
	var lineNumber = 0

	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())

		var duplicatesExist bool = false
		var hitOffset int = 0

		duplicatesExist, hitOffset, lineBuffer = checkForDuplicates(regularExpressions, regexToIdentifyIgnoredParts, lineBuffer, line, duplicateBufferSize)

		if duplicatesExist {
			hits = append(hits, "FILE/"+REGEX_LEVEL_DUPLICATE+" #"+strconv.Itoa(lineNumber-hitOffset)+"/"+strconv.Itoa(lineNumber)+" - "+line)
			result = false
		}

		lineLengthInBytes := len(line)
		lineIsOk, hitString, hitRegEx := IsLineOK(regularExpressions, regexToIdentifyIgnoredParts, line)
		if !lineIsOk {
			hits = append(hits, hitRegEx.Level+"/"+hitRegEx.Name+" #"+strconv.Itoa(lineNumber)+" - "+hitString)
			result = false
		}

		if showProgress {
			bar.Add(lineLengthInBytes)
		}

	}

	if showProgress {
		bar.Finish()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result, hits
}

func checkForDuplicates(regularExpressions []NamedRegEx, regexToIdentifyIgnoredParts []IgnoreRegEx, buffer []string, newLine string, duplicateBufferSize int) (result bool, hitOffset int, outbuffer []string) {

	hitOffset = duplicateBufferSize

	// if buffersize is samller than two there is nothing to compare
	if duplicateBufferSize < 2 {
		return false, 0, buffer
	}

	// update your checking window
	// if your window is size 5 you have to remove the
	// first line if the buffer has the size of 5
	if len(buffer) == duplicateBufferSize {
		// delete first element
		buffer = buffer[:0+copy(buffer[0:], buffer[1:])]
	}

	// should the current line be ignored ???
	for _, regEx := range regexToIdentifyIgnoredParts {
		if regEx.Level == REGEX_LEVEL_DUPLICATE {
			ignored, _ := regexp.MatchString(regEx.RegEx, newLine)
			if ignored {
				return false, 0, buffer
			}
		}
	}

	// trim the new line according to the various regEx specified by the user
	trimmedLine := newLine
	for _, regEx := range regularExpressions {
		if regEx.Level == REGEX_LEVEL_DUPLICATE_TRIM {
			r, _ := regexp.Compile(regEx.RegEx)
			trimmedLine = r.ReplaceAllString(trimmedLine, "")
		}
	}

	// check if the trimmed line
	// example from above the buffer has now the size of 4, by a window size of 5.
	// this OK because you have to add the newLine afterwards (but it is easier to check
	// for duplicates first)
	result = false
	for _, bufferLine := range buffer {
		hitOffset = hitOffset - 1
		if match(regularExpressions, bufferLine, trimmedLine) {
			result = true
			break
		}
	}

	buffer = append(buffer, trimmedLine)

	return result, hitOffset, buffer
}

func match(regularExpressions []NamedRegEx, left string, right string) (result bool) {
	foundDifferRegEx := false

	mustDifferLeft := ""
	mustDifferRight := ""

	trimmedLeft := left
	trimmedRight := right

	for _, namedRegExregEx := range regularExpressions {
		if namedRegExregEx.Level == REGEX_LEVEL_DUPLICATE_DIFFER {
			regEx, _ := regexp.Compile(namedRegExregEx.RegEx)

			findingsLeft := regEx.FindAllString(left, -1)
			mustDifferLeft = mustDifferLeft + strings.Join(findingsLeft, "-")
			trimmedLeft = regEx.ReplaceAllString(trimmedLeft, "")

			findingsRight := regEx.FindAllString(right, -1)
			mustDifferRight = mustDifferRight + strings.Join(findingsRight, "-")
			trimmedRight = regEx.ReplaceAllString(trimmedRight, "")

			foundDifferRegEx = true
		}
	}

	// no parts defined which should differ - so the result
	// is the comparision of the two trimmed parts
	if !foundDifferRegEx {
		return trimmedLeft != trimmedRight
	}

	// there IS something defined which should differ
	// but if the trimmed parts are not equal, the two
	// strings can not match
	if trimmedRight != trimmedLeft {
		return false
	}

	// the trimmed parts are equal and something should differ
	return mustDifferLeft != mustDifferRight

}
