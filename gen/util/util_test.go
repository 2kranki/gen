// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test files package

package util

import (
	"fmt"
	"testing"
)

type jsonData struct {
	Debug   bool   `json:"debug,omitempty"`
	Force   bool   `json:"force,omitempty"`
	Noop    bool   `json:"noop,omitempty"`
	Quiet   bool   `json:"quiet,omitempty"`
	Cmd     string `json:"cmd,omitempty"`
	Defines string `json:"defines,omitempty"`
	Outdir  string `json:"outdir,omitempty"`
}

func TestIsPathRegularFile(t *testing.T) {
	var path string
	var err error

	path, err = IsPathRegularFile("./util.go")
	if err != nil {
		t.Errorf("IsPathRegularFile(./files.go) failed: %s\n", err)
	}
	fmt.Println("./files.go path:", path)

	path, err = IsPathRegularFile("./xyzzy.go")
	if err == nil {
		t.Errorf("IsPathRegularFile(./xyzzy.go) should have failed!\n")
	}
	fmt.Println("./xyzzy.go path:", path)

	t.Log("\tSuccessfully completed: TestIsPathRegularFile")
}

func TestReadJson(t *testing.T) {
	var jsonOut interface{}
	var wrk interface{}
	var err error

	if jsonOut, err = ReadJsonFile("./test/test.exec.json.txt"); err != nil {
		t.Errorf("ReadJson(test.exec.json.txt) failed: %s\n", err)
	}
	m := jsonOut.(map[string]interface{})
	if wrk = m["debug"]; wrk == nil {
		t.Errorf("ReadJson(test.exec.json.txt) missing 'debug'\n")
	}
	if wrk = m["debug_not_there"]; wrk != nil {
		t.Errorf("ReadJson(test.exec.json.txt) missing 'debug'\n")
	}
	wrk = m["cmd"]
	if wrk.(string) != "sqlapp" {
		t.Errorf("ReadJson(test.exec.json.txt) missing 'cmd'\n")
	}

	t.Log("\tSuccessfully completed: TestReadJson")
}

func TestReadJsonFileToData(t *testing.T) {
	var jsonOut = jsonData{}
	var err error

	jsonOut = jsonData{}
	t.Log("&jsonOut:", &jsonOut)
	err = ReadJsonFileToData("./test/test.exec.json.txt", &jsonOut)
	if err != nil {
		t.Errorf("ReadJsonToData(test.exec.json.txt) failed: %s\n", err)
	}
	t.Log("test jsonOut:", jsonOut)
	if jsonOut.Cmd != "sqlapp" {
		t.Errorf("ReadJsonToData(test.exec.json.txt) missing or invalid 'cmd'\n")
	}
	if jsonOut.Outdir != "./test" {
		t.Errorf("ReadJson(test.exec.json.txt) missing or invalid 'outdir'\n")
	}
	t.Log("\tSuccessfully completed: TestReadJsonToData")
}
