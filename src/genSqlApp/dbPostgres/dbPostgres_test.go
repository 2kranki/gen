// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbPostgres

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
	if str != "\"github.com/lib/pq\"" {
		t.Fatalf("TestImportString() failed: %s should be \"github.com/lib/pq\"\n", str)
	}

	//t.Log(logData.String())
	t.Log("TestImportString: end of test\n")

}



