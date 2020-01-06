// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package genSqlAppGo

import (
	"genapp/pkg/sharedData"
	"github.com/2kranki/go_util"
	"testing"
	"time"
)

func TestCreateOutputPath(t *testing.T) {
	var name *util.Path
	var err error
	var dirs []string = []string{"dir"}

	t.Log("TestCreateOutputPath()")
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

	if name, err = CreateOutputPath(dirs, "dn", "tn", "fn"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn -> '%s'\n", name)
	if name.String() != "/tmp/gen/dir/fn" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	if name, err = CreateOutputPath(dirs, "dn", "tn", "fn_${DbName}.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn_${DbName}.txt -> '%s'\n", name)
	if name.String() != "/tmp/gen/dir/fn_Dn.txt" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	if name, err = CreateOutputPath(dirs, "dn", "tn", "fn_${TblName}.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn_${TblName}.txt -> '%s'\n", name)
	if name.String() != "/tmp/gen/dir/fn_Tn.txt" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	dirs = []string{"dir1", "dir2"}
	if name, err = CreateOutputPath(dirs, "dn", "tn", "fn_${TblName}.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t [dir1,dir2]/fn_${TblName}.txt -> '%s'\n", name)
	if name.String() != "/tmp/gen/dir1/dir2/fn_Tn.txt" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	t.Log("...End of TestCreateOutputPath")
}

