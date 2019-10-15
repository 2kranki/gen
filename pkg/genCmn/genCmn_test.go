// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package genCmn

import (
	// "genapp/pkg/genSqlAppGo/dbJson"
	"genapp/pkg/sharedData"

	"testing"
	"time"

	"github.com/2kranki/go_util"
)

// Setup Shared Data for the current test. We assume Debug and Noop.
// So, things don't change in the real environment.
func setupShared(t *testing.T) {

	t.Log("genCmn::SetupShared()")

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

	t.Log("...End of genCmn::SetupShared")
}

func TestCreateModelPath(t *testing.T) {
	var err error
	var name *util.Path
	var name2 *util.Path
	var gd *GenData

	t.Log("genCmn::TestCreateModelPath()")
	if gd = NewGenData(); gd == nil {
		t.Fatalf("Failed to create GenData!\n")
	}
	setupShared(t)
	gd.Name = "sqlapp"

	if name, err = gd.CreateModelPath("../sqlapp/io.table.go.tmpl.txt"); err != nil {
		t.Errorf("CreateModelPath() failed: %s\n", err)
	}
	t.Logf("\t../sqlapp/io.table.go.tmpl.txt -> '%s'\n", name)
	if name.String() != "../../models/sqlapp/io.table.go.tmpl.txt" {
		t.Errorf("CreateModelPath() file names don't match - %s!\n", name)
	}

	if name2, err = gd.CreateModelPath("io.table.go.tmpl.txt"); err != nil {
		t.Errorf("CreateModelPath() failed: %s\n", err)
	}
	t.Logf("\tio.table.go.tmpl.txt -> '%s'\n", name2)
	if name2.String() != "../../models/sqlapp/io.table.go.tmpl.txt" {
		t.Errorf("CreateModelPath() file names don't match - %s!\n", name)
	}

	if name, err = gd.CreateModelPath("dbs"); err != nil {
		t.Errorf("CreateModelPath() failed: %s\n", err)
	}
	t.Logf("\t../../Models/dbs -> '%s'\n", name)
	if name.String() != "../../models/sqlapp/dbs" {
		t.Errorf("CreateModelPath() file names don't match - %s!\n", name)
	}

	if name, err = gd.CreateModelPath("util"); err == nil {
		t.Errorf("CreateModelPath() failed: %s\n", err)
	}
	t.Logf("\t../../Models/util -> '%s'\n", name)
	if name != nil {
		t.Errorf("CreateModelPath() file names don't match - %s!\n", name)
	}

	t.Log("...End of genCmn::TestCreateModelPath")
}

func TestCreateOutputPath(t *testing.T) {
	var name *util.Path
	var err error
	var gd *GenData
	var fn string = "file.txt"

	t.Log("genCmn::TestCreateOutputPath()")
	if gd = NewGenData(); gd == nil {
		t.Fatalf("Failed to create GenData!\n")
	}
	setupShared(t)
	gd.Name = "sqlapp"
	mapper := func(placeholderName string) string {
		switch placeholderName {
		case "DbName":
			return "Movies"
		case "TblName":
			return "Movie"
		}
		return ""
	}

	if name, err = gd.CreateOutputPath(mapper, nil, fn); err != nil {
		t.Errorf("CreateModelPath() failed: %s\n", err)
	}
	t.Logf("\t nil/fn -> '%s'\n", name)
	if name.String() != ("/tmp/testgen/" + fn) {
		t.Errorf("CreateOutputPath() file path isn't correct!\n")
	}

	if name, err = gd.CreateOutputPath(mapper, nil, "fn_${DbName}.txt"); err != nil {
		t.Errorf("CreateModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn_${DbName}.txt -> '%s'\n", name)
	if name.String() != "/tmp/testgen/fn_Movies.txt" {
		t.Errorf("CreateOutputPath() file path isn't correct!\n")
	}

	if name, err = gd.CreateOutputPath(mapper, nil, "fn_${TblName}.txt"); err != nil {
		t.Errorf("CreateModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn_${TblName}.txt -> '%s'\n", name)
	if name.String() != "/tmp/testgen/fn_Movie.txt" {
		t.Errorf("CreateOutputPath() file path isn't correct!\n")
	}

	if name, err = gd.CreateOutputPath(mapper, nil, "fn_${DbName}${TblName}.txt"); err != nil {
		t.Errorf("CreateModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn_${DbName}${TblName}.txt -> '%s'\n", name)
	if name.String() != "/tmp/testgen/fn_MoviesMovie.txt" {
		t.Errorf("CreateOutputPath() file path isn't correct!\n")
	}

	t.Log("...End of genCmn::TestCreateOutputPath")
}
