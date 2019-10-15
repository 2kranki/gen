// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

// Warning: We can not check the plugin portion because it causes
// circular imports. dbJson is imported into each of the plugins.

package dbJson

import (
	"genapp/pkg/sharedData"
	"log"
	"testing"
)

const jsonTestPath = "../../../misc/test01sq/db.json.txt"

//----------------------------------------------------------------------------
//								TestTableNames
//----------------------------------------------------------------------------

func TestTableNames(t *testing.T) {
	//var err			error
	//var str			string

	log.Printf("dbJson::TableNames()..\n")
	sharedData.SetDebug(true)

	// Do some form of testing

	//t.Log(logData.String())
	t.Log("dbJson::TableNames: end of test\n")

}

//----------------------------------------------------------------------------
//								TestReadJsonFile
//----------------------------------------------------------------------------

func TestReadJsonFile(t *testing.T) {
	var err error
	var keys []string

	t.Logf("dbJson::TestReadJsonFile()..\n")
	sharedData.SetDebug(true)
	sharedData.SetMainPath(jsonTestPath)
	err = ReadJsonFile(sharedData.MainPath())
	if err != nil {
		t.Fatalf("TestReadJsonFile() Reading Main JSON failed: %s: %s\n", sharedData.MainPath(), err.Error())
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
	if keys[0] != "num" {
		t.Fatalf("TestReadJsonFile() failed: invalid keys data\n")
	}

	//t.Log(logData.String())
	t.Logf("dbJson::TestReadJsonFile: end of test\n")

}
