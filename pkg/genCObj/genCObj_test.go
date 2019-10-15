// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test C Object Generator

package genCObj

import (
	"genapp/pkg/sharedData"
	"testing"
	"time"
	//"github.com/2kranki/go_util"
)

const jsonTestPath = "../../misc/"

// Setup Shared Data for the current test. We assume Debug and Noop.
// So, things don't change in the real environment.
func setupShared(t *testing.T) {

	t.Log("genCObj::SetupShared()")

	sharedData.SetDebug(true)
	sharedData.SetForce(false)
	sharedData.SetNoop(true)
	sharedData.SetQuiet(false)
	sharedData.SetMdlDir("../../models/")
	sharedData.SetOutDir("/tmp/testgen")
	sharedData.SetTime(time.Now().Format("Mon Jan _2, 2006 15:04"))
	sharedData.SetFunc("Time", sharedData.Time)
	sharedData.SetDataPath("./test/db.json.txt")
	sharedData.SetMainPath("./test/main.json.txt")

	t.Log("...End of genCObj::SetupShared")
}

// Setup GenData and sharedData structures for testing.
func setupForCObj(t *testing.T, genData *GenData) {

	t.Log("genCObj::TestCreateModelPath()")

	genData.Name = "cobj"
	genData.Mapper = func(varSub string) string {
		switch varSub {
		case "Name":
			//return dbStruct.Name
			return "???"
		}
		return ""
	}
	//genData.FileDefs1 = &FileDefs1
	//genData.CreateOutputDirs = CreateOutputDirs
	//genData.ReadJsonData = ReadJsonFileData
	//genData.SetupFile = SetupFile
	//genData.TmplData.Data = DbStruct()

	t.Log("...End of genCObj::TestCreateModelPath")
}

func TestReadJsonFiles(t *testing.T) {
	var data *appData.Database
	var err error

	t.Log("genCObj::TestReadJsonFiles()")
	setupShared(t)

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

	t.Log("Successfully, completed: genCObj::TestReadJsonFiles")
}

func TestCreateModelPath(t *testing.T) {
	var name string
	var name2 string
	var err error

	t.Log("genCObj::TestCreateModelPath()")
	setupShared(t)

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

	t.Log("Successfully, completed: genCObj::TestCreateModelPath")
}

func TestCreateOutputPath(t *testing.T) {
	var name string
	var err error

	t.Log("genCObj::TestCreateOutputPath()")
	setupShared(t)

	if name, err = createOutputPath("tableio.go"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\ttableio.go -> '%s'\n", name)
	if name != "/tmp/gen/tableio.go" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	t.Log("Successfully, completed: genCObj::TestCreateOutputPath")
}
