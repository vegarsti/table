package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
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
		name         string
		args         []string
		inputFile    string
		expectedFile string
	}{
		{"regular", []string{}, "examples/imdb.csv", "expected-output/imdb.txt"},
		{"messy", []string{}, "examples/imdb_messy.csv", "expected-output/imdb.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./table", tt.args...)
			in, _ := cmd.StdinPipe()
			out, _ := cmd.StdoutPipe()
			cmd.Start()
			inputString := loadFixture(t, tt.inputFile)
			in.Write([]byte(inputString))
			in.Close()
			output, _ := ioutil.ReadAll(out)
			cmd.Wait()
			actual := string(output)
			expected := loadFixture(t, tt.expectedFile)
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}
}
