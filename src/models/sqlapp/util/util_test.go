// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test files package

package util

import (
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

	src := NewPath("./util.go")
	dst := NewPath("./util.go")
	if !FileCompareEqual(src, dst) {
		t.Errorf("FileCompare(%s,%s) failed comparison\n", src, dst)
	}

	src = NewPath("./test/test.exec.json.txt")
	dst = NewPath("./util.go")
	if FileCompareEqual(src, dst) {
		t.Errorf("FileCompare(%s,%s) failed comparison\n", src, dst)
	}

	t.Log("\tend: TestFileCompare")
}

func TestCopyFile(t *testing.T) {
	var err 	error

	t.Log("TestCopyFile()")

	src := NewPath("test").Append("test.exec.json.txt")
	dst := NewTempDir().Append("testout.txt")
	err = CopyFile(src, dst)
	if err != nil {
		t.Errorf("CopyFile(%s,%s) failed: %s\n", src.String(), dst.String(), err.Error())
	}

	if !FileCompareEqual(src,dst) {
		t.Errorf("CopyFile(%s,%s) failed comparison\n", src.String(), dst.String())
	}

	err = dst.DeleteFile()
	if err != nil {
		t.Errorf("DeleteFile(%s) failed: %s\n", dst.String(), err.Error())
	}

	t.Log("\tend: TestCopyFile")
}

func TestCopyDir(t *testing.T) {
	var err error

	t.Log("TestCopyDir()")

	src  := NewPath("test")
	dst  := NewTempDir().Append("test2")
	dst2 := NewTempDir().Append("test3")

	err = dst.RemoveDir()
	if err != nil {
		t.Logf("\tError: Deleting %s: %s\n", dst.String(), err.Error())
	}
	err = dst2.RemoveDir()
	if err != nil {
		t.Logf("\tError: Deleting %s: %s\n", dst2.String(), err.Error())
	}

	t.Logf("\tcopying %s -> %s\n", src, dst)
	err = CopyDir(src, dst)
	if err != nil {
		t.Fatalf("CopyDir(%s,%s) failed: %s\n", src.String(), dst.String(), err.Error())
	}

	cmd := exec.Command("diff", src.Absolute(), dst.Absolute())
	err = cmd.Run()
	if err != nil {
		t.Fatalf("CopyDir(%s,%s) comparison failed: %s\n", src, dst, err)
	}

	dst.RemoveDir()

	dst3 := dst2.Append("test")
	dst2 =  dst2.Append("")
	t.Logf("\tcopying %s -> %s\n", src.String(), dst2.String())
	err = CopyDir(src, dst3)
	if err != nil {
		t.Fatalf("CopyDir(%s,%s) failed: %s\n", src.String(), dst3.String(), err.Error())
	}

	cmd = exec.Command("diff", src.Absolute(), dst3.Absolute())
	err = cmd.Run()
	if err != nil {
		t.Fatalf("CopyDir(%s,%s) comparison failed: %s\n", src.String(), dst3.String(), err)
	}

	dst2.RemoveDir()

	t.Log("\tend: TestCopyDir")
}

func TestIsPathDir(t *testing.T) {
	var path	*Path

	t.Log("TestIsPathDir()")

	path = NewPath("./util.go")
	if path.IsPathDir() {
		t.Errorf("IsPathDir(%s) failed!\n", path.String())
	}
	t.Logf("\t%s absolute: %s\n", path.String(), path.Absolute())

	path = NewPath("./test")
	if !path.IsPathDir() {
		t.Errorf("IsPathDir(%s) failed!\n", path.String())
	}
	t.Logf("\t%s absolute: %s\n", path.String(), path.Absolute())

	t.Log("\tend: TestIsPathDir")
}

func TestIsPathRegularFile(t *testing.T) {
	var path 	*Path
	var err 	error

	t.Log("TestIsPathRegularFile()")

	path = NewPath("./util.go")
	if !path.IsPathRegularFile() {
		t.Errorf("IsPathRegularFile(%s) failed: %s\n", path.String(), err.Error())
	}
	t.Logf("\t%s Absolute: %s\n", path.String(), path.Absolute())

	path = NewPath("./xyzzy.go")
	if path.IsPathRegularFile() {
		t.Errorf("IsPathRegularFile(%s) failed: %s\n", path.String(), err.Error())
	}
	t.Logf("\t%s Absolute: %s\n", path.String(), path.Absolute())

	t.Log("\tend: TestIsPathRegularFile")
}

func TestPath(t *testing.T) {
	var err			error
	var expected	string
	var input		string
	var path 		*Path
	var pth			string
	homeDir := NewHomeDir()
	curDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error: Getting Current Directory: %s\n", err)
	}

	t.Log("TestPathClean()")

	input = "./util.go"
	expected = curDir + "/util.go"
	path = NewPath(input)
	pth = path.Clean()
	t.Logf("\t%s => %s\n", input, pth)
	if pth != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", pth, expected)
	}

	input = "./xyzzy.go"
	expected = curDir + "/xyzzy.go"
	path = NewPath(input)
	pth = path.Clean()
	t.Logf("\t%s => %s\n", input, pth)
	if pth != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", pth, expected)
	}

	input = "~"
	expected = homeDir.String()
	path = NewPath(input)
	path.Clean()
	pth = path.Absolute()
	t.Logf("\t%s => %s\n", input, pth)
	if pth != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", pth, expected)
	}

	input = "~/.ssh"
	expected = homeDir.String() + "/.ssh"
	path = NewPath(input)
	pth = path.Clean()
	t.Logf("\t%s => %s\n", input, pth)
	if pth != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", pth, expected)
	}

	input = "./test3"
	path = NewPath(input)
	if err = path.CreateDir(); err != nil {
		t.Fatalf("FATAL: create ./test3 failed: %s\n", err.Error())
	}
	if !path.IsPathDir() {
		t.Fatalf("FATAL: create ./test3 failed!\n")
	}
	if err = path.RemoveDir(); err != nil {
		t.Fatalf("FATAL: remove ./test3 failed: %s\n", err.Error())
	}
	if path.IsPathDir() {
		t.Fatalf("FATAL: remove ./test3 failed!\n")
	}

	t.Logf("\t%s => %s\n", input, pth)
	if pth != expected {
		t.Errorf("PathClean Got: %s  Expected: %s\n", pth, expected)
	}

	pwd := NewWorkDir()
	t.Logf("PWD: %s\n", pwd.String())

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
