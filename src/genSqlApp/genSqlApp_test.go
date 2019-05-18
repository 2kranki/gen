// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package genSqlApp

import (
	"log"
	"testing"
	"time"
	"../dbPkg"
	"../shared"
)

const jsonTestPath = "../misc/"

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

func TestReadJsonFiles(t *testing.T) {
	var data        *dbPkg.Database
    var err         error

	sharedData.SetDebug(true)
	sharedData.SetForce(false)
	sharedData.SetNoop(true)
	sharedData.SetQuiet(false)
	sharedData.SetMdlDir("../models/")
	sharedData.SetOutDir("/tmp/gen")
	sharedData.SetTime(time.Now().Format("Mon Jan _2, 2006 15:04"))
    sharedData.SetFunc("Time", sharedData.Time)
    sharedData.SetDataPath("../misc/test01/app.json.txt")
    sharedData.SetMainPath("../misc/test01/main.json.txt")

	if err = readJsonFiles(); err != nil {
		t.Errorf("ReadJsonFile() failed: %s\n", err)
	}
    data = appData.AppStruct()

	if data.Name != "Finances" {
		t.Errorf("ReadJsonFile() failed: Name should be 'Finances' but is '%s'\n", data.Name)
	}
	if len(data.Tables) != 2 {
		t.Errorf("ReadJsonFile() failed: should be 2 tables but is %d\n", len(data.Tables))
	}
	if len(data.Tables[0].Fields) != 8 {
		t.Errorf("ReadJsonFile() failed: should be 8 fields in table[0] but is %d\n", len(data.Tables[0].Fields))
	}

	t.Log("Successfully, completed: TestReadJsonFiles")
}

func TestCreateModelPath(t *testing.T) {
	var name        string
	var name2       string
    var err         error

	sharedData.SetDebug(true)
	sharedData.SetForce(false)
	sharedData.SetNoop(true)
	sharedData.SetQuiet(false)
	sharedData.SetMdlDir("../models/")
	sharedData.SetOutDir("/tmp/gen")
	sharedData.SetTime(time.Now().Format("Mon Jan _2, 2006 15:04"))
    sharedData.SetFunc("Time", sharedData.Time)
    sharedData.SetDataPath("../misc/test01/app.json.txt")
    sharedData.SetMainPath("../misc/test01/main.json.txt")

	if name, err = createModelPath("../models/sqlapp/tableio.go.tmpl.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t../Models/sqlapp/tableio.go.tmpl.txt -> '%s'\n", name)

	if name2, err = createModelPath("tableio.go.tmpl.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\ttableio.go.tmpl.txt -> '%s'\n", name2)

    if name != name2 {
		t.Errorf("createModelPath() file names don't match!\n")
	}

	t.Log("Successfully, completed: TestCreateModelPath")
}

func TestCreateOutputPath(t *testing.T) {
	var name        string
    var err         error

	sharedData.SetDebug(true)
	sharedData.SetForce(false)
	sharedData.SetNoop(true)
	sharedData.SetQuiet(false)
	sharedData.SetMdlDir("../models/")
	sharedData.SetOutDir("/tmp/gen")
	sharedData.SetTime(time.Now().Format("Mon Jan _2, 2006 15:04"))
    sharedData.SetFunc("Time", sharedData.Time)
    sharedData.SetDataPath("../misc/test01/app.json.txt")
    sharedData.SetMainPath("../misc/test01/main.json.txt")

	if name, err = createOutputPath("tableio.go"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\ttableio.go -> '%s'\n", name)
    if name != "/tmp/gen/tableio.go" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	t.Log("Successfully, completed: TestCreateOutputPath")
}


