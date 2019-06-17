// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbSqlite

import (
	"log"
	"testing"
	"../../shared"
	//"time"
	"../dbJson"
)

const jsonTestPath = "../../../misc/test01/db.json.txt"

var fld0sql = "\tCustNo\tINT NOT NULL PRIMARY KEY,\n"
var fld0struct = "\tCustNo\tint64\n"
var fld1sql = "\tCustName\tNVARCHAR(30),\n"
var fld1struct = "\tCustName\tstring\n"
var tbl0sql = "DROP TABLE Customer IF EXISTS;\nCREATE TABLE Customer (\n\tCustNo\tINT NOT NULL PRIMARY KEY,\n\tCustName\tNVARCHAR(30),\n\tCustAddr1\tNVARCHAR(30),\n\tCustAddr2\tNVARCHAR(30),\n\tCustCity\tNVARCHAR(20),\n\tCustState\tNVARCHAR(10),\n\tCustZip\tNVARCHAR(15),\n\tCustCurBal\tDEC(15,2)\n);\n"
var tbl0struct = "type Customer struct {\n\tCustNo\tint64\n\tCustName\tstring\n\tCustAddr1\tstring\n\tCustAddr2\tstring\n\tCustCity\tstring\n\tCustState\tstring\n\tCustZip\tstring\n\tCustCurBal\tfloat64\n}\n"

func TestCreate(t *testing.T) {
	var err			error

	log.Printf("TestCreate()..\n")
	sharedData.SetDebug(true)
	sharedData.SetMainPath(jsonTestPath)
	if err = dbJson.ReadJsonFile(sharedData.MainPath()); err != nil {
		t.Fatalf("TestCreate() Reading Main JSON failed: %s\n", sharedData.MainPath())
	}
	if err = dbJson.SetupPlugin(*pluginData); err != nil {
		t.Fatalf("TestCreate() SetupPlugin failed: %s\n", err)
	}

	t.Log("TestCreate: end of test\n")
}

func TestGenSqlOpen(t *testing.T) {
	var strs		[]string
	var plg			= Plugin{}
	var dataTest	= []string{
		"\t// dbName is a CLI argument\n",
		"\tconnStr := fmt.Sprintf(\"%s\", dbName)\n",
		"\tlog.Printf(\"\\tConnecting to %s\\n\", connStr)\n",
		"\tdb, err = sql.Open(\"sqlite3\", connStr)\n",
	}

	log.Printf("TestGenSqlOpen()..\n")

	sharedData.SetDebug(true)
	sharedData.SetDefn("GenDebugging", true)
	strs = plg.GenSqlOpen()
	if len(strs) != 4 {
		t.Fatalf("TestGenSqlOpen() Invalid generation: should be 4 lines but was %d\n", len(strs))
	}
	for i :=0; i < len(strs); i++ {
		if strs[i] != dataTest[i] {
			t.Errorf(" gen: %s", strs[i])
			t.Errorf("data: %s", dataTest[i])
			t.Fatalf("TestGenSqlOpen() generated data did not match saved data - line %d\n", i)
		}
	}

	t.Log("TestGenSqlOpen: end of test\n")
}


