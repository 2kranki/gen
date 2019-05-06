// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test C Object Generator

package genCObj

import (
	"testing"
	"time"
	"../appData"
	"../shared"
)

const jsonTestPath = "../misc/"

func TestReadJsonFiles(t *testing.T) {
	var data        *appData.Database
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


