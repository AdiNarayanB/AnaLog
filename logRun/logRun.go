package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type PatternSpec struct {
	PatternContent string `json:"patternContent"`
	Severity       int    `json:"patternSeverity"`
}

type Hit struct {
	pattern   PatternSpec
	frequency int
}

type filenameRegex struct {
	filename string
	hitList  []Hit
}

type termCount struct {
	term  string
	count int
}

type filenameTerms struct {
	filename      string
	termCountList []termCount
}

type filenameRegexTerms struct {
	filename  string
	RegexList filenameRegex
	TermList  filenameTerms
}

type aggregateResult struct {
	fileRegexTerms []filenameRegexTerms
}

func getTextFromFile(filePath string) string {
	fileBytes, err := os.ReadFile(filePath)
	fmt.Println("the file path is:" + filePath)
	if err != nil {
		return ""
	}
	fileText := string(fileBytes)
	fmt.Println(filePath)
	return fileText
}

func getAllFilePaths(directoryPath string) []string {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}
	var filePathList []string
	for _, file := range files {

		filePathList = append(filePathList, file.Name())

	}
	return filePathList
}

func parsePatternJson(patternFilePath string) map[string]int {

	jsonFileBytes, _ := os.ReadFile(patternFilePath)
	fmt.Println(string(jsonFileBytes))
	var patternSpecs []PatternSpec

	patternSpecMap := make(map[string]int)
	err := json.Unmarshal(jsonFileBytes, &patternSpecs)
	fmt.Println(patternSpecs[0].Severity)
	if err != nil {
		fmt.Println("Pattern Json Parsing failed")
	}

	for i := 0; i < len(patternSpecs); i += 1 {
		patternSpecMap[patternSpecs[i].PatternContent] = patternSpecs[i].Severity

	}
	fmt.Println(patternSpecMap)
	return patternSpecMap
}

func getTermCount(logText string, contextRange int) []termCount {

	logTextSplitLine := strings.Split(logText, "\n")
	termCountMap := make(map[string]int)

	for i := range logTextSplitLine {
		LogLine := logTextSplitLine[i]
		phrases := strings.Split(LogLine, " ")
		var termPhraseList []string
		for i := range phrases {
			if contextRange < len(phrases) {
				termPhrase := strings.Join(phrases[i:i+contextRange-1], " ")
				termPhraseList = append(termPhraseList, termPhrase)
			} else {
				termPhraseList = append(termPhraseList, strings.Join(phrases, " "))
			}
		}
		for i := range termPhraseList {
			_, ok := termCountMap[termPhraseList[i]]

			if ok {
				termCountMap[termPhraseList[i]] = termCountMap[termPhraseList[i]] + 1
			} else {
				termCountMap[termPhraseList[i]] = 1

			}
		}

	}

	var termCountList []termCount

	for phrase, count := range termCountMap {
		fmt.Println(phrase)
		var termCountItem = new(termCount)
		termCountItem.term = phrase
		termCountItem.count = count
		termCountList = append(termCountList, *termCountItem)

	}

	return termCountList

}

func aggrTermCounts(directoryPath string, contextLength int) []filenameTerms {
	filePathList := getAllFilePaths(directoryPath)
	var filenameTermList []filenameTerms
	for i := range filePathList {
		globalTermCountMap := make(map[string]int)
		filepathText := getTextFromFile(directoryPath + filePathList[i])
		termCountListFile := getTermCount(filepathText, contextLength)
		fmt.Println(termCountListFile)
		for termCountItemIndex := range termCountListFile {
			var termCountItemVal = reflect.ValueOf(termCountListFile[termCountItemIndex])

			var termCountItemPhrase = reflect.Indirect(termCountItemVal).FieldByName("term").String()
			fmt.Println(termCountItemPhrase)
			count, ok := globalTermCountMap[termCountItemPhrase]
			if ok {
				globalTermCountMap[termCountItemPhrase] = globalTermCountMap[termCountItemPhrase] + count
			} else {
				globalTermCountMap[termCountItemPhrase] = count
			}
		}

		phrases := make([]string, 0, len(globalTermCountMap))
		for phrase := range globalTermCountMap {
			phrases = append(phrases, phrase)
		}

		sort.SliceStable(phrases, func(i, j int) bool {
			return globalTermCountMap[phrases[i]] < globalTermCountMap[phrases[j]]
		})
		var termCountList []termCount
		for phrase, count := range globalTermCountMap {
			var termCountItem = new(termCount)
			termCountItem.term = phrase
			fmt.Println(count)
			termCountItem.count = count
			termCountList = append(termCountList, *termCountItem)

		}
		var fileTermCountList = new(filenameTerms)
		fileTermCountList.termCountList = termCountList
		fileTermCountList.filename = filePathList[i]

		filenameTermList = append(filenameTermList, *fileTermCountList)
	}
	fmt.Println(filenameTermList)
	return filenameTermList

}

func getSeverityFromPatternSpec(pspec *PatternSpec, prop string) string {

	var pspecval = reflect.ValueOf(pspec)
	return reflect.Indirect(pspecval.FieldByName(prop)).String()

}

func splitLogTextIntoLines(logText string) []string {
	return strings.Split(logText, "\n")

}

func applyPatternOnText(logName string, logText string, patternList []PatternSpec) filenameRegex {

	hitList := make([]Hit, 0, len(patternList))

	for i := 0; i < len(patternList); i += 1 {
		var PatternContent = reflect.ValueOf(patternList[i])
		var PatternContentVal = reflect.Indirect(PatternContent).FieldByName("PatternContent").String()

		var PatternSeverity = reflect.ValueOf(patternList[i])
		fmt.Println(reflect.Indirect(PatternSeverity).FieldByName("Severity"))
		var PatternSeverityVal = reflect.Indirect(PatternSeverity).FieldByName("Severity").Int()
		patternCompiled := regexp.MustCompile(PatternContentVal)
		fmt.Println(patternCompiled)
		fmt.Println(logText)

		logTextLines := splitLogTextIntoLines(logText)

		for i := 0; i < len(logTextLines); i++ {
			fmt.Println(strings.Split(logText, "\n"))

			fmt.Println(logTextLines[i] + " and " + PatternContentVal)
			matches := patternCompiled.FindAllStringIndex(logTextLines[i], -1)

			for i := 0; i < len(matches); i++ {

				var patternHit = new(Hit)
				var pspec = new(PatternSpec)
				pspec.PatternContent = PatternContentVal
				pspec.Severity = int(PatternSeverityVal)
				patternHit.frequency = len(matches)
				patternHit.pattern = *pspec
				hitList = append(hitList, *patternHit)
			}
		}

	}
	fmt.Println(hitList)
	sort.Slice(hitList, func(i, j int) bool {

		return getSeverityFromPatternSpec(&hitList[i].pattern, "Severity") < getSeverityFromPatternSpec(&hitList[j].pattern, "Severity")

	})
	var filenameRegexItem = new(filenameRegex)
	filenameRegexItem.filename = logName
	filenameRegexItem.hitList = hitList

	return *filenameRegexItem
}

func cvtPatternSpecMap2List(patternListMap map[string]int) []PatternSpec {
	var patternSpecList []PatternSpec
	for patternContent, patternSeverity := range patternListMap {
		var patternSpecItem = new(PatternSpec)
		patternSpecItem.PatternContent = patternContent
		patternSpecItem.Severity = patternSeverity
		patternSpecList = append(patternSpecList, *patternSpecItem)
	}
	return patternSpecList

}

func aggrMatches(directoryPath string, patternFilePath string) []filenameRegex {

	filePathList := getAllFilePaths(directoryPath)

	filenameRegexList := make([]filenameRegex, 0, len(filePathList))
	patternList := parsePatternJson(patternFilePath)

	for i := 0; i < len(filePathList); i++ {

		logText := getTextFromFile(filePathList[i])
		filenameRegexItem := applyPatternOnText(filePathList[i], logText, cvtPatternSpecMap2List(patternList))
		filenameRegexList = append(filenameRegexList, filenameRegexItem)

	}
	return filenameRegexList

}

func getfilenameTermList(filename string, filenameTermCountList []filenameTerms) (filenameTerms, error) {
	var placeHolder = filenameTermCountList[0]
	for i := 0; i < len(filenameTermCountList); i++ {
		if filenameTermCountList[i].filename == filename {
			return filenameTermCountList[i], nil
		}

	}
	return placeHolder, errors.New("file Not found")
}

func getfilenameRegexList(filename string, filenameRegexList []filenameRegex) (filenameRegex, error) {
	var placeHolder = filenameRegexList[0]

	for i := 0; i < len(filenameRegexList); i++ {
		if filenameRegexList[i].filename == filename {
			return filenameRegexList[i], nil
		}
	}
	return placeHolder, errors.New("no filename found")
}

func aggrResults(directoryPath string, patternFilePath string, contextLength int) aggregateResult {

	filePathList := getAllFilePaths(directoryPath)
	filenameRegexList := aggrMatches(directoryPath, patternFilePath)
	fmt.Println(filenameRegexList)
	filenameTermList := aggrTermCounts(directoryPath, contextLength)

	filenameRegexTermsList := make([]filenameRegexTerms, 0, len(filePathList))
	for i := 0; i < len(filePathList); i++ {
		filenameTermCountItem, _ := getfilenameTermList(filePathList[i], filenameTermList)
		filenameRegexItem, _ := getfilenameRegexList(filePathList[i], filenameRegexList)
		var filenameRegexTermItem = new(filenameRegexTerms)
		filenameRegexTermItem.RegexList = filenameRegexItem
		filenameRegexTermItem.filename = filePathList[i]
		filenameRegexTermItem.TermList = filenameTermCountItem

		filenameRegexTermsList = append(filenameRegexTermsList, *filenameRegexTermItem)
	}
	var aggRes = new(aggregateResult)
	aggRes.fileRegexTerms = filenameRegexTermsList

	return *aggRes

}

func makeResultObjIntoJSON(resultObj aggregateResult) string {

	resObj := &resultObj

	b, err := json.Marshal(resObj)
	if err != nil {
		fmt.Println(err)
		return "error"
	}

	return string(b)

}

func main() {

	directoryPath := "/home/adithya/logProc/test/Data/"
	fmt.Println(directoryPath)
	patternFilePath := "/home/adithya/logProc/test/Patterns/patternSpec.json"
	fmt.Println(patternFilePath)
	contextLength := 2
	aggResult := aggrResults(directoryPath, patternFilePath, contextLength)
	jsonStr := makeResultObjIntoJSON(aggResult)
	fmt.Println(jsonStr)

}
