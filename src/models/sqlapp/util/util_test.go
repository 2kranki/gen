// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test files package

package util

import (
	"fmt"
	"os"
	"os/exec"
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

func TestExecCmd(t *testing.T) {
	var err 	error
	var cmd		*ExecCmd

	t.Log("TestExecCmd()")

	cmd = NewExecCmd("")
	if cmd == nil {
		t.Errorf("NewExecCmd(\"\") failed to allocate!\n")
	}

	cmd = NewExecCmd("ls")
	if cmd == nil {
		t.Errorf("NewExecCmd(\"ls\") failed to allocate!\n")
	}
	out, err := cmd.RunWithOutput()
	if err != nil {
		t.Errorf("RunWithOutput(\"ls\") failed: %s\n", err.Error())
	}
	t.Logf("\tls output: %s", out)

	t.Log("\tend: TestExecCmd")
}

func TestFileCompare(t *testing.T) {

	t.Log("TestFileCompare()")

	src := "./util.go"
	dst := "./util.go"
	if !FileCompare(src,dst) {
		t.Errorf("FileCompare(%s,%s) failed comparison\n", src, dst)
	}

	src = "./test/test.exec.json.txt"
	dst = "./util.go"
	if FileCompare(src,dst) {
		t.Errorf("FileCompare(%s,%s) failed comparison\n", src, dst)
	}

	t.Log("\tend: TestFileCompare")
}

func TestCopyFile(t *testing.T) {
	var err error

	t.Log("TestCopyFile()")

	src := "./test/test.exec.json.txt"
	dst := "./testout.txt"
	err = CopyFile(src, dst)
	if err != nil {
		t.Errorf("CopyFile(%s,%s) failed: %s\n", src, dst, err)
	}

	if !FileCompare(src,dst) {
		t.Errorf("CopyFile(%s,%s) failed comparison\n", src, dst)
	}

	err = os.Remove(dst)

	t.Log("\tend: TestCopyFile")
}

func TestCopyDir(t *testing.T) {
	var err error

	t.Log("TestCopyDir()")

	src := "./test"
	dst := "./test2"
	err = CopyDir(src, dst)
	if err != nil {
		t.Errorf("CopyDir(%s,%s) failed: %s\n", src, dst, err)
	}

	cmd := exec.Command("diff", src, dst)
	err = cmd.Run()
	if err != nil {
		t.Errorf("CopyDir(%s,%s) comparison failed: %s\n", src, dst, err)
	}

	err = os.RemoveAll(dst)

	t.Log("\tend: TestCopyDir")
}

func TestIsPathDir(t *testing.T) {
	var path string
	var err error

	t.Log("TestIsPathDir()")

	path, err = IsPathDir("./util.go")
	if err == nil {
		t.Errorf("IsPathDir(./files.go) failed: %s\n", err)
	}
	fmt.Println("./files.go path:", path)

	path, err = IsPathDir("./test")
	if err != nil {
		t.Errorf("IsPathRegularFile(./xyzzy.go) should have failed!\n")
	}
	fmt.Println("./test path:", path)

	t.Log("\tend: TestIsPathDir")
}

func TestIsPathRegularFile(t *testing.T) {
	var input	string
	var path string
	var err error

	t.Log("\tend: TestIsPathRegularFile")

	path, err = IsPathRegularFile("./util.go")
	if err != nil {
		t.Errorf("IsPathRegularFile(./files.go) failed: %s\n", err)
	}
	fmt.Println("./files.go path:", path)

	input = "./xyzzy.go"
	path, err = IsPathRegularFile(input)
	if err == nil {
		t.Errorf("IsPathRegularFile(%s) failed: %s\n", path, err)
	}
	t.Logf("\t%s => %s\n", input, path)

	t.Log("\tend: TestIsPathRegularFile")
}

func TestPathClean(t *testing.T) {
	var err			error
	var expected	string
	var input		string
	var path 		string
	homeDir := HomeDir()
	curDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error: Getting Current Directory: %s\n", err)
	}

	t.Log("\tend: TestPathClean")

	input = "./util.go"
	expected = curDir + "/util.go"
	path = PathClean(input)
	t.Logf("\t%s => %s\n", input, path)
	if path != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", path, expected)
	}

	input = "./xyzzy.go"
	expected = curDir + "/xyzzy.go"
	path = PathClean(input)
	t.Logf("\t%s => %s\n", input, path)
	if path != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", path, expected)
	}

	input = "~"
	expected = homeDir
	path = PathClean(input)
	t.Logf("\t%s => %s\n", input, path)
	if path != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", path, expected)
	}

	input = "~/.ssh"
	expected = homeDir + "/.ssh"
	path = PathClean(input)
	t.Logf("\t%s => %s\n", input, path)
	if path != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", path, expected)
	}

	t.Log("\tend: TestPathClean")
}

func TestReadJson(t *testing.T) {
	var jsonOut interface{}
	var wrk interface{}
	var err error

	t.Log("TestReadJson()")

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

	t.Log("\tend: TestReadJson")
}

func TestReadJsonFileToData(t *testing.T) {
	var jsonOut = jsonData{}
	var err error

	t.Log("TestReadJsonFileToData()")

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
	t.Log("\tend: TestReadJsonToData")
}
