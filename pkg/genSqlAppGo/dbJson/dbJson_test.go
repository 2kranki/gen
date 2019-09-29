// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbJson

import (
	"genapp/pkg/sharedData"
	"log"
	"testing"
	// Include the various Database Plugins so that they will register
	// with dbPlugin.
	_ "genapp/pkg/genSqlAppGo/dbMariadb"
	_ "genapp/pkg/genSqlAppGo/dbMssql"
	_ "genapp/pkg/genSqlAppGo/dbMysql"
	_ "genapp/pkg/genSqlAppGo/dbPostgres"
	_ "genapp/pkg/genSqlAppGo/dbSqlite"
)

const jsonTestPath = "../../../misc/test01/db.json.txt"


//----------------------------------------------------------------------------
//								TestTableNames
//----------------------------------------------------------------------------

func TestTableNames(t *testing.T) {
	//var err			error
	//var str			string

	log.Printf("TableNames()..\n")
	sharedData.SetDebug(true)

	// Do some form of testing

	//t.Log(logData.String())
	t.Log("TableNames: end of test\n")

}

//----------------------------------------------------------------------------
//								TestReadJsonFile
//----------------------------------------------------------------------------

func TestReadJsonFile(t *testing.T) {
	var err			error
	var keys		[]string

	t.Logf("TestReadJsonFile()..\n")
	sharedData.SetDebug(true)
	sharedData.SetMainPath(jsonTestPath)
	err = ReadJsonFile(sharedData.MainPath())
	if err != nil {
		t.Fatalf("TestReadJsonFile() Reading Main JSON failed: %s: %s\n", sharedData.MainPath(),err.Error())
	}
	if err = ValidateData(); err != nil {
		t.Fatalf("TestReadJsonFile() Validation failed: %s'\n", sharedData.MainPath())
	}

	if len(dbStruct.Tables) != 2 {
		t.Fatalf("TestReadJsonFile() failed: len(Tables) should be 2 but is '%d'\n", len(dbStruct.Name))
	}
	if len(dbStruct.Tables[0].Fields) != 8 {
		t.Fatalf("TestReadJsonFile() failed: should be 8 Tables but is %d\n", len(dbStruct.Tables[0].Fields))
	}

	keys, err = dbStruct.Tables[0].Keys()
	if err != nil {
		t.Fatalf("TestReadJsonFile() failed: keys() error: %s\n", err.Error())
	}
	t.Logf("keys: %v\n", keys)
	if len(keys) != 1 {
		t.Fatalf("TestReadJsonFile() failed: invalid keys length\n")
	}
	if keys[0] != "Num" {
		t.Fatalf("TestReadJsonFile() failed: invalid keys data\n")
	}

	//t.Log(logData.String())
	t.Logf("TestReadJsonFile: end of test\n")

}



