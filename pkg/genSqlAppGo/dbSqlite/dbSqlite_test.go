// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbSqlite

import (
	"genapp/pkg/sharedData"
	"log"
	"testing"
	//"time"
	"genapp/pkg/genSqlAppGo/dbJson"
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
	var str		    string
	var plg			= Plugin{}
	var dataTest	= "\tconnStr := fmt.Sprintf(\"%s\", dbName)\n\tlog.Printf(\"\\tConnecting to %s\\n\", connStr)\n\tdbSql, err = sql.Open(\"sqlite3\", connStr)\n"

	log.Printf("TestGenSqlOpen()..\n")

	sharedData.SetDebug(true)
	sharedData.SetDefn("GenDebugging", true)
	str = plg.GenSqlOpen("dbSql", "dbServer", "dbPort", "dbUser", "dbPW", "dbName")
    t.Log("===")
    t.Logf("Generated: \"%s\"\n", str)
    t.Logf("Expected:  \"%s\"\n", dataTest)
    t.Log("===")
	if len(str) != len(dataTest) {
		t.Fatalf("TestGenSqlOpen() Invalid generation: length was %d vs %d\n", len(str), len(dataTest))
	}
    if str != dataTest {
        t.Errorf(" generated: %s\n", str)
        t.Errorf("expected: %s\n", dataTest)
        t.Fatalf("TestGenSqlOpen() generated data did not match saved data!\n")
	}

	t.Log("TestGenSqlOpen: end of test\n")
}


