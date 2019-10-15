// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package genSqlAppGo

import (
	"genapp/pkg/genSqlAppGo/dbJson"
	"genapp/pkg/sharedData"
	"github.com/2kranki/go_util"
	"log"
	"os"
	"testing"
	"time"
)

func TestReadJsonFiles(t *testing.T) {
	var data *dbJson.Database
	var err error

	t.Log("TestReadJsonFiles()")
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

	if data.Name != "Finances" {
		t.Errorf("ReadJsonFile() failed: Name should be 'Finances' but is '%s'\n", data.Name)
	}
	if len(data.Tables) != 2 {
		t.Errorf("ReadJsonFile() failed: should be 1 tables but is %d\n", len(data.Tables))
	}
	if len(data.Tables[0].Fields) != 8 {
		t.Errorf("ReadJsonFile() failed: should be 8 fields in table[0] but is %d\n", len(data.Tables[0].Fields))
	}

	t.Log("...End of TestReadJsonFiles")
}

func TestCreateModelPath(t *testing.T) {
	var err error
	var name util.Path
	var name2 util.Path

	t.Log("TestCreateModelPath()")
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

	if name.String() != name2.String() {
		t.Errorf("createModelPath() file names don't match!\n")
	}

	if name, err = createModelPath("util"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t../Models/util -> '%s'\n", name)

	t.Log("...End of TestCreateModelPath")
}

func TestCreateOutputPath(t *testing.T) {
	var name util.Path
	var err error

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

	if name, err = createOutputPath("dir", "dn", "tn", "fn"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn -> '%s'\n", name)
	if name.String() != "/tmp/gen/dir/fn" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	if name, err = createOutputPath("dir", "dn", "tn", "fn_${DbName}.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn_${DbName}.txt -> '%s'\n", name)
	if name.String() != "/tmp/gen/dir/fn_Dn.txt" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	if name, err = createOutputPath("dir", "dn", "tn", "fn_${TblName}.txt"); err != nil {
		t.Errorf("createModelPath() failed: %s\n", err)
	}
	t.Logf("\t dir/fn_${TblName}.txt -> '%s'\n", name)
	if name.String() != "/tmp/gen/dir/fn_Tn.txt" {
		t.Errorf("createOutputPath() file path isn't correct!\n")
	}

	t.Log("...End of TestCreateOutputPath")
}

func TestCreateDirs(t *testing.T) {
	var name string
	var err error
	var data *dbJson.Database
	var fi os.FileInfo

	t.Log("TestCreateDirs()")
	sharedData.SetDebug(true)
	sharedData.SetForce(false)
	sharedData.SetQuiet(false)
	sharedData.SetMdlDir("../models/")
	sharedData.SetOutDir("/tmp/appTest")
	sharedData.SetTime(time.Now().Format("Mon Jan _2, 2006 15:04"))
	sharedData.SetFunc("Time", sharedData.Time)
	sharedData.SetDataPath("./test/db.json.txt")
	sharedData.SetMainPath("./test/main.json.txt")

	if err = readJsonFiles(); err != nil {
		t.Fatalf("ReadJsonFile() failed: %s\n", err.Error())
	}
	data = dbJson.DbStruct()

	if data.Name != "Finances" {
		t.Fatalf("ReadJsonFile() failed: Name should be 'Finances' but is '%s'\n", data.Name)
	}
	if len(data.Tables) != 2 {
		t.Fatalf("ReadJsonFile() failed: should be 2 tables but is %d\n", len(data.Tables))
	}

	t.Logf("\tRemoving existing %s if present...\n", sharedData.OutDir())
	err = os.RemoveAll(sharedData.OutDir())
	if err != nil {
		t.Fatalf("Error: RemoveAll failed: %s\n", err.Error())
	}

	//t.Fatalf("Temporary FATAL\n")

	sharedData.SetNoop(true)
	t.Logf("\tCreating directories (NOOP == true)...\n")
	if err = createOutputDirs(data.TitledName(), data.Tables); err != nil {
		t.Fatalf("createOutputDirs() failed: %s\n", err)
	}

	//t.Fatalf("Temporary FATAL\n")

	name = "/tmp/appTest"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/html"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/static"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/style"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/tmpl"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src/hndlrFinances"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src/ioFinances"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src/FinancesCustomer"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src/FinancesVendor"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src/hndlrFinancesCustomer"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src/hndlrFinancesVendor"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src/ioFinancesCustomer"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}
	name = "/tmp/appTest/src/ioFinancesVendor"
	fi, err = os.Lstat(name)
	if err == nil {
		t.Fatalf("%s exists!\n", name)
	}

	//t.Fatalf("Temporary FATAL\n")

	t.Logf("\tRemoving existing %s if present...\n", sharedData.OutDir())
	err = os.RemoveAll(sharedData.OutDir())
	if err != nil {
		t.Fatalf("Error: RemoveAll failed: %s\n", err.Error())
	}

	sharedData.SetNoop(false)
	t.Logf("\tCreating directories (NOOP == false)...\n")
	if err = createOutputDirs(data.TitledName(), data.Tables); err != nil {
		t.Fatalf("createOutputDirs() failed: %s\n", err)
	}

	name = "/tmp/appTest"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/html"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/static"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/style"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/tmpl"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src/hndlrFinances"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src/ioFinances"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src/FinancesCustomer"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src/FinancesVendor"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src/hndlrFinancesCustomer"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src/hndlrFinancesVendor"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src/ioFinancesCustomer"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}
	name = "/tmp/appTest/src/ioFinancesVendor"
	fi, err = os.Lstat(name)
	if err != nil || !fi.Mode().IsDir() {
		t.Fatalf("%s does not exist!\n", name)
	}

	//os.RemoveAll(sharedData.OutDir())

	t.Log("...End of TestCreateDirs")
}

func TestCopyDocker(t *testing.T) {
	var name string
	var err error
	var data *dbJson.Database
	var fi os.FileInfo
	var pathIn util.Path
	var FileDefns2 []FileDefn = []FileDefn{
		{"docker",
			"/src/",
			"",
			"copyDir",
			0644,
			"single",
			0,
		},
		{"util",
			"/src/",
			"",
			"copyDir",
			0644,
			"single",
			0,
		},
	}

	t.Log("TestCreateDirs()")
	sharedData.SetDebug(true)
	sharedData.SetForce(false)
	sharedData.SetQuiet(false)
	sharedData.SetMdlDir("../models/")
	sharedData.SetOutDir("/tmp/appTest")
	sharedData.SetTime(time.Now().Format("Mon Jan _2, 2006 15:04"))
	sharedData.SetFunc("Time", sharedData.Time)
	sharedData.SetDataPath("./test/db.json.txt")
	sharedData.SetMainPath("./test/main.json.txt")

	if err = readJsonFiles(); err != nil {
		t.Fatalf("ReadJsonFile() failed: %s\n", err.Error())
	}
	data = dbJson.DbStruct()

	if data.Name != "Finances" {
		t.Fatalf("ReadJsonFile() failed: Name should be 'Finances' but is '%s'\n", data.Name)
	}
	if len(data.Tables) != 2 {
		t.Fatalf("ReadJsonFile() failed: should be 2 tables but is %d\n", len(data.Tables))
	}

	t.Logf("\tRemoving existing %s if present...\n", sharedData.OutDir())
	err = os.RemoveAll(sharedData.OutDir())
	if err != nil {
		t.Fatalf("Error: RemoveAll failed: %s\n", err.Error())
	}

	sharedData.SetNoop(false)
	t.Logf("\tCreating directories (NOOP == false)...\n")
	if err = createOutputDirs(data.TitledName(), data.Tables); err != nil {
		t.Fatalf("createOutputDirs() failed: %s\n", err)
	}

	//os.RemoveAll(sharedData.OutDir())

	if !sharedData.Quiet() {
		log.Println("Process file:", FileDefns2[0].ModelName, "generating:", FileDefns2[0].FileName, "...")
	}

	// Create the input model file path.
	if pathIn, err = createModelPath(FileDefns2[0].ModelName); err != nil {
		t.Fatalf("Error: %s: %s\n", pathIn.String(), err.Error())
	}
	if sharedData.Debug() {
		log.Println("\t\tmodelPath=", pathIn)
	}

	// Standard File
	data := TaskData{FD: &FileDefns[i], TD: &tmplData, PathIn: pathIn}
	// Create the output path
	data.PathOut, err = createOutputPath(def.FileDir, tmplData.Data.Name, "", def.FileName)
	if err != nil {
		log.Fatalln(err)
	}
	if sharedData.Debug() {
		log.Println("\t\t outPath=", data.PathOut)
	}
	// Generate the file.
	data.genFile()
	t.Log("...End of TestCreateDirs")
}
