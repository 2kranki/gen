// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbMssql

import (
	"genapp/pkg/sharedData"
	"log"
	"testing"
)

func TestImportString(t *testing.T) {
	var str string

	log.Printf("dbMssql::TestImportString()..\n")
	sharedData.SetDebug(true)

	str = plug.GenImportString()
	if str != "\"github.com/denisenkom/go-mssqldb\"" {
		t.Fatalf("TestImportString() failed: %s should be \"github.com/denisenkom/go-mssqldb\"\n", str)
	}

	//t.Log(logData.String())
	t.Log("...end of dbMssql::TestImportString\n")

}
