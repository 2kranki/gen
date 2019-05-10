// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test main

package main

import (
    "./shared"
	"testing"
)

func TestSetupShared(t *testing.T) {
	var err error

	t.Logf("shared Time = %s\n", sharedData.Time())

	t.Logf("Test that sharedData.cmd is working\n")
	sharedData.SetCmd("sql")
	if "sql" != sharedData.Cmd() {
		t.Errorf("cmd() should be 'sql', but is %s\n", sharedData.Cmd())
	}

	t.Logf("Test SetupDefns  w/o exec json file\n")
	if err = SetupShared("", "cmd"); err != nil {
		t.Errorf("SetupDefns() failed: %s\n", err)
	}
	if "cmd" != sharedData.Cmd() {
		t.Errorf("ERROR - cmd should be 'cmd', but is %s\n", sharedData.Cmd())
	}

	t.Logf("Test SetupDefns w/exec json file\n")
	if err = SetupShared("misc/test01.exec.json.txt", "cmdx"); err != nil {
		t.Errorf("SetupDefns() failed: %s\n", err)
	}
	if "sqlapp" != sharedData.Cmd() {
		t.Errorf("ERROR - cmd should be 'cmd', but is %s\n", sharedData.Cmd())
	}

}
