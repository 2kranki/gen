// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbPostgres

import (
	"genapp/pkg/sharedData"
	"log"
	"testing"
)

func TestImportString(t *testing.T) {
	var str string

	log.Printf("dbPostgres::TestImportString()..\n")
	sharedData.SetDebug(true)

	str = plug.GenImportString()
	if str != "\"github.com/lib/pq\"" {
		t.Fatalf("TestImportString() failed: %s should be \"github.com/lib/pq\"\n", str)
	}

	//t.Log(logData.String())
	t.Log("...end of dbPostgres::TestImportString\n")

}
