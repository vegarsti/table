package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func fixturePath(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}
	return filepath.Join(filepath.Dir(filename), fixture)
}

func loadFixture(t *testing.T, fixture string) string {
	content, err := ioutil.ReadFile(fixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestMain(m *testing.M) {
	build := exec.Command("go", "build")
	err := build.Run()
	if err != nil {
		fmt.Printf("could not build: %v", err)
		os.Exit(1)
	}
	m.Run()
	os.Remove("./table")
}
func TestTable(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		inputFile     string
		expectedFile  string
		expectedError string
	}{
		{"regular", []string{}, "examples/imdb.csv", "expected-output/imdb.txt", ""},
		{"messy", []string{}, "examples/imdb_messy.csv", "expected-output/imdb.txt", ""},
		{"tabs", []string{"--delimiter", `\t`}, "examples/tabs.tsv", "expected-output/tabs.txt", ""},
		{"semicolon", []string{"--delimiter", ";"}, "examples/imdb_semicolon.csv", "expected-output/imdb.txt", ""},
		{"arg", []string{"lol"}, "examples/tabs.tsv", "expected-output/empty.txt", "invalid argument(s)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./table", tt.args...)
			in, _ := cmd.StdinPipe()
			out, _ := cmd.StdoutPipe()
			errOut, _ := cmd.StderrPipe()
			cmd.Start()
			inputString := loadFixture(t, tt.inputFile)
			in.Write([]byte(inputString))
			in.Close()
			output, _ := ioutil.ReadAll(out)
			errorOutput, _ := ioutil.ReadAll(errOut)
			cmd.Wait()
			actualErr := string(errorOutput)
			if !strings.Contains(actualErr, tt.expectedError) {
				t.Fatalf("stderr did not contain expected string,\nactual stderr = %s, expected stderr = %s", actualErr, tt.expectedError)
			}
			actual := string(output)
			expected := loadFixture(t, tt.expectedFile)
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("actual stdout = %s, expected stdout = %s", actual, expected)
			}
		})
	}
}
