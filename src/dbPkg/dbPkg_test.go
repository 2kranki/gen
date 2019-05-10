// See License.txt in main repository directory

// Template Functions used in generation
// See the specific template files for how the functions
// and data are used.

package dbPkg

import (
	"../shared"
	"log"
	"testing"
)

var fld0sql = "\tCustNo\tINT NOT NULL PRIMARY KEY,\n"
var fld0struct = "\tCustNo\tint64\n"
var fld1sql = "\tCustName\tNVARCHAR(30),\n"
var fld1struct = "\tCustName\tstring\n"
var tbl0sql = "DROP TABLE Customer IF EXISTS;\nCREATE TABLE Customer (\n\tCustNo\tINT NOT NULL PRIMARY KEY,\n\tCustName\tNVARCHAR(30),\n\tCustAddr1\tNVARCHAR(30),\n\tCustAddr2\tNVARCHAR(30),\n\tCustCity\tNVARCHAR(20),\n\tCustState\tNVARCHAR(10),\n\tCustZip\tNVARCHAR(15),\n\tCustCurBal\tDEC(15,2)\n);\n"
var tbl0struct = "type Customer struct {\n\tCustNo\tint64\n\tCustName\tstring\n\tCustAddr1\tstring\n\tCustAddr2\tstring\n\tCustCity\tstring\n\tCustState\tstring\n\tCustZip\tstring\n\tCustCurBal\tfloat64\n}\n"

func TestCreate(t *testing.T) {
	var err			error
	var str			string

	log.Printf("TestCreate()..\n")
	sharedData.SetDebug(true)
	sharedData.SetMainPath("../misc/test01/db.json.txt")
	if err = ReadJsonFile(sharedData.MainPath()); err != nil {
		t.Fatalf("TestCreate() Reading Main JSON failed: %s'\n", sharedData.MainPath())
	}
	if err = ValidateData(); err != nil {
		t.Fatalf("TestCreate() Validation failed: %s'\n", sharedData.MainPath())
	}

	if len(dbStruct.Tables) != 2 {
		t.Fatalf("TestCreate() failed: len(Tables) should be 2 but is '%d'\n", len(dbStruct.Name))
	}
	if len(dbStruct.Tables[0].Fields) != 8 {
		t.Fatalf("TestCreate() failed: should be 8 Tables but is %d\n", len(dbStruct.Tables[0].Fields))
	}

	str = dbStruct.Tables[0].Fields[0].CreateSql(",")
	t.Log("Table[0].Fields[0] CreateSql =", str)
	if str != fld0sql {
		t.Fatalf("TestCreate() failed: invalid create sql generated Tables[0].Fields[0]\n")
	}

	str = dbStruct.Tables[0].Fields[1].CreateSql(",")
	t.Log("Table[0].Fields[1] CreateSql =", str)
	if str != fld1sql {
		t.Fatalf("TestCreate() failed: invalid create sql generated Tables[0].Fields[1]\n")
	}

	str = dbStruct.Tables[0].CreateSql()
	t.Log("Table[0] CreateSql =", str)
	if str != tbl0sql {
		t.Fatalf("TestCreate() failed: invalid create sql generated\n")
	}

	str = dbStruct.Tables[0].Fields[0].CreateStruct()
	t.Log("Table[0].Field[0] Struct =", str)
	if str != fld0struct {
		t.Fatalf("TestCreate() failed: invalid struct generated\n")
	}

	str = dbStruct.Tables[0].Fields[1].CreateStruct()
	t.Log("Table[0].Field[1] Struct =", str)
	if str != fld1struct {
		t.Fatalf("TestCreate() failed: invalid struct generated\n")
	}

	str = dbStruct.Tables[0].CreateStruct()
	t.Log("Table[0] Struct =", str)
	if str != tbl0struct {
		t.Fatalf("TestCreate() failed: invalid struct generated\n")
	}

	//t.Log(logData.String())
	t.Log("TestCreate: end of test\n")

}

func TestReadJsonFileDb(t *testing.T) {
	var err			error

	t.Logf("TestReadJsonFileDb()..\n")
	sharedData.SetDebug(true)
	sharedData.SetMainPath("../misc/test01/db.json.txt")
	if err = ReadJsonFile(sharedData.MainPath()); err != nil {
		t.Fatalf("ReadJsonFileDb() Reading Main JSON failed: %s'\n", sharedData.MainPath())
	}

	if len(dbStruct.Tables) != 2 {
		t.Fatalf("ReadJsonFileDb() failed: len(Tables) should be 2 but is '%d'\n", len(dbStruct.Name))
	}
	if len(dbStruct.Tables[0].Fields) != 8 {
		t.Fatalf("ReadJsonFileDb() failed: should be 3 Fields but is %d\n", len(dbStruct.Tables[0].Fields))
	}
	if len(dbStruct.Tables[1].Fields) != 8 {
		t.Fatalf("ReadJsonFileDb() failed: should be 8 Fields but is %d\n", len(dbStruct.Tables[1].Fields))
	}

	//t.Log(logData.String())
	t.Log("TestReadJsonFileDb(): end of test")

}
