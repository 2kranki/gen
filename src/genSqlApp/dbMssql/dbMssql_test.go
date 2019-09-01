// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbMssql

import (
	"genapp/shared"
	"log"
	"testing"
)

func TestImportString(t *testing.T) {
	var str			string

	log.Printf("TestImportString()..\n")
	sharedData.SetDebug(true)

	str = ImportString()
	if str != "\"github.com/denisenkom/go-mssqldb\"" {
		t.Fatalf("TestImportString() failed: %s should be \"github.com/denisenkom/go-mssqldb\"\n", str)
	}

	//t.Log(logData.String())
	t.Log("TestImportString: end of test\n")

}

