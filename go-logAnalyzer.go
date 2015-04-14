package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-logAnalyzer/logAnalyzer"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	var regExDictionrayFileUrl = flag.String("reg", "", " [MANDATORY] The url to the file containing regular expressions which are not allowed")
	var ignoreSrcFilesFileUrl = flag.String("ign", "", " [OPTIONAL] The url to the file containing the regular expressions identifing the source files which should be ignored")
	var logFileUrl = flag.String("log", "", " [MANDATORY] The url to the log-file or log-file folder")
	var outFileUrl = flag.String("out", "", " [OPTIONAL]  The url to the file containing the results")
	var duplicateBufferSize = flag.Int("buf", 0, " [OPTIONAL]  Indicates how many lines should be checked in a row for duplicates")
	var verbose = flag.Bool("v", false, "[OPTIONAL]  Verbose mode")

	flag.Parse()

	fmt.Println()
	fmt.Println("RegExFile          : " + *regExDictionrayFileUrl)
	if *ignoreSrcFilesFileUrl != "" {
		fmt.Println("IgnoreFile         : " + *ignoreSrcFilesFileUrl)
	}
	fmt.Println("LogFile/Folder     : " + *logFileUrl)
	fmt.Println()

	// try to get the regular expressions
	if *regExDictionrayFileUrl == "" || *logFileUrl == "" {
		fmt.Println("[  ERROR  ] Not all mandatory parameters are set.")
		fmt.Println()
		os.Exit(1)
	}

	regExFile, err := os.Open(*regExDictionrayFileUrl)
	if err != nil {
		fmt.Println("[  ERROR  ] Could not open regexfile.")
		fmt.Println()
		os.Exit(1)
	}
	defer regExFile.Close()

	regularExpressions := make([]logAnalyzer.NamedRegEx, 0)

	regexName := "DEFAULT"
	level := logAnalyzer.REGEX_LEVEL_LINE

	scanner := bufio.NewScanner(regExFile)
	for scanner.Scan() {
		var line = strings.TrimSpace(scanner.Text())

		// how can i loop through consts?
		if strings.HasPrefix(line, "## "+logAnalyzer.REGEX_LEVEL_LINE+"-SECTION") {
			regexName = strings.TrimSpace(strings.TrimPrefix(line, "## "+logAnalyzer.REGEX_LEVEL_LINE+"-SECTION"))
			level = logAnalyzer.REGEX_LEVEL_LINE
		}

		if strings.HasPrefix(line, "## "+logAnalyzer.REGEX_LEVEL_DUPLICATE_TRIM+"-SECTION") {
			regexName = strings.TrimSpace(strings.TrimPrefix(line, "## "+logAnalyzer.REGEX_LEVEL_DUPLICATE_TRIM+"-SECTION"))
			level = logAnalyzer.REGEX_LEVEL_DUPLICATE_TRIM
		}

		if strings.HasPrefix(line, "## "+logAnalyzer.REGEX_LEVEL_DUPLICATE_DIFFER+"-SECTION") {
			regexName = strings.TrimSpace(strings.TrimPrefix(line, "## "+logAnalyzer.REGEX_LEVEL_DUPLICATE_DIFFER+"-SECTION"))
			level = logAnalyzer.REGEX_LEVEL_DUPLICATE_DIFFER
		}

		if !strings.HasPrefix(line, "##") && line != "" && line != "\n" {
			_, err := regexp.Compile(line)
			if err != nil {
				fmt.Println("[  ERROR  ] Could not compile regular expression. " + line)
				fmt.Println("[  ERROR  ] Error was: " + err.Error())
				fmt.Println()
				os.Exit(1)
			}
			regularExpression := logAnalyzer.NamedRegEx{Level: level, Name: regexName, RegEx: line}
			regularExpressions = append(regularExpressions, regularExpression)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// try to get the files to ignore
	regexToIdentifyIgnoredParts := make([]logAnalyzer.IgnoreRegEx, 0)

	if *ignoreSrcFilesFileUrl != "" {
		logFile, err := os.Open(*ignoreSrcFilesFileUrl)
		if err != nil {
			fmt.Println("[  ERROR  ] Could not open ignoreFiles file.")
			fmt.Println()
			os.Exit(1)
		}
		defer logFile.Close()

		level = logAnalyzer.REGEX_LEVEL_LINE

		scanner := bufio.NewScanner(logFile)
		for scanner.Scan() {
			var line = strings.TrimSpace(scanner.Text())

			// how can i loop through consts?
			if strings.HasPrefix(line, "## "+logAnalyzer.REGEX_LEVEL_SOURCEFILE+"-SECTION") {
				level = logAnalyzer.REGEX_LEVEL_SOURCEFILE
			}

			if strings.HasPrefix(line, "## "+logAnalyzer.REGEX_LEVEL_DUPLICATE+"-SECTION") {
				level = logAnalyzer.REGEX_LEVEL_DUPLICATE
			}

			if !strings.HasPrefix(line, "##") && line != "" && line != "\n" {
				_, err := regexp.Compile(line)
				if err != nil {
					fmt.Println("[  ERROR  ] Could not compile regular expression. " + line)
					fmt.Println("[  ERROR  ] Error was: " + err.Error())
					fmt.Println()
					os.Exit(1)
				}
				regularExpression := logAnalyzer.IgnoreRegEx{Level: level, RegEx: line}
				regexToIdentifyIgnoredParts = append(regexToIdentifyIgnoredParts, regularExpression)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	fileToAnalyze, err := os.Open(*logFileUrl)
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

	var logsAreOk = true
	var logMsgs []string

	if fileStat.IsDir() {

		files, _ := fileToAnalyze.Readdirnames(-1)

		for _, childName := range files {

			childOutFileName := ""

			if *outFileUrl != "" {
				fileNameParts := strings.Split(*outFileUrl, ".")

				newOutFileName := ""

				if len(fileNameParts) == 1 {

					newOutFileName = *outFileUrl + "_" + childName

				} else {

					for i, part := range fileNameParts {

						if i != 0 && i != len(fileNameParts)-1 {
							newOutFileName = newOutFileName + "."
						}

						if i == len(fileNameParts)-1 {
							newOutFileName = newOutFileName + "_" + childName + "."
						}

						newOutFileName = newOutFileName + part

					}
				}
				childOutFileName = newOutFileName
			}

			logFileIsOk, logMsg := processLogFile(regularExpressions, regexToIdentifyIgnoredParts, *duplicateBufferSize, filepath.Join(fileToAnalyze.Name(), childName), *verbose, childOutFileName)

			if !logFileIsOk {
				logMsgs = append(logMsgs, logMsg)
			}

			logsAreOk = logsAreOk && logFileIsOk

		}

	} else {

		logFileIsOk, logMsg := processLogFile(regularExpressions, regexToIdentifyIgnoredParts, *duplicateBufferSize, fileToAnalyze.Name(), *verbose, *outFileUrl)

		if !logFileIsOk {
			logMsgs = append(logMsgs, logMsg)
		}

		logsAreOk = logsAreOk && logFileIsOk
	}

	if !logsAreOk {
		if !*verbose && *outFileUrl == "" {
			fmt.Println("[ WARNING ] You do not see any detailed results (no outfile is written & no screen output) because of our parameter-settings!")
		}

		for _, logMsg := range logMsgs {
			fmt.Println(logMsg)
		}

		fmt.Println()

		os.Exit(1)
	}
	fmt.Println("[   O K   ] NO positive test in logfile(s) ...")
	fmt.Println()

	os.Exit(0)
}

func processLogFile(regularExpressions []logAnalyzer.NamedRegEx, regexToIdentifyIgnoredParts []logAnalyzer.IgnoreRegEx, duplicateBufferSize int, logFileUrl string, verbose bool, outFileUrl string) (logFileIsOk bool, logMsg string) {
	fmt.Println("processing LogFile : " + logFileUrl)

	logFileIsOk = true
	logMsg = ""
	var hits []string

	logFileIsOk, hits = logAnalyzer.IsFileOK(regularExpressions, regexToIdentifyIgnoredParts, duplicateBufferSize, logFileUrl)

	if !logFileIsOk {
		logMsg = "[  ERROR  ] " + logFileUrl + " - has " + strconv.Itoa(len(hits)) + " positive test(s) ..."
	}

	printErrors(logFileUrl, verbose, hits)
	writeErrorFile(outFileUrl, hits)

	return logFileIsOk, logMsg
}

func printErrors(logfileName string, verbose bool, hits []string) {
	if verbose && len(hits) > 0 {
		for _, hit := range hits {
			fmt.Println(hit)
		}
	}
}

func writeErrorFile(errorFileUrl string, hits []string) {

	if errorFileUrl != "" && len(hits) > 0 {

		outFileUrlInternal := errorFileUrl

		// check if file exists, if so - change the name so it is unique
		// means --> append the timestamp
		if _, err := os.Stat(outFileUrlInternal); err == nil {
			now := time.Now().Local()
			outFileUrlInternal = outFileUrlInternal + "_" + now.Format(time.RFC3339)
			fmt.Println("[ WARNING ] Specified outfile exists, will write to " + outFileUrlInternal)
		}

		outFile, err := os.Create(outFileUrlInternal)
		if err != nil {
			fmt.Println("[  ERROR  ] Could not open outputFile.")
			fmt.Println()
			os.Exit(1)
		}
		defer outFile.Close()

		writer := bufio.NewWriter(outFile)
		for _, hit := range hits {
			writer.WriteString(hit + "\n")
		}
		writer.Flush()
	}
}
