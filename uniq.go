package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	options, err := validateArguments(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	strings := scanLines(options.input)

	ignoreCase := options.ignoreCase

	var transformedStrings []string
	if options.ignoreFieldsCount != 0 || options.ignoreCharsCount != 0 {
		transformedStrings = make([]string, len(strings))
		for idx, v := range strings {
			transformedStrings[idx] = transformString(v, options.ignoreFieldsCount, options.ignoreCharsCount)
		}
	}

	result := ""
	switch options.outputMode {
	case normalParam:
		result = withoutParamsAlgo(strings, transformedStrings, ignoreCase)
	case cParam:
		result = cParamAlgo(strings, transformedStrings, ignoreCase)
	case dParam:
		result = dParamAlgo(strings, transformedStrings, ignoreCase)
	case uParam:
		result = uParamAlgo(strings, transformedStrings, ignoreCase)
	}

	_, err = io.WriteString(options.output, result)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	description string = "uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]"

	normalParam     = 0
	cParam      int = 1
	dParam      int = 2
	uParam      int = 3
)

type options struct {
	outputMode        int  // -c | -d | -u
	ignoreCase        bool // -i
	ignoreFieldsCount int  // -f num
	ignoreCharsCount  int  // -s chars
	input             io.Reader
	output            io.Writer
}

func validateArguments(args []string) (options, error) {
	var flags flag.FlagSet

	cFlag := flags.Bool("c", false, "")
	dFlag := flags.Bool("d", false, "")
	uFlag := flags.Bool("u", false, "")
	iFlag := flags.Bool("i", false, "")
	fFlag := flags.Int("f", 0, "")
	sFlag := flags.Int("s", 0, "")
	flags.Parse(args)

	opts := options{
		normalParam,
		false,
		0,
		0,
		os.Stdin,
		os.Stdout,
	}

	if *cFlag {
		if *dFlag || *uFlag {
			return opts, errors.New(description)
		}
		opts.outputMode = cParam
	}
	if *dFlag {
		if *cFlag || *uFlag {
			return opts, errors.New(description)
		}
		opts.outputMode = dParam
	}
	opts.ignoreCase = *iFlag

	if *fFlag < 0 {
		return opts, errors.New("-f value should be > 0")
	}
	opts.ignoreFieldsCount = *fFlag

	if *sFlag < 0 {
		return opts, errors.New("-s value should be > 0")
	}
	opts.ignoreCharsCount = *sFlag

	if *uFlag {
		if *dFlag || *cFlag {
			return opts, errors.New(description)
		}
		opts.outputMode = uParam
	}

	for idx := 0; idx < len(args); idx++ {
		_, err := strconv.Atoi(args[idx])
		if args[idx][0] != '-' && err != nil {
			opts.input, err = os.Open(args[idx])
			if err != nil {
				return opts, err
			}
			idx++
			if idx != len(args) {
				file, err := os.Create(args[idx])
				if err != nil {
					return opts, err
				}
				opts.output = file
			}
		}
	}

	return opts, nil
}

func scanLines(reader io.Reader) []string {
	var result []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result
}

func transformString(str string, ignoreFieldsCount int, ignoreCharsCount int) string {
	if len(str) == 0 {
		return str
	}

	splittedString := strings.Fields(str)[ignoreFieldsCount:]
	newStr := ""
	for idx := range splittedString {
		newStr = fmt.Sprintf("%s %s", newStr, splittedString[idx])
	}

	if len(newStr) <= ignoreCharsCount+1 {
		return ""
	}

	return newStr[ignoreCharsCount+1:]
}

func withoutParamsAlgo(strings []string, transformedStrings []string, ignoreCase bool) (result string) {
	entries := mainAlgo(strings, transformedStrings, ignoreCase)
	for _, v := range entries {
		result += fmt.Sprintf("%s\n", *(v.Key))
	}

	return
}

func cParamAlgo(strings []string, transformedStrings []string, ignoreCase bool) (result string) {
	entries := mainAlgo(strings, transformedStrings, ignoreCase)
	for _, v := range entries {
		result += fmt.Sprintf("%d %s\n", v.Value, *(v.Key))
	}

	return
}

func dParamAlgo(strings []string, transformedStrings []string, ignoreCase bool) (result string) {
	entries := mainAlgo(strings, transformedStrings, ignoreCase)
	for _, v := range entries {
		if v.Value >= 2 {
			result += fmt.Sprintf("%s\n", *(v.Key))
		}
	}

	return
}

func uParamAlgo(strings []string, transformedStrings []string, ignoreCase bool) (result string) {
	entries := mainAlgo(strings, transformedStrings, ignoreCase)
	for _, v := range entries {
		if v.Value == 1 {
			result += fmt.Sprintf("%s\n", *(v.Key))
		}
	}

	return
}

type stringEntries struct {
	Key   *string
	Value int
}

func mainAlgo(strings []string, transformedStrings []string, ignoreCase bool) []*stringEntries {
	stringsBackup := strings
	if transformedStrings != nil {
		strings = transformedStrings
	}

	result := make([]*stringEntries, 0, 10)
	j := 0
	stringEntries_ := &stringEntries{
		&(stringsBackup)[j],
		1,
	}
	result = append(result, stringEntries_)
	for i := 1; i < len(strings); i++ {
		if StringsEquals((strings)[j], (strings)[i], ignoreCase) {
			stringEntries_.Value += 1
			continue
		}
		j++
		(strings)[j] = (strings)[i]
		(stringsBackup)[j] = (stringsBackup)[i]

		stringEntries_ = &stringEntries{
			&(stringsBackup)[j],
			1,
		}
		result = append(result, stringEntries_)
	}
	strings = strings[:j+1]

	return result
}

func StringsEquals(str1, str2 string, ignoreCase bool) bool {
	if ignoreCase {
		return strings.ToLower(str1) == strings.ToLower(str2)
	}
	return str1 == str2
}
