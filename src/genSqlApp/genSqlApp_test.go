// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package genSqlApp

import (
	"testing"
	"time"
	"./dbJson"
	"../shared"
)

func TestReadJsonFiles(t *testing.T) {
	var data        *dbJson.Database
    var err         error

	sharedData.SetDebug(true)
	sharedData.SetForce(false)
	sharedData.SetNoop(true)
	sharedData.SetQuiet(false)
	sharedData.SetMdlDir("../models/")
	sharedData.SetOutDir("/tmp/gen")
	sharedData.SetTime(time.Now().Format("Mon Jan _2, 2006 15:04"))
    sharedData.SetFunc("Time", sharedData.Time)
    sharedData.SetDataPath("./test/db.json.txt")
    sharedData.SetMainPath("./test/main.json.txt")

	if err = readJsonFiles(); err != nil {
		t.Errorf("ReadJsonFile() failed: %s\n", err)
	}

    data = dbJson.DbStruct()

	if data.Name != "MovieDB" {
		t.Errorf("ReadJsonFile() failed: Name should be 'MovieDB' but is '%s'\n", data.Name)
	}
	if len(data.Tables) != 1 {
		t.Errorf("ReadJsonFile() failed: should be 1 tables but is %d\n", len(data.Tables))
	}
	if len(data.Tables[0].Fields) != 10 {
		t.Errorf("ReadJsonFile() failed: should be 10 fields in table[0] but is %d\n", len(data.Tables[0].Fields))
	}

	t.Log("...End of TestReadJsonFiles")
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
    sharedData.SetDataPath("./test/db.json.txt")
    sharedData.SetMainPath("./test/main.json.txt")

	if name, err = createModelPath("../sqlapp/io.table.go.tmpl.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t../sqlapp/io.table.go.tmpl.txt -> '%s'\n", name)

	if name2, err = createModelPath("io.table.go.tmpl.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\tio.table.go.tmpl.txt -> '%s'\n", name2)

    if name != name2 {
		t.Errorf("createModelPath() file names don't match!\n")
	}

	if name, err = createModelPath("util"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t../Models/util -> '%s'\n", name)

	t.Log("...End of TestCreateModelPath")
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
    sharedData.SetDataPath("./test/db.json.txt")
    sharedData.SetMainPath("./test/main.json.txt")

	if name, err = createOutputPath("dir", "dn", "tn", "fn"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t -> '%s'\n", name)
    if name != "/tmp/gen/dir/fn" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	t.Log("...End of TestCreateOutputPath")
}


