package main

import (
	"flag"
	"log"
	"os"
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const splunkDateFormat string = "01/02/2006 15:04:05 -0700"
const timestampGroupName = "_timestamp"

var configFile = flag.String("c", "", "location of configuration json file")

func initialise() (*os.File, Config) {
	flag.Parse()

	if *configFile == "" {
		log.Fatal("A config file must be specified")
	}
	if len(flag.Args()) == 0 {
		log.Fatal("Log files must be specified for processing")
	}
	logFile := initLogging("log-analyse.log")
	config, err := loadConfig(*configFile)
	checkError("initialise", err)

	return logFile, config
}

func processLogLine(load []Load, line string, lineNumber int) (time.Time, map[string]string, []string) {
	var result map[string]string = make(map[string]string);
	var timestamp time.Time
	var groupNames []string

	for _, thisLoad := range load {
		re, err := regexp.Compile(thisLoad.Regexp)
		checkError("compile regexp", err)
		matches := re.FindAllStringSubmatch(line, -1)
		groupNames = thisLoad.GroupNames

		if matches != nil {
			if len(matches) > 1 {
				log.Printf("Warning: More tokens than expected (%d) at line %d", len(matches), lineNumber)
			}
			matchLine := matches[0]

			for index := 1 ; index < len(matchLine) ; index++ {
				var groupName string
				if (index > len(thisLoad.GroupNames)){
					groupName = fmt.Sprintf("Column%d", index)
				} else {
					groupName = thisLoad.GroupNames[index - 1]
				}

				result[groupName] = matchLine[index]
			}

			if timestampStr := result[timestampGroupName]; timestampStr != "" {
				timestamp, err = time.Parse(thisLoad.TimestampFormat, timestampStr)
				checkError("log line time parse", err)
			}

			break
		}
	}

	return timestamp, result, groupNames
}

func printMap(timestamp time.Time, miscOptions MiscOptions, m map[string]string, fieldList []string){
	if timestamp.IsZero() == false {
		fmt.Printf("[%s] ", timestamp.Format(splunkDateFormat))
	}
	for _, fieldName := range fieldList {
		value := m[fieldName]
		if miscOptions.OmitIfEmpty == false || len(value) > 0 {
			key := strings.Replace(fieldName, " ", miscOptions.SpaceReplacement, -1)
			fmt.Printf("%s=\"%v\" ", key, value)
		}
	}
	fmt.Printf("\n")
}

func processLogFile(config Config, filename string){
	file, err := os.Open(filename)
	checkError("open logfile", err)
	defer file.Close()

	lineCount := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		timestamp, result, groupNames := processLogLine(config.Load, scanner.Text(), lineCount)
		var outputFields []string

		if config.OutputFields == nil || len(config.OutputFields) == 0 {
			outputFields = groupNames
		} else {
			outputFields = config.OutputFields
		}

		if len(result) == 0 {
			log.Printf("Line %d: no regexp matched - line ignored", lineCount)
		} else {
			printMap(timestamp, config.MiscOptions, result, outputFields)
		}

		lineCount++
	}

	err = scanner.Err()
	checkError("scan logfile", err)
}

func main() {
	outputLogFile, config := initialise()
	defer outputLogFile.Close()

	for _, filename := range flag.Args() {
		processLogFile(config, filename)
	}
}
