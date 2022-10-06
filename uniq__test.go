package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestValidateArguments(t *testing.T) {
	type testData struct {
		args           []string
		expectedResult options
	}

	tests := []testData{
		{[]string{}, options{
			normalParam,
			false,
			0,
			0,
			os.Stdin,
			os.Stdout}},
		{[]string{"-c"}, options{
			cParam,
			false,
			0,
			0,
			os.Stdin,
			os.Stdout}},
		{[]string{"-d"}, options{
			dParam,
			false,
			0,
			0,
			os.Stdin,
			os.Stdout}},
		{[]string{"-u"}, options{
			uParam,
			false,
			0,
			0,
			os.Stdin,
			os.Stdout}},
		{[]string{"-i"}, options{
			normalParam,
			true,
			0,
			0,
			os.Stdin,
			os.Stdout}},
		{[]string{"-f", "3"}, options{
			normalParam,
			false,
			3,
			0,
			os.Stdin,
			os.Stdout}},
		{[]string{"-s", "5"}, options{
			normalParam,
			false,
			0,
			5,
			os.Stdin,
			os.Stdout}},
	}

	for _, v := range tests {
		opts, err := ValidateArguments(v.args)
		require.Equal(t, v.expectedResult, opts, err)
	}
}

func TestValidateArgumentsIncorrect(t *testing.T) {
	description := errors.New(description)

	type testData struct {
		flags   []string
		error   error
		message string
	}

	tests := []testData{
		{[]string{"-c", "-d"}, description, "This options can not use both"},
		{[]string{"-d", "-u"}, description, "This options can not use both"},
		{[]string{"-d", "-c"}, description, "This options can not use both"},
		{[]string{"-u", "-d"}, description, "This options can not use both"},
		{[]string{"-f", "-4"}, errors.New("-f value should be > 0"), "-f value should be > 0"},
		{[]string{"-s", "-4"}, errors.New("-s value should be > 0"), "-s value should be > 0"},
	}

	for _, v := range tests {
		_, err := ValidateArguments(v.flags)
		if errors.Is(err, v.error) {
			t.Errorf(v.message)
		}
	}
}

func TestWithoutParams(t *testing.T) {
	type testData struct {
		strings        []string
		expectedResult string
		message        string
	}

	tests := []testData{
		{[]string{"I love music.", "I love music.", "I love music.", "", "I love music of Kartik.", "I love music of Kartik.",
			"Thanks.", "I love music of Kartik.", "I love music of Kartik."},
			`I love music.

I love music of Kartik.
Thanks.
I love music of Kartik.
`, "Failed test withoutParams"},
	}

	for _, v := range tests {
		result := WithoutParams(v.strings, nil, false)
		if result != v.expectedResult {
			t.Errorf(v.message)
		}
	}
}

func TestI(t *testing.T) {
	type testData struct {
		strings        []string
		expectedResult string
		message        string
	}

	tests := []testData{
		{[]string{"I love music.", "I love music.", "I love music.", "", "I love music of Kartik.", "I love music of Kartik.",
			"Thanks.", "I love music of Kartik.", "I love music of Kartik."},
			`I love music.

I love music of Kartik.
Thanks.
I love music of Kartik.
`, "Failed test withoutParams"},
	}

	for _, v := range tests {
		result := WithoutParams(v.strings, nil, true)
		if result != v.expectedResult {
			t.Errorf(v.message)
		}
	}
}

func TestC(t *testing.T) {
	type testData struct {
		strings        []string
		expectedResult string
		message        string
	}

	tests := []testData{
		{[]string{"I love music.", "I love music.", "I love music.", "", "I love music of Kartik.", "I love music of Kartik.",
			"Thanks.", "I love music of Kartik.", "I love music of Kartik."},
			`3 I love music.
1 
2 I love music of Kartik.
1 Thanks.
2 I love music of Kartik.
`, "Failed test withoutParams"},
	}

	for _, v := range tests {
		result := C(v.strings, nil, true)
		if result != v.expectedResult {
			t.Errorf(v.message)
		}
	}
}

func TestD(t *testing.T) {
	type testData struct {
		strings        []string
		expectedResult string
		message        string
	}

	tests := []testData{
		{[]string{"I love music.", "I love music.", "I love music.", "", "I love music of Kartik.", "I love music of Kartik.",
			"Thanks.", "I love music of Kartik.", "I love music of Kartik."},
			`I love music.
I love music of Kartik.
I love music of Kartik.
`, "Failed test withoutParams"},
	}

	for _, v := range tests {
		result := D(v.strings, nil, true)
		if result != v.expectedResult {
			t.Errorf(v.message)
		}
	}
}

func TestU(t *testing.T) {
	type testData struct {
		strings        []string
		expectedResult string
		message        string
	}

	tests := []testData{
		{[]string{"I love music.", "I love music.", "I love music.", "", "I love music of Kartik.", "I love music of Kartik.",
			"Thanks.", "I love music of Kartik.", "I love music of Kartik."},
			`
Thanks.
`, "Failed test withoutParams"},
	}

	for _, v := range tests {
		result := U(v.strings, nil, true)
		if result != v.expectedResult {
			t.Errorf(v.message)
		}
	}
}
