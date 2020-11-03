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
		name    string
		args    []string
		fixture string
	}{
		{"small", []string{}, "expected-output/small.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./table", tt.args...)
			in, _ := cmd.StdinPipe()
			out, _ := cmd.StdoutPipe()
			cmd.Start()
			in.Write([]byte("name,age\nVegard,27"))
			in.Close()
			output, _ := ioutil.ReadAll(out)
			cmd.Wait()
			actual := string(output)
			expected := loadFixture(t, tt.fixture)
			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}
}
