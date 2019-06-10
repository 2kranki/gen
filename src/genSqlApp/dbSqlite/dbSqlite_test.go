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
	var str			string

	log.Printf("TestCreate()..\n")
	sharedData.SetDebug(true)
	sharedData.SetMainPath("../misc/test01/db.json.txt")
	if err = dbJson.ReadJsonFile(sharedData.MainPath()); err != nil {
		t.Fatalf("TestCreate() Reading Main JSON failed: %s'\n", sharedData.MainPath())
	}
	if err = dbJson.ValidateData(); err != nil {
		t.Fatalf("TestCreate() Validation failed: %s'\n", sharedData.MainPath())
	}

	//t.Log(logData.String())
	t.Log("TestCreate: end of test\n")

}


