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

	options, err := ValidateArguments(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	strings := ScanLines(options.input)

	ignoreCase := options.ignoreRegister

	var transformedStrings []string = nil
	if options.ignoreFieldsCount != 0 || options.ignoreCharsCount != 0 {
		transformedStrings = make([]string, len(strings))
		for idx, v := range strings {
			transformedStrings[idx] = transformString(v, options.ignoreFieldsCount, options.ignoreCharsCount)
		}
	}

	result := ""
	switch options.outputMode {
	case normalParam:
		result = WithoutParams(strings, transformedStrings, ignoreCase)
	case cParam:
		result = C(strings, transformedStrings, ignoreCase)
	case dParam:
		result = D(strings, transformedStrings, ignoreCase)
	case uParam:
		result = U(strings, transformedStrings, ignoreCase)
	}

	io.WriteString(options.output, result)
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
	ignoreRegister    bool // -i
	ignoreFieldsCount int  // -f num
	ignoreCharsCount  int  // -s chars
	input             io.Reader
	output            io.Writer
}

func ValidateArguments(args []string) (options, error) {
	var flags flag.FlagSet

	cFlag := flags.Bool("c", false, "")
	dFlag := flags.Bool("d", false, "")
	uFlag := flags.Bool("u", false, "")
	iFlag := flags.Bool("i", false, "")
	fFlag := flags.Int("f", 0, "")
	sFlag := flags.Int("s", 0, "")
	flags.Parse(args)

	var opts options = options{
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
	opts.ignoreRegister = *iFlag

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
			var err error
			opts.input, err = os.Open(args[idx])
			if err != nil {
				return opts, err
			}
			idx++
			if idx != len(args) {
				var file, err = os.Create(args[idx])
				if err != nil {
					return opts, err
				}
				opts.output = file
			}
		}
	}

	return opts, nil
}

func ScanLines(reader io.Reader) []string {
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

func WithoutParams(strings []string, transformedStrings []string, ignoreCase bool) string {
	entries := MainAlgo(strings, transformedStrings, ignoreCase)
	result := ""
	for _, v := range entries {
		result += *(v.Key) + "\n"
	}

	return result
}

func C(strings []string, transformedStrings []string, ignoreCase bool) string {
	entries := MainAlgo(strings, transformedStrings, ignoreCase)
	result := ""
	for _, v := range entries {
		result += strconv.Itoa(v.Value) + " " + *(v.Key) + "\n"
	}

	return result
}

func D(strings []string, transformedStrings []string, ignoreCase bool) string {
	entries := MainAlgo(strings, transformedStrings, ignoreCase)
	result := ""
	for _, v := range entries {
		if v.Value >= 2 {
			result += *(v.Key) + "\n"
		}
	}

	return result
}

func U(strings []string, transformedStrings []string, ignoreCase bool) string {
	entries := MainAlgo(strings, transformedStrings, ignoreCase)
	result := ""
	for _, v := range entries {
		if v.Value == 1 {
			result += *(v.Key) + "\n"
		}
	}

	return result
}

type StringEntries struct {
	Key   *string
	Value int
}

func MainAlgo(strings []string, transformedStrings []string, ignoreCase bool) []*StringEntries {
	stringsBackup := strings
	if transformedStrings != nil {
		strings = transformedStrings
	}

	result := make([]*StringEntries, 0, 10)
	j := 0
	stringEntries := &StringEntries{
		&(stringsBackup)[j],
		1,
	}
	result = append(result, stringEntries)
	for i := 1; i < len(strings); i++ {
		if !StringsEquals((strings)[j], (strings)[i], ignoreCase) {
			//if (*in)[j] != (*in)[i] {
			j++
			(strings)[j] = (strings)[i]
			(stringsBackup)[j] = (stringsBackup)[i]

			stringEntries = &StringEntries{
				&(stringsBackup)[j],
				1,
			}
			result = append(result, stringEntries)
		} else {
			stringEntries.Value += 1
		}
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
